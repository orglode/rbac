<script setup lang="ts">
import { onMounted, reactive, ref } from 'vue'
import { ElMessage, ElMessageBox, type FormInstance, type FormRules } from 'element-plus'
import type { MenuPage, PageInfo, Role } from '@/shared/types/api'
import { formatUnix } from '@/shared/utils/time'
import { useUserStore } from '@/stores/user'
import { fetchAllPages } from '@/features/page/api'
import {
  assignRolePages,
  createRole,
  deleteRole,
  fetchRoleList,
  fetchRolePages,
  updateRole,
} from './api'

const user = useUserStore()
const path = '/system/role'
const can = (code: string) => user.can(path, code)

const loading = ref(false)
const list = ref<Role[]>([])
const total = ref(0)
const query = reactive({ page: 1, size: 20 })

const dialogVisible = ref(false)
const formRef = ref<FormInstance>()
const form = reactive({ id: 0, role_name: '' })
const rules: FormRules = {
  role_name: [{ required: true, message: '请输入角色名', trigger: 'blur' }],
}

const permVisible = ref(false)
const permRoleId = ref(0)
const treeRef = ref()
const treeData = ref<TreeNode[]>([])
const checkedKeys = ref<number[]>([])

interface TreeNode {
  id: number
  label: string
  children?: TreeNode[]
}

function pagesToTree(pages: MenuPage[]): TreeNode[] {
  const map = new Map<number, TreeNode>()
  const roots: TreeNode[] = []
  const sorted = [...pages].sort((a, b) => a.page_level - b.page_level || a.id - b.id)
  for (const p of sorted) {
    map.set(p.id, {
      id: p.id,
      label: p.page_type === 2 ? `${p.page_name} (${p.button_code || 'btn'})` : p.page_name,
      children: [],
    })
  }
  for (const p of sorted) {
    const node = map.get(p.id)!
    if (p.parent_id && map.has(p.parent_id)) {
      map.get(p.parent_id)!.children!.push(node)
    } else if (p.page_level === 1) {
      roots.push(node)
    } else {
      roots.push(node)
    }
  }
  const prune = (nodes: TreeNode[]) => {
    for (const n of nodes) {
      if (n.children?.length) prune(n.children)
      else delete n.children
    }
  }
  prune(roots)
  return roots
}

function collectChecked(nodes: PageInfo[], set: Set<number>) {
  for (const n of nodes) {
    set.add(n.page_id)
    if (n.child_page?.length) collectChecked(n.child_page, set)
    if (n.buttons?.length) collectChecked(n.buttons, set)
  }
}

async function load() {
  loading.value = true
  try {
    const data = await fetchRoleList({ page: query.page, size: query.size })
    list.value = (data?.list as Role[]) || []
    total.value = data?.total || 0
  } finally {
    loading.value = false
  }
}

function openCreate() {
  form.id = 0
  form.role_name = ''
  dialogVisible.value = true
}

function openEdit(row: Role) {
  form.id = row.id
  form.role_name = row.role_name
  dialogVisible.value = true
}

async function submit() {
  await formRef.value?.validate()
  if (form.id) {
    await updateRole(form.id, { role_name: form.role_name })
    ElMessage.success('已更新')
  } else {
    await createRole({ role_name: form.role_name })
    ElMessage.success('已创建')
  }
  dialogVisible.value = false
  await load()
}

async function onDelete(row: Role) {
  await ElMessageBox.confirm(`确认删除角色「${row.role_name}」？`, '提示', { type: 'warning' })
  await deleteRole(row.id)
  ElMessage.success('已删除')
  await load()
}

async function openPerm(row: Role) {
  permRoleId.value = row.id
  const all = (await fetchAllPages()) || []
  treeData.value = pagesToTree(all)
  const bound = await fetchRolePages(row.id)
  const set = new Set<number>()
  collectChecked(bound?.page || [], set)
  checkedKeys.value = [...set]
  permVisible.value = true
}

async function submitPerm() {
  const tree = treeRef.value
  const keys = [
    ...(tree?.getCheckedKeys(false) || []),
    ...(tree?.getHalfCheckedKeys() || []),
  ] as number[]
  const uniq = [...new Set(keys)]
  await assignRolePages({ role_id: permRoleId.value, page_ids: uniq.join(',') })
  ElMessage.success('权限已保存')
  permVisible.value = false
}

onMounted(load)
</script>

<template>
  <el-card shadow="never">
    <div class="toolbar">
      <el-button v-if="can('create')" type="primary" @click="openCreate">新增</el-button>
      <el-button @click="load">刷新</el-button>
    </div>

    <el-table v-loading="loading" :data="list" border stripe>
      <el-table-column prop="id" label="ID" width="70" />
      <el-table-column prop="role_name" label="角色名" min-width="160" />
      <el-table-column label="创建时间" min-width="160">
        <template #default="{ row }">{{ formatUnix(row.create_time) }}</template>
      </el-table-column>
      <el-table-column label="操作" width="220" fixed="right">
        <template #default="{ row }">
          <el-button v-if="can('update')" link type="primary" @click="openEdit(row as Role)">编辑</el-button>
          <el-button v-if="can('role_page_assign')" link type="primary" @click="openPerm(row as Role)">分配权限</el-button>
          <el-button v-if="can('delete')" link type="danger" @click="onDelete(row as Role)">删除</el-button>
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

    <el-dialog v-model="dialogVisible" :title="form.id ? '编辑角色' : '新增角色'" width="420px" destroy-on-close>
      <el-form ref="formRef" :model="form" :rules="rules" label-width="80px">
        <el-form-item label="角色名" prop="role_name">
          <el-input v-model="form.role_name" />
        </el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="dialogVisible = false">取消</el-button>
        <el-button type="primary" @click="submit">确定</el-button>
      </template>
    </el-dialog>

    <el-dialog v-model="permVisible" title="分配菜单与按钮" width="520px" destroy-on-close>
      <el-tree
        ref="treeRef"
        :data="treeData"
        node-key="id"
        show-checkbox
        default-expand-all
        :default-checked-keys="checkedKeys"
        :props="{ label: 'label', children: 'children' }"
      />
      <template #footer>
        <el-button @click="permVisible = false">取消</el-button>
        <el-button type="primary" @click="submitPerm">保存</el-button>
      </template>
    </el-dialog>
  </el-card>
</template>

<style scoped>
.toolbar {
  margin-bottom: 12px;
  display: flex;
  gap: 8px;
}
.pager {
  margin-top: 16px;
  display: flex;
  justify-content: flex-end;
}
</style>
