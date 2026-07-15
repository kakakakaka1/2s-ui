<template>
  <div>
    <SectionLabel style="margin-bottom: 12px;">TProxy</SectionLabel>
    <Network :data="inbound" />
  </div>
</template>

<script lang="ts" setup>
import SectionLabel from '@/components/ui/SectionLabel.vue'
import Network from './Network.vue'

defineProps<{ inbound: any }>()
</script>
