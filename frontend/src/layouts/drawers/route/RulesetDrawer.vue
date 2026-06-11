<template>
  <MDrawer
    :open="open"
    icon="rules"
    :color="form.type === 'remote' ? 'var(--violet)' : 'var(--cyan)'"
    :title="isNew ? $t('ui.rsNew') : form.tag"
    :sub="$t('ui.rsSub')"
    :save-label="isNew ? $t('ui.create') : $t('actions.save')"
    :width="500"
    @close="$emit('close')"
    @save="saveChanges"
  >
    <Field :label="$t('objects.tag')">
      <input class="input mono" v-model="form.tag" placeholder="geosite-ads" />
    </Field>

    <!-- local / remote -->
    <Segmented v-model="kind" block :options="[['local', $t('ruleset.local')], ['remote', $t('ruleset.remote')]]" />

    <Field :label="$t('ruleset.format')">
      <select class="input" v-model="form.format">
        <option value="source">source</option>
        <option value="binary">binary</option>
      </select>
    </Field>

    <template v-if="form.type === 'local'">
      <Field :label="$t('transport.path')">
        <input class="input mono" v-model="form.path" placeholder="/etc/sing-box/rules.srs" />
      </Field>
    </template>
    <template v-else>
      <Field label="URL">
        <input class="input mono" style="font-size: 12px;" v-model="form.url" placeholder="https://…/geosite-ads.srs" />
      </Field>
      <div class="grid2">
        <Field :label="$t('objects.outbound')">
          <select class="input" v-model="downloadDetour">
            <option value="">{{ $t('ui.none') }}</option>
            <option v-for="o in outTags" :key="o" :value="o">{{ o }}</option>
          </select>
        </Field>
        <Field :label="$t('ruleset.interval') + ' (' + $t('date.d') + ')'">
          <input class="input mono" type="number" min="0" v-model.number="updateIntervals" />
        </Field>
      </div>
    </template>
  </MDrawer>
</template>

<script lang="ts" setup>
import { computed, ref, watch } from 'vue'
import MDrawer from '@/components/ui/MDrawer.vue'
import Segmented from '@/components/ui/Segmented.vue'
import Field from '@/components/ui/Field.vue'
import RandomUtil from '@/plugins/randomUtil'
import { ruleset } from '@/types/rules'

const props = defineProps<{
  open: boolean
  index: number
  data: string
  outTags: string[]
}>()
const emit = defineEmits<{ close: []; save: [data: ruleset] }>()

const isNew = computed(() => props.index === -1)

const form = ref<ruleset>(<ruleset>{})

function init() {
  if (props.index !== -1) {
    form.value = <ruleset>JSON.parse(props.data)
  } else {
    form.value = <ruleset>{ type: 'local', tag: 'rs-' + RandomUtil.randomSeq(3), format: 'binary' }
  }
}
watch(() => props.open, (v) => { if (v) init() })

// switching type clears the fields of the other kind (legacy updateType)
const kind = computed({
  get: () => form.value.type,
  set: (t) => {
    form.value.type = <'local' | 'remote'>t
    if (t === 'local') {
      delete form.value.url
      delete form.value.download_detour
      delete form.value.update_interval
    } else {
      delete form.value.path
    }
  },
})

const downloadDetour = computed({
  get: () => form.value.download_detour ?? '',
  set: (v: string) => {
    if (v.length > 0) form.value.download_detour = v
    else delete form.value.download_detour
  },
})

const updateIntervals = computed({
  get: () => (form.value.update_interval !== undefined ? parseInt(form.value.update_interval.replace('d', '')) : 0),
  set: (v: number) => { form.value.update_interval = v > 0 ? v + 'd' : undefined },
})

function saveChanges() {
  emit('save', form.value)
}
</script>
