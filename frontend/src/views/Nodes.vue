<template>
  <NodeDrawer
    :visible="drawer.visible"
    :id="drawer.id"
    :data="drawer.data"
    @close="drawer.visible = false"
  />

  <!-- delete confirmation -->
  <Modal :open="del.visible" :title="$t('actions.del')" :width="380" @close="del.visible = false">
    <div style="padding: 18px; font-size: 13.5px;">{{ $t('confirm') }}</div>
    <template #footer>
      <Btn @click="del.visible = false">{{ $t('no') }}</Btn>
      <Btn style="color: var(--rose);" :loading="deleting" @click="confirmDelete">
        <Ico name="trash" :size="15" /> {{ $t('yes') }}
      </Btn>
    </template>
  </Modal>

  <div class="page-stack fade-up">
    <div class="toolbar" style="justify-content: center;">
      <Btn variant="primary" sm @click="openDrawer(0)">
        <Ico name="plus" :size="15" /> {{ $t('actions.add') }}
      </Btn>
    </div>

    <EmptyState v-if="nodes.length === 0" icon="server" :title="$t('node.emptyTitle')" :desc="$t('node.emptyDesc')" />

    <div v-else class="entity-grid">
      <EntityCard
        v-for="item in nodes"
        :key="item.id"
        :title="item.name"
        :type="displayUrl(item)"
        :color="stateColor(item)"
        icon="server"
        :rows="cardRows(item)"
      >
        <template #chip>
          <Chip v-if="stateOf(item) === 'online'" color="emerald" dot>{{ $t('node.status.online') }}</Chip>
          <Chip v-else-if="stateOf(item) === 'core-stopped'" color="amber">{{ $t('node.status.coreStopped') }}</Chip>
          <Chip v-else-if="stateOf(item) === 'offline'" color="rose">{{ $t('node.status.offline') }}</Chip>
          <Chip v-else-if="stateOf(item) === 'disabled'">{{ $t('disable') }}</Chip>
          <Chip v-else>{{ $t('node.status.pending') }}</Chip>
        </template>
        <template #actions>
          <CardBtn icon="edit" :title="$t('actions.edit')" @click="openDrawer(item.id)" />
          <CardBtn icon="trash" border danger :title="$t('actions.del')" @click="askDelete(item.id)" />
        </template>
      </EntityCard>
    </div>
  </div>
</template>

<script lang="ts" setup>
import { computed, ref } from 'vue'
import { useI18n } from 'vue-i18n'
import Data from '@/store/modules/data'
import { locale } from '@/locales'
import { Node, NodeStatus } from '@/types/node'
import Btn from '@/components/ui/Btn.vue'
import Ico from '@/components/ui/Ico.vue'
import Chip from '@/components/ui/Chip.vue'
import Modal from '@/components/ui/Modal.vue'
import CardBtn from '@/components/ui/CardBtn.vue'
import EmptyState from '@/components/ui/EmptyState.vue'
import EntityCard, { EntityRow } from '@/components/ui/EntityCard.vue'
import NodeDrawer from '@/layouts/drawers/node/NodeDrawer.vue'

const { t } = useI18n({ useScope: 'global' })
const dataStore = Data()

const nodes = computed((): Node[] => <Node[]>dataStore.nodes)

// ---------------- live state ----------------
type CardState = 'online' | 'offline' | 'core-stopped' | 'disabled' | 'pending'

const statusOf = (n: Node): NodeStatus | undefined => dataStore.nodesStatus[n.id]

const stateOf = (n: Node): CardState => {
  if (!n.enable) return 'disabled'
  const st = statusOf(n)
  // No snapshot yet: heartbeat hasn't reached this node since panel start.
  if (!st) return 'pending'
  return st.state
}

const stateColor = (n: Node): string => {
  switch (stateOf(n)) {
    case 'online': return 'var(--emerald)'
    case 'core-stopped': return 'var(--amber)'
    case 'offline': return 'var(--rose)'
    default: return 'var(--text-3)'
  }
}

const displayUrl = (n: Node): string => n.baseUrl.replace(/^https?:\/\//, '')

// ---------------- relative time ----------------
const rtf = new Intl.RelativeTimeFormat(locale, { numeric: 'auto' })
const relTime = (unix: number): string => {
  if (!unix) return t('ui.none')
  const diff = unix - Math.floor(Date.now() / 1000)
  const abs = Math.abs(diff)
  if (abs < 60) return rtf.format(Math.trunc(diff), 'second')
  if (abs < 3600) return rtf.format(Math.trunc(diff / 60), 'minute')
  if (abs < 86400) return rtf.format(Math.trunc(diff / 3600), 'hour')
  return rtf.format(Math.trunc(diff / 86400), 'day')
}

// ---------------- card rows ----------------
const pct = (v: number): string => Math.round(v) + '%'

const cardRows = (n: Node): EntityRow[] => {
  const state = stateOf(n)
  const st = statusOf(n)
  if (state === 'online' || state === 'core-stopped') {
    const memPct = st!.mem.total > 0 ? (st!.mem.current / st!.mem.total) * 100 : 0
    return [
      { k: t('node.latency'), v: st!.latency + ' ' + t('date.ms'), mono: true },
      { k: t('ui.cpu') + ' / ' + t('ui.memory'), v: pct(st!.cpu) + ' / ' + pct(memPct), mono: true },
      { k: t('node.panelVersion'), v: st!.appVersion || t('ui.none'), mono: !!st!.appVersion },
      state === 'online'
        ? { k: t('node.coreVersion'), v: st!.coreVersion || t('ui.none'), mono: !!st!.coreVersion }
        : { k: t('node.coreVersion'), v: t('node.status.coreStopped'), color: 'var(--amber)' },
    ]
  }
  if (state === 'offline') {
    return [
      { k: t('node.lastSeen'), v: relTime(st?.lastOnline || n.lastSeen) },
      { k: t('node.error'), v: st?.error ?? t('ui.none'), color: 'var(--rose)' },
      { k: t('node.panelVersion'), v: t('ui.none') },
      { k: t('node.coreVersion'), v: t('ui.none') },
    ]
  }
  // pending / disabled
  return [
    { k: t('node.latency'), v: t('ui.none') },
    { k: t('ui.cpu') + ' / ' + t('ui.memory'), v: t('ui.none') },
    { k: t('node.lastSeen'), v: relTime(n.lastSeen) },
    { k: t('node.coreVersion'), v: t('ui.none') },
  ]
}

// ---------------- drawer ----------------
const drawer = ref({ visible: false, id: 0, data: '' })
const openDrawer = (id: number) => {
  drawer.value.id = id
  drawer.value.data = id == 0 ? '{}' : JSON.stringify(nodes.value.findLast((n) => n.id == id))
  drawer.value.visible = true
}

// ---------------- delete (with confirm) ----------------
const del = ref({ visible: false, id: 0 })
const deleting = ref(false)
const askDelete = (id: number) => {
  del.value.id = id
  del.value.visible = true
}
const confirmDelete = async () => {
  deleting.value = true
  const success = await dataStore.save('nodes', 'del', del.value.id)
  deleting.value = false
  if (success) del.value.visible = false
}
</script>
