package model

const (
	AccountTable         = "account"
	RbacRoleTable        = "role"
	RbacPageTable        = "page"
	RbacRolePageTable    = "role_page"
	RbacAccountRoleTable = "account_role"
)

// Account 账户
type Account struct {
	Id int64 `gorm:"column:id;type:int(11);primary_key;not null;auto_increment" json:"id" form:"id"`
	// UserName 账号姓名
	UserName string `gorm:"column:user_name;type:varchar(50);not null" json:"user_name" form:"user_name"`
	// Mobile 电话
	Mobile string `gorm:"column:mobile;type:varchar(50);unique_index:mobile;not null" json:"mobile" form:"mobile"`
	// AccountNumber 登录账号
	AccountNumber string `gorm:"column:account_number;type:varchar(50);unique_index:account_number;not null" json:"account_number" form:"account_number"`
	// Password 密码
	Password string `gorm:"column:password;type:varchar(50);not null" json:"password" form:"password"`
	// UserType 账户类型 // 1 管理员
	AccountType int `gorm:"column:account_type;type:tinyint(4);not null" json:"account_type" form:"account_type"`
	// LastLoginTime 最后一次登录时间
	LastLoginTime int64 `gorm:"column:last_login_time;type:int(11);not null" json:"last_login_time" form:"last_login_time"`
	// CreateTime 创建时间
	CreateTime int64 `gorm:"column:create_time;type:int(11);not null" json:"create_time" form:"create_time"`
	// UpdateTime 更新时间
	UpdateTime int64 `gorm:"column:update_time;type:int(11);not null" json:"update_time" form:"update_time"`
}

// Role 用户角色表
type Role struct {
	Id int64 `gorm:"column:id;type:int(11);primary_key;not null;auto_increment" json:"id" form:"id"`
	// RoleName 账号姓名
	RoleName string `gorm:"column:role_name;type:varchar(50);not null" json:"role_name" form:"role_name"`
	// CreateTime 创建时间
	CreateTime int64 `gorm:"column:create_time;type:int(11);not null" json:"create_time" form:"create_time"`
}

type AccountRole struct {
	Id int64 `gorm:"column:id;type:int(11);primary_key;not null;auto_increment" json:"id" form:"id"`
	// AccountID 账号ID
	AccountId int64 `gorm:"column:account_id;type:int(11);not null" json:"account_id" form:"account_id"`
	// RoleID 角色ID
	RoleId int64 `gorm:"column:role_id;type:int(11);not null" json:"role_id" form:"role_id"`
	// CreateTime 创建时间
	CreateTime int64 `gorm:"column:create_time;type:int(11);not null" json:"create_time" form:"create_time"`
}

// RolePage 页面菜单栏
type RolePage struct {
	Id int `gorm:"column:id;type:int(11);primary_key;not null;auto_increment" json:"id" form:"id"`
	// RoleID 账号ID
	RoleId int64 `gorm:"column:role_id;type:int(11);not null" json:"role_id" form:"role_id"`
	// PageID 账号ID
	PageId int64 `gorm:"column:page_id;type:int(11);not null" json:"page_id" form:"page_id"`
}
type Page struct {
	Id int64 `gorm:"column:id;type:int(11);primary_key;not null;auto_increment" json:"id" form:"id"`
	// PageName 菜单名称
	PageName string `gorm:"column:page_name;type:varchar(50);not null" json:"page_name" form:"page_name"`
	// Url 路由
	Url string `gorm:"column:url;type:varchar(255);not null" json:"url" form:"url"`
	// 父级id
	ParentId int64 `gorm:"column:parent_id;type:int(11);not null" json:"parent_id" form:"parent_id"`
	// 页面等级
	PageLevel int `json:"page_level" form:"page_level"`
	// CreateTime 创建时间
	CreateTime int64 `gorm:"column:create_time;type:int(11);not null" json:"create_time" form:"create_time"`
}
