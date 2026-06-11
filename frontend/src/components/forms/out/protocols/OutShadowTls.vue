<template>
  <div class="grid2" style="margin-bottom: 15px;">
    <Field :label="$t('version')" :mb="0">
      <select class="input" v-model.number="version">
        <option :value="1">1</option>
        <option :value="2">2</option>
        <option :value="3">3</option>
      </select>
    </Field>
    <Field v-if="data.version > 1" :label="$t('types.pw')" :mb="0">
      <input class="input mono" v-model="data.password" />
    </Field>
  </div>
</template>

<script lang="ts" setup>
import { computed } from 'vue'
import Field from '@/components/ui/Field.vue'

const props = defineProps<{ data: any }>()

const version = computed({
  get: () => props.data.version ?? 3,
  set: (v: number) => {
    props.data.version = v
    if (v == 1) {
      delete props.data.password
    } else if (props.data.password === undefined) {
      props.data.password = ''
    }
  },
})
</script>
