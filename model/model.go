package model

import "github.com/orglode/navigator/api"

type model struct {
}

type BaseResponse struct {
	*api.Code
	Data interface{} `json:"data"`
}

type BaseRequest struct {
	OperatorUid int64 `json:"operator_uid" form:"operator_uid" binding:"required" `
	Page        int   `json:"page" form:"page"`
	Size        int   `json:"size" form:"size"`
}

const (
	StatusSuccess = 2
	StatusFail    = 1
)
