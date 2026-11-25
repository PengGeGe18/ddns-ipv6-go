// Package cert 负责通过 certbot 为当前 IPv6 地址申请/更新证书，
// 部署证书到固定目录并重启 Nginx 使其生效。
package cert

import (
	"context"
	"crypto/x509"
	"encoding/hex"
	"encoding/pem"
	"fmt"
	"net"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
	"sync"
	"time"

	"github.com/penggege18/ddns-ipv6-go/v6/config"
	"github.com/penggege18/ddns-ipv6-go/v6/util"
)

const (
	// 单次 certbot 命令超时时间
	commandTimeout = 5 * time.Minute
	// Nginx 命令超时时间
	reloadTimeout = time.Minute
	// 自动补申请证书的最小间隔
	ensureCooldown = 30 * time.Minute
	// 结果消息最大长度
	maxMessageLength = 2000
)

// Status 证书任务状态
type Status struct {
	Running       bool      `json:"running"`
	LastRunAt     time.Time `json:"lastRunAt"`
	LastResult    string    `json:"lastResult"`
	LastMessage   string    `json:"lastMessage"`
	LastSuccessAt time.Time `json:"lastSuccessAt"`
	CertName      string    `json:"certName"`
	CertIPs       []string  `json:"certIPs"`
	NotAfter      time.Time `json:"notAfter"`
	NextRenewAt   time.Time `json:"nextRenewAt"`

	// 上次尝试申请证书时间(不通过接口暴露)
	lastAttemptAt time.Time
}

var (
	mu    sync.Mutex
	state = Status{}
)

// 证书任务结果
const (
	ResultSuccess = "成功"
	ResultFailed  = "失败"
	ResultSkip    = "跳过"
)

// GetStatus 获取当前证书任务状态
func GetStatus() Status {
	mu.Lock()
	defer mu.Unlock()
	s := state
	s.CertIPs = append([]string(nil), state.CertIPs...)
	return s
}

// begin 标记任务开始，已有任务在运行时返回 false
func begin() bool {
	mu.Lock()
	defer mu.Unlock()
	if state.Running {
		return false
	}
	state.Running = true
	state.LastRunAt = time.Now()
	return true
}

// finish 标记任务结束
func finish(result, message string) {
	mu.Lock()
	defer mu.Unlock()
	state.Running = false
	state.LastResult = result
	state.LastMessage = truncate(message, maxMessageLength)
	if result == ResultSuccess {
		state.LastSuccessAt = time.Now()
	}
}

// CertName 根据 IPv6 地址生成稳定的证书名
func CertName(ip string) string {
	parsed := net.ParseIP(ip)
	if parsed == nil {
		return ""
	}
	return "ip6-" + hex.EncodeToString(parsed.To16())
}

// Cover 检查当前已部署的证书是否覆盖指定 IPv6 地址且未过期
func Cover(ip string) bool {
	conf, err := config.GetConfigCached()
	if err != nil || !conf.Cert.Enable {
		// 未启用证书管理时不参与判断
		return true
	}
	ips, notAfter, err := loadCertInfo(conf.Cert.LinkDir)
	if err != nil {
		return false
	}
	if !notAfter.After(time.Now()) {
		return false
	}
	for _, v := range ips {
		if v == ip {
			return true
		}
	}
	return false
}

// EnsureAsync 当证书未覆盖当前 IPv6 地址时自动补申请（带冷却时间，避免频繁请求）
func EnsureAsync(ip string) {
	if Cover(ip) {
		return
	}
	mu.Lock()
	if state.Running || time.Since(state.lastAttemptAt) < ensureCooldown {
		mu.Unlock()
		return
	}
	state.lastAttemptAt = time.Now()
	mu.Unlock()

	IssueAsync(ip, "")
}

// IssueAsync 为指定 IPv6 地址申请证书, 任务进行中时返回 false
func IssueAsync(ip, oldIP string) bool {
	if !begin() {
		util.Log("证书任务正在运行, 跳过本次申请")
		return false
	}
	go func() {
		result, message := issue(ip, oldIP)
		finish(result, message)
	}()
	return true
}

// RenewAsync 异步更新证书, 任务进行中时返回错误
func RenewAsync(force bool) error {
	if !begin() {
		return fmt.Errorf("证书任务正在运行, 请稍后再试")
	}
	go func() {
		result, message := renew(force)
		finish(result, message)
	}()
	return nil
}

