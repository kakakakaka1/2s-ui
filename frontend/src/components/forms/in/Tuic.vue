<template>
  <div>
    <SectionLabel style="margin-bottom: 12px;">TUIC</SectionLabel>
    <!-- Outbound only -->
    <template v-if="direction == 'out'">
      <div class="grid2">
        <Field label="UUID">
          <input class="input mono" v-model="data.uuid" />
        </Field>
        <Field :label="$t('types.pw')">
          <input class="input mono" v-model="data.password" />
        </Field>
      </div>
      <div class="grid2">
        <Network :data="data" />
        <Field label="UDP Relay Mode">
          <Select v-model="udpRelayMode">
            <option value="">{{ $t('none') }}</option>
            <option v-for="m in ['native', 'quic']" :key="m" :value="m">{{ m }}</option>
          </Select>
        </Field>
      </div>
      <div style="display: flex; gap: 24px; flex-wrap: wrap; margin-bottom: 15px;">
        <SwitchLabel v-model="udpOverStream" label="UDP Over Stream" />
      </div>
    </template>
    <div class="grid2">
      <Field :label="$t('types.tuic.congControl')">
        <Select v-model="data.congestion_control">
          <option v-for="c in congestion_controls" :key="c" :value="c">{{ c }}</option>
        </Select>
      </Field>
      <Field v-if="direction == 'in'" :label="$t('types.tuic.authTimeout') + ' (' + $t('date.s') + ')'">
        <input class="input mono" type="number" min="1" v-model.number="auth_timeout" />
      </Field>
    </div>
    <div style="display: flex; gap: 24px; flex-wrap: wrap; margin-bottom: 15px;">
      <SwitchLabel v-model="zeroRtt" label="Zero-RTT Handshake" />
    </div>
    <Field :label="$t('types.tuic.hb') + ' (' + $t('date.s') + ')'">
      <input class="input mono" type="number" min="1" v-model.number="heartbeat" />
    </Field>
  </div>
</template>

<script lang="ts" setup>
import Select from '@/components/ui/Select.vue'
import { computed } from 'vue'
import Field from '@/components/ui/Field.vue'
import SectionLabel from '@/components/ui/SectionLabel.vue'
import SwitchLabel from '@/components/ui/SwitchLabel.vue'
import Network from './Network.vue'

const props = defineProps<{ data: any; direction?: string }>()

const congestion_controls = ['cubic', 'new_reno', 'bbr']

const udpOverStream = computed({
  get: (): boolean => props.data.udp_over_stream ?? false,
  set: (v: boolean) => { props.data.udp_over_stream = v },
})

const zeroRtt = computed({
  get: (): boolean => props.data.zero_rtt_handshake ?? false,
  set: (v: boolean) => { props.data.zero_rtt_handshake = v },
})

const udpRelayMode = computed({
  get: (): string => props.data.udp_relay_mode ?? '',
  set: (v: string) => {
    if (v === '') delete props.data.udp_relay_mode
    else props.data.udp_relay_mode = v
  },
})

const auth_timeout = computed({
  get: (): any => (props.data.auth_timeout ? parseInt(props.data.auth_timeout.replace('s', '')) : ''),
  set: (v: number) => { props.data.auth_timeout = v ? v + 's' : '' },
})

const heartbeat = computed({
  get: (): any => (props.data.heartbeat ? parseInt(props.data.heartbeat.replace('s', '')) : ''),
  set: (v: number) => { props.data.heartbeat = v ? v + 's' : '' },
})
</script>
