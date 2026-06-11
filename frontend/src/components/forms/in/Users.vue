<template>
  <div>
    <SectionLabel style="margin-bottom: 12px;">{{ $t('pages.clients') }}</SectionLabel>
    <Field>
      <select class="input" v-model="data.model" @change="data.values = []">
        <option v-for="m in initUsersModels" :key="m.value" :value="m.value">{{ m.title }}</option>
      </select>
    </Field>
    <Field v-if="data.model == 'group'" :label="$t('client.group')">
      <MultiPick v-model="data.values" :options="groupNames" />
    </Field>
    <Field v-if="data.model == 'client'" :label="$t('pages.clients')">
      <div style="display: flex; gap: 7px; flex-wrap: wrap;">
        <button
          v-for="c in clientItems"
          :key="c.value"
          type="button"
          class="chip mono"
          :style="{
            cursor: 'pointer', height: '30px', padding: '0 12px',
            color: data.values.includes(c.value) ? '#fff' : 'var(--text-2)',
            background: data.values.includes(c.value) ? 'var(--brand)' : 'var(--surface-3)',
            borderColor: data.values.includes(c.value) ? 'var(--brand)' : 'var(--line)',
          }"
          @click="toggleClient(c.value)"
        >{{ c.title }}</button>
      </div>
    </Field>
  </div>
</template>

<script lang="ts" setup>
import { computed } from 'vue'
import { useI18n } from 'vue-i18n'
import Field from '@/components/ui/Field.vue'
import SectionLabel from '@/components/ui/SectionLabel.vue'
import MultiPick from '@/components/ui/MultiPick.vue'

const props = defineProps<{ data: any; clients: any[] }>()
const { t } = useI18n({ useScope: 'global' })

const initUsersModels = computed(() => [
  { title: t('none'), value: 'none' },
  { title: t('all'), value: 'all' },
  { title: t('client.group'), value: 'group' },
  { title: t('pages.clients'), value: 'client' },
])

const clientItems = computed(() =>
  props.clients.map((c: any) => ({ title: c.name, value: c.id }))
)

const groupNames = computed((): string[] =>
  Array.from(new Set(props.clients.map((c: any) => c.group)))
)

const toggleClient = (id: number) => {
  props.data.values = props.data.values.includes(id)
    ? props.data.values.filter((x: number) => x !== id)
    : [...props.data.values, id]
}
</script>
