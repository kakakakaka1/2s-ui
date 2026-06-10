<template>
  <span
    class="ui-check"
    :class="{ on: checked }"
    @click="$emit('toggle')"
  >
    <Ico v-if="checked" name="check" :size="12" :stroke="2.4" />
  </span>
</template>

<script lang="ts" setup>
import Ico from './Ico.vue'

defineProps<{ checked?: boolean }>()
defineEmits<{ toggle: [] }>()
</script>

<style scoped>
.ui-check {
  width: 18px; height: 18px; flex: none;
  border-radius: 5px;
  border: 1.5px solid var(--line-2);
  background: transparent;
  display: inline-flex; align-items: center; justify-content: center;
  cursor: pointer; color: #fff;
  transition: all .12s;
}
.ui-check.on { border-color: var(--brand); background: var(--brand); }
</style>
