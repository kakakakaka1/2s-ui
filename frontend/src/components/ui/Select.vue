<template>
  <button
    ref="trigger"
    type="button"
    class="input sel-trigger"
    v-bind="$attrs"
    :disabled="disabled"
    dir="auto"
    @click="toggle"
    @keydown.esc="close"
  >
    <span class="sel-value" :class="{ placeholder: current === null }">{{ display }}</span>
    <Ico name="chevronDown" :size="15" class="sel-caret" :style="open ? 'transform: rotate(180deg)' : ''" />
  </button>
  <!-- hidden: lets callers keep writing native <option>/<optgroup> markup -->
  <select ref="probe" style="display: none;" aria-hidden="true">
    <slot />
  </select>

  <Teleport to="body">
    <template v-if="open">
      <div class="sel-scrim" @click="close" @wheel.prevent />
      <div ref="panel" class="card sel-panel hide-scroll" :style="panelStyle">
        <template v-for="(it, i) in items" :key="i">
          <div v-if="it.group !== undefined" class="sel-group">{{ it.group }}</div>
          <button
            v-else
            type="button"
            class="sel-option"
            :class="{ active: isActive(it), disabled: it.disabled }"
            :disabled="it.disabled"
            @click="pick(it)"
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
import { computed, nextTick, onBeforeUnmount, onMounted, onUpdated, ref } from 'vue'
import Ico from './Ico.vue'

type Item = { group?: string; value?: any; label: string; disabled?: boolean }

defineOptions({ inheritAttrs: false })

const props = defineProps<{
  modelValue?: any
  disabled?: boolean
}>()
const emit = defineEmits<{ 'update:modelValue': [value: any]; change: [value: any] }>()

const trigger = ref<HTMLElement>()
const probe = ref<HTMLSelectElement>()
const panel = ref<HTMLElement>()
const open = ref(false)
const items = ref<Item[]>([])
const panelStyle = ref<Record<string, string>>({})

// Read options from the hidden native select so any markup the caller
// writes (static, v-for, optgroup) keeps working untouched.
const readItems = () => {
  const out: Item[] = []
  const walk = (nodes: HTMLCollection) => {
    for (const el of Array.from(nodes)) {
      if (el.tagName === 'OPTGROUP') {
        out.push({ group: (el as HTMLOptGroupElement).label, label: '' })
        walk(el.children)
      } else if (el.tagName === 'OPTION') {
        const o = el as HTMLOptionElement
        // mirror native behavior: no value attr -> text is the value
        const v = o.hasAttribute('value') ? (o as any)._value ?? o.value : o.text
        out.push({ value: v, label: o.text, disabled: o.disabled })
      }
    }
  }
  if (probe.value) walk(probe.value.children)
  items.value = out
}

// loose equality on purpose: callers bind numbers and strings interchangeably
const isActive = (it: Item) => it.value == props.modelValue
const current = computed(() => items.value.find((it) => it.group === undefined && isActive(it)) ?? null)
const display = computed(() => current.value?.label ?? (props.modelValue !== undefined && props.modelValue !== null ? String(props.modelValue) : ''))

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

const toggle = async () => {
  if (open.value) return close()
  readItems()
  place()
  open.value = true
  await nextTick()
  panel.value?.querySelector<HTMLElement>('.sel-option.active')?.scrollIntoView({ block: 'nearest' })
}
const close = () => { open.value = false }

const pick = (it: Item) => {
  close()
  if (it.value == props.modelValue) return
  emit('update:modelValue', it.value)
  emit('change', it.value)
}

// keep the displayed label in sync with slot-rendered options
onMounted(readItems)
onUpdated(readItems)

const onWin = () => { if (open.value) close() }
window.addEventListener('resize', onWin)
window.addEventListener('scroll', onWin, true)
onBeforeUnmount(() => {
  window.removeEventListener('resize', onWin)
  window.removeEventListener('scroll', onWin, true)
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
.sel-option:hover { background: var(--surface-3); color: var(--text); }
.sel-option.active { background: var(--brand-soft); color: var(--brand); }
.sel-option.disabled { opacity: .45; cursor: not-allowed; }
</style>
