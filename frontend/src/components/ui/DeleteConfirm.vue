<template>
  <Modal :open="open" :title="$t('actions.del')" :width="380" @close="$emit('close')">
    <div style="padding: 18px; font-size: 13.5px;">{{ $t('confirm') }}</div>
    <template #footer>
      <Btn @click="$emit('close')">{{ $t('no') }}</Btn>
      <Btn style="color: var(--rose);" :loading="loading" @click="$emit('confirm')">
        <Ico name="trash" :size="15" /> {{ $t('yes') }}
      </Btn>
    </template>
  </Modal>
</template>

<script lang="ts" setup>
import Modal from './Modal.vue'
import Btn from './Btn.vue'
import Ico from './Ico.vue'

// Modal already owns the overlay stack and scroll lock — nothing here touches
// pushOverlay. loading is left optional: the views that delete locally (Rules,
// Dns) have no pending state, and Btn treats undefined as not-loading.
defineProps<{ open: boolean; loading?: boolean }>()
defineEmits<{ close: []; confirm: [] }>()
</script>
