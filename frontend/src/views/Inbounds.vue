<template>
  <InboundDrawer
    :visible="drawer.visible"
    :id="drawer.id"
    :in-tags="inTags"
    :tls-configs="tlsConfigs"
    @close="drawer.visible = false"
  />
  <StatsModal
    :visible="stats.visible"
    :resource="stats.resource"
    :tag="stats.tag"
    @close="stats.visible = false"
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
    <!-- ===================== toolbar ===================== -->
    <div class="toolbar">
      <Segmented v-model="filter" :options="filterOptions" />
      <div class="grow" />
      <Btn variant="primary" sm @click="openDrawer(0)">
        <Ico name="plus" :size="16" /> {{ $t('ui.addInbound') }}
      </Btn>
    </div>

    <!-- ===================== cards ===================== -->
    <div class="entity-grid">
      <div v-for="item in rows" :key="item.id" class="card inb-card">
        <div class="inb-head">
          <div class="inb-ico" :style="iconStyle(item.type)"><Ico name="inbound" :size="19" /></div>
          <div style="flex: 1; min-width: 0;">
            <div class="inb-tag">{{ item.tag }}</div>
            <div class="inb-type" :style="{ color: protoColor(item.type) }">{{ item.type }}</div>
          </div>
          <Chip v-if="isOnline(item.tag)" color="emerald" dot>{{ $t('ui.live') }}</Chip>
          <Chip v-else>{{ $t('ui.idle') }}</Chip>
        </div>
        <div class="inb-meta">
          <div>
            <div class="m-k">{{ $t('ui.listen') }}</div>
            <div class="m-v mono" dir="ltr">{{ listenOf(item) }}</div>
          </div>
          <div>
            <div class="m-k">{{ $t('ui.tlsLbl') }}</div>
            <Chip v-if="item.tls_id > 0" color="emerald">{{ $t('ui.enabled') }}</Chip>
            <div v-else class="m-v">{{ $t('ui.off') }}</div>
          </div>
          <div>
            <div class="m-k">{{ $t('ui.clientsCol') }}</div>
            <div class="m-v" :title="usersTitle(item)">{{ item.users ? item.users.length : '—' }}</div>
          </div>
          <div>
            <div class="m-k">{{ $t('ui.trafficCol') }}</div>
            <div class="m-v mono">{{ trafficOf(item) }}</div>
          </div>
        </div>
        <div class="inb-actions">
          <CardBtn icon="edit" :label="$t('ui.edit')" :title="$t('actions.edit')" @click="openDrawer(item.id)" />
          <CardBtn icon="clone" :label="$t('ui.clone')" border :title="$t('actions.clone')" @click="clone(item.id)" />
          <CardBtn
            v-if="dataStore.enableTraffic"
            icon="chart"
            :label="$t('ui.stats')"
            border
            :title="$t('stats.graphTitle')"
            @click="showStats(item.tag)"
          />
          <CardBtn icon="trash" border danger :title="$t('actions.del')" @click="askDelete(item.tag)" />
        </div>
      </div>
    </div>
  </div>
</template>

<script lang="ts" setup>
import { computed, ref } from 'vue'
import { useI18n } from 'vue-i18n'
import Data from '@/store/modules/data'
import { createInbound } from '@/types/inbounds'
import RandomUtil from '@/plugins/randomUtil'
import { protoColor } from '@/plugins/colors'
import { HumanReadable } from '@/plugins/utils'
import Btn from '@/components/ui/Btn.vue'
import Ico from '@/components/ui/Ico.vue'
import Chip from '@/components/ui/Chip.vue'
import Modal from '@/components/ui/Modal.vue'
import Segmented from '@/components/ui/Segmented.vue'
import CardBtn from '@/components/ui/CardBtn.vue'
import InboundDrawer from '@/layouts/drawers/inbound/InboundDrawer.vue'
import StatsModal from '@/layouts/drawers/StatsModal.vue'

const { t } = useI18n({ useScope: 'global' })
const dataStore = Data()

