// One-shot: swap native <select class="input"> for the custom <Select> component.
// Inner <option>/<optgroup> markup stays untouched (Select parses it from its slot).
import { readFileSync, writeFileSync } from 'node:fs'
import { execSync } from 'node:child_process'

const files = execSync('git grep -l "<select" -- src/', { encoding: 'utf8' })
  .split('\n').filter(Boolean)

let total = 0
for (const f of files) {
  let s = readFileSync(f, 'utf8')
  const before = s

  // opening tags (attributes may span lines); drop the now-redundant input class
  s = s.replace(/<select\b/g, '<Select')
  s = s.replace(/<\/select>/g, '</Select>')
  // class="input" is built into the Select trigger; keep extra classes like mono
  s = s.replace(/(<Select[^>]*?)\sclass="input"/g, '$1')
  s = s.replace(/(<Select[^>]*?)\sclass="input ([^"]+)"/g, '$1 class="$2"')

  if (s === before) continue

  // ensure import
  if (!s.includes("components/ui/Select.vue")) {
    s = s.replace(/(<script lang="ts" setup>\s*\n)/, `$1import Select from '@/components/ui/Select.vue'\n`)
    // options-API or non-setup scripts: fall back to inserting after <script lang="ts">
    if (!s.includes("components/ui/Select.vue")) {
      s = s.replace(/(<script lang="ts">\s*\n)/, `$1import Select from '@/components/ui/Select.vue'\n`)
    }
  }
  if (!s.includes("components/ui/Select.vue")) {
    console.error(`!! ${f}: import not injected (check script block manually)`)
  }

  writeFileSync(f, s)
  total++
  console.log(`ok ${f}`)
}
console.log(`${total} files migrated`)
