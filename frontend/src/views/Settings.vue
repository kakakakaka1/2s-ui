<template>
  <EditorModal
    :open="jsonEditor"
    :title="$t('editor') + ' - ' + $t('setting.jsonSub')"
    :content="settings.subJsonExt"
    @save="saveJsonEditor"
    @close="jsonEditor = false"
  />
  <EditorModal
    :open="clashEditor"
    :title="$t('editor') + ' - ' + $t('setting.clashSub')"
    :content="settings.subClashExt"
    @save="saveClashEditor"
    @close="clashEditor = false"
  />

  <div class="page-stack-lg fade-up" style="max-width: 760px;">
    <!-- page tabs + actions -->
    <div class="head-row">
      <Tabs v-model="tab" page :mb="0" :tabs="tabItems" style="border-bottom: none; flex: 1 1 auto; min-width: 0;" />
      <div class="head-actions">
        <Btn variant="primary" sm :loading="loading" :disabled="!stateChange" @click="save">
          <Ico name="check" :size="15" /> {{ $t('actions.save') }}
        </Btn>
        <Pop :min-width="210">
          <template #trigger="{ toggle }">
            <Btn sm style="color: var(--amber);" :loading="loading" :disabled="stateChange" @click="toggle">
              <Ico name="refresh" :size="15" /> {{ $t('actions.restartApp') }}
            </Btn>
          </template>
          <template #default="{ close }">
            <div style="padding: 8px 10px 4px; font-size: 13px; font-weight: 700;">{{ $t('actions.restartApp') }}</div>
            <div style="padding: 0 10px 8px; font-size: 12.5px; color: var(--text-3);">{{ $t('confirm') }}</div>
            <div style="display: flex; gap: 6px; padding: 2px;">
              <Btn sm style="flex: 1; color: var(--amber);" @click="close(); restartApp();">{{ $t('yes') }}</Btn>
              <Btn variant="subtle" sm style="flex: 1;" @click="close()">{{ $t('no') }}</Btn>
            </div>
          </template>
        </Pop>
      </div>
    </div>

    <!-- ===================== Interface ===================== -->
    <SettingsGroup v-if="tab === 'interface'">
      <SRow :label="$t('setting.addr')">
        <input class="input mono" v-model="settings.webListen" placeholder="0.0.0.0" />
      </SRow>
      <SRow :label="$t('setting.port')">
        <input class="input mono" type="number" min="1" v-model.number="webPort" />
      </SRow>
      <SRow :label="$t('setting.webPath')">
        <input class="input mono" v-model="settings.webPath" />
      </SRow>
      <SRow :label="$t('setting.domain')">
        <input class="input mono" v-model="settings.webDomain" placeholder="panel.example.com" />
      </SRow>
      <SRow :label="$t('setting.deployNginx')" :hint="webNginxHint">
        <select class="input" v-model="settings.webNginx">
          <option value="false">{{ $t('no') }}</option>
          <option value="true">{{ $t('yes') }}</option>
        </select>
      </SRow>
      <SRow :label="$t('setting.acmeEmail')" :hint="$t('setting.acmeHint')">
        <input class="input mono" v-model="settings.webAcmeEmail" placeholder="you@example.com" />
      </SRow>
      <div class="srow-btn">
        <Btn variant="subtle" sm :loading="loading" :disabled="!settings.webDomain" @click="issueCert">
          <Ico name="shield" :size="15" /> {{ $t('setting.issueCert') }}
        </Btn>
      </div>
      <template v-if="settings.webNginx !== 'true'">
        <ToggleRow v-model="webAcme" :label="$t('setting.autoCert')" />
        <template v-if="!webAcme">
          <SRow :label="$t('setting.sslKey')">
            <input class="input mono" v-model="settings.webKeyFile" placeholder="/path/key.pem" />
          </SRow>
          <SRow :label="$t('setting.sslCert')">
            <input class="input mono" v-model="settings.webCertFile" placeholder="/path/cert.pem" />
          </SRow>
        </template>
      </template>
      <SRow :label="$t('setting.webUri')">
        <input class="input mono" v-model="settings.webURI" placeholder="https://panel.example.com/app/" />
      </SRow>
      <SRow :label="$t('setting.sessionAge')" :hint="$t('date.m')">
        <input class="input mono" type="number" min="0" v-model.number="sessionMaxAge" />
      </SRow>
      <SRow :label="$t('setting.trafficAge')" :hint="$t('date.d')">
        <input class="input mono" type="number" min="0" v-model.number="trafficAge" />
      </SRow>
      <SRow :label="$t('setting.timeLoc')">
        <input class="input mono" v-model="settings.timeLocation" />
      </SRow>
    </SettingsGroup>

    <!-- ===================== Subscription ===================== -->
    <SettingsGroup v-else-if="tab === 'sub'">
      <ToggleRow v-model="subEncode" :label="$t('setting.subEncode')" />
      <ToggleRow v-model="subShowInfo" :label="$t('setting.subInfo')" />
      <SRow :label="$t('setting.addr')">
        <input class="input mono" v-model="settings.subListen" placeholder="0.0.0.0" />
      </SRow>
      <SRow :label="$t('setting.port')">
        <input class="input mono" type="number" min="1" v-model.number="subPort" />
      </SRow>
      <ToggleRow v-model="subAcme" :label="$t('setting.autoCert')" />
      <SRow v-if="subAcme" :label="$t('setting.acmeEmail')" :hint="$t('setting.acmeHint')">
        <input class="input mono" v-model="settings.subAcmeEmail" placeholder="you@example.com" />
      </SRow>
      <template v-else>
        <SRow :label="$t('setting.sslKey')">
          <input class="input mono" v-model="settings.subKeyFile" placeholder="/path/key.pem" />
        </SRow>
        <SRow :label="$t('setting.sslCert')">
          <input class="input mono" v-model="settings.subCertFile" placeholder="/path/cert.pem" />
        </SRow>
      </template>
      <SRow :label="$t('setting.domain')">
        <input class="input mono" v-model="settings.subDomain" placeholder="sub.example.com" />
      </SRow>
      <SRow :label="$t('setting.path')">
        <input class="input mono" v-model="settings.subPath" />
      </SRow>
      <SRow :label="$t('setting.update')" :hint="$t('date.h')">
        <input class="input mono" type="number" min="0" v-model.number="subUpdates" />
      </SRow>
      <SRow :label="$t('setting.subUri')">
        <input class="input mono" v-model="settings.subURI" placeholder="https://sub.example.com/sub/" />
      </SRow>
    </SettingsGroup>

    <!-- ===================== JSON sub ===================== -->
    <SettingsGroup v-else-if="tab === 'jsonSub'">
      <div style="font-size: 12.5px; color: var(--text-3); padding: 12px 0 14px;">{{ $t('ui.jsonExtDesc') }}</div>

      <div class="field-grid">
        <ChipSelect v-model="ruleToDirect" :options="geoList" :label="$t('setting.toDirect')" :placeholder="$t('ui.selectHint')" />
        <ChipSelect v-model="ruleToBlock" :options="geoList" :label="$t('setting.toBlock')" :placeholder="$t('ui.selectHint')" />
      </div>

      <template v-if="enableLog">
        <div class="sub-label">{{ $t('basic.log.title') }}</div>
        <div class="field-grid">
          <Field :label="$t('basic.log.level')" :mb="0">
            <select class="input" v-model="subJsonExt.log.level">
              <option v-for="l in levels" :key="l" :value="l">{{ l }}</option>
            </select>
          </Field>
          <div style="display: flex; align-items: flex-end; padding-bottom: 6px;">
            <SwitchLabel v-model="subJsonExt.log.timestamp" :label="$t('setting.timestamp')" />
          </div>
        </div>
      </template>

      <template v-if="enableDns">
        <div class="sub-label">{{ $t('ui.dnsOpt') }}</div>
        <div class="field-grid">
          <Field :label="$t('dns.final')" :mb="0">
            <select class="input" v-model="subJsonExt.dns.final">
              <option v-for="tg in dnsTags" :key="tg" :value="tg">{{ tg }}</option>
            </select>
          </Field>
          <Field :label="$t('basic.routing.defaultDns')" :mb="0">
            <select class="input" v-model="defaultResolver">
              <option value="">{{ $t('ui.none') }}</option>
              <option v-for="tg in dnsTags" :key="tg" :value="tg">{{ tg }}</option>
            </select>
          </Field>
        </div>
        <div class="grid2" style="margin-top: 14px;">
          <Field :label="$t('setting.globalDns')" :mb="0">
            <div style="display: flex; gap: 8px;">
              <select class="input" style="flex: 1; min-width: 0;" :value="proxyDns.type" @change="setDnsType(proxyDns, ($event.target as HTMLSelectElement).value)">
                <option v-for="dt in dnsTypes" :key="dt" :value="dt">{{ dt }}</option>
              </select>
              <template v-if="proxyDns.type !== 'local'">
                <input class="input mono" style="flex: 1.5; min-width: 0;" v-model="proxyDns.server" :placeholder="$t('in.addr')" />
                <input class="input mono" style="width: 76px; flex: none;" type="number" min="1" v-model.number="proxyDns.server_port" :placeholder="$t('in.port')" />
              </template>
            </div>
          </Field>
          <Field :label="$t('setting.directDns')" :mb="0">
            <div style="display: flex; gap: 8px;">
              <select class="input" style="flex: 1; min-width: 0;" :value="directDns.type" @change="setDnsType(directDns, ($event.target as HTMLSelectElement).value)">
                <option v-for="dt in dnsTypes" :key="dt" :value="dt">{{ dt }}</option>
              </select>
              <template v-if="directDns.type !== 'local'">
                <input class="input mono" style="flex: 1.5; min-width: 0;" v-model="directDns.server" :placeholder="$t('in.addr')" />
                <input class="input mono" style="width: 76px; flex: none;" type="number" min="1" v-model.number="directDns.server_port" :placeholder="$t('in.port')" />
              </template>
            </div>
          </Field>
        </div>
        <div style="margin-top: 14px;">
          <ChipSelect v-model="dnsToDirect" :options="geositeList" :label="$t('setting.toDirectDns')" :placeholder="$t('ui.selectHint')" />
        </div>
      </template>

      <template v-if="enableInb">
        <div class="sub-label">{{ $t('objects.inbound') }}</div>
        <div class="field-grid">
          <Field :label="$t('in.addr') + ' ' + $t('commaSeparated')" :mb="0">
            <input class="input mono" v-model="tunAddress" />
          </Field>
          <Field :label="$t('ui.mtu')" :mb="0">
            <input class="input mono" type="number" min="0" v-model.number="subJsonExt.inbounds[0].mtu" />
          </Field>
          <Field :label="$t('setting.excludePkg') + ' ' + $t('commaSeparated')" :mb="0">
            <input class="input mono" v-model="tunExcludePkg" />
          </Field>
          <div style="display: flex; align-items: flex-end; padding-bottom: 6px;">
            <SwitchLabel v-model="platformProxy" :label="$t('ui.platformProxy')" />
          </div>
        </div>
      </template>

      <div class="builder-foot">
        <div style="flex: 1;" />
        <Btn sm @click="jsonEditor = true"><Ico name="edit" :size="14" /> {{ $t('editor') }}</Btn>
        <Pop :min-width="220" direction="up">
          <template #trigger="{ toggle }">
            <Btn variant="subtle" sm @click="toggle">{{ $t('setting.jsonSubOptions') }} <Ico name="chevronDown" :size="14" /></Btn>
          </template>
          <label class="pop-item" @click.prevent="enableLog = !enableLog"><Toggle :model-value="enableLog" style="pointer-events: none;" /> {{ $t('basic.log.title') }}</label>
          <label class="pop-item" @click.prevent="enableDns = !enableDns"><Toggle :model-value="enableDns" style="pointer-events: none;" /> {{ $t('ui.dnsOpt') }}</label>
          <label class="pop-item" @click.prevent="enableInb = !enableInb"><Toggle :model-value="enableInb" style="pointer-events: none;" /> {{ $t('objects.inbound') }}</label>
          <label class="pop-item" @click.prevent="enableExp = !enableExp"><Toggle :model-value="enableExp" style="pointer-events: none;" /> {{ $t('ui.experimental') }}</label>
        </Pop>
      </div>
    </SettingsGroup>

    <!-- ===================== Clash sub ===================== -->
    <SettingsGroup v-else>
      <div style="font-size: 12.5px; color: var(--text-3); padding: 12px 0 14px;">{{ $t('ui.clashExtDesc') }}</div>

      <div class="field-grid">
        <template v-if="optionMixed">
          <Field :label="$t('setting.mixedPort')" :mb="0">
            <input class="input mono" type="number" min="1" max="65535" v-model.number="mixedPort" />
          </Field>
          <div style="display: flex; align-items: flex-end; padding-bottom: 6px;">
            <SwitchLabel v-model="allowLan" :label="$t('types.ts.allowLanAccess')" />
          </div>
        </template>
        <Field v-if="optionExt" :label="$t('basic.exp.extController')" :mb="0">
          <input class="input mono" v-model="externalController" />
        </Field>
        <Field v-if="optionLog" :label="$t('basic.log.title') + ' - ' + $t('basic.log.level')" :mb="0">
          <select class="input" v-model="clashLogLevel">
            <option v-for="l in clashLevels" :key="l" :value="l">{{ l }}</option>
          </select>
        </Field>
      </div>

      <div v-if="optionTun || optionDns" style="display: flex; gap: 26px; margin-top: 16px; flex-wrap: wrap;">
        <SwitchLabel v-if="optionTun" v-model="clashTun" :label="$t('setting.tun')" />
        <SwitchLabel v-if="optionDns" v-model="clashDns" :label="$t('pages.dns')" />
      </div>

      <div v-if="optionRules" style="margin-top: 16px;">
        <ChipSelect v-model="clashRules" :options="rulesIP" :label="$t('pages.rules')" :placeholder="$t('ui.selectHint')" />
      </div>

      <div class="builder-foot">
        <div style="flex: 1;" />
        <Btn sm @click="clashEditor = true"><Ico name="edit" :size="14" /> {{ $t('editor') }}</Btn>
        <Pop :min-width="220" direction="up">
          <template #trigger="{ toggle }">
            <Btn variant="subtle" sm @click="toggle">{{ $t('setting.jsonSubOptions') }} <Ico name="chevronDown" :size="14" /></Btn>
          </template>
          <label class="pop-item" @click.prevent="optionMixed = !optionMixed"><Toggle :model-value="optionMixed" style="pointer-events: none;" /> {{ $t('setting.mixedPort') }}</label>
          <label class="pop-item" @click.prevent="optionTun = !optionTun"><Toggle :model-value="optionTun" style="pointer-events: none;" /> {{ $t('setting.tun') }}</label>
          <label class="pop-item" @click.prevent="optionExt = !optionExt"><Toggle :model-value="optionExt" style="pointer-events: none;" /> {{ $t('basic.exp.extController') }}</label>
          <label class="pop-item" @click.prevent="optionLog = !optionLog"><Toggle :model-value="optionLog" style="pointer-events: none;" /> {{ $t('basic.log.title') }}</label>
          <label class="pop-item" @click.prevent="optionDns = !optionDns"><Toggle :model-value="optionDns" style="pointer-events: none;" /> {{ $t('pages.dns') }}</label>
          <label class="pop-item" @click.prevent="optionRules = !optionRules"><Toggle :model-value="optionRules" style="pointer-events: none;" /> {{ $t('pages.rules') }}</label>
        </Pop>
      </div>
    </SettingsGroup>
  </div>
