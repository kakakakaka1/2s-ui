<template>
  <Modal :open="visible" :title="$t('main.backup.title')" :width="480" @close="$emit('close')">
    <div style="padding: 18px 20px;">
      <div style="display: flex; gap: 24px; flex-wrap: wrap; margin-bottom: 18px;">
        <CheckLabel v-model="exclStats" :label="$t('main.backup.exclStats')" />
        <CheckLabel v-model="exclChanges" :label="$t('main.backup.exclChanges')" />
      </div>
      <div style="display: flex; gap: 10px; flex-wrap: wrap;">
        <Btn variant="primary" style="flex: 1;" @click="backup"><Ico name="download" :size="15" /> {{ $t('main.backup.backup') }}</Btn>
        <Btn style="flex: 1;" @click="restore"><Ico name="upload" :size="15" /> {{ $t('main.backup.restore') }}</Btn>
      </div>
      <hr class="form-divider" style="margin: 18px 0;" />
      <Btn style="width: 100%;" @click="config"><Ico name="copy" :size="15" /> {{ $t('main.backup.sbConfig') }}</Btn>
    </div>
  </Modal>
</template>

<script lang="ts" setup>
import { ref, watch } from 'vue'
import HttpUtils from '@/plugins/httputil'
import Modal from '@/components/ui/Modal.vue'
import Btn from '@/components/ui/Btn.vue'
import Ico from '@/components/ui/Ico.vue'
import CheckLabel from '@/components/ui/CheckLabel.vue'

const props = defineProps<{ visible: boolean }>()
const emit = defineEmits<{ close: [] }>()

const exclStats = ref(true)
const exclChanges = ref(true)

watch(() => props.visible, (v) => {
  if (v) {
    exclStats.value = true
    exclChanges.value = true
  }
})

const backup = () => {
  const exclude = [exclStats.value ? 'stats' : '', exclChanges.value ? 'changes' : ''].filter(Boolean)
  window.location.href = 'api/getdb' + (exclude.length > 0 ? '?exclude=' + exclude.join(',') : '')
}
const config = () => {
  window.location.href = 'api/singbox-config'
}
const restore = () => {
  const fileInput = document.createElement('input')
  fileInput.type = 'file'
  fileInput.accept = '.db'
  fileInput.addEventListener('change', async (event: Event) => {
    const inputElement = event.target as HTMLInputElement
    const dbFile = inputElement.files ? inputElement.files[0] : null
    if (dbFile) {
      const formData = new FormData()
      formData.append('db', dbFile)
      emit('close')
      const uploadMsg = await HttpUtils.post('api/importdb', formData, {
        headers: { 'Content-Type': 'multipart/form-data' },
      })
      if (uploadMsg.success) {
        await new Promise((resolve) => setTimeout(resolve, 1000))
        location.reload()
      }
    }
  })
  fileInput.click()
}
</script>
