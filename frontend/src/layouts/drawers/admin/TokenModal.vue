<template>
  <Modal :open="open" :title="$t('admin.api.title')" :width="800" @close="$emit('close')">
    <div style="padding: 18px 20px;">
      <!-- freshly created token (shown once) -->
      <div v-if="newToken.token.length > 0" class="token-alert">
        <div style="font-size: 12.5px; font-weight: 600; margin-bottom: 8px;">{{ $t('admin.api.msg') }}</div>
        <div style="display: flex; gap: 8px; align-items: center;">
          <input
            class="input mono"
            style="flex: 1;"
            readonly
            :value="newToken.token"
            @focus="($event.target as HTMLInputElement).select()"
          />
          <Btn variant="subtle" icon :title="$t('copyToClipboard')" @click="copyToClipboard(newToken.token)">
            <Ico name="copy" :size="16" />
          </Btn>
        </div>
      </div>

      <!-- token list -->
      <div style="border: 1px solid var(--line); border-radius: var(--radius-md); overflow: hidden; margin-bottom: 14px;">
        <table class="dtable">
          <thead>
            <tr>
              <th>#</th>
              <th>{{ $t('admin.api.token') }}</th>
              <th>{{ $t('client.desc') }}</th>
              <th>{{ $t('date.expiry') }}</th>
              <th style="width: 84px;">{{ $t('actions.del') }}</th>
            </tr>
          </thead>
          <tbody>
            <tr v-for="token in tokens" :key="token.id">
              <td class="mono" style="padding: 8px 16px;">{{ token.id }}</td>
              <td class="mono" style="padding: 8px 16px; word-break: break-all;">{{ token.token }}</td>
              <td style="padding: 8px 16px;">{{ token.desc }}</td>
              <td class="mono" style="padding: 8px 16px; white-space: nowrap;">{{ dateFormatted(token.expiry) }}</td>
              <td style="padding: 4px 10px;">
                <div v-if="delConfirm === token.id" style="display: flex; align-items: center; gap: 2px;">
                  <IconBtn name="check" danger :title="$t('yes')" @click="deleteToken(token.id)" />
                  <IconBtn name="close" :title="$t('no')" @click="delConfirm = null" />
                </div>
                <IconBtn v-else name="trash" danger :title="$t('actions.del')" @click="delConfirm = token.id" />
              </td>
            </tr>
            <tr v-if="tokens.length === 0">
              <td :colspan="5" style="padding: 14px 16px; color: var(--text-3); text-align: center;">{{ $t('noData') }}</td>
            </tr>
          </tbody>
        </table>
      </div>

      <!-- add token -->
      <div v-if="showNew" class="add-box">
        <div class="grid2">
          <Field :label="$t('client.desc')" :mb="0">
            <input class="input" v-model="newToken.desc" />
          </Field>
          <Field :label="$t('date.expiry')" :hint="$t('date.d')" :mb="0">
            <input class="input mono" type="number" min="0" v-model.number="newToken.expiry" />
          </Field>
        </div>
        <div style="display: flex; justify-content: flex-end; gap: 8px; margin-top: 12px;">
          <Btn sm @click="showNew = false">{{ $t('actions.close') }}</Btn>
          <Btn variant="primary" sm :loading="loading" @click="addToken">{{ $t('actions.add') }}</Btn>
        </div>
      </div>
      <Btn v-else variant="primary" sm @click="showAddToken"><Ico name="plus" :size="15" /> {{ $t('actions.add') }}</Btn>
    </div>

    <template #footer>
      <Btn sm @click="$emit('close')">{{ $t('actions.close') }}</Btn>
    </template>
  </Modal>
</template>

<script lang="ts" setup>
import { ref, watch } from 'vue'
import HttpUtils from '@/plugins/httputil'
import { i18n, locale } from '@/locales'
import { copyToClipboard } from '@/plugins/clipboard'
import Modal from '@/components/ui/Modal.vue'
import Field from '@/components/ui/Field.vue'
import Btn from '@/components/ui/Btn.vue'
import Ico from '@/components/ui/Ico.vue'
import IconBtn from '@/components/ui/IconBtn.vue'

const props = defineProps<{
  open: boolean
}>()
defineEmits<{ close: [] }>()

const loading = ref(false)
const tokens = ref<any[]>([])
const showNew = ref(false)
const delConfirm = ref<number | null>(null)
const newToken = ref({ desc: '', token: '', expiry: 0 })

const loadData = async () => {
  loading.value = true
  const data = await HttpUtils.get('api/tokens')
  if (data.success) {
    tokens.value = data.obj ?? []
    delConfirm.value = null
  }
  loading.value = false
}

const resetNewToken = () => {
  newToken.value = { desc: '', token: '', expiry: 30 }
}

const showAddToken = () => {
  resetNewToken()
  showNew.value = true
}

const addToken = async () => {
  loading.value = true
  newToken.value.expiry = newToken.value.expiry > 0 ? newToken.value.expiry : 0
  const response = await HttpUtils.post('api/addToken', { desc: newToken.value.desc, expiry: newToken.value.expiry })
  if (response.success) {
    newToken.value.token = response.obj
    loadData()
    showNew.value = false
  }
  loading.value = false
}

const deleteToken = async (id: number) => {
  loading.value = true
  const response = await HttpUtils.post('api/deleteToken', { id: id })
  if (response.success) {
    loadData()
  }
  loading.value = false
}

const dateFormatted = (expiry: number): string => {
  if (expiry == 0) return i18n.global.t('unlimited')
  const date = new Date(expiry * 1000)
  return date.toLocaleString(locale, {
    year: 'numeric',
    month: '2-digit',
    day: '2-digit',
  })
}

watch(() => props.open, (v) => {
  if (v) {
    resetNewToken()
    showNew.value = false
    delConfirm.value = null
    loadData()
  }
})
</script>

<style scoped>
.token-alert {
  border: 1px solid color-mix(in srgb, var(--emerald) 35%, transparent);
  background: color-mix(in srgb, var(--emerald) 8%, transparent);
  border-radius: var(--radius-md);
  padding: 12px 14px;
  margin-bottom: 14px;
}
.add-box {
  border: 1px dashed var(--line-2);
  background: var(--surface-3);
  border-radius: var(--radius-md);
  padding: 14px;
}
</style>
