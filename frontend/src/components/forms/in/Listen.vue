<template>
  <div>
    <div style="display: flex; align-items: center; gap: 10px; margin-bottom: 12px;">
      <SectionLabel style="flex: 1;">{{ $t('objects.listen') }}</SectionLabel>
      <Pop :min-width="240">
        <template #trigger="{ toggle }">
          <button type="button" class="btn btn-ghost btn-sm" @click="toggle">
            <Ico name="settings" :size="14" /> {{ $t('listen.options') }}
          </button>
        </template>
        <div style="display: flex; flex-direction: column; gap: 2px; padding: 5px;">
          <SwitchLabel v-model="optionDetour" :label="$t('listen.detour')" />
          <SwitchLabel v-model="optionTCP" :label="$t('listen.tcpOptions')" />
          <SwitchLabel v-model="optionUDP" :label="$t('listen.udpOptions')" />
          <SwitchLabel v-model="optionTcpKeepAlive" :label="$t('listen.tcpKeepAlive')" />
        </div>
      </Pop>
    </div>
    <div class="grid2">
      <Field :label="$t('in.addr')">
        <input class="input mono" required v-model="data.listen" />
      </Field>
      <Field :label="$t('in.port')">
        <input class="input mono" type="number" min="1" max="65535" required v-model.number="data.listen_port" />
      </Field>
    </div>
    <Field v-if="optionDetour" :label="$t('listen.detourText')">
      <Select v-model="data.detour">
        <option v-for="tg in inTags ?? []" :key="tg" :value="tg">{{ tg }}</option>
      </Select>
    </Field>
    <div v-if="optionTCP" class="fade-up" style="display: flex; gap: 24px; flex-wrap: wrap; margin-bottom: 15px;">
      <SwitchLabel v-model="data.tcp_fast_open" label="TCP Fast Open" />
      <SwitchLabel v-model="data.tcp_multi_path" label="TCP Multi Path" />
    </div>
    <div v-if="optionUDP" class="fade-up">
      <div style="display: flex; gap: 24px; flex-wrap: wrap; margin-bottom: 15px;">
        <SwitchLabel v-model="data.udp_fragment" label="UDP Fragment" />
      </div>
      <Field :label="'UDP NAT expiration' + ' (' + $t('date.m') + ')'">
        <input class="input mono" type="number" min="1" v-model.number="udpTimeout" />
      </Field>
    </div>
    <div v-if="optionTcpKeepAlive" class="fade-up">
      <div style="display: flex; gap: 24px; flex-wrap: wrap; margin-bottom: 15px;">
        <SwitchLabel v-model="data.disable_tcp_keep_alive" :label="$t('listen.disableTcpKeepAlive')" />
      </div>
      <div class="grid2">
        <Field :label="$t('listen.tcpKeepAlive')">
          <input class="input mono" v-model="data.tcp_keep_alive" />
        </Field>
        <Field :label="$t('listen.tcpKeepAliveInterval')">
          <input class="input mono" v-model="data.tcp_keep_alive_interval" />
        </Field>
      </div>
    </div>
  </div>
</template>

<script lang="ts" setup>
import Select from '@/components/ui/Select.vue'
import { computed } from 'vue'
import Field from '@/components/ui/Field.vue'
import SectionLabel from '@/components/ui/SectionLabel.vue'
import SwitchLabel from '@/components/ui/SwitchLabel.vue'
import Pop from '@/components/ui/Pop.vue'
import Ico from '@/components/ui/Ico.vue'

const props = defineProps<{ data: any; inTags?: string[] }>()

const udpTimeout = computed({
  get: (): number => (props.data.udp_timeout ? parseInt(props.data.udp_timeout.replace('m', '')) : 5),
  set: (v: number) => { props.data.udp_timeout = v > 0 ? v + 'm' : '5m' },
})

const optionTCP = computed({
  get: (): boolean =>
    props.data.tcp_fast_open != undefined &&
    props.data.tcp_multi_path != undefined,
  set: (v: boolean) => {
    props.data.tcp_fast_open = v ? false : undefined
    props.data.tcp_multi_path = v ? false : undefined
  },
})

const optionUDP = computed({
  get: (): boolean =>
    props.data.udp_fragment != undefined &&
    props.data.udp_timeout != undefined,
  set: (v: boolean) => {
    props.data.udp_fragment = v ? false : undefined
    props.data.udp_timeout = v ? '5m' : undefined
  },
})

const optionDetour = computed({
  get: (): boolean => props.data.detour != undefined,
  set: (v: boolean) => { props.data.detour = v ? props.inTags?.[0] ?? '' : undefined },
})

const optionTcpKeepAlive = computed({
  get: (): boolean =>
    props.data.disable_tcp_keep_alive != undefined ||
    props.data.tcp_keep_alive != undefined ||
    props.data.tcp_keep_alive_interval != undefined,
  set: (v: boolean) => {
    if (v) {
      props.data.tcp_keep_alive = '5m'
      props.data.tcp_keep_alive_interval = '75s'
    } else {
      delete props.data.disable_tcp_keep_alive
      delete props.data.tcp_keep_alive
      delete props.data.tcp_keep_alive_interval
    }
  },
})
</script>
