package http

import (
	"github.com/gin-gonic/gin"
	"github.com/orglode/navigator/api"
	"github.com/orglode/navigator/model"
	"net/http"
)

func pageAdd(c *gin.Context) {
	req := model.ModuleRequest{}
	if err := c.BindJSON(&req); err != nil {
		c.JSON(http.StatusOK, model.BaseResponse{
			Code: api.MissingParameter,
		})
		return
	}
	svc.PageAdd(req)
	c.JSON(http.StatusOK, svc.Response)
	return
}

func pageList(c *gin.Context) {
	svc.GetPageAll()
	c.JSON(http.StatusOK, svc.Response)
	return
}

func pageModify(c *gin.Context) {
	req := model.ModuleRequest{}
	if err := c.BindJSON(&req); err != nil {
		c.JSON(http.StatusOK, model.BaseResponse{
			Code: api.MissingParameter,
		})
		return
	}
	svc.PageModify(req)
	c.JSON(http.StatusOK, svc.Response)
	return
}

func pageDel(c *gin.Context) {
	id := c.GetInt64("id")
	operatorUid := c.GetInt64("operator_uid")
	if id <= 0 || operatorUid <= 0 {
		c.JSON(http.StatusOK, model.BaseResponse{
			Code: api.MissingParameter,
		})
		return
	}
	data, _ := svc.PageDel(id, operatorUid)
	c.JSON(http.StatusOK, data)
	return
}
