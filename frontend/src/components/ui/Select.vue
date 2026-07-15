<template>
  <button
    ref="trigger"
    type="button"
    class="input sel-trigger"
    v-bind="$attrs"
    :disabled="disabled"
    dir="auto"
    role="combobox"
    :aria-expanded="open"
    :aria-controls="open ? panelId : undefined"
    :aria-activedescendant="open && hl >= 0 ? optId(hl) : undefined"
    @click="toggle"
    @keydown="onTriggerKey"
  >
    <span class="sel-value" :class="{ placeholder: current === null }">{{ display }}</span>
    <Ico name="chevronDown" :size="15" class="sel-caret" :style="open ? 'transform: rotate(180deg)' : ''" />
  </button>

  <Teleport to="body">
    <template v-if="open">
      <div class="sel-scrim" @click="close" />
      <div :id="panelId" ref="panel" class="card sel-panel hide-scroll" role="listbox" :style="panelStyle" @mouseleave="hl = -1">
        <template v-for="(it, i) in items" :key="i">
          <div v-if="it.group !== undefined" class="sel-group">{{ it.group }}</div>
          <button
            v-else
            :id="optId(i)"
            type="button"
            class="sel-option"
            role="option"
            :class="{ active: isActive(it), hl: hl === i, disabled: it.disabled }"
            :aria-selected="isActive(it)"
            :disabled="it.disabled"
            tabindex="-1"
            @click="pick(it)"
            @mouseenter="hl = i"
          >
            <span style="flex: 1; min-width: 0; overflow: hidden; text-overflow: ellipsis;">{{ it.label }}</span>
            <Ico v-if="isActive(it)" name="check" :size="14" />
          </button>
        </template>
      </div>
    </template>
  </Teleport>
</template>

<script lang="ts" setup>
import { Fragment, computed, nextTick, onBeforeUnmount, ref, useId, useSlots, type VNode } from 'vue'
import Ico from './Ico.vue'
import { pushOverlay } from './overlay'

type Item = { group?: string; value?: any; label: string; disabled?: boolean }

defineOptions({ inheritAttrs: false })

const props = defineProps<{
  modelValue?: any
  disabled?: boolean
}>()
const emit = defineEmits<{ 'update:modelValue': [value: any]; change: [value: any] }>()

const slots = useSlots()
const uid = useId()
const panelId = `sel-${uid}`
const optId = (i: number) => `${panelId}-o${i}`

const trigger = ref<HTMLElement>()
const panel = ref<HTMLElement>()
const open = ref(false)
const hl = ref(-1)
const panelStyle = ref<Record<string, string>>({})
let releaseEsc: (() => void) | undefined

// Parse the caller's native <option>/<optgroup> markup straight from the slot
// vnode tree: bound values keep their type (numbers stay numbers) and no
// hidden-DOM probe or Vue internals are involved.
const textOf = (children: unknown): string => {
  if (children == null) return ''
  if (typeof children === 'string' || typeof children === 'number') return String(children)
  if (Array.isArray(children)) {
    return children
      .map((c) => (typeof c === 'object' && c !== null && 'children' in c ? textOf((c as VNode).children) : textOf(c)))
      .join('')
  }
  return ''
}

const collect = (nodes: VNode[], out: Item[]) => {
  for (const n of nodes) {
    if (n.type === Fragment) {
      if (Array.isArray(n.children)) collect(n.children as VNode[], out)
    } else if (n.type === 'optgroup') {
      out.push({ group: String(n.props?.label ?? ''), label: '' })
      if (Array.isArray(n.children)) collect(n.children as VNode[], out)
    } else if (n.type === 'option') {
      const label = textOf(n.children).trim()
      const hasValue = n.props !== null && 'value' in n.props
      out.push({
        // mirror native behavior: no value attr -> text is the value
        value: hasValue ? n.props!.value : label,
        label,
        disabled: n.props?.disabled !== undefined && n.props?.disabled !== false,
      })
    }
    // anything else (comments, whitespace) is ignored, like a native select
  }
}

const items = computed<Item[]>(() => {
  const out: Item[] = []
  collect((slots.default?.() ?? []) as VNode[], out)
  return out
})

// loose equality on purpose: callers bind numbers and strings interchangeably
const isActive = (it: Item) => it.value == props.modelValue
const current = computed(() => items.value.find((it) => it.group === undefined && isActive(it)) ?? null)
const display = computed(() => current.value?.label ?? (props.modelValue !== undefined && props.modelValue !== null ? String(props.modelValue) : ''))

const enabledIdx = computed(() =>
  items.value.flatMap((it, i) => (it.group === undefined && !it.disabled ? [i] : [])),
)

const place = () => {
  const r = trigger.value?.getBoundingClientRect()
  if (!r) return
  const maxH = 264
  const below = window.innerHeight - r.bottom - 10
  const up = below < Math.min(maxH, items.value.length * 34 + 12) && r.top > below
  panelStyle.value = {
    position: 'fixed',
    left: `${Math.round(r.left)}px`,
    width: `${Math.round(r.width)}px`,
    maxHeight: `${maxH}px`,
    ...(up
      ? { bottom: `${Math.round(window.innerHeight - r.top + 4)}px` }
      : { top: `${Math.round(r.bottom + 4)}px` }),
  }
}

