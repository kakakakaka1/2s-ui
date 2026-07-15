<template>
  <div style="display: flex; align-items: center; gap: 10px; margin-bottom: 12px;">
    <SectionLabel>Hysteria</SectionLabel>
    <div style="flex: 1;" />
    <Pop :min-width="200" direction="down">
      <template #trigger="{ toggle }">
        <Btn variant="subtle" sm @click="toggle">{{ $t('types.hy.hyOptions') }}</Btn>
      </template>
      <div style="display: flex; flex-direction: column; gap: 2px; padding: 4px;">
        <div class="pop-item"><SwitchLabel v-model="optionRsvConn" label="Recv window conn" /></div>
        <div class="pop-item"><SwitchLabel v-model="optionRsvWin" label="Recv window" /></div>
      </div>
    </Pop>
  </div>
  <div class="grid2" style="margin-bottom: 15px;">
    <Field :label="$t('stats.upload') + ' (' + $t('stats.Mbps') + ')'" :mb="0">
      <input class="input mono" type="number" min="0" v-model.number="upMbps" />
    </Field>
    <Field :label="$t('stats.download') + ' (' + $t('stats.Mbps') + ')'" :mb="0">
      <input class="input mono" type="number" min="0" v-model.number="downMbps" />
    </Field>
    <Field :label="$t('types.hy.obfs')" :mb="0">
      <input class="input mono" v-model="data.obfs" />
    </Field>
    <Field :label="$t('types.hy.auth')" :mb="0">
      <input class="input mono" v-model="data.auth_str" />
    </Field>
    <Network :data="data" />
    <Field v-if="data.recv_window_conn != undefined" label="Recv window conn" :mb="0">
      <input class="input mono" type="number" min="0" v-model.number="data.recv_window_conn" />
    </Field>
    <Field v-if="data.recv_window != undefined" label="Recv window" :mb="0">
      <input class="input mono" type="number" min="0" v-model.number="data.recv_window" />
    </Field>
  </div>
  <div style="margin-bottom: 15px;">
    <SwitchLabel v-model="disableMtu" label="Disable MTU discovery" />
  </div>
</template>

<script lang="ts" setup>
import { computed } from 'vue'
import Field from '@/components/ui/Field.vue'
import Btn from '@/components/ui/Btn.vue'
import Pop from '@/components/ui/Pop.vue'
import SwitchLabel from '@/components/ui/SwitchLabel.vue'
import SectionLabel from '@/components/ui/SectionLabel.vue'
import Network from '../Network.vue'

const props = defineProps<{ data: any }>()

const disableMtu = computed({
  get: () => props.data.disable_mtu_discovery ?? false,
  set: (v: boolean) => { props.data.disable_mtu_discovery = v },
})
const optionRsvConn = computed({
  get: (): boolean => props.data.recv_window_conn != undefined,
  set: (v: boolean) => { props.data.recv_window_conn = v ? 15728640 : undefined },
})
const optionRsvWin = computed({
  get: (): boolean => props.data.recv_window != undefined,
  set: (v: boolean) => { props.data.recv_window = v ? 67108864 : undefined },
})
const downMbps = computed({
  get: () => props.data.down_mbps ? props.data.down_mbps : 0,
  set: (v: any) => {
    if (v && String(v).length != 0) {
      props.data.down_mbps = v
      props.data.down = '' + v + ' Mbps'
    } else {
      props.data.down_mbps = 0
      props.data.down = '0 Mbps'
    }
  },
})
const upMbps = computed({
  get: () => props.data.up_mbps ? props.data.up_mbps : 0,
  set: (v: number) => { props.data.up_mbps = v > 0 ? v : 0 },
})
</script>
