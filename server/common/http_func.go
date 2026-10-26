package common

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	error2 "rbac/api/error"
	"rbac/conf"
	"rbac/model"
	"strconv"

	"github.com/gin-gonic/gin"
)

// InvalidParam 参数错误
var InvalidParam = error2.MissingParameter

// ErrorHandlerMiddleware 错误处理中间件
func ErrorHandlerMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Next()
		if len(c.Errors) > 0 {
			err := c.Errors.Last().Err
			var appErr *error2.AppError
			if errors.As(err, &appErr) {
				ResponseError(c, appErr)
			} else {
				ResponseError(c, &error2.AppError{
					Code:    error2.InternalError.Code,
					Message: error2.InternalError.Message,
					Err:     err,
				})
			}
			return
		}
	}
}

// ResponseSuccess 成功响应
func ResponseSuccess(c *gin.Context, data any) {
	apiSuccess := error2.Success
	c.JSON(http.StatusOK, model.HttpResponse{
		Code:    apiSuccess.Code,
		Message: apiSuccess.Message,
		Data:    data,
		TraceId: c.GetString("trace_id"),
	})
}

// ResponseError 错误响应
func ResponseError(c *gin.Context, apiErr *error2.AppError) {
	if apiErr == nil {
		apiErr = error2.InternalError
	}
	errMsg := apiErr.Message
	if apiErr.Err != nil && conf.ConfGlobal.Server.Env == model.EnvTest {
		errMsg = fmt.Sprintf("%s: %v", errMsg, apiErr.Err)
	}

	c.JSON(http.StatusOK, model.HttpResponse{
		Code:    apiErr.Code,
		Message: errMsg,
		Data:    nil,
		TraceId: c.GetString("trace_id"),
	})
}

// StandardOutput 将错误交给 ErrorHandlerMiddleware
func StandardOutput(c *gin.Context, err error) {
	_ = c.Error(err)
}

// GetInt 路径参数转 int
func GetInt(c *gin.Context, key string) int {
	valStr := c.Param(key)
	val, err := strconv.Atoi(valStr)
	if err != nil {
		return 0
	}
	return val
}

// GetInt64 路径参数转 int64
func GetInt64(c *gin.Context, key string) int64 {
	valStr := c.Param(key)
	val, err := strconv.ParseInt(valStr, 10, 64)
	if err != nil {
		return 0
	}
	return val
}

// WithOperator 将操作人写入 context
func WithOperator(c *gin.Context) context.Context {
	return context.WithValue(c.Request.Context(), model.CtxOperatorID, GetUid(c))
}

// GetUid 从 gin context 取登录用户 ID
func GetUid(c *gin.Context) int64 {
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
