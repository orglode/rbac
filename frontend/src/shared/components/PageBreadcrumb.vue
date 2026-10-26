<script setup lang="ts">
import { useBreadcrumbs } from '@/shared/composables/useBreadcrumbs'

const { breadcrumbs, onCrumbClick } = useBreadcrumbs()
</script>

<template>
  <el-breadcrumb class="page-breadcrumb" separator="/">
    <el-breadcrumb-item v-for="(item, idx) in breadcrumbs" :key="`${item.title}-${idx}`">
      <a
        v-if="item.path && idx < breadcrumbs.length - 1"
        href="javascript:;"
        @click.prevent="onCrumbClick(item)"
      >{{ item.title }}</a>
      <span v-else class="crumb-current">{{ item.title }}</span>
    </el-breadcrumb-item>
  </el-breadcrumb>
</template>

<style scoped>
.page-breadcrumb {
  margin-bottom: 12px;
}
.crumb-current {
  color: #303133;
  font-weight: 500;
}
:deep(.el-breadcrumb__inner a) {
  font-weight: 400;
  color: #909399;
}
:deep(.el-breadcrumb__inner a:hover) {
  color: #409eff;
}
</style>
