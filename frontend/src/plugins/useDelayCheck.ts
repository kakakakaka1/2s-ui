import { ref } from 'vue'
import { useI18n } from 'vue-i18n'
import HttpUtils from '@/plugins/httputil'
import { EntityRow } from '@/components/ui/EntityCard.vue'

interface CheckResult {
  loading?: boolean
  success: boolean
  data?: { OK?: boolean; Delay?: number; Error?: string } | null
  errorMessage?: string
}

// Outbounds and endpoints both probe through api/checkOutbound — endpoints
// register as outbounds inside the core — so the whole flow was duplicated
// between the two views.
//
// `list` is a getter rather than a Ref because both callers hold a readonly
// computed; taking Ref would type-check and then fail to track. It is any[]
// rather than {tag: string}[] because Outbound and Endpoint resolve to
// `{ type: string; [k: string]: any }` — their InterfaceMap is a mapped type
// that never references the per-protocol interfaces — so a required `tag`
// bound does not hold.
export function useDelayCheck(list: () => any[]) {
  const { t } = useI18n({ useScope: 'global' })

  const checkResults = ref<Record<string, CheckResult>>({})
  const testingAll = ref(false)

  const check = async (tag: string) => {
    checkResults.value = { ...checkResults.value, [tag]: { loading: true, success: false } }
    const msg = await HttpUtils.get('api/checkOutbound', { tag })
    const success = msg.success && msg.obj?.OK
    const errorMessage = success ? undefined : (msg.obj?.Error ?? msg.msg ?? '')
    checkResults.value = {
      ...checkResults.value,
      [tag]: { loading: false, success, data: msg.obj ?? null, errorMessage },
    }
  }

  const checkAll = async () => {
    const items = list()
    if (items.length === 0) return
    testingAll.value = true
    try {
      await Promise.all(items.map((i) => check(i.tag)))
    } finally {
      testingAll.value = false
    }
  }

  const delayRow = (tag: string): EntityRow => {
    const r = checkResults.value[tag]
    if (r?.loading) return { k: t('out.delay'), v: '…', mono: true }
    if (r && r.loading == false) {
      if (r.success) {
        return { k: t('out.delay'), v: (r.data?.Delay ?? 0) + ' ' + t('date.ms'), mono: true, color: 'var(--emerald)' }
      }
      return { k: t('out.delay'), v: r.errorMessage || t('failed'), color: 'var(--rose)' }
    }
    return { k: t('out.delay'), v: t('ui.none') }
  }

  return { checkResults, testingAll, check, checkAll, delayRow }
}
