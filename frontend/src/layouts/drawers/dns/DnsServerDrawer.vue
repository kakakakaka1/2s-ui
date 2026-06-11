<template>
  <MDrawer
    :open="open"
    icon="dns"
    :color="dnsColor(server.type)"
    :title="isNew ? $t('ui.dnssNew') : server.tag"
    :sub="$t('ui.dnssSub')"
    :save-label="isNew ? $t('ui.create') : $t('actions.save')"
    :width="500"
    @close="$emit('close')"
    @save="$emit('save', server)"
  >
    <div class="grid2">
      <Field :label="$t('type')">
        <Select v-model="serverType">
          <option v-for="dt in dnsTypes" :key="dt.value" :value="dt.value">{{ dt.title }}</option>
        </Select>
      </Field>
      <Field :label="$t('objects.tag')">
        <input class="input mono" v-model="server.tag" />
      </Field>
    </div>

    <!-- server address / port -->
    <div class="grid2" v-if="HasServer.includes(server.type)">
      <Field :label="$t('in.addr')">
        <input class="input mono" v-model="server.server" placeholder="1.1.1.1" />
      </Field>
      <Field :label="$t('in.port')">
        <input class="input mono" type="number" min="0" v-model.number="server.server_port" />
      </Field>
    </div>

    <!-- https / h3 path -->
    <Field v-if="HasHeaders.includes(server.type)" :label="$t('transport.path')">
      <input class="input mono" v-model="server.path" />
    </Field>

    <!-- local -->
    <div v-if="server.type === 'local'" style="margin-bottom: 15px;">
      <SwitchLabel :label="$t('dns.local.preferGo')" :model-value="!!server.prefer_go" @update:model-value="server.prefer_go = $event" />
    </div>

    <!-- dhcp -->
    <Field v-if="server.type === 'dhcp'" :label="$t('types.tun.ifName')">
      <input class="input mono" v-model="server.interface" placeholder="auto" />
    </Field>

    <!-- fakeip -->
    <div class="grid2" v-if="server.type === 'fakeip'">
      <Field :label="$t('dns.rule.inet4Range')">
        <input class="input mono" v-model="server.inet4_range" />
      </Field>
      <Field :label="$t('dns.rule.inet6Range')">
        <input class="input mono" v-model="server.inet6_range" />
      </Field>
    </div>

    <!-- hosts -->
    <template v-if="server.type === 'hosts'">
      <Field :label="$t('transport.path') + ' ' + $t('commaSeparated')">
        <input class="input mono" :value="hostsPath" @change="hostsPath = ($event.target as HTMLInputElement).value" />
      </Field>
      <div style="display: flex; align-items: center; margin-bottom: 10px;">
        <SectionLabel>Predefined</SectionLabel>
        <div style="flex: 1;" />
        <IconBtn name="plus" :title="$t('actions.add')" @click="addHostsPredefined" />
      </div>
      <div
        v-for="(pd, index) in hostsPredefined"
        :key="index"
        style="display: grid; grid-template-columns: 1fr 1fr 34px; gap: 8px; align-items: center; margin-bottom: 8px;"
      >
        <input class="input mono" style="height: 38px; font-size: 12.5px;" :value="pd.name" :placeholder="$t('setting.domain')" @change="updatePdsKey(index, ($event.target as HTMLInputElement).value)" />
        <input class="input mono" style="height: 38px; font-size: 12.5px;" :value="pd.value" :placeholder="$t('types.tun.addr') + ' ' + $t('commaSeparated')" @change="updatePdsValue(index, ($event.target as HTMLInputElement).value)" />
        <button type="button" class="btn btn-subtle btn-icon" style="height: 38px; width: 34px;" @click="delHostsPredefined(index)">
          <Ico name="close" :size="14" />
        </button>
      </div>
    </template>

    <!-- tailscale / resolved -->
    <template v-if="server.type === 'tailscale' || server.type === 'resolved'">
      <Field v-if="server.type === 'tailscale'" :label="$t('objects.endpoint')">
        <Select v-model="server.endpoint">
          <option v-for="e in tsTags" :key="e" :value="e">{{ e }}</option>
        </Select>
      </Field>
      <Field v-if="server.type === 'resolved'" :label="$t('objects.service')">
        <Select v-model="server.service">
          <option v-for="s in rslvdTags" :key="s" :value="s">{{ s }}</option>
        </Select>
      </Field>
      <MSwitchRow :label="$t('dns.rule.acceptDefault')" :model-value="!!server.accept_default_resolvers" @update:model-value="server.accept_default_resolvers = $event" />
    </template>

    <!-- ===================== Dial ===================== -->
    <template v-if="!WithoutDial.includes(server.type)">
      <hr class="form-divider" />
      <div style="display: flex; align-items: center; margin-bottom: 12px;">
        <SectionLabel>{{ $t('objects.dial') }}</SectionLabel>
        <div style="flex: 1;" />
        <Pop :width="280">
          <template #trigger="{ toggle }">
            <Btn variant="subtle" sm @click="toggle">{{ $t('dial.options') }}</Btn>
          </template>
          <div style="display: flex; flex-direction: column; gap: 2px; max-height: 300px; overflow-y: auto;">
            <div class="pop-item"><SwitchLabel :label="$t('listen.detour')" v-model="optionDetour" /></div>
            <div class="pop-item"><SwitchLabel :label="$t('dial.bindIf')" v-model="optionBind" /></div>
            <div class="pop-item"><SwitchLabel :label="$t('dial.bindIp4')" v-model="optionIPV4" /></div>
            <div class="pop-item"><SwitchLabel :label="$t('dial.bindIp6')" v-model="optionIPV6" /></div>
            <div class="pop-item"><SwitchLabel :label="$t('dial.bindNoPort')" v-model="optionBindNoPort" /></div>
            <div class="pop-item"><SwitchLabel label="Routing Mark" v-model="optionRM" /></div>
            <div class="pop-item"><SwitchLabel :label="$t('dial.reuseAddr')" v-model="optionRA" /></div>
            <div class="pop-item"><SwitchLabel :label="$t('listen.tcpOptions')" v-model="optionTCP" /></div>
            <div class="pop-item"><SwitchLabel :label="$t('listen.udpOptions')" v-model="optionUDP" /></div>
            <div class="pop-item"><SwitchLabel :label="$t('dial.connTimeout')" v-model="optionCT" /></div>
            <div class="pop-item"><SwitchLabel :label="$t('dial.tcpKeepAlive')" v-model="optionTcpKeepAlive" /></div>
            <div class="pop-item"><SwitchLabel :label="$t('dial.domainResolver')" v-model="optionDR" /></div>
          </div>
        </Pop>
      </div>

      <div class="grid2">
        <Field v-if="optionDetour" :label="$t('dial.detourText')">
          <Select v-model="server.detour">
            <option v-for="o in outTags" :key="o" :value="o">{{ o }}</option>
          </Select>
        </Field>
        <Field v-if="optionBind" :label="$t('dial.bindIf')">
          <input class="input mono" v-model="server.bind_interface" />
        </Field>
        <Field v-if="optionIPV4" :label="$t('dial.bindIp4')">
          <input class="input mono" v-model="server.inet4_bind_address" />
        </Field>
        <Field v-if="optionIPV6" :label="$t('dial.bindIp6')">
          <input class="input mono" v-model="server.inet6_bind_address" />
        </Field>
        <Field v-if="optionRM" label="Linux Routing Mark">
          <input class="input mono" type="number" min="0" v-model.number="routingMark" />
        </Field>
        <Field v-if="optionCT" :label="$t('dial.connTimeout') + ' (' + $t('date.s') + ')'">
          <input class="input mono" type="number" min="1" v-model.number="connectTimeout" />
        </Field>
        <Field v-if="optionDR" :label="$t('dial.domainResolver')">
          <Select v-model="server.domain_resolver">
            <option v-for="d in dnsTags" :key="d" :value="d">{{ d }}</option>
          </Select>
        </Field>
      </div>
      <div v-if="optionBindNoPort || optionRA || optionTCP || optionUDP" style="display: flex; gap: 20px; flex-wrap: wrap; margin-bottom: 15px;">
        <SwitchLabel v-if="optionBindNoPort" :label="$t('dial.bindNoPort')" :model-value="!!server.bind_address_no_port" @update:model-value="server.bind_address_no_port = $event" />
        <SwitchLabel v-if="optionRA" :label="$t('dial.reuseAddr')" :model-value="!!server.reuse_addr" @update:model-value="server.reuse_addr = $event" />
        <SwitchLabel v-if="optionTCP" label="TCP Fast Open" :model-value="!!server.tcp_fast_open" @update:model-value="server.tcp_fast_open = $event" />
        <SwitchLabel v-if="optionTCP" label="TCP Multi Path" :model-value="!!server.tcp_multi_path" @update:model-value="server.tcp_multi_path = $event" />
        <SwitchLabel v-if="optionUDP" label="UDP Fragment" :model-value="!!server.udp_fragment" @update:model-value="server.udp_fragment = $event" />
      </div>
      <template v-if="optionTcpKeepAlive">
        <div style="margin-bottom: 15px;">
          <SwitchLabel :label="$t('dial.disableTcpKeepAlive')" :model-value="!!server.disable_tcp_keep_alive" @update:model-value="server.disable_tcp_keep_alive = $event" />
        </div>
        <div class="grid2">
          <Field :label="$t('dial.tcpKeepAlive')">
            <input class="input mono" v-model="server.tcp_keep_alive" />
          </Field>
          <Field :label="$t('dial.tcpKeepAliveInterval')">
            <input class="input mono" v-model="server.tcp_keep_alive_interval" />
          </Field>
        </div>
      </template>
    </template>

    <!-- ===================== TLS ===================== -->
    <template v-if="HasTls.includes(server.type)">
      <hr class="form-divider" />
      <div style="display: flex; align-items: center; margin-bottom: 12px;">
        <SectionLabel>{{ $t('objects.tls') }}</SectionLabel>
        <div style="flex: 1;" />
        <Pop v-if="tls?.enabled" :width="250">
          <template #trigger="{ toggle }">
            <Btn variant="subtle" sm @click="toggle">{{ $t('tls.options') }}</Btn>
          </template>
          <div style="display: flex; flex-direction: column; gap: 2px; max-height: 300px; overflow-y: auto;">
            <div class="pop-item"><SwitchLabel :label="$t('tls.cert')" v-model="optionCert" /></div>
            <div class="pop-item"><SwitchLabel label="SNI" v-model="optionSNI" /></div>
            <div class="pop-item"><SwitchLabel label="ALPN" v-model="optionALPN" /></div>
            <div class="pop-item"><SwitchLabel :label="$t('tls.minVer')" v-model="optionMinV" /></div>
            <div class="pop-item"><SwitchLabel :label="$t('tls.maxVer')" v-model="optionMaxV" /></div>
            <div class="pop-item"><SwitchLabel :label="$t('tls.cs')" v-model="optionCS" /></div>
            <div class="pop-item"><SwitchLabel label="UTLS" v-model="optionFP" /></div>
            <div class="pop-item"><SwitchLabel label="Reality" v-model="optionReality" /></div>
            <div class="pop-item"><SwitchLabel label="ECH" v-model="optionEch" /></div>
            <div class="pop-item"><SwitchLabel :label="$t('tls.fragment')" v-model="optionFragment" /></div>
          </div>
        </Pop>
      </div>

      <MSwitchRow :label="$t('tls.enable')" v-model="tlsEnabled" />

      <template v-if="tls?.enabled">
        <div style="display: flex; gap: 20px; flex-wrap: wrap; margin-bottom: 15px;">
          <SwitchLabel :label="$t('tls.disableSni')" v-model="disableSni" />
          <SwitchLabel :label="$t('tls.insecure')" v-model="insecure" />
        </div>

        <!-- certificate -->
        <template v-if="optionCert">
          <Segmented v-model="certKind" block :options="[['path', $t('tls.usePath')], ['text', $t('tls.useText')]]" />
          <Field v-if="usePath === 0" :label="$t('tls.certPath')">
            <input class="input mono" v-model="tls.certificate_path" />
          </Field>
          <Field v-else :label="$t('tls.cert')">
            <textarea class="input mono" rows="4" style="height: auto; padding: 10px 12px; font-size: 12px;" v-model="tls.certificate"></textarea>
          </Field>
        </template>

        <div class="grid2">
          <Field v-if="tls.server_name !== undefined" label="SNI">
            <input class="input mono" v-model="tls.server_name" />
          </Field>
          <Field v-if="tls.min_version" :label="$t('tls.minVer')">
            <Select v-model="tls.min_version">
              <option v-for="v in tlsVersions" :key="v" :value="v">{{ v }}</option>
            </Select>
          </Field>
          <Field v-if="tls.max_version" :label="$t('tls.maxVer')">
            <Select v-model="tls.max_version">
              <option v-for="v in tlsVersions" :key="v" :value="v">{{ v }}</option>
            </Select>
          </Field>
          <Field v-if="tls.utls !== undefined" label="Fingerprint">
            <Select v-model="tls.utls.fingerprint">
              <option v-for="f in fingerprints" :key="f.value" :value="f.value">{{ f.title }}</option>
            </Select>
          </Field>
        </div>

        <ChipSelect
          v-if="tls.alpn"
          label="ALPN"
          :model-value="tls.alpn"
          :options="alpnOptions"
          style="margin-bottom: 15px;"
          @update:model-value="tls.alpn = $event"
        />
        <ChipSelect
          v-if="tls.cipher_suites !== undefined"
          :label="$t('tls.cs')"
          :model-value="tls.cipher_suites"
          :options="cipherSuites"
          style="margin-bottom: 15px;"
          @update:model-value="tls.cipher_suites = $event"
        />

        <!-- reality -->
        <div class="grid2" v-if="tls.reality !== undefined">
          <Field :label="$t('tls.pubKey')">
            <input class="input mono" v-model="tls.reality.public_key" />
          </Field>
          <Field label="Short ID">
            <input class="input mono" v-model="tls.reality.short_id" />
          </Field>
        </div>

        <!-- ECH -->
        <template v-if="tls.ech !== undefined">
          <SectionLabel style="margin-bottom: 10px;">ECH</SectionLabel>
          <Segmented v-model="echKind" block :options="[['path', $t('tls.usePath')], ['text', $t('tls.useText')]]" />
          <div class="grid2">
            <Field v-if="useEchPath === 0" :label="$t('tls.certPath')">
              <input class="input mono" v-model="tls.ech.config_path" />
            </Field>
            <Field v-else :label="$t('tls.cert')">
              <textarea class="input mono" rows="4" style="height: auto; padding: 10px 12px; font-size: 12px;" v-model="echConfigText"></textarea>
            </Field>
            <Field :label="$t('tls.queryServerName')">
              <input class="input mono" v-model="tls.ech.query_server_name" placeholder="ech.example.com" />
            </Field>
          </div>
        </template>

        <!-- fragment -->
        <template v-if="tls.fragment !== undefined">
          <div style="display: flex; gap: 20px; flex-wrap: wrap; margin-bottom: 15px;">
            <SwitchLabel :label="$t('tls.fragment')" v-model="tls.fragment" />
            <SwitchLabel v-if="tls.fragment" :label="$t('tls.recordFragment')" :model-value="!!tls.record_fragment" @update:model-value="tls.record_fragment = $event" />
          </div>
          <Field v-if="tls.fragment" :label="$t('tls.fragmentDelay') + ' (' + $t('date.ms') + ')'">
            <input class="input mono" type="number" min="0" v-model.number="fragmentFallbackDelay" />
          </Field>
        </template>
      </template>
    </template>

    <!-- ===================== Headers (https / h3) ===================== -->
    <template v-if="HasHeaders.includes(server.type)">
      <hr class="form-divider" />
      <div style="display: flex; align-items: center; margin-bottom: 10px;">
        <SectionLabel>{{ $t('objects.headers') }}</SectionLabel>
        <div style="flex: 1;" />
        <IconBtn name="plus" :title="$t('actions.add')" @click="addHeader" />
      </div>
      <div
        v-for="(header, index) in hdrs"
        :key="index"
        style="display: grid; grid-template-columns: 1fr 1fr 34px; gap: 8px; align-items: center; margin-bottom: 8px;"
      >
        <input class="input mono" style="height: 38px; font-size: 12.5px;" :value="header.name" :placeholder="$t('objects.key')" @change="updateHeaderKey(index, ($event.target as HTMLInputElement).value)" />
        <input class="input mono" style="height: 38px; font-size: 12.5px;" :value="header.value" :placeholder="$t('objects.value')" @change="updateHeaderValue(index, ($event.target as HTMLInputElement).value)" />
        <button type="button" class="btn btn-subtle btn-icon" style="height: 38px; width: 34px;" @click="delHeader(index)">
          <Ico name="close" :size="14" />
        </button>
      </div>
    </template>
  </MDrawer>
