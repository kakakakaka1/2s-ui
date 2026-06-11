<template>
  <Modal :open="open" :title="$t('admin.changes')" :width="860" @close="$emit('close')">
    <div style="padding: 18px 20px;">
      <!-- filters -->
      <div style="display: flex; align-items: flex-end; gap: 12px; flex-wrap: wrap; margin-bottom: 14px;">
        <Field :label="$t('admin.actor')" :mb="0" style="width: 160px;">
          <select class="input" v-model="user" @change="loadData">
            <option value="">{{ $t('ui.none') }}</option>
            <option value="DepleteJob">DepleteJob</option>
            <option v-for="a in admins" :key="a" :value="a">{{ a }}</option>
          </select>
        </Field>
        <Field :label="$t('admin.key')" :mb="0" style="width: 160px;">
          <select class="input" v-model="key" @change="loadData">
            <option value="">{{ $t('ui.none') }}</option>
            <option v-for="k in keys" :key="k" :value="k">{{ k }}</option>
          </select>
        </Field>
        <Field :label="$t('count')" :mb="0" style="width: 100px;">
          <select class="input" v-model.number="chngCount" @change="loadData">
            <option v-for="c in [10, 20, 30, 50, 100]" :key="c" :value="c">{{ c }}</option>
          </select>
        </Field>
        <div style="flex: 1;" />
        <Btn :loading="loading" @click="loadData"><Ico name="refresh" :size="15" /> {{ $t('actions.update') }}</Btn>
      </div>

      <!-- table -->
      <div style="border: 1px solid var(--line); border-radius: var(--radius-md); overflow: auto; max-height: 56vh;">
        <table class="dtable">
          <thead>
            <tr>
              <th>ID</th>
              <th>{{ $t('admin.date') }}-{{ $t('admin.time') }}</th>
              <th>{{ $t('admin.actor') }}</th>
              <th>{{ $t('admin.key') }}</th>
              <th>{{ $t('admin.action') }}</th>
              <th style="width: 40px;"></th>
            </tr>
          </thead>
          <tbody>
            <template v-for="item in changes" :key="item.id">
              <tr class="clickable" @click="toggleExpand(item.id)">
                <td class="mono" style="padding: 9px 16px;">{{ item.id }}</td>
                <td class="mono" dir="ltr" style="padding: 9px 16px; white-space: nowrap;">{{ dateFormatted(item.dateTime) }}</td>
                <td style="padding: 9px 16px;">{{ item.actor }}</td>
                <td style="padding: 9px 16px;">{{ item.key }}</td>
                <td style="padding: 9px 16px;"><Chip>{{ $t('actions.' + item.action) }}</Chip></td>
                <td style="padding: 9px 10px; color: var(--text-3);">
                  <Ico
                    name="chevronDown"
                    :size="15"
                    :style="{ transform: expanded.includes(item.id) ? 'rotate(180deg)' : 'none', transition: 'transform .15s var(--ease)' }"
                  />
                </td>
              </tr>
              <tr v-if="expanded.includes(item.id)">
                <td :colspan="6" style="padding: 0;">
                  <div class="obj-box mono" dir="ltr">
                    <div v-if="item.index > 0" style="margin-bottom: 6px; color: var(--text-3);">Index: {{ item.index }}</div>
                    <pre style="margin: 0; white-space: pre-wrap; word-break: break-all;">{{ item.obj }}</pre>
                  </div>
                </td>
              </tr>
            </template>
            <tr v-if="changes.length === 0">
              <td :colspan="6" style="padding: 14px 16px; color: var(--text-3); text-align: center;">{{ $t('noData') }}</td>
            </tr>
          </tbody>
        </table>
      </div>
    </div>
  </Modal>
</template>

<script lang="ts" setup>
import { ref, watch } from 'vue'
import HttpUtils from '@/plugins/httputil'
import { locale } from '@/locales'
import Modal from '@/components/ui/Modal.vue'
import Field from '@/components/ui/Field.vue'
import Btn from '@/components/ui/Btn.vue'
import Ico from '@/components/ui/Ico.vue'
import Chip from '@/components/ui/Chip.vue'

const props = defineProps<{
  open: boolean
  admins: string[]
  actor: string
}>()
defineEmits<{ close: [] }>()

const keys = ['inbounds', 'outbounds', 'clients', 'route', 'tls', 'experimental']

const loading = ref(false)
const changes = ref<any[]>([])
const user = ref('')
const key = ref('')
const chngCount = ref(10)
const expanded = ref<number[]>([])

const loadData = async () => {
  loading.value = true
  const data = await HttpUtils.get('api/changes', { a: user.value, k: key.value, c: chngCount.value })
  if (data.success) {
    changes.value = data.obj ?? []
  }
  loading.value = false
}

const toggleExpand = (id: number) => {
  expanded.value = expanded.value.includes(id)
    ? expanded.value.filter((i) => i !== id)
    : [...expanded.value, id]
}

const dateFormatted = (dt: number): string => {
  const date = new Date(dt * 1000)
  return date.toLocaleString(locale)
}

watch(() => props.open, (v) => {
  changes.value = []
  expanded.value = []
  user.value = props.actor
  key.value = ''
  chngCount.value = 10
  if (v) loadData()
})
</script>

<style scoped>
.obj-box {
  background: var(--surface-3);
  border-bottom: 1px solid var(--line);
  padding: 12px 16px;
  font-size: 12px;
  line-height: 1.6;
}
</style>
