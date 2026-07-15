<template>
  <div ref="el" :style="{ width: size + 'px', height: size + 'px' }" />
</template>

<script lang="ts" setup>
import { ref, watch } from 'vue'
import { useChart, resolveColor, cssVar, type EChart } from './useChart'

const props = withDefaults(defineProps<{
  data: { value: number; color: string; name?: string }[]
  size?: number
  thickness?: number
}>(), {
  size: 150,
  thickness: 18,
})

const el = ref<HTMLElement>()

const build = (chart: EChart, host: HTMLElement) => {
  const surface = cssVar(host, '--surface')
  const track = cssVar(host, '--track')
  const text = cssVar(host, '--text')
  const elevated = cssVar(host, '--elevated') || surface
  const lineCol = cssVar(host, '--line')
  const outer = props.size / 2 - 2
  const inner = Math.max(0, outer - props.thickness)
  const seg = props.data.map((d) => ({
    value: d.value,
    name: d.name ?? '',
    itemStyle: { color: resolveColor(host, d.color) },
  }))
  const option: any = {
    tooltip: {
      show: true,
      trigger: 'item',
      backgroundColor: elevated,
      borderColor: lineCol,
      borderWidth: 1,
      textStyle: { color: text, fontSize: 12 },
      formatter: (p: any) => `${p.name ? p.name + ' · ' : ''}${p.percent}%`,
    },
    series: [{
      type: 'pie',
      radius: [inner, outer],
      center: ['50%', '50%'],
      avoidLabelOverlap: false,
      label: { show: false },
      labelLine: { show: false },
      itemStyle: { borderRadius: 4, borderColor: surface, borderWidth: 2 },
      emphasis: { scaleSize: 3 },
      // 数据为空时画一圈 track 占位（与旧 SVG 版的背景环一致）
      data: seg.length ? seg : [{ value: 1, name: '', itemStyle: { color: track } }],
      silent: seg.length === 0,
    }],
  }
  chart.setOption(option, true)
}

const { render } = useChart(el, build)
watch(() => [props.data, props.size, props.thickness], render, { deep: true })
</script>
