<template>
  <svg :width="size" :height="size" style="transform: rotate(-90deg);">
    <circle :cx="size / 2" :cy="size / 2" :r="r" fill="none" stroke="var(--track)" :stroke-width="thickness" />
    <circle
      v-for="(s, i) in segments"
      :key="i"
      :cx="size / 2" :cy="size / 2" :r="r"
      fill="none"
      :stroke="s.color"
      :stroke-width="thickness"
      :stroke-dasharray="s.dash"
      :stroke-dashoffset="s.off"
      stroke-linecap="round"
      style="transition: all .8s var(--ease);"
    />
  </svg>
</template>

<script lang="ts" setup>
import { computed } from 'vue'

const props = withDefaults(defineProps<{
  data: { value: number; color: string }[]
  size?: number
  thickness?: number
}>(), {
  size: 150,
  thickness: 18,
})

const r = computed(() => props.size / 2 - props.thickness / 2 - 2)

const segments = computed(() => {
  const total = props.data.reduce((a, b) => a + b.value, 0) || 1
  const c = 2 * Math.PI * r.value
  let acc = 0
  return props.data.map((d) => {
    const frac = d.value / total
    const seg = { color: d.color, dash: `${frac * c} ${c}`, off: -acc * c }
    acc += frac
    return seg
  })
})
</script>
