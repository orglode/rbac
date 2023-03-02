package http

import (
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

}
