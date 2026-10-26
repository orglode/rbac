import { request } from '@/shared/http/request'
import type { Account, PageResult } from '@/shared/types/api'

export interface AccountListQuery {
  page: number
  size: number
  user_name?: string
  mobile?: string
  account_number?: string
  status?: number
}

export function fetchAccountList(params: AccountListQuery) {
  return request<PageResult<Account>>({
    url: '/api/admin/account',
    method: 'get',
    params,
  })
}

export function createAccount(data: Partial<Account>) {
  return request({ url: '/api/admin/account', method: 'post', data })
}

export function updateAccount(id: number, data: Partial<Account>) {
  return request({ url: `/api/admin/account/${id}`, method: 'put', data })
}

export function deleteAccount(id: number) {
  return request({ url: `/api/admin/account/${id}`, method: 'delete' })
}

export function assignAccountRole(data: { account_id: number; role_id: number }) {
  return request({ url: '/api/admin/account/role', method: 'post', data })
}

export function removeAccountRole(data: { account_id: number }) {
  return request({ url: '/api/admin/account/role', method: 'delete', data })
}
