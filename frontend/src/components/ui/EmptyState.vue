<template>
  <div class="empty">
    <div class="e-ico"><Ico :name="icon" :size="26" /></div>
    <h3>{{ title }}</h3>
    <p v-if="desc">{{ desc }}</p>
    <slot />
  </div>
</template>

<script lang="ts" setup>
import Ico from './Ico.vue'

withDefaults(defineProps<{
  icon?: string
  title: string
  desc?: string
}>(), {
  icon: 'search',
})
</script>
