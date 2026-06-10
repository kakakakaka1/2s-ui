<template>
  <svg
    :width="size"
    :height="size"
    viewBox="0 0 24 24"
    fill="none"
    stroke="currentColor"
    :stroke-width="stroke"
    stroke-linecap="round"
    stroke-linejoin="round"
    v-html="ICONS[name] ?? ''"
  />
</template>

<script lang="ts" setup>
import { ICONS } from './icons'

withDefaults(defineProps<{
  name: string
  size?: number | string
  stroke?: number | string
}>(), {
  size: 18,
  stroke: 1.75,
})
</script>
