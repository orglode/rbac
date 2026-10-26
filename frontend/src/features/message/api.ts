import { request } from '@/shared/http/request'
import type { PageResult } from '@/shared/types/api'

export interface SiteMessageInboxItem {
  inbox_id: number
  message_id: number
  title: string
  content: string
  msg_type: number
  biz_type: string
  sender_id: number
  sender_name: string
  is_read: number
  read_time: number
  create_time: number
}

export interface SiteMessageAdminItem {
  id: number
  title: string
  content: string
  msg_type: number
  biz_type: string
  sender_id: number
  sender_name: string
  create_time: number
  receiver_count: number
  read_count: number
  account_ids?: number[]
}

export function fetchUnreadCount() {
  return request<{ count: number }>({
    url: '/api/admin/message/unread-count',
    method: 'get',
  })
}

export function fetchRecentUnread() {
  return request<SiteMessageInboxItem[]>({
    url: '/api/admin/message/inbox/recent',
    method: 'get',
  })
}

export function fetchInbox(params: {
  page: number
  size: number
  msg_type?: number
  read_status?: number
}) {
  return request<PageResult<SiteMessageInboxItem>>({
    url: '/api/admin/message/inbox',
    method: 'get',
    params,
  })
}

export function markMessagesRead(data: { inbox_ids?: number[]; all?: boolean }) {
  return request({
    url: '/api/admin/message/read',
    method: 'post',
    data,
  })
}

export function fetchSentMessages(params: {
  page: number
  size: number
  msg_type?: number
  title?: string
}) {
  return request<PageResult<SiteMessageAdminItem>>({
    url: '/api/admin/message',
    method: 'get',
    params,
  })
}

export function sendMessage(data: {
  title: string
  content: string
  msg_type: number
  biz_type?: string
  account_ids?: number[]
}) {
  return request<{ id: number }>({
    url: '/api/admin/message',
    method: 'post',
    data,
  })
}
