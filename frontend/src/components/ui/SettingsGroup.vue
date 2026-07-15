<template>
  <div class="card" style="padding: 0;">
    <div v-if="title" style="padding: 16px 20px; border-bottom: 1px solid var(--line);">
      <div style="font-size: 14.5px; font-weight: 700;">{{ title }}</div>
      <div v-if="desc" style="font-size: 12.5px; color: var(--text-3); margin-top: 3px;">{{ desc }}</div>
    </div>
    <div :class="['sg-body', { 'sg-grid': grid }]">
      <slot />
    </div>
  </div>
</template>

<script lang="ts" setup>
defineProps<{
  title?: string
  desc?: string
  /** 双列网格布局：把表单行排成两列以利用宽屏空间 */
  grid?: boolean
}>()
</script>

<style scoped>
.sg-body { padding: 6px 20px 14px; }

.sg-grid {
  display: grid;
  grid-template-columns: 1fr 1fr;
  column-gap: 40px;
  padding: 8px 22px 16px;
}
/* 双列下让控件撑满列宽（覆盖 SRow 默认 280px 限制） */
.sg-grid :deep(.srow-control) { max-width: none; }
/* 加 .sg-span 的元素独占整行（按钮、分组开关等） */
.sg-grid :deep(.sg-span) { grid-column: 1 / -1; }

@media (max-width: 920px) {
  .sg-grid { grid-template-columns: 1fr; column-gap: 0; padding: 6px 20px 14px; }
}
</style>
