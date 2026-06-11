<template>
  <div>
    <SectionLabel style="margin-bottom: 12px;">ShadowTls</SectionLabel>
    <div class="grid2">
      <Field :label="$t('version')">
        <select class="input" :disabled="data.id > 0" v-model="version">
          <option v-for="v in [1, 2, 3]" :key="v" :value="v">{{ v }}</option>
        </select>
      </Field>
      <Field v-if="data.password != undefined" :label="$t('types.pw')">
        <input class="input mono" v-model="data.password" />
      </Field>
      <Field v-if="data.wildcard_sni != undefined" label="Wildcard SNI">
        <select class="input" v-model="data.wildcard_sni">
          <option value="">{{ $t('none') }}</option>
          <option v-for="w in ['off', 'authed', 'all']" :key="w" :value="w">{{ w }}</option>
        </select>
      </Field>
    </div>
    <div class="grid2">
      <Field :label="$t('types.shdwTls.hs')">
        <input class="input mono" v-model="data.handshake.server" />
      </Field>
      <Field :label="$t('out.port')">
        <input class="input mono" type="number" min="0" v-model.number="server_port" />
      </Field>
    </div>
    <Dial :dial="data.handshake" />
    <template v-if="data.handshake_for_server_name != undefined">
      <Field :label="$t('types.shdwTls.addHS')">
        <div style="display: flex; gap: 8px;">
          <input class="input mono" style="flex: 1; min-width: 0;" v-model="handshake_server" />
          <Btn sm icon :disabled="handshake_server == ''" :title="$t('actions.add')" @click="addHandshakeServer">
            <Ico name="plus" :size="14" />
          </Btn>
        </div>
      </Field>
      <div
        v-for="(value, key) in data.handshake_for_server_name"
        :key="key"
        class="card"
        style="padding: 14px; margin-bottom: 12px; background: var(--surface-2);"
      >
        <div style="display: flex; align-items: center; margin-bottom: 10px;">
          <Chip><span class="mono">{{ key }}</span></Chip>
          <IconBtn
            name="trash"
            danger
            style="margin-inline-start: auto;"
            :title="$t('actions.del')"
            @click="data.handshake_for_server_name ? delete data.handshake_for_server_name[key] : null"
          />
        </div>
        <div class="grid2">
          <Field :label="$t('types.shdwTls.hs')">
            <input class="input mono" v-model="value.server" />
          </Field>
          <Field :label="$t('out.port')">
            <input class="input mono" type="number" min="0" v-model.number="value.server_port" />
          </Field>
        </div>
        <Dial :dial="value" />
      </div>
    </template>
  </div>
</template>

<script lang="ts" setup>
import { computed, ref } from 'vue'
import Field from '@/components/ui/Field.vue'
import SectionLabel from '@/components/ui/SectionLabel.vue'
import Chip from '@/components/ui/Chip.vue'
import Btn from '@/components/ui/Btn.vue'
import Ico from '@/components/ui/Ico.vue'
import IconBtn from '@/components/ui/IconBtn.vue'
import Dial from './Dial.vue'

const props = defineProps<{ data: any }>()

const handshake_server = ref('')

const addHandshakeServer = () => {
  props.data.handshake_for_server_name[handshake_server.value] = {}
  // Clear the input field after adding the server
  handshake_server.value = ''
}

const version = computed({
  get: (): number => props.data.version,
  set: (v: number) => {
    switch (v) {
      case 1:
        delete props.data.password
        delete props.data.handshake_for_server_name
        delete props.data.wildcard_sni
        break
      case 2:
        if (!props.data.password) {
          props.data.password = ''
        }
        if (!props.data.handshake_for_server_name) {
          props.data.handshake_for_server_name = {}
        }
        delete props.data.wildcard_sni
        break
      case 3:
        delete props.data.password
        if (!props.data.handshake_for_server_name) {
          props.data.handshake_for_server_name = {}
        }
        if (!props.data.wildcard_sni) {
          props.data.wildcard_sni = ''
        }
        break
    }
    props.data.version = v
  },
})

const server_port = computed({
  get: (): any => (props.data.handshake.server_port ? props.data.handshake.server_port : 443),
  set: (v: any) => {
    props.data.handshake.server_port = v?.length == 0 || v == 0 ? 443 : parseInt(v)
  },
})
</script>
