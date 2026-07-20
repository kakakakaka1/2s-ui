<template>
  <MDrawer
    :open="visible"
    icon="server"
    color="var(--brand)"
    :title="isNew ? $t('ui.nodeNew') : nodeData.name"
    :sub="$t('ui.nodeSub')"
    :save-label="isNew ? $t('ui.create') : $t('actions.save')"
    :width="500"
    :loading="loading"
    @close="$emit('close')"
    @save="saveChanges"
  >
    <Field :label="$t('node.name')">
      <input class="input" v-model="nodeData.name" placeholder="jp-tokyo-1" />
    </Field>

    <Field :label="$t('node.baseUrl')" :hint="$t('node.baseUrlHint')">
      <input class="input mono" v-model="nodeData.baseUrl" placeholder="https://1.2.3.4:2095" />
    </Field>

    <Field :label="$t('node.webPath')">
      <input class="input mono" v-model="nodeData.webPath" placeholder="/app/" />
    </Field>

    <Field :label="$t('node.token')" :hint="$t('node.tokenHint')">
      <input
        class="input mono"
        type="password"
        v-model="nodeData.token"
        :placeholder="!isNew && nodeData.tokenSet ? $t('node.tokenKeep') : ''"
        autocomplete="new-password"
      />
    </Field>

    <template v-if="isHttps">
      <SwitchLabel v-model="nodeData.insecure" :label="$t('node.insecure')" />
      <Field :label="$t('node.certPin')" :hint="$t('node.certPinHint')">
        <input class="input mono" v-model="nodeData.certPin" placeholder="sha256 hex" />
      </Field>
    </template>
    <MHint v-else-if="nodeData.baseUrl.startsWith('http://')">{{ $t('node.httpWarn') }}</MHint>

    <Field :label="$t('node.desc')">
      <input class="input" v-model="nodeData.desc" />
    </Field>

    <SwitchLabel v-model="nodeData.enable" :label="$t('enable')" />

    <hr class="form-divider" />

    <div style="display: flex; align-items: center; gap: 10px;">
      <Btn :loading="testing" @click="testConnection">
        <Ico name="bolt" :size="15" /> {{ $t('node.test') }}
      </Btn>
      <div v-if="testResult" style="font-size: 12.5px; min-width: 0; overflow: hidden; text-overflow: ellipsis;">
        <span v-if="testResult.state === 'online'" style="color: var(--emerald); font-weight: 600;">
          {{ $t('node.status.online') }} · {{ testResult.latency }} {{ $t('date.ms') }} ·
          <span class="mono">{{ testResult.appVersion }} / {{ testResult.coreVersion }}</span>
        </span>
        <span v-else-if="testResult.state === 'core-stopped'" style="color: var(--amber); font-weight: 600;">
          {{ $t('node.testOk') }} · {{ $t('node.status.coreStopped') }}
        </span>
        <span v-else style="color: var(--rose); font-weight: 600;">{{ testResult.error }}</span>
      </div>
    </div>
  </MDrawer>
</template>

<script lang="ts" setup>
import { computed, ref, watch } from 'vue'
import Data from '@/store/modules/data'
import HttpUtils from '@/plugins/httputil'
import { Node, NodeStatus, defaultNode } from '@/types/node'
import MDrawer from '@/components/ui/MDrawer.vue'
import Field from '@/components/ui/Field.vue'
import Btn from '@/components/ui/Btn.vue'
import Ico from '@/components/ui/Ico.vue'
import MHint from '@/components/ui/MHint.vue'
import SwitchLabel from '@/components/ui/SwitchLabel.vue'

const props = defineProps<{
  visible: boolean
  id: number
  data: string
}>()

const emit = defineEmits<{ close: [] }>()

const nodeData = ref<Node>({ ...defaultNode })
const loading = ref(false)

const isNew = computed(() => props.id === 0)
const isHttps = computed(() => nodeData.value.baseUrl.toLowerCase().startsWith('https://'))

const updateData = (id: number) => {
  testResult.value = null
  if (id > 0) {
    // The server never sends the token back; an empty field on save means
    // "keep the stored one" (tokenSet drives the placeholder).
    nodeData.value = { ...defaultNode, ...(<Node>JSON.parse(props.data)), token: '' }
  } else {
    nodeData.value = { ...defaultNode }
  }
}

watch(() => props.visible, (v) => { if (v) updateData(props.id) })

const saveChanges = async () => {
  if (!props.visible) return
  loading.value = true
  const success = await Data().save('nodes', nodeData.value.id > 0 ? 'edit' : 'new', nodeData.value)
  if (success) emit('close')
  loading.value = false
}

// ---------------- test connection ----------------
const testing = ref(false)
const testResult = ref<(NodeStatus & { error?: string }) | null>(null)

const testConnection = async () => {
  testing.value = true
  testResult.value = null
  const msg = await HttpUtils.post('api/testNode', { data: JSON.stringify(nodeData.value) })
  testing.value = false
  if (msg.success && msg.obj) {
    testResult.value = msg.obj
    // probe errors surface inside the status object
    if (msg.obj.state !== 'online' && msg.obj.state !== 'core-stopped' && !msg.obj.error) {
      testResult.value = { ...msg.obj, error: msg.msg }
    }
  } else {
    testResult.value = <any>{ state: 'offline', error: msg.msg }
  }
}
</script>
