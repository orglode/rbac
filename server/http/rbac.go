package http

import (
	"strings"

	error2 "rbac/api/error"
	"rbac/api/jwt"
	"rbac/model"

	"github.com/gin-gonic/gin"
)

// 无需按钮鉴权的接口（仍需登录，login 本身不在此组）
var permissionWhitelist = map[string]struct{}{
	"GET /api/admin/heath":      {},
	"GET /api/admin/me/menu":    {},
	"GET /api/admin/me/buttons": {},
}

func registerRbacRoutes(r *gin.Engine) {
	r.POST("/api/admin/login", adminLogin)
	admin := r.Group("/api/admin/", jwt.AuthMiddleware(), PermissionMiddleware())
	{
		// 当前登录用户菜单/按钮（白名单，供前端渲染）
		admin.GET("me/menu", adminMyMenu)
		admin.GET("me/buttons", adminMyButtons)

		// 后台账号类
		admin.GET("account", adminAccountList)
		admin.POST("account", adminAddAccount)
		admin.PUT("account/:id", adminUpdateAccount)
		admin.DELETE("account/:id", adminDeleteAccount)
		// 登录日志
		admin.GET("login/log", adminLoginLogList)
		// 角色类
		admin.GET("role", adminRoleList)
		admin.GET("role/all", adminAllRole)
		admin.POST("role", adminAddRole)
		admin.PUT("role/:id", adminModifyRole)
		admin.DELETE("role/:id", adminDelRole)
		// 页面
		admin.GET("page", adminPageList)
		admin.GET("page/all", adminPage)
		admin.POST("page", adminAddPage)
		admin.PUT("page/:id", adminModifyPage)
		admin.DELETE("page/:id", adminDelPage)

		// 角色页面关联
		admin.GET("role/page/:role_id", adminRolePageList)
		admin.POST("role/page", adminAddRolePage)

		// 账户角色关联
		admin.POST("account/role", adminAddAccountRole)
		admin.DELETE("account/role", adminDelAccountRole)
	}
}

// PermissionMiddleware 按按钮绑定的 API（method + FullPath）校验权限
func PermissionMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		uid := getUid(c)
		if uid <= 0 {
			responseError(c, error2.RbacInvalidPermission)
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
			standardOutput(c, err)
			c.Abort()
			return
		}
		c.Next()
	}
}



// 删除用户
func adminAccountInfo(c *gin.Context) {
	id := getInt(c, "id")
	if id <= 0 {
		responseError(c, invalidParam)
		return
	}
	data, err := svc.RbacService.GetAccountById(c, int64(id))
	if err != nil {
		standardOutput(c, err)
		return
	}
	responseSuccess(c, data)
}

