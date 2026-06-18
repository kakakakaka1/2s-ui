<template>
  <Drawer :open="visible" :width="680" @close="$emit('close')">
    <template #title>
      <div style="display: flex; align-items: center; gap: 12px; min-width: 0;">
        <Avatar v-if="!isNew" :name="client.name" :size="38" />
        <div
          v-else
          style="width: 38px; height: 38px; border-radius: 9px; flex: none; background: var(--brand-soft); color: var(--brand); display: flex; align-items: center; justify-content: center;"
        ><Ico name="plus" :size="20" /></div>
        <div style="min-width: 0;">
          <div style="font-size: 16px; font-weight: 800; white-space: nowrap; overflow: hidden; text-overflow: ellipsis;">
            {{ isNew ? $t('ui.newClient') : client.name }}
          </div>
          <div style="font-size: 12px; color: var(--text-3); white-space: nowrap; overflow: hidden; text-overflow: ellipsis;">{{ subTitle }}</div>
        </div>
      </div>
    </template>

    <div v-if="loading" style="padding: 40px 0; text-align: center; color: var(--text-3); font-size: 13px;">
      {{ $t('loading') }}
    </div>
    <template v-else>
      <Tabs v-model="tab" :tabs="tabs" />

      <!-- ===================== General ===================== -->
      <div v-if="tab === 'general'">
        <MSwitchRow v-model="client.enable" :label="$t('ui.accountEnabled')" :desc="$t('ui.cannotConnect')" />

        <Field :label="$t('client.name')">
          <input class="input" v-model="client.name" />
        </Field>
        <Field :label="$t('client.desc')">
          <input class="input" v-model="client.desc" />
        </Field>
        <Field :label="$t('client.group')">
          <div v-if="groupItems.length > 0" style="display: flex; gap: 7px; flex-wrap: wrap; margin-bottom: 8px;">
            <button
              v-for="g in groupItems"
              :key="g"
              type="button"
              class="chip"
              :style="groupChipStyle(g)"
              @click="client.group = g"
            >{{ g }}</button>
          </div>
          <input class="input" v-model="client.group" />
        </Field>
      </div>

      <!-- ================= Inbounds + Limits ================= -->
      <div v-else-if="tab === 'access'" class="access-grid">
        <!-- Inbounds -->
        <section style="min-width: 0;">
          <div class="access-head">
            <SectionLabel>{{ $t('ui.nav.inbounds') }}</SectionLabel>
            <div style="flex: 1;" />
            <Btn sm variant="subtle" @click="setAllInbounds">{{ $t('all') }}</Btn>
          </div>
          <div style="font-size: 12.5px; color: var(--text-3); margin-bottom: 11px;">{{ $t('ui.selectInbounds') }}</div>
          <div class="inb-list hide-scroll">
            <label
              v-for="it in inboundItems"
              :key="it.id"
              class="inb-item"
              :class="{ on: client.inbounds.includes(it.id) }"
              @click.prevent="toggleInbound(it.id)"
            >
              <span class="inb-dot" :class="{ online: it.online }" />
              <div style="flex: 1; min-width: 0;">
                <div style="font-weight: 700; font-size: 13px; white-space: nowrap; overflow: hidden; text-overflow: ellipsis;">{{ it.tag }}</div>
                <div class="mono" style="font-size: 11.5px; color: var(--text-3);">{{ it.type }} · :{{ it.port }}</div>
              </div>
              <Check :checked="client.inbounds.includes(it.id)" />
            </label>
          </div>
        </section>

        <!-- Limits -->
        <section style="min-width: 0;">
          <div class="access-head">
            <SectionLabel>{{ $t('ui.limits') }}</SectionLabel>
          </div>

          <div v-if="!isNew" class="card usage-card">
            <div style="display: flex; align-items: center; margin-bottom: 10px;">
              <div style="font-size: 12px; color: var(--text-3); font-weight: 600;">{{ $t('ui.currentUsage') }}</div>
              <div v-if="percent > 0" class="mono" dir="ltr" style="margin-inline-start: auto; font-size: 12px; color: var(--text-3);">{{ percent }}%</div>
            </div>
            <div style="display: flex; gap: 22px; align-items: center;">
              <div>
                <div style="font-size: 11px; color: var(--violet); font-weight: 700;">↑ {{ $t('ui.upload') }}</div>
                <div class="mono" style="font-size: 17px; font-weight: 700;">{{ up }}</div>
              </div>
              <div>
                <div style="font-size: 11px; color: var(--cyan); font-weight: 700;">↓ {{ $t('ui.download') }}</div>
                <div class="mono" style="font-size: 17px; font-weight: 700;">{{ down }}</div>
              </div>
              <Btn sm variant="subtle" style="margin-inline-start: auto;" @click="resetUsage">
                <Ico name="refresh" :size="14" /> {{ $t('reset') }}
              </Btn>
            </div>
            <div v-if="client.volume > 0" style="margin-top: 12px;">
              <BarMini :value="percent" :color="percentColor" :height="5" />
            </div>
          </div>

          <Field :label="$t('ui.dataLimit')" :hint="$t('ui.unlimitedHint')">
            <div style="display: flex; gap: 8px;">
              <input class="input mono" type="number" min="0" v-model.number="volumeGiB" />
              <div class="input suffix-box">{{ $t('stats.GB') }}</div>
            </div>
          </Field>

          <div class="switch-stack">
            <div :style="client.up + client.down > 0 ? { opacity: 0.5, pointerEvents: 'none' } : undefined">
              <SwitchLabel v-model="delayStart" :label="$t('client.delayStart')" />
            </div>
            <SwitchLabel v-model="autoReset" :label="$t('client.autoReset')" />
          </div>

          <Field
            v-if="!(client.delayStart && !client.autoReset)"
            :label="$t('ui.expiryDate')"
            :hint="$t('ui.noExpiryHint')"
          >
            <DateTimeInput v-model="client.expiry" />
          </Field>

          <Field v-if="client.autoReset || client.delayStart" :label="$t('client.resetDays')">
            <input class="input mono" type="number" min="1" v-model.number="resetDays" />
          </Field>

          <template v-if="!isNew && client.autoReset">
            <hr class="form-divider" />
            <div class="grid2">
              <div>
                <div style="font-size: 11px; color: var(--text-3); font-weight: 600; margin-bottom: 3px;">{{ $t('client.nextReset') }}</div>
                <div class="mono" dir="ltr" style="font-size: 13px; font-weight: 600;">{{ nextResetFormatted }}</div>
              </div>
              <div>
                <div style="font-size: 11px; color: var(--text-3); font-weight: 600; margin-bottom: 3px;">{{ $t('main.stats.totalUsage') }}</div>
                <div class="mono" dir="ltr" style="font-size: 13px; font-weight: 600;">↑ {{ totalUp }} / ↓ {{ totalDown }}</div>
              </div>
            </div>
          </template>
        </section>
      </div>

      <!-- ===================== UUID / Password ===================== -->
      <div v-else-if="tab === 'keys'">
        <div class="access-head">
          <SectionLabel>{{ $t('ui.uuidPass') }}</SectionLabel>
          <div style="flex: 1;" />
          <Btn sm variant="subtle" @click="shuffle()">
            <Ico name="refresh" :size="14" /> {{ $t('reset') + ' - ' + $t('all') }}
          </Btn>
        </div>
        <div style="font-size: 12.5px; color: var(--text-3); margin-bottom: 16px;">{{ $t('ui.usedAcross') }}</div>
        <template v-for="key in configKeys" :key="key">
          <Field v-if="clientConfig[key].password !== undefined" :label="key + ' · ' + $t('ui.password')">
            <KeyInput v-model="clientConfig[key].password" @regenerate="shuffle(key)" />
          </Field>
          <Field v-if="clientConfig[key].uuid !== undefined" :label="key + ' · ' + $t('ui.uuid')">
            <KeyInput v-model="clientConfig[key].uuid" @regenerate="shuffle(key)" />
          </Field>
          <Field v-if="key === 'vless'" :label="key + ' · ' + $t('ui.flow')">
            <input class="input mono" style="font-size: 12.5px;" v-model="clientConfig[key].flow" />
          </Field>
          <Field v-if="key === 'hysteria'" :label="key + ' · ' + $t('types.hy.auth')">
            <KeyInput v-model="clientConfig[key].auth_str" @regenerate="shuffle(key)" />
          </Field>
        </template>
      </div>

      <!-- ===================== Links ===================== -->
      <div v-else-if="tab === 'links'">
        <template v-if="localLinks.length > 0">
          <div style="display: flex; flex-direction: column; gap: 7px; margin-bottom: 16px;">
            <div
              v-for="(lnk, i) in localLinks"
              :key="'local' + i"
              class="link-row"
              :title="$t('copyToClipboard')"
              @click="copyToClipboard(lnk.uri)"
            >
              <span class="mono" style="color: var(--text-3); flex: none;">{{ i + 1 }}</span>
              <span class="mono link-uri" dir="ltr">{{ lnk.uri }}</span>
              <Ico name="copy" :size="13" style="color: var(--text-3); flex: none;" />
            </div>
          </div>
          <hr class="form-divider" />
        </template>

        <div style="display: flex; align-items: center; margin-bottom: 10px;">
          <SectionLabel>{{ $t('client.external') }}</SectionLabel>
          <div style="flex: 1;" />
          <Btn variant="subtle" sm @click="extLinks.push({ type: 'external', uri: '' })">
            <Ico name="plus" :size="14" /> {{ $t('actions.add') }}
          </Btn>
        </div>
        <div v-for="(lnk, i) in extLinks" :key="'ext' + i" style="display: flex; gap: 8px; margin-bottom: 10px;">
          <input class="input mono" dir="ltr" style="font-size: 12.5px;" placeholder="<protocol>://<data>" v-model="lnk.uri" />
          <IconBtn name="trash" danger :title="$t('actions.del')" @click="extLinks.splice(i, 1)" />
        </div>

        <hr class="form-divider" />
        <div style="display: flex; align-items: center; margin-bottom: 10px;">
          <SectionLabel>{{ $t('client.sub') }}</SectionLabel>
          <div style="flex: 1;" />
          <Btn variant="subtle" sm @click="subLinks.push({ type: 'sub', uri: '' })">
            <Ico name="plus" :size="14" /> {{ $t('actions.add') }}
          </Btn>
        </div>
        <div v-for="(lnk, i) in subLinks" :key="'sub' + i" style="display: flex; gap: 8px; margin-bottom: 10px;">
          <input class="input mono" dir="ltr" style="font-size: 12.5px;" placeholder="http[s]://<domain>[:]<port>/<path>" v-model="lnk.uri" />
          <IconBtn name="trash" danger :title="$t('actions.del')" @click="subLinks.splice(i, 1)" />
        </div>
      </div>
    </template>

    <template #footer>
      <Btn style="flex: 1;" @click="$emit('close')">{{ $t('ui.cancel') }}</Btn>
      <Btn variant="primary" style="flex: 1;" :loading="loading" @click="saveChanges">
        <Ico name="check" :size="16" /> {{ isNew ? $t('ui.createClientBtn') : $t('ui.saveChanges') }}
      </Btn>
    </template>
  </Drawer>
