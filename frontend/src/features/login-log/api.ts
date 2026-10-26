import { request } from '@/shared/http/request'
import type { LoginLog, PageResult } from '@/shared/types/api'

export function fetchLoginLogs(params: {
  page: number
  size: number
  account_number?: string
  account_id?: number
  status?: number
}) {
  return request<PageResult<LoginLog>>({
    url: '/api/admin/login/log',
    method: 'get',
    params,
  })
}
