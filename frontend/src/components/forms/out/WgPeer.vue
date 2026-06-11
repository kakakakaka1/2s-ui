<template>
  <Field :label="$t('types.wg.privKey')">
    <KeyInput v-model="privateKey" :title="$t('actions.generate')" @regenerate="$emit('refreshPeerKey')" />
  </Field>
  <Field :label="$t('types.wg.pubKey')">
    <input class="input mono" style="font-size: 12px;" v-model="publicKey" />
  </Field>
  <div class="grid2">
    <Field :label="$t('out.addr')">
      <input class="input mono" v-model="address" placeholder="peer.example.net" />
    </Field>
    <Field :label="$t('out.port')">
      <input class="input mono" type="number" min="0" v-model.number="port" />
    </Field>
  </div>
  <div class="grid2">
    <Field :label="$t('types.wg.allowedIp') + ' ' + $t('commaSeparated')" :mb="0">
      <input class="input mono" v-model="allowedIps" />
    </Field>
    <Field :label="$t('types.wg.psk')" :mb="0">
      <input class="input mono" v-model="data.pre_shared_key" />
    </Field>
    <Field :label="'KeepAlive (' + $t('date.s') + ')'" :mb="0">
      <input class="input mono" type="number" min="0" v-model.number="keepAlive" />
    </Field>
    <Field :label="'Reserved ' + $t('commaSeparated')" :mb="0">
      <input class="input mono" v-model="reserved" />
    </Field>
  </div>
</template>

<script lang="ts" setup>
import { computed } from 'vue'
import Field from '@/components/ui/Field.vue'
import KeyInput from '@/components/ui/KeyInput.vue'

const props = defineProps<{ data: any; ext: any }>()
defineEmits<{ refreshPeerKey: [] }>()

const allowedIps = computed({
  get: () => props.data.allowed_ips?.join(','),
  set: (v: string) => { props.data.allowed_ips = v.length > 0 ? v.split(',') : undefined },
})
const reserved = computed({
  get: () => props.data.reserved?.join(','),
  set: (v: string) => {
    if (!v.endsWith(',')) {
      props.data.reserved = v.length > 0 ? v.split(',').map((str) => parseInt(str, 10)) : undefined
    }
  },
})
const address = computed({
  get: () => props.data.address,
  set: (v: string) => { props.data.address = v.length > 0 ? v : undefined },
})
const port = computed({
  get: () => props.data.port,
  set: (v: number) => { props.data.port = v > 0 ? v : undefined },
})
const keepAlive = computed({
  get: () => props.data.persistent_keepalive_interval ?? 0,
  set: (v: number) => { props.data.persistent_keepalive_interval = v > 0 ? v : undefined },
})
const privateKey = computed({
  get: () => {
    const indexKeys = props.ext?.keys.findIndex((key: any) => key.public_key == props.data.public_key) ?? -1
    return indexKeys > -1 ? props.ext.keys[indexKeys].private_key : ''
  },
  set: (v: string) => {
    const indexKeys = props.ext?.keys.findIndex((key: any) => key.public_key == props.data.public_key) ?? -1
    props.ext.keys[indexKeys].private_key = v
  },
})
const publicKey = computed({
  get: () => props.data.public_key,
  set: (v: string) => {
    const indexKeys = props.ext?.keys.findIndex((key: any) => key.public_key == props.data.public_key) ?? -1
    props.ext.keys[indexKeys].public_key = v
    props.data.public_key = v
  },
})
</script>
