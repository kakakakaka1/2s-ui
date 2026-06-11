<template>
  <div>
    <div class="grid2">
      <Field :label="$t('transport.hosts')"><input class="input mono" v-model="hosts" /></Field>
      <Field :label="$t('transport.path')"><input class="input mono" v-model="transport.path" /></Field>
    </div>
    <Field :label="$t('transport.httpMethod')">
      <select class="input" v-model="method">
        <option value="">{{ $t('none') }}</option>
        <option v-for="m in methodList" :key="m" :value="m">{{ m }}</option>
      </select>
    </Field>
    <div class="grid2">
      <Field :label="$t('transport.idleTimeout') + ' (' + $t('date.s') + ')'">
        <input class="input mono" type="number" min="1" v-model.number="idle_timeout" />
      </Field>
      <Field :label="$t('transport.pingTimeout') + ' (' + $t('date.s') + ')'">
        <input class="input mono" type="number" min="1" v-model.number="ping_timeout" />
      </Field>
    </div>
    <Headers :data="transport" />
  </div>
</template>

<script lang="ts" setup>
import { computed } from 'vue'
import Field from '@/components/ui/Field.vue'
import Headers from './Headers.vue'

const props = defineProps<{ transport: any }>()

const methodList = ['POST', 'GET', 'PUT', 'PATCH', 'DELETE']

const hosts = computed({
  get: (): string => (props.transport.host ? props.transport.host.join(',') : ''),
  set: (v: string) => { props.transport.host = v.length > 0 ? v.split(',') : [] },
})

const method = computed({
  get: (): string => props.transport.method ?? '',
  set: (v: string) => {
    if (v === '') delete props.transport.method
    else props.transport.method = v
  },
})

const idle_timeout = computed({
  get: (): number | string => (props.transport.idle_timeout ? parseInt(props.transport.idle_timeout.replace('s', '')) : ''),
  set: (v: number) => { props.transport.idle_timeout = v ? v + 's' : '' },
})

const ping_timeout = computed({
  get: (): number | string => (props.transport.ping_timeout ? parseInt(props.transport.ping_timeout.replace('s', '')) : ''),
  set: (v: number) => { props.transport.ping_timeout = v ? v + 's' : '' },
})
</script>
