<template>
  <v-card :loading="loading">
    <v-tabs
    v-model="tab"
    color="primary"
    align-tabs="center"
    show-arrows
  >
    <v-tab value="t1">{{ $t('setting.interface') }}</v-tab>
    <v-tab value="t2">{{ $t('setting.sub') }}</v-tab>
    <v-tab value="t3">{{ $t('setting.jsonSub') }}</v-tab>
    <v-tab value="t4">{{ $t('setting.clashSub') }}</v-tab>
  </v-tabs>
  <v-card-text>
    <v-row align="center" justify="center" style="margin-bottom: 10px;">
      <v-col cols="auto">
        <v-btn color="primary" @click="save" :loading="loading" :disabled="!stateChange">
          {{ $t('actions.save') }}
        </v-btn>
      </v-col>
      <v-col cols="auto">
        <v-btn variant="outlined" color="warning" @click="restartApp" :loading="loading" :disabled="stateChange">
          {{ $t('actions.restartApp') }}
        </v-btn>
      </v-col>
    </v-row>
    <v-window v-model="tab">
      <v-window-item value="t1">
        <v-row>
          <v-col cols="12" sm="6" md="4">
            <v-text-field v-model="settings.webListen" :label="$t('setting.addr')" hide-details></v-text-field>
          </v-col>
          <v-col cols="12" sm="6" md="4">
            <v-text-field v-model.number="webPort" min="1" type="number" :label="$t('setting.port')" hide-details></v-text-field>
          </v-col>
          <v-col cols="12" sm="6" md="4">
            <v-text-field v-model="settings.webPath" :label="$t('setting.webPath')" hide-details></v-text-field>
          </v-col>
          <v-col cols="12" sm="6" md="4">
            <v-text-field v-model="settings.webDomain" :label="$t('setting.domain')" hide-details></v-text-field>
          </v-col>
          <v-col cols="12" sm="6" md="4">
            <v-select
              v-model="settings.webNginx"
              :items="nginxItems"
              :label="$t('setting.deployNginx')"
              :hint="webNginxHint"
              persistent-hint
            ></v-select>
          </v-col>
          <v-col cols="12" sm="6" md="4">
            <v-text-field v-model="settings.webAcmeEmail" :label="$t('setting.acmeEmail')" :hint="$t('setting.acmeHint')" persistent-hint></v-text-field>
          </v-col>
          <v-col cols="12" sm="6" md="4" class="d-flex align-center">
            <v-btn color="success" variant="flat" @click="issueCert" :loading="loading" :disabled="!settings.webDomain">
              {{ $t('setting.issueCert') }}
            </v-btn>
          </v-col>
          <!-- 非 nginx 模式:面板自身加载证书,可开内置 ACME 或填证书路径(acme.sh 申请后自动回填) -->
          <template v-if="settings.webNginx !== 'true'">
            <v-col cols="12" sm="6" md="4">
              <v-switch color="primary" v-model="webAcme" :label="$t('setting.autoCert')" hide-details />
            </v-col>
            <template v-if="!webAcme">
              <v-col cols="12" sm="6" md="4">
                <v-text-field v-model="settings.webKeyFile" :label="$t('setting.sslKey')" hide-details></v-text-field>
              </v-col>
              <v-col cols="12" sm="6" md="4">
                <v-text-field v-model="settings.webCertFile" :label="$t('setting.sslCert')" hide-details></v-text-field>
              </v-col>
            </template>
          </template>
          <v-col cols="12" sm="6" md="4">
            <v-text-field v-model="settings.webURI" :label="$t('setting.webUri')" hide-details></v-text-field>
          </v-col>
          <v-col cols="12" sm="6" md="4">
            <v-text-field
              type="number"
              v-model.number="sessionMaxAge"
              min="0"
              :label="$t('setting.sessionAge')"
              :suffix="$t('date.m')"
              hide-details
              ></v-text-field>
          </v-col>
          <v-col cols="12" sm="6" md="4">
            <v-text-field
              type="number"
              v-model.number="trafficAge"
              min="0"
              :label="$t('setting.trafficAge')"
              :suffix="$t('date.d')"
              hide-details
              ></v-text-field>
          </v-col>
          <v-col cols="12" sm="6" md="4">
            <v-text-field v-model="settings.timeLocation" :label="$t('setting.timeLoc')" hide-details></v-text-field>
          </v-col>
        </v-row>
      </v-window-item>

      <v-window-item value="t2">
        <v-row>
          <v-col cols="12" sm="6" md="4">
            <v-switch color="primary" v-model="subEncode" :label="$t('setting.subEncode')" hide-details />
          </v-col>
          <v-col cols="12" sm="6" md="4">
            <v-switch color="primary" v-model="subShowInfo" :label="$t('setting.subInfo')" hide-details />
          </v-col>
        </v-row>
        <v-row>
          <v-col cols="12" sm="6" md="4">
            <v-text-field v-model="settings.subListen" :label="$t('setting.addr')" hide-details></v-text-field>
          </v-col>
          <v-col cols="12" sm="6" md="4">
            <v-text-field
              type="number"
              v-model.number="subPort"
              min="1"
              :label="$t('setting.port')"
              hide-details></v-text-field>
          </v-col>
        </v-row>
        <v-row>
          <v-col cols="12" sm="6" md="4">
            <v-switch color="primary" v-model="subAcme" :label="$t('setting.autoCert')" hide-details />
          </v-col>
          <v-col v-if="subAcme" cols="12" sm="6" md="4">
            <v-text-field v-model="settings.subAcmeEmail" :label="$t('setting.acmeEmail')" :hint="$t('setting.acmeHint')" persistent-hint></v-text-field>
          </v-col>
          <template v-else>
            <v-col cols="12" sm="6" md="4">
              <v-text-field v-model="settings.subKeyFile" :label="$t('setting.sslKey')" hide-details></v-text-field>
            </v-col>
            <v-col cols="12" sm="6" md="4">
              <v-text-field v-model="settings.subCertFile" :label="$t('setting.sslCert')" hide-details></v-text-field>
            </v-col>
          </template>
        </v-row>
        <v-row>
          <v-col cols="12" sm="6" md="4">
            <v-text-field v-model="settings.subDomain" :label="$t('setting.domain')" hide-details></v-text-field>
          </v-col>
          <v-col cols="12" sm="6" md="4">
            <v-text-field v-model="settings.subPath" :label="$t('setting.path')" hide-details></v-text-field>
          </v-col>
        </v-row>
        <v-row>
          <v-col cols="12" sm="6" md="4">
            <v-text-field
              type="number"
              v-model.number="subUpdates"
              min="0"
              :label="$t('setting.update')"
              hide-details
              ></v-text-field>
          </v-col>
          <v-col cols="12" sm="6" md="4">
            <v-text-field v-model="settings.subURI" :label="$t('setting.subUri')" hide-details></v-text-field>
          </v-col>
        </v-row>
      </v-window-item>

      <v-window-item value="t3">
        <SubJsonExtVue :settings="settings" />
      </v-window-item>

      <v-window-item value="t4">
        <SubClashExtVue :settings="settings" />
      </v-window-item>
    </v-window>
  </v-card-text>
