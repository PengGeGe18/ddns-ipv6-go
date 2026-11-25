package config

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"

	"github.com/penggege18/ddns-ipv6-go/v6/util"
)

// Webhook Webhook
type Webhook struct {
	WebhookURL         string
	WebhookRequestBody string
	WebhookHeaders     string
}

// hasJSONPrefix returns true if the string starts with a JSON open brace.
func hasJSONPrefix(s string) bool {
	return strings.HasPrefix(s, "{") || strings.HasPrefix(s, "[")
}

// Notify 发送 IPv6 变化通知
// 支持占位符: #{ipv4Addr} #{ipv6Addr} #{oldIpv6Addr} #{timestamp}
func (wh *Webhook) Notify(oldIPv6, newIPv6, ipv4 string) error {
	if wh == nil || strings.TrimSpace(wh.WebhookURL) == "" {
		return nil
	}

	method := "GET"
	postPara := ""
	timestamp := strconv.FormatInt(time.Now().UTC().Unix(), 10)
	contentType := "application/x-www-form-urlencoded"

	replace := func(org string) string {
		return strings.NewReplacer(
			"#{ipv4Addr}", ipv4,
			"#{ipv6Addr}", newIPv6,
			"#{oldIpv6Addr}", oldIPv6,
			"#{timestamp}", timestamp,
		).Replace(org)
	}

	if strings.TrimSpace(wh.WebhookRequestBody) != "" {
		method = "POST"
		postPara = replace(wh.WebhookRequestBody)
		if json.Valid([]byte(postPara)) {
			contentType = "application/json"
		} else if hasJSONPrefix(postPara) {
			util.Log("Webhook中的 RequestBody JSON 无效")
		}
	}

	requestURL := replace(wh.WebhookURL)
	u, err := url.Parse(requestURL)
	if err != nil {
		util.Log("Webhook配置中的URL不正确")
		return err
	}

	q, _ := url.ParseQuery(u.RawQuery)
	u.RawQuery = q.Encode()

	req, err := http.NewRequest(method, u.String(), strings.NewReader(postPara))
	if err != nil {
		util.Log("Webhook调用失败! 异常信息：%s", err)
		return err
	}

	for key, value := range extractHeaders(wh.WebhookHeaders) {
		req.Header.Add(key, value)
	}
	req.Header.Add("content-type", contentType)

	clt := util.CreateHTTPClient()
	resp, err := clt.Do(req)
	body, err := util.GetHTTPResponseOrg(resp, err)
	if err != nil {
		util.Log("Webhook调用失败! 异常信息：%s", err)
		return err
	}

	util.Log("Webhook调用成功! 返回数据：%s", string(body))
	return nil
}

// extractHeaders converts s into a map of headers.
func extractHeaders(s string) map[string]string {
	lines := util.SplitLines(s)
	headers := make(map[string]string, len(lines))

	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}

		parts := strings.SplitN(line, ":", 2)
		if len(parts) != 2 {
			util.Log("Webhook Header不正确: %s", line)
			continue
		}

		k, v := strings.TrimSpace(parts[0]), strings.TrimSpace(parts[1])
		headers[k] = v
	}

	return headers
}

// WebhookTest 使用假数据测试 Webhook
func (wh *Webhook) WebhookTest() error {
	if strings.TrimSpace(wh.WebhookURL) == "" {
		return fmt.Errorf("%s", util.LogStr("请输入Webhook的URL"))
	}
	return wh.Notify("2001:db8::1", "2001:db8::2", "192.168.1.100")
}
