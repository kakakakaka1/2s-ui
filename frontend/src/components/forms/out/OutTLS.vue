<template>
  <div style="margin-bottom: 15px;">
    <MSwitchRow v-if="tlsOptional" v-model="tlsEnable" :label="$t('ui.tlsTitle')" :desc="$t('ui.tlsDesc')" />
    <div v-if="tls.enabled" class="fade-up">
      <div v-if="!tlsOptional" style="display: flex; align-items: center; gap: 10px; margin-bottom: 12px;">
        <SectionLabel>{{ $t('objects.tls') }}</SectionLabel>
        <div style="flex: 1;" />
      </div>

      <div style="display: flex; gap: 20px; flex-wrap: wrap; margin-bottom: 14px;">
        <SwitchLabel v-model="disableSni" :label="$t('tls.disableSni')" />
        <SwitchLabel v-model="insecure" :label="$t('tls.insecure')" />
      </div>

      <template v-if="optionCert">
        <div style="margin-bottom: 12px;">
          <Segmented v-model="usePath" :options="[[0, $t('tls.usePath')], [1, $t('tls.useText')]]" block
            @update:model-value="switchCertMode" />
        </div>
        <Field v-if="usePath == 0" :label="$t('tls.certPath')">
          <input class="input mono" v-model="tls.certificate_path" />
        </Field>
        <Field v-else :label="$t('tls.cert')">
          <textarea class="input mono" rows="4" spellcheck="false" v-model="tls.certificate"></textarea>
        </Field>
      </template>

      <div class="grid2">
        <Field v-if="tls.server_name != undefined" label="SNI">
          <input class="input mono" v-model="tls.server_name" />
        </Field>
        <Field v-if="tls.min_version" :label="$t('tls.minVer')">
          <select class="input" v-model="tls.min_version">
            <option v-for="v in tlsVersions" :key="v" :value="v">{{ v }}</option>
          </select>
        </Field>
        <Field v-if="tls.max_version" :label="$t('tls.maxVer')">
          <select class="input" v-model="tls.max_version">
            <option v-for="v in tlsVersions" :key="v" :value="v">{{ v }}</option>
          </select>
        </Field>
        <Field v-if="tls.utls != undefined" label="Fingerprint">
          <select class="input" v-model="tls.utls.fingerprint">
            <option v-for="f in fingerprints" :key="f.value" :value="f.value">{{ f.title }}</option>
          </select>
        </Field>
      </div>

      <Field v-if="tls.alpn" label="ALPN">
        <MultiPick v-model="alpnValue" :options="alpnOptions" />
      </Field>
      <Field v-if="tls.cipher_suites != undefined" :label="$t('tls.cs')">
        <ChipSelect v-model="cipherSuites" :options="cipherSuiteItems" />
      </Field>

      <template v-if="tls.reality != undefined">
        <div class="grid2">
          <Field :label="$t('tls.pubKey')">
            <input class="input mono" v-model="tls.reality.public_key" />
          </Field>
          <Field label="Short ID">
            <input class="input mono" v-model="tls.reality.short_id" />
          </Field>
        </div>
      </template>

      <template v-if="tls.ech != undefined">
        <div style="display: flex; align-items: center; gap: 10px; margin: 4px 0 12px;">
          <SectionLabel>ECH</SectionLabel>
          <div style="flex: 1;" />
          <Segmented v-model="useEchPath" :options="[[0, $t('tls.usePath')], [1, $t('tls.useText')]]"
            @update:model-value="switchEchMode" />
        </div>
        <template v-if="useEchPath == 0">
          <Field :label="$t('tls.certPath')">
            <input class="input mono" v-model="tls.ech.config_path" />
          </Field>
        </template>
        <template v-else>
          <Field :label="$t('tls.cert')">
            <textarea class="input mono" rows="4" spellcheck="false" v-model="echConfigText"></textarea>
          </Field>
        </template>
        <Field :label="$t('tls.queryServerName')">
          <input class="input mono" v-model="tls.ech.query_server_name" placeholder="ech.example.com" />
        </Field>
      </template>

      <template v-if="tls.fragment != undefined">
        <div style="display: flex; gap: 20px; flex-wrap: wrap; margin-bottom: 14px;">
          <SwitchLabel v-model="fragment" :label="$t('tls.fragment')" />
          <SwitchLabel v-if="tls.fragment" v-model="recordFragment" :label="$t('tls.recordFragment')" />
        </div>
        <Field v-if="tls.fragment" :label="$t('tls.fragmentDelay') + ' (' + $t('date.ms') + ')'">
          <input class="input mono" type="number" min="0" v-model.number="fragmentFallbackDelay" />
        </Field>
      </template>

      <div style="display: flex; margin-bottom: 14px;">
        <div style="flex: 1;" />
        <Pop :min-width="210" direction="up">
          <template #trigger="{ toggle }">
            <Btn variant="subtle" sm @click="toggle">{{ $t('tls.options') }}</Btn>
          </template>
          <div style="display: flex; flex-direction: column; gap: 2px; padding: 4px;">
            <div class="pop-item"><SwitchLabel v-model="optionCert" :label="$t('tls.cert')" /></div>
            <div class="pop-item"><SwitchLabel v-model="optionSNI" label="SNI" /></div>
            <div class="pop-item"><SwitchLabel v-model="optionALPN" label="ALPN" /></div>
            <div class="pop-item"><SwitchLabel v-model="optionMinV" :label="$t('tls.minVer')" /></div>
            <div class="pop-item"><SwitchLabel v-model="optionMaxV" :label="$t('tls.maxVer')" /></div>
            <div class="pop-item"><SwitchLabel v-model="optionCS" :label="$t('tls.cs')" /></div>
            <div class="pop-item"><SwitchLabel v-model="optionFP" label="UTLS" /></div>
            <div class="pop-item"><SwitchLabel v-model="optionReality" label="Reality" /></div>
            <div class="pop-item"><SwitchLabel v-model="optionEch" label="ECH" /></div>
            <div class="pop-item"><SwitchLabel v-model="optionFragment" :label="$t('tls.fragment')" /></div>
          </div>
        </Pop>
      </div>
    </div>
  </div>
