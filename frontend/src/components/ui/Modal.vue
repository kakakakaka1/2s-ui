<template>
  <Teleport to="body">
    <div v-if="open" class="modal-scrim" @click="$emit('close')">
      <div
        ref="panel"
        class="card"
        role="dialog"
        aria-modal="true"
        :aria-label="title"
        :style="{ width: `min(${width}px, 100%)`, maxHeight: '90vh', padding: 0, overflow: 'hidden', display: 'flex', flexDirection: 'column' }"
        @click.stop
      >
        <div style="display: flex; align-items: center; padding: 14px 18px; border-bottom: 1px solid var(--line); flex: none;">
          <div style="flex: 1; font-weight: 700; font-size: 14px;">{{ title }}</div>
          <button type="button" class="btn btn-subtle btn-icon btn-sm" @click="$emit('close')">
            <Ico name="close" :size="16" />
          </button>
        </div>
        <div style="flex: 1; overflow-y: auto; display: flex; flex-direction: column;">
          <slot />
        </div>
        <div
          v-if="$slots.footer"
          style="display: flex; justify-content: flex-end; gap: 10px; padding: 12px 18px; border-top: 1px solid var(--line); flex: none;"
        >
          <slot name="footer" />
        </div>
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
  title?: string
  width?: number
}>(), {
  width: 720,
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
