<template>
  <!-- 删除确认：证书删掉就不再续期，且入站 TLS 可能手填了这两个路径 -->
  <Modal :open="!!pendingDelete" :title="$t('setting.certDelete')" :width="460" @close="pendingDelete = ''">
    <div style="padding: 18px; font-size: 13.5px; line-height: 1.7;">
      <div>{{ $t('setting.certDeleteConfirm', { domain: pendingDelete }) }}</div>
      <div style="margin-top: 10px; color: var(--text-3);">{{ $t('setting.certDeleteTlsHint') }}</div>
    </div>
    <template #footer>
      <Btn @click="pendingDelete = ''">{{ $t('no') }}</Btn>
      <Btn style="color: var(--rose);" :loading="busy" @click="doDelete">
        <Ico name="trash" :size="15" /> {{ $t('yes') }}
      </Btn>
    </template>
  </Modal>

  <SettingsGroup :title="$t('setting.certs')" :desc="$t('setting.certsDesc')">
    <div v-if="loading && !certs.length" style="padding: 26px 0; text-align: center; color: var(--text-3); font-size: 13px;">
      {{ $t('loading') }}
    </div>

    <EmptyState
      v-else-if="!certs.length"
      icon="shield"
      :title="$t('setting.certsEmpty')"
      :desc="$t('setting.certsEmptyDesc')"
    />

    <div v-for="c in certs" :key="c.domain" class="certrow">
      <div class="cert-main">
        <div class="cert-dom">
          <span class="mono">{{ c.domain }}</span>
          <Chip :color="statusOf(c).color">{{ statusOf(c).text }}</Chip>
        </div>
        <div class="cert-meta">{{ metaOf(c) }}</div>

        <!-- 路径摆在明面上：建入站 TLS 时要填进 certificate_path / key_path -->
        <div class="paths">
          <div class="pathrow">
            <span class="pk">{{ $t('setting.sslCertShort') }}</span>
            <code class="mono">{{ c.certFile }}</code>
            <IconBtn name="copy" :title="$t('copyToClipboard')" @click="copyToClipboard(c.certFile)" />
          </div>
          <div class="pathrow">
            <span class="pk">{{ $t('setting.sslKeyShort') }}</span>
            <code class="mono">{{ c.keyFile }}</code>
            <IconBtn name="copy" :title="$t('copyToClipboard')" @click="copyToClipboard(c.keyFile)" />
          </div>
        </div>
      </div>

      <div class="cert-acts">
        <Btn variant="subtle" sm :loading="busy" @click="renew(c.domain)">
          <Ico name="refresh" :size="15" /> {{ $t('setting.forceRenew') }}
        </Btn>
        <Btn variant="subtle" sm style="color: var(--rose);" :disabled="busy" @click="pendingDelete = c.domain">
          <Ico name="trash" :size="15" /> {{ $t('actions.del') }}
        </Btn>
      </div>
    </div>
  </SettingsGroup>

  <SettingsGroup :title="$t('setting.certIssue')" :desc="$t('setting.acmeHint')" grid>
    <SRow :label="$t('setting.domain')">
      <input class="input mono" v-model="form.domain" placeholder="panel.example.com" @keyup.enter="issue" />
    </SRow>
    <SRow :label="$t('setting.acmeEmail')">
      <input class="input mono" v-model="form.email" placeholder="you@example.com" />
    </SRow>
    <SRow :label="$t('setting.acmeMethod')" :hint="methodHint">
      <Select v-model="form.method">
        <option value="auto">{{ $t('setting.acmeMethodAuto') }}</option>
        <option value="standalone">{{ $t('setting.acmeMethodStandalone') }}</option>
        <option value="nginx">{{ $t('setting.acmeMethodNginx') }}</option>
      </Select>
    </SRow>
    <div class="sg-span" style="padding: 13px 0 2px;">
      <Btn variant="primary" sm :loading="busy" :disabled="!form.domain" @click="issue">
        <Ico name="shield" :size="15" /> {{ $t('setting.issueCert') }}
      </Btn>
    </div>
  </SettingsGroup>
</template>

<script lang="ts" setup>
import { ref, reactive, computed, onMounted } from 'vue'
import { push } from 'notivue'
import { i18n } from '@/locales'
import HttpUtils from '@/plugins/httputil'
import { copyToClipboard } from '@/plugins/clipboard'
import Btn from '@/components/ui/Btn.vue'
import Ico from '@/components/ui/Ico.vue'
import IconBtn from '@/components/ui/IconBtn.vue'
import Chip from '@/components/ui/Chip.vue'
import Modal from '@/components/ui/Modal.vue'
import Select from '@/components/ui/Select.vue'
import SRow from '@/components/ui/SRow.vue'
import SettingsGroup from '@/components/ui/SettingsGroup.vue'
import EmptyState from '@/components/ui/EmptyState.vue'

interface Cert {
  domain: string
  certFile: string
  keyFile: string
  ca: string
  keyType: string
  notAfter: number
  nextRenew: number
  managed: boolean
}

