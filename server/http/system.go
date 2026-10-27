package http

import (
	"time"

	"rbac/server/common"

	"github.com/gin-gonic/gin"
)

// SystemHandler 系统级公开接口
type SystemHandler struct{}

func NewSystemHandler() *SystemHandler {
	return &SystemHandler{}
}

func (h *SystemHandler) Health(c *gin.Context) {
	common.ResponseSuccess(c, gin.H{
		"time": time.Now().Unix(),
		"ip":   c.ClientIP(),
	})
}
