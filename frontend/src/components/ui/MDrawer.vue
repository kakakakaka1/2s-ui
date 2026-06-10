<template>
  <Drawer :open="open" :width="width" @close="$emit('close')">
    <template #title>
      <div style="display: flex; align-items: center; gap: 12px;">
        <div
          :style="{
            width: '38px', height: '38px', borderRadius: '10px', flex: 'none',
            background: `color-mix(in srgb, ${color} 14%, transparent)`, color,
            display: 'flex', alignItems: 'center', justifyContent: 'center',
          }"
        ><Ico :name="icon" :size="19" /></div>
        <div>
          <div style="font-size: 16px; font-weight: 800;">{{ title }}</div>
          <div style="font-size: 12px; color: var(--text-3);">{{ sub }}</div>
        </div>
      </div>
    </template>

    <!-- remount content each open so local form state resets -->
    <div :key="seq"><slot /></div>

    <template #footer>
      <Btn style="flex: 1;" @click="$emit('close')">{{ $t('actions.cancel') }}</Btn>
      <Btn variant="primary" style="flex: 1;" :loading="loading" @click="$emit('save')">
        <Ico name="check" :size="16" /> {{ saveLabel ?? $t('actions.save') }}
      </Btn>
    </template>
  </Drawer>
</template>

<script lang="ts" setup>
import { ref, watch } from 'vue'
import Drawer from './Drawer.vue'
import Btn from './Btn.vue'
import Ico from './Ico.vue'

const props = withDefaults(defineProps<{
  open: boolean
  icon?: string
  color?: string
  title?: string
  sub?: string
  saveLabel?: string
  width?: number
  loading?: boolean
}>(), {
  icon: 'outbound',
  color: 'var(--brand)',
  width: 500,
})
defineEmits<{ close: []; save: [] }>()

const seq = ref(0)
watch(() => props.open, (v) => { if (v) seq.value++ })
</script>
