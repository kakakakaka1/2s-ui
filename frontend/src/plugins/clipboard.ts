import Clipboard from 'clipboard'
import { push } from 'notivue'
import { i18n } from '@/locales'

// Copy text with feedback toast. Uses clipboard.js so it also works
// on plain-HTTP panels where navigator.clipboard is unavailable.
export function copyToClipboard(txt: string, containerId?: string) {
  const hiddenButton = document.createElement('button')
  hiddenButton.className = 'clipboard-btn'
  document.body.appendChild(hiddenButton)

  const clipboard = new Clipboard('.clipboard-btn', {
    text: () => txt,
    container: (containerId ? document.getElementById(containerId) : undefined) ?? undefined,
  })

  clipboard.on('success', () => {
    clipboard.destroy()
    push.success({
      message: i18n.global.t('success') + ': ' + i18n.global.t('copyToClipboard'),
      duration: 5000,
    })
  })

  clipboard.on('error', () => {
    clipboard.destroy()
    push.error({
      message: i18n.global.t('failed') + ': ' + i18n.global.t('copyToClipboard'),
      duration: 5000,
    })
  })

  hiddenButton.click()
  document.body.removeChild(hiddenButton)
}
