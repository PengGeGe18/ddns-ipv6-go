package config

import (
	"errors"
	"io"
	"log"
	"net"
	"os"
	"os/exec"
	"regexp"
	"runtime"
	"strconv"
	"strings"
	"sync"

	"github.com/penggege18/ddns-ipv6-go/v6/util"
	passwordvalidator "github.com/wagslane/go-password-validator"
	"gopkg.in/yaml.v3"
)

// Ipv4Reg IPv4正则
var Ipv4Reg = regexp.MustCompile(`((25[0-5]|(2[0-4]|1{0,1}[0-9]){0,1}[0-9])\.){3,3}(25[0-5]|(2[0-4]|1{0,1}[0-9]){0,1}[0-9])`)

// Ipv6Reg IPv6正则
// 用于从 URL / 命令输出中抓取 IPv6 字符串，最终是否合法由 net.ParseIP 保证
var Ipv6Reg = regexp.MustCompile(`([0-9A-Fa-f:.]{2,})`)

// 默认配置项
const (
	// DefaultInterval 默认 IP 检测间隔(秒)
	DefaultInterval = 60
	// DefaultCertbotPath 默认 certbot 可执行文件
	DefaultCertbotPath = "certbot"
	// DefaultLinkDir 默认证书部署目录（Nginx 应指向此目录）
	DefaultLinkDir = "/etc/ssl/ddns-ipv6"
	// DefaultNginxReloadCmd 默认 Nginx 重启命令(Linux)
	DefaultNginxReloadCmd = "systemctl restart nginx"
	// DefaultDarwinNginxReloadCmd macOS 下由 Homebrew 管理 Nginx
	DefaultDarwinNginxReloadCmd = "brew services restart nginx"
	// DefaultCertProfile 默认证书 Profile（Let's Encrypt IP 证书需要 shortlived）
	DefaultCertProfile = "shortlived"
	// LetsEncryptDir certbot 证书根目录
	LetsEncryptDir = "/etc/letsencrypt"
)

// IPConfig 单个 IP 版本的获取配置
type IPConfig struct {
	// Enable 是否启用
	Enable bool
	// GetType 获取IP类型 url/netInterface/cmd
	GetType      string
	URL          string
	NetInterface string
	Cmd          string
	// Ipv6Reg 仅 IPv6 使用: @N 表示第 N 个地址，或正则表达式
	Ipv6Reg string
}

// CertConfig certbot 证书相关配置
type CertConfig struct {
	// Enable 是否启用自动申请/更新证书
	Enable bool
	// CertbotPath certbot 可执行文件路径
	CertbotPath string
	// Email 注册邮箱（可留空）
	Email string
	// Challenge 验证方式 standalone/webroot
	Challenge string
	// Webroot webroot 验证时的网站根目录
	Webroot string
	// Profile 证书 profile，IP 证书需 shortlived
	Profile string
	// ExtraArgs 追加给 certbot 的参数
	ExtraArgs string
	// LinkDir 证书部署目录，Nginx 指向这里即可在 IP 变化后继续生效
	LinkDir string
	// Owner 证书所属用户(可选)，Nginx 以非 root 运行时填写该用户以便读取私钥
	Owner string
	// NginxReloadCmd 证书部署后执行的 Nginx 重启命令
	NginxReloadCmd string
	// DailyRenew 是否每天定时更新证书
	DailyRenew bool
	// RenewHour 每日更新小时(0-23)
	RenewHour int
	// RenewMinute 每日更新分钟(0-59)
	RenewMinute int
	// DeleteOld 申请新证书成功后删除旧 IP 的证书
	DeleteOld bool
}

// Config 配置
type Config struct {
	User
	Webhook
	IPv4 IPConfig
	IPv6 IPConfig
	// Interval IP 检测间隔(秒)
	Interval int
	Cert     CertConfig
	// 禁止公网访问
	NotAllowWanAccess bool
	// 语言
	Lang string
}

// ConfigCache ConfigCache
type cacheType struct {
	ConfigSingle *Config
	Err          error
	Lock         sync.Mutex
}

var cache = &cacheType{}

