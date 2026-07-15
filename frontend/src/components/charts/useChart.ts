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

/**
 * 仪表盘卡片显隐(v-if)会改变同一行图表所在 grid 列的宽度，但 ECharts 不会自动跟随父级
 * grid 重排；且多个图表同帧改宽时 ResizeObserver 可能因 "ResizeObserver loop" 丢通知，
 * 导致个别图表没跟着容器缩放而溢出卡片。由布局发生变化的一方(如 Home 的卡片开关)显式调用
 * 此函数，广播一次重排事件，通知所有在用图表重新测量尺寸。
 */
const RELAYOUT_EVENT = 'sui:chart-relayout'
export function requestChartRelayout() {
  window.dispatchEvent(new Event(RELAYOUT_EVENT))
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
  let rafId: number | undefined

  const render = () => {
    const el = elRef.value
    if (el && chart.value) build(chart.value, el)
  }

  // 把一连串尺寸变化合并到下一帧再 resize：切换卡片时同一行多个图表会同帧改宽，裸调 resize 易
  // 触发 ResizeObserver loop 丢通知，导致个别图表没跟随容器缩放而溢出卡片。用 rAF 等浏览器把
  // 布局算定后再量；并【显式把容器宽高传给 resize】——SVG 渲染器自测尺寸时可能读到被自身 SVG
  // 撑大的旧宽度而缩不回去（#15 表现为关掉再开回卡片后图表溢出边框），显式传 clientWidth/Height 绕开。
  const scheduleResize = () => {
    if (rafId !== undefined) cancelAnimationFrame(rafId)
    rafId = requestAnimationFrame(() => {
      rafId = undefined
      const el = elRef.value
      if (!el || !chart.value) return
      const w = el.clientWidth
      const h = el.clientHeight
      if (w > 0 && h > 0) chart.value.resize({ width: w, height: h })
    })
  }

  onMounted(() => {
    const el = elRef.value
    if (!el) return
    chart.value = echarts.init(el, undefined, { renderer: 'svg' })
    render()
    // 容器尺寸变化 -> resize（替代原 SVG viewBox 的自动响应）
    ro = new ResizeObserver(scheduleResize)
    ro.observe(el)
    // 卡片显隐引起的 grid 列宽变化，RO 不一定可靠（多图同帧重排会丢通知），由布局方显式广播兜底
    window.addEventListener(RELAYOUT_EVENT, scheduleResize)
    // 明暗主题切换（<html data-theme> 变化）-> 重取 CSS 变量并重渲染
    mo = new MutationObserver(render)
    mo.observe(document.documentElement, { attributes: true, attributeFilter: ['data-theme'] })
  })

  onBeforeUnmount(() => {
    if (rafId !== undefined) cancelAnimationFrame(rafId)
    window.removeEventListener(RELAYOUT_EVENT, scheduleResize)
    ro?.disconnect()
    mo?.disconnect()
    chart.value?.dispose()
  })

  return { render }
}
