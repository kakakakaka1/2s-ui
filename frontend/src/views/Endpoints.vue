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
import { Endpoint } from '@/types/endpoints'
import Btn from '@/components/ui/Btn.vue'
import Ico from '@/components/ui/Ico.vue'
import Chip from '@/components/ui/Chip.vue'
import Modal from '@/components/ui/Modal.vue'
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
]

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
