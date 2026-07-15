<template>
  <div>
    <Field :label="$t('transport.grpcServiceName')">
      <input class="input mono" v-model="transport.service_name" />
    </Field>
    <div style="display: flex; gap: 24px; flex-wrap: wrap; margin-bottom: 15px;">
      <SwitchLabel v-model="transport.permit_without_stream" :label="$t('transport.grpcPws')" />
    </div>
    <div class="grid2">
      <Field :label="$t('transport.idleTimeout') + ' (' + $t('date.s') + ')'">
        <input class="input mono" type="number" min="1" v-model.number="idle_timeout" />
      </Field>
      <Field :label="$t('transport.pingTimeout') + ' (' + $t('date.s') + ')'">
        <input class="input mono" type="number" min="1" v-model.number="ping_timeout" />
      </Field>
    </div>
  </div>
</template>

<script lang="ts" setup>
import { computed } from 'vue'
import Field from '@/components/ui/Field.vue'
import SwitchLabel from '@/components/ui/SwitchLabel.vue'

const props = defineProps<{ transport: any }>()

const idle_timeout = computed({
  get: (): number | string => (props.transport.idle_timeout ? parseInt(props.transport.idle_timeout.replace('s', '')) : ''),
  set: (v: number) => { props.transport.idle_timeout = v ? v + 's' : '' },
})

const ping_timeout = computed({
  get: (): number | string => (props.transport.ping_timeout ? parseInt(props.transport.ping_timeout.replace('s', '')) : ''),
  set: (v: number) => { props.transport.ping_timeout = v ? v + 's' : '' },
})
</script>
