<script setup lang="ts">
import { computed } from 'vue'
import { resolveMenuIcon } from '@/shared/icons/menu-icons'

const props = withDefaults(
  defineProps<{
    name?: string
    fallback?: string
    size?: number | string
  }>(),
  {
    name: '',
    fallback: 'Menu',
    size: 18,
  },
)

const icon = computed(() => resolveMenuIcon(props.name) || resolveMenuIcon(props.fallback))
</script>

<template>
  <el-icon v-if="icon" :size="size" class="menu-icon">
    <component :is="icon" />
  </el-icon>
</template>

<style scoped>
.menu-icon {
  margin-right: 6px;
  vertical-align: middle;
}
</style>
