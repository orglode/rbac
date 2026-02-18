package model

type PageList struct {
	Total int64       `json:"total"`
	List  interface{} `json:"list"`
}

type PageListInfo struct {
	Page
	ParentName string `json:"parent_name"`
}

type RolePageRequest struct {
	RoleId  int64  `json:"role_id"`
	PageIds string `json:"page_ids"`
}

type RolePageInfo struct {
	RoleId   int64       `json:"role_id"`
	RoleName string      `json:"role_name"`
	PageInfo []*PageInfo `json:"page"`
}
