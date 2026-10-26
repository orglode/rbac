package http

import (
	"rbac/model"
	"rbac/server/common"
	"rbac/service"

	"github.com/gin-gonic/gin"
)

// RbacHandler RBAC 模块：登录、账号、角色、页面、日志等
type RbacHandler struct {
	svc *service.Service
}

func NewRbacHandler(svc *service.Service) *RbacHandler {
	return &RbacHandler{svc: svc}
}

func (h *RbacHandler) RegisterPublic(r *gin.Engine) {
	r.POST("/api/admin/login", h.Login)
}

func (h *RbacHandler) Register(rg *gin.RouterGroup) {
	// 当前用户
	rg.GET("me/menu", h.MyMenu)
	rg.GET("me/buttons", h.MyButtons)

	// 账号
	rg.GET("account", h.AccountList)
	rg.POST("account", h.AccountCreate)
	rg.PUT("account/:id", h.AccountUpdate)
	rg.DELETE("account/:id", h.AccountDelete)
	rg.POST("account/role", h.AccountAssignRole)
	rg.DELETE("account/role", h.AccountRemoveRole)

	// 角色
	rg.GET("role", h.RoleList)
	rg.GET("role/all", h.RoleAll)
	rg.POST("role", h.RoleCreate)
	rg.PUT("role/:id", h.RoleUpdate)
	rg.DELETE("role/:id", h.RoleDelete)
	rg.GET("role/page/:role_id", h.RolePageList)
	rg.POST("role/page", h.RoleAssignPages)

	// 页面
	rg.GET("page", h.PageList)
	rg.GET("page/all", h.PageAll)
	rg.POST("page", h.PageCreate)
	rg.PUT("page/:id", h.PageUpdate)
	rg.DELETE("page/:id", h.PageDelete)

	// 日志
	rg.GET("login/log", h.LoginLogList)
	rg.GET("operation/log", h.OperationLogList)
}