</template>

<script lang="ts" setup>
import { computed, ref, watch } from 'vue'
import { useI18n } from 'vue-i18n'
import Data from '@/store/modules/data'
import { Client, Link, createClient, randomConfigs, shuffleConfigs, updateConfigs } from '@/types/clients'
import { HumanReadable } from '@/plugins/utils'
import { locale } from '@/locales'
import Drawer from '@/components/ui/Drawer.vue'
import Tabs from '@/components/ui/Tabs.vue'
import Field from '@/components/ui/Field.vue'
import Btn from '@/components/ui/Btn.vue'
import Ico from '@/components/ui/Ico.vue'
import Avatar from '@/components/ui/Avatar.vue'
import Check from '@/components/ui/Check.vue'
import MSwitchRow from '@/components/ui/MSwitchRow.vue'
import SwitchLabel from '@/components/ui/SwitchLabel.vue'
import SectionLabel from '@/components/ui/SectionLabel.vue'
import KeyInput from '@/components/ui/KeyInput.vue'
import DateTimeInput from '@/components/ui/DateTimeInput.vue'
import BarMini from '@/components/charts/BarMini.vue'
import IconBtn from '@/components/ui/IconBtn.vue'
import { copyToClipboard } from '@/plugins/clipboard'

const props = defineProps<{
  visible: boolean
  id: number
  groups: string[]
  inboundTags: { title: string; value: number }[]
}>()
const emit = defineEmits<{ close: [] }>()

