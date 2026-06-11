<template>
  <LogsModal :visible="logsOpen" @close="logsOpen = false" />
  <BackupModal :visible="backupOpen" @close="backupOpen = false" />
  <UsageStatsModal :visible="usageOpen" @close="usageOpen = false" />

  <div class="page-stack-lg fade-up">
    <!-- toolbar -->
    <div class="toolbar">
      <Btn sm @click="backupOpen = true"><Ico name="copy" :size="15" /> {{ $t('main.backup.title') }}</Btn>
      <Btn sm @click="logsOpen = true"><Ico name="list" :size="15" /> {{ $t('basic.log.title') }}</Btn>
      <Btn sm @click="usageOpen = true"><Ico name="chart" :size="15" /> {{ $t('main.stats.title') }}</Btn>
      <div class="grow" />
      <TilesMenu :show="tiles" @toggle="toggleTile" @reset="resetTiles" />
    </div>

    <!-- stat cards -->
    <div v-if="visStats.length > 0" class="stat-grid" :style="{ '--cols': visStats.length }">
      <StatCard
        v-if="tiles['stat-clients']"
        icon="clients" :label="$t('ui.onlineClients')"
        :value="String(onlineUsers)" :sub="$t('ui.ofTotal', { n: data.clients.length })"
        accent="var(--emerald)" :trend="$t('ui.live')" :spark="buf.users"
      />
      <StatCard
        v-if="tiles['stat-download']"
        icon="download" :label="$t('stats.download')"
        :value="HumanReadable.sizeFormat(totalDown)" :sub="$t('ui.inboundTraffic')"
        accent="var(--cyan)" :spark="buf.netIn"
      />
      <StatCard
        v-if="tiles['stat-upload']"
        icon="upload" :label="$t('stats.upload')"
        :value="HumanReadable.sizeFormat(totalUp)" :sub="$t('ui.outboundTraffic')"
        accent="var(--violet)" :spark="buf.netOut"
      />
      <StatCard
        v-if="tiles['stat-inbounds']"
        icon="inbound" :label="$t('ui.activeInbounds')"
        :value="`${onlineInbounds}/${data.inbounds.length}`" :sub="$t('ui.listenersHealthy')"
        accent="var(--brand)" :spark="buf.inbounds"
      />
    </div>

    <!-- traffic + protocol -->
    <div v-if="tiles.traffic || tiles.protocol" class="main-grid" :class="{ both: tiles.traffic && tiles.protocol }">
      <DPanel v-if="tiles.traffic" :title="$t('ui.traffic')" :sub="$t('ui.trafficSub')">
        <template #right>
          <Chip color="emerald" dot>{{ $t('ui.live') }}</Chip>
        </template>
        <div style="display: flex; gap: 22px; margin-bottom: 12px; flex-wrap: wrap;">
          <Legend color="var(--brand)" :label="$t('stats.download')" :value="HumanReadable.sizeFormat(netInNow) + '/s'" />
          <Legend color="var(--cyan)" :label="$t('stats.upload')" :value="HumanReadable.sizeFormat(netOutNow) + '/s'" dashed />
        </div>
        <AreaChart :data="buf.netIn" :data2="buf.netOut" :height="180" />
        <div class="mono" style="display: flex; justify-content: space-between; margin-top: 8px; font-size: 10.5px; color: var(--text-3);">
          <span v-for="l in trafficAxis" :key="l">{{ l }}</span>
        </div>
      </DPanel>

      <DPanel v-if="tiles.protocol" :title="$t('ui.protocolMix')" :sub="$t('ui.protocolSub')">
        <div v-if="protoMix.length === 0" style="font-size: 12.5px; color: var(--text-3);">{{ $t('noData') }}</div>
        <div v-else style="display: flex; align-items: center; gap: 18px; flex-wrap: wrap;">
          <div style="position: relative; flex: none;">
            <Donut :data="protoMix" :size="146" :thickness="17" />
            <div style="position: absolute; inset: 0; display: flex; flex-direction: column; align-items: center; justify-content: center;">
              <div class="mono" style="font-size: 22px; font-weight: 700; line-height: 1;">{{ onlineUsers }}</div>
              <div style="font-size: 10.5px; color: var(--text-3);">{{ $t('ui.sessions') }}</div>
            </div>
          </div>
          <div style="flex: 1; display: flex; flex-direction: column; gap: 9px; min-width: 120px;">
            <div v-for="p in protoMix" :key="p.name" style="display: flex; align-items: center; gap: 8px; font-size: 12.5px;">
              <span :style="{ width: '8px', height: '8px', borderRadius: '3px', background: p.color, flex: 'none' }" />
              <span style="flex: 1; color: var(--text-2); font-weight: 600;">{{ p.name }}</span>
              <span class="mono" style="color: var(--text-3);">{{ p.pct }}%</span>
            </div>
          </div>
        </div>
      </DPanel>
    </div>

    <!-- monitors -->
    <div v-if="monCount > 0" class="mon-grid" :style="{ '--mon-cols': monCols }">
      <DPanel v-if="tiles.resources" :title="$t('ui.systemResources')" :sub="$t('ui.realtimeUtil')">
        <div style="display: grid; grid-template-columns: 1fr 1fr; gap: 6px; justify-items: center;">
          <Gauge :value="gCpu.pct" color="var(--brand)" :label="$t('ui.cpu')" :sub="sys.cpuCount ? sys.cpuCount + ' cores' : ''" :size="120" />
          <Gauge :value="gMem.pct" color="var(--cyan)" :label="$t('ui.memory')" :sub="gMem.sub" :size="120" />
          <Gauge :value="gDsk.pct" color="var(--violet)" :label="$t('ui.disk')" :sub="gDsk.sub" :size="120" />
          <Gauge :value="gSwp.pct" color="var(--emerald)" :label="$t('ui.swap')" :sub="gSwp.sub" :size="120" />
        </div>
      </DPanel>

      <DPanel v-if="tiles.network" :title="$t('ui.networkThroughput')" :sub="`↓ ${HumanReadable.sizeFormat(netInNow)}/s · ↑ ${HumanReadable.sizeFormat(netOutNow)}/s`">
        <div style="display: flex; flex-direction: column; gap: 4px;">
          <div style="display: flex; align-items: center; gap: 8px; font-size: 12px; color: var(--text-3);">
            <span style="width: 8px; height: 8px; border-radius: 3px; background: var(--cyan);" />{{ $t('ui.inboundLbl') }}
          </div>
          <AreaChart :data="buf.netIn" color="var(--cyan)" :height="88" />
          <div style="display: flex; align-items: center; gap: 8px; font-size: 12px; color: var(--text-3); margin-top: 4px;">
            <span style="width: 8px; height: 8px; border-radius: 3px; background: var(--violet);" />{{ $t('ui.outboundLbl') }}
          </div>
          <AreaChart :data="buf.netOut" color="var(--violet)" :height="88" />
        </div>
      </DPanel>

      <DPanel v-if="tiles.server" :title="$t('ui.server')" :sub="sys.hostName ?? ''" :pad="0">
        <div style="display: flex; flex-direction: column;">
          <InfoRow :k="$t('ui.host')" :v="sys.hostName ?? '—'" />
          <InfoRow :k="$t('ui.cpu')" :v="(sys.cpuCount ?? '—') + ' cores'" chip :title="sys.cpuType" />
          <InfoRow k="IPv4" :v="sys.ipv4?.[0] ?? '—'" chip mono />
          <InfoRow :k="$t('main.info.memory')" :v="HumanReadable.sizeFormat(sbd.stats?.Alloc ?? 0)" mono />
          <InfoRow :k="$t('ui.panelLbl')" :v="'v' + (sys.appVersion ?? '—')" chip accent />
          <InfoRow :k="$t('ui.uptime')" :v="uptime" mono />
          <div style="padding: 13px 20px; display: flex; align-items: center; gap: 10px; border-top: 1px solid var(--line);">
            <Chip v-if="sbd.running" color="emerald" dot>{{ $t('ui.singboxRunning') }}</Chip>
            <Chip v-else color="rose" dot>sing-box · {{ $t('main.info.running') }}: {{ $t('no') }}</Chip>
            <Btn variant="subtle" sm style="margin-inline-start: auto;" :loading="restarting" @click="restartSb">
              <Ico name="refresh" :size="14" /> {{ $t('ui.restart') }}
            </Btn>
          </div>
        </div>
      </DPanel>
    </div>

    <!-- activity -->
    <DPanel v-if="tiles.activity" :title="$t('ui.activity')" :sub="$t('ui.activitySub')" :pad="0">
      <template #right>
        <Btn variant="subtle" sm @click="logsOpen = true">{{ $t('ui.viewLog') }} <Ico name="chevron" :size="14" /></Btn>
      </template>
      <div>
        <div v-if="changes.length === 0" style="padding: 16px 20px; font-size: 12.5px; color: var(--text-3);">{{ $t('noData') }}</div>
        <div
          v-for="(e, i) in changes"
          :key="e.id"
          :style="{ display: 'flex', alignItems: 'center', gap: '13px', padding: '12px 20px', borderTop: i ? '1px solid var(--line)' : 'none', flexWrap: 'wrap' }"
        >
          <span :style="{ width: '8px', height: '8px', borderRadius: '50%', background: actColor(e.action), flex: 'none', boxShadow: `0 0 0 4px color-mix(in srgb, ${actColor(e.action)} 18%, transparent)` }" />
          <span style="font-weight: 700; font-size: 13px;">{{ e.actor }}</span>
          <span style="font-size: 13px; color: var(--text-2); min-width: 0;">{{ actionLabel(e.action) }}</span>
          <Chip v-if="e.key" style="height: 22px;">{{ objectLabel(e.key) }}</Chip>
          <span class="mono" dir="ltr" style="margin-inline-start: auto; font-size: 11.5px; color: var(--text-3);">{{ e.time }}</span>
        </div>
      </div>
    </DPanel>

    <EmptyState
      v-if="allHidden"
      icon="settings"
      :title="$t('ui.tilesTitle')"
      :desc="$t('ui.customize')"
    />
  </div>
