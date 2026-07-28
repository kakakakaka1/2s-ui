<template>
  <!-- 删除确认。acme.sh 托管的会停掉续期并删掉文件，手动登记的只是让面板忘掉路径，
       后果差得很远，两种文案分开写。 -->
  <Modal
    :open="!!pendingDelete"
    :title="pendingDelete?.managed ? $t('setting.certDelete') : $t('setting.certRemove')"
    :width="460"
    @close="pendingDelete = null"
  >
    <div style="padding: 18px; font-size: 13.5px; line-height: 1.7;">
      <template v-if="pendingDelete?.managed">
        <div>{{ $t('setting.certDeleteConfirm', { domain: pendingDelete.domain }) }}</div>
        <div style="margin-top: 10px; color: var(--text-3);">{{ $t('setting.certDeleteTlsHint') }}</div>
      </template>
      <div v-else-if="pendingDelete">{{ $t('setting.certRemoveConfirm', { domain: pendingDelete.domain }) }}</div>
    </div>
    <template #footer>
      <Btn @click="pendingDelete = null">{{ $t('no') }}</Btn>
      <Btn style="color: var(--rose);" :loading="busy" @click="doDelete">
        <Ico name="trash" :size="15" /> {{ $t('yes') }}
      </Btn>
    </template>
  </Modal>

  <!-- 登记 / 编辑自带证书。编辑时域名锁死：换域名等于换一条记录，不是改这一条。 -->
  <Modal
    :open="manualOpen"
    :title="manualEditing ? $t('setting.certEditPaths') : $t('setting.certRegister')"
    :width="560"
    @close="closeManual"
  >
    <div class="cert-form">
      <SRow :label="$t('setting.domain')">
        <input class="input mono" v-model="manualForm.domain" :disabled="manualEditing" placeholder="cf.example.com" />
      </SRow>
      <SRow :label="$t('setting.sslCert')">
        <input class="input mono" v-model="manualForm.certFile" placeholder="/etc/ssl/cert.pem" />
      </SRow>
      <SRow :label="$t('setting.sslKey')">
        <input class="input mono" v-model="manualForm.keyFile" placeholder="/etc/ssl/key.pem" />
      </SRow>
      <div style="font-size: 12px; color: var(--text-3); line-height: 1.6; padding-top: 12px;">
        {{ $t('setting.certRegisterHint') }}
      </div>
    </div>
    <template #footer>
      <Btn @click="closeManual">{{ $t('no') }}</Btn>
      <Btn variant="primary" :loading="busy" :disabled="!manualComplete" @click="saveManual">
        <Ico name="check" :size="15" /> {{ $t('actions.save') }}
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
          <!-- 手动登记的证书正常时只挂「手动登记」这一个 chip；快过期/读不到才把状态
               也摆出来，免得平时两个 chip 抢注意力。 -->
          <Chip v-if="c.managed || needsAttention(c)" :color="statusOf(c).color">{{ statusOf(c).text }}</Chip>
          <Chip v-if="!c.managed">{{ $t('setting.certManual') }}</Chip>
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
        <!-- 通配符证书是 DNS-01 签的,面板的 standalone/nginx 验证续不了它,按钮点了
             必报错,不如不摆。强制续期会立即重签、占用 Let's Encrypt 限速额度(约 5 张/
             周),点前用 Pop 气泡确认一下防误点——不用重弹窗,保持每行紧凑。 -->
        <Pop v-if="c.managed && !c.domain.startsWith('*.')" :min-width="264">
          <template #trigger="{ toggle }">
            <Btn variant="subtle" sm :loading="busy" @click="toggle">
              <Ico name="refresh" :size="15" /> {{ $t('setting.forceRenew') }}
            </Btn>
          </template>
          <template #default="{ close }">
            <div style="padding: 8px 10px 4px; font-size: 13px; font-weight: 700;">{{ $t('setting.forceRenew') }}</div>
            <div style="padding: 0 10px 8px; font-size: 12.5px; color: var(--text-3); line-height: 1.55;">{{ $t('setting.forceRenewConfirm') }}</div>
            <div style="display: flex; gap: 6px; padding: 2px;">
              <Btn sm style="flex: 1; color: var(--amber);" @click="close(); renew(c.domain);">{{ $t('yes') }}</Btn>
              <Btn variant="subtle" sm style="flex: 1;" @click="close()">{{ $t('no') }}</Btn>
            </div>
          </template>
        </Pop>
        <Btn v-else variant="subtle" sm :disabled="busy" @click="editManual(c)">
          <Ico name="edit" :size="15" /> {{ $t('setting.certEditPaths') }}
        </Btn>
        <Btn variant="subtle" sm style="color: var(--rose);" :disabled="busy" @click="pendingDelete = c">
          <Ico name="trash" :size="15" /> {{ c.managed ? $t('actions.del') : $t('setting.certRemove') }}
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
    <!-- 自带证书（Cloudflare 源证书、公司内部 CA）照样能进来，只是续期归用户自己管 -->
    <SRow :label="$t('setting.certHave')">
      <Btn variant="subtle" sm :disabled="busy" @click="openRegister">
        <Ico name="plus" :size="15" /> {{ $t('setting.certRegisterSwitch') }}
      </Btn>
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
import { i18n, intlLocale } from '@/locales'
import HttpUtils from '@/plugins/httputil'
import { copyToClipboard } from '@/plugins/clipboard'
import Btn from '@/components/ui/Btn.vue'
import Ico from '@/components/ui/Ico.vue'
import IconBtn from '@/components/ui/IconBtn.vue'
import Chip from '@/components/ui/Chip.vue'
import Modal from '@/components/ui/Modal.vue'
import Pop from '@/components/ui/Pop.vue'
import Select from '@/components/ui/Select.vue'
import SRow from '@/components/ui/SRow.vue'
import SettingsGroup from '@/components/ui/SettingsGroup.vue'
import EmptyState from '@/components/ui/EmptyState.vue'
import { type Cert, useCerts, loadCerts, daysLeft } from '@/plugins/certs'

