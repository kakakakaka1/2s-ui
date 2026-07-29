<template>
  <div>
    <SectionLabel style="margin-bottom: 12px;">TProxy</SectionLabel>
    <Network :data="inbound" :mb="15" />
  </div>
</template>

<script lang="ts" setup>
import SectionLabel from '@/components/ui/SectionLabel.vue'
import Network from '../out/Network.vue'

defineProps<{ inbound: any }>()
</script>
