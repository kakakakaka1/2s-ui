<template>
  <div>
    <SectionLabel style="margin-bottom: 12px;">Direct</SectionLabel>
    <div class="grid2">
      <Network :data="data" />
      <Field :label="$t('types.direct.overrideAddr')">
        <input class="input mono" v-model="data.override_address" />
      </Field>
    </div>
    <Field :label="$t('types.direct.overridePort')">
      <input class="input mono" type="number" min="0" v-model.number="override_port" />
    </Field>
  </div>
</template>

<script lang="ts" setup>
import { computed } from 'vue'
import Field from '@/components/ui/Field.vue'
import SectionLabel from '@/components/ui/SectionLabel.vue'
import Network from './Network.vue'

const props = defineProps<{ data: any }>()

const override_port = computed({
  get: (): any => (props.data.override_port ? props.data.override_port : ''),
  set: (v: any) => {
    props.data.override_port = v?.length == 0 || v == 0 ? undefined : parseInt(v)
  },
})
</script>