</template>

<script lang="ts" setup>
import { computed, onMounted, ref, watch } from 'vue'
import yaml from 'yaml'
import { push } from 'notivue'
import { i18n } from '@/locales'
import HttpUtils from '@/plugins/httputil'
import { FindDiff } from '@/plugins/utils'
import EditorModal from '@/layouts/drawers/EditorModal.vue'
import Tabs from '@/components/ui/Tabs.vue'
import Btn from '@/components/ui/Btn.vue'
import Ico from '@/components/ui/Ico.vue'
import Pop from '@/components/ui/Pop.vue'
import Toggle from '@/components/ui/Toggle.vue'
import SwitchLabel from '@/components/ui/SwitchLabel.vue'
import Field from '@/components/ui/Field.vue'
import ChipSelect from '@/components/ui/ChipSelect.vue'
import SettingsGroup from '@/components/ui/SettingsGroup.vue'
import SRow from '@/components/ui/SRow.vue'
import ToggleRow from '@/components/ui/ToggleRow.vue'

const tab = ref('interface')
const loading = ref(false)
const oldSettings = ref({})
const jsonEditor = ref(false)
const clashEditor = ref(false)

const tabItems = computed((): [string, string][] => [
  ['interface', i18n.global.t('setting.interface')],
  ['sub', i18n.global.t('setting.sub')],
  ['jsonSub', i18n.global.t('setting.jsonSub')],
  ['clashSub', i18n.global.t('setting.clashSub')],
])

