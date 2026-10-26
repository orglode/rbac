import { request } from '@/shared/http/request'
import type { AccountInfo, LoginResult } from '@/shared/types/api'

export function login(data: { account_number: string; password: string }) {
  return request<LoginResult>({
    url: '/api/admin/login',
    method: 'post',
    data,
  })
}

export function fetchMyMenu() {
  return request<AccountInfo>({
    url: '/api/admin/me/menu',
    method: 'get',
  })
}

export function fetchMyButtons() {
  return request({
    url: '/api/admin/me/buttons',
    method: 'get',
  })
}
