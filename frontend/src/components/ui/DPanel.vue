<template>
  <div class="card" style="display: flex; flex-direction: column; overflow: hidden;">
    <div :style="{ display: 'flex', alignItems: 'center', gap: '10px', padding: `16px ${headPad}px`, borderBottom: '1px solid var(--line)' }">
      <div>
        <div style="font-size: 14px; font-weight: 700;">{{ title }}</div>
        <div v-if="sub" style="font-size: 11.5px; color: var(--text-3); margin-top: 2px;">{{ sub }}</div>
      </div>
      <div style="margin-inline-start: auto;"><slot name="right" /></div>
    </div>
    <div :style="{ padding: pad + 'px', flex: 1 }">
      <slot />
    </div>
  </div>
</template>

<script lang="ts" setup>
import { computed } from 'vue'

const props = withDefaults(defineProps<{
  title: string
  sub?: string
  pad?: number
}>(), {
  pad: 20,
})

const headPad = computed(() => (props.pad === 0 ? 20 : props.pad))
</script>
