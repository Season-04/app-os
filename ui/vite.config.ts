// vite.config.js
import { resolve } from 'node:path'
import { fileURLToPath, URL } from 'node:url'
import { defineConfig } from 'vite'
import vue from '@vitejs/plugin-vue'
import vueJsx from '@vitejs/plugin-vue-jsx'
import dts from 'vite-plugin-dts'

// https://vitejs.dev/config/
export default defineConfig({
  // build: {
  //   lib: {
  //     entry: resolve(__dirname, 'src/appos-ui.ts'),
  //     name: 'ApposUI',
  //     fileName: 'appos-ui',
  //     formats: ['es', 'umd', 'cjs'],
  //   },
  // },
  build: {
    rollupOptions: {
      input: {
        'appos-ui': './src/appos-ui.ts'
      },
      output: {
        entryFileNames: `[name].js`,
        chunkFileNames: `[name].js`,
        assetFileNames: `[name].[ext]`,
      }
    }
  },
  plugins: [
    vue(),
    vueJsx(),
    dts(),
  ],
  resolve: {
    alias: {
      '@': fileURLToPath(new URL('./src', import.meta.url))
    }
  }
})
