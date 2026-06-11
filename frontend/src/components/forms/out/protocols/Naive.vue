<template>
  <div class="grid2" style="margin-bottom: 15px;">
    <Field :label="$t('types.un')" :mb="0">
      <input class="input mono" v-model="data.username" />
    </Field>
    <Field :label="$t('types.pw')" :mb="0">
      <input class="input mono" type="password" v-model="data.password" />
    </Field>
    <Field :label="$t('types.naive.insecureConcurrency')" :mb="0">
      <input class="input mono" type="number" min="0" v-model.number="insecureConcurrency" />
    </Field>
    <Field v-if="data.quic" :label="$t('types.naive.quicCongestion')" :mb="0">
      <select class="input" v-model="quicCongestion">
        <option value="">{{ $t('none') }}</option>
        <option v-for="c in outCngs" :key="c.value" :value="c.value">{{ c.title }}</option>
      </select>
    </Field>
  </div>
  <div style="display: flex; gap: 20px; flex-wrap: wrap; margin-bottom: 15px;">
    <SwitchLabel v-model="udpOverTcp" :label="$t('types.naive.udpOverTcp')" />
    <SwitchLabel v-model="quic" :label="$t('types.naive.quic')" />
  </div>
  <Headers :data="extraHeaders" />
</template>

<script lang="ts" setup>
import { computed } from 'vue'
import Field from '@/components/ui/Field.vue'
import SwitchLabel from '@/components/ui/SwitchLabel.vue'
import Headers from '../Headers.vue'

const props = defineProps<{ data: any }>()

const outCngs = [
  { title: 'BBR', value: 'bbr' },
  { title: 'BBR2', value: 'bbr2' },
  { title: 'Cubic', value: 'cubic' },
  { title: 'Reno', value: 'reno' },
]

const udpOverTcp = computed({
  get: (): boolean => {
    const d = props.data
    return d?.udp_over_tcp === true || (d?.udp_over_tcp && (d.udp_over_tcp as any)?.enabled)
  },
  set: (v: boolean) => {
    if (v) props.data.udp_over_tcp = { enabled: true }
    else props.data.udp_over_tcp = false
  },
})
const quic = computed({
  get: () => props.data.quic ?? false,
  set: (v: boolean) => { props.data.quic = v },
})
const quicCongestion = computed({
  get: () => props.data.quic_congestion_control ?? '',
  set: (v: string) => {
    if (v == '') delete props.data.quic_congestion_control
    else props.data.quic_congestion_control = v
  },
})
const insecureConcurrency = computed({
  get: (): number => props.data?.insecure_concurrency ?? 0,
  set: (v: number) => { props.data.insecure_concurrency = v > 0 ? v : undefined },
})
const extraHeaders = computed((): any => {
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
