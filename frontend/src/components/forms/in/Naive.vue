<template>
  <div>
    <SectionLabel v-if="direction != 'out_json'" style="margin-bottom: 12px;">Naive</SectionLabel>
    <!-- Inbound -->
    <template v-if="direction === 'in'">
      <div class="grid2">
        <Network :data="data" />
        <Field :label="$t('types.naive.quicCongestion')">
          <Select v-model="quicCongestion">
            <option value="">{{ $t('none') }}</option>
            <option v-for="c in inbCngs" :key="c.value" :value="c.value">{{ c.title }}</option>
          </Select>
        </Field>
      </div>
    </template>
    <!-- Client side (out_json) -->
    <template v-if="direction === 'out_json'">
      <Field :label="$t('types.naive.insecureConcurrency')">
        <input class="input mono" type="number" min="0" v-model.number="insecure_concurrency" />
      </Field>
      <div style="display: flex; gap: 24px; flex-wrap: wrap; margin-bottom: 15px;">
        <SwitchLabel v-model="udpOverTcp" :label="$t('types.naive.udpOverTcp')" />
      </div>
      <Headers :data="extra_headers" />
    </template>
  </div>
</template>

<script lang="ts" setup>
import Select from '@/components/ui/Select.vue'
import { computed } from 'vue'
import Field from '@/components/ui/Field.vue'
import SectionLabel from '@/components/ui/SectionLabel.vue'
import SwitchLabel from '@/components/ui/SwitchLabel.vue'
import Network from './Network.vue'
import Headers from './Headers.vue'

const props = defineProps<{ data: any; direction?: string }>()

const inbCngs = [
  { title: 'BBR', value: 'bbr' },
  { title: 'BBR Standard', value: 'bbr_standard' },
  { title: 'BBRv2', value: 'bbr2' },
  { title: 'BBRv2 variant', value: 'bbr2_variant' },
  { title: 'Cubic', value: 'cubic' },
  { title: 'New Reno', value: 'reno' },
]

const quicCongestion = computed({
  get: (): string => props.data.quic_congestion_control ?? '',
  set: (v: string) => {
    if (v === '') delete props.data.quic_congestion_control
    else props.data.quic_congestion_control = v
  },
})

const udpOverTcp = computed({
  get: (): boolean => {
    const d = props.data
    return d?.udp_over_tcp === true || (d?.udp_over_tcp && (d.udp_over_tcp as any)?.enabled)
  },
  set: (v: boolean) => {
    if (v) {
      props.data.udp_over_tcp = { enabled: true }
    } else {
      props.data.udp_over_tcp = false
    }
  },
})

const insecure_concurrency = computed({
  get: (): number => props.data?.insecure_concurrency ?? 0,
  set: (v: number) => {
    props.data.insecure_concurrency = v > 0 ? v : undefined
  },
})

const extra_headers = computed((): any => {
  const d = props.data
  return new Proxy({}, {
    get(_, prop) {
      if (prop === 'headers') return d?.extra_headers ?? {}
      return undefined
    },
    set(_, prop, value) {
      if (prop === 'headers') {
        d.extra_headers = value
        return true
      }
      return false
    },
  })
})
</script>
