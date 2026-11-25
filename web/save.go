package web

import (
	"encoding/json"
	"net/http"
	"strings"

	"github.com/penggege18/ddns-ipv6-go/v6/config"
	"github.com/penggege18/ddns-ipv6-go/v6/monitor"
	"github.com/penggege18/ddns-ipv6-go/v6/util"
)

// Save 保存
func Save(writer http.ResponseWriter, request *http.Request) {
	result := checkAndSave(request)
	writer.Header().Set("Content-Type", "application/json; charset=utf-8")
	byt, _ := json.Marshal(map[string]string{"result": result})
	writer.Write(byt)
}

func checkAndSave(request *http.Request) string {
	conf, _ := config.GetConfigCached()

	// 从请求中读取 JSON 数据
	var data struct {
		Username           string `json:"Username"`
		Password           string `json:"Password"`
		Lang               string `json:"Lang"`
		NotAllowWanAccess  bool   `json:"NotAllowWanAccess"`
		WebhookURL         string `json:"WebhookURL"`
		WebhookRequestBody string `json:"WebhookRequestBody"`
		WebhookHeaders     string `json:"WebhookHeaders"`
		Interval           int    `json:"Interval"`
		IPv4               config.IPConfig
		IPv6               config.IPConfig
		Cert               config.CertConfig
	}

	// 解析请求中的 JSON 数据
	err := json.NewDecoder(request.Body).Decode(&data)
	if err != nil {
		return util.LogStr("数据解析失败, 请刷新页面重试")
	}
	usernameNew := strings.TrimSpace(data.Username)
	passwordNew := data.Password

	if strings.TrimSpace(data.Lang) != "" {
		conf.Lang = util.InitLogLang(data.Lang)
	} else {
		conf.Lang = util.InitLogLang(conf.Lang)
	}

	conf.NotAllowWanAccess = data.NotAllowWanAccess
	conf.WebhookURL = strings.TrimSpace(data.WebhookURL)
	conf.WebhookRequestBody = strings.TrimSpace(data.WebhookRequestBody)
	conf.WebhookHeaders = strings.TrimSpace(data.WebhookHeaders)
	conf.Interval = data.Interval

	// IPv4/IPv6 获取方式
	conf.IPv4 = data.IPv4
	conf.IPv6 = data.IPv6
	conf.IPv4.GetType = normalizeGetType(conf.IPv4.GetType)
	conf.IPv6.GetType = normalizeGetType(conf.IPv6.GetType)
	conf.IPv4.URL = strings.TrimSpace(conf.IPv4.URL)
	conf.IPv6.URL = strings.TrimSpace(conf.IPv6.URL)
	conf.IPv4.Cmd = strings.TrimSpace(conf.IPv4.Cmd)
	conf.IPv6.Cmd = strings.TrimSpace(conf.IPv6.Cmd)
	conf.IPv6.Ipv6Reg = strings.TrimSpace(conf.IPv6.Ipv6Reg)

	// 证书配置
	conf.Cert = data.Cert
	conf.Cert.CertbotPath = strings.TrimSpace(conf.Cert.CertbotPath)
	conf.Cert.Email = strings.TrimSpace(conf.Cert.Email)
	conf.Cert.Webroot = strings.TrimSpace(conf.Cert.Webroot)
	conf.Cert.LinkDir = strings.TrimSpace(conf.Cert.LinkDir)
	conf.Cert.Owner = strings.TrimSpace(conf.Cert.Owner)
	conf.Cert.NginxReloadCmd = strings.TrimSpace(conf.Cert.NginxReloadCmd)
	conf.Cert.ExtraArgs = strings.TrimSpace(conf.Cert.ExtraArgs)
	conf.Cert.Profile = strings.TrimSpace(conf.Cert.Profile)

	// 填充默认值并修正非法配置
	conf.Normalize()

	// 如果新密码不为空则检查是否够强, 内/外网要求强度不同
	conf.Username = usernameNew
	if passwordNew != "" {
		hashedPwd, err := conf.CheckPassword(passwordNew)
		if err != nil {
			return err.Error()
		}
		conf.Password = hashedPwd
	}

	// 帐号密码不能为空
	if conf.Username == "" || conf.Password == "" {
		return util.LogStr("必须输入用户名/密码")
	}

	// 保存到用户目录
	if err := conf.SaveConfig(); err != nil {
		return err.Error()
	}

	// 立即按新配置检测一次
	monitor.Trigger()

	return "ok"
}

func normalizeGetType(t string) string {
	t = strings.TrimSpace(t)
	switch t {
	case "url", "netInterface", "cmd":
		return t
	default:
		return "netInterface"
	}
}
