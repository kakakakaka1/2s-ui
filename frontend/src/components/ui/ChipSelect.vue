<template>
  <div class="field" style="gap: 6px; position: relative; margin-bottom: 0;">
    <label v-if="label">{{ label }}</label>
    <div
      ref="anchor"
      class="input"
      style="min-height: 42px; height: auto; display: flex; flex-wrap: wrap; gap: 6px; align-items: center; cursor: pointer; padding: 7px 10px;"
      tabindex="0"
      role="combobox"
      :aria-expanded="open"
      :aria-controls="open ? panelId : undefined"
      :aria-activedescendant="open && hl >= 0 ? optId(hl) : undefined"
      @click="toggleOpen"
      @keydown="onKey"
    >
      <span v-if="modelValue.length === 0" style="color: var(--text-3);">{{ placeholder ?? '…' }}</span>
      <Chip
        v-for="v in modelValue"
        :key="v"
        color="brand"
        style="gap: 5px; cursor: pointer;"
        @click.stop="toggle(v)"
      >{{ titleOf(v) }} <Ico name="close" :size="11" /></Chip>
    </div>
    <Teleport to="body">
      <template v-if="open">
        <div class="pop-scrim" @click="close" />
        <div
          :id="panelId"
          ref="panel"
          class="card fade-up"
          role="listbox"
          aria-multiselectable="true"
          :style="panelStyle"
          @mouseleave="hl = -1"
        >
          <label
            v-for="(o, i) in options"
            :id="optId(i)"
            :key="o.value"
            style="display: flex; align-items: center; gap: 9px; padding: 8px 9px; border-radius: 7px; cursor: pointer; font-size: 12.5px; font-weight: 600; color: var(--text-2);"
            class="chip-select-opt"
            :class="{ hl: hl === i }"
            role="option"
            :aria-selected="modelValue.includes(o.value)"
            @click="toggle(o.value)"
            @mouseenter="hl = i"
          >
            <Check :checked="modelValue.includes(o.value)" tabindex="-1" aria-hidden="true" /> {{ o.title }}
          </label>
        </div>
      </template>
    </Teleport>
  </div>
</template>

<script lang="ts" setup>
import { nextTick, onBeforeUnmount, ref, useId, watch } from 'vue'
import Chip from './Chip.vue'
import Check from './Check.vue'
import Ico from './Ico.vue'
import { pushOverlay } from './overlay'

const props = defineProps<{
  modelValue: string[]
  options: { title: string; value: string }[]
  label?: string
  placeholder?: string
}>()
const emit = defineEmits<{ 'update:modelValue': [value: string[]] }>()

const uid = useId()
const panelId = `cs-${uid}`
const optId = (i: number) => `${panelId}-o${i}`

const anchor = ref<HTMLElement>()
const panel = ref<HTMLElement>()
const open = ref(false)
const hl = ref(-1)
const panelStyle = ref<Record<string, string>>({})
let releaseEsc: (() => void) | undefined

// teleported to body so drawer/modal overflow can't clip the option list
const place = () => {
  const r = anchor.value?.getBoundingClientRect()
  if (!r) return
  const maxH = 240
  const below = window.innerHeight - r.bottom - 10
  const up = below < Math.min(maxH, props.options.length * 37 + 12) && r.top > below
  panelStyle.value = {
    position: 'fixed',
    left: `${Math.round(r.left)}px`,
    width: `${Math.round(r.width)}px`,
    zIndex: '99',
    padding: '6px',
    maxHeight: `${maxH}px`,
    overflowY: 'auto',
    display: 'flex',
    flexDirection: 'column',
    gap: '1px',
    boxShadow: 'var(--shadow-pop)',
    ...(up
      ? { bottom: `${Math.round(window.innerHeight - r.top + 4)}px` }
      : { top: `${Math.round(r.bottom + 4)}px` }),
  }
}

// picking chips rewraps the anchor (height: auto) — keep the panel anchored
watch(() => props.modelValue, async () => {
  if (!open.value) return
  await nextTick()
  place()
})

const toggleOpen = async () => {
  if (open.value) return close()
  place()
  hl.value = -1
  open.value = true
  releaseEsc = pushOverlay(close)
  await nextTick()
}

const close = () => {
  if (!open.value) return
  open.value = false
  releaseEsc?.()
  releaseEsc = undefined
  anchor.value?.focus({ preventScroll: true })
}

// combobox pattern: focus stays on the anchor, keys drive the highlight —
// the teleported panel never needs to receive focus (works inside focus traps)
const onKey = (e: KeyboardEvent) => {
  if (!open.value) {
    if (e.key === 'Enter' || e.key === ' ' || e.key === 'ArrowDown') {
      e.preventDefault()
      toggleOpen()
    }
    return
  }
  switch (e.key) {
    case 'ArrowDown':
    case 'ArrowUp': {
      e.preventDefault()
      const n = props.options.length
      if (n === 0) return
      const dir = e.key === 'ArrowDown' ? 1 : -1
      hl.value = hl.value < 0
        ? (dir === 1 ? 0 : n - 1)
        : Math.min(Math.max(hl.value + dir, 0), n - 1)
      nextTick(() => document.getElementById(optId(hl.value))?.scrollIntoView({ block: 'nearest' }))
      break
    }
    case 'Enter':
    case ' ': {
      e.preventDefault()
      const o = props.options[hl.value]
      if (o) toggle(o.value)
      break
    }
    case 'Tab': close(); break
  }
}

const toggle = (v: string) => {
  emit('update:modelValue', props.modelValue.includes(v)
    ? props.modelValue.filter((x) => x !== v)
    : [...props.modelValue, v])
}
const titleOf = (v: string) => props.options.find((o) => o.value === v)?.title ?? v

const onWin = (e?: Event) => {
  if (!open.value) return
  if (e?.type === 'scroll' && panel.value && e.target instanceof Node && panel.value.contains(e.target)) return
  close()
}
window.addEventListener('resize', onWin)
window.addEventListener('scroll', onWin, true)
onBeforeUnmount(() => {
  window.removeEventListener('resize', onWin)
  window.removeEventListener('scroll', onWin, true)
  releaseEsc?.()
})
</script>

<style>
.chip-select-opt:hover, .chip-select-opt.hl { background: var(--surface-3); }
</style>
