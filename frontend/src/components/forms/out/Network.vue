<template>
  <Field :label="$t('network')" :mb="mb">
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

// Defaults to 0 because every out/protocols caller sits in a grid2 that owns
// the spacing; standalone callers pass 15 to get Field's usual gap.
const props = withDefaults(defineProps<{ data: any; mb?: number }>(), { mb: 0 })

const network = computed({
  get: (): string => props.data.network ?? '',
  set: (v: string) => { props.data.network = v != '' ? v : undefined },
})
</script>
