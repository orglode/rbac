package main

import (
	"github.com/orglode/rbac/conf"
	"github.com/orglode/rbac/server/http"
	"github.com/orglode/rbac/service"
)

func main() {
	//初始化配置
	cfg := conf.Init()
	//初始化service
	svc := service.NewService(cfg)
	//初始化http控件
	http.Init(svc, cfg)
}
