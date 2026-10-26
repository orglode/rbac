<script setup lang="ts">
import { onMounted, reactive, ref } from 'vue'
import { ElMessage, ElMessageBox, type FormInstance, type FormRules } from 'element-plus'
import type { MenuPage, PageInfo } from '@/shared/types/api'
import IconPicker from '@/shared/components/IconPicker.vue'
import MenuIcon from '@/shared/components/MenuIcon.vue'
import PageBreadcrumb from '@/shared/components/PageBreadcrumb.vue'
import { useUserStore } from '@/stores/user'
import { createPage, deletePage, fetchAllPages, fetchPageList, updatePage } from './api'

const user = useUserStore()
const path = '/system/page'
const can = (code: string) => user.can(path, code)

interface TreeRow {
  id: number
  page_name: string
  icon?: string
  page_level: number
  page_type: number
  page_path: string
  button_code?: string
  api_method?: string
  api_path?: string
  sort: number
  children?: TreeRow[]
}

const loading = ref(false)
const list = ref<TreeRow[]>([])
const total = ref(0)
const query = reactive({
  page: 1,
  size: 20,
  page_name: '',
  api: '',
})
const parents = ref<MenuPage[]>([])

const dialogVisible = ref(false)
const formRef = ref<FormInstance>()
const form = reactive({
  id: 0,
  page_name: '',
  icon: '',
  url: '',
  parent_id: 0,
  page_level: 1,
  page_type: 1,
  button_code: '',
  api_method: '',
  api_path: '',
  api_prefix: '',
  sort: 0,
})

const rules: FormRules = {
  page_name: [{ required: true, message: '请输入名称', trigger: 'blur' }],
  page_level: [{ required: true, message: '请选择层级', trigger: 'change' }],
}

const levelLabel: Record<number, string> = { 1: '一级菜单', 2: '二级菜单', 3: '按钮' }
const typeLabel: Record<number, string> = { 1: '菜单', 2: '按钮' }

function toTreeRow(node: PageInfo): TreeRow {
  const children: TreeRow[] = []
  for (const c of node.child_page || []) {
    children.push(toTreeRow(c))
  }
  for (const b of node.buttons || []) {
    children.push(toTreeRow(b))
  }
  return {
    id: node.page_id,
    page_name: node.page_name,
    icon: node.icon,
    page_level: node.page_level,
    page_type: node.page_type,
    page_path: node.page_path,
    button_code: node.button_code,
    api_method: node.api_method,
    api_path: node.api_path,
    sort: node.sort ?? 0,
    children: children.length ? children : undefined,
  }
}

