<template>
  <div class="grid2">
    <Field :label="$t('transport.hosts')">
      <input class="input mono" v-model="transport.host" />
    </Field>
    <Field :label="$t('transport.path')">
      <input class="input mono" v-model="transport.path" />
    </Field>
  </div>
  <Headers :data="transport" />
</template>

<script lang="ts" setup>
import Field from '@/components/ui/Field.vue'
import Headers from '../Headers.vue'

defineProps<{ transport: any }>()
</script>
