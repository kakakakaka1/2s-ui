<template>
  <div class="pop-wrap">
    <slot name="trigger" :toggle="toggle" :open="open" />
    <template v-if="open">
      <div class="pop-scrim" @click="open = false" />
      <div
        class="card pop-menu fade-up"
        :class="{ start: align === 'start', up: direction === 'up' }"
        :style="{ width: width ? width + 'px' : undefined, minWidth: minWidth ? minWidth + 'px' : undefined }"
      >
        <slot :close="() => (open = false)" />
      </div>
    </template>
  </div>
</template>

<script lang="ts" setup>
import { ref } from 'vue'

withDefaults(defineProps<{
  width?: number
  minWidth?: number
  align?: 'start' | 'end'
  direction?: 'down' | 'up'
}>(), {
  align: 'end',
  direction: 'down',
})

const open = ref(false)
const toggle = () => (open.value = !open.value)
defineExpose({ close: () => (open.value = false) })
</script>
