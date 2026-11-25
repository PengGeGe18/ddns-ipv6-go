package web

import (
	"embed"
	"fmt"
	"html/template"
	"net/http"
	"os"

	"github.com/penggege18/ddns-ipv6-go/v6/config"
)

//go:embed writing.html
var writingEmbedFile embed.FS

const VersionEnv = "DDNS_GO_VERSION"

// Writing 填写信息
func Writing(writer http.ResponseWriter, request *http.Request) {
	tmpl, err := template.ParseFS(writingEmbedFile, "writing.html")
	if err != nil {
		fmt.Println("Error happened..")
		fmt.Println(err)
		return
	}

	conf, err := config.GetConfigCached()
	// 默认禁止公网访问
	if err != nil {
		conf.NotAllowWanAccess = true
	}

	ipv4, ipv6, _ := config.GetNetInterface()

	err = tmpl.Execute(writer, struct {
		Conf    config.Config
		Version string
		Ipv4    []config.NetInterface
		Ipv6    []config.NetInterface
	}{
		Conf:    conf,
		Version: os.Getenv(VersionEnv),
		Ipv4:    ipv4,
		Ipv6:    ipv6,
	})
	if err != nil {
		fmt.Println("Error happened..")
		fmt.Println(err)
	}
}
