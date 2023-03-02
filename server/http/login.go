package http

import (
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"github.com/orglode/navigator/api"
	"github.com/orglode/navigator/model"
	"net/http"
)

func backgroundList(c *gin.Context) {
	req := model.UserPageList{}
	if err := c.ShouldBindQuery(&req); err != nil {
		c.JSON(http.StatusOK, model.BaseResponse{
			Code: api.MissingParameter,
		})
		return
	}
	svc.GetUserPageShow(req)
	c.JSON(http.StatusOK, svc.Response)
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
