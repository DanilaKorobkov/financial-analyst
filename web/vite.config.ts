import path from 'node:path';
import { defineConfig } from 'vite';
import react from '@vitejs/plugin-react';
import tailwindcss from '@tailwindcss/vite';

// Connect-сервер слушает :8080. Проксируем сюда вызовы, чтобы избежать CORS
// preflight для Content-Type: application/json.
export default defineConfig({
  plugins: [react(), tailwindcss()],
  resolve: {
    alias: {
      '@': path.resolve(__dirname, './src'),
    },
  },
  server: {
    port: 5173,
    proxy: {
      '/company.v1.': {
        target: 'http://localhost:8080',
        changeOrigin: false,
      },
    },
  },
});
