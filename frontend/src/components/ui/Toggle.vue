<template>
  <button
    type="button"
    class="toggle"
    :class="{ on: modelValue }"
    :style="scale !== 1 ? { transform: `scale(${scale})` } : undefined"
    @click.prevent="$emit('update:modelValue', !modelValue)"
  />
</template>

<script lang="ts" setup>
withDefaults(defineProps<{
  modelValue?: boolean
  scale?: number
}>(), {
  modelValue: false,
  scale: 1,
})
defineEmits<{ 'update:modelValue': [value: boolean] }>()
</script>
