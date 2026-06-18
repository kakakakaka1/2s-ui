<template>
  <div ref="el" :style="{ width: '100%', minWidth: 0, height: height + 'px' }" />
</template>

<script lang="ts" setup>
import { ref, watch } from 'vue'
import echarts from '@/plugins/echarts'
import { useChart, resolveColor, withAlpha, cssVar } from './useChart'

const props = withDefaults(defineProps<{
  data: number[]
  data2?: number[]
  color?: string
  color2?: string
  height?: number
  width?: number
  pad?: number
  /** tooltip 数值格式化（如流量传 HumanReadable.sizeFormat）；不传则千分位展示 */
  valueFormatter?: (v: number) => string
  /** x 轴每个点的标签（日期/时间），用于 tooltip 标题行 */
  labels?: string[]
}>(), {
  color: 'var(--brand)',
  color2: 'var(--emerald)',
  height: 160,
  width: 600,
  pad: 8,
})

const el = ref<HTMLElement>()

const build = (chart: ReturnType<typeof echarts.init>, host: HTMLElement) => {
  const c1 = resolveColor(host, props.color)
  const c2 = resolveColor(host, props.color2)
  const lineCol = cssVar(host, '--line')
  const muted = cssVar(host, '--text-3')
  const surface = cssVar(host, '--surface')
  const elevated = cssVar(host, '--elevated') || surface
  const text = cssVar(host, '--text')
  const fmt = props.valueFormatter ?? ((v: number) => Math.round(v).toLocaleString())

  const all = props.data2 ? [...props.data, ...props.data2] : props.data
  const max = (Math.max(...all, 0) || 1) * 1.12
  const last = props.data.length - 1

  const series: any[] = [
    {
      type: 'line',
      smooth: 0.4,
      showSymbol: false,
      clip: false, // 末端圆点贴边时不被裁剪
      lineStyle: { color: c1, width: 2.4, cap: 'round' },
      areaStyle: {
        color: new echarts.graphic.LinearGradient(0, 0, 0, 1, [
          { offset: 0, color: withAlpha(c1, 0.32) },
          { offset: 1, color: withAlpha(c1, 0) },
        ]),
      },
      data: props.data,
      z: 3,
      markPoint: {
        symbol: 'circle',
        symbolSize: 7,
        silent: true,
        label: { show: false },
        itemStyle: { color: c1, borderColor: surface, borderWidth: 2 },
        data: last >= 0 ? [{ coord: [last, props.data[last]] }] : [],
      },
    },
  ]

  if (props.data2) {
    series.push({
      type: 'line',
      smooth: 0.4,
      showSymbol: false,
      clip: false,
      lineStyle: { color: c2, width: 2, type: 'dashed', opacity: 0.85 },
      data: props.data2,
      z: 2,
    })
  }

  const option: any = {
    animationDurationUpdate: 300,
    grid: { left: props.pad, right: props.pad + 4, top: 14, bottom: 10, containLabel: false },
    tooltip: {
      trigger: 'axis',
      backgroundColor: elevated,
      borderColor: lineCol,
      borderWidth: 1,
      padding: [6, 10],
      textStyle: { color: text, fontSize: 12 },
      axisPointer: { type: 'line', lineStyle: { color: muted, width: 1, type: 'dashed' } },
      formatter: (ps: any[]) => {
        const idx = ps[0]?.dataIndex ?? 0
        const lbl = props.labels?.[idx]
        const head = lbl ? `<div style="font-size:11px;color:${muted};margin-bottom:5px;">${lbl}</div>` : ''
        return head + ps
          .map(
            (p) =>
              `<div style="display:flex;align-items:center;gap:7px;line-height:1.65;"><span style="width:8px;height:8px;border-radius:2px;background:${p.color};flex:none;"></span><span class="mono">${fmt(Number(p.value))}</span></div>`,
          )
          .join('')
      },
    },
    xAxis: {
      type: 'category',
      boundaryGap: false,
      show: false,
      data: props.labels ?? props.data.map((_, i) => i),
    },
    yAxis: {
      // show:true 但关掉轴线/刻度/标签 —— 否则 show:false 会连带隐藏 splitLine 网格线
      type: 'value',
      min: 0,
      max,
      show: true,
      axisLine: { show: false },
      axisTick: { show: false },
      axisLabel: { show: false },
      splitNumber: 3,
      splitLine: { show: true, lineStyle: { color: lineCol, width: 1 } },
    },
    series,
  }

  chart.setOption(option) // merge 模式：数据变化平滑过渡
}

const { render } = useChart(el, build)

// 实时数据变化时重渲染
watch(() => [props.data, props.data2], render, { deep: true })
</script>
