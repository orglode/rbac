package http

import (
	"rbac/api/jwt"
	"rbac/model"

	"github.com/gin-gonic/gin"
)

func registerRbacRoutes(r *gin.Engine) {
	r.POST("/api/admin/login", adminLogin)
	admin := r.Group("/api/admin/", jwt.AuthMiddleware())
	{
		// 后台账号类
		admin.GET("account", adminAccountList)
		admin.POST("account", adminAddAccount)
		admin.PUT("account/:id", adminUpdateAccount)
		admin.DELETE("account/:id", adminDeleteAccount)
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
	account, err := svc.RbacService.LoginAccount(c, req.AccountNumber, req.Password)
	if err != nil {
		standardOutput(c, err)
		return
	}
	responseSuccess(c, account)
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
	if err := c.ShouldBindJSON(&req); err != nil || req.Id <= 0 {
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