const settings = ref({
  webListen: "",
  webDomain: "",
  webPort: "2095",
  webCertFile: "",
  webKeyFile: "",
  webCertMode: "",
  webNginx: "",
  webAcmeEmail: "",
  webPath: "/app/",
  webURI: "",
  sessionMaxAge: "0",
  trafficAge: "30",
  timeLocation: "Asia/Tehran",
  subListen: "",
  subPort: "2096",
  subPath: "/sub/",
  subDomain: "",
  subCertFile: "",
  subKeyFile: "",
  subCertMode: "",
  subAcmeEmail: "",
  subUpdates: "12",
  subEncode: "true",
  subShowInfo: "false",
  subURI: "",
  subJsonExt: "",
  subClashExt: "",
})

onMounted(async () => {
  loading.value = true
  await loadData()
  await detectNginx()
  loading.value = false
})

const loadData = async () => {
  loading.value = true
  const msg = await HttpUtils.get('api/settings')
  loading.value = false
  if (msg.success) {
    setData(msg.obj)
  }
}

const setData = (data: any) => {
  settings.value = data
  loadSubJsonExt()
  oldSettings.value = { ...settings.value }
}

const save = async () => {
  loading.value = true
  // 开了「自动申请证书 (ACME)」时,保存前先真去申请一次校验;失败就拦下不保存,
  // 并弹出真实原因(如 80 端口连不上)。成功的话证书也顺便申请好了。
  const acmeChecks = [
    { on: settings.value.webCertMode === 'acme', domain: settings.value.webDomain, email: settings.value.webAcmeEmail },
    { on: settings.value.subCertMode === 'acme', domain: settings.value.subDomain, email: settings.value.subAcmeEmail },
  ]
  for (const chk of acmeChecks) {
    if (!chk.on) continue
    const r = await HttpUtils.post('api/testAcme', { domain: chk.domain, email: chk.email })
    if (!r.success) {
      loading.value = false
      return
    }
  }
  const msg = await HttpUtils.post('api/save', { object: 'settings', action: 'set', data: JSON.stringify(settings.value) })
  if (msg.success) {
    push.success({
      title: i18n.global.t('success'),
      duration: 5000,
      message: i18n.global.t('actions.set') + " " + i18n.global.t('pages.settings')
    })
    setData(msg.obj.settings)
  }
  loading.value = false
}

