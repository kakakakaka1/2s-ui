/**
 * 图表通用逻辑：初始化 echarts 实例、容器响应式 resize、主题切换重渲染、销毁。
 * 所有基于 echarts 的图表组件复用此 composable，避免每个组件重复踩主题/resize 的坑。
 */
import { onBeforeUnmount, onMounted, shallowRef, type Ref } from 'vue'
import echarts from '@/plugins/echarts'

/** 取元素继承到的 CSS 变量计算值（会随 <html data-theme> 主题变化） */
export function cssVar(el: HTMLElement, name: string): string {
  return getComputedStyle(el).getPropertyValue(name).trim()
}

/** 解析颜色：'var(--brand)' -> 实际色值；已是具体色值则原样返回（echarts 不认 CSS 变量） */
export function resolveColor(el: HTMLElement, c: string): string {
  const m = c.match(/var\((--[\w-]+)\)/)
  return m ? cssVar(el, m[1]) : c
}

/** 给颜色叠加透明度，支持 #rgb / #rrggbb / rgb() / rgba() */
export function withAlpha(color: string, alpha: number): string {
  if (color.startsWith('#')) {
    const h = color.slice(1)
    const v = h.length === 3 ? h.replace(/./g, '$&$&') : h
    const n = parseInt(v, 16)
    return `rgba(${(n >> 16) & 255}, ${(n >> 8) & 255}, ${n & 255}, ${alpha})`
  }
  const nums = color.match(/[\d.]+/g)
  if (nums && nums.length >= 3) return `rgba(${nums[0]}, ${nums[1]}, ${nums[2]}, ${alpha})`
  return color
}

export type EChart = ReturnType<typeof echarts.init>

/**
 * @param elRef  图表容器 ref
 * @param build  构建/更新 option 的回调，参数为 (chart, host)，内部用 cssVar/resolveColor 取色
 * @returns render —— 手动触发重渲染（如数据变化时调用）
 */
export function useChart(
  elRef: Ref<HTMLElement | undefined>,
  build: (chart: EChart, host: HTMLElement) => void,
) {
  const chart = shallowRef<EChart>()
  let ro: ResizeObserver | undefined
  let mo: MutationObserver | undefined

  const render = () => {
    const el = elRef.value
    if (el && chart.value) build(chart.value, el)
  }

  onMounted(() => {
    const el = elRef.value
    if (!el) return
    chart.value = echarts.init(el, undefined, { renderer: 'svg' })
    render()
    // 容器尺寸变化 -> resize（替代原 SVG viewBox 的自动响应）
    ro = new ResizeObserver(() => chart.value?.resize())
    ro.observe(el)
    // 明暗主题切换（<html data-theme> 变化）-> 重取 CSS 变量并重渲染
    mo = new MutationObserver(render)
    mo.observe(document.documentElement, { attributes: true, attributeFilter: ['data-theme'] })
  })

  onBeforeUnmount(() => {
    ro?.disconnect()
    mo?.disconnect()
    chart.value?.dispose()
  })

  return { render }
}
