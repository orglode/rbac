package model

type UserPageList struct {
	BaseRequest
}

type AccountLoginRequest struct {
	Account  string `json:"account"`
	PassWord string `json:"pass_word"`
}
