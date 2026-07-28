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

export default defineConfig({
  base: '',
  plugins: [
    vue(),
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