</template>

<script lang="ts" setup>
import { computed, ref } from 'vue'
import { oTls, defaultOutTls } from '@/types/tls'
import Field from '@/components/ui/Field.vue'
import Btn from '@/components/ui/Btn.vue'
import Pop from '@/components/ui/Pop.vue'
import MSwitchRow from '@/components/ui/MSwitchRow.vue'
import SwitchLabel from '@/components/ui/SwitchLabel.vue'
import Segmented from '@/components/ui/Segmented.vue'
import SectionLabel from '@/components/ui/SectionLabel.vue'
import MultiPick from '@/components/ui/MultiPick.vue'
import ChipSelect from '@/components/ui/ChipSelect.vue'

const props = defineProps<{ outbound: any }>()

const usePath = ref<string | number>(props.outbound?.tls?.certificate ? 1 : 0)
const useEchPath = ref<string | number>(props.outbound?.tls?.ech?.config ? 1 : 0)

const alpnOptions = ['h3', 'h2', 'http/1.1']
const tlsVersions = ['1.0', '1.1', '1.2', '1.3']
const cipherSuiteItems = [
  { title: 'RSA-AES128-CBC-SHA', value: 'TLS_RSA_WITH_AES_128_CBC_SHA' },
  { title: 'RSA-AES256-CBC-SHA', value: 'TLS_RSA_WITH_AES_256_CBC_SHA' },
  { title: 'RSA-AES128-GCM-SHA256', value: 'TLS_RSA_WITH_AES_128_GCM_SHA256' },
  { title: 'RSA-AES256-GCM-SHA384', value: 'TLS_RSA_WITH_AES_256_GCM_SHA384' },
  { title: 'AES128-GCM-SHA256', value: 'TLS_AES_128_GCM_SHA256' },
  { title: 'AES256-GCM-SHA384', value: 'TLS_AES_256_GCM_SHA384' },
  { title: 'CHACHA20-POLY1305-SHA256', value: 'TLS_CHACHA20_POLY1305_SHA256' },
  { title: 'ECDHE-ECDSA-AES128-CBC-SHA', value: 'TLS_ECDHE_ECDSA_WITH_AES_128_CBC_SHA' },
  { title: 'ECDHE-ECDSA-AES256-CBC-SHA', value: 'TLS_ECDHE_ECDSA_WITH_AES_256_CBC_SHA' },
  { title: 'ECDHE-RSA-AES128-CBC-SHA', value: 'TLS_ECDHE_RSA_WITH_AES_128_CBC_SHA' },
  { title: 'ECDHE-RSA-AES256-CBC-SHA', value: 'TLS_ECDHE_RSA_WITH_AES_256_CBC_SHA' },
  { title: 'ECDHE-ECDSA-AES128-GCM-SHA256', value: 'TLS_ECDHE_ECDSA_WITH_AES_128_GCM_SHA256' },
  { title: 'ECDHE-ECDSA-AES256-GCM-SHA384', value: 'TLS_ECDHE_ECDSA_WITH_AES_256_GCM_SHA384' },
  { title: 'ECDHE-RSA-AES128-GCM-SHA256', value: 'TLS_ECDHE_RSA_WITH_AES_128_GCM_SHA256' },
  { title: 'ECDHE-RSA-AES256-GCM-SHA384', value: 'TLS_ECDHE_RSA_WITH_AES_256_GCM_SHA384' },
  { title: 'ECDHE-ECDSA-CHACHA20-POLY1305-SHA256', value: 'TLS_ECDHE_ECDSA_WITH_CHACHA20_POLY1305_SHA256' },
  { title: 'ECDHE-RSA-CHACHA20-POLY1305-SHA256', value: 'TLS_ECDHE_RSA_WITH_CHACHA20_POLY1305_SHA256' },
]
const fingerprints = [
  { title: 'Chrome', value: 'chrome' },
  { title: 'Firefox', value: 'firefox' },
  { title: 'Microsoft Edge', value: 'edge' },
  { title: 'Apple Safari', value: 'safari' },
  { title: '360', value: '360' },
  { title: 'QQ', value: 'qq' },
  { title: 'Apple IOS', value: 'ios' },
  { title: 'Android', value: 'android' },
  { title: 'Random', value: 'random' },
  { title: 'Randomized', value: 'randomized' },
]