// ---------------- store data ----------------
const inbounds = computed((): any[] => dataStore.inbounds ?? [])
const tlsConfigs = computed((): any[] => dataStore.tlsConfigs ?? [])
const inTags = computed((): string[] => [
  ...(inbounds.value?.map((i: any) => i.tag) ?? []),
  ...(dataStore.endpoints?.filter((e: any) => e.listen_port > 0).map((e: any) => e.tag) ?? []),
])
const onlines = computed((): string[] => dataStore.onlines.inbound ?? [])
const isOnline = (tag: string): boolean => onlines.value.includes(tag)

// ---------------- filter ----------------
const filter = ref<string | number>('all')
const filterOptions = computed((): [string | number, string][] => [
  ['all', t('ui.any')],
  ['tls', t('ui.tlsOnly')],
  ['online', t('ui.online')],
])
const rows = computed((): any[] =>
  inbounds.value.filter((i: any) => {
    if (filter.value === 'tls') return i.tls_id > 0
    if (filter.value === 'online') return isOnline(i.tag)
    return true
  }),
)

// ---------------- cells ----------------
const iconStyle = (type: string) => {
  const col = protoColor(type)
  return { color: col, background: `color-mix(in srgb, ${col} 14%, transparent)` }
}
const listenOf = (item: any): string =>
  item.listen_port ? `${item.listen ?? ''}:${item.listen_port}` : '—'
const usersTitle = (item: any): string | undefined =>
  item.users && item.users.length > 0 ? item.users.join('\n') : undefined
const trafficOf = (item: any): string => {
  if (!item.users) return '—'
  const total = (dataStore.clients ?? [])
    .filter((c: any) => item.users.includes(c.name))
    .reduce((sum: number, c: any) => sum + (c.up ?? 0) + (c.down ?? 0), 0)
  return HumanReadable.sizeFormat(total)
}

// ---------------- drawer ----------------
const drawer = ref({ visible: false, id: 0 })
const openDrawer = (id: number) => {
  drawer.value.id = id
  drawer.value.visible = true
}

// ---------------- clone ----------------
const cloning = ref(false)
const clone = async (id: number) => {
  if (cloning.value) return
  cloning.value = true
  const inboundArray = await Data().loadInbounds([id])
  const inbound = inboundArray[0]
  const newTag = inbound.type + '-' + RandomUtil.randomSeq(3)
  const newInbound = createInbound(inbound.type, {
    ...inbound,
    id: 0,
    tag: newTag,
    listen_port: RandomUtil.randomIntRange(10000, 60000),
  })
  await Data().save('inbounds', 'new', newInbound)
  cloning.value = false
}

// ---------------- delete (with confirm) ----------------
const del = ref({ visible: false, tag: '' })
const deleting = ref(false)
const askDelete = (tag: string) => {
  del.value = { visible: true, tag }
}
const confirmDelete = async () => {
  if (del.value.tag === '') return
  deleting.value = true
  const success = await Data().save('inbounds', 'del', del.value.tag)
  if (success) del.value.visible = false
  deleting.value = false
}

// ---------------- stats ----------------
const stats = ref({ visible: false, resource: 'inbound', tag: '' })
const showStats = (tag: string) => {
  stats.value.tag = tag
  stats.value.visible = true
}
</script>

<style scoped>
.inb-card {
  padding: 0;
  overflow: hidden;
  display: flex;
  flex-direction: column;
}
.inb-head {
  padding: 15px 16px 13px;
  display: flex;
  align-items: center;
  gap: 11px;
}
.inb-ico {
  width: 38px;
  height: 38px;
  border-radius: 10px;
  flex: none;
  display: flex;
  align-items: center;
  justify-content: center;
}
.inb-tag {
  font-weight: 700;
  font-size: 14px;
  white-space: nowrap;
  overflow: hidden;
  text-overflow: ellipsis;
}
.inb-type {
  font-size: 11.5px;
  font-weight: 600;
}
.inb-meta {
  padding: 0 16px 14px;
  display: grid;
  grid-template-columns: 1fr 1fr;
  gap: 10px 14px;
}
.m-k {
  font-size: 11px;
  color: var(--text-3);
  font-weight: 600;
  margin-bottom: 3px;
}
.m-v {
  font-size: 13px;
  font-weight: 600;
}
.inb-actions {
  display: flex;
  border-top: 1px solid var(--line);
  margin-top: auto;
}
</style>