const props = defineProps<{
  /** 从域名框「去申请」跳过来时带着域名，省得再输一遍 */
  initialDomain?: string
  /**
   * 设置页【表单当前值】里开着反代的域名。申请/续期要据此上传 behindProxy——
   * 后端读库拿到的是还没保存的旧值，会装上一条对不上号的续期重载钩子。
   */
  proxiedDomains?: string[]
}>()

// 共享清单：loadCerts(true) 写的就是这一份,别处(域名框建议、回执)也跟着更新
const certs = useCerts()
const loading = ref(false)
const busy = ref(false)
const pendingDelete = ref<Cert | null>(null)
const form = reactive({ domain: '', email: '', method: 'auto' })
const detected = ref({ installed: false, active: false, port80Busy: false })

const manualOpen = ref(false)
const manualEditing = ref(false)
const manualForm = reactive({ domain: '', certFile: '', keyFile: '' })
const manualComplete = computed(() =>
  !!manualForm.domain.trim() && !!manualForm.certFile.trim() && !!manualForm.keyFile.trim())

// 跟 TokenModal/ChangesModal 一样按应用语言格式化,不跟浏览器语言——同一个面板里
// 两种日期写法只会让人怀疑数据不对
const fmtDate = (unixSec: number) =>
  unixSec ? new Date(unixSec * 1000).toLocaleDateString(intlLocale(), { year: 'numeric', month: '2-digit', day: '2-digit' }) : '—'

// color 只能取 Chip 的预设名，别的值会被它当成 CSS 颜色直接用
const statusOf = (c: Cert): { text: string; color: 'emerald' | 'amber' | 'rose' } => {
  if (!c.notAfter) return { text: i18n.global.t('setting.certUnreadable'), color: 'rose' }
  const d = daysLeft(c.notAfter)
  if (d < 0) return { text: i18n.global.t('setting.certExpired'), color: 'rose' }
  if (d <= 14) return { text: i18n.global.t('setting.certExpiring', { days: d }), color: 'amber' }
  return { text: i18n.global.t('setting.certValid'), color: 'emerald' }
}

// 手动登记的证书没人替它续期，所以「快过期」这件事只有这里会说
const needsAttention = (c: Cert) => !c.notAfter || daysLeft(c.notAfter) <= 14

