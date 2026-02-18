package http

import (
	"rbac/api/jwt"

	"github.com/gin-gonic/gin"
)

func initRouter(r *gin.Engine) {
	// 注册Rbac路由
	registerRbacRoutes(r)

	admin := r.Group("/api/admin/", jwt.AuthMiddleware())
	{
		admin.GET("/heath", func(c *gin.Context) {})
	}

}
