<template>
  <div class="field" :style="{ marginBottom: mb + 'px' }">
    <label v-if="label">{{ label }}</label>
    <slot />
    <span v-if="hint" style="font-size: 11px; color: var(--text-3);">{{ hint }}</span>
  </div>
</template>

<script lang="ts" setup>
withDefaults(defineProps<{
  label?: string
  hint?: string
  mb?: number
}>(), {
  mb: 15,
})
</script>
