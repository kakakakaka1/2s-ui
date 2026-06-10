<template>
  <div style="display: flex; align-items: center; padding: 11px 20px; border-top: 1px solid var(--line); font-size: 13px;" :title="title">
    <span style="color: var(--text-3); width: 80px; flex: none; font-weight: 600;">{{ k }}</span>
    <Chip v-if="chip" :color="accent ? 'brand' : undefined">
      <span :class="{ mono }" dir="ltr">{{ v }}</span>
    </Chip>
    <span v-else :class="{ mono }" dir="ltr" style="font-weight: 600; min-width: 0; overflow: hidden; text-overflow: ellipsis; white-space: nowrap;">{{ v }}</span>
  </div>
</template>

<script lang="ts" setup>
import Chip from '@/components/ui/Chip.vue'

defineProps<{
  k: string
  v: string | number
  chip?: boolean
  mono?: boolean
  accent?: boolean
  title?: string
}>()
</script>
