<template>
  <div style="display: flex; align-items: center; gap: 10px; margin-bottom: 12px;">
    <SectionLabel>SSH</SectionLabel>
    <div style="flex: 1;" />
    <Pop :min-width="210" direction="down">
      <template #trigger="{ toggle }">
        <Btn variant="subtle" sm @click="toggle">{{ $t('types.ssh.options') }}</Btn>
      </template>
      <div style="display: flex; flex-direction: column; gap: 2px; padding: 4px;">
        <div class="pop-item"><SwitchLabel v-model="optionKey" label="SSH Key" /></div>
        <div class="pop-item"><SwitchLabel v-model="optionHostKey" :label="$t('types.ssh.hostKey')" /></div>
        <div class="pop-item"><SwitchLabel v-model="optionAlgorithms" :label="$t('types.ssh.algorithm')" /></div>
        <div class="pop-item"><SwitchLabel v-model="optionVer" :label="$t('types.ssh.clientVer')" /></div>
      </div>
    </Pop>
  </div>

  <template v-if="optionKey">
    <div style="margin-bottom: 12px;">
      <Segmented v-model="usePath" :options="[[0, $t('tls.usePath')], [1, $t('tls.useText')]]" block
        @update:model-value="switchKeyMode" />
    </div>
    <Field v-if="usePath == 0" :label="$t('tls.keyPath')">
      <input class="input mono" v-model="data.private_key_path" />
    </Field>
    <Field v-else :label="$t('tls.key')">
      <textarea class="input mono" rows="4" spellcheck="false" v-model="data.private_key"></textarea>
    </Field>
    <Field :label="$t('types.ssh.passphrase')">
      <input class="input mono" type="password" v-model="data.private_key_passphrase" />
    </Field>
  </template>
  <template v-else>
    <div class="grid2" style="margin-bottom: 15px;">
      <Field :label="$t('types.un')" :mb="0">
        <input class="input mono" v-model="data.user" />
      </Field>
      <Field :label="$t('types.pw')" :mb="0">
        <input class="input mono" type="password" v-model="data.password" />
      </Field>
    </div>
  </template>

  <Field v-if="optionHostKey" :label="$t('types.ssh.hostKey')">
    <textarea class="input mono" rows="4" spellcheck="false" v-model="hostKey"></textarea>
  </Field>
  <div class="grid2" style="margin-bottom: 15px;" v-if="data.host_key_algorithms != undefined || data.client_version != undefined">
    <Field v-if="data.host_key_algorithms != undefined" :label="$t('types.ssh.algorithm') + ' ' + $t('commaSeparated')" :mb="0">
      <input class="input mono" v-model="algorithms" />
    </Field>
    <Field v-if="data.client_version != undefined" :label="$t('types.ssh.clientVer')" :mb="0">
      <input class="input mono" v-model="data.client_version" />
    </Field>
  </div>
</template>

<script lang="ts" setup>
import { computed, ref } from 'vue'
import Field from '@/components/ui/Field.vue'
import Btn from '@/components/ui/Btn.vue'
import Pop from '@/components/ui/Pop.vue'
import SwitchLabel from '@/components/ui/SwitchLabel.vue'
import Segmented from '@/components/ui/Segmented.vue'
import SectionLabel from '@/components/ui/SectionLabel.vue'

const props = defineProps<{ data: any }>()

const usePath = ref<string | number>(props.data.private_key != undefined ? 1 : 0)

const switchKeyMode = (v: string | number) => {
  if (v == 0) {
    props.data.private_key = undefined
    props.data.private_key_path = ''
  } else {
    props.data.private_key_path = undefined
    props.data.private_key = ''
  }
}

const optionKey = computed({
  get: (): boolean => props.data.private_key != undefined || props.data.private_key_path != undefined,
  set: (v: boolean) => {
    usePath.value = 0
    if (v) {
      props.data.private_key_path = ''
      delete props.data.user
      delete props.data.password
    } else {
      delete props.data.private_key_path
      delete props.data.private_key
      delete props.data.private_key_passphrase
    }
  },
})
const optionHostKey = computed({
  get: (): boolean => props.data.host_key != undefined,
  set: (v: boolean) => { props.data.host_key = v ? '' : undefined },
})
const optionAlgorithms = computed({
  get: (): boolean => props.data.host_key_algorithms != undefined,
  set: (v: boolean) => { props.data.host_key_algorithms = v ? [] : undefined },
})
const optionVer = computed({
  get: (): boolean => props.data.client_version != undefined,
  set: (v: boolean) => { props.data.client_version = v ? 'SSH-2.0-OpenSSH_7.4p1' : undefined },
})
const hostKey = computed({
  get: (): string => props.data.host_key ? (Array.isArray(props.data.host_key) ? props.data.host_key.join('\n') : props.data.host_key) : '',
  set: (v: string) => { props.data.host_key = v.split('\n') },
})
const algorithms = computed({
  get: () => props.data.host_key_algorithms ? props.data.host_key_algorithms.join(',') : '',
  set: (v: string) => { props.data.host_key_algorithms = v.length > 0 ? v.split(',') : undefined },
})
</script>
