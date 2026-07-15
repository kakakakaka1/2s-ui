<template>
  <label
    style="display: inline-flex; align-items: center; gap: 9px; cursor: pointer; font-size: 13px; font-weight: 600; color: var(--text-2);"
    @click.prevent="$emit('update:modelValue', !modelValue)"
  >
    <Check :checked="modelValue" />
    {{ label }}
  </label>
</template>

<script lang="ts" setup>
import Check from './Check.vue'

defineProps<{
  modelValue: boolean
  label: string
}>()
defineEmits<{ 'update:modelValue': [value: boolean] }>()
</script>
