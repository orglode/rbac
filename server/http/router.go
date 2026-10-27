package http

import (
	"rbac/api/jwt"
	"rbac/server/common"
	"rbac/service"

	"github.com/gin-gonic/gin"
)

func initRouter(r *gin.Engine, svc *service.Service) {
	rbac := NewRbacHandler(svc)
	message := NewMessageHandler(svc)
	sys := NewSystemHandler()

	// 公开接口
	r.GET("/api/admin/health", sys.Health)
	r.POST("/api/admin/login", rbac.Login)

	admin := r.Group("/api/admin/", jwt.AuthMiddleware(), common.PermissionMiddleware(svc))

	// 当前用户
	admin.GET("me/menu", rbac.MyMenu)
	admin.GET("me/buttons", rbac.MyButtons)

	// 账号
	admin.GET("account", rbac.AccountList)
	admin.POST("account", rbac.AccountCreate)
	admin.PUT("account/:id", rbac.AccountUpdate)
	admin.DELETE("account/:id", rbac.AccountDelete)
	admin.POST("account/role", rbac.AccountAssignRole)
	admin.DELETE("account/role", rbac.AccountRemoveRole)

	// 角色
	admin.GET("role", rbac.RoleList)
	admin.GET("role/all", rbac.RoleAll)
	admin.POST("role", rbac.RoleCreate)
	admin.PUT("role/:id", rbac.RoleUpdate)
	admin.DELETE("role/:id", rbac.RoleDelete)
	admin.GET("role/page/:role_id", rbac.RolePageList)
	admin.POST("role/page", rbac.RoleAssignPages)

	// 页面
	admin.GET("page", rbac.PageList)
	admin.GET("page/all", rbac.PageAll)
	admin.POST("page", rbac.PageCreate)
	admin.PUT("page/:id", rbac.PageUpdate)
	admin.DELETE("page/:id", rbac.PageDelete)

	// 日志
	admin.GET("login/log", rbac.LoginLogList)
	admin.GET("operation/log", rbac.OperationLogList)

	// 站内信
	admin.GET("message/unread-count", message.UnreadCount)
	admin.GET("message/inbox/recent", message.InboxRecent)
	admin.GET("message/inbox", message.Inbox)
	admin.POST("message/read", message.Read)
	admin.GET("message", message.List)
	admin.POST("message", message.Send)
}
