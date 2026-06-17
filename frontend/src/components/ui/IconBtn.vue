<template>
  <button type="button" class="row-ico-btn" :class="{ danger }" :title="title">
    <Ico :name="name" :size="16" />
  </button>
</template>

<script lang="ts" setup>
import Ico from './Ico.vue'

defineProps<{
  name: string
  title?: string
  danger?: boolean
}>()
</script>

<style scoped>
.row-ico-btn {
  width: 32px; height: 32px; flex: none;
  border-radius: 8px; border: none;
  background: transparent; cursor: pointer;
  color: var(--text-3);
  display: inline-flex; align-items: center; justify-content: center;
  transition: all .14s;
}
.row-ico-btn:hover { background: var(--surface-3); color: var(--text); }
.row-ico-btn.danger:hover {
  background: color-mix(in srgb, var(--rose) 14%, transparent);
  color: var(--rose);
}
</style>
