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
  </div>
</template>

<script lang="ts" setup>
import Sidebar from './Sidebar.vue'
import Topbar from './Topbar.vue'
</script>
