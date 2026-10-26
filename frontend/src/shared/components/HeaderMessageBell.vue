<script setup lang="ts">
import { onMounted, onUnmounted, ref } from 'vue'
import { useRouter } from 'vue-router'
import { Bell } from '@element-plus/icons-vue'
import { formatUnix } from '@/shared/utils/time'
import {
  fetchRecentUnread,
  fetchUnreadCount,
  markMessagesRead,
  type SiteMessageInboxItem,
} from '@/features/message/api'

const router = useRouter()
const unread = ref(0)
const recent = ref<SiteMessageInboxItem[]>([])
const loading = ref(false)

const msgTypeLabel: Record<number, string> = {
  1: '广播',
  2: '业务',
  3: '私信',
}

async function refresh() {
  loading.value = true
  try {
    const [countRes, list] = await Promise.all([fetchUnreadCount(), fetchRecentUnread()])
    unread.value = countRes?.count || 0
    recent.value = list || []
  } catch {
    unread.value = 0
    recent.value = []
  } finally {
    loading.value = false
  }
}

async function openItem(item: SiteMessageInboxItem) {
  if (item.is_read === 0) {
    await markMessagesRead({ inbox_ids: [item.inbox_id] })
  }
  await refresh()
  router.push({ path: '/system/message', query: { inbox_id: String(item.inbox_id) } })
}

function goInbox() {
  router.push('/system/message')
}

async function markAll() {
  await markMessagesRead({ all: true })
  await refresh()
}

function onChanged() {
  void refresh()
}

let timer: number | undefined
onMounted(() => {
  void refresh()
  timer = window.setInterval(() => void refresh(), 60000)
  window.addEventListener('site-message-changed', onChanged)
})
onUnmounted(() => {
  if (timer) window.clearInterval(timer)
  window.removeEventListener('site-message-changed', onChanged)
})
</script>

<template>
  <el-dropdown trigger="click" placement="bottom-end" @visible-change="(v: boolean) => v && refresh()">
    <div class="bell-entry">
      <el-badge :value="unread" :hidden="unread <= 0" :max="99">
        <el-icon :size="20"><Bell /></el-icon>
      </el-badge>
    </div>
    <template #dropdown>
      <div class="bell-panel">
        <div class="bell-head">
          <span>站内信</span>
          <div class="bell-actions">
            <el-button link type="primary" size="small" @click="markAll">全部已读</el-button>
            <el-button link type="primary" size="small" @click="goInbox">查看全部</el-button>
          </div>
        </div>
        <el-scrollbar max-height="320px">
          <div v-if="loading && !recent.length" class="bell-empty">加载中...</div>
          <div v-else-if="!recent.length" class="bell-empty">暂无未读消息</div>
          <button
            v-for="item in recent"
            :key="item.inbox_id"
            type="button"
            class="bell-item"
            @click="openItem(item)"
          >
            <div class="bell-item-title">
              <el-tag size="small" type="warning">{{ msgTypeLabel[item.msg_type] || '消息' }}</el-tag>
              <span>{{ item.title }}</span>
            </div>
            <div class="bell-item-meta">{{ item.sender_name || '系统' }} · {{ formatUnix(item.create_time) }}</div>
          </button>
        </el-scrollbar>
      </div>
    </template>
  </el-dropdown>
</template>

<style scoped>
.bell-entry {
  display: flex;
  align-items: center;
  justify-content: center;
  width: 36px;
  height: 36px;
  margin-right: 8px;
  border-radius: 6px;
  cursor: pointer;
  color: #606266;
  outline: none;
}
.bell-entry:hover {
  background: #f5f7fa;
  color: #409eff;
}
.bell-panel {
  width: 340px;
  padding: 8px 0;
}
.bell-head {
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding: 4px 14px 10px;
  font-weight: 600;
  color: #303133;
  border-bottom: 1px solid #ebeef5;
}
.bell-actions {
  display: flex;
  gap: 4px;
}
.bell-empty {
  padding: 28px 16px;
  text-align: center;
  color: #909399;
  font-size: 13px;
}
.bell-item {
  display: block;
  width: 100%;
  text-align: left;
  border: 0;
  background: transparent;
  padding: 10px 14px;
  cursor: pointer;
  border-bottom: 1px solid #f2f3f5;
}
.bell-item:hover {
  background: #f5f7fa;
}
.bell-item-title {
  display: flex;
  align-items: center;
  gap: 6px;
  color: #303133;
  font-size: 13px;
  margin-bottom: 4px;
}
.bell-item-title span {
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
}
.bell-item-meta {
  color: #909399;
  font-size: 12px;
  padding-left: 2px;
}
</style>
