// Package monitor 实时检测本机 IPv4/IPv6 地址。
// 当 IPv6 地址发生变化时, 触发 Webhook 通知并使用 certbot 为新的地址申请证书。
package monitor

import (
	"sync"
	"time"

	"github.com/penggege18/ddns-ipv6-go/v6/cert"
	"github.com/penggege18/ddns-ipv6-go/v6/config"
	"github.com/penggege18/ddns-ipv6-go/v6/util"
)

// 最小检测间隔
const minInterval = 5 * time.Second

// Snapshot 当前监测状态
type Snapshot struct {
	IPv4       string    `json:"ipv4"`
	IPv6       string    `json:"ipv6"`
	LastCheck  time.Time `json:"lastCheck"`
	LastChange time.Time `json:"lastChange"`
	LastError  string    `json:"lastError"`
	Checking   bool      `json:"checking"`
}

var (
	mu        sync.Mutex
	snap      Snapshot
	lastError string
	checkCh   = make(chan struct{}, 1)
)

// Start 启动实时检测(阻塞运行)
func Start(fallback time.Duration) {
	loop(fallback)
}

// Trigger 立即触发一次检测
func Trigger() {
	select {
	case checkCh <- struct{}{}:
	default:
	}
}

// GetStatus 获取当前监测状态
func GetStatus() Snapshot {
	mu.Lock()
	defer mu.Unlock()
	return snap
}

func loop(fallback time.Duration) {
	for {
		conf, err := config.GetConfigCached()
		interval := fallback
		if err == nil && conf.Interval > 0 {
			interval = time.Duration(conf.Interval) * time.Second
		}
		if interval < minInterval {
			interval = minInterval
		}

		checkOnce()

		select {
		case <-checkCh:
			// 立即检测
		case <-time.After(interval):
		}
	}
}

func checkOnce() {
	conf, err := config.GetConfigCached()

	mu.Lock()
	snap.Checking = true
	mu.Unlock()
	defer func() {
		mu.Lock()
		snap.Checking = false
		snap.LastCheck = time.Now()
		mu.Unlock()
	}()

	if err != nil {
		setError("读取配置失败, 无法检测 IP: " + err.Error())
		return
	}

	var v4, v6 string
	if conf.IPv4.Enable {
		v4 = conf.GetIpv4Addr()
	}
	if conf.IPv6.Enable {
		v6 = conf.GetIpv6Addr()
	}

	mu.Lock()
	prevV6 := snap.IPv6
	if v4 != "" {
		snap.IPv4 = v4
	}
	if v6 != "" {
		snap.IPv6 = v6
	}
	changed := prevV6 != "" && v6 != "" && prevV6 != v6
	first := prevV6 == "" && v6 != ""
	if changed {
		snap.LastChange = time.Now()
	}
	mu.Unlock()

	if v6 == "" && conf.IPv6.Enable {
		setError("获取 IPv6 地址失败")
	} else if v4 == "" && conf.IPv4.Enable {
		setError("获取 IPv4 地址失败")
	} else {
		setError("")
	}

	if v6 == "" {
		return
	}

	if changed {
		util.Log("检测到 IPv6 地址变化: %s -> %s", prevV6, v6)

		// Webhook 通知
		wh := conf.Webhook
		ipv4 := v4
		go func() {
			if err := wh.Notify(prevV6, v6, ipv4); err != nil {
				util.Log("Webhook 通知失败: %s", err)
			}
		}()

		// 为新地址申请证书
		cert.IssueAsync(v6, prevV6)
		return
	}

	if first {
		util.Log("当前 IPv6 地址: %s", v6)
	}

	// 证书缺失或与当前地址不匹配时自动补申请
	cert.EnsureAsync(v6)
}

// setError 记录并输出错误信息(相同错误只输出一次)
func setError(msg string) {
	mu.Lock()
	snap.LastError = msg
	changed := msg != lastError
	lastError = msg
	mu.Unlock()

	if msg != "" && changed {
		util.Log("%s", msg)
	}
}
