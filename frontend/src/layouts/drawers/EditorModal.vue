<template>
  <Modal :open="open" :title="title" :width="760" @close="$emit('close')">
    <div style="display: flex; flex-direction: column; flex: 1;">
      <textarea
        ref="ta"
        class="input mono editor-area"
        spellcheck="false"
        :value="draft"
        @input="draft = ($event.target as HTMLTextAreaElement).value"
      />
    </div>
    <template #footer>
      <Btn sm @click="$emit('close')">{{ $t('actions.cancel') }}</Btn>
      <Btn variant="primary" sm @click="$emit('save', draft); $emit('close')">{{ $t('actions.save') }}</Btn>
    </template>
  </Modal>
</template>

<script lang="ts" setup>
import { ref, watch } from 'vue'
import Modal from '@/components/ui/Modal.vue'
import Btn from '@/components/ui/Btn.vue'

const props = defineProps<{
  open: boolean
  title: string
  content: string
}>()
defineEmits<{ close: []; save: [value: string] }>()

const draft = ref('')
watch(() => props.open, (v) => { if (v) draft.value = props.content })
</script>

<style scoped>
.editor-area {
  border: none;
  border-radius: 0;
  min-height: 380px;
  resize: vertical;
  line-height: 1.6;
  padding: 16px;
  font-size: 12.5px;
  background: var(--surface-3);
}
.editor-area:focus { box-shadow: none; border: none; background: var(--surface-3); }
</style>
