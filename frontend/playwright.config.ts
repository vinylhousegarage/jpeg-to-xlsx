import { defineConfig } from '@playwright/test';

const isCI = !!process.env.CI;

export default defineConfig({
  testMatch: 'tests/**/*.spec.ts',

  use: {
    baseURL: process.env.BASE_URL || (isCI ? 'http://frontend:5173' : 'http://127.0.0.1:5173'),
  },

  webServer: {
    command: isCI ? 'echo "Server already running"' : 'npm run preview',
    url: isCI ? 'http://frontend:5173' : 'http://127.0.0.1:5173',
    reuseExistingServer: true,
    timeout: 120 * 1000,
  },
});
