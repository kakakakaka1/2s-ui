<template>
  <div style="margin-bottom: 15px;">
    <div style="display: flex; align-items: center; gap: 10px; margin-bottom: 12px;">
      <SectionLabel>{{ $t('objects.dial') }}</SectionLabel>
      <div class="grow" style="flex: 1;" />
      <Pop :min-width="230" direction="up">
        <template #trigger="{ toggle }">
          <Btn variant="subtle" sm @click="toggle">{{ $t('dial.options') }}</Btn>
        </template>
        <div style="display: flex; flex-direction: column; gap: 2px; padding: 4px;">
          <div class="pop-item"><SwitchLabel v-model="optionDetour" :label="$t('listen.detour')" /></div>
          <div class="pop-item"><SwitchLabel v-model="optionBind" :label="$t('dial.bindIf')" /></div>
          <div class="pop-item"><SwitchLabel v-model="optionIPV4" :label="$t('dial.bindIp4')" /></div>
          <div class="pop-item"><SwitchLabel v-model="optionIPV6" :label="$t('dial.bindIp6')" /></div>
          <div class="pop-item"><SwitchLabel v-model="optionBindNoPort" :label="$t('dial.bindNoPort')" /></div>
          <div class="pop-item"><SwitchLabel v-model="optionRM" label="Routing Mark" /></div>
          <div class="pop-item"><SwitchLabel v-model="optionRA" :label="$t('dial.reuseAddr')" /></div>
          <div class="pop-item"><SwitchLabel v-model="optionTCP" :label="$t('listen.tcpOptions')" /></div>
          <div class="pop-item"><SwitchLabel v-model="optionUDP" :label="$t('listen.udpOptions')" /></div>
          <div class="pop-item"><SwitchLabel v-model="optionCT" :label="$t('dial.connTimeout')" /></div>
          <div class="pop-item"><SwitchLabel v-model="optionTcpKeepAlive" :label="$t('dial.tcpKeepAlive')" /></div>
          <div class="pop-item"><SwitchLabel v-model="optionDR" :label="$t('dial.domainResolver')" /></div>
        </div>
      </Pop>
    </div>

    <div class="grid2">
      <Field v-if="optionDetour" :label="$t('dial.detourText')">
        <select class="input" v-model="dial.detour">
          <option v-for="tag in outTags" :key="tag" :value="tag">{{ tag }}</option>
        </select>
      </Field>
      <Field v-if="optionBind" :label="$t('dial.bindIf')">
        <input class="input mono" v-model="dial.bind_interface" />
      </Field>
      <Field v-if="optionIPV4" :label="$t('dial.bindIp4')">
        <input class="input mono" v-model="dial.inet4_bind_address" />
      </Field>
      <Field v-if="optionIPV6" :label="$t('dial.bindIp6')">
        <input class="input mono" v-model="dial.inet6_bind_address" />
      </Field>
      <Field v-if="optionRM" label="Linux Routing Mark">
        <input class="input mono" type="number" min="0" v-model.number="routingMark" />
      </Field>
      <Field v-if="optionCT" :label="$t('dial.connTimeout') + ' (' + $t('date.s') + ')'">
        <input class="input mono" type="number" min="1" v-model.number="connectTimeout" />
      </Field>
      <Field v-if="optionDR" :label="$t('dial.domainResolver')">
        <select class="input" v-model="dial.domain_resolver">
          <option v-for="tag in dnsTags" :key="tag" :value="tag">{{ tag }}</option>
        </select>
      </Field>
    </div>

    <div
      v-if="optionBindNoPort || optionRA || optionTCP || optionUDP"
      style="display: flex; gap: 20px; flex-wrap: wrap; margin-bottom: 14px;"
    >
      <SwitchLabel v-if="optionBindNoPort" v-model="bindNoPort" :label="$t('dial.bindNoPort')" />
      <SwitchLabel v-if="optionRA" v-model="reuseAddr" :label="$t('dial.reuseAddr')" />
      <SwitchLabel v-if="optionTCP" v-model="tcpFastOpen" label="TCP Fast Open" />
      <SwitchLabel v-if="optionTCP" v-model="tcpMultiPath" label="TCP Multi Path" />
      <SwitchLabel v-if="optionUDP" v-model="udpFragment" label="UDP Fragment" />
    </div>

    <template v-if="optionTcpKeepAlive">
      <div style="margin-bottom: 12px;">
        <SwitchLabel v-model="disableTcpKeepAlive" :label="$t('dial.disableTcpKeepAlive')" />
      </div>
      <div class="grid2">
        <Field :label="$t('dial.tcpKeepAlive')">
          <input class="input mono" v-model="dial.tcp_keep_alive" />
        </Field>
        <Field :label="$t('dial.tcpKeepAliveInterval')">
          <input class="input mono" v-model="dial.tcp_keep_alive_interval" />
        </Field>
      </div>
    </template>
  </div>
</template>

<script lang="ts" setup>
import { computed } from 'vue'
import Data from '@/store/modules/data'
import Field from '@/components/ui/Field.vue'
import Btn from '@/components/ui/Btn.vue'
import Pop from '@/components/ui/Pop.vue'
import SwitchLabel from '@/components/ui/SwitchLabel.vue'
import SectionLabel from '@/components/ui/SectionLabel.vue'

