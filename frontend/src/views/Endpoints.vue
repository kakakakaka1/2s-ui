<template>
  <EndpointDrawer
    :visible="drawer.visible"
    :id="drawer.id"
    :data="drawer.data"
    :tags="endpointTags"
    @close="drawer.visible = false"
  />
  <StatsModal
    :visible="stats.visible"
    :resource="stats.resource"
    :tag="stats.tag"
    @close="stats.visible = false"
  />
  <WgQrModal :visible="qrcode.visible" :data="qrcode.data" @close="qrcode.visible = false" />

  <!-- delete confirmation -->
  <DeleteConfirm :open="del.visible" :loading="deleting" @close="del.visible = false" @confirm="confirmDelete" />

  <div class="page-stack fade-up">
    <div class="toolbar" style="justify-content: center;">
      <Btn variant="primary" sm @click="openDrawer(0)">
        <Ico name="plus" :size="15" /> {{ $t('actions.add') }}
      </Btn>
      <Btn
        sm
        :loading="testingAll"
        :disabled="testingAll || endpoints.length === 0"
        @click="checkAllEndpoints"
      >
        <Ico name="chart" :size="15" /> {{ $t('actions.testAll') }}
      </Btn>
    </div>

    <div class="entity-grid">
      <EntityCard
        v-for="item in endpoints"
        :key="item.tag"
        :title="item.tag"
        :type="item.type"
        :color="item.type === 'tailscale' ? 'var(--violet)' : 'var(--cyan)'"
        icon="endpoint"
        :rows="cardRows(item)"
      >
        <template #chip>
          <Chip v-if="onlines.includes(item.tag)" color="emerald" dot>{{ $t('online') }}</Chip>
          <Chip v-else>{{ $t('ui.none') }}</Chip>
        </template>
        <template #actions>
          <CardBtn icon="edit" :title="$t('actions.edit')" @click="openDrawer(item.id)" />
          <CardBtn
            icon="bolt"
            border
            :title="$t('ui.delay')"
            :disabled="checkResults[item.tag]?.loading"
            @click="checkEndpoint(item.tag)"
          />
          <CardBtn icon="trash" border danger :title="$t('actions.del')" @click="askDelete(item.tag)" />
          <CardBtn
            v-if="item.type == 'wireguard' && item.peers?.length > 0"
            icon="qr"
            border
            title="QR"
            @click="showQrCode(item.id)"
          />
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
import { Endpoint } from '@/types/endpoints'
import Btn from '@/components/ui/Btn.vue'
import Ico from '@/components/ui/Ico.vue'
import Chip from '@/components/ui/Chip.vue'
import DeleteConfirm from '@/components/ui/DeleteConfirm.vue'
import CardBtn from '@/components/ui/CardBtn.vue'
import EntityCard, { EntityRow } from '@/components/ui/EntityCard.vue'
import EndpointDrawer from '@/layouts/drawers/endpoint/EndpointDrawer.vue'
import WgQrModal from '@/layouts/drawers/endpoint/WgQrModal.vue'
import StatsModal from '@/layouts/drawers/StatsModal.vue'

const { t } = useI18n({ useScope: 'global' })
const dataStore = Data()

// ---------------- store data ----------------
const endpoints = computed((): Endpoint[] => <Endpoint[]>dataStore.endpoints)

const endpointTags = computed((): string[] => endpoints.value?.map((o: Endpoint) => o.tag))

const onlines = computed(() => [
  ...dataStore.onlines.inbound ?? [],
  ...dataStore.onlines.outbound ?? [],
])

// ---------------- delay check ----------------
// Endpoints register as outbounds in the core, so the same api/checkOutbound
// latency probe used on the Outbounds page works here unchanged.
interface CheckResult {
  loading?: boolean
  success: boolean
  data?: { OK?: boolean; Delay?: number; Error?: string } | null
  errorMessage?: string
}

const checkResults = ref<Record<string, CheckResult>>({})

const checkEndpoint = async (tag: string) => {
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

const checkAllEndpoints = async () => {
  const list = endpoints.value
  if (list.length === 0) return
  testingAll.value = true
  try {
    await Promise.all(list.map((e) => checkEndpoint(e.tag)))
  } finally {
    testingAll.value = false
  }
}

// ---------------- card rows ----------------
const cardRows = (item: any): EntityRow[] => [
  {
    k: t('in.addr'),
    v: item.address?.length > 0 ? item.address[0] : t('ui.none'),
    mono: item.address?.length > 0,
  },
  {
    k: t('in.port'),
    v: item.listen_port > 0 ? item.listen_port : t('ui.none'),
    mono: item.listen_port > 0,
  },
  { k: t('types.wg.peers'), v: item.peers?.length ?? t('ui.none') },
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
  drawer.value.data = id == 0 ? '' : JSON.stringify(endpoints.value.findLast((o) => o.id == id))
  drawer.value.visible = true
}

const stats = ref({ visible: false, resource: 'endpoint', tag: '' })
const showStats = (tag: string) => {
  stats.value.tag = tag
  stats.value.visible = true
}

const qrcode = ref({ visible: false, data: <any>{} })
const showQrCode = (id: number) => {
  qrcode.value.data = endpoints.value.findLast((o) => o.id == id)
  qrcode.value.visible = true
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
  const success = await Data().save('endpoints', 'del', del.value.tag)
  if (success) del.value.visible = false
  deleting.value = false
}
</script>
