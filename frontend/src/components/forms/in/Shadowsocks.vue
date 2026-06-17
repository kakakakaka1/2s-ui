<template>
  <div>
    <SectionLabel style="margin-bottom: 12px;">Shadowsocks</SectionLabel>
    <div class="grid2">
      <Field :label="$t('in.ssMethod')">
        <Select v-model="data.method" @change="onMethodChange">
          <option v-for="m in ssMethods" :key="m" :value="m">{{ m }}</option>
        </Select>
      </Field>
      <Network :data="data" />
    </div>
    <div style="display: flex; gap: 24px; flex-wrap: wrap; margin-bottom: 15px;">
      <SwitchLabel v-model="data.managed" :label="$t('in.ssManageable')" />
    </div>
    <Field v-if="data.method != 'none'" :label="$t('types.pw')">
      <KeyInput v-model="data.password" :title="$t('actions.generate')" @regenerate="changeMethod(data.method)" />
    </Field>
  </div>
</template>

<script lang="ts" setup>
import Select from '@/components/ui/Select.vue'
import RandomUtil from '@/plugins/randomUtil'
import Field from '@/components/ui/Field.vue'
import SectionLabel from '@/components/ui/SectionLabel.vue'
import SwitchLabel from '@/components/ui/SwitchLabel.vue'
import KeyInput from '@/components/ui/KeyInput.vue'
import Network from './Network.vue'

const props = defineProps<{ data: any; direction?: string }>()

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

const changeMethod = (ssMethod: string) => {
  if (ssMethod.startsWith('2022')) {
    props.data.password = ssMethod == '2022-blake3-aes-128-gcm'
      ? RandomUtil.randomShadowsocksPassword(16)
      : RandomUtil.randomShadowsocksPassword(32)
  } else if (ssMethod == 'none') {
    delete props.data.password
  } else {
    props.data.password = RandomUtil.randomSeq(10)
  }
}

const onMethodChange = () => {
  changeMethod(props.data.method)
}
</script>
