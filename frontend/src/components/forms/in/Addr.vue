<template>
  <div>
    <div style="display: flex; justify-content: flex-end; margin-bottom: 8px;">
      <Pop :min-width="220">
        <template #trigger="{ toggle }">
          <button type="button" class="btn btn-ghost btn-sm" @click="toggle">
            <Ico name="settings" :size="14" /> {{ $t('in.mdOption') }}
          </button>
        </template>
        <div style="display: flex; flex-direction: column; gap: 2px; padding: 5px;">
          <SwitchLabel v-model="optionRemark" :label="$t('in.remark')" />
          <SwitchLabel v-if="hasTls" v-model="optionTLS" :label="$t('objects.tls')" />
        </div>
      </Pop>
    </div>
    <div class="grid2">
      <Field :label="$t('out.addr')">
        <input class="input mono" required v-model="addr.server" />
      </Field>
      <Field :label="$t('out.port')">
        <input class="input mono" type="number" required v-model.number="addr.server_port" />
      </Field>
    </div>
    <Field v-if="optionRemark" :label="$t('in.remark')">
      <input class="input" v-model="addr.remark" />
    </Field>
    <OutTLS v-if="optionTLS" :outbound="addr" />
  </div>
</template>

<script lang="ts" setup>
import { computed } from 'vue'
import Field from '@/components/ui/Field.vue'
import SwitchLabel from '@/components/ui/SwitchLabel.vue'
import Pop from '@/components/ui/Pop.vue'
import Ico from '@/components/ui/Ico.vue'
import OutTLS from './OutTLS.vue'

const props = defineProps<{ addr: any; hasTls?: boolean }>()

const optionTLS = computed({
  get: (): boolean => props.addr.tls != undefined,
  set: (v: boolean) => { props.addr.tls = v ? { enabled: true } : undefined },
})

const optionRemark = computed({
  get: (): boolean => props.addr.remark != undefined,
  set: (v: boolean) => { props.addr.remark = v ? '' : undefined },
})
</script>
