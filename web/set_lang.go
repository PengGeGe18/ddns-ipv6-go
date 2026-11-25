package web

import (
	"encoding/json"
	"net/http"
	"strings"

	"github.com/penggege18/ddns-ipv6-go/v6/config"
	"github.com/penggege18/ddns-ipv6-go/v6/util"
)

func SetLang(writer http.ResponseWriter, request *http.Request) {
	conf, _ := config.GetConfigCached()

	var data struct {
		Lang string `json:"Lang"`
	}

	if err := json.NewDecoder(request.Body).Decode(&data); err != nil {
		returnError(writer, util.LogStr("数据解析失败, 请刷新页面重试"))
		return
	}

	lang := strings.TrimSpace(data.Lang)
	if lang == "" {
		lang = conf.Lang
	}

	conf.Lang = util.InitLogLang(lang)
	if err := conf.SaveConfig(); err != nil {
		returnError(writer, err.Error())
		return
	}

	byt, _ := json.Marshal(map[string]string{"result": "ok", "lang": conf.Lang})
	writer.Write(byt)
}
