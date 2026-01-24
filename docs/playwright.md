# Playwright Integration

Smokepod runs Playwright tests in Docker containers using the official Microsoft Playwright images.

## Overview

The Playwright runner:
1. Creates a container with the Playwright image
2. Mounts your project directory at `/app`
3. Runs `npm ci` to install dependencies
4. Executes `npx playwright test --reporter=json`
5. Parses JSON output and reports results

## Configuration

```yaml
tests:
  - name: e2e-smoke
    type: playwright
    path: ./e2e                 # required: path to playwright project
    image: mcr.microsoft.com/playwright:v1.45.0-jammy  # optional
    args: ["--grep", "@smoke"]  # optional: pass-through to playwright
```

### Required Fields

| Field | Description |
|-------|-------------|
| `path` | Path to your Playwright project (contains `playwright.config.ts` and `package.json`) |

### Optional Fields

| Field | Default | Description |
|-------|---------|-------------|
| `image` | `mcr.microsoft.com/playwright:latest` | Docker image |
| `args` | `[]` | Arguments passed to `npx playwright test` |

## Project Requirements

Your Playwright project must have:

```
e2e/
├── package.json           # with playwright dependency
├── playwright.config.ts   # playwright configuration
└── tests/
    └── example.spec.ts    # test files
```

### package.json

```json
{
  "name": "e2e-tests",
  "private": true,
  "devDependencies": {
    "@playwright/test": "^1.45.0"
  }
}
```

### playwright.config.ts

```typescript
import { defineConfig } from '@playwright/test';

export default defineConfig({
  testDir: './tests',
  reporter: [['json', { outputFile: 'results.json' }]],
  use: {
    baseURL: 'http://host.docker.internal:3000',
  },
});
```

Note: Use `host.docker.internal` to reach services running on the host machine.

## Docker Images

Available Playwright images:

| Image | Description |
|-------|-------------|
| `mcr.microsoft.com/playwright:latest` | Latest stable version |
| `mcr.microsoft.com/playwright:v1.45.0-jammy` | Specific version on Ubuntu 22.04 |
| `mcr.microsoft.com/playwright:v1.45.0-focal` | Specific version on Ubuntu 20.04 |

Use specific versions for reproducible builds.

## Pass-through Arguments

Arguments in `args` are passed directly to `npx playwright test`:

```yaml
# Run only tests matching @smoke tag
args: ["--grep", "@smoke"]

# Run specific test file
args: ["tests/login.spec.ts"]

# Run in headed mode (requires X11 forwarding)
args: ["--headed"]

# Retry failed tests
args: ["--retries", "2"]

# Run specific project
args: ["--project", "chromium"]
```

## Example Configurations

### Basic Setup

```yaml
name: e2e-tests
version: "1"

tests:
  - name: e2e
    type: playwright
    path: ./e2e
```

### Smoke Tests Only

```yaml
tests:
  - name: e2e-smoke
    type: playwright
    path: ./e2e
    args: ["--grep", "@smoke"]
```

### Multiple Test Suites

```yaml
tests:
  - name: e2e-smoke
    type: playwright
    path: ./e2e
    args: ["--grep", "@smoke"]

  - name: e2e-critical
    type: playwright
    path: ./e2e
    args: ["--grep", "@critical"]

  - name: e2e-full
    type: playwright
    path: ./e2e
```

### Pinned Version

```yaml
tests:
  - name: e2e
    type: playwright
    path: ./e2e
    image: mcr.microsoft.com/playwright:v1.45.0-jammy
```

## Troubleshooting

### "Cannot find module '@playwright/test'"

Your `package.json` doesn't include Playwright:

```json
{
  "devDependencies": {
    "@playwright/test": "^1.45.0"
  }
}
```

### Tests timeout

1. Increase the global timeout:
   ```yaml
   settings:
     timeout: 15m
   ```

2. Or add timeout to playwright config:
   ```typescript
   export default defineConfig({
     timeout: 60000,
   });
   ```

### Cannot connect to localhost

Inside Docker, `localhost` refers to the container itself. Use `host.docker.internal`:

```typescript
use: {
  baseURL: 'http://host.docker.internal:3000',
}
```

### Image pull slow

Pre-pull the image before running tests:

```bash
docker pull mcr.microsoft.com/playwright:v1.45.0-jammy
```

### Tests pass locally but fail in smokepod

1. Ensure CI environment is handled in your tests
2. The container runs with `CI=true` environment variable
3. Check for browser-specific issues (container may not have GPU)

### Permission denied on mounted files

The container runs as root by default. If your tests create files, they may have root ownership on the host.

## JSON Output

Smokepod parses Playwright's JSON reporter output. Test results include:
- Pass/fail status
- Test duration
- Error messages for failures

Example result:
```json
{
  "name": "e2e-smoke",
  "type": "playwright",
  "passed": true,
  "duration": 45000000000
}
```