const sleep = (ms: number) => new Promise(resolve => setTimeout(resolve, ms))

const restartApp = async () => {
  loading.value = true
  const msg = await HttpUtils.post('api/restartApp', {})
  if (msg.success) {
    let url = settings.value.webURI
    if (url !== "") {
      const isTLS = settings.value.webCertMode === "acme" || settings.value.webCertFile !== "" || settings.value.webKeyFile !== ""
      url = buildURL(settings.value.webDomain, settings.value.webPort.toString(), isTLS, settings.value.webPath)
    }
    await sleep(3000)
    window.location.replace(url)
  }
  loading.value = false
}

const buildURL = (host: string, port: string, isTLS: boolean, path: string) => {
  if (!host || host.length == 0) host = window.location.hostname
  if (!port || port.length == 0) port = window.location.port

  const protocol = isTLS ? "https:" : "http:"

  if (port === "" || (isTLS && port === "443") || (!isTLS && port === "80")) {
    port = ""
  } else {
    port = `:${port}`
  }

  return `${protocol}//${host}${port}${path}settings`
}

const detectedNginx = ref({ installed: false, active: false })

const webNginxHint = computed(() => {
  return detectedNginx.value.installed
    ? i18n.global.t('setting.nginxDetected')
    : i18n.global.t('setting.deployNginxHint')
})

// 页面加载时检测 nginx;仅当 webNginx 尚未设置过(空)时用检测结果作默认值
const detectNginx = async () => {
  const r = await HttpUtils.get('api/detectNginx')
  if (r.success && r.obj) {
    detectedNginx.value = r.obj
    if (!settings.value.webNginx) {
      settings.value.webNginx = r.obj.installed ? 'true' : 'false'
      oldSettings.value = { ...settings.value }
    }
  }
}

// 用 acme.sh 申请证书:无 nginx 走 standalone、面板自身加载证书;有 nginx 走 nginx 模式、证书给 nginx
const issueCert = async () => {
  if (!settings.value.webDomain) return
  loading.value = true
  const r = await HttpUtils.post('api/issueCert', {
    domain: settings.value.webDomain,
    email: settings.value.webAcmeEmail,
    nginx: settings.value.webNginx === 'true' ? 'true' : 'false',
  })
  if (r.success && r.obj) {
    if (settings.value.webNginx !== 'true') {
      settings.value.webCertFile = r.obj.certFile
      settings.value.webKeyFile = r.obj.keyFile
      settings.value.webCertMode = ''
    }
    push.success({
      title: i18n.global.t('success'),
      duration: 5000,
      message: i18n.global.t('setting.issueCertOk'),
    })
  }
  loading.value = false
}

const webAcme = computed({
  get: () => { return settings.value.webCertMode == "acme" },
  set: (v: boolean) => { settings.value.webCertMode = v ? "acme" : "" }
})

const subAcme = computed({
  get: () => { return settings.value.subCertMode == "acme" },
  set: (v: boolean) => { settings.value.subCertMode = v ? "acme" : "" }
})

const subEncode = computed({
  get: () => { return settings.value.subEncode == "true" },
  set: (v: boolean) => { settings.value.subEncode = v ? "true" : "false" }
})

