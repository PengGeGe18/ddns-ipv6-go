package main

import (
	"embed"
	"errors"
	"flag"
	"fmt"
	"log"
	"net"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strconv"
	"time"

	"github.com/penggege18/ddns-ipv6-go/v6/cert"
	"github.com/penggege18/ddns-ipv6-go/v6/config"
	"github.com/penggege18/ddns-ipv6-go/v6/monitor"
	"github.com/penggege18/ddns-ipv6-go/v6/util"
	"github.com/penggege18/ddns-ipv6-go/v6/util/osutil"
	"github.com/penggege18/ddns-ipv6-go/v6/web"
	"github.com/kardianos/service"
)

// ddns-ipv6-go 版本
// ddns-ipv6-go version
var versionFlag = flag.Bool("v", false, "ddns-ipv6-go version")

// 监听地址
var listen = flag.String("l", ":9876", "Listen address")

// IP 检测频率(秒), 配置文件中未设置时使用
var every = flag.Int("f", 60, "IP check frequency(seconds)")

// 服务管理
var serviceType = flag.String("s", "", "Service management (install|uninstall|restart)")

// 配置文件路径
var configFilePath = flag.String("c", util.GetConfigFilePathDefault(), "Custom configuration file path")

// Web 服务
var noWebService = flag.Bool("noweb", false, "No web service")

// 跳过验证证书
var skipVerify = flag.Bool("skipVerify", false, "Skip certificate verification")

// 自定义 DNS 服务器
var customDNS = flag.String("dns", "", "Custom DNS server address, example: 8.8.8.8")

// 重置密码
var newPassword = flag.String("resetPassword", "", "Reset password to the one entered")

// 后台运行
var daemonize = flag.Bool("d", false, "Run in background (daemon/detached)")

//go:embed static
var staticEmbeddedFiles embed.FS

//go:embed favicon.ico
var faviconEmbeddedFile embed.FS

// version
var version = "DEV"

func main() {
	flag.Parse()
	if *versionFlag {
		fmt.Println(version)
		return
	}

	if *daemonize && os.Getenv("DDNS_GO_DAEMON") != "1" {
		if err := runAsDaemon(); err != nil {
			log.Fatalf("Daemonize failed: %v", err)
		}
		return
	}

	// 安卓 go/src/time/zoneinfo_android.go 固定localLoc 为 UTC
	if runtime.GOOS == "android" {
		util.FixTimezone()
	}
	// 检查监听地址
	if _, err := net.ResolveTCPAddr("tcp", *listen); err != nil {
		log.Fatalf("Parse listen address failed! Exception: %s", err)
	}
	// 设置版本号
	os.Setenv(web.VersionEnv, version)
	// 设置配置文件路径
	if *configFilePath != "" {
		absPath, _ := filepath.Abs(*configFilePath)
		os.Setenv(util.ConfigFilePathENV, absPath)
	}
	// 重置密码
	if *newPassword != "" {
		if _, statErr := os.Stat(*configFilePath); os.IsNotExist(statErr) {
			util.Log("配置文件 %s 不存在, 可通过-c指定配置文件", *configFilePath)
		} else {
			conf, _ := config.GetConfigCached()
			conf.ResetPassword(*newPassword)
		}
		return
	}
	// 设置跳过证书验证
	if *skipVerify {
		util.SetInsecureSkipVerify()
	}
	// 设置自定义DNS
	if *customDNS != "" {
		util.SetDNS(*customDNS)
	}
	switch *serviceType {
	case "install":
		installService()
	case "uninstall":
		uninstallService()
	case "restart":
		restartService()
	default:
		if util.IsRunInDocker() || os.Getenv("DDNS_GO_DAEMON") == "1" {
			run()
		} else {
			s := getService()
			status, _ := s.Status()
			if status != service.StatusUnknown {
				// 以服务方式运行
				s.Run()
			} else {
				// 非服务方式运行
				switch s.Platform() {
				case "windows-service":
					util.Log("可使用 .\\ddns-ipv6-go.exe -s install 安装服务运行")
				default:
					util.Log("可使用 sudo ./ddns-ipv6-go -s install 安装服务运行")
				}
				run()
			}
		}
	}
}

