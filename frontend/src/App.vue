<template>
  <div v-if="loading" class="app-loading">
    <span class="app-loading-spinner" />
    {{ $t('loading') }}
  </div>
  <Message />
  <router-view />
</template>

<script lang="ts" setup>
import Message from '@/components/message.vue'
import { inject, ref, Ref } from 'vue'

const loading: Ref = inject('loading') ?? ref(false)

// Change page title
document.title = "2S-UI " + document.location.hostname
</script>

<style>
.app-loading {
  position: fixed; inset: 0; z-index: 10000;
  display: flex; flex-direction: column; align-items: center; justify-content: center; gap: 16px;
  background: rgba(6, 8, 16, .55);
  backdrop-filter: blur(3px);
  color: #fff; font-size: 13.5px; font-weight: 600;
}
.app-loading-spinner {
  width: 42px; height: 42px;
  border: 3px solid rgba(255, 255, 255, .25);
  border-top-color: var(--brand);
  border-radius: 50%;
  animation: spin .8s linear infinite;
}
</style>
