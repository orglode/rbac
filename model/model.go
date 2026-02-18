package model

type model struct {
}

// Response 通用响应结构
type HttpResponse struct {
	Code    int         `json:"code"`
	Message string      `json:"msg"`
	Data    interface{} `json:"data,omitempty"`
	TraceId string      `json:"trace_id"`
}
type Paging struct {
	Page int `json:"page" form:"page"`
	Size int `json:"size" form:"size"`
}

const (
	EnvProduction = "production"
	EnvTest       = "test"

	StatusSuccess = 1
	StatusFail    = 2
)
