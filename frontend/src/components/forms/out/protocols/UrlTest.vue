<template>
  <div style="display: flex; align-items: center; gap: 10px; margin-bottom: 12px;">
    <SectionLabel>URL Test</SectionLabel>
    <div style="flex: 1;" />
    <Pop :min-width="210" direction="down">
      <template #trigger="{ toggle }">
        <Btn variant="subtle" sm @click="toggle">{{ $t('types.lb.urlTestOptions') }}</Btn>
      </template>
      <div style="display: flex; flex-direction: column; gap: 2px; padding: 4px;">
        <div class="pop-item"><SwitchLabel v-model="optionUrl" :label="$t('types.lb.testUrl')" /></div>
        <div class="pop-item"><SwitchLabel v-model="optionInterval" :label="$t('types.lb.interval')" /></div>
        <div class="pop-item"><SwitchLabel v-model="optionTolerance" :label="$t('types.lb.tolerance')" /></div>
        <div class="pop-item"><SwitchLabel v-model="optionIdle" :label="$t('transport.idleTimeout')" /></div>
      </div>
    </Pop>
  </div>

  <Field :label="$t('pages.outbounds')">
    <MultiPick v-model="outboundsValue" :options="tags ?? []" />
  </Field>
  <Field v-if="optionUrl" :label="$t('types.lb.testUrl')">
    <input class="input mono" v-model="data.url" />
  </Field>
  <div class="grid2" style="margin-bottom: 15px;" v-if="optionInterval || optionTolerance || optionIdle">
    <Field v-if="optionInterval" :label="$t('types.lb.interval') + ' (' + $t('date.s') + ')'" :mb="0">
      <input class="input mono" type="number" min="3" v-model.number="interval" />
    </Field>
    <Field v-if="optionTolerance" :label="$t('types.lb.tolerance') + ' (' + $t('date.ms') + ')'" :mb="0">
      <input class="input mono" type="number" min="0" v-model.number="tolerance" />
    </Field>
    <Field v-if="optionIdle" :label="$t('transport.idleTimeout') + ' (' + $t('date.m') + ')'" :mb="0">
      <input class="input mono" type="number" min="0" v-model.number="idleTimeout" />
    </Field>
  </div>
  <div style="margin-bottom: 15px;">
    <CheckLabel v-model="interrupt" :label="$t('types.lb.interruptConn')" />
  </div>
</template>

<script lang="ts" setup>
import { computed } from 'vue'
import Field from '@/components/ui/Field.vue'
import Btn from '@/components/ui/Btn.vue'
import Pop from '@/components/ui/Pop.vue'
import SwitchLabel from '@/components/ui/SwitchLabel.vue'
import CheckLabel from '@/components/ui/CheckLabel.vue'
import SectionLabel from '@/components/ui/SectionLabel.vue'
import MultiPick from '@/components/ui/MultiPick.vue'

const props = defineProps<{ data: any; tags?: string[] }>()

const outboundsValue = computed({
  get: (): string[] => props.data.outbounds ?? [],
  set: (v: string[]) => { props.data.outbounds = v },
})
const interrupt = computed({
  get: () => props.data.interrupt_exist_connections ?? false,
  set: (v: boolean) => { props.data.interrupt_exist_connections = v },
})
const optionUrl = computed({
  get: (): boolean => props.data.url != undefined,
  set: (v: boolean) => { props.data.url = v ? 'https://www.gstatic.com/generate_204' : undefined },
})
const optionInterval = computed({
  get: (): boolean => props.data.interval != undefined,
  set: (v: boolean) => { props.data.interval = v ? '3s' : undefined },
})
const optionTolerance = computed({
  get: (): boolean => props.data.tolerance != undefined,
  set: (v: boolean) => { props.data.tolerance = v ? 50 : undefined },
})
const optionIdle = computed({
  get: (): boolean => props.data.idle_timeout != undefined,
  set: (v: boolean) => { props.data.idle_timeout = v ? '30m' : undefined },
})
const interval = computed({
  get: () => props.data.interval ? parseInt(props.data.interval.replace('s', '')) : 3,
  set: (v: number) => { props.data.interval = v > 0 ? v + 's' : '3s' },
})
const tolerance = computed({
  get: () => props.data.tolerance ? parseInt(props.data.tolerance) : 0,
  set: (v: number) => { props.data.tolerance = v > 0 ? v : 0 },
})
const idleTimeout = computed({
  get: () => props.data.idle_timeout ? parseInt(props.data.idle_timeout.replace('m', '')) : 30,
  set: (v: number) => { props.data.idle_timeout = v > 0 ? v + 'm' : '0m' },
})
</script>
