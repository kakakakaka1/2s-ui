<template>
  <div ref="wrap" class="pop-wrap">
    <slot name="trigger" :toggle="toggle" :open="open" />
    <Teleport to="body">
      <template v-if="open">
        <div class="pop-scrim" @click="close" />
        <div ref="menu" class="card pop-menu fade-up" :style="menuStyle">
          <slot :close="close" />
        </div>
      </template>
    </Teleport>
  </div>
</template>

<script lang="ts" setup>
import { nextTick, onBeforeUnmount, ref } from 'vue'
import { pushOverlay } from './overlay'

const props = withDefaults(defineProps<{
  width?: number
  minWidth?: number
  align?: 'start' | 'end'
  direction?: 'down' | 'up'
}>(), {
  align: 'end',
  direction: 'down',
})

const wrap = ref<HTMLElement>()
const menu = ref<HTMLElement>()
const open = ref(false)
const menuStyle = ref<Record<string, string>>({})
let releaseEsc: (() => void) | undefined
let restoreFocus: HTMLElement | null = null

// Teleported to body so overflow ancestors (drawer bodies, table wraps) can't
// clip the menu; position is computed from the wrap rect, honoring align as a
// logical (RTL-aware) edge and flipping direction when the viewport is tight.
const place = () => {
  const r = wrap.value?.getBoundingClientRect()
  const m = menu.value
  if (!r || !m) return
  const gap = 6
  const margin = 8
  const mw = m.offsetWidth
  const mh = m.offsetHeight
  const rtl = getComputedStyle(wrap.value!).direction === 'rtl'

  // align=start: start edges line up; align=end: end edges line up
  const alignLeft = (props.align === 'start') !== rtl
  let left = alignLeft ? r.left : r.right - mw
  left = Math.min(Math.max(left, margin), Math.max(window.innerWidth - mw - margin, margin))

  const below = window.innerHeight - r.bottom - gap - margin
  const above = r.top - gap - margin
  const up = props.direction === 'up'
    ? above >= mh || above >= below
    : below < mh && above > below
  let top = up ? r.top - gap - mh : r.bottom + gap
  top = Math.min(Math.max(top, margin), Math.max(window.innerHeight - mh - margin, margin))

  menuStyle.value = {
    ...sizeStyle(),
    position: 'fixed',
    left: `${Math.round(left)}px`,
    top: `${Math.round(top)}px`,
    visibility: 'visible',
  }
}

const sizeStyle = (): Record<string, string> => ({
  ...(props.width ? { width: `${props.width}px` } : {}),
  ...(props.minWidth ? { minWidth: `${props.minWidth}px` } : {}),
  // min-width menus grow with content; never wider than the viewport
  maxWidth: `${Math.max(window.innerWidth - 16, 0)}px`,
})

// menu actions can reflow the trigger (e.g. picking a language resizes the
// button) — re-place while open whenever the anchor's rect changes
let raf = 0
let lastRect = ''
const track = () => {
  if (!open.value) return
  const r = wrap.value?.getBoundingClientRect()
  if (r) {
    const key = `${r.left},${r.top},${r.width},${r.height}`
    if (lastRect && key !== lastRect) place()
    lastRect = key
  }
  raf = requestAnimationFrame(track)
}

const toggle = async () => {
  if (open.value) return close()
  restoreFocus = document.activeElement as HTMLElement | null
  menuStyle.value = { ...sizeStyle(), position: 'fixed', left: '-9999px', top: '0', visibility: 'hidden' }
  open.value = true
  releaseEsc = pushOverlay(close)
  await nextTick()
  place()
  lastRect = ''
  cancelAnimationFrame(raf)
  raf = requestAnimationFrame(track)
}

const close = () => {
  if (!open.value) return
  open.value = false
  cancelAnimationFrame(raf)
  releaseEsc?.()
  releaseEsc = undefined
  if (restoreFocus && restoreFocus.isConnected) restoreFocus.focus({ preventScroll: true })
  restoreFocus = null
}

// close when the page moves under us, but not when the menu itself scrolls
const onWin = (e?: Event) => {
  if (!open.value) return
  if (e?.type === 'scroll' && menu.value && e.target instanceof Node && menu.value.contains(e.target)) return
  close()
}
window.addEventListener('resize', onWin)
window.addEventListener('scroll', onWin, true)
onBeforeUnmount(() => {
  window.removeEventListener('resize', onWin)
  window.removeEventListener('scroll', onWin, true)
  cancelAnimationFrame(raf)
  releaseEsc?.()
})

defineExpose({ close })
</script>
