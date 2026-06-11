<template>
  <div>
    <MSwitchRow v-model="tlsOn" :label="$t('ui.tlsReality')" :desc="$t('ui.encryptListener')" />
    <div v-if="tlsOn" class="fade-up">
      <Field :label="$t('template')">
        <select class="input" v-model="inbound.tls_id">
          <option v-for="it in tlsItems" :key="it.value" :value="it.value">{{ it.title }}</option>
        </select>
      </Field>
    </div>
  </div>
</template>

<script lang="ts" setup>
import { computed } from 'vue'
import { useI18n } from 'vue-i18n'
import Field from '@/components/ui/Field.vue'
import MSwitchRow from '@/components/ui/MSwitchRow.vue'

const props = defineProps<{ inbound: any; tlsConfigs?: any[] }>()
const { t } = useI18n({ useScope: 'global' })

const tlsItems = computed((): any[] => [
  { title: t('none'), value: 0 },
  ...(props.tlsConfigs ?? []).map((c: any) => ({ title: c.name, value: c.id })),
])

const tlsOn = computed({
  get: (): boolean => (props.inbound.tls_id ?? 0) > 0,
  set: (v: boolean) => {
    props.inbound.tls_id = v ? props.tlsConfigs?.[0]?.id ?? 0 : 0
  },
})
</script>
