<template>
  <div class="grid2">
    <Field :label="$t('transport.hosts') + ' ' + $t('commaSeparated')">
      <input class="input mono" v-model="hosts" />
    </Field>
    <Field :label="$t('transport.path')">
      <input class="input mono" v-model="transport.path" />
    </Field>
    <Field :label="$t('transport.httpMethod')">
      <Select v-model="method">
        <option value="">{{ $t('none') }}</option>
        <option v-for="m in methodList" :key="m" :value="m">{{ m }}</option>
      </Select>
    </Field>
    <Field :label="$t('transport.idleTimeout') + ' (' + $t('date.s') + ')'">
      <input class="input mono" type="number" min="1" v-model.number="idleTimeout" />
    </Field>
    <Field :label="$t('transport.pingTimeout') + ' (' + $t('date.s') + ')'">
      <input class="input mono" type="number" min="1" v-model.number="pingTimeout" />
    </Field>
  </div>
  <Headers :data="transport" />
</template>

<script lang="ts" setup>
import Select from '@/components/ui/Select.vue'
import { computed } from 'vue'
import { HTTP } from '@/types/transport'
import Field from '@/components/ui/Field.vue'
import Headers from '../Headers.vue'

const props = defineProps<{ transport: any }>()

const methodList = ['POST', 'GET', 'PUT', 'PATCH', 'DELETE']

const http = computed((): HTTP => props.transport ?? {})

const hosts = computed({
  get: () => http.value.host ? http.value.host.join(',') : '',
  set: (v: string) => { props.transport.host = v.length > 0 ? v.split(',') : [] },
})
const method = computed({
  get: () => http.value.method ?? '',
  set: (v: string) => {
    if (v == '') delete props.transport.method
    else props.transport.method = v
  },
})
const idleTimeout = computed({
  get: () => http.value.idle_timeout ? parseInt(http.value.idle_timeout.replace('s', '')) : '',
  set: (v: number) => { props.transport.idle_timeout = v ? v + 's' : '' },
})
const pingTimeout = computed({
  get: () => http.value.ping_timeout ? parseInt(http.value.ping_timeout.replace('s', '')) : '',
  set: (v: number) => { props.transport.ping_timeout = v ? v + 's' : '' },
})
</script>
