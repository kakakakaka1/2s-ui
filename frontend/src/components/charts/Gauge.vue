<template>
  <div style="display: flex; flex-direction: column; align-items: center;">
    <div ref="el" :style="{ width: size + 'px', height: size * 0.82 + 'px' }" />
    <div style="text-align: center; margin-top: -2px;">
      <div style="font-size: 13px; font-weight: 700;">{{ label }}</div>
      <div class="mono" style="font-size: 11px; color: var(--text-3); margin-top: 2px;">{{ sub }}</div>
    </div>
  </div>
</template>

<script lang="ts" setup>
import { ref, watch } from 'vue'
import { useChart, resolveColor, cssVar, type EChart } from './useChart'

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

const el = ref<HTMLElement>()

const build = (chart: EChart, host: HTMLElement) => {
  const col = resolveColor(host, props.color)
  const track = cssVar(host, '--track')
  const text = cssVar(host, '--text')
  const w = Math.max(7, props.size * 0.07)
  const v = Math.min(100, Math.max(0, props.value || 0))
  const option: any = {
    series: [{
      type: 'gauge',
      startAngle: 220,
      endAngle: -40,
      radius: '96%',
      center: ['50%', '52%'],
      min: 0,
      max: 100,
      progress: { show: true, width: w, roundCap: true, itemStyle: { color: col } },
      axisLine: { roundCap: true, lineStyle: { width: w, color: [[1, track]] } },
      axisTick: { show: false },
      splitLine: { show: false },
      axisLabel: { show: false },
      pointer: { show: false },
      anchor: { show: false },
      title: { show: false },
      detail: {
        valueAnimation: true,
        offsetCenter: [0, '6%'],
        formatter: (val: number) => `${Math.round(val)}%`,
        fontSize: props.size * 0.2,
        fontWeight: 700,
        color: text,
      },
      data: [{ value: v }],
    }],
  }
  chart.setOption(option, true)
}

const { render } = useChart(el, build)
watch(() => [props.value, props.color, props.size], render)
</script>