const metaOf = (c: Cert) => {
  const parts = [c.ca, c.keyType].filter(Boolean)
  if (c.notAfter) {
    parts.push(i18n.global.t('setting.certExpiresOn', { date: fmtDate(c.notAfter), days: daysLeft(c.notAfter) }))
  }
  if (c.managed) {
    parts.push(c.nextRenew
      ? i18n.global.t('setting.certNextRenew', { date: fmtDate(c.nextRenew) })
      : i18n.global.t('setting.certNoAutoRenew'))
  } else {
    parts.push(i18n.global.t('setting.certManualRenew'))
  }
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
  await loadCerts(true)
  loading.value = false
}

const detect = async () => {
  const r = await HttpUtils.get('api/detectNginx')
  if (r.success && r.obj) detected.value = r.obj
}

// 该域名是否由反代终结 TLS——按设置页的表单当前值判,与后端回退读库的口径一致但
// 时点更准(开关刚打开还没保存时,库里还是旧值)
const isProxied = (domain: string) =>
  (props.proxiedDomains ?? []).some((d) => d.trim().toLowerCase() === domain.trim().toLowerCase())

const issue = async () => {
  if (!form.domain) return
  const domain = form.domain.trim()
  const proxied = isProxied(domain)
  busy.value = true
  const r = await HttpUtils.post('api/issueCert', {
    domain,
    email: form.email.trim(),
    method: form.method,
    force: 'false',
    behindProxy: proxied ? 'true' : 'false',
  })
  busy.value = false
  if (!r.success) return
  let message = i18n.global.t('setting.issuedVia', { method: r.obj?.method ?? form.method })
  // 后端只在检测到 nginx 时才装续期重载钩子。Caddy/Traefik/HAProxy 拿不到钩子:
  // 续期只覆盖磁盘文件,代理仍用内存里的旧证书,到期才暴雷——必须当场明说。
  if (proxied && !r.obj?.reloadCmd) {
    message += ' ' + i18n.global.t('setting.proxyRenewHint')
  }
  push.success({ title: i18n.global.t('success'), duration: 8000, message })
  form.domain = ''
  await load()
}

// 强制续期会立刻重签并占用 Let's Encrypt 限速额度（同域名约 5 张/周），
// 但它是修「续期卡住了」的唯一手段，所以留在每一行上，不额外加确认。
const renew = async (domain: string) => {
  const proxied = isProxied(domain)
  busy.value = true
  const r = await HttpUtils.post('api/issueCert', {
    domain,
    method: 'auto',
    force: 'true',
    behindProxy: proxied ? 'true' : 'false',
  })
  busy.value = false
  if (!r.success) return
  let message = i18n.global.t('setting.issueCertOk')
  if (proxied && !r.obj?.reloadCmd) {
    message += ' ' + i18n.global.t('setting.proxyRenewHint')
  }
  push.success({ title: i18n.global.t('success'), duration: 6000, message })
  await load()
}

const openRegister = () => {
  manualEditing.value = false
  // 申请表单里已经填了域名的话带过来，省得再输一遍
  Object.assign(manualForm, { domain: form.domain.trim(), certFile: '', keyFile: '' })
  manualOpen.value = true
}

const editManual = (c: Cert) => {
  manualEditing.value = true
  Object.assign(manualForm, { domain: c.domain, certFile: c.certFile, keyFile: c.keyFile })
  manualOpen.value = true
}

const closeManual = () => {
  manualOpen.value = false
}

const saveManual = async () => {
  busy.value = true
  const r = await HttpUtils.post('api/saveManualCert', {
    domain: manualForm.domain.trim(),
    certFile: manualForm.certFile.trim(),
    keyFile: manualForm.keyFile.trim(),
  })
  busy.value = false
  if (!r.success) return
  manualOpen.value = false
  push.success({ title: i18n.global.t('success'), duration: 5000, message: i18n.global.t('setting.certRegistered') })
  await load()
}

const doDelete = async () => {
  if (!pendingDelete.value) return
  busy.value = true
  const r = await HttpUtils.post('api/deleteCert', { domain: pendingDelete.value.domain })
  busy.value = false
  pendingDelete.value = null
  if (r.success) await load()
}

onMounted(async () => {
  if (props.initialDomain) form.domain = props.initialDomain
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

/* Modal 里没有 sg-grid 那层，SRow 的控件会被压到 280px —— 路径填不下 */
.cert-form { padding: 6px 20px 16px; }
.cert-form :deep(.srow-control) { max-width: none; }

@media (max-width: 820px) {
  .certrow { flex-direction: column; gap: 12px; }
  .cert-acts { width: 100%; }
}
</style>
