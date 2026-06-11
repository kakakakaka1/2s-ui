<template>
  <Modal :open="visible" :title="$t('basic.log.title')" :width="1100" @close="$emit('close')">
    <div style="padding: 18px 20px;">
      <div style="display: flex; align-items: flex-end; gap: 12px; flex-wrap: wrap; margin-bottom: 14px;">
        <Field :label="$t('basic.log.level')" :mb="0" style="width: 160px;">
          <Select v-model="logLevel" @change="loadData">
            <option value="debug">DEBUG</option>
            <option value="info">INFO</option>
            <option value="warning">WARNING</option>
            <option value="err">ERROR</option>
          </Select>
        </Field>
        <Field :label="$t('count')" :mb="0" style="width: 120px;">
          <Select v-model.number="logCount" @change="loadData">
            <option v-for="c in [10, 20, 30, 50, 100]" :key="c" :value="c">{{ c }}</option>
          </Select>
        </Field>
        <div style="flex: 1;" />
        <Btn :loading="loading" @click="loadData"><Ico name="refresh" :size="15" /> {{ $t('actions.update') }}</Btn>
      </div>
      <div class="log-box mono" dir="ltr">
        <div v-if="lines.length === 0" style="color: var(--text-3);">—</div>
        <div v-for="(line, i) in lines" :key="i" class="log-line">{{ line }}</div>
      </div>
    </div>
  </Modal>
</template>

<script lang="ts" setup>
import Select from '@/components/ui/Select.vue'
import { ref, watch } from 'vue'
import HttpUtils from '@/plugins/httputil'
import Modal from '@/components/ui/Modal.vue'
import Field from '@/components/ui/Field.vue'
import Btn from '@/components/ui/Btn.vue'
import Ico from '@/components/ui/Ico.vue'

const props = defineProps<{ visible: boolean }>()
defineEmits<{ close: [] }>()

const loading = ref(false)
const lines = ref<string[]>([])
const logLevel = ref('info')
const logCount = ref(10)

const loadData = async () => {
  loading.value = true
  const data = await HttpUtils.get('api/logs', { c: logCount.value, l: logLevel.value })
  if (data.success) lines.value = data.obj ?? []
  loading.value = false
}

watch(() => props.visible, (v) => {
  if (v) {
    lines.value = []
    logLevel.value = 'info'
    logCount.value = 10
    loadData()
  }
})
</script>

<style scoped>
.log-box {
  background: var(--surface-3);
  border: 1px solid var(--line);
  border-radius: var(--radius-md);
  padding: 14px 16px;
  font-size: 12px;
  line-height: 1.7;
  max-height: 50vh;
  overflow: auto;
  white-space: pre-wrap;
  word-break: break-all;
}
.log-line + .log-line { border-top: 1px dashed var(--line); }
</style>
