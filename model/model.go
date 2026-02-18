package model

import (
	"github.com/golang-jwt/jwt/v4"
)

type model struct {
}

// Response 通用响应结构
type HttpResponse struct {
	Code    int         `json:"code"`
	Message string      `json:"msg"`
	Data    interface{} `json:"data,omitempty"`
}
type BaseCrmResponse struct {
	Total int64       `json:"total"`
	List  interface{} `json:"list"`
}

type Paging struct {
	Page int `json:"page" form:"page"`
	Size int `json:"size" form:"size"`
}

const (
	EnvProduction = "production"

	StatusSuccess = 1
	StatusFail    = 2
)

type MyClaims struct {
	UserId int64 `json:"user_id"`
	jwt.RegisteredClaims
}
