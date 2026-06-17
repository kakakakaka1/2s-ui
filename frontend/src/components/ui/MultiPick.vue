<template>
  <div style="display: flex; gap: 7px; flex-wrap: wrap;">
    <button
      v-for="o in options"
      :key="o"
      type="button"
      class="chip mono"
      :style="{
        cursor: 'pointer', height: '30px', padding: '0 12px',
        color: modelValue.includes(o) ? '#fff' : 'var(--text-2)',
        background: modelValue.includes(o) ? 'var(--brand)' : 'var(--surface-3)',
        borderColor: modelValue.includes(o) ? 'var(--brand)' : 'var(--line)',
      }"
      @click="toggle(o)"
    >{{ o }}</button>
  </div>
</template>

<script lang="ts" setup>
const props = defineProps<{
  modelValue: string[]
  options: string[]
}>()
const emit = defineEmits<{ 'update:modelValue': [value: string[]] }>()

const toggle = (v: string) => {
  emit('update:modelValue', props.modelValue.includes(v)
    ? props.modelValue.filter((x) => x !== v)
    : [...props.modelValue, v])
}
</script>
