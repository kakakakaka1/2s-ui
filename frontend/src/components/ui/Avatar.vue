<template>
  <div
    :style="{
      width: size + 'px', height: size + 'px',
      borderRadius: Math.round(size * 0.27) + 'px',
      flex: 'none',
      display: 'flex', alignItems: 'center', justifyContent: 'center',
      fontSize: size * 0.38 + 'px', fontWeight: 700, color: '#fff',
      background: `linear-gradient(135deg, hsl(${h} 62% 55%), hsl(${(h + 40) % 360} 62% 46%))`,
    }"
  >{{ letters }}</div>
</template>

<script lang="ts" setup>
import { computed } from 'vue'

const props = withDefaults(defineProps<{
  name: string
  size?: number
}>(), {
  size: 30,
})

const letters = computed(() =>
  props.name.replace(/[^a-z]/gi, '').slice(0, 2).toUpperCase() || '··'
)
const h = computed(() => {
  let v = 0
  for (const c of props.name) v = (v * 31 + c.charCodeAt(0)) % 360
  return v
})
</script>
