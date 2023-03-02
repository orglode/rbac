package main

import (
	"encoding/gob"
	"github.com/orglode/navigator/conf"
	"github.com/orglode/navigator/model"
	"github.com/orglode/navigator/server/http"
	"github.com/orglode/navigator/service"
)

func main() {
	//注册session结构体-- 无状态服务可去掉
	gob.Register(model.UserSessionInfo{})
	//初始化配置
	cfg := conf.Init()
	//初始化service
	svc := service.NewService(cfg)
	//初始化http控件
	http.Init(svc, cfg)
}
