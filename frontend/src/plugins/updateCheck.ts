// 面板新版本检查:浏览器直连 GitHub releases API(CORS 开放),面板服务器无需出网。
// 每次调用都直查,不做缓存;网络失败/访问不了 GitHub 时静默返回 null,不打扰使用。
const RELEASE_API = 'https://api.github.com/repos/shenaba/2s-ui/releases/latest'
export const RELEASE_URL = 'https://github.com/shenaba/2s-ui/releases/latest'

export const latestRelease = async (): Promise<string | null> => {
  try {
    const resp = await fetch(RELEASE_API)
    if (!resp.ok) return null
    const tag = (await resp.json())?.tag_name
    if (typeof tag !== 'string' || tag.length === 0) return null
    return tag
  } catch (e) {
    return null
  }
}

// 版本比较:数字段逐段比;同基线时正式版新于预发布(1.5.4 > 1.5.4-alpha.1),
// 预发布之间按数字感知的字典序(alpha.10 > alpha.9)
export const isNewerVersion = (latest: string, current: string): boolean => {
  const parse = (v: string) => {
    const [base, ...pre] = v.replace(/^v/i, '').split('-')
    return { nums: base.split('.').map((n) => parseInt(n) || 0), pre: pre.join('-') }
  }
  const a = parse(latest)
  const b = parse(current)
  for (let i = 0; i < Math.max(a.nums.length, b.nums.length); i++) {
    const x = a.nums[i] ?? 0
    const y = b.nums[i] ?? 0
    if (x !== y) return x > y
  }
  if (a.pre === b.pre) return false
  if (!a.pre) return true
  if (!b.pre) return false
  return a.pre.localeCompare(b.pre, undefined, { numeric: true }) > 0
}
