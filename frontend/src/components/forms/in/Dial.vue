<template>
  <div>
    <div style="display: flex; align-items: center; gap: 10px; margin-bottom: 12px;">
      <SectionLabel style="flex: 1;">{{ $t('objects.dial') }}</SectionLabel>
      <Pop :min-width="250">
        <template #trigger="{ toggle }">
          <button type="button" class="btn btn-ghost btn-sm" @click="toggle">
            <Ico name="settings" :size="14" /> {{ $t('dial.options') }}
          </button>
        </template>
        <div style="display: flex; flex-direction: column; gap: 2px; padding: 5px; max-height: 300px; overflow-y: auto;">
          <SwitchLabel v-if="mode != 'client'" v-model="optionDetour" :label="$t('listen.detour')" />
          <SwitchLabel v-if="mode != 'client'" v-model="optionBind" :label="$t('dial.bindIf')" />
          <SwitchLabel v-if="mode != 'client'" v-model="optionIPV4" :label="$t('dial.bindIp4')" />
          <SwitchLabel v-if="mode != 'client'" v-model="optionIPV6" :label="$t('dial.bindIp6')" />
          <SwitchLabel v-if="mode != 'client'" v-model="optionBindNoPort" :label="$t('dial.bindNoPort')" />
          <SwitchLabel v-if="mode != 'client'" v-model="optionRM" label="Routing Mark" />
          <SwitchLabel v-if="mode != 'client'" v-model="optionRA" :label="$t('dial.reuseAddr')" />
          <SwitchLabel v-model="optionTCP" :label="$t('listen.tcpOptions')" />
          <SwitchLabel v-model="optionUDP" :label="$t('listen.udpOptions')" />
          <SwitchLabel v-model="optionCT" :label="$t('dial.connTimeout')" />
          <SwitchLabel v-model="optionTcpKeepAlive" :label="$t('dial.tcpKeepAlive')" />
          <SwitchLabel v-if="mode != 'client'" v-model="optionDR" :label="$t('dial.domainResolver')" />
        </div>
      </Pop>
    </div>
    <div class="grid2">
      <Field v-if="optionDetour" :label="$t('dial.detourText')">
        <select class="input" v-model="dial.detour">
          <option v-for="tg in outTags" :key="tg" :value="tg">{{ tg }}</option>
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
    </div>
    <div
      v-if="optionBindNoPort || optionRA || optionUDP"
      style="display: flex; gap: 24px; flex-wrap: wrap; margin-bottom: 15px;"
    >
      <SwitchLabel v-if="optionBindNoPort" v-model="dial.bind_address_no_port" :label="$t('dial.bindNoPort')" />
      <SwitchLabel v-if="optionRA" v-model="dial.reuse_addr" :label="$t('dial.reuseAddr')" />
      <SwitchLabel v-if="optionUDP" v-model="dial.udp_fragment" label="UDP Fragment" />
    </div>
    <div v-if="optionTCP" class="fade-up" style="display: flex; gap: 24px; flex-wrap: wrap; margin-bottom: 15px;">
      <SwitchLabel v-model="dial.tcp_fast_open" label="TCP Fast Open" />
      <SwitchLabel v-model="dial.tcp_multi_path" label="TCP Multi Path" />
    </div>
    <div v-if="optionTcpKeepAlive" class="fade-up">
      <div style="display: flex; gap: 24px; flex-wrap: wrap; margin-bottom: 15px;">
        <SwitchLabel v-model="dial.disable_tcp_keep_alive" :label="$t('dial.disableTcpKeepAlive')" />
      </div>
      <div class="grid2">
        <Field :label="$t('dial.tcpKeepAlive')">
          <input class="input mono" v-model="dial.tcp_keep_alive" />
        </Field>
        <Field :label="$t('dial.tcpKeepAliveInterval')">
          <input class="input mono" v-model="dial.tcp_keep_alive_interval" />
        </Field>
      </div>
    </div>
    <Field v-if="optionDR" :label="$t('dial.domainResolver')">
      <select class="input" v-model="dial.domain_resolver">
        <option v-for="tg in dnsTags" :key="tg" :value="tg">{{ tg }}</option>
      </select>
    </Field>
  </div>
</template>

<script lang="ts" setup>
import { computed } from 'vue'
import Data from '@/store/modules/data'
import Field from '@/components/ui/Field.vue'
import SectionLabel from '@/components/ui/SectionLabel.vue'
import SwitchLabel from '@/components/ui/SwitchLabel.vue'
import Pop from '@/components/ui/Pop.vue'
import Ico from '@/components/ui/Ico.vue'

const props = defineProps<{ dial: any; mode?: string }>()

const outTags = computed((): string[] => [
  ...Data().outbounds?.map((o: any) => o.tag),
  ...Data().endpoints?.map((e: any) => e.tag),
])

const dnsTags = computed((): string[] => Data().config.dns?.servers?.map((d: any) => d.tag) ?? [])

const connectTimeout = computed({
  get: (): number => (props.dial.connect_timeout ? parseInt(props.dial.connect_timeout.replace('s', '')) : 5),
  set: (v: number) => { props.dial.connect_timeout = v > 0 ? v + 's' : '5s' },
})

const routingMark = computed({
  get: (): number => props.dial.routing_mark ?? 0,
  set: (v: number) => { props.dial.routing_mark = v > 0 ? v : 0 },
})

const optionDetour = computed({
  get: (): boolean => props.dial.detour != undefined,
  set: (v: boolean) => { v ? (props.dial.detour = outTags.value[0] ?? '') : delete props.dial.detour },
})

const optionBind = computed({
  get: (): boolean => props.dial.bind_interface != undefined,
  set: (v: boolean) => { v ? (props.dial.bind_interface = '') : delete props.dial.bind_interface },
})

const optionIPV4 = computed({
  get: (): boolean => props.dial.inet4_bind_address != undefined,
  set: (v: boolean) => { v ? (props.dial.inet4_bind_address = '') : delete props.dial.inet4_bind_address },
})

const optionIPV6 = computed({
  get: (): boolean => props.dial.inet6_bind_address != undefined,
  set: (v: boolean) => { v ? (props.dial.inet6_bind_address = '') : delete props.dial.inet6_bind_address },
})

const optionBindNoPort = computed({
  get: (): boolean => props.dial.bind_address_no_port != undefined,
  set: (v: boolean) => { v ? (props.dial.bind_address_no_port = true) : delete props.dial.bind_address_no_port },
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
  set: (v: boolean) => { v ? (props.dial.routing_mark = 0) : delete props.dial.routing_mark },
})

const optionRA = computed({
  get: (): boolean => props.dial.reuse_addr != undefined,
  set: (v: boolean) => { v ? (props.dial.reuse_addr = true) : delete props.dial.reuse_addr },
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
  set: (v: boolean) => { v ? (props.dial.udp_fragment = true) : delete props.dial.udp_fragment },
})

const optionCT = computed({
  get: (): boolean => props.dial.connect_timeout != undefined,
  set: (v: boolean) => { v ? (props.dial.connect_timeout = '5s') : delete props.dial.connect_timeout },
})

const optionDR = computed({
  get: (): boolean => props.dial.domain_resolver != undefined,
  set: (v: boolean) => { props.dial.domain_resolver = v ? dnsTags.value[0] ?? '' : undefined },
})
</script>