const { t } = useI18n({ useScope: 'global' })

const isNew = computed(() => props.id === 0)
const tab = ref('general')
const loading = ref(false)
const client = ref<Client>(createClient())
const clientConfig = ref<any>({})
const extLinks = ref<Link[]>([])
const subLinks = ref<Link[]>([])
const localLinks = computed(() => client.value.links?.filter((l) => l.type == 'local') ?? [])

const tabs = computed((): [string, string][] => [
  ['general', t('ui.general')],
  ['access', t('ui.accessLimits')],
  ['keys', t('ui.uuidPass')],
  ['links', t('client.links')],
])

// ---------- group chips (same palette logic as Clients.vue) ----------
const FIXED_GROUP_COLORS: Record<string, string> = {
  vip: 'var(--amber)',
  standard: 'var(--brand)',
  trial: 'var(--text-3)',
  staff: 'var(--emerald)',
  reseller: 'var(--cyan)',
}
const GROUP_PALETTE = ['var(--amber)', 'var(--brand)', 'var(--emerald)', 'var(--cyan)', 'var(--violet)']
const groupColor = (g: string): string => {
  if (!g) return 'var(--text-3)'
  if (FIXED_GROUP_COLORS[g]) return FIXED_GROUP_COLORS[g]
  let h = 0
  for (const ch of g) h = (h * 31 + ch.charCodeAt(0)) % 997
  return GROUP_PALETTE[h % GROUP_PALETTE.length]
}
const groupItems = computed(() => Array.from(new Set(props.groups.filter((g) => g.length > 0))))
const groupChipStyle = (g: string) => {
  const col = groupColor(g)
  const on = client.value.group === g
  return {
    cursor: 'pointer',
    height: '32px',
    padding: '0 13px',
    textTransform: 'capitalize' as const,
    color: on ? '#fff' : col,
    background: on ? col : `color-mix(in srgb, ${col} 12%, transparent)`,
    borderColor: `color-mix(in srgb, ${col} 30%, transparent)`,
  }
}

