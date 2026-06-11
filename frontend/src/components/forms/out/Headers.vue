<template>
  <div style="margin-bottom: 15px;">
    <div style="display: flex; align-items: center; gap: 10px; margin-bottom: 10px;">
      <SectionLabel>{{ $t('objects.headers') }}</SectionLabel>
      <Btn variant="subtle" sm @click="addHeader"><Ico name="plus" :size="14" /> {{ $t('actions.add') }}</Btn>
    </div>
    <div v-for="(header, index) in hdrs" :key="index" style="display: flex; gap: 10px; align-items: center; margin-bottom: 10px;">
      <input
        class="input mono"
        :placeholder="$t('objects.key')"
        :value="header.name"
        @input="updateKey(index, ($event.target as HTMLInputElement).value)"
      />
      <input
        class="input mono"
        :placeholder="$t('objects.value')"
        :value="header.value"
        @input="updateValue(index, ($event.target as HTMLInputElement).value)"
      />
      <IconBtn name="trash" danger :title="$t('actions.del')" @click="delHeader(index)" />
    </div>
  </div>
</template>

<script lang="ts" setup>
import { computed } from 'vue'
import SectionLabel from '@/components/ui/SectionLabel.vue'
import Btn from '@/components/ui/Btn.vue'
import Ico from '@/components/ui/Ico.vue'
import IconBtn from '@/components/ui/IconBtn.vue'

type Header = { name: string; value: string }

const props = defineProps<{ data: any }>()

const hdrs = computed({
  get: (): Header[] => {
    const headers: Header[] = []
    const h = props.data.headers
    if (h) {
      Object.keys(h).forEach((key) => {
        if (Array.isArray(h[key])) {
          h[key].forEach((v: string) => headers.push({ name: key, value: v }))
        } else {
          headers.push({ name: key, value: h[key] })
        }
      })
    }
    return headers
  },
  set: (v: Header[]) => {
    if (v.length > 0) {
      const headers: any = {}
      v.forEach((h: Header) => {
        if (headers[h.name]) {
          if (Array.isArray(headers[h.name])) {
            headers[h.name].push(h.value)
          } else {
            headers[h.name] = [headers[h.name], h.value]
          }
        } else {
          headers[h.name] = h.value
        }
      })
      props.data.headers = headers
    } else {
      props.data.headers = undefined
    }
  },
})

const addHeader = () => { hdrs.value = [...hdrs.value, { name: 'Host', value: '' }] }
const delHeader = (i: number) => {
  const h = hdrs.value
  h.splice(i, 1)
  hdrs.value = h
}
const updateKey = (i: number, k: string) => {
  const h = hdrs.value
  h[i].name = k
  hdrs.value = h
}
const updateValue = (i: number, v: string) => {
  const h = hdrs.value
  h[i].value = v
  hdrs.value = h
}
</script>
