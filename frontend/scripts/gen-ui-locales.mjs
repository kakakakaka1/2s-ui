// One-shot generator: builds src/locales/ui/<lang>.ts from the design handoff dictionaries.
// Usage: node scripts/gen-ui-locales.mjs <path-to-design-handoff-dir>
import { readFileSync, writeFileSync, mkdirSync } from 'node:fs'
import { join, dirname } from 'node:path'
import { fileURLToPath } from 'node:url'

const here = dirname(fileURLToPath(import.meta.url))
const handoff = process.argv[2] || join(here, '..', '..', 'design_handoff_2s_ui_console')
const outDir = join(here, '..', 'src', 'locales', 'ui')

// --- extract a top-level `const NAME = {...}` object literal from a JSX file ---
function extractObject(file, name) {
  const src = readFileSync(join(handoff, file), 'utf8')
  const marker = `const ${name} = {`
  const start = src.indexOf(marker)
  if (start < 0) throw new Error(`${name} not found in ${file}`)
  let i = src.indexOf('{', start)
  let depth = 0
  for (let j = i; j < src.length; j++) {
    const ch = src[j]
    if (ch === '{') depth++
    else if (ch === '}') {
      depth--
      if (depth === 0) {
        const literal = src.slice(i, j + 1)
        return (0, eval)(`(${literal})`)
      }
    }
  }
  throw new Error(`unbalanced braces for ${name} in ${file}`)
}

const T = extractObject('i18n.jsx', 'T')        // 6 languages, authoritative
const TX = extractObject('pages3.jsx', 'TX')    // en + zhcn
const MXD = extractObject('modals.jsx', 'MXD')  // en + zhcn
const RXD = extractObject('modals2.jsx', 'RXD') // en + zhcn

const LANGS = ['en', 'zhcn', 'zhtw', 'ru', 'vi', 'fa']

// merge order: RXD < MXD < TX < T (T wins)
function flatFor(lang) {
  const merged = {}
  for (const dict of [RXD, MXD, TX]) {
    const en = dict.en || {}
    const loc = dict[lang] || {}
    for (const k of Object.keys(en)) {
      merged[k] = loc[k] ?? en[k] // fallback to English when the dict lacks this language
    }
  }
  const t = T[lang] || {}
  for (const k of Object.keys(T.en)) {
    merged[k] = t[k] ?? T.en[k]
  }
  return merged
}

// nest dotted keys: "nav.home" -> { nav: { home } }
function nest(flat) {
  const out = {}
  for (const [k, v] of Object.entries(flat)) {
    const parts = k.split('.')
    let cur = out
    for (let i = 0; i < parts.length - 1; i++) {
      cur = cur[parts[i]] ??= {}
    }
    cur[parts[parts.length - 1]] = v
  }
  return out
}

// vue-i18n message-format: escape literal @ and | (special chars), keep {n} placeholders
function esc(s) {
  return s.replace(/@/g, "{'@'}").replace(/\|/g, "{'|'}")
}
function emit(obj, indent = '  ') {
  const lines = ['{']
  for (const [k, v] of Object.entries(obj)) {
    const key = /^[a-zA-Z_$][\w$]*$/.test(k) ? k : JSON.stringify(k)
    if (typeof v === 'object') {
      lines.push(`${indent}${key}: ${emit(v, indent + '  ')},`)
    } else {
      lines.push(`${indent}${key}: ${JSON.stringify(esc(String(v)))},`)
    }
  }
  lines.push(indent.slice(2) + '}')
  return lines.join('\n')
}

mkdirSync(outDir, { recursive: true })
for (const lang of LANGS) {
  const body = emit(nest(flatFor(lang)))
  writeFileSync(join(outDir, `${lang}.ts`), `// Generated from design_handoff_2s_ui_console dictionaries (gen-ui-locales.mjs)\nexport default ${body}\n`)
  console.log(`ui/${lang}.ts written (${Object.keys(flatFor(lang)).length} keys)`)
}
