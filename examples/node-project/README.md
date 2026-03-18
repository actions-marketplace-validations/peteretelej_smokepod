# Node Project Example

Using smokepod as an npm devDependency for Docker-based smoke tests.

## Setup

```bash
npm install
```

## Run

```bash
npm run smoke
```

## Files

- `package.json` - Declares smokepod as a devDependency with a `smoke` script
- `smokepod.yaml` - Test configuration
- `tests/api.test` - API smoke test file

## How It Works

The `smokepod` npm package installs the native binary automatically. You can call it from package scripts just like any other dev tool.
