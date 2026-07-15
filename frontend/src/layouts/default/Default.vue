<template>
  <div class="app-shell">
    <Sidebar />
    <div class="main-col">
      <Topbar />
      <main class="page">
        <div class="page-inner">
          <router-view />
        </div>
      </main>
    </div>
    <UpdateModal
      :visible="ver.updateOpen"
      :current="ver.appVersion"
      :target="ver.newVersion ?? ''"
      @close="ver.updateOpen = false"
    />
  </div>
</template>

<script lang="ts" setup>
import Sidebar from './Sidebar.vue'
import Topbar from './Topbar.vue'
import UpdateModal from '../drawers/UpdateModal.vue'
import VersionStore from '@/store/modules/version'

const ver = VersionStore()
</script>
