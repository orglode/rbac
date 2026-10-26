<script setup lang="ts">
import { useUserStore } from '@/stores/user'

const user = useUserStore()
</script>

<template>
  <el-card shadow="never">
    <template #header>
      <span>工作台</span>
    </template>
    <p>你好，{{ user.profile?.user_name || '管理员' }}。</p>
    <p v-if="user.roleName">当前角色：{{ user.roleName }}</p>
    <p class="tip">请从左侧二级菜单进入账号、角色、页面等管理功能。</p>
  </el-card>
</template>

<style scoped>
p {
  margin: 0 0 8px;
  color: #606266;
}
.tip {
  margin-top: 16px;
  color: #909399;
  font-size: 13px;
}
</style>
