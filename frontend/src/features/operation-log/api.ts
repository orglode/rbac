import { request } from '@/shared/http/request'
import type { OperationLog, PageResult } from '@/shared/types/api'

export function fetchOperationLogs(params: {
  page: number
  size: number
  module?: number
  action?: number
  operator_account?: string
  target_name?: string
}) {
  return request<PageResult<OperationLog>>({
    url: '/api/admin/operation/log',
    method: 'get',
    params,
  })
}
