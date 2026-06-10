<template>
  <div style="display: flex; gap: 8px;">
    <input
      class="input mono"
      :value="modelValue"
      :placeholder="placeholder"
      :type="secret ? 'password' : 'text'"
      style="font-size: 12.5px;"
      @input="$emit('update:modelValue', ($event.target as HTMLInputElement).value)"
    />
    <button type="button" class="btn btn-ghost btn-icon" :title="title" @click="$emit('regenerate')">
      <Ico name="refresh" :size="15" />
    </button>
  </div>
</template>

<script lang="ts" setup>
import Ico from './Ico.vue'

withDefaults(defineProps<{
  modelValue?: string
  placeholder?: string
  secret?: boolean
  title?: string
}>(), {
  title: 'Generate',
})
defineEmits<{ 'update:modelValue': [value: string]; regenerate: [] }>()
</script>
