<template>
  <div class="grid2">
    <Field :label="$t('transport.path')">
      <input class="input mono" v-model="transport.path" />
    </Field>
    <Field :label="$t('transport.host')">
      <input class="input mono" v-model="host" />
    </Field>
    <Field label="Max Early Data">
      <input class="input mono" type="number" min="0" v-model.number="maxEarlyData" />
    </Field>
    <Field label="Early Data Header Name">
      <input class="input mono" v-model="transport.early_data_header_name" />
    </Field>
  </div>
  <Headers :data="transport" />
</template>

<script lang="ts" setup>
import { computed, onMounted } from 'vue'
import { WebSocket } from '@/types/transport'
import Field from '@/components/ui/Field.vue'
import Headers from '../Headers.vue'

const props = defineProps<{ transport: any }>()

const ws = computed((): WebSocket => props.transport)

const maxEarlyData = computed({
  get: () => ws.value.max_early_data ? ws.value.max_early_data : '',
  set: (v: number) => { props.transport.max_early_data = v != 0 ? v : undefined },
})
const host = computed({
  get: () => ws.value.headers?.Host ? ws.value.headers.Host : '',
  set: (v: string) => { props.transport.headers = v != '' ? { Host: v } : undefined },
})

onMounted(() => {
  props.transport.early_data_header_name ??= 'Sec-WebSocket-Protocol'
  props.transport.path ??= '/'
})
</script>