// 登录
func adminLogin(c *gin.Context) {
	var req struct {
		AccountNumber string `json:"account_number"`
		Password      string `json:"password"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		responseError(c, invalidParam)
		return
	}
	meta := model.LoginMeta{
		LoginIp:   c.ClientIP(),
		UserAgent: c.Request.UserAgent(),
	}
	account, err := svc.RbacService.LoginAccount(c.Request.Context(), req.AccountNumber, req.Password, meta)
	if err != nil {
		standardOutput(c, err)
		return
	}
	responseSuccess(c, account)
}

// 登录日志列表
func adminLoginLogList(c *gin.Context) {
	req := model.LoginLogQuery{}
	if err := c.ShouldBindQuery(&req); err != nil {
		responseError(c, invalidParam)
		return
	}
	data, err := svc.RbacService.GetLoginLogList(c.Request.Context(), req)
	if err != nil {
		standardOutput(c, err)
		return
	}
	responseSuccess(c, data)
}

// 用户列表
func adminAccountList(c *gin.Context) {
	req := model.Paging{}
	if err := c.ShouldBindQuery(&req); err != nil {
		responseError(c, invalidParam)
		return
	}
	list, err := svc.RbacService.GetAccountList(c, req)
	if err != nil {
		standardOutput(c, err)
		return
	}
	responseSuccess(c, list)
}

// 添加用户
func adminAddAccount(c *gin.Context) {
	var req model.Account
	if err := c.ShouldBindJSON(&req); err != nil {
		responseError(c, invalidParam)
		return
	}
	account, err := svc.RbacService.AddAccount(c, &req)
	if err != nil {
		standardOutput(c, err)
		return
	}
	responseSuccess(c, account)
}

// 更新用户
func adminUpdateAccount(c *gin.Context) {
	var req model.Account
	if err := c.ShouldBindJSON(&req); err != nil {
		responseError(c, invalidParam)
		return
	}
	req.Id = getInt64(c, "id")
	if req.Id <= 0 {
		responseError(c, invalidParam)
		return
	}
	account, err := svc.RbacService.UpdateAccount(c, &req)
	if err != nil {
		standardOutput(c, err)
		return
	}
	responseSuccess(c, account)
}

// 当前用户菜单树（含二级菜单下的 buttons）
func adminMyMenu(c *gin.Context) {
	uid := getUid(c)
	data, err := svc.RbacService.GetUserRolePage(c.Request.Context(), uid)
	if err != nil {
		standardOutput(c, err)
		return
	}
	responseSuccess(c, data)
}

// 当前用户按钮权限列表
func adminMyButtons(c *gin.Context) {
	uid := getUid(c)
	data, err := svc.RbacService.GetMyButtons(c.Request.Context(), uid)
	if err != nil {
		standardOutput(c, err)
		return
	}
	responseSuccess(c, data)
}

// 删除用户
func adminDeleteAccount(c *gin.Context) {
	id := getInt(c, "id")
	if id <= 0 {
		responseError(c, invalidParam)
		return
	}
	err := svc.RbacService.DeleteAccount(c, int64(id))
	if err != nil {
		standardOutput(c, err)
		return
	}
	responseSuccess(c, nil)
}

func adminRoleList(c *gin.Context) {
	ctx := c.Request.Context()
	req := model.Paging{}
	if err := c.ShouldBindQuery(&req); err != nil {
		responseError(c, invalidParam)
		return
	}
	data, err := svc.RbacService.GetRoleList(ctx, req)
	if err != nil {
		standardOutput(c, err)
		return
	}
	responseSuccess(c, data)
}

func adminAllRole(c *gin.Context) {
	ctx := c.Request.Context()
	data, err := svc.RbacService.GetAllRole(ctx)
	if err != nil {
		standardOutput(c, err)
		return
	}
	responseSuccess(c, data)
}

func adminAddRole(c *gin.Context) {
	ctx := c.Request.Context()
	req := &model.Role{}
	if err := c.ShouldBindJSON(req); err != nil {
		responseError(c, invalidParam)
		return
	}
	err := svc.RbacService.AddRole(ctx, req)
	if err != nil {
		standardOutput(c, err)
		return
	}
	responseSuccess(c, nil)
}

func adminModifyRole(c *gin.Context) {
	ctx := c.Request.Context()
	req := &model.Role{}
	if err := c.ShouldBindJSON(req); err != nil {
		responseError(c, invalidParam)
		return
	}
	req.Id = getInt64(c, "id")
	if req.Id <= 0 {
		responseError(c, invalidParam)
		return
	}
	err := svc.RbacService.ModifyRole(ctx, req)
	if err != nil {
		standardOutput(c, err)
		return
	}
	responseSuccess(c, nil)
}

func adminDelRole(c *gin.Context) {
	ctx := c.Request.Context()
	id := getInt64(c, "id")
	if id <= 0 {
		responseError(c, invalidParam)
		return
	}
	err := svc.RbacService.DeleteRole(ctx, id)
	if err != nil {
		standardOutput(c, err)
		return
	}
	responseSuccess(c, nil)
}

func adminPageList(c *gin.Context) {
	ctx := c.Request.Context()
	req := model.Paging{}
	if err := c.ShouldBindQuery(&req); err != nil {
		responseError(c, invalidParam)
		return
	}
	data, err := svc.RbacService.GetPageList(ctx, req)
	if err != nil {
		standardOutput(c, err)
		return
	}
	responseSuccess(c, data)

}
func adminPage(c *gin.Context) {
	ctx := c.Request.Context()
	data, err := svc.RbacService.GetAllPage(ctx)
	if err != nil {
		standardOutput(c, err)
		return
	}
	responseSuccess(c, data)
}

func adminAddPage(c *gin.Context) {
	ctx := c.Request.Context()
	req := &model.Page{}
	if err := c.ShouldBindJSON(req); err != nil {
		responseError(c, invalidParam)
		return
	}
	err := svc.RbacService.AddPage(ctx, req)
	if err != nil {
		standardOutput(c, err)
		return
	}
	responseSuccess(c, nil)
}

func adminModifyPage(c *gin.Context) {
	ctx := c.Request.Context()
	req := &model.Page{}
	if err := c.ShouldBindJSON(req); err != nil {
		responseError(c, invalidParam)
		return
	}
	req.Id = getInt64(c, "id")
	if req.Id <= 0 {
		responseError(c, invalidParam)
		return
	}
	err := svc.RbacService.ModifyPage(ctx, req)
	if err != nil {
		standardOutput(c, err)
		return
	}
	responseSuccess(c, nil)
}

func adminDelPage(c *gin.Context) {
	ctx := c.Request.Context()
	id := getInt64(c, "id")
	if id <= 0 {
		responseError(c, invalidParam)
		return
	}
	err := svc.RbacService.DeletePage(ctx, id)
	if err != nil {
		standardOutput(c, err)
		return
	}
	responseSuccess(c, nil)
}

func adminRolePageList(c *gin.Context) {
	ctx := c.Request.Context()
	req := model.Paging{}
	if err := c.ShouldBindQuery(&req); err != nil {
		responseError(c, invalidParam)
		return
	}
	roleId := getInt64(c, "role_id")
	data, err := svc.RbacService.GetRolePage(ctx, roleId)
	if err != nil {
		standardOutput(c, err)
		return
	}
	responseSuccess(c, data)
}

func adminAddRolePage(c *gin.Context) {
	ctx := c.Request.Context()
	req := &model.RolePageRequest{}
	if err := c.ShouldBindJSON(req); err != nil {
		responseError(c, invalidParam)
		return
	}
	err := svc.RbacService.AddPageRole(ctx, req)
	if err != nil {
		standardOutput(c, err)
		return
	}
	responseSuccess(c, nil)
}

func adminAddAccountRole(c *gin.Context) {
	ctx := c.Request.Context()
	req := &model.AccountRole{}
	if err := c.ShouldBindJSON(req); err != nil {
		responseError(c, invalidParam)
		return
	}
	err := svc.RbacService.AddAccountRole(ctx, req)
	if err != nil {
		standardOutput(c, err)
		return
	}
	responseSuccess(c, nil)
}

func adminDelAccountRole(c *gin.Context) {
	ctx := c.Request.Context()
	req := &model.AccountRoleReq{}
	if err := c.ShouldBindJSON(req); err != nil {
		responseError(c, invalidParam)
		return
	}
	err := svc.RbacService.DelAccountRole(ctx, req.AccountId)
	if err != nil {
		standardOutput(c, err)
		return
	}
	responseSuccess(c, nil)
}
