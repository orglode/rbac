package model

// 站内信类型
const (
	MsgTypeBroadcast = 1 // 系统广播（全员）
	MsgTypeBiz       = 2 // 业务通知（可全员或指定人，由 receiver 决定）
	MsgTypeDirect    = 3 // 指定用户
)

// 已读状态
const (
	MsgUnread = 0
	MsgRead   = 1
)

// SiteMessage 站内信正文
type SiteMessage struct {
	Id         int64  `gorm:"column:id;primary_key;auto_increment" json:"id"`
	Title      string `gorm:"column:title;type:varchar(120);not null;default:''" json:"title"`
	Content    string `gorm:"column:content;type:varchar(2000);not null;default:''" json:"content"`
	MsgType    int    `gorm:"column:msg_type;type:tinyint;not null;default:1" json:"msg_type"`
	BizType    string `gorm:"column:biz_type;type:varchar(32);not null;default:''" json:"biz_type"`
	SenderId   int64  `gorm:"column:sender_id;type:int;not null;default:0" json:"sender_id"`
	SenderName string `gorm:"column:sender_name;type:varchar(50);not null;default:''" json:"sender_name"`
	CreateTime int64  `gorm:"column:create_time;type:int;not null;default:0" json:"create_time"`
}

// SiteMessageInbox 收件箱（含已读）
type SiteMessageInbox struct {
	Id         int64 `gorm:"column:id;primary_key;auto_increment" json:"id"`
	MessageId  int64 `gorm:"column:message_id;type:int;not null;default:0" json:"message_id"`
	AccountId  int64 `gorm:"column:account_id;type:int;not null;default:0" json:"account_id"`
	IsRead     int   `gorm:"column:is_read;type:tinyint;not null;default:0" json:"is_read"`
	ReadTime   int64 `gorm:"column:read_time;type:int;not null;default:0" json:"read_time"`
	CreateTime int64 `gorm:"column:create_time;type:int;not null;default:0" json:"create_time"`
}

// SiteMessageSendReq 发送站内信
type SiteMessageSendReq struct {
	Title      string  `json:"title" binding:"required"`
	Content    string  `json:"content" binding:"required"`
	MsgType    int     `json:"msg_type" binding:"required"` // 1广播 2业务 3指定
	BizType    string  `json:"biz_type"`
	AccountIds []int64 `json:"account_ids"` // msg_type=2/3 时指定接收人；广播忽略
}

// SiteMessageInboxQuery 我的收件箱
type SiteMessageInboxQuery struct {
	Paging
	MsgType    int `json:"msg_type" form:"msg_type"`
	ReadStatus int `json:"read_status" form:"read_status"` // 0不限 1未读 2已读
}

// SiteMessageAdminQuery 发送记录
type SiteMessageAdminQuery struct {
	Paging
	MsgType int    `json:"msg_type" form:"msg_type"`
	Title   string `json:"title" form:"title"`
}

// SiteMessageInboxItem 收件箱列表项
type SiteMessageInboxItem struct {
	InboxId    int64  `json:"inbox_id"`
	MessageId  int64  `json:"message_id"`
	Title      string `json:"title"`
	Content    string `json:"content"`
	MsgType    int    `json:"msg_type"`
	BizType    string `json:"biz_type"`
	SenderId   int64  `json:"sender_id"`
	SenderName string `json:"sender_name"`
	IsRead     int    `json:"is_read"`
	ReadTime   int64  `json:"read_time"`
	CreateTime int64  `json:"create_time"`
}

// SiteMessageReadReq 标记已读
type SiteMessageReadReq struct {
	InboxIds []int64 `json:"inbox_ids"`
	All      bool    `json:"all"`
}

// SiteMessageAdminItem 发送记录项
type SiteMessageAdminItem struct {
	SiteMessage
	ReceiverCount int   `json:"receiver_count"`
	ReadCount     int   `json:"read_count"`
	AccountIds    []int64 `json:"account_ids,omitempty"`
}
