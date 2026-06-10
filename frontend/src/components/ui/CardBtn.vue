<template>
  <button type="button" class="card-btn" :class="{ border, danger, 'icon-only': !label }" :title="title">
    <Ico :name="icon" :size="16" />{{ label }}
  </button>
</template>

<script lang="ts" setup>
import Ico from './Ico.vue'

defineProps<{
  icon: string
  label?: string
  border?: boolean
  danger?: boolean
  title?: string
}>()
</script>

<style scoped>
.card-btn {
  flex: 1; height: 44px;
  border: none; background: transparent; cursor: pointer;
  color: var(--text-2);
  display: inline-flex; align-items: center; justify-content: center; gap: 7px;
  font-family: var(--font-ui); font-size: 12.5px; font-weight: 600;
  transition: all .14s;
}
.card-btn.icon-only { flex: none; width: 46px; }
.card-btn.border { border-inline-start: 1px solid var(--line); }
.card-btn.danger { color: var(--text-3); }
.card-btn:hover { background: var(--surface-3); color: var(--text); }
.card-btn.danger:hover {
  background: color-mix(in srgb, var(--rose) 12%, transparent);
  color: var(--rose);
}
</style>
