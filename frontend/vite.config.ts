import { defineConfig } from 'vite'
import react from '@vitejs/plugin-react'

export default defineConfig({
  plugins: [react()],
  base: '/',
  build: {
    outDir: 'dist',
    assetsDir: 'assets',
    sourcemap: true,
    rollupOptions: {
      output: {
        manualChunks: undefined,
      },
    },
  },
  server: {
    port: 3001,
    host: true,
    proxy: {
      '/api': {
        target: 'http://127.0.0.1:8084',
        changeOrigin: true,
        secure: false,
        rewrite: (path) => path,
      },
      '/health': {
        target: 'http://127.0.0.1:8084',
        changeOrigin: true,
      },
    },
  },
})
