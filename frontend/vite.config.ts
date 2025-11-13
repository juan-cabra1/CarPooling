import { defineConfig } from 'vite'
import react from '@vitejs/plugin-react'
import path from 'path'

// https://vite.dev/config/
export default defineConfig({
  plugins: [react()],
  resolve: {
    alias: {
      '@': path.resolve(__dirname, './src'),
      '@/components': path.resolve(__dirname, './src/components'),
      '@/pages': path.resolve(__dirname, './src/pages'),
      '@/services': path.resolve(__dirname, './src/services'),
      '@/contexts': path.resolve(__dirname, './src/contexts'),
      '@/hooks': path.resolve(__dirname, './src/hooks'),
      '@/utils': path.resolve(__dirname, './src/utils'),
      '@/types': path.resolve(__dirname, './src/types'),
      '@/assets': path.resolve(__dirname, './src/assets'),
    },
  },
  server: {
    port: 3000,
    proxy: {
      '/api/users': {
        target: 'http://localhost:8001',
        changeOrigin: true,
      },
      '/api/trips': {
        target: 'http://localhost:8002',
        changeOrigin: true,
      },
      '/api/bookings': {
        target: 'http://localhost:8003',
        changeOrigin: true,
      },
      '/api/search': {
        target: 'http://localhost:8004',
        changeOrigin: true,
      },
    },
  },
})
