<template>
  <div class="grid2" style="margin-bottom: 15px;">
    <Field :label="$t('in.ssMethod')" :mb="0">
      <select class="input" v-model="data.method">
        <option v-for="m in ssMethods" :key="m" :value="m">{{ m }}</option>
      </select>
    </Field>
    <Network :data="data" />
    <UoT :data="data" />
  </div>
  <Field :label="$t('types.pw')">
    <input class="input mono" v-model="data.password" />
  </Field>
</template>

<script lang="ts" setup>
import Field from '@/components/ui/Field.vue'
import Network from '../Network.vue'
import UoT from '../UoT.vue'

defineProps<{ data: any }>()

const ssMethods = [
  'none',
  'aes-128-gcm',
  'aes-192-gcm',
  'aes-256-gcm',
  'chacha20-ietf-poly1305',
  'xchacha20-ietf-poly1305',
  '2022-blake3-aes-128-gcm',
  '2022-blake3-aes-256-gcm',
  '2022-blake3-chacha20-poly1305',
]
</script>
