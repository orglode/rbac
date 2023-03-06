package service

import (
	"github.com/orglode/rbac/api"
	"github.com/orglode/rbac/conf"
	"github.com/orglode/rbac/dao"
	"github.com/orglode/rbac/manager"
	"github.com/orglode/rbac/model"
)

type Service struct {
	c        *conf.Config
	mgr      *manager.Manager
	dao      *dao.Dao
	Response *model.BaseResponse
}

func NewService(conf *conf.Config) *Service {
	return &Service{
		c:   conf,
		mgr: manager.NewManager(conf),
		dao: dao.NewDao(conf),
		Response: &model.BaseResponse{
			Code: api.Success,
		},
	}
}
