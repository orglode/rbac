package api

type Code struct {
	Code int    `json:"code"`
	Msg  string `json:"msg"`
}

var (
	Success           = &Code{0, "success"}
	MissingParameter  = &Code{499, "缺少参数"}
	SystemErr         = &Code{500, "服务器错误"}
	UserNoTRegister   = &Code{800, "用户未注册"}
	PassWordErr       = &Code{800, "用户密码错误"}
	UserNotFindRole   = &Code{1000, "该用户无角色"}
	RoleNotFindModule = &Code{1000, "该角色未配置权限"}
)