// GetConfigCached 获得缓存的配置
func GetConfigCached() (conf Config, err error) {
	cache.Lock.Lock()
	defer cache.Lock.Unlock()

	if cache.ConfigSingle != nil {
		return *cache.ConfigSingle, cache.Err
	}

	// init config
	cache.ConfigSingle = &Config{}

	configFilePath := util.GetConfigFilePath()
	_, err = os.Stat(configFilePath)
	if err != nil {
		// 首次运行, 使用默认配置
		cache.ConfigSingle.SetDefaults()
		cache.Err = nil
		return *cache.ConfigSingle, nil
	}

	byt, err := os.ReadFile(configFilePath)
	if err != nil {
		util.Log("异常信息: %s", err)
		cache.Err = err
		return *cache.ConfigSingle, err
	}

	err = yaml.Unmarshal(byt, cache.ConfigSingle)
	if err != nil {
		util.Log("异常信息: %s", err)
		cache.Err = err
		return *cache.ConfigSingle, err
	}

	// 填充默认值
	cache.ConfigSingle.Normalize()

	// 未填写登录信息, 确保不能从公网访问
	if cache.ConfigSingle.Username == "" && cache.ConfigSingle.Password == "" {
		cache.ConfigSingle.NotAllowWanAccess = true
	}

	// remove err
	cache.Err = nil
	return *cache.ConfigSingle, err
}

// SetDefaults 首次运行时的默认配置
func (conf *Config) SetDefaults() {
	conf.Interval = DefaultInterval
	conf.IPv4.Enable = true
	conf.IPv4.GetType = "netInterface"
	conf.IPv6.Enable = true
	conf.IPv6.GetType = "netInterface"
	conf.Cert.Enable = true
	conf.Cert.DailyRenew = true
	conf.Cert.DeleteOld = true
	conf.Cert.RenewHour = 3
	conf.Cert.RenewMinute = 0
	conf.normalizeCert()
}

// Normalize 填充缺省值并修正非法配置
func (conf *Config) Normalize() {
	if conf.Interval <= 0 {
		conf.Interval = DefaultInterval
	}
	if conf.Interval < 5 {
		conf.Interval = 5
	}
	if conf.Interval > 86400 {
		conf.Interval = 86400
	}
	if conf.IPv4.GetType == "" {
		conf.IPv4.GetType = "netInterface"
	}
	if conf.IPv6.GetType == "" {
		conf.IPv6.GetType = "netInterface"
	}
	conf.normalizeCert()
}

func (conf *Config) normalizeCert() {
	if strings.TrimSpace(conf.Cert.CertbotPath) == "" {
		conf.Cert.CertbotPath = DefaultCertbotPath
	}
	if strings.TrimSpace(conf.Cert.LinkDir) == "" {
		conf.Cert.LinkDir = defaultLinkDir()
	}
	if strings.TrimSpace(conf.Cert.NginxReloadCmd) == "" {
		conf.Cert.NginxReloadCmd = defaultNginxReloadCmd()
	}
	conf.Cert.Owner = strings.TrimSpace(conf.Cert.Owner)
	if strings.TrimSpace(conf.Cert.Profile) == "" {
		conf.Cert.Profile = DefaultCertProfile
	}
	conf.Cert.Challenge = strings.TrimSpace(conf.Cert.Challenge)
	if conf.Cert.Challenge != "webroot" {
		conf.Cert.Challenge = "standalone"
	}
	if conf.Cert.Challenge == "webroot" && strings.TrimSpace(conf.Cert.Webroot) == "" {
		conf.Cert.Webroot = "/var/www/html"
	}
	if conf.Cert.RenewHour < 0 || conf.Cert.RenewHour > 23 {
		conf.Cert.RenewHour = 3
	}
	if conf.Cert.RenewMinute < 0 || conf.Cert.RenewMinute > 59 {
		conf.Cert.RenewMinute = 0
	}
}

// defaultLinkDir 按平台推导默认证书目录: macOS 使用 Homebrew Nginx 的 certs 目录
func defaultLinkDir() string {
	if runtime.GOOS != "darwin" {
		return DefaultLinkDir
	}
	if _, err := os.Stat("/opt/homebrew"); err == nil {
		return "/opt/homebrew/etc/nginx/certs"
	}
	return "/usr/local/etc/nginx/certs"
}

// defaultNginxReloadCmd 按平台推导默认 Nginx 重启命令
func defaultNginxReloadCmd() string {
	if runtime.GOOS == "darwin" {
		return DefaultDarwinNginxReloadCmd
	}
	return DefaultNginxReloadCmd
}