func issue(ip, oldIP string) (string, string) {
	conf, err := config.GetConfigCached()
	if err != nil {
		util.Log("读取配置失败, 无法申请证书: %s", err)
		return ResultFailed, "读取配置失败: " + err.Error()
	}
	if !conf.Cert.Enable {
		return ResultSkip, "证书自动申请未启用"
	}

	parsed := net.ParseIP(ip)
	if parsed == nil || parsed.To4() != nil {
		util.Log("无效的 IPv6 地址: %s", ip)
		return ResultFailed, "无效的 IPv6 地址: " + ip
	}

	name := CertName(ip)
	util.Log("开始为 IPv6 地址 %s 申请证书, cert-name: %s", ip, name)

	ctx, cancel := context.WithTimeout(context.Background(), commandTimeout)
	defer cancel()

	out, err := runCommand(ctx, conf.Cert.CertbotPath, buildIssueArgs(&conf, ip)...)
	logOutput("certbot", out)
	if err != nil {
		msg := fmt.Sprintf("证书申请失败: %s\n%s", err, truncate(out, 800))
		util.Log("证书申请失败: %s", err)
		return ResultFailed, msg
	}

	if err := deploy(&conf, ip); err != nil {
		util.Log("证书部署失败: %s", err)
		return ResultFailed, "证书部署失败: " + err.Error()
	}

	if err := reloadNginx(&conf); err != nil {
		util.Log("Nginx 重启失败: %s", err)
		return ResultFailed, "证书申请成功, 但 Nginx 重启失败: " + err.Error()
	}

	// 删除旧 IP 的证书
	if oldIP != "" && oldIP != ip && conf.Cert.DeleteOld {
		deleteLineage(&conf, oldIP)
	}

	refreshCertInfo(&conf)
	mu.Lock()
	state.CertName = name
	mu.Unlock()
	util.Log("IPv6 地址 %s 的证书申请并部署成功, Nginx 已重启", ip)
	return ResultSuccess, "证书申请成功, 已部署并重启 Nginx"
}

func buildIssueArgs(conf *config.Config, ip string) []string {
	args := []string{
		"certonly",
		"--non-interactive",
		"--agree-tos",
		"--cert-name", CertName(ip),
		"--ip-address", ip,
	}

	if email := strings.TrimSpace(conf.Cert.Email); email != "" {
		args = append(args, "-m", email)
	} else {
		args = append(args, "--register-unsafely-without-email")
	}

	if profile := strings.TrimSpace(conf.Cert.Profile); profile != "" {
		args = append(args, "--preferred-profile", profile)
	}

	if conf.Cert.Challenge == "webroot" {
		args = append(args, "--webroot", "-w", conf.Cert.Webroot)
	} else {
		args = append(args, "--standalone")
	}

	if extra := strings.Fields(conf.Cert.ExtraArgs); len(extra) > 0 {
		args = append(args, extra...)
	}
	return args
}

func renew(force bool) (string, string) {
	conf, err := config.GetConfigCached()
	if err != nil {
		return ResultFailed, "读取配置失败: " + err.Error()
	}
	if !conf.Cert.Enable {
		return ResultSkip, "证书自动更新未启用"
	}

	name := currentCertName()
	if name == "" {
		return ResultSkip, "未找到当前证书, 无需更新"
	}

	args := []string{"renew", "--non-interactive", "--cert-name", name}
	if force {
		args = append(args, "--force-renewal")
	}

	util.Log("开始更新证书(force=%v), cert-name: %s", force, name)
	ctx, cancel := context.WithTimeout(context.Background(), commandTimeout)
	defer cancel()

	out, err := runCommand(ctx, conf.Cert.CertbotPath, args...)
	logOutput("certbot", out)
	if err != nil {
		msg := fmt.Sprintf("证书更新失败: %s\n%s", err, truncate(out, 800))
		util.Log("证书更新失败: %s", err)
		return ResultFailed, msg
	}

	if err := reloadNginx(&conf); err != nil {
		util.Log("Nginx 重启失败: %s", err)
		return ResultFailed, "证书更新完成, 但 Nginx 重启失败: " + err.Error()
	}

	refreshCertInfo(&conf)
	util.Log("证书更新完成, Nginx 已重启")
	return ResultSuccess, "证书更新完成, 已重启 Nginx"
}

// currentCertName 获取当前证书名
func currentCertName() string {
	mu.Lock()
	defer mu.Unlock()
	if state.CertName != "" {
		return state.CertName
	}
	if len(state.CertIPs) > 0 {
		return CertName(state.CertIPs[0])
	}
	return ""
}

// deploy 将 certbot 生成的证书通过软链接(或复制)部署到固定目录
func deploy(conf *config.Config, ip string) error {
	liveDir := filepath.Join(config.LetsEncryptDir, "live", CertName(ip))
	if _, err := os.Stat(filepath.Join(liveDir, "fullchain.pem")); err != nil {
		return fmt.Errorf("未找到证书文件: %s", err)
	}

	if err := os.MkdirAll(conf.Cert.LinkDir, 0755); err != nil {
		return err
	}

	for _, name := range []string{"fullchain.pem", "privkey.pem", "chain.pem"} {
		src := filepath.Join(liveDir, name)
		if _, err := os.Stat(src); err != nil {
			// chain.pem 部分场景可能不存在, 跳过
			if name == "chain.pem" {
				continue
			}
			return err
		}
		if err := replaceLink(src, filepath.Join(conf.Cert.LinkDir, name)); err != nil {
			return err
		}
	}
	return nil
}

