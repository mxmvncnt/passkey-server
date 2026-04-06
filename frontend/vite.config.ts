import { defineConfig } from 'vite'
import vue from '@vitejs/plugin-vue'

// https://vite.dev/config/
export default defineConfig({
  base: "/passkey-demo/",
  plugins: [vue()],
  server: {
    proxy: {
      '/passkey': {
        target: 'http://localhost:9090',
        changeOrigin: true,
      },
      '/users': {
        target: 'http://localhost:9090',
        changeOrigin: true,
      },
    },
  },
})
