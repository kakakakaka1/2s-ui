<template>
  <div v-if="app.sidebarOpen" class="sidebar-scrim" @click="app.sidebarOpen = false" />
  <aside class="sidebar" :class="{ 'mobile-open': app.sidebarOpen }">
    <div class="sidebar-head">
      <Logo :size="34" />
      <div style="flex: 1;">
        <div style="font-weight: 800; font-size: 16px; letter-spacing: -.02em; line-height: 1;">
          2S<span style="color: var(--text-3); font-weight: 700;">-UI</span>
        </div>
        <VersionMenu />
      </div>
      <Btn v-if="app.sidebarOpen" variant="subtle" icon sm @click="app.sidebarOpen = false">
        <Ico name="close" :size="17" />
      </Btn>
    </div>

    <nav class="sidebar-nav hide-scroll">
      <template v-for="(n, i) in nav" :key="i">
        <div v-if="'sep' in n" class="nav-sep">{{ $t(n.sep) }}</div>
        <button
          v-else
          type="button"
          class="nav-item"
          :class="{ active: isActive(n.path) }"
          @click="go(n.path)"
        >
          <span v-if="isActive(n.path)" class="nav-rail" />
          <Ico :name="n.icon" :size="18" />
          <span class="nav-label">{{ $t(n.title) }}</span>
          <span v-if="n.count !== undefined && n.count() > 0" class="nav-count">{{ n.count() }}</span>
        </button>
      </template>
    </nav>

    <div class="sidebar-foot">
      <a
        href="https://github.com/shenaba/2s-ui"
        target="_blank"
        rel="noreferrer"
        class="gh-link"
      >
        <Ico name="github" :size="15" />
        <span style="flex: 1;">GitHub</span>
      </a>
    </div>
  </aside>
</template>

<script lang="ts" setup>
import { onMounted } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import Logo from '@/components/ui/Logo.vue'
import Ico from '@/components/ui/Ico.vue'
import Btn from '@/components/ui/Btn.vue'
import VersionMenu from './VersionMenu.vue'
import AppStore from '@/store/modules/app'
import Data from '@/store/modules/data'
import VersionStore from '@/store/modules/version'

const app = AppStore()
const data = Data()
const route = useRoute()
const router = useRouter()

onMounted(() => { VersionStore().init() })

type NavItem = { title: string; icon: string; path: string; count?: () => number }
type NavSep = { sep: string }

const nav: (NavItem | NavSep)[] = [
  { title: 'ui.nav.home', icon: 'home', path: '/' },
  { title: 'ui.nav.inbounds', icon: 'inbound', path: '/inbounds', count: () => data.inbounds.length },
  { title: 'ui.nav.clients', icon: 'clients', path: '/clients', count: () => data.clients.length },
  { title: 'ui.nav.outbounds', icon: 'outbound', path: '/outbounds' },
  { title: 'ui.nav.endpoints', icon: 'endpoint', path: '/endpoints' },
  { title: 'ui.nav.services', icon: 'services', path: '/services' },
  { sep: 'ui.sec.configuration' },
  { title: 'ui.nav.tls', icon: 'tls', path: '/tls' },
  { title: 'ui.nav.basics', icon: 'basics', path: '/basics' },
  { title: 'ui.nav.rules', icon: 'rules', path: '/rules' },
  { title: 'ui.nav.dns', icon: 'dns', path: '/dns' },
  { sep: 'ui.sec.administration' },
  { title: 'ui.nav.admins', icon: 'admins', path: '/admins' },
  { title: 'ui.nav.settings', icon: 'settings', path: '/settings' },
]

const isActive = (p: string) => route.path === p
const go = (p: string) => {
  router.push(p)
  app.sidebarOpen = false
}
</script>

<style scoped>
.gh-link {
  display: flex; align-items: center; gap: 9px;
  padding: 7px 9px; border-radius: 9px;
  text-decoration: none;
  font-family: var(--font-ui); font-size: 12.5px; font-weight: 600;
  color: var(--text-3);
  transition: all .15s var(--ease);
}
.gh-link:hover { background: var(--surface-3); color: var(--text); }
</style>
