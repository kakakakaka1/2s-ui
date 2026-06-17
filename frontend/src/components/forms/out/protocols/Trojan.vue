<template>
  <div class="grid2" style="margin-bottom: 15px;">
    <Field :label="$t('types.pw')" :mb="0">
      <input class="input mono" v-model="data.password" />
    </Field>
    <Network :data="data" />
  </div>
</template>

<script lang="ts" setup>
import Field from '@/components/ui/Field.vue'
import Network from '../Network.vue'

defineProps<{ data: any }>()
</script>
