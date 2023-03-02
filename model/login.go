package model

type UserSessionInfo struct {
	Account string `json:"account"`
	Phone   string `json:"phone"`
	UserId  int64  `json:"user_id"`
	RoleId  int64  `json:"role_id"`
}

type UserPageList struct {
	BaseRequest
}

type AccountLoginRequest struct {
	Account  string `form:"account"`
	PassWord string `form:"pass_word"`
}
