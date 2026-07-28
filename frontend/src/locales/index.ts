import { createI18n } from 'vue-i18n'
import en from './en'
import fa from './fa'
import vi from './vi'
import zhcn from './zhcn'
import zhtw from './zhtw'
import ru from './ru'

export const i18n = createI18n({
  legacy: false,
  locale: localStorage.getItem("locale") ?? 'en',
  fallbackLocale: 'en',
  messages: {
    en: en,
    fa: fa,
    vi: vi,
    zhHans: zhcn,
    zhHant: zhtw,
    ru: ru
  },
})

// Must be read per call, not snapshotted: setLocale() only assigns to
// i18n.global.locale, so a value captured at module load leaves every Intl and
// toLocale* consumer formatting in whichever language the tab was opened in.
export const intlLocale = () => {
  const l = i18n.global.locale.value
  switch (l) {
    case "zhHans":
      return "zh-cn"
    case "zhHant":
      return "zh-tw"
    default:
      return l
  }
}
