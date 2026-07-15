<template>
  <div class="tabs" :class="{ 'page-tabs': page, 'hide-scroll': true }" :style="{ marginBottom: mb + 'px', overflowX: 'auto' }">
    <button
      v-for="[key, label] in tabs"
      :key="key"
      type="button"
      :class="{ active: modelValue === key }"
      @click="$emit('update:modelValue', key)"
    >{{ label }}</button>
  </div>
</template>

<script lang="ts" setup>
withDefaults(defineProps<{
  modelValue: string
  tabs: [string, string][]
  page?: boolean
  mb?: number
}>(), {
  mb: 18,
})
defineEmits<{ 'update:modelValue': [value: string] }>()
</script>
