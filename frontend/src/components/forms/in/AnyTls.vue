<template>
  <div>
    <SectionLabel v-if="direction != 'out_json'" style="margin-bottom: 12px;">AnyTls</SectionLabel>
    <!-- Inbound -->
    <template v-if="direction == 'in'">
      <Field label="Padding scheme">
        <textarea class="input mono" rows="9" v-model="padding_scheme"></textarea>
      </Field>
    </template>
    <!-- Outbound / client side (out_json) -->
    <template v-else>
      <Field v-if="direction == 'out'" :label="$t('types.pw')">
        <input class="input mono" v-model="data.password" />
      </Field>
      <div class="grid2">
        <Field :label="$t('types.anytls.idleInterval') + ' (' + $t('date.s') + ')'">
          <input class="input mono" type="number" min="0" v-model.number="idleInterval" />
        </Field>
        <Field :label="$t('types.anytls.idleTimeout') + ' (' + $t('date.s') + ')'">
          <input class="input mono" type="number" min="0" v-model.number="idleTimeout" />
        </Field>
      </div>
      <Field :label="$t('types.anytls.minIdle')">
        <input class="input mono" type="number" min="0" v-model.number="minIdle" />
      </Field>
    </template>
  </div>
</template>

<script lang="ts" setup>
import { computed } from 'vue'
import Field from '@/components/ui/Field.vue'
import SectionLabel from '@/components/ui/SectionLabel.vue'

const props = defineProps<{ data: any; direction?: string }>()

const padding_scheme = computed({
  get: (): string => (props.data.padding_scheme?.length > 0 ? props.data.padding_scheme.join('\n') : ''),
  set: (v: string) => { props.data.padding_scheme = v.length > 0 ? v.split('\n') : undefined },
})

const idleInterval = computed({
  get: (): number =>
    props.data.idle_session_check_interval?.length > 0
      ? parseInt(props.data.idle_session_check_interval.replace('s', ''))
      : 30,
  set: (v: number) => { props.data.idle_session_check_interval = v && v >= 0 ? `${v}s` : undefined },
})

const idleTimeout = computed({
  get: (): number =>
    props.data.idle_session_timeout?.length > 0
      ? parseInt(props.data.idle_session_timeout.replace('s', ''))
      : 30,
  set: (v: number) => { props.data.idle_session_timeout = v && v >= 0 ? `${v}s` : undefined },
})

const minIdle = computed({
  get: (): number => (props.data.min_idle_session != undefined ? props.data.min_idle_session : 0),
  set: (v: number) => { props.data.min_idle_session = v > 0 ? v : undefined },
})
</script>
