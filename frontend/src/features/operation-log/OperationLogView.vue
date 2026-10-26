<script setup lang="ts">
import { onMounted, reactive, ref } from 'vue'
import type { OperationLog } from '@/shared/types/api'
import { formatUnix } from '@/shared/utils/time'
import PageBreadcrumb from '@/shared/components/PageBreadcrumb.vue'
import { fetchOperationLogs } from './api'

const moduleOptions = [
  { value: '0', module: 0, label: '全部' },
  { value: '1', module: 1, label: '角色变更' },
  { value: '2', module: 2, label: '页面权限' },
  { value: '3', module: 3, label: '账号角色' },
  { value: '4', module: 4, label: '页面资源' },
]

const actionOptions = [
  { value: 1, label: '新增' },
  { value: 2, label: '修改' },
  { value: 3, label: '删除' },
  { value: 4, label: '分配' },
  { value: 5, label: '移除' },
]

const moduleLabelMap: Record<number, string> = {
  1: '角色变更',
  2: '页面权限',
  3: '账号角色',
  4: '页面资源',
}

const actionLabelMap: Record<number, string> = {
  1: '新增',
  2: '修改',
  3: '删除',
  4: '分配',
  5: '移除',
}

const actionTagType: Record<number, 'success' | 'warning' | 'danger' | 'primary' | 'info'> = {
  1: 'success',
  2: 'warning',
  3: 'danger',
  4: 'primary',
  5: 'info',
}

const loading = ref(false)
const list = ref<OperationLog[]>([])
const total = ref(0)
const activeModule = ref('0')
const query = reactive({
  page: 1,
  size: 20,
  action: undefined as number | undefined,
  operator_account: '',
  target_name: '',
})

async function load() {
  if (loading.value) return
  loading.value = true
  try {
    const params: Record<string, unknown> = {
      page: query.page,
      size: query.size,
    }
    const module = Number(activeModule.value) || 0
    if (module > 0) params.module = module
    if (query.action) params.action = query.action
    if (query.operator_account.trim()) params.operator_account = query.operator_account.trim()
    if (query.target_name.trim()) params.target_name = query.target_name.trim()
    const data = await fetchOperationLogs(params as Parameters<typeof fetchOperationLogs>[0])
    list.value = Array.isArray(data?.list) ? (data.list as OperationLog[]) : []
    total.value = Number(data?.total) || 0
  } catch {
    list.value = []
    total.value = 0
  } finally {
    loading.value = false
  }
}

function onSearch() {
  query.page = 1
  void load()
}

function onTabChange() {
  query.page = 1
  void load()
}

function onPageChange(page: number) {
  query.page = page
  void load()
}

onMounted(() => {
  void load()
})
</script>

<template>
  <div>
    <PageBreadcrumb />
    <el-card shadow="never">
    <el-tabs v-model="activeModule" class="module-tabs" @tab-change="onTabChange">
      <el-tab-pane
        v-for="item in moduleOptions"
        :key="item.value"
        :label="item.label"
        :name="item.value"
      />
    </el-tabs>

    <div class="filter-panel">
      <el-form :model="query" inline class="filter-form" label-width="72px" @submit.prevent="onSearch">
        <el-form-item label="操作">
          <el-select v-model="query.action" clearable placeholder="全部" style="width: 120px">
            <el-option
              v-for="item in actionOptions"
              :key="item.value"
              :value="item.value"
              :label="item.label"
            />
          </el-select>
        </el-form-item>
        <el-form-item label="操作人">
          <el-input
            v-model="query.operator_account"
            clearable
            placeholder="请输入账号"
            @keyup.enter="onSearch"
          />
        </el-form-item>
        <el-form-item label="目标">
          <el-input
            v-model="query.target_name"
            clearable
            placeholder="目标名称"
            @keyup.enter="onSearch"
          />
        </el-form-item>
        <el-form-item label-width="0">
          <el-button type="primary" native-type="submit">查询</el-button>
          <el-button :disabled="loading" @click="load">刷新</el-button>
        </el-form-item>
      </el-form>
    </div>

    <el-table v-loading="loading" :data="list" border stripe empty-text="暂无日志">
      <el-table-column prop="id" label="ID" width="70" />
      <el-table-column label="模块" width="110">
        <template #default="{ row }">
          {{ moduleLabelMap[row.module] || row.module }}
        </template>
      </el-table-column>
      <el-table-column label="操作" width="90">
        <template #default="{ row }">
          <el-tag :type="actionTagType[row.action] || 'info'" size="small">
            {{ actionLabelMap[row.action] || row.action }}
          </el-tag>
        </template>
      </el-table-column>
      <el-table-column prop="target_name" label="目标" min-width="120" show-overflow-tooltip />
      <el-table-column prop="content" label="说明" min-width="240" show-overflow-tooltip />
      <el-table-column label="操作人" min-width="150">
        <template #default="{ row }">
          <span>{{ row.operator_name || '-' }}</span>
          <span v-if="row.operator_account" class="muted">({{ row.operator_account }})</span>
        </template>
      </el-table-column>
      <el-table-column label="时间" width="170">
        <template #default="{ row }">{{ formatUnix(row.create_time) }}</template>
      </el-table-column>
    </el-table>

    <div class="pager">
      <el-pagination
        :current-page="query.page"
        :page-size="query.size"
        layout="total, prev, pager, next"
        :total="total"
        @current-change="onPageChange"
      />
    </div>
  </el-card>
  </div>
</template>

<style scoped>
.module-tabs {
  margin-bottom: 4px;
}
.pager {
  margin-top: 16px;
  display: flex;
  justify-content: flex-end;
}
.muted {
  color: var(--el-text-color-secondary);
  margin-left: 4px;
}
</style>
