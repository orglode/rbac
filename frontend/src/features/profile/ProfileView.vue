<script setup lang="ts">
import { computed } from 'vue'
import { useUserStore } from '@/stores/user'
import defaultAvatar from '@/assets/default-avatar.svg'

const user = useUserStore()

const displayRole = computed(() => user.roleName || user.profile?.role_name || '未分配角色')
</script>

<template>
  <el-card shadow="never" class="profile-card">
    <template #header>
      <span>用户资料</span>
    </template>

    <div class="profile-body">
      <el-avatar :size="96" :src="defaultAvatar" class="avatar" />
      <el-descriptions :column="1" border class="desc">
        <el-descriptions-item label="姓名">
          {{ user.profile?.user_name || '-' }}
        </el-descriptions-item>
        <el-descriptions-item label="账号">
          {{ user.profile?.account_number || '-' }}
        </el-descriptions-item>
        <el-descriptions-item label="手机号">
          {{ user.profile?.mobile || '-' }}
        </el-descriptions-item>
        <el-descriptions-item label="角色">
          <el-tag type="primary" effect="plain">{{ displayRole }}</el-tag>
        </el-descriptions-item>
        <el-descriptions-item label="用户 ID">
          {{ user.profile?.id ?? '-' }}
        </el-descriptions-item>
      </el-descriptions>
    </div>
  </el-card>
</template>

<style scoped>
.profile-body {
  display: flex;
  align-items: flex-start;
  gap: 32px;
  padding: 8px 0 16px;
}
.avatar {
  flex-shrink: 0;
  background: #dce3ec;
}
.desc {
  flex: 1;
  max-width: 520px;
}
</style>
