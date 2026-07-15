<template>
  <div style="display: flex; align-items: center; gap: 10px; margin-bottom: 12px;">
    <SectionLabel>Warp</SectionLabel>
    <div style="flex: 1;" />
    <Pop :min-width="190" direction="down">
      <template #trigger="{ toggle }">
        <Btn variant="subtle" sm @click="toggle">{{ $t('types.wg.options') }}</Btn>
      </template>
      <div style="display: flex; flex-direction: column; gap: 2px; padding: 4px;">
        <div class="pop-item"><SwitchLabel v-model="optionUdp" label="UDP Timeout" /></div>
        <div class="pop-item"><SwitchLabel v-model="optionWorker" :label="$t('types.wg.worker')" /></div>
        <div class="pop-item"><SwitchLabel v-model="optionMtu" label="MTU" /></div>
      </div>
    </Pop>
  </div>

  <template v-if="data.id > 0">
    <div style="display: flex; flex-direction: column; margin-bottom: 15px;" dir="ltr">
      <div class="kv-row">
        <span class="k">Device ID</span>
        <span class="v mono">{{ data.ext.device_id }}</span>
      </div>
      <div class="kv-row">
        <span class="k">Access Token</span>
        <span class="v mono">{{ data.ext.access_token }}</span>
      </div>
      <div class="kv-row">
        <span class="k">{{ $t('types.wg.privKey') }}</span>
        <span class="v mono">{{ data.private_key }}</span>
      </div>
      <div class="kv-row">
        <span class="k">{{ $t('types.wg.localIp') }}</span>
        <span class="v mono">{{ data.address.join(',') }}</span>
      </div>
    </div>
    <Field label="License Key">
      <input class="input mono" v-model="data.ext.license_key" />
    </Field>

    <div style="display: flex; align-items: center; margin: 4px 0 12px;">
      <SectionLabel>{{ $t('types.wg.peer') }}</SectionLabel>
    </div>
    <div class="grid2">
      <Field :label="$t('out.addr')">
        <input class="input mono" v-model="data.peers[0].address" />
      </Field>
      <Field :label="$t('out.port')">
        <input class="input mono" type="number" min="1" v-model.number="data.peers[0].port" />
      </Field>
    </div>
    <div style="display: flex; flex-direction: column; margin-bottom: 15px;" dir="ltr">
      <div class="kv-row">
        <span class="k">{{ $t('types.wg.pubKey') }}</span>
        <span class="v mono">{{ data.peers[0].public_key }}</span>
      </div>
      <div class="kv-row">
        <span class="k">{{ $t('types.wg.allowedIp') }}</span>
        <span class="v mono">{{ data.peers[0].allowed_ips.join(',') }}</span>
      </div>
      <div class="kv-row">
        <span class="k">Reserved</span>
        <span class="v mono">[{{ data.peers[0].reserved.join(',') }}]</span>
      </div>
    </div>
  </template>

  <div class="grid2" v-if="data.udp_timeout != undefined || data.workers != undefined || data.mtu != undefined">
    <Field v-if="data.udp_timeout != undefined" :label="'UDP Timeout (' + $t('date.m') + ')'">
      <input class="input mono" type="number" min="0" v-model.number="udpTimeout" />
    </Field>
    <Field v-if="data.workers != undefined" :label="$t('types.wg.worker')">
      <input class="input mono" type="number" min="1" v-model.number="data.workers" />
    </Field>
    <Field v-if="data.mtu != undefined" :label="$t('ui.mtu')">
      <input class="input mono" type="number" min="0" v-model.number="data.mtu" />
    </Field>
  </div>
  <div style="display: flex; gap: 20px; flex-wrap: wrap; margin-bottom: 15px;">
    <SwitchLabel v-model="system" :label="$t('types.wg.sysIf')" />
  </div>
  <Field v-if="data.system" :label="$t('types.wg.ifName')">
    <input class="input mono" v-model="ifName" />
  </Field>
</template>

<script lang="ts" setup>
import { computed } from 'vue'
import Field from '@/components/ui/Field.vue'
import Btn from '@/components/ui/Btn.vue'
import Pop from '@/components/ui/Pop.vue'
import SwitchLabel from '@/components/ui/SwitchLabel.vue'
import SectionLabel from '@/components/ui/SectionLabel.vue'

const props = defineProps<{ data: any }>()

const system = computed({
  get: () => props.data.system ?? false,
  set: (v: boolean) => { props.data.system = v },
})
const optionUdp = computed({
  get: (): boolean => props.data.udp_timeout != undefined,
  set: (v: boolean) => { props.data.udp_timeout = v ? '5m' : undefined },
})
const optionWorker = computed({
  get: (): boolean => props.data.workers != undefined,
  set: (v: boolean) => { props.data.workers = v ? 2 : undefined },
})
const optionMtu = computed({
  get: (): boolean => props.data.mtu != undefined,
  set: (v: boolean) => { props.data.mtu = v ? 1408 : undefined },
})
const ifName = computed({
  get: () => props.data.name ?? '',
  set: (v: string) => { props.data.name = v.length > 0 ? v : undefined },
})
const udpTimeout = computed({
  get: () => props.data.udp_timeout ? parseInt(props.data.udp_timeout.replace('m', '')) : 5,
  set: (v: number) => { props.data.udp_timeout = v > 0 ? v + 'm' : '5m' },
})
</script>
