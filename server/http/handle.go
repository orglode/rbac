package http

import (
	"errors"
	"net/http"
	error2 "rbac/api/error"
	"rbac/model"
	"strconv"

	"github.com/gin-gonic/gin"
)

const (
	successCode = 0

	invalidParam = 499

	systemErr = 500
)

// ErrorHandlerMiddleware 错误处理中间件
func ErrorHandlerMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Next()
		if len(c.Errors) > 0 {
			err := c.Errors.Last().Err
			var appErr *error2.AppError
			if errors.As(err, &appErr) {
				responseError(c, appErr.Code, appErr.Message)
			} else {
				responseError(c, systemErr, err.Error())
			}
			return
		}
	}
}

// ResponseSuccess 成功响应
func responseSuccess(c *gin.Context, data interface{}) {
	c.JSON(200, model.HttpResponse{
		Code:    successCode,
		Message: "success",
		Data:    data,
	})
}

func standardOutput(c *gin.Context, err error) {
	_ = c.Error(err)
}

// ResponseError 错误响应
func responseError(c *gin.Context, code int, message ...string) {
	var msg string
	if len(message) > 0 {
		msg = message[0]
	}
	if code == invalidParam {
		msg = "缺少参数"
	}
	c.JSON(http.StatusOK, model.HttpResponse{
		Code:    code,
		Message: msg,
		Data:    nil,
	})
}

// getInt 获取整数参数
func getInt(c *gin.Context, key string) int {
	valStr := c.Param(key)
	val, err := strconv.Atoi(valStr)
	if err != nil {
		return 0
	}
	return val
}

// getInt 获取整数参数
func getInt64(c *gin.Context, key string) int64 {
	valStr := c.Param(key)
	val, err := strconv.ParseInt(valStr, 10, 64)
	if err != nil {
		return 0
	}
	return val
}

// getUserId 获取用户ID
func getUid(c *gin.Context) int64 {
	// 从上下文中获取用户ID，具体实现依赖于你的认证中间件
	val, exists := c.Get("uid")
	if !exists {
		return 0
	}
	userId, ok := val.(int64)
	if !ok {
		return 0
	}
	return userId
}
