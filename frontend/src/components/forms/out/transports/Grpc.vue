<template>
  <div class="grid2">
    <Field :label="$t('transport.grpcServiceName')">
      <input class="input mono" v-model="transport.service_name" />
    </Field>
    <Field :label="$t('transport.idleTimeout') + ' (' + $t('date.s') + ')'">
      <input class="input mono" type="number" min="1" v-model.number="idleTimeout" />
    </Field>
    <Field :label="$t('transport.pingTimeout') + ' (' + $t('date.s') + ')'">
      <input class="input mono" type="number" min="1" v-model.number="pingTimeout" />
    </Field>
  </div>
  <div style="margin-bottom: 14px;">
    <SwitchLabel v-model="permitWithoutStream" :label="$t('transport.grpcPws')" />
  </div>
</template>

<script lang="ts" setup>
import { computed } from 'vue'
import { gRPC } from '@/types/transport'
import Field from '@/components/ui/Field.vue'
import SwitchLabel from '@/components/ui/SwitchLabel.vue'

const props = defineProps<{ transport: any }>()

const grpc = computed((): gRPC => props.transport ?? {})

const permitWithoutStream = computed({
  get: () => props.transport.permit_without_stream ?? false,
  set: (v: boolean) => { props.transport.permit_without_stream = v },
})
const idleTimeout = computed({
  get: () => grpc.value.idle_timeout ? parseInt(grpc.value.idle_timeout.replace('s', '')) : '',
  set: (v: number) => { props.transport.idle_timeout = v ? v + 's' : '' },
})
const pingTimeout = computed({
  get: () => grpc.value.ping_timeout ? parseInt(grpc.value.ping_timeout.replace('s', '')) : '',
  set: (v: number) => { props.transport.ping_timeout = v ? v + 's' : '' },
})
</script>