// ---------- header subtitle: protocol(s) · group ----------
const subTitle = computed(() => {
  if (isNew.value) return t('ui.createAccount')
  const types = Array.from(new Set(
    (client.value.inbounds ?? [])
      .map((id) => Data().inbounds?.find((i: any) => i.id == id)?.type)
      .filter((x: any) => !!x),
  ))
  const parts = [types.join(', '), client.value.group].filter((s) => !!s && String(s).length > 0)
  return parts.length > 0 ? parts.join(' · ') : t('objects.client')
})

// ---------- config ----------
const configKeys = computed(() => Object.keys(clientConfig.value ?? {}))
const shuffle = (k?: string) => {
  shuffleConfigs(clientConfig.value, k)
}

// ---------- inbounds ----------
const inboundItems = computed(() => props.inboundTags.map((it) => {
  const inb = Data().inbounds?.find((i: any) => i.id == it.value)
  return {
    id: it.value,
    tag: it.title,
    type: inb?.type ?? '',
    port: inb?.listen_port ?? '',
    online: Data().onlines?.inbound ? Data().onlines.inbound.includes(it.title) : false,
  }
}))
const toggleInbound = (id: number) => {
  const list = client.value.inbounds
  client.value.inbounds = list.includes(id) ? list.filter((i) => i !== id) : [...list, id].sort()
}
const setAllInbounds = () => {
  client.value.inbounds = props.inboundTags.map((i) => i.value).sort()
}

// ---------- limits ----------
const volumeGiB = computed({
  get: () => (client.value.volume == 0 ? 0 : client.value.volume / (1024 ** 3)),
  set: (v: number | string) => {
    const n = Number(v)
    client.value.volume = n > 0 ? n * (1024 ** 3) : 0
  },
})
const delayStart = computed({
  get: () => client.value.delayStart ?? false,
  set: (v: boolean) => {
    client.value.delayStart = v
    client.value.resetDays = v ? 1 : 0
    if (v && !autoReset.value) client.value.expiry = 0
  },
})
const autoReset = computed({
  get: () => client.value.autoReset ?? false,
  set: (v: boolean) => {
    client.value.autoReset = v
    client.value.resetDays = v ? 1 : 0
    if (!v) client.value.nextReset = 0
  },
})
const resetDays = computed({
  get: () => client.value.resetDays ?? 1,
  set: (v: number | string | null) => {
    let n = typeof v === 'number' ? v : Number(v)
    if (!n) n = 1
    if (client.value.nextReset && client.value.nextReset > 0) {
      client.value.nextReset += (n - (client.value.resetDays ?? 0)) * 24 * 60 * 60
    }
    client.value.resetDays = n
  },
})

