/// <reference types="vitest" />
import { defineConfig, type UserConfig } from 'vite'
import vue from '@vitejs/plugin-vue'
import istanbul from 'vite-plugin-istanbul'

// https://vite.dev/config/
export default defineConfig({
  plugins: [
    vue(),
    istanbul({
      include: 'src/*',
      exclude: ['node_modules', 'test/'],
      extension: ['.js', '.ts', '.vue'],
      requireEnv: false,
    }),
  ],
  test: {
    environment: 'happy-dom',
    globals: true,
  },
// eslint-disable-next-line @typescript-eslint/no-explicit-any
} as UserConfig & { test: any })
