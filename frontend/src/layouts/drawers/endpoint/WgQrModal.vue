<template>
  <Modal :open="visible" title="Wireguard QrCode" :width="400" @close="$emit('close')">
    <div id="wg-qrcode-modal" class="qr-list">
      <template v-for="(l, i) in wgLinks" :key="i">
        <div v-if="l.length > 0" class="qr-item">
          <div style="display: flex; align-items: center; gap: 8px;">
            <Chip>{{ $t('types.wg.peer') + ' ' + (i + 1) }}</Chip>
            <IconBtn name="download" :title="$t('ui.download')" @click="download(l, i)" />
          </div>
          <QrcodeVue
            :value="l"
            :size="size"
            :margin="1"
            style="border-radius: .5rem; cursor: copy;"
            @click="copy(l)"
          />
        </div>
      </template>
    </div>
  </Modal>
</template>

<script lang="ts" setup>
import { computed, ref, watch } from 'vue'
import QrcodeVue from 'qrcode.vue'
import { copyToClipboard } from '@/plugins/clipboard'
import Modal from '@/components/ui/Modal.vue'
import Chip from '@/components/ui/Chip.vue'
import IconBtn from '@/components/ui/IconBtn.vue'

const props = defineProps<{
  visible: boolean
  data: any
}>()
defineEmits<{ close: [] }>()

const wgData = ref<any>({})
const wgLinks = ref<string[]>([])

const load = () => {
  wgData.value = props.data
  wgLinks.value = []
  const address = location.hostname
  wgData.value.peers.forEach((_: any, index: number) => {
    wgLinks.value.push(getWireguardLink(index, address))
  })
}

const getWireguardLink = (peerId: number, address: string): string => {
  const peerData = wgData.value.peers[peerId]
  const keys = wgData.value.ext?.keys?.find((key: any) => key.public_key == peerData.public_key)
  if (!keys || !wgData.value.ext?.public_key) return ''
  let txt = `[Interface]\n`
  txt += `PrivateKey = ${keys.private_key}\n`
  txt += `Address = ${peerData.allowed_ips.join(',')}\n`
  txt += `DNS = ${wgData.value.ext?.dns?.length > 0 ? wgData.value.ext.dns : '1.1.1.1, 9.9.9.9'}\n`
  if (wgData.value.mtu) {
    txt += `MTU = ${wgData.value.mtu}\n`
  }
  txt += `\n# ${wgData.value.tag} - ${peerId}\n`
  txt += `[Peer]\n`
  txt += `PublicKey = ${wgData.value.ext.public_key}\n`
  txt += `AllowedIPs = 0.0.0.0/0, ::/0\n`
  txt += `Endpoint = ${address}:${wgData.value.listen_port}\n`
  if (peerData.pre_shared_key) {
    txt += `\nPresharedKey = ${peerData.pre_shared_key}`
  }
  if (peerData.persistent_keepalive_interval) {
    txt += `\nPersistentKeepalive = ${peerData.persistent_keepalive_interval}\n`
  }
  return txt
}

const copy = (txt: string) => copyToClipboard(txt, 'wg-qrcode-modal')

const download = (text: string, i: number) => {
  let filename = wgData.value.tag + '_peer_' + (i + 1) + '.conf'
  let element = document.createElement('a')
  element.setAttribute('href', 'data:application/json;charset=utf-8,' + encodeURIComponent(text))
  element.setAttribute('download', filename)

  element.style.display = 'none'
  document.body.appendChild(element)

  element.click()
  document.body.removeChild(element)
}

const size = computed(() => {
  if (window.innerWidth > 380) return 300
  if (window.innerWidth > 330) return 280
  return 250
})

watch(() => props.visible, (v) => {
  if (v) load()
})
</script>

<style scoped>
.qr-list {
  padding: 18px 20px;
  display: flex;
  flex-direction: column;
  gap: 18px;
}
.qr-item {
  display: flex;
  flex-direction: column;
  align-items: center;
  gap: 9px;
}
</style>
