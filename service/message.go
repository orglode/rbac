package service

import (
	"context"
	"time"

	error2 "rbac/api/error"
	"rbac/model"
)

func (s *RbacService) SendSiteMessage(ctx context.Context, req *model.SiteMessageSendReq, senderId int64) (int64, error) {
	if req.Title == "" || req.Content == "" {
		return 0, error2.MissingParameter
	}
	if req.MsgType != model.MsgTypeBroadcast && req.MsgType != model.MsgTypeBiz && req.MsgType != model.MsgTypeDirect {
		return 0, error2.MissingParameter
	}

	receiverIDs := make([]int64, 0)
	switch req.MsgType {
	case model.MsgTypeBroadcast:
		ids, err := s.dao.Rbac.ListEnabledAccountIDs(ctx)
		if err != nil {
			return 0, err
		}
		receiverIDs = ids
	case model.MsgTypeBiz, model.MsgTypeDirect:
		if len(req.AccountIds) == 0 {
			// 业务通知未指定人时默认全员
			if req.MsgType == model.MsgTypeBiz {
				ids, err := s.dao.Rbac.ListEnabledAccountIDs(ctx)
				if err != nil {
					return 0, err
				}
				receiverIDs = ids
			} else {
				return 0, error2.MissingParameter
			}
		} else {
			receiverIDs = uniqueInt64(req.AccountIds)
		}
	}
	if len(receiverIDs) == 0 {
		return 0, error2.MissingParameter
	}

	senderName := ""
	if senderId > 0 {
		if acc, err := s.dao.Rbac.GetAccountById(ctx, senderId); err == nil && acc != nil {
			senderName = acc.UserName
			if senderName == "" {
				senderName = acc.AccountNumber
			}
		}
	}

	now := time.Now().Unix()
	msg := &model.SiteMessage{
		Title:      req.Title,
		Content:    req.Content,
		MsgType:    req.MsgType,
		BizType:    req.BizType,
		SenderId:   senderId,
		SenderName: senderName,
		CreateTime: now,
	}
	if err := s.dao.Rbac.CreateSiteMessage(ctx, msg); err != nil {
		return 0, err
	}

	rows := make([]*model.SiteMessageInbox, 0, len(receiverIDs))
	for _, aid := range receiverIDs {
		rows = append(rows, &model.SiteMessageInbox{
			MessageId:  msg.Id,
			AccountId:  aid,
			IsRead:     model.MsgUnread,
			CreateTime: now,
		})
	}
	if err := s.dao.Rbac.CreateSiteMessageInboxBatch(ctx, rows); err != nil {
		return 0, err
	}
	return msg.Id, nil
}

func (s *RbacService) GetMyInbox(ctx context.Context, accountId int64, query model.SiteMessageInboxQuery) (*model.AccountListRest, error) {
	query.Normalize()
	list, total, err := s.dao.Rbac.GetMyInboxList(ctx, accountId, query)
	if err != nil {
		return nil, err
	}
	return &model.AccountListRest{List: list, Total: total}, nil
}

func (s *RbacService) GetUnreadMessageCount(ctx context.Context, accountId int64) (int64, error) {
	return s.dao.Rbac.GetInboxUnreadCount(ctx, accountId)
}

func (s *RbacService) GetRecentUnreadMessages(ctx context.Context, accountId int64, limit int) ([]*model.SiteMessageInboxItem, error) {
	return s.dao.Rbac.GetRecentUnreadInbox(ctx, accountId, limit)
}

func (s *RbacService) MarkMessagesRead(ctx context.Context, accountId int64, req *model.SiteMessageReadReq) error {
	if req == nil {
		return error2.MissingParameter
	}
	now := time.Now().Unix()
	if req.All {
		return s.dao.Rbac.MarkInboxRead(ctx, accountId, nil, now)
	}
	if len(req.InboxIds) == 0 {
		return error2.MissingParameter
	}
	return s.dao.Rbac.MarkInboxRead(ctx, accountId, req.InboxIds, now)
}

func (s *RbacService) GetAdminMessageList(ctx context.Context, query model.SiteMessageAdminQuery) (*model.AccountListRest, error) {
	query.Normalize()
	list, total, err := s.dao.Rbac.GetAdminMessageList(ctx, query)
	if err != nil {
		return nil, err
	}
	items := make([]*model.SiteMessageAdminItem, 0, len(list))
	for _, m := range list {
		totalRecv, readCnt, err := s.dao.Rbac.CountInboxByMessage(ctx, m.Id)
		if err != nil {
			return nil, err
		}
		item := &model.SiteMessageAdminItem{
			SiteMessage:   *m,
			ReceiverCount: int(totalRecv),
			ReadCount:     int(readCnt),
		}
		if m.MsgType == model.MsgTypeDirect || (m.MsgType == model.MsgTypeBiz && totalRecv < 50) {
			ids, _ := s.dao.Rbac.ListInboxAccountIDs(ctx, m.Id)
			item.AccountIds = ids
		}
		items = append(items, item)
	}
	return &model.AccountListRest{List: items, Total: total}, nil
}

func uniqueInt64(ids []int64) []int64 {
	seen := make(map[int64]struct{}, len(ids))
	out := make([]int64, 0, len(ids))
	for _, id := range ids {
		if id <= 0 {
			continue
		}
		if _, ok := seen[id]; ok {
			continue
		}
		seen[id] = struct{}{}
		out = append(out, id)
	}
	return out
}
