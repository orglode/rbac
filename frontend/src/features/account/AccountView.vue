<script setup lang="ts">
import { onMounted, reactive, ref } from 'vue'
import { ElMessage, ElMessageBox, type FormInstance, type FormRules } from 'element-plus'
import type { Account, Role } from '@/shared/types/api'
import { formatUnix } from '@/shared/utils/time'
import { useUserStore } from '@/stores/user'
import {
  assignAccountRole,
  createAccount,
  deleteAccount,
  fetchAccountList,
  removeAccountRole,
  updateAccount,
} from './api'
import { fetchAllRoles } from '@/features/role/api'

const user = useUserStore()
const path = '/system/account'
const can = (code: string) => user.can(path, code)

const loading = ref(false)
const list = ref<Account[]>([])
const total = ref(0)
const query = reactive({
  page: 1,
  size: 20,
  user_name: '',
  mobile: '',
  account_number: '',
  status: undefined as number | undefined,
})

const dialogVisible = ref(false)
const dialogTitle = ref('新增账号')
const formRef = ref<FormInstance>()
const form = reactive({
  id: 0,
  user_name: '',
  mobile: '',
  account_number: '',
  password: '',
  status: 1,
})

const rules: FormRules = {
  user_name: [{ required: true, message: '请输入姓名', trigger: 'blur' }],
  mobile: [{ required: true, message: '请输入手机号', trigger: 'blur' }],
  account_number: [{ required: true, message: '请输入账号', trigger: 'blur' }],
}

const roleDialog = ref(false)
const roleAccountId = ref(0)
const roleId = ref<number>()
const roles = ref<Role[]>([])

async function load() {
  loading.value = true
  try {
    const data = await fetchAccountList({
      page: query.page,
      size: query.size,
      user_name: query.user_name || undefined,
      mobile: query.mobile || undefined,
      account_number: query.account_number || undefined,
      status: query.status || undefined,
    })
    list.value = data?.list || []
    total.value = data?.total || 0
  } finally {
    loading.value = false
  }
}

function onSearch() {
  query.page = 1
  load()
}

function onReset() {
  query.user_name = ''
  query.mobile = ''
  query.account_number = ''
  query.status = undefined
  query.page = 1
  load()
}

function openCreate() {
  dialogTitle.value = '新增账号'
  Object.assign(form, {
    id: 0,
    user_name: '',
    mobile: '',
    account_number: '',
    password: '',
    status: 1,
  })
  dialogVisible.value = true
}

function openEdit(row: Account) {
  dialogTitle.value = '编辑账号'
  Object.assign(form, {
    id: row.id,
    user_name: row.user_name,
    mobile: row.mobile,
    account_number: row.account_number,
    password: '',
    status: row.status || 1,
  })
  dialogVisible.value = true
}

async function submit() {
  await formRef.value?.validate()
  if (form.id) {
    await updateAccount(form.id, { ...form })
    ElMessage.success('已更新')
  } else {
    if (!form.password) {
      ElMessage.warning('请填写密码')
      return
    }
    await createAccount({ ...form })
    ElMessage.success('已创建')
  }
  dialogVisible.value = false
  await load()
}

async function onDelete(row: Account) {
  await ElMessageBox.confirm(`确认删除账号「${row.account_number}」？`, '提示', { type: 'warning' })
  await deleteAccount(row.id)
  ElMessage.success('已删除')
  await load()
}

async function openRole(row: Account) {
  roleAccountId.value = row.id
  roleId.value = undefined
  roles.value = (await fetchAllRoles()) || []
  roleDialog.value = true
}

async function submitRole() {
  if (!roleId.value) {
    ElMessage.warning('请选择角色')
    return
  }
  await assignAccountRole({ account_id: roleAccountId.value, role_id: roleId.value })
  ElMessage.success('已分配角色')
  roleDialog.value = false
}

async function clearRole(row: Account) {
  await ElMessageBox.confirm(`移除账号「${row.account_number}」的角色？`, '提示', { type: 'warning' })
  await removeAccountRole({ account_id: row.id })
  ElMessage.success('已移除')
}

onMounted(load)
</script>

