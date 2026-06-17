<template>
  <div>
    <div style="display: flex; align-items: center; gap: 6px; margin-bottom: 7px;">
      <span :style="{ width: '7px', height: '7px', borderRadius: '2px', background: color, flex: 'none' }" />
      <span style="font-size: 12px; color: var(--text-2); font-weight: 600;">{{ label }}</span>
      <span class="mono" style="margin-inline-start: auto; font-size: 15px; font-weight: 700;">{{ value }}%</span>
    </div>
    <BarMini :value="value" :color="color" />
    <div v-if="sub" class="mono" style="font-size: 10px; color: var(--text-3); margin-top: 5px;">{{ sub }}</div>
  </div>
</template>

<script lang="ts" setup>
import BarMini from '@/components/charts/BarMini.vue'

defineProps<{
  label: string
  value: number
  color: string
  sub?: string
}>()
</script>
