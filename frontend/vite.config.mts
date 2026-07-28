// Plugins
import vue from '@vitejs/plugin-vue'

// Utilities
import { defineConfig } from 'vite'
import { fileURLToPath, URL } from 'node:url'
import { randomBytes } from 'crypto'

// web.go serves assets/ with max-age=31536000, so every build must produce filenames
// the previous build never used. One id per build (not per file) is enough for that,
// and keeping [name] alongside it means split chunks stay tellable apart.
const BUILD_ID = randomBytes(8).toString('hex')

// @fontsource emits every @font-face as `woff2, woff`. Anything that can run this
// bundle supports woff2, so the fallback is never requested — but it still lands in
// dist/, and from there into the binary via web.go's //go:embed.
function stripWoffFallback() {
  return {
    name: 'strip-woff-fallback',
    enforce: 'pre' as const,
    transform(code: string, id: string) {
      if (!id.includes('fontsource') || !id.includes('.css')) return null
      // Only ever drops a *fallback* (the comma is required), so a woff-only
      // @font-face — if @fontsource ever ships one — is left intact.
      const out = code.replace(/,\s*url\([^)]*\.woff\)\s*format\('woff'\)/g, '')
      if (out === code) {
        this.warn(`no woff fallback found in ${id}; @fontsource may have changed its CSS shape`)
        return null
      }
      return out
    },
  }
}

export default defineConfig({
  base: '',
  plugins: [
    vue(),
    stripWoffFallback(),
  ],
  build: {
    manifest: false,
    outDir: 'dist',
    chunkSizeWarningLimit: 2000,
    rollupOptions: {
      output: {
        entryFileNames: `assets/[name]-${BUILD_ID}.js`,
        chunkFileNames: `assets/[name]-${BUILD_ID}.js`,
        assetFileNames: (assetInfo) => {
          if (assetInfo.names.some(name => name.endsWith('.css')))
            return `assets/[name]-${BUILD_ID}.css`
          return 'assets/' + assetInfo.names[0]
        },
      },
    }
  },
  define: { 'process.env': {} },
  resolve: {
    alias: {
      '@': fileURLToPath(new URL('./src', import.meta.url)),
    },
    extensions: ['.js', '.json', '.jsx', '.mjs', '.ts', '.tsx', '.vue'],
  },
  server: {
    port: 3000,
    proxy: {
      '/app/api': {
        target: 'http://localhost:2095',
        changeOrigin: true,
      },
    },
  }
})
