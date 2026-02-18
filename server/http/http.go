package http

import (
	"errors"
	"net/http"
	error2 "rbac/api/error"
	"rbac/conf"
	"rbac/model"
	"rbac/service"
	"strconv"

	"github.com/gin-gonic/gin"
	logger "github.com/orglode/hades/logger_v2"
	"github.com/orglode/hades/trace"
)

var (
	svc *service.Service
)

func Init(s *service.Service, conf *conf.Config) {
	svc = s
	router := gin.New()

	//禁用调式终端颜色
	gin.DisableConsoleColor()

	//判断环境 是否开启debug模式
	if conf.Server.Env == model.EnvProduction {
		gin.SetMode(gin.ReleaseMode)
	}

	// 注册trace中间件
	router.Use(trace.TraceIDMiddleware())

	// gin日志模式
	router.Use(logger.GinLogger())

	// 报错输出中间件
	router.Use(ErrorHandlerMiddleware())

	// gin恢复模式
	router.Use(gin.Recovery())

	// 初始化路由
	initRouter(router)

	// 启动服务
	router.Run(conf.Server.Addr)
}

// ================================================================//

// ErrorHandlerMiddleware 错误处理中间件
func ErrorHandlerMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Next()
		if len(c.Errors) > 0 {
			err := c.Errors.Last().Err
			var appErr *error2.AppError
			if errors.As(err, &appErr) {
				responseError(c, appErr)
			} else {
				systemErr := error2.InternalError
				if conf.ConfGlobal.Server.Env == model.EnvTest {
					systemErr.Err = err
					responseError(c, systemErr)
				} else {
					responseError(c, systemErr)
				}
			}
			return
		}
	}
}

// ResponseSuccess 成功响应
func responseSuccess(c *gin.Context, data interface{}) {
	apiSuccess := error2.Success
	c.JSON(http.StatusOK, model.HttpResponse{
		Code:    apiSuccess.Code,
		Message: apiSuccess.Message,
		Data:    data,
		TraceId: c.GetString("trace_id"),
	})
}

// ResponseError 错误响应
func responseError(c *gin.Context, apiErr *error2.AppError) {
	var code int
	var errMsg string
	if apiErr != nil && apiErr.Err == nil {
		code = error2.InternalError.Code
		errMsg = error2.InternalError.Message
	} else {
		code = error2.InternalError.Code
		errMsg = error2.InternalError.Err.Error()
	}

	c.JSON(http.StatusOK, model.HttpResponse{
		Code:    code,
		Message: errMsg,
		Data:    nil,
		TraceId: c.GetString("trace_id"),
	})
}

func standardOutput(c *gin.Context, err error) {
	_ = c.Error(err)
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