const subShowInfo = computed({
  get: () => { return settings.value.subShowInfo == "true" },
  set: (v: boolean) => { settings.value.subShowInfo = v ? "true" : "false" }
})

const webPort = computed({
  get: () => { return settings.value.webPort.length > 0 ? parseInt(settings.value.webPort) : 2095 },
  set: (v: number) => { settings.value.webPort = v > 0 ? v.toString() : "2095" }
})

const sessionMaxAge = computed({
  get: () => { return settings.value.sessionMaxAge.length > 0 ? parseInt(settings.value.sessionMaxAge) : 0 },
  set: (v: number) => { settings.value.sessionMaxAge = v > 0 ? v.toString() : "0" }
})

const trafficAge = computed({
  get: () => { return settings.value.trafficAge.length > 0 ? parseInt(settings.value.trafficAge) : 0 },
  set: (v: number) => { settings.value.trafficAge = v > 0 ? v.toString() : "0" }
})

const subPort = computed({
  get: () => { return settings.value.subPort.length > 0 ? parseInt(settings.value.subPort) : 2096 },
  set: (v: number) => { settings.value.subPort = v > 0 ? v.toString() : "2096" }
})

const subUpdates = computed({
  get: () => { return settings.value.subUpdates.length > 0 ? parseInt(settings.value.subUpdates) : 12 },
  set: (v: number) => { settings.value.subUpdates = v > 0 ? v.toString() : "12" }
})

const stateChange = computed(() => {
  return !FindDiff.deepCompare(settings.value, oldSettings.value)
})

/* ===================================================================
 * JSON subscription extension (逻辑平移自旧 components/SubJsonExt.vue)
 * =================================================================== */

const subJsonExt = ref<any>({})

const levels = ["trace", "debug", "info", "warn", "error", "fatal", "panic"]
const dnsTypes = ['udp', 'tcp', 'local', 'tls', 'quic', 'h3']

const defaultLog = {
  "level": "info",
  "timestamp": true
}

const defaultInb = [
  {
    "type": "tun",
    "address": [
      "172.19.0.1/30",
      "fdfe:dcba:9876::1/126"
    ],
    "mtu": 9000,
    "auto_route": true,
    "strict_route": false,
    "endpoint_independent_nat": false,
    "stack": "mixed",
    "exclude_package": <string[]>[],
    "platform": {
      "http_proxy": {
        "enabled": true,
        "server": "127.0.0.1",
        "server_port": 2080
      }
    }
  },
  {
    "type": "mixed",
    "listen": "127.0.0.1",
    "listen_port": 2080,
    "users": []
  }
]

const defaultExp = {
  "clash_api": {
    "external_controller": "127.0.0.1:9090",
    "external_ui": "ui",
    "secret": "",
    "external_ui_download_url": "https://mirror.ghproxy.com/https://github.com/MetaCubeX/Yacd-meta/archive/gh-pages.zip",
    "external_ui_download_detour": "direct",
    "default_mode": "rule"
  },
  "cache_file": {
    "enabled": true,
    "store_fakeip": false
  }
}

const defaultDns = {
  "servers": [
    {
      "type": "tcp",
      "tag": "proxy-dns",
      "server": "8.8.8.8",
      "server_port": 53,
      "detour": "proxy",
      "domain_resolver": "local-dns",
    },
    {
      "tag": "direct-dns",
      "type": "local",
    },
    {
      "tag": "local-dns",
      "type": "local",
    }
  ],
  "rules": [
    {
      "clash_mode": "Global",
      "source_ip_cidr": [
        "172.19.0.0/30",
        "fdfe:dcba:9876::1/126"
      ],
      "action": "route",
      "server": "proxy-dns"
    },
    {
      "clash_mode": "Direct",
      "action": "route",
      "server": "direct-dns"
    },
    {
      "source_ip_cidr": [
        "172.19.0.0/30",
        "fdfe:dcba:9876::1/126"
      ],
      "action": "route",
      "server": "proxy-dns"
    },
  ],
  "final": "local-dns",
  "strategy": "prefer_ipv4"
}

const geositeList = [
  { title: "Private", value: "geosite-private" },
  { title: "Ads", value: "geosite-ads" },
  { title: "🇮🇷 Iran", value: "geosite-ir" },
  { title: "🇨🇳 China", value: "geosite-cn" },
  { title: "🇻🇳 Vietnam", value: "geosite-vn" },
]

const geoList = [
  { title: "Site-Private", value: "geoip-private" },
  { title: "IP-Private", value: "geosite-private" },
  { title: "Site-Ads", value: "geosite-ads" },
  { title: "🇮🇷 Site-Iran", value: "geosite-ir" },
  { title: "🇮🇷 IP-Iran", value: "geoip-ir" },
  { title: "🇨🇳 Site-China", value: "geosite-cn" },
  { title: "🇨🇳 IP-China", value: "geoip-cn" },
  { title: "🇻🇳 Site-Vietnam", value: "geosite-vn" },
  { title: "🇻🇳 IP-Vietnam", value: "geoip-vn" },
]

