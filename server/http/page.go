package http

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/orglode/navigator/model"
)

func pageAdd(c *gin.Context) {
	req := model.TestAReq{}
	if err := c.ShouldBindQuery(&req); err != nil {
		fmt.Println("asd")
	}
	data, _ := svc.Test(req)
	c.JSON(200, data)
	return
}

func pageList(c *gin.Context) {
	req := model.TestAReq{}
	if err := c.ShouldBindQuery(&req); err != nil {
		fmt.Println("asd")
	}
	data, _ := svc.Test(req)
	c.JSON(200, data)
	return
}

func pageModify(c *gin.Context) {
	req := model.TestAReq{}
	if err := c.ShouldBindQuery(&req); err != nil {
		fmt.Println("asd")
	}
	data, _ := svc.Test(req)
	c.JSON(200, data)
	return
}

func pageDel(c *gin.Context) {
	req := model.TestAReq{}
	if err := c.ShouldBindQuery(&req); err != nil {
		fmt.Println("asd")
	}
	data, _ := svc.Test(req)
	c.JSON(200, data)
	return
}
