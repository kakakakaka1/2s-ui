<template>
  <div class="grid2" style="margin-bottom: 15px;">
    <Field :label="$t('types.tor.execPath')" :mb="0">
      <input class="input mono" v-model="data.executable_path" />
    </Field>
    <Field :label="$t('types.tor.dataDir')" :mb="0">
      <input class="input mono" v-model="data.data_directory" />
    </Field>
  </div>
  <Field :label="$t('types.tor.extArgs') + ' ' + $t('commaSeparated')">
    <input class="input mono" v-model="extraArgs" />
  </Field>
  <div style="display: flex; align-items: center; gap: 10px; margin-bottom: 10px;">
    <SectionLabel>Torrc</SectionLabel>
    <Btn variant="subtle" sm @click="addTorrcOption"><Ico name="plus" :size="14" /> {{ $t('actions.add') }}</Btn>
  </div>
  <div v-for="(torrc, index) in torrcOptions" :key="index" style="display: flex; gap: 10px; align-items: center; margin-bottom: 10px;">
    <input
      class="input mono"
      :placeholder="$t('objects.key')"
      :value="torrc.name"
      @input="updateKey(index, ($event.target as HTMLInputElement).value)"
    />
    <input
      class="input mono"
      :placeholder="$t('objects.value')"
      :value="torrc.value"
      @input="updateValue(index, ($event.target as HTMLInputElement).value)"
    />
    <IconBtn name="trash" danger :title="$t('actions.del')" @click="delTorrcOption(index)" />
  </div>
</template>

<script lang="ts" setup>
import { computed } from 'vue'
import Field from '@/components/ui/Field.vue'
import Btn from '@/components/ui/Btn.vue'
import Ico from '@/components/ui/Ico.vue'
import IconBtn from '@/components/ui/IconBtn.vue'
import SectionLabel from '@/components/ui/SectionLabel.vue'

type TorrcOption = { name: string; value: string }

const props = defineProps<{ data: any }>()

const torrcOptions = computed({
  get: (): TorrcOption[] => {
    const options: TorrcOption[] = []
    const h = props.data.torrc
    if (h) {
      Object.keys(h).forEach((key) => {
        if (Array.isArray(h[key])) {
          h[key].forEach((v: string) => options.push({ name: key, value: v }))
        } else {
          options.push({ name: key, value: h[key] })
        }
      })
    }
    return options
  },
  set: (v: TorrcOption[]) => {
    if (v.length > 0) {
      const torrc: any = {}
      v.forEach((h: TorrcOption) => {
        if (torrc[h.name]) {
          if (Array.isArray(torrc[h.name])) {
            torrc[h.name].push(h.value)
          } else {
            torrc[h.name] = [torrc[h.name], h.value]
          }
        } else {
          torrc[h.name] = h.value
        }
      })
      props.data.torrc = torrc
    } else {
      props.data.torrc = undefined
    }
  },
})

const extraArgs = computed({
  get: () => props.data.extra_args?.join(','),
  set: (v: string) => { props.data.extra_args = v.length > 0 ? v.split(',') : undefined },
})

const addTorrcOption = () => { torrcOptions.value = [...torrcOptions.value, { name: '', value: '' }] }
const delTorrcOption = (i: number) => {
  const h = torrcOptions.value
  h.splice(i, 1)
  torrcOptions.value = h
}
const updateKey = (i: number, k: string) => {
  const h = torrcOptions.value
  h[i].name = k
  torrcOptions.value = h
}
const updateValue = (i: number, v: string) => {
  const h = torrcOptions.value
  h[i].value = v
  torrcOptions.value = h
}
</script>
