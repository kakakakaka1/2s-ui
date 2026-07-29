import { ref } from 'vue'

// Drag-and-drop plus move-up/down for the rule lists on Rules and Dns.
//
// Takes a getter and only ever splices. Both callers pass a readonly computed
// (`rules`, `dnsRules`) that resolves to the array hanging off the config
// object, so reassigning would fail at runtime — Vue logs "Write operation
// failed: computed value is readonly" and the click silently does nothing. A
// Ref<any[]> parameter would accept those computeds at compile time and hide
// exactly that, which is why this is a getter.
export function useReorder(list: () => any[]) {
  const draggedItemIndex = ref<number | null>(null)

  const onDragStart = (index: number) => { draggedItemIndex.value = index }

  const onDrop = (index: number) => {
    if (draggedItemIndex.value !== null) {
      const arr = list()
      const draggedItem = arr[draggedItemIndex.value]
      arr.splice(draggedItemIndex.value, 1)
      arr.splice(index, 0, draggedItem)
      draggedItemIndex.value = null
    }
  }

  const move = (index: number, dir: number) => {
    const arr = list()
    const to = index + dir
    if (to < 0 || to >= arr.length) return
    const item = arr[index]
    arr.splice(index, 1)
    arr.splice(to, 0, item)
  }

  return { onDragStart, onDrop, move }
}
