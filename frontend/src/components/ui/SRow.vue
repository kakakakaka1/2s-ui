<template>
  <div class="srow">
    <div class="srow-label">
      <div style="font-size: 13.5px; font-weight: 600;">{{ label }}</div>
      <div v-if="hint" style="font-size: 11.5px; color: var(--text-3); margin-top: 2px;">{{ hint }}</div>
    </div>
    <div class="srow-control"><slot /></div>
  </div>
</template>

<script lang="ts" setup>
defineProps<{
  label: string
  hint?: string
}>()
</script>

<style scoped>
.srow {
  display: flex; align-items: center; gap: 16px;
  padding: 13px 0; border-bottom: 1px solid var(--line);
}
.srow:last-child { border-bottom: none; }
.srow-label { width: 200px; flex: none; }
.srow-control { flex: 1; max-width: 280px; }
@media (max-width: 820px) {
  .srow { flex-direction: column; align-items: stretch; gap: 8px; }
  .srow-label { width: auto; }
  .srow-control { max-width: none; }
}
</style>
