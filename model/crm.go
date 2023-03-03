package model

type CrmAccountLogin struct {
	UserName string `json:"user_name" form:"user_name"`
	Pd       string `json:"pd" form:"user_name"`
}

type CrmUserListRequest struct {
	BaseRequest
	Phone    string `form:"phone"`
	Username string `gorm:"column:username" json:"username" form:"username"` //账号姓名
	Status   int8   `gorm:"column:status" json:"status" form:"status"`       //状态 2启用 1:禁用
}

type CrmUserListInfo struct {
	Users
	RoleName     string `json:"role_name"`
	RoleTypeName string `json:"role_type_name"`
}

type CrmUserListRes struct {
	List  interface{} `json:"list"`
	Total int64       `json:"total"`
}

type CrmUserRequest struct {
	BaseRequest
	Users
	RoleId int64 `form:"role_id" json:"role_id"`
}

type RoleRequest struct {
	BaseRequest
	Role
	RoleModuleId []int64 `json:"role_module_id"`
}

type RoleTypeRequest struct {
	BaseRequest
	RoleType
}

type RoleListRequest struct {
	BaseRequest
	RoleName string `form:"role_name"`
	Status   int8   `form:"status"`
	RoleType int64  `form:"role_type"`
}

type RoleListRes struct {
	List  interface{} `json:"list"`
	Total int64       `json:"total"`
}

type RoleListInfo struct {
	Role
	TypeName string `json:"type_name"`
}

type RoleTypeListRequest struct {
	BaseRequest
}

type RoleTypeListRes struct {
	List  interface{} `json:"list"`
	Total int64       `json:"total"`
}

type ModuleRequest struct {
	BaseRequest
	Module
	RoleId int64 `json:"role_id"`
}

type ModuleAllRes struct {
	List interface{} `json:"list"`
}
