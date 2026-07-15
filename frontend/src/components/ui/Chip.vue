<template>
  <span class="chip" :class="presetClass" :style="customStyle">
    <span v-if="dot" class="dot" />
    <slot />
  </span>
</template>

<script lang="ts" setup>
import { computed } from 'vue'

const PRESETS = ['emerald', 'rose', 'amber', 'brand', 'cyan']

const props = defineProps<{
  // preset name (emerald/rose/amber/brand/cyan) or any CSS color → tinted chip
  color?: string
  dot?: boolean
}>()

const presetClass = computed(() =>
  props.color && PRESETS.includes(props.color) ? `chip-${props.color}` : ''
)
const customStyle = computed(() => {
  if (!props.color || PRESETS.includes(props.color)) return undefined
  return {
    color: props.color,
    background: `color-mix(in srgb, ${props.color} 12%, transparent)`,
    borderColor: `color-mix(in srgb, ${props.color} 30%, transparent)`,
  }
})
</script>
