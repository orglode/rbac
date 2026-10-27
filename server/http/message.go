package http

import (
	"rbac/model"
	"rbac/server/common"
	"rbac/service"

	"github.com/gin-gonic/gin"
)

// MessageHandler 站内信
type MessageHandler struct {
	svc *service.Service
}

func NewMessageHandler(svc *service.Service) *MessageHandler {
	return &MessageHandler{svc: svc}
}

func (h *MessageHandler) UnreadCount(c *gin.Context) {
	count, err := h.svc.RbacService.GetUnreadMessageCount(c.Request.Context(), common.GetUid(c))
	if err != nil {
		common.StandardOutput(c, err)
		return
	}
	common.ResponseSuccess(c, gin.H{"count": count})
}

func (h *MessageHandler) InboxRecent(c *gin.Context) {
	list, err := h.svc.RbacService.GetRecentUnreadMessages(c.Request.Context(), common.GetUid(c), 8)
	if err != nil {
		common.StandardOutput(c, err)
		return
	}
	common.ResponseSuccess(c, list)
}

func (h *MessageHandler) Inbox(c *gin.Context) {
	req := model.SiteMessageInboxQuery{}
	if err := c.ShouldBindQuery(&req); err != nil {
		common.ResponseError(c, common.InvalidParam)
		return
	}
	data, err := h.svc.RbacService.GetMyInbox(c.Request.Context(), common.GetUid(c), req)
	if err != nil {
		common.StandardOutput(c, err)
		return
	}
	common.ResponseSuccess(c, data)
}

func (h *MessageHandler) Read(c *gin.Context) {
	req := &model.SiteMessageReadReq{}
	if err := c.ShouldBindJSON(req); err != nil {
		common.ResponseError(c, common.InvalidParam)
		return
	}
	if err := h.svc.RbacService.MarkMessagesRead(c.Request.Context(), common.GetUid(c), req); err != nil {
		common.StandardOutput(c, err)
		return
	}
	common.ResponseSuccess(c, nil)
}

func (h *MessageHandler) List(c *gin.Context) {
	req := model.SiteMessageAdminQuery{}
	if err := c.ShouldBindQuery(&req); err != nil {
		common.ResponseError(c, common.InvalidParam)
		return
	}
	data, err := h.svc.RbacService.GetAdminMessageList(c.Request.Context(), req)
	if err != nil {
		common.StandardOutput(c, err)
		return
	}
	common.ResponseSuccess(c, data)
}

func (h *MessageHandler) Send(c *gin.Context) {
	req := &model.SiteMessageSendReq{}
	if err := c.ShouldBindJSON(req); err != nil {
		common.ResponseError(c, common.InvalidParam)
		return
	}
	id, err := h.svc.RbacService.SendSiteMessage(c.Request.Context(), req, common.GetUid(c))
	if err != nil {
		common.StandardOutput(c, err)
		return
	}
	common.ResponseSuccess(c, gin.H{"id": id})
}