const tls = computed((): oTls => props.outbound.tls)

const tlsOptional = computed((): boolean =>
  !['hysteria', 'hysteria2', 'tuic', 'shadowtls', 'anytls', 'naive'].includes(props.outbound.type))

const tlsEnable = computed({
  get: () => Object.hasOwn(tls.value, 'enabled') ? tls.value.enabled ?? false : false,
  set: (v: boolean) => { props.outbound.tls = v ? { enabled: true } : { enabled: false } },
})
const disableSni = computed({
  get: () => tls.value.disable_sni ?? false,
  set: (v: boolean) => { props.outbound.tls.disable_sni = v ? true : undefined },
})
const insecure = computed({
  get: () => tls.value.insecure ?? false,
  set: (v: boolean) => { props.outbound.tls.insecure = v ? true : undefined },
})
const fragment = computed({
  get: () => tls.value.fragment ?? false,
  set: (v: boolean) => { props.outbound.tls.fragment = v },
})
const recordFragment = computed({
  get: () => tls.value.record_fragment ?? false,
  set: (v: boolean) => { props.outbound.tls.record_fragment = v },
})
const alpnValue = computed({
  get: () => tls.value.alpn ?? [],
  set: (v: string[]) => { props.outbound.tls.alpn = v },
})
const cipherSuites = computed({
  get: () => tls.value.cipher_suites ?? [],
  set: (v: string[]) => { props.outbound.tls.cipher_suites = v },
})
const echConfigText = computed({
  get: (): string => tls.value.ech?.config ? tls.value.ech.config.join('\n') : '',
  set: (v: string) => { if (tls.value.ech) tls.value.ech.config = v.split('\n') },
})

const switchCertMode = (v: string | number) => {
  if (v == 0) {
    tls.value.certificate = undefined
    props.outbound.tls.certificate_path = ''
  } else {
    tls.value.certificate_path = undefined
    props.outbound.tls.certificate = ''
  }
}
const switchEchMode = (v: string | number) => {
  if (!tls.value.ech) return
  if (v == 0) delete tls.value.ech.config
  else delete tls.value.ech.config_path
}

const optionCert = computed({
  get: (): boolean => tls.value.certificate != undefined || tls.value.certificate_path != undefined,
  set: (v: boolean) => {
    usePath.value = 0
    if (v) {
      props.outbound.tls.certificate_path = ''
    } else {
      delete props.outbound.tls.certificate_path
      delete props.outbound.tls.certificate
    }
  },
})
const optionSNI = computed({
  get: (): boolean => tls.value.server_name != undefined,
  set: (v: boolean) => { props.outbound.tls.server_name = v ? '' : undefined },
})
const optionALPN = computed({
  get: (): boolean => tls.value.alpn != undefined,
  set: (v: boolean) => { props.outbound.tls.alpn = v ? defaultOutTls.alpn : undefined },
})
const optionMinV = computed({
  get: (): boolean => tls.value.min_version != undefined,
  set: (v: boolean) => { props.outbound.tls.min_version = v ? defaultOutTls.min_version : undefined },
})
const optionMaxV = computed({
  get: (): boolean => tls.value.max_version != undefined,
  set: (v: boolean) => { props.outbound.tls.max_version = v ? defaultOutTls.max_version : undefined },
})
const optionCS = computed({
  get: (): boolean => tls.value.cipher_suites != undefined,
  set: (v: boolean) => { props.outbound.tls.cipher_suites = v ? defaultOutTls.cipher_suites : undefined },
})
const optionFP = computed({
  get: (): boolean => tls.value.utls != undefined,
  set: (v: boolean) => { props.outbound.tls.utls = v ? defaultOutTls.utls : undefined },
})
const optionReality = computed({
  get: (): boolean => tls.value.reality != undefined,
  set: (v: boolean) => { props.outbound.tls.reality = v ? defaultOutTls.reality : undefined },
})
const optionEch = computed({
  get: (): boolean => tls.value.ech != undefined,
  set: (v: boolean) => { props.outbound.tls.ech = v ? defaultOutTls.ech : undefined },
})
const optionFragment = computed({
  get: (): boolean => tls.value.fragment != undefined,
  set: (v: boolean) => {
    if (v) {
      props.outbound.tls.fragment = false
    } else {
      delete props.outbound.tls.fragment
      delete props.outbound.tls.fragment_fallback_delay
      delete props.outbound.tls.record_fragment
    }
  },
})
const fragmentFallbackDelay = computed({
  get: (): number => parseInt(tls.value.fragment_fallback_delay?.replace('ms', '') ?? '500') ?? 500,
  set: (v: number) => { props.outbound.tls.fragment_fallback_delay = v > 0 ? `${v}ms` : undefined },
})
</script>
