import { createRouter, createWebHistory, type RouteRecordRaw } from 'vue-router'
import { getToken } from '@/shared/utils/auth'
import { useUserStore } from '@/stores/user'
import { resolveView } from './view-map'

const staticRoutes: RouteRecordRaw[] = [
  {
    path: '/login',
    name: 'login',
    component: () => import('@/features/auth/LoginView.vue'),
    meta: { public: true },
  },
  {
    path: '/',
    name: 'layout',
    component: () => import('@/app/layout/AdminLayout.vue'),
    redirect: '/dashboard',
    children: [
      {
        path: 'dashboard',
        name: 'dashboard',
        component: () => import('@/features/dashboard/DashboardView.vue'),
        meta: { title: '首页' },
      },
      {
        path: 'profile',
        name: 'profile',
        component: () => import('@/features/profile/ProfileView.vue'),
        meta: { title: '用户资料' },
      },
    ],
  },
  {
    path: '/:pathMatch(.*)*',
    name: 'not-found',
    component: () => import('@/features/auth/NotFoundView.vue'),
    meta: { public: true },
  },
]

const router = createRouter({
  history: createWebHistory(),
  routes: staticRoutes,
})

let dynamicReady = false

export function resetDynamicRoutes() {
  dynamicReady = false
}

async function ensureDynamicRoutes() {
  if (dynamicReady) return false
  const user = useUserStore()
  if (!user.loaded) {
    await user.loadMenus()
  }
  for (const top of user.menus) {
    for (const child of top.child_page || []) {
      const path = child.page_path
      if (!path || !resolveView(path)) continue
      const name = `menu-${child.page_id}`
      if (router.hasRoute(name)) continue
      router.addRoute('layout', {
        path: path.replace(/^\//, ''),
        name,
        component: resolveView(path)!,
        meta: { title: child.page_name, menuPath: path },
      })
    }
  }
  dynamicReady = true
  return true
}

router.beforeEach(async (to) => {
  if (to.path === '/login') {
    if (getToken()) return { path: '/' }
    return true
  }

  if (!getToken()) {
    if (to.meta.public) return true
    return { path: '/login', query: { redirect: to.fullPath } }
  }

  const user = useUserStore()
  try {
    const justAdded = await ensureDynamicRoutes()
    if (justAdded) {
      return { path: to.fullPath, replace: true }
    }
  } catch {
    user.logout()
    resetDynamicRoutes()
    return { path: '/login' }
  }

  if (to.name === 'layout' || to.path === '/') {
    const leaves = user.collectLeafPaths()
    if (leaves.length) return leaves[0]
  }

  return true
})

export default router
