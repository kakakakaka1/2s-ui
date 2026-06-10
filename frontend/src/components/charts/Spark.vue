<template>
  <svg :width="width" :height="height" style="display: block;">
    <path :d="d" fill="none" :stroke="color" stroke-width="2" stroke-linecap="round" />
  </svg>
</template>

<script lang="ts" setup>
import { computed } from 'vue'
import { smoothPath } from './path'

const props = withDefaults(defineProps<{
  data: number[]
  color?: string
  width?: number
  height?: number
}>(), {
  color: 'var(--brand)',
  width: 110,
  height: 34,
})

const d = computed(() => {
  const max = Math.max(...props.data) || 1
  const min = Math.min(...props.data)
  const xStep = props.width / Math.max(1, props.data.length - 1)
  const pts: [number, number][] = props.data.map((v, i) => [
    i * xStep,
    props.height - 3 - ((v - min) / (max - min || 1)) * (props.height - 8),
  ])
  return smoothPath(pts)
})
</script>
