<template>
  <div style="display: flex; align-items: center; gap: 10px; padding: 11px 14px; border-radius: 11px; border: 1px dashed var(--line-2); font-size: 12.5px; color: var(--text-3);">
    <Ico name="bolt" :size="16" style="flex: none;" />
    <slot />
  </div>
</template>

<script lang="ts" setup>
import Ico from './Ico.vue'
</script>
