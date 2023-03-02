package model

type CrmAccountLogin struct {
	UserName string `json:"user_name" form:"user_name"`
	Pd       string `json:"pd" form:"user_name"`
}

type CrmUserListRequest struct {
	BaseRequest
}

type CrmUserListRes struct {
	List  interface{} `json:"list"`
	Total int64       `json:"total"`
}

type CrmUserRequest struct {
	BaseRequest
	Users
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
	Total int64       `json:"total"`
}

type RoleTypeListRequest struct {
	BaseRequest
}

type RoleTypeListRes struct {
	List  interface{} `json:"list"`
	Total int64       `json:"total"`
}
