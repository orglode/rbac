package http

import (
	"encoding/gob"
	"fmt"
	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/cookie"
	"github.com/gin-gonic/gin"
	"github.com/orglode/navigator/conf"
	"github.com/orglode/navigator/model"
	"github.com/orglode/navigator/service"
	"io"
	"os"
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
	if conf.Env == model.EnvProduction {
		//注册session结构体
		gob.Register(model.UserSessionInfo{})
		gin.SetMode(gin.ReleaseMode)
	}
	//gin日志模式
	router.Use(gin.LoggerWithConfig(initGinLog()))
	router.Use(gin.Recovery())
	//初始化session
	initSession(router)
	//初始化路由
	initRouter(router)
	router.Run(conf.HttpAddr)
}
func initSession(router *gin.Engine) {
	store := cookie.NewStore([]byte("secret"))
	router.Use(sessions.Sessions("session", store))
}

// 初始化gin日志库
func initGinLog() gin.LoggerConfig {
	f, _ := os.OpenFile("./log/app.log", os.O_CREATE|os.O_APPEND|os.O_RDWR, 0777)
	var logConf = gin.LoggerConfig{
		Formatter: func(param gin.LogFormatterParams) string {
			return fmt.Sprintf("客户端IP:%s,请求时间:[%s],请求方式:%s,请求地址:%s,响应时间:%s,错误信息:%s\n",
				param.ClientIP,
				param.TimeStamp.Format("2006年01月02日 15:03:04"),
				param.Method,
				param.Path,
				param.Latency,
				param.ErrorMessage,
			)
		},
		Output: io.MultiWriter(os.Stdout, f),
	}
	return logConf
}
