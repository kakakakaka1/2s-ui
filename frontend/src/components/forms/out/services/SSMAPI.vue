<template>
  <div style="display: flex; align-items: center; gap: 10px; margin-bottom: 12px;">
    <SectionLabel>Shadowsocks API</SectionLabel>
    <Btn variant="subtle" sm @click="addServer"><Ico name="plus" :size="14" /> {{ $t('actions.add') }}</Btn>
  </div>
  <div v-for="(server, index) in servers" :key="index" style="display: flex; gap: 10px; align-items: flex-end; margin-bottom: 10px;">
    <Field :label="$t('transport.path')" :mb="0" style="flex: 1;">
      <input
        class="input mono"
        :value="server.name"
        @input="updateKey(index, ($event.target as HTMLInputElement).value)"
      />
    </Field>
    <Field :label="$t('objects.inbound')" :mb="0" style="flex: 1;">
      <select class="input" :value="server.value" @change="updateValue(index, ($event.target as HTMLSelectElement).value)">
        <option v-for="tag in ssTags ?? []" :key="tag" :value="tag">{{ tag }}</option>
      </select>
    </Field>
    <IconBtn name="trash" danger :title="$t('actions.del')" style="margin-bottom: 7px;" @click="delServer(index)" />
  </div>
</template>

<script lang="ts" setup>
import { computed } from 'vue'
import Field from '@/components/ui/Field.vue'
import Btn from '@/components/ui/Btn.vue'
import Ico from '@/components/ui/Ico.vue'
import IconBtn from '@/components/ui/IconBtn.vue'
import SectionLabel from '@/components/ui/SectionLabel.vue'

type Server = { name: string; value: string }

const props = defineProps<{ data: any; ssTags?: string[] }>()

const servers = computed({
  get: (): Server[] => {
    const list: Server[] = []
    const h = props.data.servers
    if (h) {
      Object.keys(h).forEach((key) => {
        if (Array.isArray(h[key])) {
          h[key].forEach((v: string) => list.push({ name: key, value: v }))
        } else {
          list.push({ name: key, value: h[key] })
        }
      })
    }
    return list
  },
  set: (v: Server[]) => {
    if (v.length > 0) {
      const obj: any = {}
      v.forEach((h: Server) => {
        if (obj[h.name]) {
          if (Array.isArray(obj[h.name])) {
            obj[h.name].push(h.value)
          } else {
            obj[h.name] = [obj[h.name], h.value]
          }
        } else {
          obj[h.name] = h.value
        }
      })
      props.data.servers = obj
    } else {
      props.data.servers = undefined
    }
  },
})

const addServer = () => {
  servers.value = [...servers.value, { name: '/ss' + servers.value.length, value: (props.ssTags ?? [])[0] || '' }]
}
const delServer = (i: number) => {
  const h = servers.value
  h.splice(i, 1)
  servers.value = h
}
const updateKey = (i: number, k: string) => {
  const h = servers.value
  h[i].name = k
  servers.value = h
}
const updateValue = (i: number, v: string) => {
  const h = servers.value
  h[i].value = v
  servers.value = h
}
</script>
