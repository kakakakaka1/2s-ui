<template>
  <div>
    <div style="display: flex; align-items: center; gap: 10px; margin-bottom: 12px;">
      <SectionLabel style="flex: 1;">Hysteria</SectionLabel>
      <Pop :min-width="220">
        <template #trigger="{ toggle }">
          <button type="button" class="btn btn-ghost btn-sm" @click="toggle">
            <Ico name="settings" :size="14" /> {{ $t('types.hy.hyOptions') }}
          </button>
        </template>
        <div style="display: flex; flex-direction: column; gap: 2px; padding: 5px;">
          <SwitchLabel v-model="optionRsvConn" label="Recv window conn" />
          <SwitchLabel v-model="optionRsvClnt" label="Recv window client" />
          <SwitchLabel v-model="optionMaxConn" label="Max conn client" />
        </div>
      </Pop>
    </div>
    <div class="grid2">
      <Field :label="$t('stats.upload') + ' (' + $t('stats.Mbps') + ')'">
        <input class="input mono" type="number" v-model.number="up_mbps" />
      </Field>
      <Field :label="$t('stats.download') + ' (' + $t('stats.Mbps') + ')'">
        <input class="input mono" type="number" min="0" v-model.number="down_mbps" />
      </Field>
    </div>
    <Field :label="$t('types.hy.obfs')">
      <input class="input mono" v-model="data.obfs" />
    </Field>
    <div style="display: flex; gap: 24px; flex-wrap: wrap; margin-bottom: 15px;">
      <SwitchLabel v-model="data.disable_mtu_discovery" label="Disable MTU discovery" />
    </div>
    <div class="grid2">
      <Field v-if="data.recv_window_conn != undefined" label="Recv window conn">
        <input class="input mono" type="number" min="0" v-model.number="data.recv_window_conn" />
      </Field>
      <Field v-if="data.recv_window_client != undefined" label="Recv window client">
        <input class="input mono" type="number" min="0" v-model.number="data.recv_window_client" />
      </Field>
      <Field v-if="data.max_conn_client != undefined" label="Max conn client">
        <input class="input mono" type="number" min="0" v-model.number="data.max_conn_client" />
      </Field>
    </div>
  </div>
</template>

<script lang="ts" setup>
import { computed } from 'vue'
import Field from '@/components/ui/Field.vue'
import SectionLabel from '@/components/ui/SectionLabel.vue'
import SwitchLabel from '@/components/ui/SwitchLabel.vue'
import Pop from '@/components/ui/Pop.vue'
import Ico from '@/components/ui/Ico.vue'

const props = defineProps<{ data: any; direction?: string }>()

const optionRsvConn = computed({
  get: (): boolean => props.data.recv_window_conn != undefined,
  set: (v: boolean) => { props.data.recv_window_conn = v ? 15728640 : undefined },
})

const optionRsvClnt = computed({
  get: (): boolean => props.data.recv_window_client != undefined,
  set: (v: boolean) => { props.data.recv_window_client = v ? 67108864 : undefined },
})

const optionMaxConn = computed({
  get: (): boolean => props.data.max_conn_client != undefined,
  set: (v: boolean) => { props.data.max_conn_client = v ? 1024 : undefined },
})

const down_mbps = computed({
  get: (): any => (props.data.down_mbps ? props.data.down_mbps : 0),
  set: (v: any) => {
    if (v?.length != 0) {
      props.data.down_mbps = v
      props.data.down = '' + v + ' Mbps'
    } else {
      props.data.down_mbps = 0
      props.data.down = '0 Mbps'
    }
  },
})

const up_mbps = computed({
  get: (): number => (props.data.up_mbps ? props.data.up_mbps : 0),
  set: (v: number) => { props.data.up_mbps = v > 0 ? v : 0 },
})
</script>
