package http

import (
	"rbac/api/jwt"
	"rbac/server/common"
	"rbac/service"

	"github.com/gin-gonic/gin"
)

func initRouter(r *gin.Engine, svc *service.Service) {
	r.GET("/api/admin/health", health)

	rbac := NewRbacHandler(svc)
	rbac.RegisterPublic(r)
	
	message := NewMessageHandler(svc)


	routers := []common.RouteRegistrar{
		rbac,
		message,
	}

	admin := r.Group("/api/admin/", jwt.AuthMiddleware(), common.PermissionMiddleware(svc))
	for _, h := range routers {
		h.Register(admin)
	}
}