func (h *RbacHandler) Login(c *gin.Context) {
	var req struct {
		AccountNumber string `json:"account_number"`
		Password      string `json:"password"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		common.ResponseError(c, common.InvalidParam)
		return
	}
	meta := model.LoginMeta{
		LoginIp:   c.ClientIP(),
		UserAgent: c.Request.UserAgent(),
	}
	account, err := h.svc.RbacService.LoginAccount(c.Request.Context(), req.AccountNumber, req.Password, meta)
	if err != nil {
		common.StandardOutput(c, err)
		return
	}
	common.ResponseSuccess(c, account)
}

func (h *RbacHandler) MyMenu(c *gin.Context) {
	data, err := h.svc.RbacService.GetUserRolePage(c.Request.Context(), common.GetUid(c))
	if err != nil {
		common.StandardOutput(c, err)
		return
	}
	common.ResponseSuccess(c, data)
}

func (h *RbacHandler) MyButtons(c *gin.Context) {
	data, err := h.svc.RbacService.GetMyButtons(c.Request.Context(), common.GetUid(c))
	if err != nil {
		common.StandardOutput(c, err)
		return
	}
	common.ResponseSuccess(c, data)
}

func (h *RbacHandler) AccountList(c *gin.Context) {
	req := model.AccountQuery{}
	if err := c.ShouldBindQuery(&req); err != nil {
		common.ResponseError(c, common.InvalidParam)
		return
	}
	list, err := h.svc.RbacService.GetAccountList(c, req)
	if err != nil {
		common.StandardOutput(c, err)
		return
	}
	common.ResponseSuccess(c, list)
}

func (h *RbacHandler) AccountCreate(c *gin.Context) {
	var req model.Account
	if err := c.ShouldBindJSON(&req); err != nil {
		common.ResponseError(c, common.InvalidParam)
		return
	}
	account, err := h.svc.RbacService.AddAccount(c, &req)
	if err != nil {
		common.StandardOutput(c, err)
		return
	}
	common.ResponseSuccess(c, account)
}

func (h *RbacHandler) AccountUpdate(c *gin.Context) {
	var req model.Account
	if err := c.ShouldBindJSON(&req); err != nil {
		common.ResponseError(c, common.InvalidParam)
		return
	}
	req.Id = common.GetInt64(c, "id")
	if req.Id <= 0 {
		common.ResponseError(c, common.InvalidParam)
		return
	}
	account, err := h.svc.RbacService.UpdateAccount(c, &req)
	if err != nil {
		common.StandardOutput(c, err)
		return
	}
	common.ResponseSuccess(c, account)
}

func (h *RbacHandler) AccountDelete(c *gin.Context) {
	id := common.GetInt(c, "id")
	if id <= 0 {
		common.ResponseError(c, common.InvalidParam)
		return
	}
	if err := h.svc.RbacService.DeleteAccount(c, int64(id)); err != nil {
		common.StandardOutput(c, err)
		return
	}
	common.ResponseSuccess(c, nil)
}

func (h *RbacHandler) AccountAssignRole(c *gin.Context) {
	ctx := common.WithOperator(c)
	req := &model.AccountRole{}
	if err := c.ShouldBindJSON(req); err != nil {
		common.ResponseError(c, common.InvalidParam)
		return
	}
	if err := h.svc.RbacService.AddAccountRole(ctx, req); err != nil {
		common.StandardOutput(c, err)
		return
	}
	common.ResponseSuccess(c, nil)
}

func (h *RbacHandler) AccountRemoveRole(c *gin.Context) {
	ctx := common.WithOperator(c)
	req := &model.AccountRoleReq{}
	if err := c.ShouldBindJSON(req); err != nil {
		common.ResponseError(c, common.InvalidParam)
		return
	}
	if err := h.svc.RbacService.DelAccountRole(ctx, req.AccountId); err != nil {
		common.StandardOutput(c, err)
		return
	}
	common.ResponseSuccess(c, nil)
}

func (h *RbacHandler) RoleList(c *gin.Context) {
	req := model.Paging{}
	if err := c.ShouldBindQuery(&req); err != nil {
		common.ResponseError(c, common.InvalidParam)
		return
	}
	data, err := h.svc.RbacService.GetRoleList(c.Request.Context(), req)
	if err != nil {
		common.StandardOutput(c, err)
		return
	}
	common.ResponseSuccess(c, data)
}

func (h *RbacHandler) RoleAll(c *gin.Context) {
	data, err := h.svc.RbacService.GetAllRole(c.Request.Context())
	if err != nil {
		common.StandardOutput(c, err)
		return
	}
	common.ResponseSuccess(c, data)
}

func (h *RbacHandler) RoleCreate(c *gin.Context) {
	ctx := common.WithOperator(c)
	req := &model.Role{}
	if err := c.ShouldBindJSON(req); err != nil {
		common.ResponseError(c, common.InvalidParam)
		return
	}
	if err := h.svc.RbacService.AddRole(ctx, req); err != nil {
		common.StandardOutput(c, err)
		return
	}
	common.ResponseSuccess(c, nil)
}

func (h *RbacHandler) RoleUpdate(c *gin.Context) {
	ctx := common.WithOperator(c)
	req := &model.Role{}
	if err := c.ShouldBindJSON(req); err != nil {
		common.ResponseError(c, common.InvalidParam)
		return
	}
	req.Id = common.GetInt64(c, "id")
	if req.Id <= 0 {
		common.ResponseError(c, common.InvalidParam)
		return
	}
	if err := h.svc.RbacService.ModifyRole(ctx, req); err != nil {
		common.StandardOutput(c, err)
		return
	}
	common.ResponseSuccess(c, nil)
}

func (h *RbacHandler) RoleDelete(c *gin.Context) {
	ctx := common.WithOperator(c)
	id := common.GetInt64(c, "id")
	if id <= 0 {
		common.ResponseError(c, common.InvalidParam)
		return
	}
	if err := h.svc.RbacService.DeleteRole(ctx, id); err != nil {
		common.StandardOutput(c, err)
		return
	}
	common.ResponseSuccess(c, nil)
}

func (h *RbacHandler) RolePageList(c *gin.Context) {
	roleId := common.GetInt64(c, "role_id")
	data, err := h.svc.RbacService.GetRolePage(c.Request.Context(), roleId)
	if err != nil {
		common.StandardOutput(c, err)
		return
	}
	common.ResponseSuccess(c, data)
}

func (h *RbacHandler) RoleAssignPages(c *gin.Context) {
	ctx := common.WithOperator(c)
	req := &model.RolePageRequest{}
	if err := c.ShouldBindJSON(req); err != nil {
		common.ResponseError(c, common.InvalidParam)
		return
	}
	if err := h.svc.RbacService.AddPageRole(ctx, req); err != nil {
		common.StandardOutput(c, err)
		return
	}
	common.ResponseSuccess(c, nil)
}

func (h *RbacHandler) PageList(c *gin.Context) {
	req := model.PageQuery{}
	if err := c.ShouldBindQuery(&req); err != nil {
		common.ResponseError(c, common.InvalidParam)
		return
	}
	data, err := h.svc.RbacService.GetPageList(c.Request.Context(), req)
	if err != nil {
		common.StandardOutput(c, err)
		return
	}
	common.ResponseSuccess(c, data)
}

func (h *RbacHandler) PageAll(c *gin.Context) {
	data, err := h.svc.RbacService.GetAllPage(c.Request.Context())
	if err != nil {
		common.StandardOutput(c, err)
		return
	}
	common.ResponseSuccess(c, data)
}

func (h *RbacHandler) PageCreate(c *gin.Context) {
	ctx := common.WithOperator(c)
	req := &model.Page{}
	if err := c.ShouldBindJSON(req); err != nil {
		common.ResponseError(c, common.InvalidParam)
		return
	}
	if err := h.svc.RbacService.AddPage(ctx, req); err != nil {
		common.StandardOutput(c, err)
		return
	}
	common.ResponseSuccess(c, nil)
}

func (h *RbacHandler) PageUpdate(c *gin.Context) {
	ctx := common.WithOperator(c)
	req := &model.Page{}
	if err := c.ShouldBindJSON(req); err != nil {
		common.ResponseError(c, common.InvalidParam)
		return
	}
	req.Id = common.GetInt64(c, "id")
	if req.Id <= 0 {
		common.ResponseError(c, common.InvalidParam)
		return
	}
	if err := h.svc.RbacService.ModifyPage(ctx, req); err != nil {
		common.StandardOutput(c, err)
		return
	}
	common.ResponseSuccess(c, nil)
}

func (h *RbacHandler) PageDelete(c *gin.Context) {
	ctx := common.WithOperator(c)
	id := common.GetInt64(c, "id")
	if id <= 0 {
		common.ResponseError(c, common.InvalidParam)
		return
	}
	if err := h.svc.RbacService.DeletePage(ctx, id); err != nil {
		common.StandardOutput(c, err)
		return
	}
	common.ResponseSuccess(c, nil)
}

func (h *RbacHandler) LoginLogList(c *gin.Context) {
	req := model.LoginLogQuery{}
	if err := c.ShouldBindQuery(&req); err != nil {
		common.ResponseError(c, common.InvalidParam)
		return
	}
	data, err := h.svc.RbacService.GetLoginLogList(c.Request.Context(), req)
	if err != nil {
		common.StandardOutput(c, err)
		return
	}
	common.ResponseSuccess(c, data)
}

func (h *RbacHandler) OperationLogList(c *gin.Context) {
	req := model.OperationLogQuery{}
	if err := c.ShouldBindQuery(&req); err != nil {
		common.ResponseError(c, common.InvalidParam)
		return
	}
	data, err := h.svc.RbacService.GetOperationLogList(c.Request.Context(), req)
	if err != nil {
		common.StandardOutput(c, err)
		return
	}
	common.ResponseSuccess(c, data)
}
