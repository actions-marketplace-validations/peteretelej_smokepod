import { defineConfig } from '@playwright/test';

export default defineConfig({
  testDir: './tests',
  timeout: 30000,
  retries: 0,
  reporter: [['json', { outputFile: 'results.json' }]],
  use: {
    // Use host.docker.internal to reach services on the host
    baseURL: 'https://example.com',
    trace: 'off',
  },
  projects: [
    {
      name: 'chromium',
      use: { browserName: 'chromium' },
    },
  ],
});
