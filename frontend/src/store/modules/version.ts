import { defineStore } from 'pinia'
import HttpUtils from '@/plugins/httputil'
import { isNewerVersion, latestRelease } from '@/plugins/updateCheck'

// 面板版本与自更新状态:侧边栏版本弹层和更新弹窗共用,避免各页面重复拉取
const VersionStore = defineStore('Version', {
  state: () => ({
    appVersion: '',
    latestTag: null as string | null,
    // 是否支持面板内一键更新(Linux 裸机 systemd 或 Docker 容器),由后端判定
    canSelfUpdate: false,
    // Docker 环境:更新写在容器可写层,重建容器会回退,确认弹窗需要提示
    isDocker: false,
    checking: false,
    updateOpen: false,
    loaded: false,
  }),
  getters: {
    newVersion: (s): string | null =>
      s.appVersion && s.latestTag && isNewerVersion(s.latestTag, s.appVersion)
        ? s.latestTag.replace(/^v/i, '')
        : null,
  },
  actions: {
    async init() {
      if (this.loaded) return
      // 只有面板信息拿到才落锁:首次失败(网络抖动/会话刚过期)时不锁死,
      // 侧边栏重新挂载或用户点刷新都还能补回,否则 pill 会一直停在 v—
      await Promise.all([this.loadPanel(), this.loadLatest()])
    },
    // 版本号取 updateInfo 的 current:它和 api/status 的 appVersion 同源
    // (后端两处都是 config.GetVersion()),一个请求就能拿齐版本 + 更新能力,
    // 不必再打首页正在 2s 轮询的 api/status
    async loadPanel() {
      const msg = await HttpUtils.get('api/updateInfo')
      if (!msg.success || !msg.obj) return
      if (msg.obj.current) {
        this.appVersion = msg.obj.current
        this.loaded = true
      }
      this.canSelfUpdate = !!msg.obj.canSelfUpdate
      this.isDocker = !!msg.obj.docker
    },
    async loadLatest(force = false) {
      const tag = await latestRelease(force)
      if (tag) this.latestTag = tag
    },
    // 弹层里的手动刷新:绕过 24h 缓存直查 GitHub,顺带补回首次没拉到的面板信息
    async recheck() {
      if (this.checking) return
      this.checking = true
      try {
        await Promise.all([
          this.loadLatest(true),
          this.appVersion ? Promise.resolve() : this.loadPanel(),
        ])
      } finally {
        this.checking = false
      }
    },
  },
})

export default VersionStore
