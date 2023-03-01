package model

type CrmAccountLogin struct {
	UserName string `json:"user_name" form:"user_name"`
	Pd       string `json:"pd" form:"user_name"`
}

type RoleRequest struct {
	BaseRequest
	Role
}

type RoleTypeRequest struct {
	BaseRequest
	RoleType
}

type RoleListRequest struct {
	BaseRequest
}

type RoleListRes struct {
	List  interface{} `json:"list"`
	Total int         `json:"total"`
}

type RoleTypeListRequest struct {
	BaseRequest
}

type RoleTypeListRes struct {
	List  interface{} `json:"list"`
	Total int         `json:"total"`
}