// replaceLink 优先创建软链接, 失败时退化为复制文件
func replaceLink(src, dst string) error {
	tmp := dst + ".tmp"
	_ = os.Remove(tmp)

	if err := os.Symlink(src, tmp); err == nil {
		_ = os.Remove(dst)
		if err := os.Rename(tmp, dst); err == nil {
			return nil
		}
		_ = os.Remove(tmp)
	}

	// 复制文件（Windows 或无权限创建软链接时）
	byt, err := os.ReadFile(src)
	if err != nil {
		return err
	}
	mode := os.FileMode(0644)
	if strings.HasSuffix(dst, "privkey.pem") {
		mode = 0600
	}
	return os.WriteFile(dst, byt, mode)
}

// reloadNginx 执行配置的 Nginx 重启命令
func reloadNginx(conf *config.Config) error {
	cmdLine := strings.TrimSpace(conf.Cert.NginxReloadCmd)
	if cmdLine == "" {
		return nil
	}

	ctx, cancel := context.WithTimeout(context.Background(), reloadTimeout)
	defer cancel()

	var out string
	var err error
	if runtime.GOOS == "windows" {
		out, err = runCommand(ctx, "cmd", "/c", cmdLine)
	} else {
		out, err = runCommand(ctx, "sh", "-c", cmdLine)
	}
	logOutput("nginx", out)
	if err != nil {
		return fmt.Errorf("%s: %s", err, truncate(out, 200))
	}
	return nil
}

// deleteLineage 删除旧 IP 对应的证书
func deleteLineage(conf *config.Config, oldIP string) {
	name := CertName(oldIP)
	if name == "" {
		return
	}
	ctx, cancel := context.WithTimeout(context.Background(), commandTimeout)
	defer cancel()
	out, err := runCommand(ctx, conf.Cert.CertbotPath, "delete", "--cert-name", name, "--non-interactive")
	if err != nil {
		util.Log("删除旧证书 %s 失败: %s", name, err)
		return
	}
	logOutput("certbot", out)
	util.Log("已删除旧证书: %s", name)
}

// loadCertInfo 读取已部署证书中包含的 IP 地址及过期时间
func loadCertInfo(linkDir string) ([]string, time.Time, error) {
	byt, err := os.ReadFile(filepath.Join(linkDir, "fullchain.pem"))
	if err != nil {
		return nil, time.Time{}, err
	}
	block, _ := pem.Decode(byt)
	if block == nil {
		return nil, time.Time{}, fmt.Errorf("证书内容解析失败")
	}
	cert, err := x509.ParseCertificate(block.Bytes)
	if err != nil {
		return nil, time.Time{}, err
	}

	ips := make([]string, 0, len(cert.IPAddresses))
	for _, ip := range cert.IPAddresses {
		ips = append(ips, ip.String())
	}
	return ips, cert.NotAfter, nil
}

// refreshCertInfo 刷新内存中的证书信息
func refreshCertInfo(conf *config.Config) {
	ips, notAfter, err := loadCertInfo(conf.Cert.LinkDir)
	if err != nil {
		return
	}
	mu.Lock()
	state.CertIPs = ips
	state.NotAfter = notAfter
	if state.CertName == "" && len(ips) > 0 {
		state.CertName = CertName(ips[0])
	}
	mu.Unlock()
}

// StartScheduler 启动每日定时更新证书任务
func StartScheduler() {
	// 启动时读取已有证书信息
	if conf, err := config.GetConfigCached(); err == nil {
		refreshCertInfo(&conf)
	}

	go func() {
		for {
			conf, err := config.GetConfigCached()
			if err != nil || !conf.Cert.Enable || !conf.Cert.DailyRenew {
				mu.Lock()
				state.NextRenewAt = time.Time{}
				mu.Unlock()
				time.Sleep(time.Minute)
				continue
			}

			next := nextDaily(conf.Cert.RenewHour, conf.Cert.RenewMinute)
			mu.Lock()
			state.NextRenewAt = next
			mu.Unlock()

			util.Log("下次自动更新证书时间: %s", next.Format("2006-01-02 15:04:05"))
			time.Sleep(time.Until(next))

			if !begin() {
				util.Log("证书任务正在运行, 跳过本次定时更新")
				continue
			}
			result, message := renew(false)
			finish(result, message)
		}
	}()
}

// nextDaily 计算下一次执行时间
func nextDaily(hour, minute int) time.Time {
	now := time.Now()
	next := time.Date(now.Year(), now.Month(), now.Day(), hour, minute, 0, 0, now.Location())
	if !next.After(now) {
		next = next.Add(24 * time.Hour)
	}
	return next
}

func runCommand(ctx context.Context, name string, args ...string) (string, error) {
	cmd := exec.CommandContext(ctx, name, args...)
	out, err := cmd.CombinedOutput()
	return string(out), err
}

func logOutput(prefix, out string) {
	lines := strings.Split(strings.TrimSpace(out), "\n")
	start := 0
	if len(lines) > 30 {
		start = len(lines) - 30
	}
	for _, line := range lines[start:] {
		if s := strings.TrimSpace(line); s != "" {
			util.Log("[%s] %s", prefix, s)
		}
	}
}

func truncate(s string, n int) string {
	s = strings.TrimSpace(s)
	if len(s) <= n {
		return s
	}
	return s[:n] + "..."
}