const geo = [
  {
    tag: "geosite-ads",
    type: "remote",
    format: "binary",
    url: "https://testingcf.jsdelivr.net/gh/MetaCubeX/meta-rules-dat@sing/geo/geosite/category-ads-all.srs",
    download_detour: "direct"
  },
  {
    tag: "geosite-private",
    type: "remote",
    format: "binary",
    url: "https://testingcf.jsdelivr.net/gh/MetaCubeX/meta-rules-dat@sing/geo/geosite/private.srs",
    download_detour: "direct"
  },
  {
    tag: "geosite-ir",
    type: "remote",
    format: "binary",
    url: "https://testingcf.jsdelivr.net/gh/MetaCubeX/meta-rules-dat@sing/geo/geosite/category-ir.srs",
    download_detour: "direct"
  },
  {
    tag: "geosite-cn",
    type: "remote",
    format: "binary",
    url: "https://testingcf.jsdelivr.net/gh/MetaCubeX/meta-rules-dat@sing/geo/geosite/cn.srs",
    download_detour: "direct"
  },
  {
    tag: "geosite-vn",
    type: "remote",
    format: "binary",
    url: "https://github.com/Thaomtam/Geosite-vn/raw/rule-set/Geosite-vn.srs",
    download_detour: "direct"
  },
  {
    tag: "geoip-private",
    type: "remote",
    format: "binary",
    url: "https://testingcf.jsdelivr.net/gh/MetaCubeX/meta-rules-dat@sing/geo/geoip/private.srs",
    download_detour: "direct"
  },
  {
    tag: "geoip-ir",
    type: "remote",
    format: "binary",
    url: "https://testingcf.jsdelivr.net/gh/MetaCubeX/meta-rules-dat@sing/geo/geoip/ir.srs",
    download_detour: "direct"
  },
  {
    tag: "geoip-cn",
    type: "remote",
    format: "binary",
    url: "https://testingcf.jsdelivr.net/gh/MetaCubeX/meta-rules-dat@sing/geo/geoip/cn.srs",
    download_detour: "direct"
  },
  {
    tag: "geoip-vn",
    type: "remote",
    format: "binary",
    url: "https://testingcf.jsdelivr.net/gh/MetaCubeX/meta-rules-dat@sing/geo/geoip/vn.srs",
    download_detour: "direct"
  }
]

// 解析 settings.subJsonExt → 本地对象;并按旧实现把规范化后的 JSON 回写
const loadSubJsonExt = () => {
  const s: string = settings.value.subJsonExt ?? ''
  if (s.length > 0) {
    try {
      const parsed = JSON.parse(s)
      subJsonExt.value = parsed
      settings.value.subJsonExt = Object.keys(parsed).length > 0 ? JSON.stringify(parsed, null, 2) : ''
    } catch (e) {
      // 数据库里是坏 JSON:保留原始字符串不动(与旧版一致,不在这里破坏数据)
    }
  } else {
    subJsonExt.value = {}
  }
}

watch(subJsonExt, (v) => {
  settings.value.subJsonExt = Object.keys(v).length > 0 ? JSON.stringify(v, null, 2) : ''
}, { deep: true })

const enableLog = computed({
  get: (): boolean => subJsonExt.value?.log != undefined,
  set: (v: boolean) => { v ? subJsonExt.value.log = defaultLog : delete subJsonExt.value.log }
})

const enableDns = computed({
  get: (): boolean => subJsonExt.value?.dns != undefined,
  set: (v: boolean) => {
    if (v) {
      subJsonExt.value.dns = defaultDns
      if (rules.value == undefined) subJsonExt.value.rules = [{ action: 'sniff' }]
      subJsonExt.value.rules.unshift({ protocol: "dns", action: "hijack-dns" })
    } else {
      delete subJsonExt.value.dns
      const newRules = subJsonExt.value?.rules?.filter((r: any) => r.protocol != "dns") ?? []
      if (newRules.length >= 0) subJsonExt.value.rules = newRules
      if (rules.value?.length == 0) delete subJsonExt.value.rules
    }
  }
})

const enableInb = computed({
  get: (): boolean => subJsonExt.value?.inbounds != undefined,
  set: (v: boolean) => { v ? subJsonExt.value.inbounds = defaultInb.slice() : delete subJsonExt.value.inbounds }
})

const enableExp = computed({
  get: (): boolean => subJsonExt.value?.experimental != undefined,
  set: (v: boolean) => { v ? subJsonExt.value.experimental = defaultExp : delete subJsonExt.value.experimental }
})

const dns = computed((): any => subJsonExt.value?.dns ?? undefined)

const proxyDns = computed((): any => dns.value?.servers?.findLast((d: any) => d.tag == "proxy-dns") ?? {})

const directDns = computed((): any => dns.value?.servers?.findLast((d: any) => d.tag == "direct-dns") ?? {})

// 平移旧 SimpleDNS:切到 local 时去掉 server / server_port
const setDnsType = (server: any, t: string) => {
  server.type = t
  if (t == 'local') {
    delete server.server
    delete server.server_port
  }
}

const dnsTags = computed((): string[] => dns.value?.servers?.map((d: any) => d.tag) ?? [])

const defaultResolver = computed({
  get: (): string => subJsonExt.value?.default_domain_resolver ?? '',
  set: (v: string) => {
    if (v) subJsonExt.value.default_domain_resolver = v
    else delete subJsonExt.value.default_domain_resolver
  }
})

const dnsToDirect = computed({
  get: (): string[] => {
    const ruleIndex = dns.value?.rules?.findIndex((r: any) => r.server == "direct-dns" && Object.hasOwn(r, 'rule_set')) ?? -1
    return ruleIndex >= 0 ? dns.value.rules[ruleIndex].rule_set : []
  },
  set: (v: string[]) => {
    const ruleIndex = dns.value?.rules?.findIndex((r: any) => r.server == "direct-dns" && Object.hasOwn(r, 'rule_set')) ?? -1
    if (v.length > 0) {
      if (ruleIndex >= 0) {
        dns.value.rules[ruleIndex].rule_set = v
      } else {
        dns.value.rules.push({ rule_set: v, action: "route", server: "direct-dns" })
      }
    } else {
      if (ruleIndex != -1) dns.value.rules.splice(ruleIndex, 1)
    }
    updateRuleSets()
  }
})

