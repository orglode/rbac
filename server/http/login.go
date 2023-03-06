package http

import (
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"github.com/orglode/rbac/api"
	"github.com/orglode/rbac/model"
	"net/http"
)

func backgroundLoginOut(c *gin.Context) {
	req := model.BaseRequest{}
	if err := c.ShouldBindQuery(&req); err != nil {
		c.JSON(http.StatusOK, model.BaseResponse{
			Code: api.MissingParameter,
		})
		return
	}
	sInfo := sessions.Default(c)
	sInfo.Delete("uid")
	sInfo.Save()
	c.JSON(http.StatusOK, model.BaseResponse{
		Code: api.Success,
	})
	return
}

func backgroundLogin(c *gin.Context) {
	req := model.AccountLoginRequest{}
	if err := c.ShouldBindQuery(&req); err != nil {
		c.JSON(http.StatusOK, model.BaseResponse{
			Code: api.MissingParameter,
		})
		return
	}
	if req.Account == "" || req.PassWord == "" {
		c.JSON(http.StatusOK, model.BaseResponse{
			Code: api.MissingParameter,
		})
		return
	}
	sInfo := sessions.Default(c)
	svc.AccountLogin(req, sInfo)
	c.JSON(http.StatusOK, svc.Response)
	return
}

func checkUserLoginSession(c *gin.Context) {
	sInfo := sessions.Default(c)
	userInfo := sInfo.Get("uid")
	if userInfo == nil {
		c.AbortWithStatusJSON(http.StatusOK, model.BaseResponse{
			Code: api.UserNoLogin,
		})
		return
	}
}
