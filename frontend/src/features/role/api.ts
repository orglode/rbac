import { request } from '@/shared/http/request'
import type { PageInfo, PageResult, Role } from '@/shared/types/api'

export interface RolePageInfo {
  role_id: number
  role_name: string
  page: PageInfo[]
}

export function fetchRoleList(params: { page: number; size: number }) {
  return request<PageResult<Role>>({
    url: '/api/admin/role',
    method: 'get',
    params,
  })
}

export function fetchAllRoles() {
  return request<Role[]>({ url: '/api/admin/role/all', method: 'get' })
}

export function createRole(data: { role_name: string }) {
  return request({ url: '/api/admin/role', method: 'post', data })
}

export function updateRole(id: number, data: { role_name: string }) {
  return request({ url: `/api/admin/role/${id}`, method: 'put', data })
}

export function deleteRole(id: number) {
  return request({ url: `/api/admin/role/${id}`, method: 'delete' })
}

export function fetchRolePages(roleId: number) {
  return request<RolePageInfo>({ url: `/api/admin/role/page/${roleId}`, method: 'get' })
}

export function assignRolePages(data: { role_id: number; page_ids: string }) {
  return request({ url: '/api/admin/role/page', method: 'post', data })
}
