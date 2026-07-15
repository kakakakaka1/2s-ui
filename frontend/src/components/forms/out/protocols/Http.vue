<template>
  <div class="grid2" style="margin-bottom: 15px;">
    <Field :label="$t('types.un')" :mb="0">
      <input class="input mono" v-model="username" />
    </Field>
    <Field :label="$t('types.pw')" :mb="0">
      <input class="input mono" type="password" v-model="password" />
    </Field>
    <Field :label="$t('transport.path')" :mb="0">
      <input class="input mono" v-model="data.path" />
    </Field>
  </div>
  <Headers :data="data" />
</template>

<script lang="ts" setup>
import { computed } from 'vue'
import Field from '@/components/ui/Field.vue'
import Headers from '../Headers.vue'

const props = defineProps<{ data: any }>()

const username = computed({
  get: (): string => props.data.username?.length > 0 ? props.data.username : '',
  set: (v: string) => { props.data.username = v.length > 0 ? v : undefined },
})
const password = computed({
  get: (): string => props.data.password?.length > 0 ? props.data.password : '',
  set: (v: string) => { props.data.password = v.length > 0 ? v : undefined },
})
</script>