</template>

<script lang="ts" setup>
import { computed, onBeforeUnmount, onMounted, reactive, ref } from 'vue'
import { useI18n } from 'vue-i18n'
import HttpUtils from '@/plugins/httputil'
import { HumanReadable } from '@/plugins/utils'
import Data from '@/store/modules/data'
import Btn from '@/components/ui/Btn.vue'
import Ico from '@/components/ui/Ico.vue'
import Chip from '@/components/ui/Chip.vue'
import DPanel from '@/components/ui/DPanel.vue'
import EmptyState from '@/components/ui/EmptyState.vue'
import AreaChart from '@/components/charts/AreaChart.vue'
import Donut from '@/components/charts/Donut.vue'
import Gauge from '@/components/charts/Gauge.vue'
import StatCard from './dashboard/StatCard.vue'
import TilesMenu from './dashboard/TilesMenu.vue'
import InfoRow from './dashboard/InfoRow.vue'
import Legend from './dashboard/Legend.vue'
import LogsModal from '@/layouts/drawers/LogsModal.vue'
import BackupModal from '@/layouts/drawers/BackupModal.vue'
import UsageStatsModal from '@/layouts/drawers/UsageStatsModal.vue'
import { loadTiles, saveTiles, defaultTiles } from './dashboard/tiles'