const props = defineProps<{ dial: any }>()

const outTags = computed(() => [
  ...Data().outbounds?.map((o: any) => o.tag),
  ...Data().endpoints?.map((e: any) => e.tag),
])
const dnsTags = computed(() => Data().config.dns?.servers?.map((d: any) => d.tag) ?? [])

const connectTimeout = computed({
  get: () => props.dial.connect_timeout ? parseInt(props.dial.connect_timeout.replace('s', '')) : 5,
  set: (v: number) => { props.dial.connect_timeout = v > 0 ? v + 's' : '5s' },
})
const routingMark = computed({
  get: () => props.dial.routing_mark ?? 0,
  set: (v: number) => { props.dial.routing_mark = v > 0 ? v : 0 },
})
const bindNoPort = computed({
  get: () => props.dial.bind_address_no_port ?? false,
  set: (v: boolean) => { props.dial.bind_address_no_port = v },
})
const reuseAddr = computed({
  get: () => props.dial.reuse_addr ?? false,
  set: (v: boolean) => { props.dial.reuse_addr = v },
})
const tcpFastOpen = computed({
  get: () => props.dial.tcp_fast_open ?? false,
  set: (v: boolean) => { props.dial.tcp_fast_open = v },
})
const tcpMultiPath = computed({
  get: () => props.dial.tcp_multi_path ?? false,
  set: (v: boolean) => { props.dial.tcp_multi_path = v },
})
const udpFragment = computed({
  get: () => props.dial.udp_fragment ?? false,
  set: (v: boolean) => { props.dial.udp_fragment = v },
})
const disableTcpKeepAlive = computed({
  get: () => props.dial.disable_tcp_keep_alive ?? false,
  set: (v: boolean) => { props.dial.disable_tcp_keep_alive = v },
})

const optionDetour = computed({
  get: (): boolean => props.dial.detour != undefined,
  set: (v: boolean) => { v ? props.dial.detour = outTags.value[0] ?? '' : delete props.dial.detour },
})
const optionBind = computed({
  get: (): boolean => props.dial.bind_interface != undefined,
  set: (v: boolean) => { v ? props.dial.bind_interface = '' : delete props.dial.bind_interface },
})
const optionIPV4 = computed({
  get: (): boolean => props.dial.inet4_bind_address != undefined,
  set: (v: boolean) => { v ? props.dial.inet4_bind_address = '' : delete props.dial.inet4_bind_address },
})
const optionIPV6 = computed({
  get: (): boolean => props.dial.inet6_bind_address != undefined,
  set: (v: boolean) => { v ? props.dial.inet6_bind_address = '' : delete props.dial.inet6_bind_address },
})
const optionBindNoPort = computed({
  get: (): boolean => props.dial.bind_address_no_port != undefined,
  set: (v: boolean) => { v ? props.dial.bind_address_no_port = true : delete props.dial.bind_address_no_port },
})
const optionTcpKeepAlive = computed({
  get: (): boolean =>
    props.dial.disable_tcp_keep_alive != undefined ||
    props.dial.tcp_keep_alive != undefined ||
    props.dial.tcp_keep_alive_interval != undefined,
  set: (v: boolean) => {
    if (v) {
      props.dial.tcp_keep_alive = '5m'
      props.dial.tcp_keep_alive_interval = '75s'
    } else {
      delete props.dial.disable_tcp_keep_alive
      delete props.dial.tcp_keep_alive
      delete props.dial.tcp_keep_alive_interval
    }
  },
})
const optionRM = computed({
  get: (): boolean => props.dial.routing_mark != undefined,
  set: (v: boolean) => { v ? props.dial.routing_mark = 0 : delete props.dial.routing_mark },
})
const optionRA = computed({
  get: (): boolean => props.dial.reuse_addr != undefined,
  set: (v: boolean) => { v ? props.dial.reuse_addr = true : delete props.dial.reuse_addr },
})
const optionTCP = computed({
  get: (): boolean =>
    props.dial.tcp_fast_open != undefined &&
    props.dial.tcp_multi_path != undefined,
  set: (v: boolean) => {
    if (v) {
      props.dial.tcp_fast_open = false
      props.dial.tcp_multi_path = false
    } else {
      delete props.dial.tcp_fast_open
      delete props.dial.tcp_multi_path
    }
  },
})
const optionUDP = computed({
  get: (): boolean => props.dial.udp_fragment != undefined,
  set: (v: boolean) => { v ? props.dial.udp_fragment = true : delete props.dial.udp_fragment },
})
const optionCT = computed({
  get: (): boolean => props.dial.connect_timeout != undefined,
  set: (v: boolean) => { v ? props.dial.connect_timeout = '5s' : delete props.dial.connect_timeout },
})
const optionDR = computed({
  get: (): boolean => props.dial.domain_resolver != undefined,
  set: (v: boolean) => { props.dial.domain_resolver = v ? dnsTags.value[0] ?? '' : undefined },
})
</script>
