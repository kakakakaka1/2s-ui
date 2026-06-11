<template>
  <Field :label="$t('network')">
    <Select v-model="network">
      <option v-for="n in networks" :key="n.value" :value="n.value">{{ n.title }}</option>
    </Select>
  </Field>
</template>

<script lang="ts" setup>
import Select from '@/components/ui/Select.vue'
import { computed } from 'vue'
import Field from '@/components/ui/Field.vue'

const props = defineProps<{ data: any }>()

const networks = [
  { title: 'TCP/UDP', value: '' },
  { title: 'TCP', value: 'tcp' },
  { title: 'UDP', value: 'udp' },
]

const network = computed({
  get: (): string => props.data.network ?? '',
  set: (v: string) => { props.data.network = v != '' ? v : undefined },
})
</script>
