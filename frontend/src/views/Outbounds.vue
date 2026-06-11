<template>
  <OutboundDrawer
    :visible="drawer.visible"
    :id="drawer.id"
    :data="drawer.data"
    :tags="outboundTags"
    @close="drawer.visible = false"
  />
  <BulkDrawer
    :visible="bulkVisible"
    :outbound-tags="outboundTags"
    @close="bulkVisible = false"
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
    <div class="toolbar" style="justify-content: center;">
      <Btn variant="primary" sm @click="openDrawer(0)">
        <Ico name="plus" :size="15" /> {{ $t('actions.add') }}
      </Btn>
      <Btn variant="primary" sm @click="bulkVisible = true">
        <Ico name="upload" :size="15" /> {{ $t('actions.addbulk') }}
      </Btn>
      <Btn
        sm
        :loading="testingAll"
        :disabled="testingAll || outbounds.length === 0"
        @click="checkAllOutbounds"
      >
        <Ico name="chart" :size="15" /> {{ $t('actions.testAll') }}
      </Btn>
    </div>

    <div class="entity-grid">
      <EntityCard
        v-for="item in outbounds"
        :key="item.tag"
        :title="item.tag"
        :type="item.type"
        :color="outColor(item.type)"
        icon="outbound"
        :rows="cardRows(item)"
      >
        <template #chip>
          <Chip v-if="onlines.includes(item.tag)" color="emerald" dot>{{ $t('online') }}</Chip>
          <Chip v-else-if="item.server == undefined">{{ $t('ui.coreLbl') }}</Chip>
          <Chip v-else>{{ $t('ui.none') }}</Chip>
        </template>
        <template #actions>
          <CardBtn icon="edit" :title="$t('actions.edit')" @click="openDrawer(item.id)" />
          <CardBtn icon="trash" border danger :title="$t('actions.del')" @click="askDelete(item.tag)" />
          <CardBtn
            v-if="dataStore.enableTraffic"
            icon="chart"
            border
            :title="$t('stats.graphTitle')"
            @click="showStats(item.tag)"
          />
        </template>
      </EntityCard>
    </div>
  </div>
</template>

<script lang="ts" setup>
import { computed, ref } from 'vue'
import { useI18n } from 'vue-i18n'
import Data from '@/store/modules/data'
import HttpUtils from '@/plugins/httputil'
import { outColor } from '@/plugins/colors'
import { Outbound } from '@/types/outbounds'
import Btn from '@/components/ui/Btn.vue'
import Ico from '@/components/ui/Ico.vue'
import Chip from '@/components/ui/Chip.vue'
import Modal from '@/components/ui/Modal.vue'
import CardBtn from '@/components/ui/CardBtn.vue'
import EntityCard, { EntityRow } from '@/components/ui/EntityCard.vue'
import OutboundDrawer from '@/layouts/drawers/outbound/OutboundDrawer.vue'
import BulkDrawer from '@/layouts/drawers/outbound/BulkDrawer.vue'
import StatsModal from '@/layouts/drawers/StatsModal.vue'

const { t } = useI18n({ useScope: 'global' })
const dataStore = Data()

// ---------------- store data ----------------
const outbounds = computed((): Outbound[] => <Outbound[]>dataStore.outbounds)

const outboundTags = computed((): string[] => [
  ...dataStore.outbounds?.map((o: Outbound) => o.tag),
  ...dataStore.endpoints?.map((e: any) => e.tag),
])

const onlines = computed(() => dataStore.onlines.outbound ?? [])

// ---------------- delay check (legacy logic) ----------------
interface CheckResult {
  loading?: boolean
  success: boolean
  data?: { OK?: boolean; Delay?: number; Error?: string } | null
  errorMessage?: string
}

const checkResults = ref<Record<string, CheckResult>>({})

const checkOutbound = async (tag: string) => {
  checkResults.value = { ...checkResults.value, [tag]: { loading: true, success: false } }
  const msg = await HttpUtils.get('api/checkOutbound', { tag })
  const success = msg.success && msg.obj?.OK
  const errorMessage = success ? undefined : (msg.obj?.Error ?? msg.msg ?? '')
  checkResults.value = {
    ...checkResults.value,
    [tag]: { loading: false, success, data: msg.obj ?? null, errorMessage },
  }
}

const testingAll = ref(false)

const checkAllOutbounds = async () => {
  const list = outbounds.value
  if (list.length === 0) return
  testingAll.value = true
  try {
    await Promise.all(list.map((o) => checkOutbound(o.tag)))
  } finally {
    testingAll.value = false
  }
}

// ---------------- card rows ----------------
const cardRows = (item: any): EntityRow[] => [
  { k: t('in.addr'), v: item.server ?? t('ui.none'), mono: !!item.server },
  { k: t('in.port'), v: item.server_port ?? t('ui.none'), mono: !!item.server_port },
  {
    k: t('objects.tls'),
    v: Object.hasOwn(item, 'tls') ? t(item.tls?.enabled ? 'enable' : 'disable') : t('ui.none'),
    color: item.tls?.enabled ? 'var(--emerald)' : undefined,
  },
  delayRow(item),
]

const delayRow = (item: any): EntityRow => {
  const r = checkResults.value[item.tag]
  if (r?.loading) return { k: t('out.delay'), v: '…', mono: true }
  if (r && r.loading == false) {
    if (r.success) {
      return { k: t('out.delay'), v: (r.data?.Delay ?? 0) + ' ' + t('date.ms'), mono: true, color: 'var(--emerald)' }
    }
    return { k: t('out.delay'), v: r.errorMessage || t('failed'), color: 'var(--rose)' }
  }
  return { k: t('out.delay'), v: t('ui.none') }
}

// ---------------- drawers / modals ----------------
const drawer = ref({ visible: false, id: 0, data: '' })
const openDrawer = (id: number) => {
  drawer.value.id = id
  drawer.value.data = id == 0 ? '' : JSON.stringify(outbounds.value.findLast((o) => o.id == id))
  drawer.value.visible = true
}

const bulkVisible = ref(false)

const stats = ref({ visible: false, resource: 'outbound', tag: '' })
const showStats = (tag: string) => {
  stats.value.tag = tag
  stats.value.visible = true
}

// ---------------- delete (with confirm) ----------------
const del = ref({ visible: false, tag: '' })
const deleting = ref(false)
const askDelete = (tag: string) => {
  del.value = { visible: true, tag }
}
const confirmDelete = async () => {
  if (del.value.tag.length === 0) return
  deleting.value = true
  const success = await Data().save('outbounds', 'del', del.value.tag)
  if (success) del.value.visible = false
  deleting.value = false
}
</script>
