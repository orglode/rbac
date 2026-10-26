<script setup lang="ts">
import { computed, onMounted, reactive, ref } from 'vue'
import { ElMessage, type FormInstance, type FormRules } from 'element-plus'
import type { Account } from '@/shared/types/api'
import { formatUnix } from '@/shared/utils/time'
import PageBreadcrumb from '@/shared/components/PageBreadcrumb.vue'
import { useUserStore } from '@/stores/user'
import { fetchAccountList } from '@/features/account/api'
import {
  fetchInbox,
  fetchSentMessages,
  markMessagesRead,
  sendMessage,
  type SiteMessageAdminItem,
  type SiteMessageInboxItem,
} from './api'

const user = useUserStore()
const path = '/system/message'
const can = (code: string) => user.can(path, code)

const msgTypeLabel: Record<number, string> = {
  1: '系统广播',
  2: '业务通知',
  3: '指定用户',
}

const activeTab = ref('inbox')
const loading = ref(false)
const inboxList = ref<SiteMessageInboxItem[]>([])
const inboxTotal = ref(0)
const inboxQuery = reactive({
  page: 1,
  size: 20,
  msg_type: undefined as number | undefined,
  read_status: undefined as number | undefined,
})

const sentList = ref<SiteMessageAdminItem[]>([])
const sentTotal = ref(0)
const sentQuery = reactive({
  page: 1,
  size: 20,
  msg_type: undefined as number | undefined,
  title: '',
})

const detailVisible = ref(false)
const detail = ref<SiteMessageInboxItem | null>(null)

const sendFormRef = ref<FormInstance>()
const sendForm = reactive({
  title: '',
  content: '',
  msg_type: 1,
  biz_type: '',
  account_ids: [] as number[],
})
const sendRules: FormRules = {
  title: [{ required: true, message: '请输入标题', trigger: 'blur' }],
  content: [{ required: true, message: '请输入内容', trigger: 'blur' }],
  msg_type: [{ required: true, message: '请选择类型', trigger: 'change' }],
}
const accountOptions = ref<Account[]>([])
const needReceivers = computed(() => sendForm.msg_type === 3 || sendForm.msg_type === 2)

async function loadInbox() {
  loading.value = true
  try {
    const params: Record<string, unknown> = {
      page: inboxQuery.page,
      size: inboxQuery.size,
    }
    if (inboxQuery.msg_type) params.msg_type = inboxQuery.msg_type
    if (inboxQuery.read_status) params.read_status = inboxQuery.read_status
    const data = await fetchInbox(params as Parameters<typeof fetchInbox>[0])
    inboxList.value = (data?.list as SiteMessageInboxItem[]) || []
    inboxTotal.value = data?.total || 0
  } finally {
    loading.value = false
  }
}

async function loadSent() {
  loading.value = true
  try {
    const params: Record<string, unknown> = {
      page: sentQuery.page,
      size: sentQuery.size,
    }
    if (sentQuery.msg_type) params.msg_type = sentQuery.msg_type
    if (sentQuery.title.trim()) params.title = sentQuery.title.trim()
    const data = await fetchSentMessages(params as Parameters<typeof fetchSentMessages>[0])
    sentList.value = (data?.list as SiteMessageAdminItem[]) || []
    sentTotal.value = data?.total || 0
  } finally {
    loading.value = false
  }
}

async function loadAccounts() {
  const data = await fetchAccountList({ page: 1, size: 200, status: 1 })
  accountOptions.value = (data?.list as Account[]) || []
}

async function openDetail(row: SiteMessageInboxItem) {
  detail.value = row
  detailVisible.value = true
  if (row.is_read === 0) {
    await markMessagesRead({ inbox_ids: [row.inbox_id] })
    row.is_read = 1
    row.read_time = Math.floor(Date.now() / 1000)
    window.dispatchEvent(new CustomEvent('site-message-changed'))
  }
}

async function markAllRead() {
  await markMessagesRead({ all: true })
  ElMessage.success('已全部标为已读')
  window.dispatchEvent(new CustomEvent('site-message-changed'))
  await loadInbox()
}

