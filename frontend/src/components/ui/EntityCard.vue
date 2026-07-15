<template>
  <div class="card" style="padding: 0; overflow: hidden; display: flex; flex-direction: column;">
    <div style="padding: 14px 15px 11px; display: flex; align-items: center; gap: 11px;">
      <div
        :style="{
          width: '36px', height: '36px', borderRadius: '10px', flex: 'none',
          display: 'flex', alignItems: 'center', justifyContent: 'center',
          color, background: `color-mix(in srgb, ${color} 14%, transparent)`,
        }"
      ><Ico :name="icon" :size="18" /></div>
      <div style="flex: 1; min-width: 0;">
        <div style="font-weight: 700; font-size: 14px; white-space: nowrap; overflow: hidden; text-overflow: ellipsis;">{{ title }}</div>
        <div :style="{ fontSize: '11.5px', color, fontWeight: 600 }">{{ type }}</div>
      </div>
      <slot name="chip" />
    </div>
    <div style="padding: 0 15px 12px; display: flex; flex-direction: column;">
      <div v-for="(r, i) in rows" :key="i" class="kv-row">
        <span class="k">{{ r.k }}</span>
        <span class="v" :class="{ mono: r.mono }" :style="r.color ? { color: r.color } : undefined">{{ r.v }}</span>
      </div>
    </div>
    <div style="display: flex; border-top: 1px solid var(--line); margin-top: auto;">
      <slot name="actions" />
    </div>
  </div>
</template>

<script lang="ts" setup>
import Ico from './Ico.vue'

export interface EntityRow {
  k: string
  v: string | number
  mono?: boolean
  color?: string
}

withDefaults(defineProps<{
  title: string
  type?: string
  color?: string
  icon?: string
  rows: EntityRow[]
}>(), {
  color: 'var(--brand)',
  icon: 'outbound',
})
</script>