const certs = ref<Cert[]>([])
const loading = ref(false)
const busy = ref(false)
const pendingDelete = ref('')
const form = reactive({ domain: '', email: '', method: 'auto' })
const detected = ref({ installed: false, active: false, port80Busy: false })

// 后端只给 unix 秒，天数在这里算：服务器时区不一定是用户的，后端算好会差一天
const daysLeft = (unixSec: number) => Math.floor((unixSec * 1000 - Date.now()) / 86400000)
const fmtDate = (unixSec: number) =>
  unixSec ? new Date(unixSec * 1000).toLocaleDateString(undefined, { year: 'numeric', month: '2-digit', day: '2-digit' }) : '—'

// color 只能取 Chip 的预设名，别的值会被它当成 CSS 颜色直接用
const statusOf = (c: Cert): { text: string; color: 'emerald' | 'amber' | 'rose' } => {
  if (!c.notAfter) return { text: i18n.global.t('setting.certUnreadable'), color: 'rose' }
  const d = daysLeft(c.notAfter)
  if (d < 0) return { text: i18n.global.t('setting.certExpired'), color: 'rose' }
  if (d <= 14) return { text: i18n.global.t('setting.certExpiring', { days: d }), color: 'amber' }
  return { text: i18n.global.t('setting.certValid'), color: 'emerald' }
}

const metaOf = (c: Cert) => {
  const parts = [c.ca, c.keyType].filter(Boolean)
  if (c.notAfter) {
    parts.push(i18n.global.t('setting.certExpiresOn', { date: fmtDate(c.notAfter), days: daysLeft(c.notAfter) }))
  }
  parts.push(c.nextRenew
    ? i18n.global.t('setting.certNextRenew', { date: fmtDate(c.nextRenew) })
    : i18n.global.t('setting.certNoAutoRenew'))
  return parts.join(' · ')
}

// 「自动」会走哪条路的实时提示，判定顺序与后端 resolveMethod 一致
const methodHint = computed(() => {
  if (!detected.value.port80Busy) return i18n.global.t('setting.acmeHintFree')
  if (detected.value.active) return i18n.global.t('setting.acmeHintNginx')
  return i18n.global.t('setting.acmeHint80Busy')
})

const load = async () => {
  loading.value = true
  const r = await HttpUtils.get('api/certs')
  if (r.success && Array.isArray(r.obj)) certs.value = r.obj
  loading.value = false
}

const detect = async () => {
  const r = await HttpUtils.get('api/detectNginx')
  if (r.success && r.obj) detected.value = r.obj
}

const issue = async () => {
  if (!form.domain) return
  busy.value = true
  const r = await HttpUtils.post('api/issueCert', {
    domain: form.domain.trim(),
    email: form.email.trim(),
    method: form.method,
    force: 'false',
  })
  busy.value = false
  if (!r.success) return
  push.success({
    title: i18n.global.t('success'),
    duration: 8000,
    message: i18n.global.t('setting.issuedVia', { method: r.obj?.method ?? form.method }),
  })
  form.domain = ''
  await load()
}

// 强制续期会立刻重签并占用 Let's Encrypt 限速额度（同域名约 5 张/周），
// 但它是修「续期卡住了」的唯一手段，所以留在每一行上，不额外加确认。
const renew = async (domain: string) => {
  busy.value = true
  const r = await HttpUtils.post('api/issueCert', { domain, method: 'auto', force: 'true' })
  busy.value = false
  if (!r.success) return
  push.success({ title: i18n.global.t('success'), duration: 6000, message: i18n.global.t('setting.issueCertOk') })
  await load()
}

const doDelete = async () => {
  busy.value = true
  const r = await HttpUtils.post('api/deleteCert', { domain: pendingDelete.value })
  busy.value = false
  pendingDelete.value = ''
  if (r.success) await load()
}

onMounted(async () => {
  await load()
  await detect()
})
</script>

<style scoped>
.certrow {
  display: flex; gap: 18px; align-items: flex-start;
  padding: 16px 0; border-bottom: 1px solid var(--line);
}
.certrow:last-child { border-bottom: none; }
.cert-main { flex: 1; min-width: 0; }
.cert-dom { display: flex; align-items: center; gap: 8px; flex-wrap: wrap; font-size: 14px; font-weight: 700; }
.cert-meta { font-size: 11.5px; color: var(--text-3); margin-top: 6px; font-variant-numeric: tabular-nums; }
.cert-acts { display: flex; gap: 6px; flex: none; }

.paths { margin-top: 11px; display: flex; flex-direction: column; gap: 6px; max-width: 640px; }
.pathrow {
  display: flex; align-items: center; gap: 9px;
  background: var(--surface-3); border: 1px solid var(--line);
  border-radius: 9px; padding: 4px 4px 4px 11px;
}
.pathrow .pk { font-size: 11px; font-weight: 700; color: var(--text-3); flex: none; }
.pathrow code {
  flex: 1; min-width: 0; font-size: 12px; color: var(--text-2);
  overflow-x: auto; white-space: nowrap; scrollbar-width: none;
}
.pathrow code::-webkit-scrollbar { display: none; }

@media (max-width: 820px) {
  .certrow { flex-direction: column; gap: 12px; }
  .cert-acts { width: 100%; }
}
</style>