</v-card>
</template>

<script lang="ts" setup>
import { i18n } from '@/locales'
import { Ref, computed, inject, onMounted, ref } from 'vue'
import HttpUtils from '@/plugins/httputil'
import { FindDiff } from '@/plugins/utils'
import SubJsonExtVue from '@/components/SubJsonExt.vue'
import SubClashExtVue from '@/components/SubClashExt.vue'
import { push } from 'notivue'
const tab = ref("t1")
const loading:Ref = inject('loading')?? ref(false)
const oldSettings = ref({})

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
  oldSettings.value = { ...data }
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
  const msg = await HttpUtils.post('api/restartApp',{})
  if (msg.success) {
    let url = settings.value.webURI
    if (url !== "") {
      const isTLS = settings.value.webCertMode === "acme" || settings.value.webCertFile !== "" || settings.value.webKeyFile !== ""
      url = buildURL(settings.value.webDomain,settings.value.webPort.toString(),isTLS, settings.value.webPath)
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

const nginxItems = computed(() => [
  { title: i18n.global.t('no'), value: 'false' },
  { title: i18n.global.t('yes'), value: 'true' },
])

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
  set: (v:boolean) => { settings.value.webCertMode = v ? "acme" : "" }
})

const subAcme = computed({
  get: () => { return settings.value.subCertMode == "acme" },
  set: (v:boolean) => { settings.value.subCertMode = v ? "acme" : "" }
})

const subEncode = computed({
  get: () => { return settings.value.subEncode == "true" },
  set: (v:boolean) => { settings.value.subEncode = v ? "true" : "false" }
})

const subShowInfo = computed({
  get: () => { return settings.value.subShowInfo == "true" },
  set: (v:boolean) => { settings.value.subShowInfo = v ? "true" : "false" }
})

const webPort = computed({
  get: () => { return settings.value.webPort.length>0 ? parseInt(settings.value.webPort) : 2095 },
  set: (v:number) => { settings.value.webPort = v>0 ? v.toString() : "2095" }
})

const sessionMaxAge = computed({
  get: () => { return settings.value.sessionMaxAge.length>0 ? parseInt(settings.value.sessionMaxAge) : 0 },
  set: (v:number) => { settings.value.sessionMaxAge = v>0 ? v.toString() : "0" }
})

const trafficAge = computed({
  get: () => { return settings.value.trafficAge.length>0 ? parseInt(settings.value.trafficAge) : 0 },
  set: (v:number) => { settings.value.trafficAge = v>0 ? v.toString() : "0" }
})

const subPort = computed({
  get: () => { return settings.value.subPort.length>0 ? parseInt(settings.value.subPort) : 2096 },
  set: (v:number) => { settings.value.subPort = v>0 ? v.toString() : "2096" }
})

const subUpdates = computed({
  get: () => { return settings.value.subUpdates.length>0 ? parseInt(settings.value.subUpdates) : 12 },
  set: (v:number) => { settings.value.subUpdates = v>0 ? v.toString() : "12" }
})

const stateChange = computed(() => {
  return !FindDiff.deepCompare(settings.value,oldSettings.value)
})
</script>
