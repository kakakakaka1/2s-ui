<template>
  <div class="grid2">
    <Network :data="data" />
    <Field :label="$t('types.direct.overrideAddr')" :mb="0">
      <input class="input mono" v-model="data.override_address" />
    </Field>
    <Field :label="$t('types.direct.overridePort')" :mb="0">
      <input class="input mono" type="number" min="0" v-model.number="overridePort" />
    </Field>
  </div>
  <div style="height: 15px;" />
</template>

<script lang="ts" setup>
import { computed } from 'vue'
import Field from '@/components/ui/Field.vue'
import Network from '../Network.vue'

const props = defineProps<{ data: any }>()

const overridePort = computed({
  get: () => props.data.override_port ? props.data.override_port : '',
  set: (v: any) => { props.data.override_port = !v || v == 0 ? undefined : parseInt(v) },
})
</script>