// SaveConfig 保存配置
func (conf *Config) SaveConfig() (err error) {
	cache.Lock.Lock()
	defer cache.Lock.Unlock()

	byt, err := yaml.Marshal(conf)
	if err != nil {
		log.Println(err)
		return err
	}

	configFilePath := util.GetConfigFilePath()
	err = os.WriteFile(configFilePath, byt, 0600)
	if err != nil {
		log.Println(err)
		return
	}

	util.Log("配置文件已保存在: %s", configFilePath)

	// 清空配置缓存
	cache.ConfigSingle = nil

	return
}

// ResetPassword 重置密码
func (conf *Config) ResetPassword(newPassword string) {
	// 初始化语言
	util.InitLogLang(conf.Lang)

	// 先检查密码是否安全
	hashedPwd, err := conf.CheckPassword(newPassword)
	if err != nil {
		util.Log(err.Error())
		return
	}

	// 保存配置
	conf.Password = hashedPwd
	conf.SaveConfig()
	util.Log("用户名 %s 的密码已重置成功! 请重启服务", conf.Username)
}

// CheckPassword 检查密码
func (conf *Config) CheckPassword(newPassword string) (hashedPwd string, err error) {
	var minEntropyBits float64 = 30
	if conf.NotAllowWanAccess {
		minEntropyBits = 25
	}
	err = passwordvalidator.Validate(newPassword, minEntropyBits)
	if err != nil {
		return "", errors.New(util.LogStr("密码不安全！尝试使用更复杂的密码"))
	}

	// 加密密码
	hashedPwd, err = util.HashPassword(newPassword)
	if err != nil {
		return "", errors.New(util.LogStr("异常信息: %s", err.Error()))
	}
	return
}

func (conf *Config) getIpv4AddrFromInterface() string {
	ipv4, _, err := GetNetInterface()
	if err != nil {
		util.Log("从网卡获得IPv4失败")
		return ""
	}

	name := conf.IPv4.NetInterface
	if name == "" && len(ipv4) > 0 {
		// 未指定网卡时自动使用第一个有地址的网卡
		name = ipv4[0].Name
	}

	for _, netInterface := range ipv4 {
		if netInterface.Name == name && len(netInterface.Address) > 0 {
			return netInterface.Address[0]
		}
	}

	util.Log("从网卡中获得IPv4失败! 网卡名: %s", name)
	return ""
}

func (conf *Config) getIpv4AddrFromUrl() string {
	client := util.CreateBoundNoProxyHTTPClient("tcp4", "")
	urls := strings.Split(conf.IPv4.URL, ",")
	for _, url := range urls {
		url = strings.TrimSpace(url)
		if url == "" {
			continue
		}
		resp, err := client.Get(url)
		if err != nil {
			util.Log("通过接口获取IPv4失败! 接口地址: %s", url)
			util.Log("异常信息: %s", err)
			continue
		}
		defer resp.Body.Close()
		lr := io.LimitReader(resp.Body, 1024000)
		body, err := io.ReadAll(lr)
		if err != nil {
			util.Log("异常信息: %s", err)
			continue
		}
		result := Ipv4Reg.FindString(string(body))
		if result == "" {
			util.Log("获取IPv4结果失败! 接口: %s ,返回值: %s", url, string(body))
		}
		return result
	}
	return ""
}

func findIPv6InText(text string) string {
	for _, candidate := range Ipv6Reg.FindAllString(text, -1) {
		if ip := net.ParseIP(candidate); ip != nil && ip.To4() == nil {
			return candidate
		}
	}
	return ""
}

func (conf *Config) getAddrFromCmd(addrType string) string {
	var cmd string
	var comp *regexp.Regexp
	if addrType == "IPv4" {
		cmd = conf.IPv4.Cmd
		comp = Ipv4Reg
	} else {
		cmd = conf.IPv6.Cmd
		comp = Ipv6Reg
	}
	// cmd is empty
	if cmd == "" {
		return ""
	}
	// run cmd with proper shell
	var execCmd *exec.Cmd
	if runtime.GOOS == "windows" {
		execCmd = exec.Command("powershell", "-Command", cmd)
	} else {
		// If Bash does not exist, use sh
		_, err := exec.LookPath("bash")
		if err != nil {
			execCmd = exec.Command("sh", "-c", cmd)
		} else {
			execCmd = exec.Command("bash", "-c", cmd)
		}
	}
	// run cmd
	out, err := execCmd.CombinedOutput()
	if err != nil {
		util.Log("获取%s结果失败! 未能成功执行命令：%s, 错误：%q, 退出状态码：%s", addrType, execCmd.String(), out, err)
		return ""
	}
	str := string(out)
	// get result
	result := ""
	if addrType == "IPv4" {
		result = comp.FindString(str)
	} else {
		result = findIPv6InText(str)
	}
	if result == "" {
		util.Log("获取%s结果失败! 命令: %s, 标准输出: %q", addrType, execCmd.String(), str)
	}
	return result
}

