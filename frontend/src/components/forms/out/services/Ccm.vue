<template>
  <div style="display: flex; align-items: center; gap: 10px; margin-bottom: 12px;">
    <SectionLabel>CCM (Claude Code Multiplexer)</SectionLabel>
  </div>
  <div class="grid2">
    <Field :label="$t('types.ccm.credentialPath')">
      <input class="input mono" v-model="data.credential_path" />
    </Field>
    <Field :label="$t('types.ccm.usagesPath')">
      <input class="input mono" v-model="data.usages_path" />
    </Field>
  </div>
  <Field :label="$t('dial.detourText')">
    <select class="input" v-model="data.detour">
      <option v-for="tag in outTags" :key="tag" :value="tag">{{ tag }}</option>
    </select>
  </Field>

  <div style="display: flex; align-items: center; gap: 10px; margin-bottom: 12px;">
    <SectionLabel>{{ $t('types.ccm.users') }}</SectionLabel>
    <Btn variant="subtle" sm @click="addUser"><Ico name="plus" :size="14" /> {{ $t('actions.add') }}</Btn>
  </div>
  <div
    v-for="(user, index) in (data.users || [])"
    :key="index"
    class="card"
    style="padding: 14px; margin-bottom: 12px; background: var(--surface-2);"
  >
    <div style="display: flex; align-items: center; margin-bottom: 10px;">
      <Chip>{{ Number(index) + 1 }}</Chip>
      <div style="flex: 1;" />
      <Btn variant="subtle" icon sm :title="$t('actions.del')" @click="delUser(index)">
        <Ico name="trash" :size="14" />
      </Btn>
    </div>
    <div class="grid2">
      <Field :label="$t('types.ccm.userName')" :mb="0">
        <input class="input mono" v-model="user.name" />
      </Field>
      <Field :label="$t('types.ccm.userToken')" :mb="0">
        <input class="input mono" type="password" v-model="user.token" />
      </Field>
    </div>
  </div>
</template>

<script lang="ts" setup>
import { computed } from 'vue'
import Data from '@/store/modules/data'
import Field from '@/components/ui/Field.vue'
import Btn from '@/components/ui/Btn.vue'
import Ico from '@/components/ui/Ico.vue'
import Chip from '@/components/ui/Chip.vue'
import SectionLabel from '@/components/ui/SectionLabel.vue'

const props = defineProps<{ data: any }>()

const outTags = computed(() => [
  ...Data().outbounds?.map((o: any) => o.tag) ?? [],
  ...Data().endpoints?.map((e: any) => e.tag) ?? [],
])

const addUser = () => {
  if (!props.data.users) props.data.users = []
  props.data.users.push({ name: '', token: '' })
}
const delUser = (i: number | string) => {
  props.data.users?.splice(Number(i), 1)
}
</script>
