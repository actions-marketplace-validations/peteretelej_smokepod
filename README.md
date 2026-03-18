# Smokepod

[![CI](https://github.com/peteretelej/smokepod/actions/workflows/ci.yml/badge.svg)](https://github.com/peteretelej/smokepod/actions/workflows/ci.yml)
[![Go Report Card](https://goreportcard.com/badge/github.com/peteretelej/smokepod)](https://goreportcard.com/report/github.com/peteretelej/smokepod)
[![Go Reference](https://pkg.go.dev/badge/github.com/peteretelej/smokepod.svg)](https://pkg.go.dev/github.com/peteretelej/smokepod)

Smoke test runner for CLI applications. Record expected outputs, verify against fixtures, and optionally run smoke tests in Docker containers.

## Quick Start

Write a test file, record expected outputs, then verify against any target:

```text
# tests/basics.test
## echo
$ echo "hello world"
```

```bash
npx smokepod record --target /bin/bash --tests tests/ --fixtures fixtures/
npx smokepod verify --target ./my-shell --tests tests/ --fixtures fixtures/
```

## Three Modes

### Record

Execute commands from `.test` files using a reference target, save results to fixture JSON:

```bash
npx smokepod record --target /bin/bash --tests tests/ --fixtures fixtures/
```

Pass fixed arguments to the target:

```bash
npx smokepod record --target /bin/bash --target-arg --norc --target-arg --noprofile \
  --tests tests/ --fixtures fixtures/
```

### Verify

Re-execute commands and compare output against recorded fixtures:

```bash
npx smokepod verify --target ./my-shell --tests tests/ --fixtures fixtures/
```

Use process mode for targets that communicate via JSONL (no shell wrapping):

```bash
npx smokepod verify --target ./my-adapter --mode process \
  --tests tests/ --fixtures fixtures/
```

### Run

Execute tests in Docker containers using a YAML config:

```bash
npx smokepod run smokepod.yaml
```

```yaml
# smokepod.yaml
name: myproject-smoke
version: "1"

tests:
  - name: api-smoke
    type: cli
    image: curlimages/curl:latest
    file: tests/api.test
    run: [health]  # optional: run specific sections
```

## Installation

```bash
# npx (no install needed)
npx smokepod --help

# npm devDependency
npm install --save-dev smokepod

# Go
go install github.com/peteretelej/smokepod/cmd/smokepod@latest
```

The npm package downloads the matching native binary during `postinstall`, so no Go toolchain is required.

## GitHub Action

```yaml
# Verify
- uses: peteretelej/smokepod@v1
  with:
    mode: verify
    target: /bin/bash
    tests: tests/
    fixtures: fixtures/

# Run Docker smoke tests
- uses: peteretelej/smokepod@v1
  with:
    mode: run
    config: smokepod.yaml
```

See [docs/github-action.md](docs/github-action.md) for inputs, multi-platform matrix, and more examples.

## Test File Format

| Syntax | Meaning |
|--------|---------|
| `## name` | Named test section |
| `$ command` | Command to execute |
| Following lines | Expected output |
| `(re)` suffix | Regex matching |
| `(stderr)` suffix | Match against stderr |
| `[exit:N]` | Expected exit code |
| `#` | Comment |

See [docs/test-format.md](docs/test-format.md) for full syntax, multi-line commands, and examples.

## Playwright Tests

```yaml
tests:
  - name: e2e
    type: playwright
    path: ./e2e
    image: mcr.microsoft.com/playwright:v1.45.0-jammy
```

See [docs/playwright.md](docs/playwright.md) for setup and configuration.

## Requirements

- Docker (only needed for `run` mode)
- Go 1.21+ (only for building from source)

## Documentation

- [Configuration Reference](docs/config-reference.md) - All config options
- [Test File Format](docs/test-format.md) - `.test` file syntax
- [GitHub Action](docs/github-action.md) - CI integration
- [Playwright Integration](docs/playwright.md) - Browser testing
- [Go Library Usage](docs/library.md) - Using smokepod as a library
- [Troubleshooting](docs/troubleshooting.md) - Common issues and solutions

## License

MIT
