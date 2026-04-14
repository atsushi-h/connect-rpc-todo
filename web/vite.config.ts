import { tanstackStart } from '@tanstack/react-start/plugin/vite'
import viteReact from '@vitejs/plugin-react'
import { defineConfig } from 'vite'

export default defineConfig({
  plugins: [
    tanstackStart({
      srcDirectory: 'app',
      vite: {
        // @ts-expect-error -- TanStack Start の型定義に plugins が未定義だが動作する
        plugins: [viteReact()],
      },
    }),
  ],
  server: {
    port: 4000,
  },
  resolve: {
    alias: {
      '~': `${import.meta.dirname}/app`,
    },
  },
})
