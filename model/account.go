package model

type AccountListRest struct {
	List  interface{} `json:"list"`
	Total int64       `json:"total"`
}

type AccountLoginResult struct {
	Id            int64  `json:"id"`
	Token         string `json:"token"`
	UserName      string `json:"user_name"`
	Mobile        string `json:"mobile"`
	AccountNumber string `json:"account_number"`
	Status        int    `json:"status"`
	LastLoginTime int64  `json:"last_login_time"`
	CreateTime    int64  `json:"create_time"`
}

type AccountQuery struct {
	Paging
	UserName      string `json:"user_name" form:"user_name"`
	Mobile        string `json:"mobile" form:"mobile"`
	AccountNumber string `json:"account_number" form:"account_number"`
	Status        int    `json:"status" form:"status"` // 0 不限；1启用 2禁用
}

type PageQuery struct {
	Paging
	PageName string `json:"page_name" form:"page_name"` // 名称模糊
	Api      string `json:"api" form:"api"`             // api_method / api_path 模糊
}

type LoginLogQuery struct {
	Paging
	AccountId     int64  `json:"account_id" form:"account_id"`
	AccountNumber string `json:"account_number" form:"account_number"`
	Status        int    `json:"status" form:"status"`
}

type LoginMeta struct {
	LoginIp   string
	UserAgent string
}

type AccountInfo struct {
	AccountId int64       `json:"account_id"`
	RoleId    int64       `json:"role_id"`
	RoleName  string      `json:"role_name"`
	PageInfo  []*PageInfo `json:"page"`
}

type PageInfo struct {
	PageId     int64       `json:"page_id"`
	PageLevel  int         `json:"page_level"`
	PageName   string      `json:"page_name"`
	PagePath   string      `json:"page_path"`
	PageType   int         `json:"page_type"`
	ButtonCode string      `json:"button_code,omitempty"`
	ApiMethod  string      `json:"api_method,omitempty"`
	ApiPath    string      `json:"api_path,omitempty"`
	ChildPage  []*PageInfo `json:"child_page"`
	Buttons    []*PageInfo `json:"buttons,omitempty"`
}

type AccountRoleReq struct {
	AccountId int64 `json:"account_id"`
	RoleId    int64 `json:"role_id"`
}

type ChinaCityRequestParam struct {
	CityId   int64 `json:"city_id" form:"city_id" binding:"required"`
	CityType int   `json:"city_type" form:"city_type" binding:"required"`
}
