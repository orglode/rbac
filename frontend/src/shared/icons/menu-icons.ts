import * as ElementPlusIconsVue from '@element-plus/icons-vue'
import type { Component } from 'vue'

/** Element Plus Icons 全量组件表（组件名 → Vue 组件） */
export const iconMap = ElementPlusIconsVue as unknown as Record<string, Component>

/** 全部可用图标名（PascalCase，与 @element-plus/icons-vue 导出一致） */
export const allIconNames: string[] = Object.keys(iconMap)
  .filter((name) => /^[A-Z]/.test(name))
  .sort()

/** 侧栏常用图标（优先展示在选择器） */
export const preferredIconNames = [
  'Odometer',
  'House',
  'Setting',
  'User',
  'UserFilled',
  'Avatar',
  'Menu',
  'Grid',
  'Key',
  'Lock',
  'Document',
  'Notebook',
  'Tickets',
  'List',
  'Files',
  'Folder',
  'FolderOpened',
  'Monitor',
  'Platform',
  'DataAnalysis',
  'DataBoard',
  'DataLine',
  'Histogram',
  'PieChart',
  'TrendCharts',
  'Bell',
  'Message',
  'ChatDotRound',
  'Search',
  'Edit',
  'Delete',
  'Plus',
  'Refresh',
  'Tools',
  'Operation',
  'Management',
  'Stamp',
  'Collection',
  'Star',
  'Flag',
  'Location',
  'Calendar',
  'Clock',
  'Timer',
  'Warning',
  'InfoFilled',
  'CircleCheck',
  'View',
  'Hide',
  'Link',
  'Share',
  'Download',
  'Upload',
  'Printer',
  'Picture',
  'Camera',
  'Phone',
  'Cellphone',
  'Postcard',
  'CreditCard',
  'Wallet',
  'Goods',
  'ShoppingCart',
  'Box',
  'Briefcase',
  'Suitcase',
  'OfficeBuilding',
  'School',
  'Guide',
].filter((name) => name in iconMap)

export function resolveMenuIcon(name?: string | null): Component | undefined {
  if (!name) return undefined
  return iconMap[name]
}
