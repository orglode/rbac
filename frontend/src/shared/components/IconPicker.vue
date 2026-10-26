<script setup lang="ts">
import { computed, ref } from 'vue'
import { allIconNames, preferredIconNames, resolveMenuIcon } from '@/shared/icons/menu-icons'

const model = defineModel<string>({ default: '' })

const keyword = ref('')

const displayNames = computed(() => {
  const q = keyword.value.trim().toLowerCase()
  const source = q
    ? allIconNames.filter((n) => n.toLowerCase().includes(q))
    : preferredIconNames
  return source
})

function pick(name: string) {
  model.value = name === model.value ? '' : name
}
</script>

<template>
  <div class="icon-picker">
    <div class="picker-toolbar">
      <el-input
        v-model="keyword"
        clearable
        placeholder="搜索图标名，如 Setting / User"
        style="flex: 1"
      />
      <el-button v-if="model" link type="danger" @click="model = ''">清除</el-button>
    </div>
    <div v-if="model" class="current">
      <span>已选：</span>
      <el-icon :size="18"><component :is="resolveMenuIcon(model)" /></el-icon>
      <code>{{ model }}</code>
    </div>
    <div class="icon-grid">
      <button
        v-for="name in displayNames"
        :key="name"
        type="button"
        class="icon-item"
        :class="{ active: model === name }"
        :title="name"
        @click="pick(name)"
      >
        <el-icon :size="20"><component :is="resolveMenuIcon(name)" /></el-icon>
        <span class="icon-name">{{ name }}</span>
      </button>
    </div>
    <div class="hint">
      图标库：@element-plus/icons-vue（约 {{ allIconNames.length }} 个，MIT）。不搜时展示常用集，搜索可匹配全部。
    </div>
  </div>
</template>

<style scoped>
.icon-picker {
  width: 100%;
}
.picker-toolbar {
  display: flex;
  align-items: center;
  gap: 8px;
  margin-bottom: 10px;
}
.current {
  display: flex;
  align-items: center;
  gap: 6px;
  margin-bottom: 8px;
  color: #606266;
  font-size: 13px;
}
.current code {
  background: #f5f7fa;
  padding: 1px 6px;
  border-radius: 4px;
}
.icon-grid {
  display: grid;
  grid-template-columns: repeat(auto-fill, minmax(88px, 1fr));
  gap: 6px;
  max-height: 220px;
  overflow: auto;
  padding: 4px;
  border: 1px solid #ebeef5;
  border-radius: 6px;
  background: #fafafa;
}
.icon-item {
  display: flex;
  flex-direction: column;
  align-items: center;
  gap: 4px;
  padding: 8px 4px;
  border: 1px solid transparent;
  border-radius: 6px;
  background: #fff;
  cursor: pointer;
  color: #606266;
}
.icon-item:hover {
  border-color: #c6e2ff;
  color: #409eff;
}
.icon-item.active {
  border-color: #409eff;
  background: #ecf5ff;
  color: #409eff;
}
.icon-name {
  font-size: 11px;
  line-height: 1.2;
  max-width: 100%;
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
}
.hint {
  margin-top: 8px;
  color: #909399;
  font-size: 12px;
  line-height: 1.4;
}
</style>