const { t, te, locale } = useI18n({ useScope: 'global' })
const data = Data()

/* ---------- tiles visibility ---------- */
const tiles = reactive(loadTiles())
const toggleTile = (id: string) => { tiles[id] = !tiles[id]; saveTiles(tiles) }
const resetTiles = () => { Object.assign(tiles, defaultTiles()); saveTiles(tiles) }
const visStats = computed(() => ['stat-clients', 'stat-download', 'stat-upload', 'stat-inbounds'].filter((i) => tiles[i]))
const monCount = computed(() => ['resources', 'network', 'server'].filter((i) => tiles[i]).length)
const monCols = computed(() => (monCount.value === 3 ? '1.1fr 1.1fr 1fr' : Array(monCount.value).fill('1fr').join(' ')))
const allHidden = computed(() => Object.values(tiles).every((v) => !v))

/* ---------- modals ---------- */
const logsOpen = ref(false)
const backupOpen = ref(false)
const usageOpen = ref(false)

/* ---------- live status polling (2s, like the old dashboard) ---------- */
const status = ref<any>({})
const sys = ref<any>({})
const sbd = computed(() => status.value.sbd ?? {})
const BUF = 40
const buf = reactive({
  netIn: [] as number[],
  netOut: [] as number[],
  users: [] as number[],
  inbounds: [] as number[],
})
let lastNet: { recv: number; sent: number } | null = null

