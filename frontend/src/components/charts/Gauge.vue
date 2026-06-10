<template>
  <div style="display: flex; flex-direction: column; align-items: center;">
    <div :style="{ position: 'relative', width: size + 'px', height: size * 0.82 + 'px' }">
      <svg :width="size" :height="size" :style="{ position: 'absolute', top: -size * 0.06 + 'px' }">
        <path :d="arc(START, END)" fill="none" stroke="var(--track)" stroke-width="9" stroke-linecap="round" />
        <path
          :d="arc(START, vAng)"
          fill="none"
          :stroke="color"
          stroke-width="9"
          stroke-linecap="round"
          :style="{ transition: 'all .9s var(--ease)', filter: `drop-shadow(0 2px 6px color-mix(in srgb, ${color} 50%, transparent))` }"
        />
      </svg>
      <div :style="{ position: 'absolute', inset: 0, display: 'flex', flexDirection: 'column', alignItems: 'center', justifyContent: 'center', paddingTop: size * 0.06 + 'px' }">
        <div class="mono" :style="{ fontSize: size * 0.24 + 'px', fontWeight: 700, letterSpacing: '-.02em', lineHeight: 1 }">
          {{ value }}<span :style="{ fontSize: size * 0.12 + 'px', color: 'var(--text-3)', marginLeft: '1px' }">%</span>
        </div>
      </div>
    </div>
    <div style="text-align: center; margin-top: -2px;">
      <div style="font-size: 13px; font-weight: 700;">{{ label }}</div>
      <div class="mono" style="font-size: 11px; color: var(--text-3); margin-top: 2px;">{{ sub }}</div>
    </div>
  </div>
</template>

<script lang="ts" setup>
import { computed } from 'vue'

const props = withDefaults(defineProps<{
  value: number
  color?: string
  label?: string
  sub?: string
  size?: number
}>(), {
  color: 'var(--brand)',
  label: '',
  sub: '',
  size: 132,
})

const START = -220
const END = 40

const r = computed(() => props.size / 2 - 12)
const c = computed(() => props.size / 2)

const toXY = (ang: number): [number, number] => {
  const a = (ang * Math.PI) / 180
  return [c.value + r.value * Math.cos(a), c.value + r.value * Math.sin(a)]
}
const arc = (a0: number, a1: number) => {
  const [x0, y0] = toXY(a0), [x1, y1] = toXY(a1)
  const large = a1 - a0 > 180 ? 1 : 0
  return `M ${x0} ${y0} A ${r.value} ${r.value} 0 ${large} 1 ${x1} ${y1}`
}
const vAng = computed(() => START + ((END - START) * Math.min(100, Math.max(0, props.value))) / 100)
</script>
