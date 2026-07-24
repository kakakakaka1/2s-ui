<template>
  <!-- 可输入的域名框：值就是输入框里的文本，下拉只是建议。
       域名和证书是两件事——反代由 Cloudflare 终结、或者服务本来就跑 HTTP 时，域名
       照样要填（后端拿它校验 Host），但根本不需要证书，所以不能强制从列表里选。 -->
  <div class="dom-wrap">
    <div class="input dom-input" :class="{ focused: open }">
      <input
        ref="field"
        class="dom-field mono"
        type="text"
        dir="ltr"
        role="combobox"
        aria-autocomplete="list"
        autocapitalize="none"
        autocorrect="off"
        autocomplete="off"
        spellcheck="false"
        :aria-expanded="open"
        :aria-controls="open ? panelId : undefined"
        :aria-activedescendant="open && hl >= 0 ? optId(hl) : undefined"
        :value="modelValue"
        :placeholder="placeholder"
        :disabled="disabled"
        @input="onInput"
        @focus="show"
        @keydown="onKey"
      />
      <button
        type="button"
        class="dom-chev"
        tabindex="-1"
        :disabled="disabled"
        :aria-label="$t('setting.certsExisting')"
        @mousedown.prevent
        @click="toggle"
      >
        <Ico name="chevronDown" :size="15" :style="open ? 'transform: rotate(180deg)' : ''" />
      </button>
    </div>

    <Teleport to="body">
      <template v-if="open">
        <div class="dom-scrim" @mousedown="close" />
        <div :id="panelId" ref="panel" class="card dom-pop hide-scroll" role="listbox" :style="panelStyle">
          <div v-if="suggestions.length" class="dom-head">{{ $t('setting.certsExisting') }}</div>
          <button
            v-for="(c, i) in suggestions"
            :id="optId(i)"
            :key="c.domain"
            type="button"
            class="dom-item"
            role="option"
            :aria-selected="c.domain === modelValue"
            :class="{ on: c.domain === modelValue, hl: hl === i }"
            @mousedown.prevent="pick(c.domain)"
            @mouseenter="hl = i"
          >
            <span class="mono dom-item-text">{{ c.domain }}</span>
            <Chip :color="hintOf(c).color">{{ hintOf(c).text }}</Chip>
          </button>

          <div v-if="suggestions.length" class="dom-sep" />
          <button type="button" class="dom-add" @mousedown.prevent="emitIssue">
            {{ trimmed ? $t('setting.certIssueFor', { domain: trimmed }) : $t('setting.certIssueNew') }}
          </button>
        </div>
      </template>
    </Teleport>
  </div>
</template>

<script lang="ts" setup>
import { computed, nextTick, onBeforeUnmount, ref, useId } from 'vue'
import { i18n } from '@/locales'
import Ico from './Ico.vue'
import Chip from './Chip.vue'
import { pushOverlay } from './overlay'
import { type Cert, useCerts, loadCerts, daysLeft } from '@/plugins/certs'

const props = defineProps<{
  modelValue: string
  placeholder?: string
  disabled?: boolean
}>()
const emit = defineEmits<{ 'update:modelValue': [value: string]; issue: [domain: string] }>()

const certs = useCerts()
const uid = useId()
const panelId = `dom-${uid}`
const optId = (i: number) => `${panelId}-o${i}`

const field = ref<HTMLInputElement>()
const panel = ref<HTMLElement>()
const open = ref(false)
const hl = ref(-1)
const panelStyle = ref<Record<string, string>>({})
let releaseEsc: (() => void) | undefined

const trimmed = computed(() => (props.modelValue ?? '').trim())

// 通配符证书的域名不能当建议：点它会把字面 `*.example.com` 写进 webDomain，那既过
// 不了 Host 校验也存不进去（后端已拒绝）。它的证书照样能用——findCert 会把具体子域
// 名泛匹配到通配符证书上，用户直接输入 foo.example.com 即可。
const pickable = computed<Cert[]>(() => certs.value.filter((c) => !c.domain.startsWith('*.')))

// 输入即过滤。已经完整输入了某个域名时不过滤——那时列表该继续显示全部，方便改主意
// 换一个，否则下拉里只剩当前这一条，等于没有。
const suggestions = computed<Cert[]>(() => {
  const q = trimmed.value.toLowerCase()
  if (!q || pickable.value.some((c) => c.domain.toLowerCase() === q)) return pickable.value
  return pickable.value.filter((c) => c.domain.toLowerCase().includes(q))
})

const hintOf = (c: Cert): { text: string; color: string } => {
  if (!c.managed) return { text: i18n.global.t('setting.certManual'), color: '' }
  if (!c.notAfter) return { text: i18n.global.t('setting.certUnreadable'), color: 'rose' }
  const d = daysLeft(c.notAfter)
  if (d < 0) return { text: i18n.global.t('setting.certExpired'), color: 'rose' }
  return { text: i18n.global.t('setting.certDaysLeft', { days: d }), color: d <= 14 ? 'amber' : 'emerald' }
}