func run() {
	conf, _ := config.GetConfigCached()
	// 初始化语言
	util.InitLogLang(conf.Lang)

	if !*noWebService {
		go func() {
			// 启动web服务
			err := runWebServer()
			if err != nil {
				log.Println(err)
				time.Sleep(time.Minute)
				os.Exit(1)
			}
		}()
	}

	// 初始化备用DNS
	util.InitBackupDNS(*customDNS, conf.Lang)

	// 等待网络连接
	util.WaitInternet([]string{"https://www.baidu.com", "https://api.ipify.org"})

	// 启动 certbot 每日定时更新证书
	cert.StartScheduler()

	// 实时检测 IP 地址
	monitor.Start(time.Duration(*every) * time.Second)
}

func staticFsFunc(writer http.ResponseWriter, request *http.Request) {
	http.FileServer(http.FS(staticEmbeddedFiles)).ServeHTTP(writer, request)
}

func faviconFsFunc(writer http.ResponseWriter, request *http.Request) {
	http.FileServer(http.FS(faviconEmbeddedFile)).ServeHTTP(writer, request)
}

func runWebServer() error {
	// 启动静态文件服务
	http.HandleFunc("/static/", web.AuthAssert(staticFsFunc))
	http.HandleFunc("/favicon.ico", web.AuthAssert(faviconFsFunc))
	http.HandleFunc("/login", web.AuthAssert(web.Login))
	http.HandleFunc("/loginFunc", web.AuthAssert(web.LoginFunc))

	http.HandleFunc("/", web.Auth(web.Writing))
	http.HandleFunc("/save", web.Auth(web.Save))
	http.HandleFunc("/setLang", web.Auth(web.SetLang))
	http.HandleFunc("/status", web.Auth(web.Status))
	http.HandleFunc("/checkNow", web.Auth(web.CheckNow))
	http.HandleFunc("/certRenew", web.Auth(web.CertRenew))
	http.HandleFunc("/logs", web.Auth(web.Logs))
	http.HandleFunc("/clearLog", web.Auth(web.ClearLog))
	http.HandleFunc("/webhookTest", web.Auth(web.WebhookTest))
	http.HandleFunc("/logout", web.Auth(web.Logout))

	util.Log("监听 %s", *listen)

	l, err := net.Listen("tcp", *listen)
	if err != nil {
		return errors.New(util.LogStr("监听端口发生异常, 请检查端口是否被占用! %s", err))
	}

	return http.Serve(l, nil)
}

// 以守护/分离进程方式运行（Unix 使用 setsid，Windows 使用 DETACHED_PROCESS）
func runAsDaemon() error {
	exe, err := os.Executable()
	if err != nil {
		return err
	}

	// 过滤掉 -d 参数
	args := make([]string, 0, len(os.Args))
	args = append(args, exe)
	for i := 1; i < len(os.Args); i++ {
		if os.Args[i] == "-d" {
			continue
		}
		args = append(args, os.Args[i])
	}

	// 重定向到系统空设备
	nullFile, err := os.OpenFile(os.DevNull, os.O_RDWR, 0)
	if err != nil {
		return err
	}
	defer nullFile.Close()

	proc, err := osutil.StartDetachedProcess(exe, args, nullFile)
	if err != nil {
		return err
	}

	return proc.Release()
}

type program struct{}

func (p *program) Start(s service.Service) error {
	// Start should not block. Do the actual work async.
	go p.run()
	return nil
}
func (p *program) run() {
	run()
}
func (p *program) Stop(s service.Service) error {
	// Stop should not block. Return with a few seconds.
	return nil
}

func getService() service.Service {
	options := make(service.KeyValue)
	var depends []string

	// 确保服务等待网络就绪后再启动
	switch service.ChosenSystem().String() {
	case "unix-systemv":
		options["SysvScript"] = sysvScript
	case "windows-service":
		// 将 Windows 服务的启动类型设为自动(延迟启动)
		options["DelayedAutoStart"] = true
	default:
		// 向 Systemd 添加网络依赖
		depends = append(depends, "Requires=network.target",
			"After=network-online.target")
	}

	svcConfig := &service.Config{
		Name:         "ddns-ipv6-go",
		DisplayName:  "ddns-ipv6-go",
		Description:  "Monitor local IPv4/IPv6 address, notify via Webhook and renew IP certificates with certbot when IPv6 changes",
		Arguments:    []string{"-l", *listen, "-f", strconv.Itoa(*every), "-c", *configFilePath},
		Dependencies: depends,
		Option:       options,
	}

	if *noWebService {
		svcConfig.Arguments = append(svcConfig.Arguments, "-noweb")
	}

	if *skipVerify {
		svcConfig.Arguments = append(svcConfig.Arguments, "-skipVerify")
	}

	if *customDNS != "" {
		svcConfig.Arguments = append(svcConfig.Arguments, "-dns", *customDNS)
	}

	prg := &program{}
	s, err := service.New(prg, svcConfig)
	if err != nil {
		log.Fatalln(err)
	}
	return s
}

