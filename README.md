# DDNS-IPv6-GO

实时监测本机 IPv4/IPv6 地址，当 IPv6 地址变化时通过 Webhook 通知，并使用 certbot 自动为新的 IPv6 地址申请证书、替换旧证书，随后重启 Nginx 使证书生效。

## 功能特性

- 实时获取本机 IPv4 / IPv6 地址，检测频率可调（Web 界面或 `-f` 参数）
- 支持接口、网卡、命令三种方式获取 IP，IPv6 支持正则表达式匹配指定地址
- IPv6 地址变化时触发 Webhook 通知，支持自定义 URL、RequestBody 与请求头
- Webhook 内容支持 `#{ipv4Addr}`、`#{ipv6Addr}`、`#{oldIpv6Addr}`、`#{timestamp}` 变量
- 集成 certbot：为新的 IPv6 地址自动申请 IP 证书并替换旧证书
- 证书统一放置于固定目录（默认 `/etc/ssl/ddns-ipv6`）供 Nginx 8443 端口使用，更换后自动重启 Nginx
- 支持 certbot 每日定时续期，也可在 Web 界面手动更新证书
- Web 界面提供状态面板：当前 IP、证书地址与到期时间、最近任务结果、下次续期时间
- Web 界面可配置 IP 检测方式、Webhook、certbot、账号密码，并可查看日志
- 支持以系统服务方式运行，支持 Windows / Linux / macOS
- 支持中英文界面与日志
