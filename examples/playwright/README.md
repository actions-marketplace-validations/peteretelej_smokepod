# Playwright Example

A minimal smokepod example using Playwright tests.

## Setup

```bash
npm install
```

## Run with Smokepod

```bash
npx smokepod run smokepod.yaml
```

## Run Locally

```bash
npx playwright test
```

## Files

- `smokepod.yaml` - Smokepod configuration
- `playwright.config.ts` - Playwright configuration
- `package.json` - Node.js dependencies
- `tests/example.spec.ts` - Example test file

## Notes

- Tests run against https://example.com (public test site)
- For testing your own services, update `baseURL` in playwright.config.ts
- Use `host.docker.internal` to reach services on your host machine