const push = (arr: number[], v: number) => { arr.push(v); if (arr.length > BUF) arr.shift() }

const poll = async () => {
  const r = ['net', 'sbd']
  if (tiles.resources) r.push('cpu', 'mem', 'dsk', 'swp')
  const msg = await HttpUtils.get('api/status', { r: r.join(',') })
  if (!msg.success || !msg.obj) return
  status.value = { ...status.value, ...msg.obj }
  const net = msg.obj.net
  if (net) {
    if (lastNet) {
      push(buf.netIn, Math.max(0, (net.recv - lastNet.recv) / 2))
      push(buf.netOut, Math.max(0, (net.sent - lastNet.sent) / 2))
    }
    lastNet = { recv: net.recv, sent: net.sent }
  }
  push(buf.users, data.onlines.user?.length ?? 0)
  push(buf.inbounds, data.onlines.inbound?.length ?? 0)
}

const loadSys = async () => {
  const msg = await HttpUtils.get('api/status', { r: 'sys' })
  if (msg.success && msg.obj?.sys) sys.value = msg.obj.sys
}

let pollId: ReturnType<typeof setInterval> | null = null
onMounted(() => {
  loadSys()
  poll()
  loadChanges()
  pollId = setInterval(poll, 2000)
})
onBeforeUnmount(() => { if (pollId) clearInterval(pollId) })

/* ---------- derived live values ---------- */
const netInNow = computed(() => buf.netIn[buf.netIn.length - 1] ?? 0)
const netOutNow = computed(() => buf.netOut[buf.netOut.length - 1] ?? 0)
const onlineUsers = computed(() => data.onlines.user?.length ?? 0)
const onlineInbounds = computed(() => data.onlines.inbound?.length ?? 0)
const totalDown = computed(() => data.clients.reduce((a: number, c: any) => a + (c.down ?? 0), 0))
const totalUp = computed(() => data.clients.reduce((a: number, c: any) => a + (c.up ?? 0), 0))
const trafficAxis = computed(() => {
  const span = (buf.netIn.length - 1) * 2
  if (span <= 0) return []
  return [`-${span}s`, `-${Math.round(span * 0.75)}s`, `-${Math.round(span * 0.5)}s`, `-${Math.round(span * 0.25)}s`, 'now']
})

const pctOf = (d: any) => (d && d.total ? Math.min(100, Math.ceil((d.current * 100) / d.total)) : 0)
const subOf = (d: any) => (d && d.total ? `${HumanReadable.sizeFormat(d.current, 0)} / ${HumanReadable.sizeFormat(d.total, 0)}` : '')
const gCpu = computed(() => ({ pct: Math.ceil(status.value.cpu ?? 0) }))
const gMem = computed(() => ({ pct: pctOf(status.value.mem), sub: subOf(status.value.mem) }))
const gDsk = computed(() => ({ pct: pctOf(status.value.dsk), sub: subOf(status.value.dsk) }))
const gSwp = computed(() => ({ pct: pctOf(status.value.swp), sub: subOf(status.value.swp) }))

