import { defineConfig } from 'vite'
import vue from '@vitejs/plugin-vue'

// https://vite.dev/config/
export default defineConfig({
  plugins: [vue()],
  server: {
    host: '0.0.0.0', // ホストを0.0.0.0に設定して外部からアクセス可能にする
  },
  envDir: './../',
  build: {
    outDir: './../public/',
    emptyOutDir: true,
  }
})