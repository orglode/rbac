package error

import (
	"fmt"
)

var (
	PleaseLogin  = &AppError{Code: 490, Message: "please login"}
	InvalidParam = &AppError{Code: 499, Message: "缺少参数"}

	// RBAC
	RbacAccountAlreadyExists = &AppError{Code: 600, Message: "账号已存在"}
	RbacRoleNotFound         = &AppError{Code: 601, Message: "角色不存在"}
	RbacInvalidPermission    = &AppError{Code: 602, Message: "权限无效"}
	RbacMobileExists         = &AppError{Code: 603, Message: "手机号已存在"}
	LoginPassWord            = &AppError{Code: 604, Message: "账号密码错误"}
)

type AppError struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
	Err     error  `json:"error"`
}

func (e *AppError) Error() string {
	return fmt.Sprintf("code=%d, msg=%s", e.Code, e.Message)
}

func (e *AppError) Unwrap() error {
	return e.Err
}
