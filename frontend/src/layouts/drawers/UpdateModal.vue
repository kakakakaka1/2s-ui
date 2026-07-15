<template>
  <Modal :open="visible" :title="$t('ui.selfUpdate.title')" :width="460" @close="onClose">
    <div style="padding: 20px;">
      <!-- confirm -->
      <template v-if="state === 'confirm'">
        <div style="display: flex; align-items: center; gap: 10px; margin-bottom: 14px;">
          <span class="mono" style="font-size: 13px; color: var(--text-2);">v{{ current }}</span>
          <Ico name="chevron" :size="14" />
          <span class="mono" style="font-size: 13px; font-weight: 700; color: var(--brand);">v{{ target }}</span>
        </div>
        <p style="font-size: 13px; color: var(--text-2); line-height: 1.6; margin-bottom: 20px;">
          {{ $t('ui.selfUpdate.confirmDesc') }}
        </p>
        <p v-if="ver.isDocker" class="docker-note">
          <Ico name="bolt" :size="14" style="flex: none; margin-top: 1px;" />
          {{ $t('ui.selfUpdate.dockerNote') }}
        </p>
        <div style="display: flex; gap: 10px;">
          <Btn style="flex: 1;" @click="onClose">{{ $t('ui.selfUpdate.cancel') }}</Btn>
          <Btn variant="primary" style="flex: 1;" @click="start"><Ico name="download" :size="15" /> {{ $t('ui.selfUpdate.now') }}</Btn>
        </div>
      </template>

      <!-- running / restarting -->
      <template v-else-if="state === 'running'">
        <div style="display: flex; flex-direction: column; align-items: center; gap: 14px; padding: 10px 0;">
          <span class="upd-spinner" />
          <div style="font-size: 13.5px; font-weight: 600;">{{ phaseLabel }}</div>
          <div v-if="reconnecting" style="font-size: 12px; color: var(--text-3); text-align: center;">
            {{ $t('ui.selfUpdate.reconnecting') }}
          </div>
          <div v-else style="font-size: 12px; color: var(--text-3);">v{{ current }} → v{{ target }}</div>
        </div>
      </template>

      <!-- done (already up to date) -->
      <template v-else-if="state === 'done'">
        <div style="display: flex; flex-direction: column; align-items: center; gap: 12px; padding: 8px 0;">
          <Ico name="check" :size="30" style="color: var(--emerald);" />
          <div style="font-size: 13.5px;">{{ doneMessage }}</div>
          <Btn variant="primary" style="width: 100%; margin-top: 8px;" @click="onClose">{{ $t('close') }}</Btn>
        </div>
      </template>

      <!-- failed -->
      <template v-else-if="state === 'failed'">
        <div style="display: flex; flex-direction: column; gap: 12px;">
          <div style="display: flex; align-items: center; gap: 8px; color: var(--rose); font-weight: 600; font-size: 13.5px;">
            <Ico name="close" :size="18" /> {{ $t('ui.selfUpdate.failed') }}
          </div>
          <div class="mono" style="font-size: 11.5px; color: var(--text-3); background: var(--bg-2); border-radius: 8px; padding: 10px; word-break: break-all; max-height: 140px; overflow-y: auto;">
            {{ errorMessage }}
          </div>
          <Btn style="width: 100%;" @click="onClose">{{ $t('close') }}</Btn>
        </div>
      </template>
    </div>
  </Modal>
</template>

<script lang="ts" setup>
import { ref, computed, watch, onBeforeUnmount } from 'vue'
import { useI18n } from 'vue-i18n'
import HttpUtils from '@/plugins/httputil'
import Modal from '@/components/ui/Modal.vue'
import Btn from '@/components/ui/Btn.vue'
import Ico from '@/components/ui/Ico.vue'
import VersionStore from '@/store/modules/version'

const ver = VersionStore()
const props = defineProps<{ visible: boolean; current: string; target: string }>()
const emit = defineEmits<{ close: [] }>()
const { t } = useI18n({ useScope: 'global' })

type State = 'confirm' | 'running' | 'done' | 'failed'
const state = ref<State>('confirm')
const phase = ref('')
const reconnecting = ref(false)
const doneMessage = ref('')
const errorMessage = ref('')

let timer: ReturnType<typeof setTimeout> | null = null
let sawDown = false

const stopPoll = () => { if (timer) { clearTimeout(timer); timer = null } }

watch(() => props.visible, (v) => {
  if (v) {
    // reset to a clean confirm screen every time it opens
    state.value = 'confirm'
    phase.value = ''
    reconnecting.value = false
    doneMessage.value = ''
    errorMessage.value = ''
    sawDown = false
  } else {
    stopPoll()
  }
})
onBeforeUnmount(stopPoll)

const phaseLabel = computed(() => {
  const key = 'ui.selfUpdate.phase.' + (phase.value || 'checking')
  const label = t(key)
  return label === key ? t('ui.selfUpdate.phase.checking') : label
})

const start = async () => {
  state.value = 'running'
  phase.value = 'checking'
  const msg = await HttpUtils.post('api/updatePanel', null)
  if (!msg.success) {
    errorMessage.value = msg.msg || t('ui.selfUpdate.failed')
    state.value = 'failed'
    return
  }
  poll()
}

// Poll updateStatus with raw fetch (not HttpUtils) so the expected failures
// while the panel restarts don't spam error toasts. A successful request after
// the panel has been down means the new version is up — reload into it.
const poll = () => {
  timer = setTimeout(async () => {
    try {
      const resp = await fetch('api/updateStatus', {
        headers: { 'X-Requested-With': 'XMLHttpRequest' },
        cache: 'no-store',
      })
      if (!resp.ok) throw new Error('bad status')
      const body = await resp.json()

      if (sawDown) {
        // panel answered again after going down => restarted on the new build
        location.reload()
        return
      }

      const obj = body?.obj ?? {}
      phase.value = obj.phase || phase.value
      if (obj.phase === 'restarting') reconnecting.value = true

      if (obj.phase === 'done') {
        doneMessage.value = obj.message === 'already up to date'
          ? t('ui.selfUpdate.upToDate')
          : (obj.message || t('ui.selfUpdate.upToDate'))
        state.value = 'done'
        stopPoll()
        return
      }
      if (obj.phase === 'failed') {
        errorMessage.value = obj.error || obj.message || t('ui.selfUpdate.failed')
        state.value = 'failed'
        stopPoll()
        return
      }
      poll()
    } catch {
      // request failed: the panel is restarting. Mark it down and keep pinging;
      // the next success triggers the reload above.
      sawDown = true
      reconnecting.value = true
      poll()
    }
  }, 1500)
}

const onClose = () => {
  stopPoll()
  emit('close')
}
</script>

<style scoped>
.docker-note {
  display: flex; gap: 8px;
  margin: -8px 0 20px; padding: 9px 11px;
  font-size: 12px; line-height: 1.55; color: var(--amber);
  background: color-mix(in srgb, var(--amber) 12%, transparent);
  border: 1px solid color-mix(in srgb, var(--amber) 26%, transparent);
  border-radius: var(--radius-sm);
}
.upd-spinner {
  width: 28px; height: 28px;
  border: 3px solid var(--line-2);
  border-top-color: var(--brand);
  border-radius: 50%;
  animation: spin .7s linear infinite;
}
</style>
