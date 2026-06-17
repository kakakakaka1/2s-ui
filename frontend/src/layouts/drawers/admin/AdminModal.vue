<template>
  <MDrawer
    :open="open"
    icon="admins"
    color="var(--brand)"
    :title="$t('admin.changeCred')"
    :sub="user?.username ?? ''"
    :width="420"
    :loading="loading"
    @close="$emit('close')"
    @save="save"
  >
    <Field :label="$t('admin.oldPass')">
      <input class="input" type="password" v-model="form.oldPass" autocomplete="current-password" />
      <span v-if="tried && form.oldPass === ''" class="err">{{ $t('login.pwRules') }}</span>
    </Field>
    <Field :label="$t('admin.newUname')">
      <input class="input" v-model="form.newUsername" autocomplete="username" />
      <span v-if="tried && form.newUsername === ''" class="err">{{ $t('login.unRules') }}</span>
    </Field>
    <Field :label="$t('admin.newPass')">
      <input class="input" type="password" v-model="form.newPass" autocomplete="new-password" />
      <span v-if="tried && form.newPass === ''" class="err">{{ $t('login.pwRules') }}</span>
    </Field>
  </MDrawer>
</template>

<script lang="ts" setup>
import { ref, watch } from 'vue'
import MDrawer from '@/components/ui/MDrawer.vue'
import Field from '@/components/ui/Field.vue'

const props = defineProps<{
  open: boolean
  user: any
  loading?: boolean
}>()
const emit = defineEmits<{ close: []; save: [data: any] }>()

const tried = ref(false)
const form = ref({ id: 0, oldPass: '', newUsername: '', newPass: '' })

watch(() => props.open, (v) => {
  if (v) {
    form.value = { id: props.user?.id ?? 0, oldPass: '', newUsername: '', newPass: '' }
    tried.value = false
  }
})

const save = () => {
  tried.value = true
  if (form.value.oldPass === '' || form.value.newUsername === '' || form.value.newPass === '') return
  emit('save', { ...form.value })
}
</script>

<style scoped>
.err { font-size: 11px; color: var(--rose); }
</style>
