<template>
  <Teleport to="body">
    <div
      class="drawer-scrim"
      :style="{ opacity: open ? 1 : 0, pointerEvents: open ? 'auto' : 'none' }"
      @click="$emit('close')"
    />
    <div class="drawer-panel" :class="{ closed: !open }" :style="{ width: width + 'px' }">
      <div class="drawer-head">
        <slot name="title" />
        <button type="button" class="btn btn-subtle btn-icon" style="margin-inline-start: auto;" @click="$emit('close')">
          <Ico name="close" :size="18" />
        </button>
      </div>
      <div class="drawer-body hide-scroll">
        <slot />
      </div>
      <div v-if="$slots.footer" class="drawer-foot">
        <slot name="footer" />
      </div>
    </div>
  </Teleport>
</template>

<script lang="ts" setup>
import Ico from './Ico.vue'

withDefaults(defineProps<{
  open: boolean
  width?: number
}>(), {
  width: 460,
})
defineEmits<{ close: [] }>()
</script>
