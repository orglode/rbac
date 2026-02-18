package manager

import (
	"github.com/kirinlabs/HttpRequest"
	"net/http"
	"rbac/conf"
	"time"
)

type Manager struct {
	c          *conf.Config
	httpClient *HttpRequest.Request
}

func NewManager(conf *conf.Config) *Manager {
	transport := &http.Transport{
		MaxIdleConns:        100,              // 最大空闲连接数
		MaxIdleConnsPerHost: 10,               // 每个主机最大空闲连接
		IdleConnTimeout:     90 * time.Second, // 空闲连接超时时间
		DisableKeepAlives:   false,            // 启用 KeepAlive
	}

	// 初始化客户端并应用配置
	client := HttpRequest.NewRequest().
		Transport(transport).        // 应用自定义 Transport
		SetTimeout(10 * time.Second) // 设置全局默认超时
	return &Manager{
		c:          conf,
		httpClient: client,
	}
}
