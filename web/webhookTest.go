package web

import (
	"encoding/json"
	"net/http"

	"github.com/penggege18/ddns-ipv6-go/v6/config"
)

// WebhookTest 使用假数据测试 Webhook
func WebhookTest(writer http.ResponseWriter, request *http.Request) {
	var data struct {
		URL         string `json:"URL"`
		RequestBody string `json:"RequestBody"`
		Headers     string `json:"Headers"`
	}
	if err := json.NewDecoder(request.Body).Decode(&data); err != nil {
		writeResult(writer, "数据解析失败, 请刷新页面重试", "")
		return
	}

	wh := config.Webhook{
		WebhookURL:         data.URL,
		WebhookRequestBody: data.RequestBody,
		WebhookHeaders:     data.Headers,
	}

	if err := wh.WebhookTest(); err != nil {
		writeResult(writer, err.Error(), "")
		return
	}
	writeResult(writer, "ok", "测试请求已发送, 请查看日志确认结果")
}
