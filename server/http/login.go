package http

import (
	"github.com/gin-gonic/gin"
	"github.com/orglode/navigator/api"
	"github.com/orglode/navigator/model"
	"net/http"
)

func backgroundLogin(c *gin.Context) {
	req := model.CrmAccountLogin{}
	if err := c.ShouldBindQuery(&req); err != nil {
		c.JSON(http.StatusOK, model.BaseResponse{
			Code: api.MissingParameter,
		})
		return
	}
	//data, _ := svc.Test(req)
	//c.JSON(200, data)
	return
}

func backgroundList(c *gin.Context) {
	req := model.CrmAccountLogin{}
	if err := c.ShouldBindQuery(&req); err != nil {
		c.JSON(http.StatusOK, model.BaseResponse{
			Code: api.MissingParameter,
		})
		return
	}
	return
}
