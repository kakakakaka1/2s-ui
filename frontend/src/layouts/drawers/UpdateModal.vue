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
          <Btn variant="primary" style="width: 100%; margin-top: 8px;" @click="onClose">{{ $t('actions.close') }}</Btn>
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
          <Btn style="width: 100%;" @click="onClose">{{ $t('actions.close') }}</Btn>
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
let downSince = 0

// 面板掉线后等它回来的上限。正常重启只要几秒,超过这个还连不上基本就是起不来了,
// 不能让"会自动重新连接"的承诺无限转圈下去。只从第一次不可达起算,所以下载慢
// 不会被误判成超时。
const RECONNECT_TIMEOUT = 120_000

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
    downSince = 0
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

// 用裸 fetch 而不是 HttpUtils 轮询 updateStatus:面板重启期间请求必然失败,
// 走 HttpUtils 会刷一屏错误 toast。
const poll = () => {
  timer = setTimeout(async () => {
    try {
      const resp = await fetch('api/updateStatus', {
        headers: { 'X-Requested-With': 'XMLHttpRequest' },
        cache: 'no-store',
      })
      if (!resp.ok) throw new Error('bad status')
      const body = await resp.json()
      const obj = body?.obj ?? {}

      // 判定"该载入新版面板了":后端答出的已经不是本次更新的进行态。
      //  - phase === 'idle':更新状态机是后端包级变量,StartUpdate 一进来就推到
      //    checking 且只进不退,同一进程绝不可能再答出 idle —— 能答出 idle 的只会是
      //    重启后的全新进程。这修的正是主 bug:重启不可达窗口常短于 1.5s 轮询间隔,
      //    错过掉线后下一拍就撞上新进程的 idle,原来只当它"既非 done 也非 failed"
      //    继续空转,更新其实早成了界面却一直转圈。
      //  - 没有 phase:正常响应必带 phase(至少是 idle);缺了只可能是 session 中途
      //    失效后 checkLogin 返回的 {success:false}(HTTP 200、无 obj)—— 本会话再也
      //    拿不到状态,reload 让它落到登录页,不至于无限空转。
      // 不再用 sawDown 判定:它只表示"掉线过",下载期一次网络抖动就会误置,恢复后
      // phase 明明还在 downloading 却被当成"已重启"而误 reload。看后端上报的 phase 更准。
      if (obj.phase === 'idle' || !obj.phase) {
        location.reload()
        return
      }

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
      // 请求失败:面板正在重启。记下第一次掉线的时间用于超时计时,继续 ping;
      // 面板起回来后会在上面按 phase 判定 reload。
      if (!sawDown) {
        sawDown = true
        downSince = Date.now()
      }
      reconnecting.value = true
      if (Date.now() - downSince > RECONNECT_TIMEOUT) {
        errorMessage.value = t('ui.selfUpdate.timeout')
        state.value = 'failed'
        stopPoll()
        return
      }
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