const up = computed(() => HumanReadable.sizeFormat(client.value.up))
const down = computed(() => HumanReadable.sizeFormat(client.value.down))
const totalUp = computed(() => HumanReadable.sizeFormat((client.value.totalUp ?? 0) + client.value.up))
const totalDown = computed(() => HumanReadable.sizeFormat((client.value.totalDown ?? 0) + client.value.down))
const nextResetFormatted = computed(() => {
  const ts = client.value.nextReset ?? 0
  if (ts == 0) return '-'
  return new Date(ts * 1000).toLocaleString(locale)
})
const percent = computed(() =>
  client.value.volume > 0 ? Math.round(((client.value.up + client.value.down) * 100) / client.value.volume) : 0,
)
const percentColor = computed(() => {
  if (client.value.up + client.value.down >= client.value.volume) return 'var(--rose)'
  return percent.value > 90 ? 'var(--amber)' : 'var(--brand)'
})
const resetUsage = () => {
  client.value.totalUp = (client.value.totalUp ?? 0) + client.value.up
  client.value.totalDown = (client.value.totalDown ?? 0) + client.value.down
  client.value.up = 0
  client.value.down = 0
}

// ---------- lifecycle / save (payload identical to legacy modals/Client.vue) ----------
const updateData = async (id: number) => {
  tab.value = 'general'
  if (id > 0) {
    loading.value = true
    const newData = await Data().loadClients(id)
    client.value = createClient(newData)
    clientConfig.value = client.value.config
    loading.value = false
  } else {
    client.value = createClient()
    // 新建默认：开启自动重置，周期 30 天
    client.value.autoReset = true
    client.value.resetDays = 30
    clientConfig.value = randomConfigs('client')
  }
  // links are kept untouched and merged back into the payload on save
  extLinks.value = client.value.links?.filter((l) => l.type == 'external') ?? []
  subLinks.value = client.value.links?.filter((l) => l.type == 'sub') ?? []
}

const saveChanges = async () => {
  if (!props.visible || loading.value) return
  // check duplicate name
  const isDuplicateName = Data().checkClientName(props.id, client.value.name)
  if (isDuplicateName) return

  // check if delayStart is true and autoReset is false, set expiry to 0
  if (client.value.delayStart && !client.value.autoReset) client.value.expiry = 0

  loading.value = true
  client.value.config = updateConfigs(clientConfig.value, client.value.name)
  client.value.links = [
    ...extLinks.value.filter((l) => l.uri != ''),
    ...subLinks.value.filter((l) => l.uri != ''),
  ]
  const success = await Data().save('clients', props.id == 0 ? 'new' : 'edit', client.value)
  if (success) emit('close')
  loading.value = false
}

watch(() => props.visible, (v) => {
  if (v) updateData(props.id)
})
</script>

<style scoped>
.access-grid {
  display: grid;
  grid-template-columns: minmax(0, 0.92fr) minmax(0, 1.08fr);
  gap: 26px;
  align-items: start;
}
@media (max-width: 700px) {
  .access-grid { grid-template-columns: 1fr; gap: 26px; }
}
.access-head {
  display: flex;
  align-items: center;
  gap: 10px;
  min-height: 30px;
  margin-bottom: 10px;
}
.inb-list {
  display: flex;
  flex-direction: column;
  gap: 9px;
  max-height: 372px;
  overflow-y: auto;
}
.switch-stack {
  display: flex;
  flex-direction: column;
  gap: 13px;
  margin: 2px 0 16px;
}
.inb-item {
  display: flex;
  align-items: center;
  gap: 12px;
  padding: 11px 14px;
  border-radius: 11px;
  border: 1px solid var(--line);
  background: var(--surface);
  cursor: pointer;
  transition: all 0.15s var(--ease);
}
.inb-item.on { background: var(--brand-soft); }
.inb-dot { width: 8px; height: 8px; border-radius: 50%; background: var(--text-3); flex: none; }
.inb-dot.online { background: var(--emerald); }
.usage-card {
  padding: 16px;
  margin-bottom: 18px;
  background: var(--surface-3);
  border: none;
  box-shadow: none;
}
.suffix-box {
  width: 80px;
  flex: none;
  display: flex;
  align-items: center;
  justify-content: center;
  color: var(--text-2);
}
</style>
