package http

import (
	"rbac/conf"
	"rbac/model"
	"rbac/service"

	"github.com/gin-gonic/gin"
	logger "github.com/orglode/hades/logger_v2"
	"github.com/orglode/hades/trace"
)

var (
	svc *service.Service
)

func Init(s *service.Service, conf *conf.Config) {
	svc = s
	router := gin.New()

	//禁用调式终端颜色
	gin.DisableConsoleColor()

	//判断环境 是否开启debug模式
	if conf.Server.Env == model.EnvProduction {
		gin.SetMode(gin.ReleaseMode)
	}

	// 注册trace中间件
	router.Use(trace.TraceIDMiddleware())

	// gin日志模式
	router.Use(logger.GinLogger())

	// 报错输出中间件
	router.Use(ErrorHandlerMiddleware())

	// gin恢复模式
	router.Use(gin.Recovery())

	// 初始化路由
	initRouter(router)

	// 启动服务
	router.Run(conf.Server.Addr)
}
