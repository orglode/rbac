package model

import "time"

// 权限模块表
type Module struct {
	Id         int64     `gorm:"column:id" json:"id"`
	Name       string    `gorm:"column:name" json:"name"`               //模块名称
	ApiPath    string    `gorm:"column:api_path" json:"api_path"`       //后端接口
	Parent     int64     `gorm:"column:parent" json:"parent"`           //所属父级信息
	Sort       int       `gorm:"column:sort" json:"sort"`               //排序 数值越大越靠后
	Code       int       `gorm:"column:code" json:"code"`               //按钮标记 0--菜单 1--查询 2-添加 3--修改 4-删除 5-修改其他状态，6–导出，其他标记 ...
	Type       string    `gorm:"column:type" json:"type"`               //前端区分按钮 页面 1–菜单 2--按钮
	Root       string    `gorm:"column:root" json:"root"`               //前端路由使用
	Status     int8      `gorm:"column:status" json:"status"`           //1--启用 2--禁用 3--菜单上不显示
	OperatorId int64     `gorm:"column:operator_id" json:"operator_id"` //操作者id
	Operator   string    `gorm:"column:operator" json:"operator"`       //操作者
	CreateTime time.Time `gorm:"column:create_time" json:"create_time"` //创建时间
	UpdateTime time.Time `gorm:"column:update_time" json:"update_time"` //更新时间
}

// 角色类型表
type RoleType struct {
	Id         int64  `gorm:"column:id" json:"id" form:"id"`
	TypeName   string `gorm:"column:type_name" json:"type_name" form:"type_name"`       //角色类型名称
	Status     int8   `gorm:"column:status" json:"status" form:"status"`                //状态 1：有效 2：删除
	CreateTime int64  `gorm:"column:create_time" json:"create_time" form:"create_time"` //创建时间
	UpdateTime int64  `gorm:"column:update_time" json:"update_time" form:"update_time"` //更新时间
}

// 角色表（含数据范围）
type Role struct {
	Id          int64     `gorm:"column:id" json:"id"`
	Pid         string    `gorm:"column:pid" json:"pid"`                   //父级角色id，多个逗号分隔
	TypeId      int16     `gorm:"column:type_id" json:"type_id"`           //角色类型id
	Name        string    `gorm:"column:name" json:"name"`                 //角色名称
	Description string    `gorm:"column:description" json:"description"`   //角色描述
	OperatorId  int64     `gorm:"column:operator_id" json:"operator_id"`   //后台操作人id
	Operator    string    `gorm:"column:operator" json:"operator"`         //操作人名称
	PageTypeId  int       `gorm:"column:page_type_id" json:"page_type_id"` //页面分类id
	Status      int8      `gorm:"column:status" json:"status"`             //状态 1--启用 2--禁用
	CreateTime  time.Time `gorm:"column:create_time" json:"create_time"`   //创建时间
	UpdateTime  time.Time `gorm:"column:update_time" json:"update_time"`   //更新时间
}

// 角色权限映射关系表
type RoleModule struct {
	Id       int64 `gorm:"column:id" json:"id"`
	RoleId   int16 `gorm:"column:role_id" json:"role_id"`     //角色ID
	ModuleId int   `gorm:"column:module_id" json:"module_id"` //模块id
}

// 页面分类表
type PageType struct {
	Id         int64     `gorm:"column:id" json:"id"`
	Name       string    `gorm:"column:name" json:"name"`               //页面分类名称
	ModuleIds  string    `gorm:"column:module_ids" json:"module_ids"`   //模块ids
	Status     int8      `gorm:"column:status" json:"status"`           //状态 1：有效 2：删除
	Operator   string    `gorm:"column:operator" json:"operator"`       //操作者
	CreateTime time.Time `gorm:"column:create_time" json:"create_time"` //创建时间
	UpdateTime time.Time `gorm:"column:update_time" json:"update_time"` //更新时间
}

// 系统下需要跳过校验的特殊url表
type UrlSkip struct {
	Id      int64  `gorm:"column:id" json:"id"`
	SkipUrl string `gorm:"column:skip_url" json:"skip_url"` //后端需要跳过校验的接口
	Status  int8   `gorm:"column:status" json:"status"`     //状态 1：有效 2：无效
	Desc    string `gorm:"column:desc" json:"desc"`         //url描述
}

// 用户角色关系表
type InkeUserRole struct {
	Id     int64 `gorm:"column:id" json:"id"`
	UserId int64 `gorm:"column:user_id" json:"user_id"` //用户id
	RoleId int64 `gorm:"column:role_id" json:"role_id"` //角色id
}

// 系统账号表
type Users struct {
	Id            uint      `gorm:"column:id" json:"id"`
	Username      string    `gorm:"column:username" json:"username"`               //账号姓名
	Email         string    `gorm:"column:email" json:"email"`                     //账号
	Phone         string    `gorm:"column:phone" json:"phone"`                     //电话
	Description   string    `gorm:"column:description" json:"description"`         //描述
	Status        int8      `gorm:"column:status" json:"status"`                   //状态 1:启用 2:禁用
	Operator      string    `gorm:"column:operator" json:"operator"`               //操作者
	LastLoginTime time.Time `gorm:"column:last_login_time" json:"last_login_time"` //最后一次登录时间
	CreateTime    time.Time `gorm:"column:create_time" json:"create_time"`         //创建时间
	UpdateTime    time.Time `gorm:"column:update_time" json:"update_time"`         //更新时间
}
