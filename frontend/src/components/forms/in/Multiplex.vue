<template>
  <div>
    <MSwitchRow v-model="muxEnable" :label="$t('objects.multiplex')" :desc="$t('mux.enable')" />
    <div v-if="muxEnable" class="fade-up">
      <template v-if="direction == 'out'">
        <Field :label="$t('protocol')">
          <select class="input" v-model="muxProtocol">
            <option value="">{{ $t('none') }}</option>
            <option v-for="p in ['smux', 'yamux', 'h2mux']" :key="p" :value="p">{{ p }}</option>
          </select>
        </Field>
        <div class="grid2">
          <Field :label="$t('mux.maxConn')">
            <input class="input mono" type="number" min="0" v-model.number="max_connections" />
          </Field>
          <Field :label="$t('mux.minStr')">
            <input class="input mono" type="number" min="0" v-model.number="min_streams" />
          </Field>
        </div>
        <Field :label="$t('mux.maxStr')">
          <input class="input mono" type="number" :min="min_streams" v-model.number="max_streams" />
        </Field>
      </template>
      <div style="display: flex; gap: 24px; flex-wrap: wrap; margin-bottom: 15px;">
        <SwitchLabel v-model="padding" :label="$t('mux.padding')" />
        <SwitchLabel v-model="brutalEnable" :label="$t('mux.enableBrutal')" />
      </div>
      <div v-if="mux?.brutal?.enabled" class="grid2 fade-up">
        <Field :label="$t('stats.upload') + ' (' + $t('stats.Mbps') + ')'">
          <input class="input mono" type="number" v-model.number="up_mbps" />
        </Field>
        <Field :label="$t('stats.download') + ' (' + $t('stats.Mbps') + ')'">
          <input class="input mono" type="number" min="0" v-model.number="down_mbps" />
        </Field>
      </div>
    </div>
  </div>
</template>

<script lang="ts" setup>
import { computed } from 'vue'
import { oMultiplex } from '@/types/multiplex'
import Field from '@/components/ui/Field.vue'
import MSwitchRow from '@/components/ui/MSwitchRow.vue'
import SwitchLabel from '@/components/ui/SwitchLabel.vue'

const props = defineProps<{ data: any; direction?: string }>()

const mux = computed((): oMultiplex => <oMultiplex>props.data.multiplex ?? null)

const muxEnable = computed({
  get: (): boolean => (mux.value ? mux.value.enabled : false),
  set: (v: boolean) => { props.data.multiplex = v ? { enabled: v } : undefined },
})

const muxProtocol = computed({
  get: (): string => mux.value?.protocol ?? '',
  set: (v: string) => {
    if (v === '') delete mux.value?.protocol
    else mux.value.protocol = <any>v
  },
})

const max_connections = computed({
  get: (): number => (mux.value?.max_connections ? mux.value.max_connections : 0),
  set: (v: number) => { mux.value.max_connections = v > 0 ? v : undefined },
})

const min_streams = computed({
  get: (): number => (mux.value?.min_streams ? mux.value.min_streams : 0),
  set: (v: number) => { mux.value.min_streams = v > 0 ? v : undefined },
})

const max_streams = computed({
  get: (): number => (mux.value?.max_streams ? mux.value.max_streams : 0),
  set: (v: number) => { mux.value.max_streams = v > 0 ? v : undefined },
})

const padding = computed({
  get: (): boolean => (mux.value?.padding ? mux.value.padding : false),
  set: (v: boolean) => { mux.value.padding = v ? true : undefined },
})

const brutalEnable = computed({
  get: (): boolean => (mux.value?.brutal ? mux.value.brutal.enabled : false),
  set: (v: boolean) => { mux.value.brutal = v ? { enabled: v, up_mbps: 100, down_mbps: 100 } : undefined },
})

const down_mbps = computed({
  get: (): any => (mux.value?.brutal && mux.value.brutal.down_mbps ? mux.value.brutal.down_mbps : 0),
  set: (v: any) => {
    if (mux.value.brutal) {
      mux.value.brutal.down_mbps = v?.length != 0 ? v : 0
    }
  },
})

const up_mbps = computed({
  get: (): any => (mux.value?.brutal && mux.value.brutal.up_mbps ? mux.value.brutal.up_mbps : 0),
  set: (v: any) => {
    if (mux.value.brutal) {
      mux.value.brutal.up_mbps = v?.length != 0 ? v : 0
    }
  },
})
</script>
