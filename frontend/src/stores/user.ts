import { defineStore } from 'pinia'
import { computed, ref } from 'vue'
import type { AccountInfo, LoginResult, PageInfo } from '@/shared/types/api'
import { clearAuth, getStoredUser, getToken, setStoredUser, setToken, type StoredUser } from '@/shared/utils/auth'
import { fetchMyMenu, login as loginApi } from '@/features/auth/api'

/** 仅用户会话与菜单权限，不做业务列表缓存 */
export const useUserStore = defineStore('user', () => {
  const token = ref(getToken())
  const profile = ref<StoredUser | null>(getStoredUser())
  const menus = ref<PageInfo[]>([])
  const roleName = ref(profile.value?.role_name || '')
  const loaded = ref(false)

  const isLogin = computed(() => !!token.value)

  function applyLogin(data: LoginResult) {
    token.value = data.token
    setToken(data.token)
    const user: StoredUser = {
      id: data.id,
      user_name: data.user_name,
      mobile: data.mobile,
      account_number: data.account_number,
      role_name: profile.value?.role_name,
    }
    profile.value = user
    setStoredUser(user)
  }

  function syncProfile(partial: Partial<StoredUser>) {
    if (!profile.value) return
    profile.value = { ...profile.value, ...partial }
    setStoredUser(profile.value)
  }

  async function login(account_number: string, password: string) {
    const data = await loginApi({ account_number, password })
    applyLogin(data)
    await loadMenus()
  }

  async function loadMenus() {
    const info: AccountInfo = await fetchMyMenu()
    menus.value = info.page || []
    roleName.value = info.role_name || ''
    if (profile.value) {
      syncProfile({ role_name: roleName.value })
    }
    loaded.value = true
  }

  function logout() {
    token.value = ''
    profile.value = null
    menus.value = []
    roleName.value = ''
    loaded.value = false
    clearAuth()
  }

  /** 当前路由对应二级菜单下的按钮权限 */
  function can(path: string, buttonCode: string): boolean {
    for (const top of menus.value) {
      for (const child of top.child_page || []) {
        if (child.page_path === path) {
          return (child.buttons || []).some((b) => b.button_code === buttonCode)
        }
      }
    }
    return false
  }

  function collectLeafPaths(nodes: PageInfo[] = menus.value): string[] {
    const paths: string[] = []
    for (const n of nodes) {
      if (n.page_level === 2 && n.page_path) paths.push(n.page_path)
      if (n.child_page?.length) paths.push(...collectLeafPaths(n.child_page))
    }
    return paths
  }

  return {
    token,
    profile,
    menus,
    roleName,
    loaded,
    isLogin,
    login,
    loadMenus,
    logout,
    syncProfile,
    can,
    collectLeafPaths,
  }
})
