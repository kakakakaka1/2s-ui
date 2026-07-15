<template>
  <Pop :width="184">
    <template #trigger="{ toggle }">
      <button type="button" class="btn btn-ghost lang-btn" @click="toggle">
        <Ico class="lang-globe" name="globe" :size="16" />
        <span style="font-weight: 700;">{{ current.flag }}</span>
        <Ico class="lang-caret" name="chevronDown" :size="14" />
      </button>
    </template>
    <template #default="{ close }">
      <button
        v-for="l in LANGS"
        :key="l.value"
        type="button"
        class="pop-item"
        :class="{ active: l.value === locale }"
        @click="pick(l.value); close()"
      >
        <span class="lang-flag">{{ l.flag }}</span>
        <span style="flex: 1;">{{ l.title }}</span>
        <Ico v-if="l.value === locale" name="check" :size="15" />
      </button>
    </template>
  </Pop>
</template>

<script lang="ts" setup>
import { computed } from 'vue'
import { useI18n } from 'vue-i18n'
import Pop from '@/components/ui/Pop.vue'
import Ico from '@/components/ui/Ico.vue'
import AppStore from '@/store/modules/app'

const app = AppStore()
const { locale } = useI18n({ useScope: 'global' })

const LANGS = [
  { title: 'English', value: 'en', flag: 'EN' },
  { title: 'فارسی', value: 'fa', flag: 'FA' },
  { title: 'Tiếng Việt', value: 'vi', flag: 'VI' },
  { title: '简体中文', value: 'zhHans', flag: '中' },
  { title: '繁體中文', value: 'zhHant', flag: '繁' },
  { title: 'Русский', value: 'ru', flag: 'RU' },
]

const current = computed(() => LANGS.find((l) => l.value === locale.value) ?? LANGS[0])
const pick = (l: string) => app.setLocale(l)
</script>

<style scoped>
.lang-btn { gap: 7px; height: 38px; padding: 0 12px; }
.lang-flag {
  width: 26px; height: 26px; border-radius: 7px; flex: none;
  display: flex; align-items: center; justify-content: center;
  font-size: 12px; font-weight: 700;
  background: var(--surface-3); color: var(--text-2);
}
@media (max-width: 820px) {
  .lang-globe, .lang-caret { display: none; }
  .lang-btn { padding: 0 9px; }
}
</style>
