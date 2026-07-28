<template>
  <Field label="UDP over TCP" :mb="mb">
    <Select v-model.number="udpOverTcp">
      <option :value="0">{{ $t('disable') }}</option>
      <option :value="1">1</option>
      <option :value="2">2</option>
    </Select>
  </Field>
</template>

<script lang="ts" setup>
import Select from '@/components/ui/Select.vue'
import { computed } from 'vue'
import Field from '@/components/ui/Field.vue'

// See Network.vue: 0 for the grid2 callers, 15 for standalone ones.
const props = withDefaults(defineProps<{ data: any; mb?: number }>(), { mb: 0 })

const udpOverTcp = computed({
  get: (): number => props.data.udp_over_tcp?.version ?? 0,
  set: (v: number) => { props.data.udp_over_tcp = v > 0 ? { enabled: true, version: v } : undefined },
})
</script>
