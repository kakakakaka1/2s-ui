<template>
  <div style="display: flex; align-items: center; gap: 8px;">
    <span
      :style="{
        width: '16px', height: '3px', borderRadius: '2px',
        background: dashed ? 'none' : color,
        borderTop: dashed ? `2px dashed ${color}` : 'none',
      }"
    />
    <div>
      <div style="font-size: 11.5px; color: var(--text-3);">{{ label }}</div>
      <div class="mono" style="font-size: 14px; font-weight: 700;" dir="ltr">{{ value }}</div>
    </div>
  </div>
</template>

<script lang="ts" setup>
defineProps<{
  color: string
  label: string
  value: string
  dashed?: boolean
}>()
</script>
