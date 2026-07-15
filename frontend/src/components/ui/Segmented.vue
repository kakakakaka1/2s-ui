<template>
  <div class="segmented hide-scroll" :class="{ block }">
    <button
      v-for="[v, l] in options"
      :key="String(v)"
      type="button"
      :class="{ active: modelValue === v }"
      @click="$emit('update:modelValue', v)"
    >{{ l }}</button>
  </div>
</template>

<script lang="ts" setup>
defineProps<{
  modelValue: string | number
  // [value, label] pairs
  options: [string | number, string][]
  block?: boolean
}>()
defineEmits<{ 'update:modelValue': [value: string | number] }>()
</script>
