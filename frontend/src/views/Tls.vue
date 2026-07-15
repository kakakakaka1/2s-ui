<template>
  <TlsDrawer
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

    <div class="entity-grid">
      <EntityCard
        v-for="item in tlsConfigs"
        :key="item.id"
        :title="item.name"
        :type="item.server?.server_name?.length > 0 ? item.server.server_name : $t('ui.none')"
        color="var(--brand)"
        icon="tls"
        :rows="cardRows(item)"
      >
        <template #actions>
          <CardBtn icon="edit" :title="$t('actions.edit')" @click="openDrawer(item.id)" />
          <CardBtn
            v-if="tlsInbounds(item.id).length == 0"
            icon="trash"
            border
            danger
            :title="$t('actions.del')"
            @click="askDelete(item.id)"
          />
          <CardBtn icon="clone" border :title="$t('actions.clone')" @click="clone(item)" />
        </template>
      </EntityCard>
    </div>
  </div>
</template>

<script lang="ts" setup>
import { computed, ref } from 'vue'
import { useI18n } from 'vue-i18n'
import Data from '@/store/modules/data'
import { Inbound } from '@/types/inbounds'
import Btn from '@/components/ui/Btn.vue'
import Ico from '@/components/ui/Ico.vue'
import Modal from '@/components/ui/Modal.vue'
import CardBtn from '@/components/ui/CardBtn.vue'
import EntityCard, { EntityRow } from '@/components/ui/EntityCard.vue'
import TlsDrawer from '@/layouts/drawers/tls/TlsDrawer.vue'

const { t } = useI18n({ useScope: 'global' })
const dataStore = Data()

// ---------------- store data ----------------
const tlsConfigs = computed((): any[] => dataStore.tlsConfigs)

const inbounds = computed((): Inbound[] => dataStore.inbounds)

const tlsInbounds = (id: number): string[] =>
  inbounds.value.filter((i) => i.tls_id == id).map((i) => i.tag)

// ---------------- card rows ----------------
const cardRows = (item: any): EntityRow[] => [
  {
    k: t('pages.inbounds'),
    v: tlsInbounds(item.id).length > 0 ? tlsInbounds(item.id).length : t('ui.none'),
  },
  {
    k: 'ACME',
    v: t(item.server?.acme == undefined ? 'no' : 'yes'),
    color: item.server?.acme != undefined ? 'var(--emerald)' : undefined,
  },
  { k: 'ECH', v: t(item.server?.ech == undefined ? 'no' : 'yes') },
  {
    k: 'Reality',
    v: t(item.server?.reality == undefined ? 'no' : 'yes'),
    color: item.server?.reality != undefined ? 'var(--brand)' : undefined,
  },
]

// ---------------- drawer / clone ----------------
const drawer = ref({ visible: false, id: 0, data: '' })
const openDrawer = (id: number) => {
  drawer.value.id = id
  drawer.value.data = id == 0 ? '{}' : JSON.stringify(tlsConfigs.value.findLast((c) => c.id == id))
  drawer.value.visible = true
}

const clone = async (obj: any) => {
  const data = JSON.parse(JSON.stringify(obj))
  data.id = 0
  while (tlsConfigs.value.findIndex((c) => c.name == data.name) != -1) {
    data.name += '-copy'
  }
  await Data().save('tls', 'new', data)
}

// ---------------- delete (with confirm) ----------------
const del = ref({ visible: false, id: 0 })
const deleting = ref(false)
const askDelete = (id: number) => {
  del.value = { visible: true, id }
}
const confirmDelete = async () => {
  if (del.value.id === 0) return
  deleting.value = true
  const success = await Data().save('tls', 'del', del.value.id)
  if (success) del.value.visible = false
  deleting.value = false
}
</script>
