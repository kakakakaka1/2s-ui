<template>
  <div :style="{ height: height + 'px', borderRadius: '10px', background: 'var(--track)', overflow: 'hidden', width: '100%' }">
    <div :style="{ height: '100%', width: Math.min(100, value) + '%', background: color, borderRadius: '10px', transition: 'width .6s var(--ease)' }" />
  </div>
</template>

<script lang="ts" setup>
withDefaults(defineProps<{
  value: number
  color?: string
  height?: number
}>(), {
  color: 'var(--brand)',
  height: 6,
})
</script>
