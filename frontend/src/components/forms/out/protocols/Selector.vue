<template>
  <Field :label="$t('pages.outbounds')">
    <MultiPick v-model="outboundsValue" :options="tags ?? []" />
  </Field>
  <Field :label="$t('types.lb.defaultOut')">
    <select class="input" v-model="defaultOut">
      <option value="">{{ $t('none') }}</option>
      <option v-for="o in data.outbounds ?? []" :key="o" :value="o">{{ o }}</option>
    </select>
  </Field>
  <div style="margin-bottom: 15px;">
    <CheckLabel v-model="interrupt" :label="$t('types.lb.interruptConn')" />
  </div>
</template>

<script lang="ts" setup>
import { computed } from 'vue'
import Field from '@/components/ui/Field.vue'
import MultiPick from '@/components/ui/MultiPick.vue'
import CheckLabel from '@/components/ui/CheckLabel.vue'

const props = defineProps<{ data: any; tags?: string[] }>()

const outboundsValue = computed({
  get: (): string[] => props.data.outbounds ?? [],
  set: (v: string[]) => {
    props.data.outbounds = v
    if (!props.data.outbounds?.includes(props.data.default)) {
      delete props.data.default
    }
  },
})
const defaultOut = computed({
  get: () => props.data.default ?? '',
  set: (v: string) => {
    if (v == '') delete props.data.default
    else props.data.default = v
  },
})
const interrupt = computed({
  get: () => props.data.interrupt_exist_connections ?? false,
  set: (v: boolean) => { props.data.interrupt_exist_connections = v },
})
</script>
