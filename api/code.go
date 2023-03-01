package api

type Code struct {
	Code int    `json:"code"`
	Msg  string `json:"msg"`
}

var (
	Success          = &Code{0, "success"}
	MissingParameter = &Code{499, "缺少参数"}
	SystemErr        = &Code{500, "服务器错误"}
)
