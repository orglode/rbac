package common

import (
	"strings"

	error2 "rbac/api/error"
	"rbac/service"

	"github.com/gin-gonic/gin"
)

// 无需按钮鉴权的接口（仍需登录）
var permissionWhitelist = map[string]struct{}{
	"GET /api/admin/me/menu":              {},
	"GET /api/admin/me/buttons":           {},
	"GET /api/admin/message/unread-count": {},
	"GET /api/admin/message/inbox":        {},
	"GET /api/admin/message/inbox/recent": {},
	"POST /api/admin/message/read":        {},
}

// PermissionMiddleware 按按钮绑定的 API（method + FullPath）校验权限
func PermissionMiddleware(svc *service.Service) gin.HandlerFunc {
	return func(c *gin.Context) {
		uid := GetUid(c)
		if uid <= 0 {
			ResponseError(c, error2.RbacInvalidPermission)
			c.Abort()
			return
		}

		method := strings.ToUpper(c.Request.Method)
		fullPath := c.FullPath()
		if fullPath == "" {
			c.Next()
			return
		}

		key := method + " " + fullPath
		if _, ok := permissionWhitelist[key]; ok {
			c.Next()
			return
		}

		if err := svc.RbacService.CheckApiPermission(c.Request.Context(), uid, method, fullPath); err != nil {
			StandardOutput(c, err)
			c.Abort()
			return
		}
		c.Next()
	}
}
