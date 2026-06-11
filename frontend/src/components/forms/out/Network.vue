<template>
  <Field :label="$t('network')" :mb="0">
    <Select v-model="network">
      <option value="">TCP/UDP</option>
      <option value="tcp">TCP</option>
      <option value="udp">UDP</option>
    </Select>
  </Field>
</template>

<script lang="ts" setup>
import Select from '@/components/ui/Select.vue'
import { computed } from 'vue'
import Field from '@/components/ui/Field.vue'

const props = defineProps<{ data: any }>()

const network = computed({
  get: (): string => props.data.network ?? '',
  set: (v: string) => { props.data.network = v != '' ? v : undefined },
})
</script>
