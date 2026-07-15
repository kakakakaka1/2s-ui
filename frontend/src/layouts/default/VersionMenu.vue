<template>
  <Pop :width="248" align="start">
    <template #trigger="{ toggle }">
      <button type="button" class="ver-pill" :title="pillTip" @click="toggle">
        <span class="mono">v{{ ver.appVersion || '—' }}</span>
        <span v-if="ver.newVersion" class="ver-dot" />
      </button>
    </template>

    <template #default="{ close }">
      <div class="ver-card">
        <div class="ver-head">
          <span class="ver-label">{{ $t('ui.version.current') }}</span>
          <Btn variant="subtle" icon sm :loading="ver.checking" @click="ver.recheck()">
            <Ico name="refresh" :size="14" />
          </Btn>
        </div>

        <div class="ver-big mono">v{{ ver.appVersion || '—' }}</div>
        <div class="ver-latest">
          {{ $t('ui.version.latest') }}: <span class="mono">{{ ver.latestTag ? 'v' + ver.latestTag.replace(/^v/i, '') : '—' }}</span>
        </div>

        <template v-if="ver.newVersion">
          <div class="ver-alert">
            <span class="ver-alert-ico"><Ico name="download" :size="15" /></span>
            <span>
              <span class="ver-alert-title">{{ $t('ui.newVersion') }}!</span>
              <span class="mono" style="display: block; font-size: 12px; margin-top: 2px;">v{{ ver.newVersion }}</span>
            </span>
          </div>
          <Btn v-if="ver.canSelfUpdate" variant="primary" class="ver-cta" @click="ver.updateOpen = true; close()">
            <Ico name="download" :size="15" /> {{ $t('ui.selfUpdate.now') }}
          </Btn>
          <a v-else class="btn btn-primary ver-cta" :href="RELEASE_URL" target="_blank" rel="noopener">
            <Ico name="download" :size="15" /> {{ $t('ui.selfUpdate.now') }}
          </a>
        </template>
        <div v-else-if="ver.latestTag" class="ver-uptodate">
          <Ico name="check" :size="14" /> {{ $t('ui.selfUpdate.upToDate') }}
        </div>

        <a class="ver-changelog" :href="RELEASE_URL" target="_blank" rel="noopener">
          {{ $t('ui.version.changelog') }} <Ico name="link" :size="12" />
        </a>
      </div>
    </template>
  </Pop>
</template>

<script lang="ts" setup>
import { computed } from 'vue'
import { useI18n } from 'vue-i18n'
import Pop from '@/components/ui/Pop.vue'
import Btn from '@/components/ui/Btn.vue'
import Ico from '@/components/ui/Ico.vue'
import VersionStore from '@/store/modules/version'
import { RELEASE_URL } from '@/plugins/updateCheck'

const ver = VersionStore()
const { t } = useI18n({ useScope: 'global' })

const pillTip = computed(() =>
  ver.newVersion
    ? `${t('ui.newVersion')}: v${ver.newVersion}`
    : `${t('ui.version.current')}: v${ver.appVersion || '—'}`
)
</script>

<style scoped>
.ver-pill {
  display: inline-flex; align-items: center; gap: 5px;
  margin-top: 3px; padding: 1px 7px;
  border: 1px solid var(--line); border-radius: 999px;
  background: var(--surface-2);
  font-size: 10.5px; font-weight: 700; color: var(--text-3);
  cursor: pointer;
  transition: all .15s var(--ease);
}
.ver-pill:hover { color: var(--text); border-color: var(--line-2); background: var(--surface-3); }
.ver-dot {
  width: 6px; height: 6px; border-radius: 50%; flex: none;
  background: var(--amber);
  animation: ver-breathe 1.8s ease-in-out infinite;
}
/* 呼吸灯:点体明暗呼吸 + 光晕扩散渐隐;光晕用 amber 0% 而非 transparent,避免向黑色插值 */
@keyframes ver-breathe {
  0%, 100% { opacity: 1; box-shadow: 0 0 0 0 color-mix(in srgb, var(--amber) 55%, transparent); }
  50% { opacity: .45; box-shadow: 0 0 0 5px color-mix(in srgb, var(--amber) 0%, transparent); }
}
@media (prefers-reduced-motion: reduce) {
  .ver-dot { animation: none; box-shadow: 0 0 0 3px color-mix(in srgb, var(--amber) 22%, transparent); }
}

.ver-card { padding: 14px; display: flex; flex-direction: column; }
.ver-head { display: flex; align-items: center; justify-content: space-between; }
.ver-label { font-size: 12px; font-weight: 600; color: var(--text-3); }
.ver-big {
  text-align: center; font-size: 19px; font-weight: 800;
  letter-spacing: -.02em; margin-top: 6px;
}
.ver-latest {
  text-align: center; font-size: 11.5px; color: var(--text-3); margin-top: 4px;
}

.ver-alert {
  display: flex; align-items: center; gap: 10px;
  margin-top: 12px; padding: 10px 12px;
  border-radius: var(--radius-sm);
  color: var(--amber);
  background: color-mix(in srgb, var(--amber) 14%, transparent);
  border: 1px solid color-mix(in srgb, var(--amber) 28%, transparent);
}
.ver-alert-ico {
  width: 30px; height: 30px; border-radius: 50%; flex: none;
  display: flex; align-items: center; justify-content: center;
  background: color-mix(in srgb, var(--amber) 18%, transparent);
}
.ver-alert-title { display: block; font-size: 12.5px; font-weight: 700; }

.ver-cta {
  width: 100%; margin-top: 12px; gap: 7px;
  text-decoration: none;
}
.ver-uptodate {
  display: flex; align-items: center; justify-content: center; gap: 6px;
  margin-top: 12px; font-size: 12px; font-weight: 600; color: var(--emerald);
}

.ver-changelog {
  display: inline-flex; align-items: center; justify-content: center; gap: 4px;
  margin-top: 12px; align-self: center;
  font-size: 11.5px; font-weight: 600; color: var(--text-3);
  text-decoration: none;
  transition: color .15s var(--ease);
}
.ver-changelog:hover { color: var(--text); }
</style>
