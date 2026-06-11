<template>
  <div style="display: flex; align-items: center; gap: 10px; margin-bottom: 12px;">
    <SectionLabel>DERP</SectionLabel>
    <div style="flex: 1;" />
    <Pop :min-width="220" direction="down">
      <template #trigger="{ toggle }">
        <Btn variant="subtle" sm @click="toggle">{{ $t('types.derp.options') }}</Btn>
      </template>
      <div style="display: flex; flex-direction: column; gap: 2px; padding: 4px;">
        <div class="pop-item"><SwitchLabel v-model="optionVerifyCE" :label="$t('types.derp.verifyClientEndpoint')" /></div>
        <div class="pop-item"><SwitchLabel v-model="optionVerifyCU" :label="$t('types.derp.verifyClientUrl')" /></div>
        <div class="pop-item"><SwitchLabel v-model="optionHome" :label="$t('pages.home')" /></div>
        <div class="pop-item"><SwitchLabel v-model="optionMesh" :label="$t('types.derp.meshWith')" /></div>
        <div class="pop-item"><SwitchLabel v-model="optionStun" :label="$t('types.derp.stun')" /></div>
      </div>
    </Pop>
  </div>

  <Field :label="$t('types.derp.configPath')">
    <input class="input mono" v-model="data.config_path" />
  </Field>
  <Field v-if="optionHome" :label="$t('pages.home')">
    <input class="input mono" v-model="data.home" placeholder="blank | http[s]://example.com:port/path" />
  </Field>
  <Field v-if="optionVerifyCE" :label="$t('types.derp.verifyClientEndpoint')">
    <ChipSelect v-model="verifyClientEndpoint" :options="(tsTags ?? []).map((t: string) => ({ title: t, value: t }))" />
  </Field>

  <template v-if="optionVerifyCU">
    <div style="display: flex; align-items: center; gap: 10px; margin-bottom: 12px;">
      <SectionLabel>{{ $t('types.derp.verifyClientUrl') }}</SectionLabel>
      <Btn variant="subtle" sm @click="data.verify_client_url.push({ url: '' })">
        <Ico name="plus" :size="14" /> {{ $t('actions.add') }}
      </Btn>
    </div>
    <div
      v-for="(clientUrl, index) in data.verify_client_url"
      :key="index"
      class="card"
      style="padding: 14px; margin-bottom: 12px; background: var(--surface-2);"
    >
      <div style="display: flex; align-items: center; margin-bottom: 10px;">
        <Chip>{{ Number(index) + 1 }}</Chip>
        <div style="flex: 1;" />
        <Btn variant="subtle" icon sm :title="$t('actions.del')" @click="data.verify_client_url.splice(index, 1)">
          <Ico name="trash" :size="14" />
        </Btn>
      </div>
      <Field :label="$t('types.derp.verifyClientUrl')">
        <input class="input mono" v-model="clientUrl.url" />
      </Field>
      <Dial :dial="clientUrl" />
    </div>
  </template>

  <template v-if="optionMesh">
    <div style="display: flex; align-items: center; gap: 10px; margin-bottom: 12px;">
      <SectionLabel>{{ $t('types.derp.meshWith') }}</SectionLabel>
      <Btn variant="subtle" sm @click="data.mesh_with.push({ tls: {} })">
        <Ico name="plus" :size="14" /> {{ $t('actions.add') }}
      </Btn>
    </div>
    <div
      v-for="(mesh, index) in data.mesh_with"
      :key="index"
      class="card"
      style="padding: 14px; margin-bottom: 12px; background: var(--surface-2);"
    >
      <div style="display: flex; align-items: center; margin-bottom: 10px;">
        <Chip>{{ Number(index) + 1 }}</Chip>
        <div style="flex: 1;" />
        <Btn variant="subtle" icon sm :title="$t('actions.del')" @click="data.mesh_with.splice(index, 1)">
          <Ico name="trash" :size="14" />
        </Btn>
      </div>
      <div class="grid2">
        <Field :label="$t('out.addr')">
          <input class="input mono" v-model="mesh.server" />
        </Field>
        <Field :label="$t('out.port')">
          <input class="input mono" type="number" v-model.number="mesh.server_port" />
        </Field>
      </div>
      <Field :label="$t('transport.host')">
        <input class="input mono" v-model="mesh.host" />
      </Field>
      <Dial :dial="mesh" />
      <OutTLS :outbound="mesh" />
    </div>
    <div style="margin-bottom: 12px;">
      <Segmented
        v-model="usePskText"
        :options="[[0, $t('types.derp.meshPsk')], [1, $t('types.derp.meshPskFile')]]"
        block
        @update:model-value="switchPskMode"
      />
    </div>
    <Field v-if="usePskText == 1" :label="$t('types.derp.meshPskFile')">
      <input class="input mono" v-model="data.mesh_psk_file" />
    </Field>
    <Field v-else :label="$t('types.derp.meshPsk')">
      <input class="input mono" v-model="data.mesh_psk" />
    </Field>
  </template>

  <template v-if="optionStun">
    <div style="display: flex; align-items: center; gap: 10px; margin-bottom: 12px;">
      <SectionLabel>{{ $t('types.derp.stun') }}</SectionLabel>
    </div>
    <div class="card" style="padding: 14px; margin-bottom: 12px; background: var(--surface-2);">
      <Listen :data="data.stun" :inTags="inTags" />
    </div>
  </template>