// 卸载服务
func uninstallService() {
	s := getService()
	s.Stop()
	if service.ChosenSystem().String() == "unix-systemv" {
		if _, err := exec.Command("/etc/init.d/ddns-ipv6-go", "stop").Output(); err != nil {
			log.Println(err)
		}
	}
	if err := s.Uninstall(); err == nil {
		util.Log("ddns-ipv6-go 服务卸载成功")
	} else {
		util.Log("ddns-ipv6-go 服务卸载失败, 异常信息: %s", err)
	}
}

// 安装服务
func installService() {
	s := getService()

	status, err := s.Status()
	if err != nil && status == service.StatusUnknown {
		// 服务未知，创建服务
		if err = s.Install(); err == nil {
			s.Start()
			util.Log("安装 ddns-ipv6-go 服务成功! 请打开浏览器并进行配置")
			if service.ChosenSystem().String() == "unix-systemv" {
				if _, err := exec.Command("/etc/init.d/ddns-ipv6-go", "enable").Output(); err != nil {
					log.Println(err)
				}
				if _, err := exec.Command("/etc/init.d/ddns-ipv6-go", "start").Output(); err != nil {
					log.Println(err)
				}
			}
			return
		}
		util.Log("安装 ddns-ipv6-go 服务失败, 异常信息: %s", err)
	}

	if status != service.StatusUnknown {
		util.Log("ddns-ipv6-go 服务已安装, 无需再次安装")
	}
}

// 重启服务
func restartService() {
	s := getService()
	status, err := s.Status()
	if err == nil {
		if status == service.StatusRunning {
			if err = s.Restart(); err == nil {
				util.Log("重启 ddns-ipv6-go 服务成功")
			}
		} else if status == service.StatusStopped {
			if err = s.Start(); err == nil {
				util.Log("启动 ddns-ipv6-go 服务成功")
			}
		}
	} else {
		util.Log("ddns-ipv6-go 服务未安装, 请先安装服务")
	}
}

const sysvScript = `#!/bin/sh /etc/rc.common
DESCRIPTION="{{.Description}}"
cmd="{{.Path}}{{range .Arguments}} {{.|cmd}}{{end}}"
name="ddns-ipv6-go"
pid_file="/var/run/$name.pid"
stdout_log="/var/log/$name.log"
stderr_log="/var/log/$name.err"
START=99
get_pid() {
    cat "$pid_file"
}
is_running() {
    [ -f "$pid_file" ] && cat /proc/$(get_pid)/stat > /dev/null 2>&1
}
start() {
	if is_running; then
		echo "Already started"
	else
		echo "Starting $name"
		{{if .WorkingDirectory}}cd '{{.WorkingDirectory}}'{{end}}
		$cmd >> "$stdout_log" 2>> "$stderr_log" &
		echo $! > "$pid_file"
		if ! is_running; then
			echo "Unable to start, see $stdout_log and $stderr_log"
			exit 1
		fi
	fi
}
stop() {
	if is_running; then
		echo -n "Stopping $name.."
		kill $(get_pid)
		for i in $(seq 1 10)
		do
			if ! is_running; then
				break
			fi
			echo -n "."
			sleep 1
		done
		echo
		if is_running; then
			echo "Not stopped; may still be shutting down or shutdown may have failed"
			exit 1
		else
			echo "Stopped"
			if [ -f "$pid_file" ]; then
				rm "$pid_file"
			fi
		fi
	else
		echo "Not running"
	fi
}
restart() {
	stop
	if is_running; then
		echo "Unable to stop, will not attempt to start"
		exit 1
	fi
	start
}
`
