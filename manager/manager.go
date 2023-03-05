package manager

import (
	"github.com/kirinlabs/HttpRequest"
	"github.com/orglode/rbac/conf"
)

type Manager struct {
	c          *conf.Config
	httpClient *HttpRequest.Request
}

func NewManager(conf *conf.Config) *Manager {
	return &Manager{
		c:          conf,
		httpClient: HttpRequest.NewRequest(),
	}
}
