<template>
  <div>
    <div class="grid2">
      <Field :label="$t('transport.path')"><input class="input mono" v-model="transport.path" /></Field>
      <Field :label="$t('transport.host')"><input class="input mono" v-model="host" /></Field>
    </div>
    <div class="grid2">
      <Field label="Max Early Data">
        <input class="input mono" type="number" min="0" v-model.number="max_early_data" />
      </Field>
      <Field label="Early Data Header Name">
        <input class="input mono" v-model="transport.early_data_header_name" />
      </Field>
    </div>
    <Headers :data="transport" />
  </div>
</template>

<script lang="ts" setup>
import { computed, onMounted } from 'vue'
import Field from '@/components/ui/Field.vue'
import Headers from './Headers.vue'

const props = defineProps<{ transport: any }>()

const max_early_data = computed({
  get: (): number | string => (props.transport.max_early_data ? props.transport.max_early_data : ''),
  set: (v: number) => { props.transport.max_early_data = v != 0 ? v : undefined },
})

const host = computed({
  get: (): string => (props.transport.headers?.Host ? props.transport.headers.Host : ''),
  set: (v: string) => {
    props.transport.headers = v != '' ? { Host: v } : undefined
  },
})

onMounted(() => {
  props.transport.early_data_header_name ??= 'Sec-WebSocket-Protocol'
  props.transport.path ??= '/'
})
</script>
