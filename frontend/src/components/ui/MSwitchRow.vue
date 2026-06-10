<template>
  <div style="display: flex; align-items: center; justify-content: space-between; gap: 12px; padding: 12px 14px; border-radius: 12px; background: var(--surface-3); margin-bottom: 14px;">
    <div>
      <div style="font-weight: 700; font-size: 13.5px;">{{ label }}</div>
      <div v-if="desc" style="font-size: 11.5px; color: var(--text-3);">{{ desc }}</div>
    </div>
    <Toggle :model-value="modelValue" @update:model-value="$emit('update:modelValue', $event)" />
  </div>
</template>

<script lang="ts" setup>
import Toggle from './Toggle.vue'

defineProps<{
  modelValue: boolean
  label: string
  desc?: string
}>()
defineEmits<{ 'update:modelValue': [value: boolean] }>()
</script>
