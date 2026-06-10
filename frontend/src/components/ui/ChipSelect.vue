<template>
  <div class="field" style="gap: 6px; position: relative; margin-bottom: 0;">
    <label v-if="label">{{ label }}</label>
    <div
      class="input"
      style="min-height: 42px; height: auto; display: flex; flex-wrap: wrap; gap: 6px; align-items: center; cursor: pointer; padding: 7px 10px;"
      @click="open = !open"
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
    <template v-if="open">
      <div class="pop-scrim" @click="open = false" />
      <div
        class="card fade-up"
        style="position: absolute; top: 100%; inset-inline: 0; margin-top: 4px; z-index: 31; padding: 6px; max-height: 240px; overflow-y: auto; display: flex; flex-direction: column; gap: 1px; box-shadow: var(--shadow-pop);"
      >
        <label
          v-for="o in options"
          :key="o.value"
          style="display: flex; align-items: center; gap: 9px; padding: 8px 9px; border-radius: 7px; cursor: pointer; font-size: 12.5px; font-weight: 600; color: var(--text-2);"
          class="chip-select-opt"
          @click="toggle(o.value)"
        >
          <Check :checked="modelValue.includes(o.value)" /> {{ o.title }}
        </label>
      </div>
    </template>
  </div>
</template>

<script lang="ts" setup>
import { ref } from 'vue'
import Chip from './Chip.vue'
import Check from './Check.vue'
import Ico from './Ico.vue'

const props = defineProps<{
  modelValue: string[]
  options: { title: string; value: string }[]
  label?: string
  placeholder?: string
}>()
const emit = defineEmits<{ 'update:modelValue': [value: string[]] }>()

const open = ref(false)
const toggle = (v: string) => {
  emit('update:modelValue', props.modelValue.includes(v)
    ? props.modelValue.filter((x) => x !== v)
    : [...props.modelValue, v])
}
const titleOf = (v: string) => props.options.find((o) => o.value === v)?.title ?? v
</script>

<style scoped>
.chip-select-opt:hover { background: var(--surface-3); }
</style>
