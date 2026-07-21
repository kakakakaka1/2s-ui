<template>
  <div class="grid2" style="margin-bottom: 15px;">
    <Field :label="$t('in.ssMethod')" :mb="0">
      <Select v-model="data.method">
        <option v-for="m in ssMethods" :key="m" :value="m">{{ m }}</option>
      </Select>
    </Field>
    <Network :data="data" />
    <UoT :data="data" />
  </div>
  <Field :label="$t('types.pw')">
    <input class="input mono" v-model="data.password" />
  </Field>
  <div class="grid2">
    <Field label="Plugin" :mb="0">
      <input class="input mono" v-model="plugin" />
    </Field>
    <Field v-if="data.plugin" label="Plugin Options" :mb="0">
      <input class="input mono" v-model="pluginOpts" />
    </Field>
  </div>
</template>

<script lang="ts" setup>
import { computed } from 'vue'
import Select from '@/components/ui/Select.vue'
import Field from '@/components/ui/Field.vue'
import Network from '../Network.vue'
import UoT from '../UoT.vue'

const props = defineProps<{ data: any }>()

// Empty strings must collapse to undefined: the value is spread straight into
// the sing-box outbound JSON, and an empty `plugin` key is not the same as none.
const plugin = computed({
  get: (): string => props.data.plugin ?? '',
  set: (v: string) => {
    props.data.plugin = v.length > 0 ? v : undefined
    if (!props.data.plugin) props.data.plugin_opts = undefined
  },
})
const pluginOpts = computed({
  get: (): string => props.data.plugin_opts ?? '',
  set: (v: string) => { props.data.plugin_opts = v.length > 0 ? v : undefined },
})

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
