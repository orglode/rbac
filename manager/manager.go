package manager

import (
	"github.com/kirinlabs/HttpRequest"
	"github.com/orglode/navigator/api/logger"
	"github.com/orglode/navigator/conf"
	"go.uber.org/zap"
)

type Manager struct {
	c          *conf.Config
	logger     *zap.Logger
	httpClient *HttpRequest.Request
}

func NewManager(conf *conf.Config) *Manager {
	return &Manager{
		c:          conf,
		httpClient: HttpRequest.NewRequest(),
		logger:     logger.InitLogger(),
	}
}
