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

func (p *Paging) Normalize() {
	if p.Page <= 0 {
		p.Page = 1
	}
	if p.Size <= 0 {
		p.Size = 20
	}
}

const (
	EnvProduction = "production"
	EnvTest       = "test"

	StatusSuccess = 1
	StatusFail    = 2
)
