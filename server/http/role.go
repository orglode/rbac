package http

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/orglode/navigator/model"
)

func roleList(c *gin.Context) {
	req := model.TestAReq{}
	if err := c.ShouldBindQuery(&req); err != nil {
		fmt.Println("asd")
	}
	data, _ := svc.Test(req)
	c.JSON(200, data)
	return
}

func roleAdd(c *gin.Context) {
	req := model.TestAReq{}
	if err := c.ShouldBindQuery(&req); err != nil {
		fmt.Println("asd")
	}
	data, _ := svc.Test(req)
	c.JSON(200, data)
	return
}

func roleModify(c *gin.Context) {
	req := model.TestAReq{}
	if err := c.ShouldBindQuery(&req); err != nil {
		fmt.Println("asd")
	}
	data, _ := svc.Test(req)
	c.JSON(200, data)
	return
}

func roleDel(c *gin.Context) {
	req := model.TestAReq{}
	if err := c.ShouldBindQuery(&req); err != nil {
		fmt.Println("asd")
	}
	data, _ := svc.Test(req)
	c.JSON(200, data)
	return
}

func roleTypeList(c *gin.Context) {
	req := model.TestAReq{}
	if err := c.ShouldBindQuery(&req); err != nil {
		fmt.Println("asd")
	}
	data, _ := svc.Test(req)
	c.JSON(200, data)
	return
}

func roleTypeAdd(c *gin.Context) {
	req := model.TestAReq{}
	if err := c.ShouldBindQuery(&req); err != nil {
		fmt.Println("asd")
	}
	data, _ := svc.Test(req)
	c.JSON(200, data)
	return
}

func roleTypeModify(c *gin.Context) {
	req := model.TestAReq{}
	if err := c.ShouldBindQuery(&req); err != nil {
		fmt.Println("asd")
	}
	data, _ := svc.Test(req)
	c.JSON(200, data)
	return
}

func roleTypeDel(c *gin.Context) {
	req := model.TestAReq{}
	if err := c.ShouldBindQuery(&req); err != nil {
		fmt.Println("asd")
	}
	data, _ := svc.Test(req)
	c.JSON(200, data)
	return
}
