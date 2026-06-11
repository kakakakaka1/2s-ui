<template>
  <Field label="UDP over TCP">
    <select class="input" v-model="udpOverTcp">
      <option v-for="v in versions" :key="v.value" :value="v.value">{{ v.title }}</option>
    </select>
  </Field>
</template>

<script lang="ts" setup>
import { computed } from 'vue'
import { useI18n } from 'vue-i18n'
import Field from '@/components/ui/Field.vue'

const props = defineProps<{ data: any }>()
const { t } = useI18n({ useScope: 'global' })

const versions = computed(() => [
  { title: t('disable'), value: 0 },
  { title: '1', value: 1 },
  { title: '2', value: 2 },
])

const udpOverTcp = computed({
  get: (): number => props.data.udp_over_tcp?.version ?? 0,
  set: (v: number) => { props.data.udp_over_tcp = v > 0 ? { enabled: true, version: v } : undefined },
})
</script>
