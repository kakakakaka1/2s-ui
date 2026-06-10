<template>
  <header class="topbar">
    <Btn class="hamburger" icon @click="app.sidebarOpen = true">
      <Ico name="menu" :size="19" />
    </Btn>

    <div style="min-width: 0; flex: 0 1 auto;">
      <h1>{{ $t(`ui.title.${pageId}`) }}</h1>
      <div class="page-sub">{{ $t(`ui.sub.${pageId}`) }}</div>
    </div>

    <div class="grow" />

    <LangMenu />

    <Btn icon :title="$t('ui.theme')" @click="app.toggleTheme()">
      <Ico :name="app.theme === 'dark' ? 'sun' : 'moon'" :size="17" />
    </Btn>

    <div class="topbar-divider" />

    <Pop :width="168">
      <template #trigger="{ toggle }">
        <button type="button" class="btn btn-ghost user-chip" @click="toggle">
          <span class="user-avatar">{{ initials }}</span>
          <span class="user-meta">
            <span class="user-name">{{ username }}</span>
            <span class="user-role">{{ $t('ui.superAdmin') }}</span>
          </span>
          <Ico class="user-caret" name="chevronDown" :size="14" />
        </button>
      </template>
      <template #default="{ close }">
        <button type="button" class="pop-item danger" @click="close(); doLogout()">
          <Ico name="logout" :size="15" /> {{ $t('menu.logout') }}
        </button>
      </template>
    </Pop>
  </header>
</template>

<script lang="ts" setup>
import { computed } from 'vue'
import { useRoute } from 'vue-router'
import Btn from '@/components/ui/Btn.vue'
import Ico from '@/components/ui/Ico.vue'
import Pop from '@/components/ui/Pop.vue'
import LangMenu from './LangMenu.vue'
import AppStore from '@/store/modules/app'
import { logout } from '@/plugins/httputil'

const app = AppStore()
const route = useRoute()

const pageId = computed(() => String(route.name).replace('pages.', '') || 'home')

const username = computed(() => localStorage.getItem('2sui-user') || 'admin')
const initials = computed(() => username.value.slice(0, 2).toUpperCase())

const doLogout = () => logout()
</script>

<style scoped>
.hamburger { display: none; flex: none; }
@media (max-width: 820px) {
  .hamburger { display: inline-flex; }
}
.user-chip { gap: 9px; height: 42px; padding: 0 10px 0 6px; }
[dir='rtl'] .user-chip { padding: 0 6px 0 10px; }
.user-avatar {
  width: 28px; height: 28px; border-radius: 9px; flex: none;
  background: linear-gradient(135deg, var(--brand), var(--violet));
  display: flex; align-items: center; justify-content: center;
  color: #fff; font-weight: 700; font-size: 11.5px;
}
.user-meta { text-align: start; line-height: 1.15; }
.user-name { display: block; font-size: 12.5px; font-weight: 700; color: var(--text); }
.user-role { display: block; font-size: 10.5px; font-weight: 600; color: var(--text-3); }
@media (max-width: 820px) {
  .user-meta, .user-caret { display: none; }
  .user-chip { height: 38px; padding: 0 5px; }
}
</style>