const scrollToHl = () => {
  nextTick(() => {
    document.getElementById(optId(hl.value))?.scrollIntoView({ block: 'nearest' })
  })
}

const move = (dir: 1 | -1) => {
  const ids = enabledIdx.value
  if (ids.length === 0) return
  const pos = ids.indexOf(hl.value)
  hl.value = pos < 0
    ? dir === 1 ? ids[0] : ids[ids.length - 1]
    : ids[Math.min(Math.max(pos + dir, 0), ids.length - 1)]
  scrollToHl()
}

const toggle = async () => {
  if (open.value) return close()
  place()
  // highlight mirrors hover: only the selected option starts highlighted,
  // first ArrowDown lands on the first enabled option via move()
  hl.value = items.value.findIndex((it) => it.group === undefined && isActive(it))
  open.value = true
  bindWin()
  releaseEsc = pushOverlay(close)
  await nextTick()
  panel.value?.querySelector<HTMLElement>('.sel-option.active')?.scrollIntoView({ block: 'nearest' })
}

const close = () => {
  if (!open.value) return
  open.value = false
  unbindWin()
  releaseEsc?.()
  releaseEsc = undefined
  // mouse picks focus the option button; hand focus back so an enclosing
  // Modal/Drawer focus trap isn't left pointing at <body>
  trigger.value?.focus({ preventScroll: true })
}

const pick = (it: Item) => {
  close()
  if (it.value == props.modelValue) return
  emit('update:modelValue', it.value)
  emit('change', it.value)
}

// focus stays on the trigger while open; keyboard drives the highlighted option
const onTriggerKey = (e: KeyboardEvent) => {
  if (!open.value) {
    // Enter/Space already click the button natively; arrows open too
    if (e.key === 'ArrowDown' || e.key === 'ArrowUp') {
      e.preventDefault()
      toggle()
    }
    return
  }
  switch (e.key) {
    case 'ArrowDown': e.preventDefault(); move(1); break
    case 'ArrowUp': e.preventDefault(); move(-1); break
    case 'Home': e.preventDefault(); hl.value = enabledIdx.value[0] ?? -1; scrollToHl(); break
    case 'End': e.preventDefault(); hl.value = enabledIdx.value[enabledIdx.value.length - 1] ?? -1; scrollToHl(); break
    case 'Enter':
    case ' ': {
      e.preventDefault()
      const it = items.value[hl.value]
      if (it && it.group === undefined && !it.disabled) pick(it)
      break
    }
    case 'Tab': close(); break
  }
}

// close when the page moves under us, but not when the panel itself scrolls
// (keyboard scrollIntoView and wheel-scrolling long lists fire scroll too)
const onWin = (e?: Event) => {
  if (!open.value) return
  if (e?.type === 'scroll' && panel.value && e.target instanceof Node && panel.value.contains(e.target)) return
  close()
}
// only listen while the panel is open: this component is reused 100+ times,
// and a permanently-bound capture-phase scroll handler per instance would fire
// on every page scroll for nothing
const bindWin = () => {
  window.addEventListener('resize', onWin)
  window.addEventListener('scroll', onWin, true)
}
const unbindWin = () => {
  window.removeEventListener('resize', onWin)
  window.removeEventListener('scroll', onWin, true)
}
onBeforeUnmount(() => {
  unbindWin()
  releaseEsc?.()
})
</script>

<style scoped>
.sel-trigger {
  display: inline-flex;
  align-items: center;
  gap: 8px;
  text-align: start;
  cursor: pointer;
}
.sel-trigger:disabled { opacity: .55; cursor: not-allowed; }
.sel-value {
  flex: 1;
  min-width: 0;
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
}
.sel-caret { flex: none; color: var(--text-3); transition: transform .18s var(--ease); }
</style>

<style>
/* panel is teleported to body - unscoped */
.sel-scrim { position: fixed; inset: 0; z-index: 98; }
.sel-panel {
  z-index: 99;
  padding: 5px;
  overflow-y: auto;
  box-shadow: var(--shadow-pop);
  animation: fadeIn .14s var(--ease);
}
.sel-group {
  font-size: 10.5px;
  font-weight: 700;
  letter-spacing: .07em;
  text-transform: uppercase;
  color: var(--text-3);
  padding: 8px 9px 4px;
}
.sel-option {
  display: flex;
  align-items: center;
  gap: 8px;
  width: 100%;
  border: none;
  background: transparent;
  cursor: pointer;
  border-radius: 7px;
  padding: 8px 9px;
  font-family: var(--font-ui);
  font-size: 13px;
  font-weight: 600;
  color: var(--text-2);
  text-align: start;
}
.sel-option:hover, .sel-option.hl { background: var(--surface-3); color: var(--text); }
.sel-option.active { background: var(--brand-soft); color: var(--brand); }
.sel-option.disabled { opacity: .45; cursor: not-allowed; }
</style>
