import axios from 'axios'

axios.defaults.headers.post['Content-Type'] = 'application/x-www-form-urlencoded; charset=UTF-8'
axios.defaults.headers.common['X-Requested-With'] = 'XMLHttpRequest'

axios.defaults.baseURL = "./"
const pendingRequests = new Map()

axios.interceptors.request.use(
    (config) => {
        // Generate a unique key for the request
        const requestKey = `${config.method}:${config.url}`
        
        // Check if there is already a pending request with the same key
        if (pendingRequests.has(requestKey)) {
            const cancelSource = pendingRequests.get(requestKey)
            cancelSource.cancel('Duplicate request cancelled')
        }
        
        // Create a new cancel token for the request
        const cancelSource = axios.CancelToken.source()
        config.cancelToken = cancelSource.token
        
        // Store the cancel token in the pending requests map
        pendingRequests.set(requestKey, cancelSource)
        
        if (config.data instanceof FormData) {
            config.headers['Content-Type'] = 'multipart/form-data'
        }
        return config
    },
    (error) => Promise.reject(error),
)

axios.interceptors.response.use(
    (response) => {
        // Remove the request from the pending requests map
        const requestKey = `${response.config.method}:${response.config.url}`
        pendingRequests.delete(requestKey)
        return response
    },
    (error) => {
        if (axios.isCancel(error)) {
            // Handle duplicate request cancellation here if needed
            console.warn(error.message)
        } else {
            // Remove the request from the pending requests map on error
            const requestKey = `${error.config.method}:${error.config.url}`
            pendingRequests.delete(requestKey)
        }
        return Promise.reject(error)
    }
)

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
      // 只在网络层失败（拿不到响应，如后端没起）时用 mock 兜底；
      // HTTP 4xx/5xx 是有响应的真实故障，照常抛出，避免被假数据掩盖
      if (axios.isAxiosError(err) && err.response) throw err
      const { getMock } = await import('./mock')
      const m = getMock(config.url || '', config.params)
      if (m) return { data: m, status: 200, statusText: 'OK (mock)', headers: {}, config, request: {} } as any
      throw err
    }
  }
}

export default api
