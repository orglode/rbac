package model

const (
	AccountTable         = "account"
	LoginLogTable        = "login_log"
	RbacRoleTable        = "role"
	RbacPageTable        = "page"
	RbacRolePageTable    = "role_page"
	RbacAccountRoleTable = "account_role"
)

// 账号状态
const (
	AccountStatusEnable  = 1 // 启用
	AccountStatusDisable = 2 // 禁用
)

// 登录结果
const (
	LoginStatusSuccess = 1 // 成功
	LoginStatusFail    = 2 // 失败
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
	// AccountType 账户类型 // 1 超级管理员
	AccountType int `gorm:"column:account_type;type:tinyint(4);not null" json:"account_type" form:"account_type"`
	// Status 账号状态 1启用 2禁用
	Status int `gorm:"column:status;type:tinyint(4);not null;default:1" json:"status" form:"status"`
	// LastLoginTime 最后一次登录时间
	LastLoginTime int64 `gorm:"column:last_login_time;type:int(11);not null" json:"last_login_time" form:"last_login_time"`
	// CreateTime 创建时间
	CreateTime int64 `gorm:"column:create_time;type:int(11);not null" json:"create_time" form:"create_time"`
	// UpdateTime 更新时间
	UpdateTime int64 `gorm:"column:update_time;type:int(11);not null" json:"update_time" form:"update_time"`
}

// LoginLog 登录日志
type LoginLog struct {
	Id            int64  `gorm:"column:id;type:int(11);primary_key;not null;auto_increment" json:"id"`
	AccountId     int64  `gorm:"column:account_id;type:int(11);not null;default:0" json:"account_id"`
	AccountNumber string `gorm:"column:account_number;type:varchar(50);not null;default:''" json:"account_number"`
	LoginIp       string `gorm:"column:login_ip;type:varchar(64);not null;default:''" json:"login_ip"`
	UserAgent     string `gorm:"column:user_agent;type:varchar(255);not null;default:''" json:"user_agent"`
	Status        int    `gorm:"column:status;type:tinyint(4);not null;default:0" json:"status"` // 1成功 2失败
	Message       string `gorm:"column:message;type:varchar(100);not null;default:''" json:"message"`
	CreateTime    int64  `gorm:"column:create_time;type:int(11);not null;default:0" json:"create_time"`
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

// RolePage 角色与页面/按钮关联
type RolePage struct {
	Id int `gorm:"column:id;type:int(11);primary_key;not null;auto_increment" json:"id" form:"id"`
	// RoleID 角色ID
	RoleId int64 `gorm:"column:role_id;type:int(11);not null" json:"role_id" form:"role_id"`
	// PageID 页面或按钮ID
	PageId int64 `gorm:"column:page_id;type:int(11);not null" json:"page_id" form:"page_id"`
}

// 页面类型
const (
	PageTypeMenu   = 1 // 菜单
	PageTypeButton = 2 // 按钮
)

// 账户类型
const (
	AccountTypeAdmin = 1 // 超级管理员，跳过按钮鉴权
)

// 默认 CRUD 按钮编码
const (
	ButtonQuery  = "query"
	ButtonCreate = "create"
	ButtonUpdate = "update"
	ButtonDelete = "delete"
)

type Page struct {
	Id int64 `gorm:"column:id;type:int(11);primary_key;not null;auto_increment" json:"id" form:"id"`
	// PageName 菜单/按钮名称
	PageName string `gorm:"column:page_name;type:varchar(50);not null" json:"page_name" form:"page_name"`
	// Url 前端路由（菜单用）
	Url string `gorm:"column:url;type:varchar(255);not null" json:"url" form:"url"`
	// 父级id
	ParentId int64 `gorm:"column:parent_id;type:int(11);not null" json:"parent_id" form:"parent_id"`
	// 页面等级 1一级菜单 2二级菜单 3按钮
	PageLevel int `gorm:"column:page_level;type:int(11);not null" json:"page_level" form:"page_level"`
	// PageType 1菜单 2按钮
	PageType int `gorm:"column:page_type;type:tinyint(4);not null;default:1" json:"page_type" form:"page_type"`
	// ButtonCode 按钮编码 query/create/update/delete
	ButtonCode string `gorm:"column:button_code;type:varchar(50);not null;default:''" json:"button_code" form:"button_code"`
	// ApiMethod 对应 HTTP 方法，按钮用
	ApiMethod string `gorm:"column:api_method;type:varchar(16);not null;default:''" json:"api_method" form:"api_method"`
	// ApiPath 对应 Gin FullPath，如 /api/admin/account/:id
	ApiPath string `gorm:"column:api_path;type:varchar(255);not null;default:''" json:"api_path" form:"api_path"`
	// CreateTime 创建时间
	CreateTime int64 `gorm:"column:create_time;type:int(11);not null" json:"create_time" form:"create_time"`
	// ApiPrefix 创建二级菜单时传入，自动挂 CRUD 按钮；不入库
	ApiPrefix string `gorm:"-" json:"api_prefix,omitempty" form:"api_prefix"`
}