const place = () => {
  const r = field.value?.closest('.dom-input')?.getBoundingClientRect()
  if (!r) return
  const maxH = 280
  const below = window.innerHeight - r.bottom - 10
  const rows = suggestions.value.length + 2
  const up = below < Math.min(maxH, rows * 36 + 12) && r.top > below
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

const show = async () => {
  if (open.value || props.disabled) return
  // 先置 open 再 await（与 Select.vue 同序）。反过来的话，开头那个重入守卫在 await
  // 期间是虚设的：点 chevron 会同步触发 focus，两个 show() 都穿过守卫，pushOverlay
  // 被压两次而只记住最后一个释放句柄——先压的那个永远留在共享 Esc 栈里，整页的 Esc
  // 从此被捕获阶段的监听吃掉。
  open.value = true
  hl.value = -1
  place()
  bindWin()
  releaseEsc = pushOverlay(close)
  // 第一次展开才回源；之后用共享缓存，切 tab 来回点不会反复打接口。
  // 列表到位后重新定位一次——条数变了，向上翻转的判定会跟着变。
  await loadCerts()
  await nextTick()
  if (open.value) place()
}

const close = () => {
  if (!open.value) return
  open.value = false
  hl.value = -1
  unbindWin()
  releaseEsc?.()
  releaseEsc = undefined
}

const toggle = () => {
  if (open.value) {
    close()
    return
  }
  field.value?.focus()
  show()
}

const onInput = (e: Event) => {
  emit('update:modelValue', (e.target as HTMLInputElement).value)
  hl.value = -1
  if (!open.value) show()
  else nextTick(place)
}

const pick = (domain: string) => {
  emit('update:modelValue', domain)
  close()
  field.value?.focus()
}

const emitIssue = () => {
  close()
  emit('issue', trimmed.value)
}

const onKey = (e: KeyboardEvent) => {
  if (!open.value) {
    if (e.key === 'ArrowDown') {
      e.preventDefault()
      show()
    }
    return
  }
  const n = suggestions.value.length
  switch (e.key) {
    case 'ArrowDown':
      e.preventDefault()
      if (n) hl.value = hl.value < n - 1 ? hl.value + 1 : 0
      break
    case 'ArrowUp':
      e.preventDefault()
      if (n) hl.value = hl.value > 0 ? hl.value - 1 : n - 1
      break
    case 'Enter':
      // 没高亮任何一条时不拦截：用户输的就是最终值，回车该按原样提交
      if (hl.value >= 0 && hl.value < n) {
        e.preventDefault()
        pick(suggestions.value[hl.value].domain)
      } else {
        close()
      }
      break
    case 'Tab':
      close()
      break
  }
  if (hl.value >= 0) {
    nextTick(() => document.getElementById(optId(hl.value))?.scrollIntoView({ block: 'nearest' }))
  }
}

// 页面在浮层底下滚动时它会留在原地，跟着关掉最省事；浮层自身滚动要放行
const onWin = (e?: Event) => {
  if (!open.value) return
  if (e?.type === 'scroll' && panel.value && e.target instanceof Node && panel.value.contains(e.target)) return
  close()
}
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
.dom-wrap { position: relative; }
.dom-input { display: flex; align-items: center; gap: 6px; padding-right: 4px; }
.dom-input.focused { border-color: var(--brand); box-shadow: 0 0 0 4px var(--brand-soft); }
.dom-field {
  flex: 1; min-width: 0; height: 100%;
  border: none; background: transparent; outline: none;
  color: inherit; font: inherit; padding: 0;
}
.dom-field::placeholder { color: var(--text-3); }
.dom-chev {
  flex: none; width: 28px; height: 28px; padding: 0; border: none; cursor: pointer;
  display: grid; place-items: center; border-radius: 7px;
  background: transparent; color: var(--text-3);
  transition: background .18s var(--ease), color .18s var(--ease);
}
.dom-chev:hover { background: var(--surface); color: var(--text-2); }
.dom-chev:disabled { cursor: not-allowed; opacity: .5; }
.dom-chev :deep(svg) { transition: transform .18s var(--ease); }
</style>

<style>
/* 浮层 teleport 到 body，不能用 scoped */
.dom-scrim { position: fixed; inset: 0; z-index: 98; }
.dom-pop {
  z-index: 99;
  padding: 5px;
  overflow-y: auto;
  box-shadow: var(--shadow-pop);
  animation: fadeIn .14s var(--ease);
}
.dom-head {
  font-size: 10.5px; font-weight: 700; letter-spacing: .07em; text-transform: uppercase;
  color: var(--text-3); padding: 8px 9px 4px;
}
.dom-item {
  display: flex; align-items: center; justify-content: space-between; gap: 10px;
  width: 100%; border: none; background: transparent; cursor: pointer;
  border-radius: 7px; padding: 8px 9px;
  font-size: 13px; font-weight: 600; color: var(--text-2); text-align: start;
}
.dom-item-text { flex: 1; min-width: 0; overflow: hidden; text-overflow: ellipsis; white-space: nowrap; }
.dom-item:hover, .dom-item.hl { background: var(--surface-3); color: var(--text); }
.dom-item.on { background: var(--brand-soft); color: var(--brand); }
.dom-sep { height: 1px; background: var(--line); margin: 5px 3px; }
.dom-add {
  width: 100%; border: none; background: transparent; cursor: pointer;
  border-radius: 7px; padding: 9px;
  font-family: var(--font-ui); font-size: 12.5px; font-weight: 600;
  color: var(--text-3); text-align: start;
}
.dom-add:hover { background: var(--surface-3); color: var(--text-2); }
</style>
