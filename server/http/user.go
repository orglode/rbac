package http

import (
	"github.com/gin-gonic/gin"
	"github.com/orglode/rbac/api"
	"github.com/orglode/rbac/model"
	"net/http"
)

func userList(c *gin.Context) {
	req := model.CrmUserListRequest{}
	if err := c.ShouldBindQuery(&req); err != nil {
		c.JSON(http.StatusOK, model.BaseResponse{
			Code: api.MissingParameter,
		})
		return
	}
	svc.CrmUserListInfo(req)
	c.JSON(http.StatusOK, svc.Response)
	return
}

func userAdd(c *gin.Context) {
	req := model.CrmUserRequest{}
	if err := c.BindJSON(&req); err != nil {
		c.JSON(http.StatusOK, model.BaseResponse{
			Code: api.MissingParameter,
		})
		return
	}
	svc.UserAdd(req)
	c.JSON(http.StatusOK, svc.Response)
	return
}

func userModify(c *gin.Context) {
	req := model.CrmUserRequest{}
	if err := c.BindJSON(&req); err != nil {
		c.JSON(http.StatusOK, model.BaseResponse{
			Code: api.MissingParameter,
		})
		return
	}
	svc.UserModify(req)
	c.JSON(http.StatusOK, svc.Response)
	return
}

func userDel(c *gin.Context) {
	id := c.GetInt64("id")
	operatorUid := c.GetInt64("operator_uid")
	if id <= 0 || operatorUid <= 0 {
		c.JSON(http.StatusOK, model.BaseResponse{
			Code: api.MissingParameter,
		})
		return
	}
	svc.UserDel(id, operatorUid)
	c.JSON(http.StatusOK, svc.Response)
	return
}
