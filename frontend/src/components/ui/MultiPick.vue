<template>
  <div style="display: flex; gap: 7px; flex-wrap: wrap;">
    <button
      v-for="o in allOptions"
      :key="o"
      type="button"
      class="chip mono"
      :title="isStale(o) ? o + ' — ' + $t('actions.del') : undefined"
      :style="chipStyle(o)"
      @click="toggle(o)"
    >{{ o }}</button>
  </div>
</template>

<script lang="ts" setup>
import { computed } from 'vue'

const props = defineProps<{
  modelValue: string[]
  options: string[]
}>()
const emit = defineEmits<{ 'update:modelValue': [value: string[]] }>()

// 已选中但已不在候选项里的「孤儿值」(例如指向已删除出站的引用)也要渲染出来，
// 否则用户既看不到也无法点掉它，残留进配置后会让 sing-box 因引用不存在的出站而启动失败。
const allOptions = computed(() => [
  ...props.options,
  ...props.modelValue.filter((v) => !props.options.includes(v)),
])

const isStale = (o: string) => props.modelValue.includes(o) && !props.options.includes(o)

const chipStyle = (o: string) => {
  const base = { cursor: 'pointer', height: '30px', padding: '0 12px' }
  if (isStale(o)) {
    return {
      ...base,
      color: 'var(--rose)',
      background: 'color-mix(in srgb, var(--rose) 14%, transparent)',
      borderColor: 'color-mix(in srgb, var(--rose) 45%, transparent)',
      textDecoration: 'line-through',
    }
  }
  const selected = props.modelValue.includes(o)
  return {
    ...base,
    color: selected ? '#fff' : 'var(--text-2)',
    background: selected ? 'var(--brand)' : 'var(--surface-3)',
    borderColor: selected ? 'var(--brand)' : 'var(--line)',
  }
}

const toggle = (v: string) => {
  emit('update:modelValue', props.modelValue.includes(v)
    ? props.modelValue.filter((x) => x !== v)
    : [...props.modelValue, v])
}
</script>
