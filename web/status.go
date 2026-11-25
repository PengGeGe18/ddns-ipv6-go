package web

import (
	"encoding/json"
	"net/http"

	"github.com/penggege18/ddns-ipv6-go/v6/cert"
	"github.com/penggege18/ddns-ipv6-go/v6/config"
	"github.com/penggege18/ddns-ipv6-go/v6/monitor"
)

// Status 返回运行状态(IP、证书任务等)
func Status(writer http.ResponseWriter, request *http.Request) {
	conf, _ := config.GetConfigCached()

	writer.Header().Set("Content-Type", "application/json; charset=utf-8")
	json.NewEncoder(writer).Encode(map[string]interface{}{
		"monitor":     monitor.GetStatus(),
		"cert":        cert.GetStatus(),
		"certEnabled": conf.Cert.Enable,
		"interval":    conf.Interval,
	})
}
