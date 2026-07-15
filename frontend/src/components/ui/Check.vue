<template>
  <span
    class="ui-check"
    :class="{ on: checked }"
    role="checkbox"
    :aria-checked="!!checked"
    tabindex="0"
    @click="$emit('toggle')"
    @keydown.space.prevent="kbd"
    @keydown.enter.prevent="kbd"
  >
    <Ico v-if="checked" name="check" :size="12" :stroke="2.4" />
  </span>
</template>

<script lang="ts" setup>
import Ico from './Ico.vue'

defineProps<{ checked?: boolean }>()
defineEmits<{ toggle: [] }>()

// synthesize a real click so keyboard rides the exact same event path as the
// mouse (bubbles to wrapping <label>/row handlers) — no double-toggle risk
const kbd = (e: KeyboardEvent) => (e.currentTarget as HTMLElement).click()
</script>

<style scoped>
.ui-check {
  width: 18px; height: 18px; flex: none;
  border-radius: 5px;
  border: 1.5px solid var(--line-2);
  background: transparent;
  display: inline-flex; align-items: center; justify-content: center;
  cursor: pointer; color: #fff;
  transition: all .12s;
}
.ui-check.on { border-color: var(--brand); background: var(--brand); }
.ui-check:focus-visible { outline: 2px solid var(--brand); outline-offset: 2px; }
</style>
