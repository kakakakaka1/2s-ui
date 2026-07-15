<template>
  <Field :label="$t('types.pw')">
    <input class="input mono" v-model="data.password" />
  </Field>
  <div class="grid2" style="margin-bottom: 15px;">
    <Field :label="$t('types.anytls.idleInterval') + ' (' + $t('date.s') + ')'" :mb="0">
      <input class="input mono" type="number" min="0" v-model.number="idleInterval" />
    </Field>
    <Field :label="$t('types.anytls.idleTimeout') + ' (' + $t('date.s') + ')'" :mb="0">
      <input class="input mono" type="number" min="0" v-model.number="idleTimeout" />
    </Field>
    <Field :label="$t('types.anytls.minIdle')" :mb="0">
      <input class="input mono" type="number" min="0" v-model.number="minIdle" />
    </Field>
  </div>
</template>

<script lang="ts" setup>
import { computed } from 'vue'
import Field from '@/components/ui/Field.vue'

const props = defineProps<{ data: any }>()

const idleInterval = computed({
  get: () => props.data.idle_session_check_interval?.length > 0 ? parseInt(props.data.idle_session_check_interval.replace('s', '')) : 30,
  set: (v: number) => { props.data.idle_session_check_interval = v && v >= 0 ? `${v}s` : undefined },
})
const idleTimeout = computed({
  get: () => props.data.idle_session_timeout?.length > 0 ? parseInt(props.data.idle_session_timeout.replace('s', '')) : 30,
  set: (v: number) => { props.data.idle_session_timeout = v && v >= 0 ? `${v}s` : undefined },
})
const minIdle = computed({
  get: () => props.data.min_idle_session != undefined ? props.data.min_idle_session : 0,
  set: (v: number) => { props.data.min_idle_session = v > 0 ? v : undefined },
})
</script>
