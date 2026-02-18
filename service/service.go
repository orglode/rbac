package service

import (
	"rbac/conf"
	"rbac/dao"
	"rbac/manager"
)

type Service struct {
	c           *conf.Config
	mgr         *manager.Manager
	dao         *dao.Dao
	RbacService *RbacService
}

func NewService(conf *conf.Config) *Service {
	service := &Service{
		c:   conf,
		mgr: manager.NewManager(conf),
		dao: dao.NewDao(conf),
	}
	service.RbacService = NewRbacService(service)
	return service
}
