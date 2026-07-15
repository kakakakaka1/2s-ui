<template>
  <Field :label="$t('ui.uuid')">
    <KeyInput v-model="data.uuid" @regenerate="data.uuid = RandomUtil.randomUUID()" />
  </Field>
  <div class="grid2" style="margin-bottom: 15px;">
    <Field label="Alter ID" :mb="0">
      <input class="input mono" type="number" min="0" v-model.number="data.alter_id" />
    </Field>
    <Field :label="$t('types.vmess.security')" :mb="0">
      <Select v-model="data.security">
        <option v-for="s in securities" :key="s" :value="s">{{ s }}</option>
      </Select>
    </Field>
    <Field :label="$t('types.vless.udpEnc')" :mb="0">
      <Select v-model="packetEncoding">
        <option value="none">none</option>
        <option value="packetaddr">packetaddr</option>
        <option value="xudp">xudp</option>
      </Select>
    </Field>
    <Network :data="data" />
  </div>
  <div style="display: flex; gap: 20px; flex-wrap: wrap; margin-bottom: 15px;">
    <SwitchLabel v-model="globalPadding" :label="$t('types.vmess.globalPadding')" />
    <SwitchLabel v-model="authenticatedLength" :label="$t('types.vmess.authLen')" />
  </div>
</template>

<script lang="ts" setup>
import Select from '@/components/ui/Select.vue'
import { computed } from 'vue'
import RandomUtil from '@/plugins/randomUtil'
import Field from '@/components/ui/Field.vue'
import KeyInput from '@/components/ui/KeyInput.vue'
import SwitchLabel from '@/components/ui/SwitchLabel.vue'
import Network from '../Network.vue'

const props = defineProps<{ data: any }>()

const securities = ['auto', 'none', 'zero', 'aes-128-gcm', 'aes-128-ctr', 'chacha20-poly1305']

const packetEncoding = computed({
  get: () => props.data.packet_encoding != undefined ? props.data.packet_encoding : 'none',
  set: (v: string) => { props.data.packet_encoding = v != 'none' ? v : undefined },
})
const globalPadding = computed({
  get: () => props.data.global_padding ?? false,
  set: (v: boolean) => { props.data.global_padding = v },
})
const authenticatedLength = computed({
  get: () => props.data.authenticated_length ?? false,
  set: (v: boolean) => { props.data.authenticated_length = v },
})
</script>