</template>

<script lang="ts" setup>
import Select from '@/components/ui/Select.vue'
import { computed, ref, watch } from 'vue'
import MDrawer from '@/components/ui/MDrawer.vue'
import Field from '@/components/ui/Field.vue'
import Btn from '@/components/ui/Btn.vue'
import Ico from '@/components/ui/Ico.vue'
import IconBtn from '@/components/ui/IconBtn.vue'
import Pop from '@/components/ui/Pop.vue'
import Segmented from '@/components/ui/Segmented.vue'
import ChipSelect from '@/components/ui/ChipSelect.vue'
import SwitchLabel from '@/components/ui/SwitchLabel.vue'
import SectionLabel from '@/components/ui/SectionLabel.vue'
import MSwitchRow from '@/components/ui/MSwitchRow.vue'
import Data from '@/store/modules/data'
import RandomUtil from '@/plugins/randomUtil'
import { dnsColor } from '@/plugins/colors'
import { DnsTypes, createDnsServer, DnsServer } from '@/types/dns'
import { defaultOutTls } from '@/types/tls'

const props = defineProps<{
  open: boolean
  index: number
  data: string
  tsTags: string[]
  rslvdTags: string[]
}>()
defineEmits<{ close: []; save: [data: any] }>()

const isNew = computed(() => props.index === -1)

