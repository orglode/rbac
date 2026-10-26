import { request } from '@/shared/http/request'
import type { MenuPage, PageInfo, PageResult } from '@/shared/types/api'

export interface PageListQuery {
  page: number
  size: number
  page_name?: string
  api?: string
}

export function fetchPageList(params: PageListQuery) {
  return request<PageResult<PageInfo>>({
    url: '/api/admin/page',
    method: 'get',
    params,
  })
}

export function fetchAllPages() {
  return request<MenuPage[]>({ url: '/api/admin/page/all', method: 'get' })
}

export function createPage(data: Record<string, unknown>) {
  return request({ url: '/api/admin/page', method: 'post', data })
}

export function updatePage(id: number, data: Record<string, unknown>) {
  return request({ url: `/api/admin/page/${id}`, method: 'put', data })
}

export function deletePage(id: number) {
  return request({ url: `/api/admin/page/${id}`, method: 'delete' })
}