<template>
  <el-card shadow="never">
    <div class="filter-panel">
      <el-form :model="query" inline class="filter-form" label-width="64px" @submit.prevent="onSearch">
        <el-form-item label="姓名">
          <el-input v-model="query.user_name" clearable placeholder="请输入姓名" @keyup.enter="onSearch" />
        </el-form-item>
        <el-form-item label="账号">
          <el-input v-model="query.account_number" clearable placeholder="请输入账号" @keyup.enter="onSearch" />
        </el-form-item>
        <el-form-item label="手机号">
          <el-input v-model="query.mobile" clearable placeholder="请输入手机号" @keyup.enter="onSearch" />
        </el-form-item>
        <el-form-item label="状态">
          <el-select v-model="query.status" clearable placeholder="全部">
            <el-option :value="1" label="启用" />
            <el-option :value="2" label="禁用" />
          </el-select>
        </el-form-item>
        <el-form-item label-width="0">
          <el-button type="primary" native-type="submit">查询</el-button>
          <el-button @click="onReset">重置</el-button>
        </el-form-item>
      </el-form>
    </div>
    <div class="filter-actions">
      <el-button v-if="can('create')" type="primary" @click="openCreate">新增</el-button>
      <el-button @click="load">刷新</el-button>
    </div>

    <el-table v-loading="loading" :data="list" border stripe>
      <el-table-column prop="id" label="ID" width="70" />
      <el-table-column prop="user_name" label="姓名" min-width="100" />
      <el-table-column prop="account_number" label="账号" min-width="120" />
      <el-table-column prop="mobile" label="手机号" min-width="120" />
      <el-table-column label="状态" width="90">
        <template #default="{ row }">
          <el-tag :type="row.status === 1 ? 'success' : 'danger'" size="small">
            {{ row.status === 1 ? '启用' : '禁用' }}
          </el-tag>
        </template>
      </el-table-column>
      <el-table-column label="最后登录" min-width="160">
        <template #default="{ row }">{{ formatUnix(row.last_login_time) }}</template>
      </el-table-column>
      <el-table-column label="操作" width="260" fixed="right">
        <template #default="{ row }">
          <el-button v-if="can('update')" link type="primary" @click="openEdit(row as Account)">编辑</el-button>
          <el-button v-if="can('assign_role')" link type="primary" @click="openRole(row as Account)">分配角色</el-button>
          <el-button v-if="can('remove_role')" link type="warning" @click="clearRole(row as Account)">移除角色</el-button>
          <el-button v-if="can('delete')" link type="danger" @click="onDelete(row as Account)">删除</el-button>
        </template>
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

    <el-dialog v-model="dialogVisible" :title="dialogTitle" width="480px" destroy-on-close>
      <el-form ref="formRef" :model="form" :rules="rules" label-width="90px">
        <el-form-item label="姓名" prop="user_name">
          <el-input v-model="form.user_name" />
        </el-form-item>
        <el-form-item label="账号" prop="account_number">
          <el-input v-model="form.account_number" />
        </el-form-item>
        <el-form-item label="手机号" prop="mobile">
          <el-input v-model="form.mobile" />
        </el-form-item>
        <el-form-item :label="form.id ? '新密码' : '密码'">
          <el-input v-model="form.password" type="password" show-password :placeholder="form.id ? '不填则不改' : ''" />
        </el-form-item>
        <el-form-item label="状态">
          <el-radio-group v-model="form.status">
            <el-radio :value="1">启用</el-radio>
            <el-radio :value="2">禁用</el-radio>
          </el-radio-group>
        </el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="dialogVisible = false">取消</el-button>
        <el-button type="primary" @click="submit">确定</el-button>
      </template>
    </el-dialog>

    <el-dialog v-model="roleDialog" title="分配角色" width="420px" destroy-on-close>
      <el-select v-model="roleId" placeholder="选择角色" style="width: 100%" filterable>
        <el-option v-for="r in roles" :key="r.id" :label="r.role_name" :value="r.id" />
      </el-select>
      <template #footer>
        <el-button @click="roleDialog = false">取消</el-button>
        <el-button type="primary" @click="submitRole">确定</el-button>
      </template>
    </el-dialog>
  </el-card>
</template>

<style scoped>
.pager {
  margin-top: 16px;
  display: flex;
  justify-content: flex-end;
}
</style>
