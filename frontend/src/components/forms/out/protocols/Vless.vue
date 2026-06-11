<template>
  <Field :label="$t('ui.uuid')">
    <KeyInput v-model="data.uuid" @regenerate="data.uuid = RandomUtil.randomUUID()" />
  </Field>
  <div class="grid2" style="margin-bottom: 15px;">
    <Field :label="$t('types.vless.flow')" :mb="0">
      <select class="input" v-model="data.flow">
        <option value="">{{ $t('none') }}</option>
        <option value="xtls-rprx-vision">xtls-rprx-vision</option>
      </select>
    </Field>
    <Field :label="$t('types.vless.udpEnc')" :mb="0">
      <select class="input" v-model="packetEncoding">
        <option value="none">none</option>
        <option value="packetaddr">packetaddr</option>
        <option value="xudp">xudp</option>
      </select>
    </Field>
    <Network :data="data" />
  </div>
</template>

<script lang="ts" setup>
import { computed } from 'vue'
import RandomUtil from '@/plugins/randomUtil'
import Field from '@/components/ui/Field.vue'
import KeyInput from '@/components/ui/KeyInput.vue'
import Network from '../Network.vue'

const props = defineProps<{ data: any }>()

const packetEncoding = computed({
  get: () => props.data.packet_encoding != undefined ? props.data.packet_encoding : 'none',
  set: (v: string) => { props.data.packet_encoding = v != 'none' ? v : undefined },
})
</script>
