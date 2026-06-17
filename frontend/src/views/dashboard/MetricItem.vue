<template>
  <div style="display: flex; align-items: center; gap: 9px; min-width: 0;">
    <div
      :style="{
        width: '30px', height: '30px', borderRadius: '9px', flex: 'none',
        display: 'flex', alignItems: 'center', justifyContent: 'center',
        color: accent, background: `color-mix(in srgb, ${accent} 14%, transparent)`,
      }"
    ><Ico :name="icon" :size="16" /></div>
    <div style="min-width: 0;">
      <div style="font-size: 11px; color: var(--text-3);">{{ label }}</div>
      <div class="mono" style="font-size: 17px; font-weight: 700; line-height: 1.15; white-space: nowrap; overflow: hidden; text-overflow: ellipsis;">{{ value }}</div>
    </div>
  </div>
</template>

<script lang="ts" setup>
import Ico from '@/components/ui/Ico.vue'

defineProps<{
  icon: string
  label: string
  value: string
  accent: string
}>()
</script>
