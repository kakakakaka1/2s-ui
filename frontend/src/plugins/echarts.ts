/**
 * echarts —— 按需引入，仅注册用到的模块以控制打包体积。
 * 新增图表类型/组件时在此 use()，不要直接 import 整个 'echarts'。
 */
import * as echarts from 'echarts/core'
import { LineChart, PieChart, GaugeChart } from 'echarts/charts'
import { GridComponent, TooltipComponent, MarkPointComponent } from 'echarts/components'
import { SVGRenderer } from 'echarts/renderers'

// 注意：markPoint/markLine/markArea 是独立组件，按需引入时必须显式注册，否则配置被静默忽略
echarts.use([LineChart, PieChart, GaugeChart, GridComponent, TooltipComponent, MarkPointComponent, SVGRenderer])

// 开发期暴露到 window，方便控制台调试图表实例（生产构建不会执行此分支）
if ((import.meta as any).env?.DEV) {
  ;(window as any).echarts = echarts
}

export default echarts
