package dao

import (
	"context"
	"rbac/model"
)

func (d *Rabc) CreateSiteMessage(ctx context.Context, msg *model.SiteMessage) error {
	return d.db.Master(ctx).Table(model.SiteMessageTable).Create(msg).Error
}

func (d *Rabc) CreateSiteMessageInboxBatch(ctx context.Context, rows []*model.SiteMessageInbox) error {
	if len(rows) == 0 {
		return nil
	}
	return d.db.Master(ctx).Table(model.SiteMessageInboxTable).CreateInBatches(rows, 200).Error
}

func (d *Rabc) ListEnabledAccountIDs(ctx context.Context) ([]int64, error) {
	ids := make([]int64, 0)
	err := d.db.Slave(ctx).Table(model.AccountTable).
		Where("status=?", model.AccountStatusEnable).
		Pluck("id", &ids).Error
	return ids, err
}

func (d *Rabc) GetSiteMessageById(ctx context.Context, id int64) (*model.SiteMessage, error) {
	var msg model.SiteMessage
	if err := d.db.Slave(ctx).Table(model.SiteMessageTable).Where("id=?", id).First(&msg).Error; err != nil {
		return nil, err
	}
	return &msg, nil
}

func (d *Rabc) GetInboxUnreadCount(ctx context.Context, accountId int64) (int64, error) {
	var count int64
	err := d.db.Slave(ctx).Table(model.SiteMessageInboxTable).
		Where("account_id=? AND is_read=?", accountId, model.MsgUnread).
		Count(&count).Error
	return count, err
}

func (d *Rabc) GetMyInboxList(ctx context.Context, accountId int64, query model.SiteMessageInboxQuery) ([]*model.SiteMessageInboxItem, int64, error) {
	var count int64
	result := make([]*model.SiteMessageInboxItem, 0)

	db := d.db.Slave(ctx).Table(model.SiteMessageInboxTable+" AS i").
		Joins("INNER JOIN "+model.SiteMessageTable+" AS m ON m.id = i.message_id").
		Where("i.account_id=?", accountId)
	if query.MsgType > 0 {
		db = db.Where("m.msg_type=?", query.MsgType)
	}
	if query.ReadStatus == 1 {
		db = db.Where("i.is_read=?", model.MsgUnread)
	} else if query.ReadStatus == 2 {
		db = db.Where("i.is_read=?", model.MsgRead)
	}
	if err := db.Count(&count).Error; err != nil {
		return nil, 0, err
	}

	page := Paging{Page: query.Page, Size: query.Size}
	page.Normalize()
	err := db.Select(`i.id AS inbox_id, i.message_id, m.title, m.content, m.msg_type, m.biz_type,
		m.sender_id, m.sender_name, i.is_read, i.read_time, m.create_time`).
		Order("i.is_read ASC, m.id DESC").
		Offset(page.Offset()).Limit(page.Size).
		Scan(&result).Error
	return result, count, err
}

func (d *Rabc) MarkInboxRead(ctx context.Context, accountId int64, inboxIds []int64, readTime int64) error {
	db := d.db.Master(ctx).Table(model.SiteMessageInboxTable).
		Where("account_id=? AND is_read=?", accountId, model.MsgUnread)
	if len(inboxIds) > 0 {
		db = db.Where("id IN ?", inboxIds)
	}
	return db.Updates(map[string]interface{}{
		"is_read":   model.MsgRead,
		"read_time": readTime,
	}).Error
}

func (d *Rabc) GetAdminMessageList(ctx context.Context, query model.SiteMessageAdminQuery) ([]*model.SiteMessage, int64, error) {
	var count int64
	result := make([]*model.SiteMessage, 0)
	db := d.db.Slave(ctx).Table(model.SiteMessageTable)
	if query.MsgType > 0 {
		db = db.Where("msg_type=?", query.MsgType)
	}
	if query.Title != "" {
		db = db.Where("title LIKE ?", "%"+query.Title+"%")
	}
	if err := db.Count(&count).Error; err != nil {
		return nil, 0, err
	}
	page := Paging{Page: query.Page, Size: query.Size}
	page.Normalize()
	err := db.Order("id DESC").Offset(page.Offset()).Limit(page.Size).Find(&result).Error
	return result, count, err
}

func (d *Rabc) CountInboxByMessage(ctx context.Context, messageId int64) (total int64, read int64, err error) {
	err = d.db.Slave(ctx).Table(model.SiteMessageInboxTable).
		Where("message_id=?", messageId).Count(&total).Error
	if err != nil {
		return
	}
	err = d.db.Slave(ctx).Table(model.SiteMessageInboxTable).
		Where("message_id=? AND is_read=?", messageId, model.MsgRead).Count(&read).Error
	return
}

func (d *Rabc) ListInboxAccountIDs(ctx context.Context, messageId int64) ([]int64, error) {
	ids := make([]int64, 0)
	err := d.db.Slave(ctx).Table(model.SiteMessageInboxTable).
		Where("message_id=?", messageId).
		Pluck("account_id", &ids).Error
	return ids, err
}

func (d *Rabc) GetRecentUnreadInbox(ctx context.Context, accountId int64, limit int) ([]*model.SiteMessageInboxItem, error) {
	if limit <= 0 {
		limit = 5
	}
	result := make([]*model.SiteMessageInboxItem, 0)
	err := d.db.Slave(ctx).Table(model.SiteMessageInboxTable+" AS i").
		Joins("INNER JOIN "+model.SiteMessageTable+" AS m ON m.id = i.message_id").
		Where("i.account_id=? AND i.is_read=?", accountId, model.MsgUnread).
		Select(`i.id AS inbox_id, i.message_id, m.title, m.content, m.msg_type, m.biz_type,
			m.sender_id, m.sender_name, i.is_read, i.read_time, m.create_time`).
		Order("m.id DESC").Limit(limit).
		Scan(&result).Error
	return result, err
}
