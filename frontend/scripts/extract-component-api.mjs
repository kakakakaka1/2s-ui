// Dump props / emits / slots for every design-system component so the
// docs can be written against real signatures rather than from memory.
import { readFileSync, readdirSync } from 'node:fs'
import { join } from 'node:path'

const dirs = ['src/components/ui', 'src/components/charts']
const out = []

for (const dir of dirs) {
  for (const f of readdirSync(dir).filter((n) => n.endsWith('.vue'))) {
    const src = readFileSync(join(dir, f), 'utf8')
    const name = f.replace('.vue', '')

    let props = ''
    let defaults = ''
    const wd = src.match(/withDefaults\(\s*defineProps<\{([\s\S]*?)\}>\(\),\s*\{([\s\S]*?)\}\s*\)/)
    if (wd) {
      props = wd[1].trim()
      defaults = wd[2].replace(/\s+/g, ' ').trim()
    } else {
      const m = src.match(/defineProps<\{([\s\S]*?)\}>\(\)/)
      if (m) props = m[1].trim()
    }

    let emits = ''
    const e = src.match(/defineEmits<\{([\s\S]*?)\}>\(\)/)
    if (e) emits = e[1].replace(/\s+/g, ' ').trim()
    else {
      const e2 = src.match(/defineEmits\(\[([\s\S]*?)\]\)/)
      if (e2) emits = e2[1].replace(/\s+/g, ' ').trim()
    }

    const slots = [...new Set([...src.matchAll(/<slot(?:\s+name="([^"]+)")?/g)].map((s) => s[1] || 'default'))]

    out.push({ dir: dir.split('/').pop(), name, props, defaults, emits, slots })
  }
}

for (const c of out) {
  console.log(`\n### ${c.name}  [${c.dir}]`)
  if (c.props) console.log('PROPS:\n' + c.props.split('\n').map((l) => '  ' + l.trim()).join('\n'))
  if (c.defaults) console.log('DEFAULTS: ' + c.defaults)
  if (c.emits) console.log('EMITS: ' + c.emits)
  if (c.slots.length) console.log('SLOTS: ' + c.slots.join(', '))
}
