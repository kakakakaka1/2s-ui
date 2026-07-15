import axios from 'axios'

axios.defaults.headers.post['Content-Type'] = 'application/x-www-form-urlencoded; charset=UTF-8'
axios.defaults.headers.common['X-Requested-With'] = 'XMLHttpRequest'

axios.defaults.baseURL = "./"

// 注意：defaults 必须在 create() 之前设好——实例是在创建时把 defaults 合进去的
const api = axios.create()

// Dev-only：本地无后端时，用假数据喂 api/load、api/status、api/changes，让首页可演示。
// 生产构建中 import.meta.env.DEV 为 false，整块连同 ./mock 一起被 tree-shake 掉。
if ((import.meta as any).env?.DEV) {
  // 先打真实后端；仅当请求失败（如本地没起后端）时用 mock 兜底，避免干扰有后端的正常开发
  const builtin = axios.getAdapter(['xhr', 'http', 'fetch'])
  api.defaults.adapter = async (config) => {
    try {
      return await builtin(config)
    } catch (err) {
      if (axios.isCancel(err)) throw err
      // 只在"后端不在"时用 mock 兜底,真实故障照常抛出,避免被假数据掩盖。
      // 注意 vite 代理会把 ECONNREFUSED 转成 502 响应(而非网络层失败),
      // 所以 dev 下 502 网关错误也按"后端没起"处理,否则 mock 永远不触发
      if (axios.isAxiosError(err) && err.response && err.response.status !== 502) throw err
      const { getMock } = await import('./mock')
      const m = getMock(config.url || '', config.params)
      if (m) return { data: m, status: 200, statusText: 'OK (mock)', headers: {}, config, request: {} } as any
      throw err
    }
  }
}

export default api
