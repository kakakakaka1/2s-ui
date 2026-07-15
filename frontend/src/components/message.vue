<template>
  <Notivue v-slot="item">
    <NotivueSwipe :item="item">
      <Notification
        :item="item"
        :theme="theme"
        :dir="direction"
        :icons="outlinedIcons"
        :hideClose="true"
        @click="item.clear"
      />
    </NotivueSwipe>
  </Notivue>
</template>

<script lang="ts" setup>
import { Notivue, Notification, NotivueSwipe, outlinedIcons, pastelTheme, darkTheme } from 'notivue'
import { computed } from 'vue'
import AppStore from '@/store/modules/app'

const app = AppStore()

const theme = computed(() => {
  let currentTheme = app.theme === 'light' ? pastelTheme : darkTheme
  currentTheme = {
    ...currentTheme,
    '--nv-width': 'auto',
  }
  return currentTheme
})

const direction = computed(() => (app.isRtl ? 'rtl' : 'ltr'))
</script>

<style>
:root {
  --nv-z: 10020;
}
</style>