const dnsTypes = Object.keys(DnsTypes).map((key, index) => ({ title: key, value: Object.values(DnsTypes)[index] }))
const HasServer: string[] = [DnsTypes.TCP, DnsTypes.UDP, DnsTypes.TLS, DnsTypes.QUIC, DnsTypes.HTTPS, DnsTypes.HTTP3]
const HasHeaders: string[] = [DnsTypes.HTTPS, DnsTypes.HTTP3]
const HasTls: string[] = [DnsTypes.TLS, DnsTypes.QUIC, DnsTypes.HTTPS, DnsTypes.HTTP3]
const WithoutDial: string[] = [DnsTypes.Hosts, DnsTypes.Tailscale, DnsTypes.FakeIP, DnsTypes.Resolved]

const alpnOptions = [
  { title: 'H3', value: 'h3' },
  { title: 'H2', value: 'h2' },
  { title: 'Http/1.1', value: 'http/1.1' },
]
const tlsVersions = ['1.0', '1.1', '1.2', '1.3']
const cipherSuites = [
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

const server = ref<any>(createDnsServer('local', { tag: 'dns-' + RandomUtil.randomSeq(3) }))
const usePath = ref(0)
const useEchPath = ref(0)

function init() {
  if (props.index !== -1) {
    server.value = JSON.parse(props.data)
  } else {
    server.value = createDnsServer('local', { tag: 'dns-' + RandomUtil.randomSeq(3) })
  }
  usePath.value = server.value?.tls?.certificate ? 1 : 0
  useEchPath.value = server.value?.tls?.ech?.config ? 1 : 0
}
watch(() => props.open, (v) => { if (v) init() })

// changing the type rebuilds the server with that type's defaults (legacy changeType)
const serverType = computed({
  get: () => server.value.type,
  set: (t: string) => {
    server.value = <DnsServer>createDnsServer(t, { tag: server.value.tag })
    usePath.value = 0
    useEchPath.value = 0
  },
})

const outTags = computed((): string[] => [
  ...Data().outbounds?.map((o: any) => o.tag),
  ...Data().endpoints?.map((e: any) => e.tag),
])
const dnsTags = computed((): string[] => Data().config.dns?.servers?.map((d: any) => d.tag) ?? [])

// ---- hosts: path list + predefined records (legacy computeds) ----
const hostsPath = computed({
  get: () => (Array.isArray(server.value.path) ? server.value.path.join(',') : server.value.path ?? ''),
  set: (v: string) => {
    server.value.path = v.length > 0 ? v.split(',').map((item: string) => item.trim()) : undefined
  },
})

const hostsPredefined = computed<{ name: string; value: string }[]>({
  get: () => {
    const pds: { name: string; value: string }[] = []
    const h = server.value.predefined
    if (h) {
      Object.keys(h).forEach((key) => {
        if (Array.isArray(h[key])) {
          pds.push({ name: key, value: h[key].join(',') })
        } else {
          pds.push({ name: key, value: h[key] })
        }
      })
    }
    return pds
  },
  set: (v: { name: string; value: string }[]) => {
    if (v.length > 0) {
      const pds: any = {}
      v.forEach((pd) => {
        pds[pd.name] = pd.value.split(',').map((item: string) => item.trim())
      })
      server.value.predefined = pds
    } else {
      server.value.predefined = undefined
    }
  },
})
const addHostsPredefined = () => { hostsPredefined.value = [...hostsPredefined.value, { name: 'localhost', value: '127.0.0.1,::1' }] }
const delHostsPredefined = (i: number) => { const pds = [...hostsPredefined.value]; pds.splice(i, 1); hostsPredefined.value = pds }
const updatePdsKey = (i: number, k: string) => { const pds = [...hostsPredefined.value]; pds[i] = { ...pds[i], name: k }; hostsPredefined.value = pds }
const updatePdsValue = (i: number, v: string) => { const pds = [...hostsPredefined.value]; pds[i] = { ...pds[i], value: v }; hostsPredefined.value = pds }

// ---- headers (https / h3, legacy Headers.vue computed) ----
const hdrs = computed<{ name: string; value: string }[]>({
  get: () => {
    const headers: { name: string; value: string }[] = []
    const h = server.value.headers
    if (h) {
      Object.keys(h).forEach((key) => {
        if (Array.isArray(h[key])) {
          h[key].forEach((v: string) => headers.push({ name: key, value: v }))
        } else {
          headers.push({ name: key, value: h[key] })
        }
      })
    }
    return headers
  },
  set: (v: { name: string; value: string }[]) => {
    if (v.length > 0) {
      const headers: any = {}
      v.forEach((h) => {
        if (headers[h.name]) {
          if (Array.isArray(headers[h.name])) {
            headers[h.name].push(h.value)
          } else {
            headers[h.name] = [headers[h.name], h.value]
          }
        } else {
          headers[h.name] = h.value
        }
      })
      server.value.headers = headers
    } else {
      server.value.headers = undefined
    }
  },
})
const addHeader = () => { hdrs.value = [...hdrs.value, { name: 'Host', value: '' }] }
const delHeader = (i: number) => { const h = [...hdrs.value]; h.splice(i, 1); hdrs.value = h }
const updateHeaderKey = (i: number, k: string) => { const h = [...hdrs.value]; h[i] = { ...h[i], name: k }; hdrs.value = h }
const updateHeaderValue = (i: number, v: string) => { const h = [...hdrs.value]; h[i] = { ...h[i], value: v }; hdrs.value = h }

// ---- dial options (legacy Dial.vue computeds, dial == server object) ----
const connectTimeout = computed({
  get: () => (server.value.connect_timeout ? parseInt(server.value.connect_timeout.replace('s', '')) : 5),
  set: (v: number) => { server.value.connect_timeout = v > 0 ? v + 's' : '5s' },
})
const routingMark = computed({
  get: () => server.value.routing_mark ?? 0,
  set: (v: number) => { server.value.routing_mark = v > 0 ? v : 0 },
})
const optionDetour = computed({
  get: () => server.value.detour !== undefined,
  set: (v: boolean) => { if (v) server.value.detour = outTags.value[0] ?? ''; else delete server.value.detour },
})
const optionBind = computed({
  get: () => server.value.bind_interface !== undefined,
  set: (v: boolean) => { if (v) server.value.bind_interface = ''; else delete server.value.bind_interface },
})
const optionIPV4 = computed({
  get: () => server.value.inet4_bind_address !== undefined,
  set: (v: boolean) => { if (v) server.value.inet4_bind_address = ''; else delete server.value.inet4_bind_address },
})
const optionIPV6 = computed({
  get: () => server.value.inet6_bind_address !== undefined,
  set: (v: boolean) => { if (v) server.value.inet6_bind_address = ''; else delete server.value.inet6_bind_address },
})
const optionBindNoPort = computed({
  get: () => server.value.bind_address_no_port !== undefined,
  set: (v: boolean) => { if (v) server.value.bind_address_no_port = true; else delete server.value.bind_address_no_port },
})
const optionRM = computed({
  get: () => server.value.routing_mark !== undefined,
  set: (v: boolean) => { if (v) server.value.routing_mark = 0; else delete server.value.routing_mark },
})
const optionRA = computed({
  get: () => server.value.reuse_addr !== undefined,
  set: (v: boolean) => { if (v) server.value.reuse_addr = true; else delete server.value.reuse_addr },
})
const optionTCP = computed({
  get: () => server.value.tcp_fast_open !== undefined && server.value.tcp_multi_path !== undefined,
  set: (v: boolean) => {
    if (v) {
      server.value.tcp_fast_open = false
      server.value.tcp_multi_path = false
    } else {
      delete server.value.tcp_fast_open
      delete server.value.tcp_multi_path
    }
  },
})
const optionUDP = computed({
  get: () => server.value.udp_fragment !== undefined,
  set: (v: boolean) => { if (v) server.value.udp_fragment = true; else delete server.value.udp_fragment },
})
const optionCT = computed({
  get: () => server.value.connect_timeout !== undefined,
  set: (v: boolean) => { if (v) server.value.connect_timeout = '5s'; else delete server.value.connect_timeout },
})
const optionTcpKeepAlive = computed({
  get: () =>
    server.value.disable_tcp_keep_alive !== undefined ||
    server.value.tcp_keep_alive !== undefined ||
    server.value.tcp_keep_alive_interval !== undefined,
  set: (v: boolean) => {
    if (v) {
      server.value.tcp_keep_alive = '5m'
      server.value.tcp_keep_alive_interval = '75s'
    } else {
      delete server.value.disable_tcp_keep_alive
      delete server.value.tcp_keep_alive
      delete server.value.tcp_keep_alive_interval
    }
  },
})
const optionDR = computed({
  get: () => server.value.domain_resolver !== undefined,
  set: (v: boolean) => { server.value.domain_resolver = v ? dnsTags.value[0] ?? '' : undefined },
})

// ---- TLS (legacy OutTLS.vue computeds, outbound == server object) ----
const tls = computed<any>(() => server.value.tls)

const tlsEnabled = computed({
  get: () => (tls.value && Object.hasOwn(tls.value, 'enabled') ? !!tls.value.enabled : false),
  set: (v: boolean) => { server.value.tls = v ? { enabled: true } : { enabled: false } },
})
const disableSni = computed({
  get: () => tls.value?.disable_sni ?? false,
  set: (v: boolean) => { server.value.tls.disable_sni = v ? true : undefined },
})
const insecure = computed({
  get: () => tls.value?.insecure ?? false,
  set: (v: boolean) => { server.value.tls.insecure = v ? true : undefined },
})
const echConfigText = computed({
  get: () => (tls.value?.ech?.config ? tls.value.ech.config.join('\n') : ''),
  set: (v: string) => { if (tls.value?.ech) tls.value.ech.config = v.split('\n') },
})
const certKind = computed({
  get: () => (usePath.value === 0 ? 'path' : 'text'),
  set: (v) => {
    if (v === 'path') {
      usePath.value = 0
      tls.value.certificate = undefined
      tls.value.certificate_path = ''
    } else {
      usePath.value = 1
      tls.value.certificate_path = undefined
      tls.value.certificate = ''
    }
  },
})
const echKind = computed({
  get: () => (useEchPath.value === 0 ? 'path' : 'text'),
  set: (v) => {
    if (v === 'path') {
      useEchPath.value = 0
      delete tls.value.ech?.config
    } else {
      useEchPath.value = 1
      delete tls.value.ech?.config_path
    }
  },
})
const optionCert = computed({
  get: () => tls.value?.certificate !== undefined || tls.value?.certificate_path !== undefined,
  set: (v: boolean) => {
    usePath.value = 0
    if (v) {
      server.value.tls.certificate_path = ''
    } else {
      delete server.value.tls.certificate_path
      delete server.value.tls.certificate
    }
  },
})
const optionSNI = computed({
  get: () => tls.value?.server_name !== undefined,
  set: (v: boolean) => { server.value.tls.server_name = v ? '' : undefined },
})
const optionALPN = computed({
  get: () => tls.value?.alpn !== undefined,
  set: (v: boolean) => { server.value.tls.alpn = v ? defaultOutTls.alpn : undefined },
})
const optionMinV = computed({
  get: () => tls.value?.min_version !== undefined,
  set: (v: boolean) => { server.value.tls.min_version = v ? defaultOutTls.min_version : undefined },
})
const optionMaxV = computed({
  get: () => tls.value?.max_version !== undefined,
  set: (v: boolean) => { server.value.tls.max_version = v ? defaultOutTls.max_version : undefined },
})
const optionCS = computed({
  get: () => tls.value?.cipher_suites !== undefined,
  set: (v: boolean) => { server.value.tls.cipher_suites = v ? defaultOutTls.cipher_suites : undefined },
})
const optionFP = computed({
  get: () => tls.value?.utls !== undefined,
  set: (v: boolean) => { server.value.tls.utls = v ? defaultOutTls.utls : undefined },
})
const optionReality = computed({
  get: () => tls.value?.reality !== undefined,
  set: (v: boolean) => { server.value.tls.reality = v ? defaultOutTls.reality : undefined },
})
const optionEch = computed({
  get: () => tls.value?.ech !== undefined,
  set: (v: boolean) => { server.value.tls.ech = v ? defaultOutTls.ech : undefined },
})
const optionFragment = computed({
  get: () => tls.value?.fragment !== undefined,
  set: (v: boolean) => {
    if (v) {
      server.value.tls.fragment = false
    } else {
      delete server.value.tls.fragment
      delete server.value.tls.fragment_fallback_delay
      delete server.value.tls.record_fragment
    }
  },
})
const fragmentFallbackDelay = computed({
  get: () => parseInt(tls.value?.fragment_fallback_delay?.replace('ms', '') ?? '500') ?? 500,
  set: (v: number) => { server.value.tls.fragment_fallback_delay = v > 0 ? `${v}ms` : undefined },
})
</script>
