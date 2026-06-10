<template>
  <div class="card" style="padding: 0; overflow: hidden; display: flex; flex-direction: column;">
    <div style="padding: 13px 15px 10px; display: flex; align-items: center; gap: 11px;">
      <div
        class="mono"
        style="width: 30px; height: 30px; border-radius: 8px; flex: none; display: flex; align-items: center; justify-content: center; font-weight: 700; font-size: 13px; background: var(--surface-3); color: var(--text-2);"
      >{{ n }}</div>
      <div style="flex: 1; font-size: 12.5px; color: var(--text-3); font-weight: 600;">{{ modeLabel }}</div>
      <Chip :color="actColor">{{ action }}</Chip>
    </div>
    <div style="padding: 0 15px 12px; display: flex; flex-direction: column;">
      <div v-for="([k, v], i) in kvRows" :key="i" class="kv-row">
        <span class="k">{{ k }}</span>
        <span class="v">{{ v }}</span>
      </div>
    </div>
    <div style="display: flex; border-top: 1px solid var(--line); margin-top: auto;">
      <slot name="actions">
        <CardBtn icon="edit" @click="$emit('edit')" />
        <CardBtn icon="trash" border danger @click="$emit('del')" />
      </slot>
    </div>
  </div>
</template>

<script lang="ts" setup>
import { computed } from 'vue'
import { useI18n } from 'vue-i18n'
import Chip from './Chip.vue'
import CardBtn from './CardBtn.vue'

const props = defineProps<{
  n: number
  logical?: boolean
  mode?: string
  action: string
  targetKey: string
  target?: string
  rules: number | string
  invert?: boolean
}>()
defineEmits<{ edit: []; del: [] }>()

const { t } = useI18n()

const ACT_COLORS: Record<string, string> = {
  route: 'var(--brand)',
  reject: 'var(--rose)',
  'hijack-dns': 'var(--cyan)',
  'route-options': 'var(--violet)',
}
const actColor = computed(() => ACT_COLORS[props.action] || 'var(--text-3)')
const modeLabel = computed(() =>
  props.logical ? `${t('rule.logical')} (${props.mode})` : t('rule.simple')
)
const kvRows = computed<[string, string][]>(() => [
  [t('rule.action'), props.action],
  [props.targetKey, props.target || '—'],
  [t('rule.rules'), String(props.rules)],
  [t('rule.invert'), props.invert ? t('yes') : t('no')],
])
</script>
