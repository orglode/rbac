package model

type AccountListRest struct {
	List  interface{} `json:"list"`
	Total int64       `json:"total"`
}

type AccountLoginResult struct {
	Id            int64  `json:"id" `
	Token         string `json:"token"`
	UserName      string `json:"user_name"`
	Mobile        string `json:"mobile"`
	AccountType   int    `json:"account_type"`
	LastLoginTime int64  ` json:"last_login_time"`
	CreateTime    int64  `json:"create_time"`
}

type AccountInfo struct {
	AccountId int64       `json:"account_id"`
	RoleId    int64       `json:"role_id"`
	RoleName  string      `json:"role_name"`
	PageInfo  []*PageInfo `json:"page"`
}

type PageInfo struct {
	PageId    int64       `json:"page_id"`
	PageLevel int         `json:"page_level"`
	PageName  string      `json:"page_name"`
	PagePath  string      `json:"page_path"`
	ChildPage []*PageInfo `json:"child_page"`
}

type AccountRoleReq struct {
	AccountId int64 `json:"account_id"`
	RoleId    int64 `json:"role_id"`
}

type ChinaCityRequestParam struct {
	CityId   int64 `json:"city_id" form:"city_id" binding:"required"`
	CityType int   `json:"city_type" form:"city_type" binding:"required"`
}
