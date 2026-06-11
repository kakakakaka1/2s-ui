<template>
  <ServiceDrawer
    :visible="drawer.visible"
    :id="drawer.id"
    :data="drawer.data"
    :in-tags="inTags"
    :ts-tags="tsTags"
    :ss-tags="ssTags"
    :tls-configs="tlsConfigs"
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

    <div class="entity-grid">
      <EntityCard
        v-for="item in services"
        :key="item.tag"
        :title="item.tag"
        :type="item.type"
        color="var(--brand)"
        icon="services"
        :rows="cardRows(item)"
      >
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
import { Srv } from '@/types/services'
import Btn from '@/components/ui/Btn.vue'
import Ico from '@/components/ui/Ico.vue'
import Modal from '@/components/ui/Modal.vue'
import CardBtn from '@/components/ui/CardBtn.vue'
import EntityCard, { EntityRow } from '@/components/ui/EntityCard.vue'
import ServiceDrawer from '@/layouts/drawers/service/ServiceDrawer.vue'

const { t } = useI18n({ useScope: 'global' })
const dataStore = Data()

// ---------------- store data ----------------
const services = computed((): Srv[] => <Srv[]>dataStore.services)

const tsTags = computed((): string[] =>
  dataStore.endpoints?.filter((o: any) => o.type == 'tailscale')?.map((o: any) => o.tag))

const ssTags = computed((): string[] =>
  dataStore.inbounds?.filter((o: any) => o.type == 'shadowsocks' && !o.users)?.map((o: any) => o.tag))

const inTags = computed((): string[] => [
  ...dataStore.inbounds?.map((o: any) => o.tag).filter((tag: any) => tag != null),
  ...dataStore.endpoints?.filter((e: any) => e.listen_port > 0).map((e: any) => e.tag),
])

const tlsConfigs = computed((): any[] => <any[]>dataStore.tlsConfigs)

// ---------------- card rows ----------------
const cardRows = (item: any): EntityRow[] => [
  { k: t('in.addr'), v: item.listen, mono: true },
  { k: t('in.port'), v: item.listen_port, mono: true },
  {
    k: t('objects.tls'),
    v: t(item.tls_id > 0 ? 'enable' : 'disable'),
    color: item.tls_id > 0 ? 'var(--emerald)' : undefined,
  },
]

// ---------------- drawer ----------------
const drawer = ref({ visible: false, id: 0, data: '' })
const openDrawer = (id: number) => {
  drawer.value.id = id
  drawer.value.data = id == 0 ? '' : JSON.stringify(services.value.findLast((o) => o.id == id))
  drawer.value.visible = true
}

// ---------------- delete (with confirm) ----------------
const del = ref({ visible: false, id: 0 })
const deleting = ref(false)
const askDelete = (id: number) => {
  del.value = { visible: true, id }
}
const confirmDelete = async () => {
  const item = services.value.find((i) => i.id == del.value.id)
  if (!item) return
  deleting.value = true
  const success = await Data().save('services', 'del', item.tag)
  if (success) del.value.visible = false
  deleting.value = false
}
</script>
