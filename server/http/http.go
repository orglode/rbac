package http

import (
	"github.com/gin-gonic/gin"
	"github.com/orglode/navigator/conf"
	"github.com/orglode/navigator/service"
)

var (
	svc *service.Service
)

func Init(s *service.Service, conf *conf.Config) {
	svc = s
	rbac := gin.Default()
	initRouter(rbac)
	rbac.Run(conf.HttpAddr)
}