const inbounds = computed((): any => subJsonExt.value?.inbounds ?? undefined)

const tunAddress = computed({
  get: (): string => (inbounds.value?.[0]?.address ?? []).join(','),
  set: (v: string) => {
    if (!inbounds.value?.[0]) return
    inbounds.value[0].address = v.split(',').map((s: string) => s.trim()).filter((s: string) => s.length > 0)
  }
})

const tunExcludePkg = computed({
  get: (): string => (inbounds.value?.[0]?.exclude_package ?? []).join(','),
  set: (v: string) => {
    if (!inbounds.value?.[0]) return
    inbounds.value[0].exclude_package = v.split(',').map((s: string) => s.trim()).filter((s: string) => s.length > 0)
  }
})

const platformProxy = computed({
  get: (): boolean => inbounds.value?.[0]?.platform != undefined,
  set: (v: boolean) => { subJsonExt.value.inbounds[0].platform = v ? defaultInb[0].platform : undefined }
})

const rules = computed((): any => subJsonExt.value?.rules ?? undefined)

const ruleToDirect = computed({
  get: (): string[] => {
    const ruleIndex = rules.value?.findIndex((r: any) => r.outbound == "direct" && Object.hasOwn(r, 'rule_set')) ?? -1
    return ruleIndex >= 0 ? rules.value[ruleIndex].rule_set : []
  },
  set: (v: string[]) => {
    const ruleIndex = rules.value?.findIndex((r: any) => r.outbound == "direct" && Object.hasOwn(r, 'rule_set')) ?? -1
    if (v.length > 0) {
      if (ruleIndex >= 0) {
        rules.value[ruleIndex].rule_set = v
      } else {
        if (rules.value == undefined) subJsonExt.value.rules = []
        rules.value.push({ rule_set: v, action: "route", outbound: "direct" })
      }
    } else {
      if (ruleIndex != -1) rules.value.splice(ruleIndex, 1)
    }
    updateRuleSets()
  }
})

const ruleToBlock = computed({
  get: (): string[] => {
    const ruleIndex = rules.value?.findIndex((r: any) => r.action == "reject" && Object.hasOwn(r, 'rule_set')) ?? -1
    return ruleIndex >= 0 ? rules.value[ruleIndex].rule_set : []
  },
  set: (v: string[]) => {
    const ruleIndex = rules.value?.findIndex((r: any) => r.action == "reject" && Object.hasOwn(r, 'rule_set')) ?? -1
    if (v.length > 0) {
      if (ruleIndex >= 0) {
        rules.value[ruleIndex].rule_set = v
      } else {
        if (rules.value == undefined) subJsonExt.value.rules = []
        rules.value.push({ rule_set: v, action: "reject" })
      }
    } else {
      if (ruleIndex != -1) rules.value.splice(ruleIndex, 1)
    }
    updateRuleSets()
  }
})

const updateRuleSets = () => {
  let tags = <string[]>[]
  if ((dns.value?.rules?.length ?? 0) > 0) dns.value.rules.forEach((r: any) => { if (r.rule_set) tags.push(...r.rule_set) })
  if ((rules.value?.length ?? 0) > 0) rules.value.forEach((r: any) => { if (r.rule_set) tags.push(...r.rule_set) })
  if (tags.length > 0) {
    subJsonExt.value.rule_set = geo.filter((g: any) => tags.includes(g.tag))
  } else {
    delete subJsonExt.value.rule_set
  }
  if (rules.value?.length == 0) delete subJsonExt.value.rules
}

const saveJsonEditor = (data: string) => {
  try {
    subJsonExt.value = JSON.parse(data)
  } catch (e) {
    push.error({
      message: i18n.global.t('failed') + ": " + i18n.global.t('error.invalidData'),
      duration: 5000,
    })
  }
}

/* ===================================================================
 * Clash subscription extension (逻辑平移自旧 components/SubClashExt.vue)
 * =================================================================== */

const defaultConfig: any = {
  "mixed-port": 7890,
  "allow-lan": false,
  "mode": "rule",
  "log-level": "info",
  "external-controller": "127.0.0.1:9090",
  "tun": {
    "enable": true,
    "stack": "system",
    "auto-route": true,
    "auto-detect-interface": true,
    "dns-hijack": ["any:53"],
  },
  "dns": {
    "enable": true,
    "ipv6": false,
    "enhanced-mode": "fake-ip",
    "fake-ip-range": "198.18.0.1/16",
    "default-nameserver": ["8.8.8.8", "1.1.1.1"],
    "nameserver": [
      "https://doh.pub/dns-query",
      "https://1.0.0.1/dns-query"
    ],
    "fallback": ["tcp://9.9.9.9:53"],
    "fake-ip-filter": ["*.lan", "localhost", "*.local"]
  },
  "rules": [
    "GEOIP,Private,DIRECT",
    "MATCH,Proxy"
  ]
}

const clashLevels = ['debug', 'info', 'warning', 'error']