async function submitSend() {
  await sendFormRef.value?.validate()
  if (sendForm.msg_type === 3 && sendForm.account_ids.length === 0) {
    ElMessage.warning('请选择接收人')
    return
  }
  await sendMessage({
    title: sendForm.title,
    content: sendForm.content,
    msg_type: sendForm.msg_type,
    biz_type: sendForm.biz_type || undefined,
    account_ids: needReceivers.value ? sendForm.account_ids : undefined,
  })
  ElMessage.success('发送成功')
  sendForm.title = ''
  sendForm.content = ''
  sendForm.biz_type = ''
  sendForm.account_ids = []
  window.dispatchEvent(new CustomEvent('site-message-changed'))
  if (can('query')) {
    activeTab.value = 'sent'
    await loadSent()
  }
}

function onTabChange(name: string | number) {
  if (name === 'inbox') void loadInbox()
  if (name === 'sent') void loadSent()
  if (name === 'send') void loadAccounts()
}

onMounted(() => {
  void loadInbox()
})
</script>

<template>
  <div>
    <PageBreadcrumb />
    <el-card shadow="never">
      <el-tabs v-model="activeTab" @tab-change="onTabChange">
        <el-tab-pane label="收件箱" name="inbox">
          <div class="filter-panel">
            <el-form :model="inboxQuery" inline class="filter-form" label-width="64px" @submit.prevent="() => { inboxQuery.page = 1; loadInbox() }">
              <el-form-item label="类型">
                <el-select v-model="inboxQuery.msg_type" clearable placeholder="全部" style="width: 140px">
                  <el-option :value="1" label="系统广播" />
                  <el-option :value="2" label="业务通知" />
                  <el-option :value="3" label="指定用户" />
                </el-select>
              </el-form-item>
              <el-form-item label="状态">
                <el-select v-model="inboxQuery.read_status" clearable placeholder="全部" style="width: 120px">
                  <el-option :value="1" label="未读" />
                  <el-option :value="2" label="已读" />
                </el-select>
              </el-form-item>
              <el-form-item label-width="0">
                <el-button type="primary" native-type="submit">查询</el-button>
                <el-button @click="loadInbox">刷新</el-button>
                <el-button @click="markAllRead">全部已读</el-button>
              </el-form-item>
            </el-form>
          </div>

          <el-table v-loading="loading" :data="inboxList" border stripe>
            <el-table-column label="状态" width="80">
              <template #default="{ row }">
                <el-tag :type="row.is_read === 0 ? 'danger' : 'info'" size="small">
                  {{ row.is_read === 0 ? '未读' : '已读' }}
                </el-tag>
              </template>
            </el-table-column>
            <el-table-column label="类型" width="100">
              <template #default="{ row }">{{ msgTypeLabel[row.msg_type] || row.msg_type }}</template>
            </el-table-column>
            <el-table-column prop="title" label="标题" min-width="160" show-overflow-tooltip />
            <el-table-column prop="sender_name" label="发送人" width="120" />
            <el-table-column label="时间" width="170">
              <template #default="{ row }">{{ formatUnix(row.create_time) }}</template>
            </el-table-column>
            <el-table-column label="操作" width="100" fixed="right">
              <template #default="{ row }">
                <el-button link type="primary" @click="openDetail(row as SiteMessageInboxItem)">查看</el-button>
              </template>
            </el-table-column>
          </el-table>
          <div class="pager">
            <el-pagination
              v-model:current-page="inboxQuery.page"
              v-model:page-size="inboxQuery.size"
              layout="total, prev, pager, next"
              :total="inboxTotal"
              @current-change="loadInbox"
            />
          </div>
        </el-tab-pane>

        <el-tab-pane v-if="can('create')" label="发送消息" name="send">
          <el-form ref="sendFormRef" :model="sendForm" :rules="sendRules" label-width="96px" style="max-width: 640px">
            <el-form-item label="类型" prop="msg_type">
              <el-radio-group v-model="sendForm.msg_type">
                <el-radio :value="1">系统广播</el-radio>
                <el-radio :value="2">业务通知</el-radio>
                <el-radio :value="3">指定用户</el-radio>
              </el-radio-group>
            </el-form-item>
            <el-form-item v-if="sendForm.msg_type === 2" label="业务类型">
              <el-input v-model="sendForm.biz_type" placeholder="如 order / audit（可选）" />
            </el-form-item>
            <el-form-item v-if="needReceivers" label="接收人">
              <el-select
                v-model="sendForm.account_ids"
                multiple
                filterable
                clearable
                placeholder="业务通知可不选（默认全员）；指定用户必选"
                style="width: 100%"
              >
                <el-option
                  v-for="a in accountOptions"
                  :key="a.id"
                  :label="`${a.user_name} (${a.account_number})`"
                  :value="a.id"
                />
              </el-select>
            </el-form-item>
            <el-form-item label="标题" prop="title">
              <el-input v-model="sendForm.title" maxlength="120" show-word-limit />
            </el-form-item>
            <el-form-item label="内容" prop="content">
              <el-input v-model="sendForm.content" type="textarea" :rows="6" maxlength="2000" show-word-limit />
            </el-form-item>
            <el-form-item>
              <el-button type="primary" @click="submitSend">发送</el-button>
            </el-form-item>
          </el-form>
        </el-tab-pane>

        <el-tab-pane v-if="can('query')" label="发送记录" name="sent">
          <div class="filter-panel">
            <el-form :model="sentQuery" inline class="filter-form" label-width="64px" @submit.prevent="() => { sentQuery.page = 1; loadSent() }">
              <el-form-item label="类型">
                <el-select v-model="sentQuery.msg_type" clearable placeholder="全部" style="width: 140px">
                  <el-option :value="1" label="系统广播" />
                  <el-option :value="2" label="业务通知" />
                  <el-option :value="3" label="指定用户" />
                </el-select>
              </el-form-item>
              <el-form-item label="标题">
                <el-input v-model="sentQuery.title" clearable placeholder="模糊搜索" />
              </el-form-item>
              <el-form-item label-width="0">
                <el-button type="primary" native-type="submit">查询</el-button>
                <el-button @click="loadSent">刷新</el-button>
              </el-form-item>
            </el-form>
          </div>
          <el-table v-loading="loading" :data="sentList" border stripe>
            <el-table-column prop="id" label="ID" width="70" />
            <el-table-column label="类型" width="100">
              <template #default="{ row }">{{ msgTypeLabel[row.msg_type] || row.msg_type }}</template>
            </el-table-column>
            <el-table-column prop="title" label="标题" min-width="160" show-overflow-tooltip />
            <el-table-column prop="sender_name" label="发送人" width="120" />
            <el-table-column label="已读/接收" width="110">
              <template #default="{ row }">{{ row.read_count }}/{{ row.receiver_count }}</template>
            </el-table-column>
            <el-table-column label="时间" width="170">
              <template #default="{ row }">{{ formatUnix(row.create_time) }}</template>
            </el-table-column>
          </el-table>
          <div class="pager">
            <el-pagination
              v-model:current-page="sentQuery.page"
              v-model:page-size="sentQuery.size"
              layout="total, prev, pager, next"
              :total="sentTotal"
              @current-change="loadSent"
            />
          </div>
        </el-tab-pane>
      </el-tabs>
    </el-card>

    <el-dialog v-model="detailVisible" title="消息详情" width="560px" destroy-on-close>
      <template v-if="detail">
        <el-descriptions :column="1" border>
          <el-descriptions-item label="类型">{{ msgTypeLabel[detail.msg_type] }}</el-descriptions-item>
          <el-descriptions-item v-if="detail.biz_type" label="业务类型">{{ detail.biz_type }}</el-descriptions-item>
          <el-descriptions-item label="发送人">{{ detail.sender_name || '-' }}</el-descriptions-item>
          <el-descriptions-item label="时间">{{ formatUnix(detail.create_time) }}</el-descriptions-item>
          <el-descriptions-item label="标题">{{ detail.title }}</el-descriptions-item>
          <el-descriptions-item label="内容">
            <div class="detail-content">{{ detail.content }}</div>
          </el-descriptions-item>
        </el-descriptions>
      </template>
    </el-dialog>
  </div>
</template>

<style scoped>
.pager {
  margin-top: 16px;
  display: flex;
  justify-content: flex-end;
}
.detail-content {
  white-space: pre-wrap;
  line-height: 1.6;
}
</style>