</template>

<script lang="ts" setup>
import { computed, ref } from 'vue'
import Field from '@/components/ui/Field.vue'
import Btn from '@/components/ui/Btn.vue'
import Ico from '@/components/ui/Ico.vue'
import Chip from '@/components/ui/Chip.vue'
import Pop from '@/components/ui/Pop.vue'
import SwitchLabel from '@/components/ui/SwitchLabel.vue'
import Segmented from '@/components/ui/Segmented.vue'
import SectionLabel from '@/components/ui/SectionLabel.vue'
import ChipSelect from '@/components/ui/ChipSelect.vue'
import Dial from '../Dial.vue'
import OutTLS from '../OutTLS.vue'
import Listen from '../Listen.vue'

const props = defineProps<{ data: any; tsTags?: string[]; inTags?: string[] }>()

const usePskText = ref<string | number>(props.data.mesh_psk == undefined ? 1 : 0)

const switchPskMode = (v: string | number) => {
  if (v == 0) delete props.data.mesh_psk_file
  else delete props.data.mesh_psk
}

const verifyClientEndpoint = computed({
  get: (): string[] => props.data.verify_client_endpoint ?? [],
  set: (v: string[]) => { props.data.verify_client_endpoint = v },
})

const optionVerifyCE = computed({
  get: () => props.data.verify_client_endpoint != undefined,
  set: (v: boolean) => { props.data.verify_client_endpoint = v ? [] : undefined },
})
const optionVerifyCU = computed({
  get: () => props.data.verify_client_url != undefined,
  set: (v: boolean) => { props.data.verify_client_url = v ? [{ url: '' }] : undefined },
})
const optionHome = computed({
  get: () => props.data.home != undefined,
  set: (v: boolean) => { props.data.home = v ? '' : undefined },
})
const optionMesh = computed({
  get: () => props.data.mesh_with != undefined,
  set: (v: boolean) => {
    if (v) {
      props.data.mesh_with = [{ tls: {} }]
      delete props.data.mesh_psk_file
      props.data.mesh_psk = ''
      usePskText.value = 0
    } else {
      delete props.data.mesh_with
      delete props.data.mesh_psk_file
      delete props.data.mesh_psk
    }
  },
})
const optionStun = computed({
  get: () => props.data.stun != undefined,
  set: (v: boolean) => { props.data.stun = v ? { enabled: true } : undefined },
})
</script>
