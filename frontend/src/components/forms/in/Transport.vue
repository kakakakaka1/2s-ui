<template>
  <div>
    <MSwitchRow v-model="tpEnable" :label="$t('objects.transport')" :desc="$t('transport.enable')" />
    <div v-if="tpEnable" class="fade-up">
      <Field :label="$t('type')">
        <select class="input" v-model="transportType">
          <option v-for="(value, key) in trspTypes" :key="value" :value="value">{{ key }}</option>
        </select>
      </Field>
      <TransportHttp v-if="trsp.type == trspTypes.HTTP" :transport="trsp" />
      <TransportWs v-if="trsp.type == trspTypes.WebSocket" :transport="trsp" />
      <TransportGrpc v-if="trsp.type == trspTypes.gRPC" :transport="trsp" />
      <TransportHttpUpgrade v-if="trsp.type == trspTypes.HTTPUpgrade" :transport="trsp" />
    </div>
  </div>
</template>

<script lang="ts" setup>
import { computed } from 'vue'
import { TrspTypes } from '@/types/transport'
import Field from '@/components/ui/Field.vue'
import MSwitchRow from '@/components/ui/MSwitchRow.vue'
import TransportHttp from './TransportHttp.vue'
import TransportWs from './TransportWs.vue'
import TransportGrpc from './TransportGrpc.vue'
import TransportHttpUpgrade from './TransportHttpUpgrade.vue'

const props = defineProps<{ data: any }>()

const trspTypes = TrspTypes

const trsp = computed<any>(() => props.data.transport)

const tpEnable = computed({
  get: (): boolean => Object.hasOwn(props.data.transport, 'type'),
  set: (v: boolean) => { props.data.transport = v ? { type: 'http' } : {} },
})

const transportType = computed({
  get: (): string => trsp.value.type,
  set: (v: string) => { props.data.transport = { type: v } },
})
</script>
