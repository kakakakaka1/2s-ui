<template>
  <div class="card stat-card">
    <div style="display: flex; align-items: center; gap: 10px;">
      <div
        :style="{
          width: '34px', height: '34px', borderRadius: '9px', flex: 'none',
          display: 'flex', alignItems: 'center', justifyContent: 'center',
          color: accent, background: `color-mix(in srgb, ${accent} 14%, transparent)`,
        }"
      ><Ico :name="icon" :size="18" /></div>
      <div style="font-size: 12.5px; font-weight: 600; color: var(--text-2);">{{ label }}</div>
      <Chip v-if="trend" color="emerald" dot style="margin-inline-start: auto;">{{ trend }}</Chip>
    </div>
    <div style="display: flex; align-items: flex-end; justify-content: space-between; gap: 8px;">
      <div style="min-width: 0;">
        <div class="mono stat-value">{{ value }}</div>
        <div v-if="sub" style="font-size: 11.5px; color: var(--text-3); margin-top: 6px; white-space: nowrap; overflow: hidden; text-overflow: ellipsis;">{{ sub }}</div>
      </div>
      <div v-if="spark && spark.length > 1" class="stat-spark" style="margin-bottom: 2px; flex: none;">
        <Spark :data="spark" :color="accent" :width="88" />
      </div>
    </div>
  </div>
</template>

<script lang="ts" setup>
import Ico from '@/components/ui/Ico.vue'
import Chip from '@/components/ui/Chip.vue'
import Spark from '@/components/charts/Spark.vue'

defineProps<{
  icon: string
  label: string
  value: string
  sub?: string
  accent: string
  spark?: number[]
  trend?: string
}>()
</script>

<style scoped>
.stat-card { padding: 16px 18px; display: flex; flex-direction: column; gap: 10px; }
.stat-value { font-size: 24px; font-weight: 700; letter-spacing: -.03em; line-height: 1; white-space: nowrap; }
@media (max-width: 480px) {
  .stat-value { font-size: 20px; }
  .stat-spark { display: none; }
}
</style>
