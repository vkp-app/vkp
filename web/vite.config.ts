import { defineConfig } from 'vite'
import react from '@vitejs/plugin-react'

// https://vitejs.dev/config/
export default defineConfig({
  plugins: [react()],
  server: {
    https: false,
    host: true,
    port: 3000,
    hmr: {
      protocol: "wss",
      clientPort: 443
    }
  }
})