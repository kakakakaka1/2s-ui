<template>
  <Modal :open="visible" :title="$t('main.stats.title')" :width="420" @close="$emit('close')">
    <div style="padding: 8px 20px 14px;">
      <div v-for="row in rows" :key="row.key" style="display: flex; align-items: center; gap: 12px; padding: 11px 0; border-bottom: 1px solid var(--line); font-size: 13px;">
        <span :style="{ color: row.color || 'var(--text-3)', display: 'flex' }"><Ico :name="row.icon" :size="16" /></span>
        <span style="color: var(--text-2); font-weight: 600; flex: 1;">{{ row.label }}</span>
        <span class="mono" dir="ltr" style="font-weight: 700;">{{ row.value }}</span>
      </div>
      <div style="display: flex; justify-content: flex-end; margin-top: 14px;">
        <Btn sm :loading="loading" @click="refresh"><Ico name="refresh" :size="14" /> {{ $t('actions.update') }}</Btn>
      </div>
    </div>
  </Modal>
</template>

<script lang="ts" setup>
import { computed, ref, watch } from 'vue'
import { useI18n } from 'vue-i18n'
import HttpUtils from '@/plugins/httputil'
import { HumanReadable } from '@/plugins/utils'
import Modal from '@/components/ui/Modal.vue'
import Btn from '@/components/ui/Btn.vue'
import Ico from '@/components/ui/Ico.vue'

const props = defineProps<{ visible: boolean }>()
defineEmits<{ close: [] }>()

const { t } = useI18n({ useScope: 'global' })
const loading = ref(false)
const info = ref<any>({})

const rows = computed(() => [
  { key: 'clients', icon: 'clients', label: t('pages.clients'), value: info.value.clients ?? 0, color: '' },
  { key: 'inbounds', icon: 'inbound', label: t('pages.inbounds'), value: info.value.inbounds ?? 0, color: '' },
  { key: 'outbounds', icon: 'outbound', label: t('pages.outbounds'), value: info.value.outbounds ?? 0, color: '' },
  { key: 'services', icon: 'services', label: t('pages.services'), value: info.value.services ?? 0, color: '' },
  { key: 'endpoints', icon: 'endpoint', label: t('pages.endpoints'), value: info.value.endpoints ?? 0, color: '' },
  { key: 'clientUp', icon: 'upload', label: t('stats.upload'), value: HumanReadable.sizeFormat(info.value.clientUp ?? 0), color: 'var(--violet)' },
  { key: 'clientDown', icon: 'download', label: t('stats.download'), value: HumanReadable.sizeFormat(info.value.clientDown ?? 0), color: 'var(--cyan)' },
  { key: 'total', icon: 'chart', label: t('main.stats.totalUsage'), value: HumanReadable.sizeFormat((info.value.clientUp ?? 0) + (info.value.clientDown ?? 0)), color: 'var(--brand)' },
])

const refresh = async () => {
  loading.value = true
  const data = await HttpUtils.get('api/status', { r: 'db' })
  if (data.success && data.obj) info.value = data.obj.db ?? data.obj
  loading.value = false
}

watch(() => props.visible, (v) => { if (v) refresh() })
</script>
