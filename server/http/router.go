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
	message := NewMessageHandler(svc)
	rbac.RegisterPublic(r)

	admin := r.Group("/api/admin/", jwt.AuthMiddleware(), common.PermissionMiddleware(svc))
	for _, h := range []common.RouteRegistrar{rbac, message} {
		h.Register(admin)
	}
}
