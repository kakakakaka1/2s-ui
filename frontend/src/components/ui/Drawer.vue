<template>
  <Teleport to="body">
    <div
      class="drawer-scrim"
      :style="{ opacity: open ? 1 : 0, pointerEvents: open ? 'auto' : 'none' }"
      @click="$emit('close')"
    />
    <div
      ref="panel"
      class="drawer-panel"
      :class="{ closed: !open }"
      :style="{ width: width + 'px' }"
      role="dialog"
      aria-modal="true"
      :aria-hidden="!open"
      :inert="!open"
    >
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
import { nextTick, onBeforeUnmount, ref, watch } from 'vue'
import Ico from './Ico.vue'
import { lockScroll, pushOverlay, trapFocus } from './overlay'

const props = withDefaults(defineProps<{
  open: boolean
  width?: number
}>(), {
  width: 460,
})
const emit = defineEmits<{ close: [] }>()

const panel = ref<HTMLElement>()
let releases: (() => void)[] = []

const engage = async () => {
  releases.push(pushOverlay(() => emit('close'), { closeChildren: true }), lockScroll())
  await nextTick()
  if (panel.value) releases.push(trapFocus(panel.value))
}
const disengage = () => {
  releases.forEach((fn) => fn())
  releases = []
}

watch(() => props.open, (v) => (v ? engage() : disengage()), { immediate: true })
onBeforeUnmount(disengage)
</script>
