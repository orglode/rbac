import { computed } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { useUserStore } from '@/stores/user'

export interface Crumb {
  title: string
  path?: string
}

export function useBreadcrumbs() {
  const route = useRoute()
  const router = useRouter()
  const user = useUserStore()

  const breadcrumbs = computed<Crumb[]>(() => {
    const crumbs: Crumb[] = [{ title: '首页', path: '/dashboard' }]
    const current = (route.meta.menuPath as string) || route.path

    if (current === '/dashboard' || route.name === 'dashboard') {
      return crumbs
    }
    if (route.name === 'profile') {
      crumbs.push({ title: '用户资料', path: '/profile' })
      return crumbs
    }

    for (const top of user.menus) {
      for (const child of top.child_page || []) {
        if (child.page_path === current) {
          crumbs.push({ title: top.page_name })
          crumbs.push({ title: child.page_name, path: child.page_path })
          return crumbs
        }
      }
    }

    const title = (route.meta.title as string) || '当前页'
    crumbs.push({ title })
    return crumbs
  })

  function onCrumbClick(item: Crumb) {
    if (item.path && item.path !== route.path) {
      router.push(item.path)
    }
  }

  return { breadcrumbs, onCrumbClick }
}