async function load() {
  loading.value = true
  try {
    const data = await fetchPageList({
      page: query.page,
      size: query.size,
      page_name: query.page_name || undefined,
      api: query.api || undefined,
    })
    const roots = data?.list || []
    list.value = roots.map(toTreeRow)
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
  query.page_name = ''
  query.api = ''
  query.page = 1
  load()
}

async function loadParents() {
  parents.value = (await fetchAllPages()) || []
}

function openCreate() {
  Object.assign(form, {
    id: 0,
    page_name: '',
    icon: '',
    url: '',
    parent_id: 0,
    page_level: 1,
    page_type: 1,
    button_code: '',
    api_method: '',
    api_path: '',
    api_prefix: '',
    sort: 0,
  })
  loadParents()
  dialogVisible.value = true
}

function openEdit(row: TreeRow) {
  Object.assign(form, {
    id: row.id,
    page_name: row.page_name,
    icon: row.icon || '',
    url: row.page_path || '',
    parent_id: 0,
    page_level: row.page_level,
    page_type: row.page_type,
    button_code: row.button_code || '',
    api_method: row.api_method || '',
    api_path: row.api_path || '',
    api_prefix: '',
    sort: row.sort ?? 0,
  })
  // 从全量列表回填 parent_id
  loadParents().then(() => {
    const self = parents.value.find((p) => p.id === row.id)
    if (self) {
      form.parent_id = self.parent_id
      form.sort = self.sort ?? row.sort ?? 0
      form.icon = self.icon || row.icon || ''
    }
  })
  dialogVisible.value = true
}

function onLevelChange(level: number) {
  if (level === 3) {
    form.page_type = 2
  } else {
    form.page_type = 1
    form.button_code = ''
    form.api_method = ''
    form.api_path = ''
  }
}

async function submit() {
  await formRef.value?.validate()
  const payload: Record<string, unknown> = {
    page_name: form.page_name,
    icon: form.page_level <= 2 ? form.icon : '',
    url: form.url,
    parent_id: form.parent_id,
    page_level: form.page_level,
    page_type: form.page_type,
    button_code: form.button_code,
    api_method: form.api_method,
    api_path: form.api_path,
    sort: form.sort,
  }
  if (!form.id && form.page_level === 2 && form.api_prefix) {
    payload.api_prefix = form.api_prefix
  }
  if (form.id) {
    await updatePage(form.id, payload)
    ElMessage.success('已更新')
  } else {
    await createPage(payload)
    ElMessage.success('已创建')
  }
  dialogVisible.value = false
  await load()
}

async function onDelete(row: TreeRow) {
  await ElMessageBox.confirm(`确认删除「${row.page_name}」及其子节点？`, '提示', { type: 'warning' })
  await deletePage(row.id)
  ElMessage.success('已删除')
  await load()
}

onMounted(load)
</script>

<template>
  <div>
    <PageBreadcrumb />
    <el-card shadow="never">
    <div class="filter-panel">
      <el-form :model="query" inline class="filter-form" label-width="64px" @submit.prevent="onSearch">
        <el-form-item label="名称">
          <el-input v-model="query.page_name" clearable placeholder="请输入名称" @keyup.enter="onSearch" />
        </el-form-item>
        <el-form-item label="API">
          <el-input
            v-model="query.api"
            clearable
            placeholder="method / path 模糊匹配"
            style="width: 240px"
            @keyup.enter="onSearch"
          />
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

    <el-table
      v-loading="loading"
      :data="list"
      border
      row-key="id"
      :tree-props="{ children: 'children' }"
    >
      <el-table-column prop="page_name" label="名称" min-width="180" />
      <el-table-column label="图标" width="90">
        <template #default="{ row }">
          <div v-if="row.page_type === 1 && row.icon" class="icon-cell">
            <MenuIcon :name="row.icon" :size="16" />
            <span class="icon-code">{{ row.icon }}</span>
          </div>
          <span v-else>-</span>
        </template>
      </el-table-column>
      <el-table-column label="层级" width="100">
        <template #default="{ row }">{{ levelLabel[row.page_level] || row.page_level }}</template>
      </el-table-column>
      <el-table-column label="类型" width="80">
        <template #default="{ row }">{{ typeLabel[row.page_type] || row.page_type }}</template>
      </el-table-column>
      <el-table-column prop="page_path" label="页面路由" min-width="140" />
      <el-table-column prop="sort" label="排序" width="80" />
      <el-table-column prop="button_code" label="按钮编码" width="110" />
      <el-table-column label="API" min-width="220">
        <template #default="{ row }">
          <span v-if="row.api_path">{{ row.api_method }} {{ row.api_path }}</span>
          <span v-else>-</span>
        </template>
      </el-table-column>
      <el-table-column label="操作" width="140" fixed="right">
        <template #default="{ row }">
          <el-button v-if="can('update')" link type="primary" @click="openEdit(row as TreeRow)">编辑</el-button>
          <el-button v-if="can('delete')" link type="danger" @click="onDelete(row as TreeRow)">删除</el-button>
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

    <el-dialog v-model="dialogVisible" :title="form.id ? '编辑页面' : '新增页面'" width="640px" destroy-on-close>
      <el-form ref="formRef" :model="form" :rules="rules" label-width="100px">
        <el-form-item label="名称" prop="page_name">
          <el-input v-model="form.page_name" />
        </el-form-item>
        <el-form-item label="层级" prop="page_level">
          <el-select v-model="form.page_level" style="width: 100%" @change="onLevelChange">
            <el-option :value="1" label="一级菜单" />
            <el-option :value="2" label="二级菜单" />
            <el-option :value="3" label="按钮" />
          </el-select>
        </el-form-item>
        <el-form-item v-if="form.page_level <= 2" label="菜单图标">
          <IconPicker v-model="form.icon" />
        </el-form-item>
        <el-form-item v-if="form.page_level > 1" label="父级">
          <el-select v-model="form.parent_id" style="width: 100%" filterable>
            <el-option
              v-for="p in parents.filter((x) => x.page_type === 1 && x.page_level < form.page_level)"
              :key="p.id"
              :label="`${p.page_name} (#${p.id})`"
              :value="p.id"
            />
          </el-select>
        </el-form-item>
        <el-form-item v-if="form.page_type === 1" label="页面路由">
          <el-input v-model="form.url" placeholder="/system/xxx" />
        </el-form-item>
        <el-form-item label="排序">
          <el-input-number v-model="form.sort" :step="1" controls-position="right" style="width: 100%" />
          <div class="form-tip">数字越大越靠前；相同则按创建时间倒序</div>
        </el-form-item>
        <el-form-item v-if="!form.id && form.page_level === 2" label="API 前缀">
          <el-input v-model="form.api_prefix" placeholder="/api/admin/xxx（自动生成 CRUD 按钮）" />
        </el-form-item>
        <template v-if="form.page_level === 3">
          <el-form-item label="按钮编码">
            <el-input v-model="form.button_code" placeholder="query/create/update/delete" />
          </el-form-item>
          <el-form-item label="HTTP 方法">
            <el-select v-model="form.api_method" style="width: 100%">
              <el-option value="GET" />
              <el-option value="POST" />
              <el-option value="PUT" />
              <el-option value="DELETE" />
            </el-select>
          </el-form-item>
          <el-form-item label="API 路径">
            <el-input v-model="form.api_path" placeholder="/api/admin/xxx" />
          </el-form-item>
        </template>
      </el-form>
      <template #footer>
        <el-button @click="dialogVisible = false">取消</el-button>
        <el-button type="primary" @click="submit">确定</el-button>
      </template>
    </el-dialog>
  </el-card>
  </div>
</template>

<style scoped>
.pager {
  margin-top: 16px;
  display: flex;
  justify-content: flex-end;
}
.form-tip {
  margin-top: 4px;
  color: #909399;
  font-size: 12px;
  line-height: 1.4;
}
.icon-cell {
  display: inline-flex;
  align-items: center;
  gap: 2px;
}
.icon-code {
  color: #909399;
  font-size: 12px;
}
.icon-cell :deep(.menu-icon) {
  margin-right: 2px;
}
</style>