const uptime = computed(() =>
  sys.value.bootTime ? HumanReadable.formatSecond(Date.now() / 1000 - sys.value.bootTime) : '—'
)

/* ---------- protocol mix (clients per inbound protocol) ---------- */
const PALETTE = ['var(--brand)', 'var(--cyan)', 'var(--violet)', 'var(--emerald)', 'var(--amber)', 'var(--text-3)']
const protoMix = computed(() => {
  const typeById = new Map<number, string>(data.inbounds.map((i: any) => [i.id, i.type]))
  const byType: Record<string, number> = {}
  for (const c of data.clients) {
    if (!c.enable) continue
    const types = new Set((c.inbounds ?? []).map((id: number) => typeById.get(id)).filter(Boolean))
    for (const tp of types) byType[tp as string] = (byType[tp as string] ?? 0) + 1
  }
  const entries = Object.entries(byType).sort((a, b) => b[1] - a[1])
  const top = entries.slice(0, 5)
  const rest = entries.slice(5).reduce((a, e) => a + e[1], 0)
  if (rest > 0) top.push([t('ui.none'), rest] as any)
  const total = top.reduce((a, e) => a + e[1], 0) || 1
  return top.map(([name, value], i) => ({
    name,
    value,
    pct: Math.round((value / total) * 100),
    color: PALETTE[i % PALETTE.length],
  }))
})

/* ---------- activity (changes log) ---------- */
const changes = ref<any[]>([])
const loadChanges = async () => {
  const msg = await HttpUtils.get('api/changes', { a: '', k: '', c: 6 })
  if (msg.success && Array.isArray(msg.obj)) {
    const l = String(locale.value) == 'fa' ? 'fa-IR' : 'en-US'
    changes.value = msg.obj.map((c: any) => ({
      id: c.id,
      actor: c.Actor,
      action: c.action,
      key: c.key,
      time: new Date(Number(c.dateTime) * 1000).toLocaleTimeString(l, { hour: '2-digit', minute: '2-digit' }),
    }))
  }
}

const ACT_COLORS: Record<string, string> = {
  add: 'var(--emerald)', new: 'var(--emerald)',
  del: 'var(--rose)', delete: 'var(--rose)',
  edit: 'var(--brand)', set: 'var(--brand)', update: 'var(--brand)',
  restart: 'var(--amber)', reset: 'var(--amber)',
}
const actColor = (a: string) => ACT_COLORS[a] ?? 'var(--brand)'
const actionLabel = (a: string) => (te('actions.' + a) ? t('actions.' + a) : a)
const objectLabel = (k: string) => (te('objects.' + k) ? t('objects.' + k) : k)

/* ---------- restart sing-box ---------- */
const restarting = ref(false)
const restartSb = async () => {
  restarting.value = true
  await HttpUtils.post('api/restartSb', {})
  restarting.value = false
}
</script>

<style scoped>
.stat-grid {
  display: grid;
  grid-template-columns: repeat(var(--cols, 4), 1fr);
  gap: 14px;
}
.main-grid { display: grid; grid-template-columns: 1fr; gap: 18px; }
.main-grid.both { grid-template-columns: 1.55fr 1fr; }
.mon-grid { display: grid; grid-template-columns: var(--mon-cols, 1fr); gap: 18px; }
@media (max-width: 1180px) {
  .main-grid.both { grid-template-columns: 1fr; }
  .mon-grid { grid-template-columns: 1fr 1fr !important; }
}
@media (max-width: 820px) {
  .stat-grid { grid-template-columns: 1fr 1fr; gap: 10px; }
  .mon-grid { grid-template-columns: 1fr !important; }
}
</style>
