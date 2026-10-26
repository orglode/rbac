<script setup lang="ts">
import { computed, ref } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { ArrowDown, Fold, Expand, SwitchButton, User } from '@element-plus/icons-vue'
import { useUserStore } from '@/stores/user'
import { resetDynamicRoutes } from '@/app/router'
import defaultAvatar from '@/assets/default-avatar.svg'

const route = useRoute()
const router = useRouter()
const user = useUserStore()
const collapsed = ref(false)

const activeMenu = computed(() => (route.meta.menuPath as string) || route.path)
const displayRole = computed(() => user.roleName || user.profile?.role_name || '未分配角色')

interface Crumb {
  title: string
  path?: string
}

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

function onUserCommand(cmd: string) {
  if (cmd === 'profile') {
    router.push('/profile')
    return
  }
  if (cmd === 'logout') {
    onLogout()
  }
}

function onLogout() {
  user.logout()
  resetDynamicRoutes()
  router.replace('/login')
}
</script>

<template>
  <el-container class="layout">
    <el-aside :width="collapsed ? '64px' : '220px'" class="aside">
      <div class="brand">
        <span v-if="!collapsed">RBAC Admin</span>
        <span v-else>R</span>
      </div>
      <el-scrollbar>
        <el-menu
          :default-active="activeMenu"
          :collapse="collapsed"
          router
          background-color="#1f2d3d"
          text-color="#bfcbd9"
          active-text-color="#409eff"
        >
          <el-menu-item index="/dashboard">
            <span>首页</span>
          </el-menu-item>
          <el-sub-menu v-for="top in user.menus" :key="top.page_id" :index="top.page_path || String(top.page_id)">
            <template #title>
              <span>{{ top.page_name }}</span>
            </template>
            <el-menu-item
              v-for="child in top.child_page"
              :key="child.page_id"
              :index="child.page_path"
            >
              {{ child.page_name }}
            </el-menu-item>
          </el-sub-menu>
        </el-menu>
      </el-scrollbar>
    </el-aside>

    <el-container>
      <el-header class="header">
        <div class="header-left">
          <el-icon class="collapse-btn" @click="collapsed = !collapsed">
            <Fold v-if="!collapsed" />
            <Expand v-else />
          </el-icon>
          <el-breadcrumb separator="/">
            <el-breadcrumb-item v-for="(item, idx) in breadcrumbs" :key="`${item.title}-${idx}`">
              <a
                v-if="item.path && idx < breadcrumbs.length - 1"
                href="javascript:;"
                @click.prevent="onCrumbClick(item)"
              >{{ item.title }}</a>
              <span v-else class="crumb-current">{{ item.title }}</span>
            </el-breadcrumb-item>
          </el-breadcrumb>
        </div>
        <div class="header-right">
          <el-dropdown trigger="click" @command="onUserCommand">
            <div class="user-entry">
              <el-avatar :size="32" :src="defaultAvatar" />
              <div class="user-meta">
                <span class="user-name">{{ user.profile?.user_name || '用户' }}</span>
                <span class="user-role">{{ displayRole }}</span>
              </div>
              <el-icon class="user-arrow"><ArrowDown /></el-icon>
            </div>
            <template #dropdown>
              <el-dropdown-menu>
                <el-dropdown-item command="profile" :icon="User">用户资料</el-dropdown-item>
                <el-dropdown-item divided command="logout" :icon="SwitchButton">退出登录</el-dropdown-item>
              </el-dropdown-menu>
            </template>
          </el-dropdown>
        </div>
      </el-header>
      <el-main class="main">
        <router-view />
      </el-main>
    </el-container>
  </el-container>
</template>

<style scoped>
.layout {
  height: 100vh;
  background: #f0f2f5;
}
.aside {
  background: #1f2d3d;
  transition: width 0.2s;
  overflow: hidden;
}
.brand {
  height: 56px;
  display: flex;
  align-items: center;
  justify-content: center;
  color: #fff;
  font-weight: 600;
  letter-spacing: 0.5px;
  border-bottom: 1px solid rgba(255, 255, 255, 0.08);
}
.header {
  display: flex;
  align-items: center;
  justify-content: space-between;
  background: #fff;
  border-bottom: 1px solid #ebeef5;
  height: 56px;
}
.header-left {
  display: flex;
  align-items: center;
  gap: 12px;
}
.collapse-btn {
  cursor: pointer;
  font-size: 18px;
}
.crumb-current {
  color: #303133;
  font-weight: 500;
}
:deep(.el-breadcrumb__inner a) {
  font-weight: 400;
  color: #909399;
}
:deep(.el-breadcrumb__inner a:hover) {
  color: #409eff;
}
.header-right {
  display: flex;
  align-items: center;
}
.user-entry {
  display: flex;
  align-items: center;
  gap: 10px;
  cursor: pointer;
  padding: 4px 8px;
  border-radius: 6px;
  outline: none;
}
.user-entry:hover {
  background: #f5f7fa;
}
.user-meta {
  display: flex;
  flex-direction: column;
  line-height: 1.2;
  max-width: 120px;
}
.user-name {
  color: #303133;
  font-size: 14px;
  font-weight: 500;
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
}
.user-role {
  color: #909399;
  font-size: 12px;
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
}
.user-arrow {
  color: #909399;
  font-size: 12px;
}
.main {
  padding: 16px;
}
</style>
