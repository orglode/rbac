package service

import (
	"github.com/orglode/navigator/api"
	"github.com/orglode/navigator/conf"
	"github.com/orglode/navigator/dao"
	"github.com/orglode/navigator/manager"
	"github.com/orglode/navigator/model"
	"go.uber.org/zap"
)

type Service struct {
	c        *conf.Config
	mgr      *manager.Manager
	dao      *dao.Dao
	logger   *zap.Logger
	Response *model.BaseResponse
}

func NewService(conf *conf.Config) *Service {
	return &Service{
		c:      conf,
		mgr:    manager.NewManager(conf),
		dao:    dao.NewDao(conf),
		logger: api.InitLogger(),
		Response: &model.BaseResponse{
			Code: api.Success,
		},
	}
}
