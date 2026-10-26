package http

import (
	"time"

	"rbac/conf"
	"rbac/model"
	"rbac/server/common"
	"rbac/service"

	"github.com/gin-gonic/gin"
	logger "github.com/orglode/hades/logger_v2"
	"github.com/orglode/hades/trace"
)

// Init 初始化引擎、注册路由并启动 HTTP 服务
func Init(s *service.Service, cfg *conf.Config) {
	r := gin.New()
	gin.DisableConsoleColor()
	if cfg.Server.Env == model.EnvProduction {
		gin.SetMode(gin.ReleaseMode)
	}

	r.Use(trace.TraceIDMiddleware(), logger.GinLogger(), common.ErrorHandlerMiddleware(), gin.Recovery())
	initRouter(r, s)
	_ = r.Run(cfg.Server.Addr)
}

func health(c *gin.Context) {
	common.ResponseSuccess(c, gin.H{
		"time": time.Now().Unix(),
		"ip":   c.ClientIP(),
	})
}
