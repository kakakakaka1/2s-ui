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
import { useLocale, useTheme } from 'vuetify'

const Theme = useTheme()
const { isRtl } = useLocale()

const theme = computed(() =>{
  let currenTheme = Theme.global.name.value == "light" ? pastelTheme : darkTheme
  currenTheme = {
    ...currenTheme,
    '--nv-width': 'auto',
  }
  return currenTheme
})

const direction = computed(() => {
  return isRtl.value ? 'rtl' : 'ltr'
})
</script>

<style>
:root {
  --nv-z: 10020;
}

</style>