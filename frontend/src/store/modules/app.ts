import { defineStore } from 'pinia'
import { i18n } from '@/locales'

export type Theme = 'dark' | 'light'

const THEME_KEY = '2sui-theme'
const RTL_LOCALES = new Set(['fa'])

function applyTheme(theme: Theme) {
  document.documentElement.dataset.theme = theme
}

function applyDir(locale: string) {
  document.documentElement.dir = RTL_LOCALES.has(locale) ? 'rtl' : 'ltr'
}

const AppStore = defineStore('App', {
  state: () => ({
    theme: (localStorage.getItem(THEME_KEY) as Theme) || 'dark',
    sidebarOpen: false, // mobile overlay sidebar
  }),
  getters: {
    isRtl: () => RTL_LOCALES.has(i18n.global.locale.value),
  },
  actions: {
    init() {
      applyTheme(this.theme)
      applyDir(i18n.global.locale.value)
    },
    toggleTheme() {
      this.theme = this.theme === 'dark' ? 'light' : 'dark'
      localStorage.setItem(THEME_KEY, this.theme)
      applyTheme(this.theme)
    },
    setLocale(l: string) {
      i18n.global.locale.value = l as any
      localStorage.setItem('locale', l)
      applyDir(l)
    },
  },
})

export default AppStore