// GetIpv4Addr 获得IPv4地址
func (conf *Config) GetIpv4Addr() string {
	// 判断从哪里获取IP
	switch conf.IPv4.GetType {
	case "netInterface":
		// 从网卡获取 IP
		return conf.getIpv4AddrFromInterface()
	case "url":
		// 从 URL 获取 IP
		return conf.getIpv4AddrFromUrl()
	case "cmd":
		// 从命令行获取 IP
		return conf.getAddrFromCmd("IPv4")
	default:
		log.Println("IPv4's get IP method is unknown")
		return "" // unknown type
	}
}

func (conf *Config) getIpv6AddrFromInterface() string {
	_, ipv6, err := GetNetInterface()
	if err != nil {
		util.Log("从网卡获得IPv6失败")
		return ""
	}

	name := conf.IPv6.NetInterface
	if name == "" && len(ipv6) > 0 {
		// 未指定网卡时自动使用第一个有地址的网卡
		name = ipv6[0].Name
	}

	for _, netInterface := range ipv6 {
		if netInterface.Name == name && len(netInterface.Address) > 0 {
			if conf.IPv6.Ipv6Reg != "" {
				// 匹配第几个IPv6
				if match, err := regexp.MatchString("@\\d", conf.IPv6.Ipv6Reg); err == nil && match {
					num, err := strconv.Atoi(conf.IPv6.Ipv6Reg[1:])
					if err == nil {
						if num > 0 {
							if num <= len(netInterface.Address) {
								return netInterface.Address[num-1]
							}
							util.Log("未找到第 %d 个IPv6地址! 将使用第一个IPv6地址", num)
							return netInterface.Address[0]
						}
						util.Log("IPv6匹配表达式 %s 不正确! 最小从1开始", conf.IPv6.Ipv6Reg)
						return ""
					}
				}
				// 正则表达式匹配
				util.Log("IPv6将使用正则表达式 %s 进行匹配", conf.IPv6.Ipv6Reg)
				for i := 0; i < len(netInterface.Address); i++ {
					matched, err := regexp.MatchString(conf.IPv6.Ipv6Reg, netInterface.Address[i])
					if matched && err == nil {
						util.Log("匹配成功! 匹配到地址: %s", netInterface.Address[i])
						return netInterface.Address[i]
					}
				}
				util.Log("没有匹配到任何一个IPv6地址, 将使用第一个地址")
			}
			return netInterface.Address[0]
		}
	}

	util.Log("从网卡中获得IPv6失败! 网卡名: %s", name)
	return ""
}

func (conf *Config) getIpv6AddrFromUrl() string {
	client := util.CreateBoundNoProxyHTTPClient("tcp6", "")
	urls := strings.Split(conf.IPv6.URL, ",")
	for _, url := range urls {
		url = strings.TrimSpace(url)
		if url == "" {
			continue
		}
		resp, err := client.Get(url)
		if err != nil {
			util.Log("通过接口获取IPv6失败! 接口地址: %s", url)
			util.Log("异常信息: %s", err)
			continue
		}

		defer resp.Body.Close()
		lr := io.LimitReader(resp.Body, 1024000)
		body, err := io.ReadAll(lr)
		if err != nil {
			util.Log("异常信息: %s", err)
			continue
		}
		result := findIPv6InText(string(body))
		if result == "" {
			util.Log("获取IPv6结果失败! 接口: %s ,返回值: %s", url, string(body))
		}
		return result
	}
	return ""
}

// GetIpv6Addr 获得IPv6地址
func (conf *Config) GetIpv6Addr() (result string) {
	// 判断从哪里获取IP
	switch conf.IPv6.GetType {
	case "netInterface":
		// 从网卡获取 IP
		return conf.getIpv6AddrFromInterface()
	case "url":
		// 从 URL 获取 IP
		return conf.getIpv6AddrFromUrl()
	case "cmd":
		// 从命令行获取 IP
		return conf.getAddrFromCmd("IPv6")
	default:
		log.Println("IPv6's get IP method is unknown")
		return "" // unknown type
	}
}
