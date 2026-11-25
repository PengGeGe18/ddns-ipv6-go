package web

import (
	"encoding/json"
	"net/http"

	"github.com/penggege18/ddns-ipv6-go/v6/cert"
	"github.com/penggege18/ddns-ipv6-go/v6/monitor"
)

// CheckNow 立即检测 IP
func CheckNow(writer http.ResponseWriter, request *http.Request) {
	monitor.Trigger()
	writeResult(writer, "ok", "已触发检测")
}

// CertRenew 手动更新证书
// 当前证书不包含当前 IPv6 时重新申请, 否则执行强制续期
func CertRenew(writer http.ResponseWriter, request *http.Request) {
	ip := monitor.GetStatus().IPv6
	if ip != "" && !cert.Cover(ip) {
		if !cert.IssueAsync(ip, "") {
			writeResult(writer, "证书任务正在运行, 请稍后再试", "")
			return
		}
		writeResult(writer, "ok", "当前证书不包含该地址, 已开始重新申请")
		return
	}

	if err := cert.RenewAsync(true); err != nil {
		writeResult(writer, err.Error(), "")
		return
	}
	writeResult(writer, "ok", "已开始更新证书")
}

func writeResult(writer http.ResponseWriter, result, msg string) {
	writer.Header().Set("Content-Type", "application/json; charset=utf-8")
	byt, _ := json.Marshal(map[string]string{"result": result, "msg": msg})
	writer.Write(byt)
}
