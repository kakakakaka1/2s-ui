<template>
  <div class="card" style="padding: 0;">
    <div v-if="title" style="padding: 16px 20px; border-bottom: 1px solid var(--line);">
      <div style="font-size: 14.5px; font-weight: 700;">{{ title }}</div>
      <div v-if="desc" style="font-size: 12.5px; color: var(--text-3); margin-top: 3px;">{{ desc }}</div>
    </div>
    <div style="padding: 6px 20px 14px;">
      <slot />
    </div>
  </div>
</template>

<script lang="ts" setup>
defineProps<{
  title?: string
  desc?: string
}>()
</script>
