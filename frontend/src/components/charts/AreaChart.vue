<template>
  <svg width="100%" :viewBox="`0 0 ${width} ${height}`" preserveAspectRatio="none" style="display: block;">
    <defs>
      <linearGradient :id="gid" x1="0" y1="0" x2="0" y2="1">
        <stop offset="0%" :stop-color="color" stop-opacity="0.32" />
        <stop offset="100%" :stop-color="color" stop-opacity="0" />
      </linearGradient>
    </defs>
    <line
      v-for="f in [0.25, 0.5, 0.75]"
      :key="f"
      :x1="pad" :x2="width - pad"
      :y1="10 + f * (height - 26)" :y2="10 + f * (height - 26)"
      stroke="var(--line)" stroke-width="1"
    />
    <path :d="area" :fill="`url(#${gid})`" />
    <path v-if="line2" :d="line2" fill="none" :stroke="color2" stroke-width="2" stroke-opacity="0.85" stroke-dasharray="3 3" />
    <path :d="line" fill="none" :stroke="color" stroke-width="2.4" stroke-linecap="round" />
    <circle v-if="pts.length" :cx="pts[pts.length - 1][0]" :cy="pts[pts.length - 1][1]" r="3.5" :fill="color" stroke="var(--surface)" stroke-width="2" />
  </svg>
</template>

<script lang="ts" setup>
import { computed, useId } from 'vue'
import { smoothPath } from './path'

const props = withDefaults(defineProps<{
  data: number[]
  data2?: number[]
  color?: string
  color2?: string
  height?: number
  width?: number
  pad?: number
}>(), {
  color: 'var(--brand)',
  color2: 'var(--cyan)',
  height: 160,
  width: 600,
  pad: 8,
})

const gid = `g${useId().replace(/[:-]/g, '')}`

const max = computed(() => {
  const all = props.data2 ? [...props.data, ...props.data2] : props.data
  return Math.max(...all) * 1.12 || 1
})
const xStep = computed(() => (props.width - props.pad * 2) / Math.max(1, props.data.length - 1))
const toPts = (arr: number[]): [number, number][] =>
  arr.map((v, i) => [props.pad + i * xStep.value, props.height - 10 - (v / max.value) * (props.height - 26)])

const pts = computed(() => toPts(props.data))
const line = computed(() => smoothPath(pts.value))
const area = computed(() => {
  const p = pts.value
  if (!p.length) return ''
  return `${line.value} L ${p[p.length - 1][0]},${props.height - 6} L ${p[0][0]},${props.height - 6} Z`
})
const line2 = computed(() => (props.data2 ? smoothPath(toPts(props.data2)) : ''))
</script>
