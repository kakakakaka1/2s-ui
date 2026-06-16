<template>
  <div ref="el" :style="{ width: width + 'px', height: height + 'px' }" />
</template>

<script lang="ts" setup>
import { ref, watch } from 'vue'
import { useChart, resolveColor, type EChart } from './useChart'

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

const el = ref<HTMLElement>()

const build = (chart: EChart, host: HTMLElement) => {
  const col = resolveColor(host, props.color)
  const min = Math.min(...props.data, 0)
  const max = Math.max(...props.data, 1)
  const option: any = {
    grid: { left: 2, right: 2, top: 3, bottom: 3 },
    xAxis: { type: 'category', show: false, boundaryGap: false, data: props.data.map((_, i) => i) },
    yAxis: { type: 'value', show: false, min, max },
    series: [{
      type: 'line',
      data: props.data,
      showSymbol: false,
      smooth: 0.4,
      lineStyle: { color: col, width: 2, cap: 'round' },
      silent: true,
    }],
    tooltip: { show: false },
  }
  chart.setOption(option)
}

const { render } = useChart(el, build)
watch(() => [props.data, props.color], render, { deep: true })
</script>
