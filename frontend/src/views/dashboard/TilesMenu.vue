<template>
  <Pop :width="280" align="end">
    <template #trigger="{ toggle }">
      <Btn sm @click="toggle"><Ico name="settings" :size="15" /> {{ $t('ui.customize') }}</Btn>
    </template>
    <div style="padding: 8px;">
      <div style="display: flex; align-items: center; margin-bottom: 10px;">
        <div style="font-size: 13.5px; font-weight: 700;">{{ $t('ui.tilesTitle') }}</div>
        <Btn variant="subtle" sm style="margin-inline-start: auto; height: 26px;" @click="$emit('reset')">
          <Ico name="refresh" :size="13" /> {{ $t('ui.reset') }}
        </Btn>
      </div>
      <div v-for="g in TILE_GROUPS" :key="g" style="margin-bottom: 8px;">
        <div style="font-size: 10.5px; font-weight: 700; color: var(--text-3); letter-spacing: .05em; padding: 6px 2px 4px; text-transform: uppercase;">
          {{ $t('ui.' + g) }}
        </div>
        <label
          v-for="w in TILES.filter((x) => x.group === g)"
          :key="w.id"
          class="tile-row"
          @click.prevent="$emit('toggle', w.id)"
        >
          <span :style="{ flex: 1, color: show[w.id] ? 'var(--text)' : 'var(--text-3)', fontSize: '13px', fontWeight: 600 }">{{ $t('ui.' + w.k) }}</span>
          <Toggle :model-value="!!show[w.id]" :scale="0.85" />
        </label>
      </div>
    </div>
  </Pop>
</template>

<script lang="ts" setup>
import Pop from '@/components/ui/Pop.vue'
import Btn from '@/components/ui/Btn.vue'
import Ico from '@/components/ui/Ico.vue'
import Toggle from '@/components/ui/Toggle.vue'
import { TILES, TILE_GROUPS } from './tiles'

defineProps<{ show: Record<string, boolean> }>()
defineEmits<{ toggle: [id: string]; reset: [] }>()
</script>

<style scoped>
.tile-row {
  display: flex; align-items: center; gap: 10px;
  padding: 7px 8px; border-radius: 9px; cursor: pointer;
}
.tile-row:hover { background: var(--surface-3); }
</style>
