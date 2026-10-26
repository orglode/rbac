<script setup lang="ts">
import { onMounted, reactive, ref } from 'vue'
import type { LoginLog } from '@/shared/types/api'
import { formatUnix } from '@/shared/utils/time'
import { fetchLoginLogs } from './api'

const loading = ref(false)
const list = ref<LoginLog[]>([])
const total = ref(0)
const query = reactive({
  page: 1,
  size: 20,
  account_number: '',
  status: undefined as number | undefined,
})

async function load() {
  loading.value = true
  try {
    const params: Record<string, unknown> = {
      page: query.page,
      size: query.size,
    }
    if (query.account_number) params.account_number = query.account_number
    if (query.status) params.status = query.status
    const data = await fetchLoginLogs(params as Parameters<typeof fetchLoginLogs>[0])
    list.value = (data?.list as LoginLog[]) || []
    total.value = data?.total || 0
  } finally {
    loading.value = false
  }
}

function onSearch() {
  query.page = 1
  load()
}

onMounted(load)
</script>

<template>
  <el-card shadow="never">
    <div class="filter-panel">
      <el-form :model="query" inline class="filter-form" label-width="64px" @submit.prevent="onSearch">
        <el-form-item label="账号">
          <el-input
            v-model="query.account_number"
            clearable
            placeholder="请输入账号"
            @keyup.enter="onSearch"
          />
        </el-form-item>
        <el-form-item label="状态">
          <el-select v-model="query.status" clearable placeholder="全部">
            <el-option :value="1" label="成功" />
            <el-option :value="2" label="失败" />
          </el-select>
        </el-form-item>
        <el-form-item label-width="0">
          <el-button type="primary" native-type="submit">查询</el-button>
          <el-button @click="load">刷新</el-button>
        </el-form-item>
      </el-form>
    </div>

    <el-table v-loading="loading" :data="list" border stripe>
      <el-table-column prop="id" label="ID" width="70" />
      <el-table-column prop="account_number" label="账号" min-width="120" />
      <el-table-column prop="login_ip" label="IP" min-width="120" />
      <el-table-column label="状态" width="90">
        <template #default="{ row }">
          <el-tag :type="row.status === 1 ? 'success' : 'danger'" size="small">
            {{ row.status === 1 ? '成功' : '失败' }}
          </el-tag>
        </template>
      </el-table-column>
      <el-table-column prop="message" label="说明" min-width="140" />
      <el-table-column prop="user_agent" label="UA" min-width="200" show-overflow-tooltip />
      <el-table-column label="时间" width="170">
        <template #default="{ row }">{{ formatUnix(row.create_time) }}</template>
      </el-table-column>
    </el-table>

    <div class="pager">
      <el-pagination
        v-model:current-page="query.page"
        v-model:page-size="query.size"
        layout="total, prev, pager, next"
        :total="total"
        @current-change="load"
      />
    </div>
  </el-card>
</template>

<style scoped>
.pager {
  margin-top: 16px;
  display: flex;
  justify-content: flex-end;
}
</style>
