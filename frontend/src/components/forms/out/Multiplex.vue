<template>
  <div style="margin-bottom: 15px;">
    <MSwitchRow v-model="muxEnable" :label="$t('objects.multiplex')" :desc="$t('mux.enable')" />
    <div v-if="muxEnable" class="fade-up">
      <div class="grid2">
        <Field :label="$t('protocol')">
          <select class="input" v-model="protocol">
            <option value="">{{ $t('none') }}</option>
            <option value="smux">smux</option>
            <option value="yamux">yamux</option>
            <option value="h2mux">h2mux</option>
          </select>
        </Field>
        <Field :label="$t('mux.maxConn')">
          <input class="input mono" type="number" min="0" v-model.number="maxConnections" />
        </Field>
        <Field :label="$t('mux.minStr')">
          <input class="input mono" type="number" min="0" v-model.number="minStreams" />
        </Field>
        <Field :label="$t('mux.maxStr')">
          <input class="input mono" type="number" :min="minStreams" v-model.number="maxStreams" />
        </Field>
      </div>
      <div style="display: flex; gap: 20px; flex-wrap: wrap; margin-bottom: 14px;">
        <SwitchLabel v-model="padding" :label="$t('mux.padding')" />
        <SwitchLabel v-model="brutalEnable" :label="$t('mux.enableBrutal')" />
      </div>
      <div v-if="mux?.brutal?.enabled" class="grid2 fade-up">
        <Field :label="$t('stats.upload') + ' (' + $t('stats.Mbps') + ')'">
          <input class="input mono" type="number" min="0" v-model.number="upMbps" />
        </Field>
        <Field :label="$t('stats.download') + ' (' + $t('stats.Mbps') + ')'">
          <input class="input mono" type="number" min="0" v-model.number="downMbps" />
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

const props = defineProps<{ data: any }>()

const mux = computed((): oMultiplex => props.data.multiplex ?? null)

const muxEnable = computed({
  get: (): boolean => mux.value ? mux.value.enabled : false,
  set: (v: boolean) => { props.data.multiplex = v ? { enabled: v } : undefined },
})
const protocol = computed({
  get: (): string => mux.value?.protocol ?? '',
  set: (v: string) => {
    if (v == '') delete props.data.multiplex.protocol
    else props.data.multiplex.protocol = v
  },
})
const maxConnections = computed({
  get: (): number => mux.value?.max_connections ? mux.value.max_connections : 0,
  set: (v: number) => { mux.value.max_connections = v > 0 ? v : undefined },
})
const minStreams = computed({
  get: (): number => mux.value?.min_streams ? mux.value.min_streams : 0,
  set: (v: number) => { mux.value.min_streams = v > 0 ? v : undefined },
})
const maxStreams = computed({
  get: (): number => mux.value?.max_streams ? mux.value.max_streams : 0,
  set: (v: number) => { mux.value.max_streams = v > 0 ? v : undefined },
})
const padding = computed({
  get: (): boolean => mux.value?.padding ? mux.value.padding : false,
  set: (v: boolean) => { mux.value.padding = v ? true : undefined },
})
const brutalEnable = computed({
  get: (): boolean => mux.value?.brutal ? mux.value.brutal.enabled : false,
  set: (v: boolean) => { mux.value.brutal = v ? { enabled: v, up_mbps: 100, down_mbps: 100 } : undefined },
})
const downMbps = computed({
  get: () => mux.value?.brutal && mux.value.brutal.down_mbps ? mux.value.brutal.down_mbps : 0,
  set: (v: any) => {
    if (mux.value.brutal) {
      mux.value.brutal.down_mbps = v && String(v).length != 0 ? v : 0
    }
  },
})
const upMbps = computed({
  get: () => mux.value?.brutal && mux.value.brutal.up_mbps ? mux.value.brutal.up_mbps : 0,
  set: (v: any) => {
    if (mux.value.brutal) {
      mux.value.brutal.up_mbps = v && String(v).length != 0 ? v : 0
    }
  },
})
</script>
