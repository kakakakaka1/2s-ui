<template>
  <div style="display: flex; align-items: center; gap: 16px; padding: 13px 0; border-bottom: 1px solid var(--line);">
    <div style="flex: 1;">
      <div style="font-size: 13.5px; font-weight: 600;">{{ label }}</div>
      <div v-if="desc" style="font-size: 11.5px; color: var(--text-3); margin-top: 2px;">{{ desc }}</div>
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