const rulesIP = [
  { title: 'Private-Direct', value: 'GEOIP,Private,DIRECT' },
  { title: 'Private-Block', value: 'GEOIP,Private,REJECT' },
  { title: 'LAN-Direct', value: 'GEOIP,LAN,DIRECT' },
  { title: 'LAN-Block', value: 'GEOIP,LAN,REJECT' },
  { title: 'Ads-Direct', value: 'GEOIP,Ads,DIRECT' },
  { title: 'Ads-Block', value: 'GEOIP,Ads,REJECT' },
  { title: '🇨🇳 China-Direct', value: 'GEOIP,CN,DIRECT' },
  { title: '🇨🇳 China-Block', value: 'GEOIP,CN,REJECT' },
  { title: '🇮🇷 Iran-Direct', value: 'GEOIP,CATEGORY-IR,DIRECT' },
  { title: '🇮🇷 Iran-Block', value: 'GEOIP,CATEGORY-IR,REJECT' },
  { title: '🇻🇳 Vietnam-Direct', value: 'GEOIP,CATEGORY-VN,DIRECT' },
  { title: '🇻🇳 Vietnam-Block', value: 'GEOIP,CATEGORY-VN,REJECT' },
  { title: '🇯🇵 Japan-Direct', value: 'GEOIP,JP,DIRECT' },
  { title: '🇯🇵 Japan-Block', value: 'GEOIP,JP,REJECT' },
]

const metaJson = computed({
  get: (): any => {
    try {
      return yaml.parse(settings.value.subClashExt) ?? {}
    } catch (e) {
      return {}
    }
  },
  set: (v: any) => {
    settings.value.subClashExt = Object.keys(v).length == 0 ? "" : yaml.stringify(v)
  }
})

const updateMetaJson = (data: any, key: string) => {
  let newMetaJson = metaJson.value
  if (data == null) {
    delete newMetaJson[key]
  } else {
    newMetaJson[key] = data
  }
  metaJson.value = newMetaJson
}

const optionMixed = computed({
  get: (): boolean => (metaJson.value['mixed-port'] ?? 0) > 0,
  set: (v: boolean) => {
    updateMetaJson(v ? defaultConfig['mixed-port'] : null, 'mixed-port')
    updateMetaJson(v ? defaultConfig['allow-lan'] : null, 'allow-lan')
  }
})

const optionTun = computed({
  get: (): boolean => metaJson.value['tun']?.['enable'] ?? false,
  set: (v: boolean) => { updateMetaJson(v ? defaultConfig['tun'] : null, 'tun') }
})

const optionExt = computed({
  get: (): boolean => (metaJson.value['external-controller']?.length ?? 0) > 0,
  set: (v: boolean) => { updateMetaJson(v ? defaultConfig['external-controller'] : null, 'external-controller') }
})

const optionLog = computed({
  get: (): boolean => (metaJson.value['log-level']?.length ?? 0) > 0,
  set: (v: boolean) => { updateMetaJson(v ? defaultConfig['log-level'] : null, 'log-level') }
})

const optionDns = computed({
  get: (): boolean => metaJson.value['dns']?.['enable'] ?? false,
  set: (v: boolean) => { updateMetaJson(v ? defaultConfig['dns'] : null, 'dns') }
})

const optionRules = computed({
  get: (): boolean => (metaJson.value['rules']?.length ?? 0) > 0,
  set: (v: boolean) => {
    updateMetaJson(v ? defaultConfig['rules'] : null, 'rules')
    updateMetaJson(v ? defaultConfig['mode'] : null, 'mode')
  }
})

const mixedPort = computed({
  get: (): number => metaJson.value['mixed-port'],
  set: (v: number) => { updateMetaJson(v, 'mixed-port') }
})

const allowLan = computed({
  get: (): boolean => metaJson.value['allow-lan'] ?? false,
  set: (v: boolean) => { updateMetaJson(v, 'allow-lan') }
})

const externalController = computed({
  get: (): string => metaJson.value['external-controller'] ?? '',
  set: (v: string) => { updateMetaJson(v, 'external-controller') }
})

const clashLogLevel = computed({
  get: (): string => metaJson.value['log-level'] ?? '',
  set: (v: string) => { updateMetaJson(v, 'log-level') }
})

const clashDns = computed({
  get: (): boolean => metaJson.value['dns']?.['enable'] ?? false,
  set: (v: boolean) => { updateMetaJson({ ...metaJson.value['dns'], 'enable': v }, 'dns') }
})

const clashTun = computed({
  get: (): boolean => metaJson.value['tun']?.['enable'] ?? false,
  set: (v: boolean) => { updateMetaJson({ ...metaJson.value['tun'], 'enable': v }, 'tun') }
})

const clashRules = computed({
  get: (): string[] => (metaJson.value.rules?.length ?? 0) > 0 ? metaJson.value.rules.filter((r: string) => r != "MATCH,Proxy") : [],
  set: (v: string[]) => {
    let newRules = <string[]>[]
    v.forEach((r: string) => { newRules.push(r) })
    updateMetaJson([...newRules, "MATCH,Proxy"], 'rules')
  }
})

const saveClashEditor = (data: string) => {
  try {
    const result = yaml.parse(data)
    if (typeof result != 'object' || Array.isArray(result)) throw new Error()
  } catch (e) {
    push.error({
      message: i18n.global.t('failed') + ": " + i18n.global.t('error.invalidData'),
      duration: 5000,
    })
    return
  }
  settings.value.subClashExt = data
}
</script>

<style scoped>
.head-row {
  display: flex; align-items: center; gap: 8px;
  border-bottom: 1px solid var(--line);
  flex-wrap: wrap;
}
.head-actions {
  display: flex; gap: 8px; align-items: center;
  margin-inline-start: auto;
}
@media (max-width: 820px) {
  .head-actions { margin-inline-start: 0; padding-bottom: 10px; }
}
.srow-btn {
  display: flex; align-items: center;
  padding: 13px 0; border-bottom: 1px solid var(--line);
}
.sub-label {
  font-size: 12px; font-weight: 700; color: var(--text-3);
  letter-spacing: .04em; text-transform: uppercase;
  margin: 16px 0 10px;
}
.builder-foot {
  display: flex; align-items: center; gap: 10px;
  margin-top: 18px; padding-top: 14px;
  border-top: 1px solid var(--line);
}
</style>
