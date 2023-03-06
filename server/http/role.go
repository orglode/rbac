package http

import (
	"github.com/gin-gonic/gin"
	"github.com/orglode/rbac/api"
	"github.com/orglode/rbac/model"
	"net/http"
)

func roleList(c *gin.Context) {
	req := model.RoleListRequest{}
	if err := c.ShouldBindQuery(&req); err != nil {
		c.JSON(http.StatusOK, model.BaseResponse{
			Code: api.MissingParameter,
		})
		return
	}
	svc.CrmRoleListInfo(req)
	c.JSON(http.StatusOK, svc.Response)
	return
}

func roleAll(c *gin.Context) {
	data, _ := svc.GetRoleAll()
	c.JSON(http.StatusOK, data)
	return
}

func roleAdd(c *gin.Context) {
	req := model.RoleRequest{}
	if err := c.BindJSON(&req); err != nil {
		c.JSON(http.StatusOK, model.BaseResponse{
			Code: api.MissingParameter,
		})
		return
	}
	data, _ := svc.RoleAdd(req)
	c.JSON(http.StatusOK, data)
	return
}

func roleModify(c *gin.Context) {
	req := model.RoleRequest{}
	if err := c.ShouldBindQuery(&req); err != nil {
		c.JSON(http.StatusOK, model.BaseResponse{
			Code: api.MissingParameter,
		})
		return
	}
	if req.Id <= 0 {
		return
	}
	data, _ := svc.RoleModify(req)
	c.JSON(http.StatusOK, data)
	return
}

func roleDel(c *gin.Context) {
	id := c.GetInt64("id")
	operatorUid := c.GetInt64("operator_uid")
	if id <= 0 || operatorUid <= 0 {
		c.JSON(http.StatusOK, model.BaseResponse{
			Code: api.MissingParameter,
		})
		return
	}
	data, _ := svc.RoleDel(id, operatorUid)
	c.JSON(http.StatusOK, data)
	return
}

func roleTypeList(c *gin.Context) {
	req := model.RoleTypeListRequest{}
	if err := c.ShouldBindQuery(&req); err != nil {
		c.JSON(http.StatusOK, model.BaseResponse{
			Code: api.MissingParameter,
		})
		return
	}
	svc.CrmRoleTypeListInfo(req)
	c.JSON(http.StatusOK, svc.Response)
	return
}

func roleTypeAll(c *gin.Context) {
	data, _ := svc.GetRoleTypeAll()
	c.JSON(http.StatusOK, data)
	return
}

func roleTypeAdd(c *gin.Context) {
	req := model.RoleTypeRequest{}
	if err := c.BindJSON(&req); err != nil {
		c.JSON(http.StatusOK, model.BaseResponse{
			Code: api.MissingParameter,
		})
		return
	}
	data, _ := svc.RoleTypeAdd(req)
	c.JSON(http.StatusOK, data)
	return
}

func roleTypeModify(c *gin.Context) {
	req := model.RoleTypeRequest{}
	if err := c.BindJSON(&req); err != nil {
		c.JSON(http.StatusOK, model.BaseResponse{
			Code: api.MissingParameter,
		})
		return
	}
	if req.Id <= 0 {
		c.JSON(http.StatusOK, model.BaseResponse{
			Code: api.MissingParameter,
		})
		return
	}
	data, _ := svc.RoleTypeModify(req)
	c.JSON(http.StatusOK, data)
	return
}

func roleTypeDel(c *gin.Context) {
	id := c.GetInt64("id")
	operatorUid := c.GetInt64("operator_uid")
	if id <= 0 || operatorUid <= 0 {
		c.JSON(http.StatusOK, model.BaseResponse{
			Code: api.MissingParameter,
		})
		return
	}
	data, _ := svc.DelRoleTypeInfo(id, operatorUid)
	c.JSON(http.StatusOK, data)
	return
}
