<template>
  <Field :label="$t('objects.tls')">
    <Select v-model="inbound.tls_id">
      <option v-for="it in tlsItems" :key="it.value" :value="it.value">{{ it.title }}</option>
    </Select>
  </Field>
</template>

<script lang="ts" setup>
import Select from '@/components/ui/Select.vue'
import { computed } from 'vue'
import { useI18n } from 'vue-i18n'
import Field from '@/components/ui/Field.vue'

const props = defineProps<{ inbound: any; tlsConfigs?: any[] }>()
const { t } = useI18n({ useScope: 'global' })

const tlsItems = computed((): any[] => [
  { title: t('none'), value: 0 },
  ...(props.tlsConfigs ?? []).map((c: any) => ({ title: c.name, value: c.id })),
])
</script>
