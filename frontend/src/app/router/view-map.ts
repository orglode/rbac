import type { Component } from 'vue'

/** 菜单 url → 懒加载页面组件（按功能目录） */
const viewLoaders: Record<string, () => Promise<{ default: Component }>> = {
  '/system/account': () => import('@/features/account/AccountView.vue'),
  '/system/role': () => import('@/features/role/RoleView.vue'),
  '/system/page': () => import('@/features/page/PageView.vue'),
  '/system/login-log': () => import('@/features/login-log/LoginLogView.vue'),
  '/system/operation-log': () => import('@/features/operation-log/OperationLogView.vue'),
  '/system/message': () => import('@/features/message/MessageView.vue'),
}

export function resolveView(path: string) {
  return viewLoaders[path]
}

export function knownMenuPaths() {
  return Object.keys(viewLoaders)
}
