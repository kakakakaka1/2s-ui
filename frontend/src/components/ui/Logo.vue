<template>
  <svg :width="size" :height="size" viewBox="0 0 40 40" fill="none">
    <defs>
      <linearGradient :id="gid" x1="4" y1="2" x2="36" y2="38" gradientUnits="userSpaceOnUse">
        <stop stop-color="#7b73ff" />
        <stop offset="1" stop-color="#1bc7d8" />
      </linearGradient>
    </defs>
    <rect x="2" y="2" width="36" height="36" rx="11" :fill="`url(#${gid})`" />
    <rect x="2" y="2" width="36" height="36" rx="11" fill="none" stroke="rgba(255,255,255,.25)" />
    <path
      d="M26.5 14.2c-1.2-1.7-3.3-2.7-5.8-2.7-3.5 0-6 1.8-6 4.6 0 2.5 1.8 3.7 5.2 4.4 2.6.5 3.4 1 3.4 2.1 0 1.2-1.2 1.9-3.1 1.9-2 0-3.4-.8-4.3-2.2"
      stroke="#fff" stroke-width="2.6" stroke-linecap="round" fill="none"
    />
    <circle cx="13.4" cy="26.4" r="2" fill="#fff" />
  </svg>
</template>

<script lang="ts" setup>
import { useId } from 'vue'

withDefaults(defineProps<{ size?: number }>(), { size: 30 })
const gid = `lg${useId().replace(/[:-]/g, '')}`
</script>
