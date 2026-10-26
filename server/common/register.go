package common

import (
	"github.com/gin-gonic/gin"
)

// RouteRegistrar 模块路由注册器：后续新菜单/功能实现此接口并挂到 router 即可
type RouteRegistrar interface {
	Register(rg *gin.RouterGroup)
}
