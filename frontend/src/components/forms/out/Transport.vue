<template>
  <div style="margin-bottom: 15px;">
    <MSwitchRow v-model="tpEnable" :label="$t('objects.transport')" :desc="$t('transport.enable')" />
    <div v-if="tpEnable" class="fade-up">
      <Field :label="$t('type')">
        <Select v-model="transportType">
          <option v-for="(value, key) in trspTypes" :key="key" :value="value">{{ key }}</option>
        </Select>
      </Field>
      <TransportHttp v-if="transport.type == trspTypes.HTTP" :transport="transport" />
      <TransportWs v-if="transport.type == trspTypes.WebSocket" :transport="transport" />
      <TransportGrpc v-if="transport.type == trspTypes.gRPC" :transport="transport" />
      <TransportHttpUpgrade v-if="transport.type == trspTypes.HTTPUpgrade" :transport="transport" />
    </div>
  </div>
</template>

<script lang="ts" setup>
import Select from '@/components/ui/Select.vue'
import { computed } from 'vue'
import { TrspTypes, Transport } from '@/types/transport'
import Field from '@/components/ui/Field.vue'
import MSwitchRow from '@/components/ui/MSwitchRow.vue'
import TransportHttp from './transports/Http.vue'
import TransportWs from './transports/WebSocket.vue'
import TransportGrpc from './transports/Grpc.vue'
import TransportHttpUpgrade from './transports/HttpUpgrade.vue'

const props = defineProps<{ data: any }>()

const trspTypes = TrspTypes

const transport = computed((): Transport => props.data.transport)

const tpEnable = computed({
  get: () => Object.hasOwn(props.data.transport, 'type'),
  set: (v: boolean) => { props.data.transport = v ? { type: 'http' } : {} },
})
const transportType = computed({
  get: () => transport.value.type,
  set: (v: string) => { props.data.transport = { type: v } },
})
</script>
