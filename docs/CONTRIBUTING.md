# Contributing to Smokepod

## Prerequisites

- Go 1.21+
- Docker
- golangci-lint (optional, for linting)

## Local Setup

Clone the repository:

```bash
git clone https://github.com/peteretelej/smokepod.git
cd smokepod
```

Install dependencies:

```bash
go mod download
```

Build:

```bash
go build ./cmd/smokepod
```

## Git Hooks

Set up the pre-push hook to run CI checks locally before pushing:

```bash
./scripts/setup-hooks.sh
```

This installs a pre-push hook that runs:
- `go test -race ./...`
- `go vet ./...`
- `golangci-lint run` (if installed)
- `go build ./cmd/smokepod`

## Development Workflow

### Running Tests

```bash
go test ./...

# With race detection
go test -race ./...

# Verbose output
go test -v ./...
```

### Linting

Install golangci-lint:

```bash
# macOS
brew install golangci-lint

# Linux/other
go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest
```

Run:

```bash
golangci-lint run
```

### Building

```bash
go build ./cmd/smokepod
./smokepod --version
```

## Project Structure

```
smokepod/
├── cmd/smokepod/           # CLI entrypoint
├── pkg/smokepod/           # Public library
│   ├── config.go           # Config types and parsing
│   ├── executor.go         # Test orchestration
│   ├── docker.go           # Container management
│   ├── reporter.go         # JSON output
│   └── runners/            # CLI and Playwright runners
├── internal/testfile/      # .test file parser
├── testdata/               # Test fixtures
├── examples/               # Usage examples
└── docs/                   # Documentation
```

## Making Changes

1. Create a branch for your changes
2. Make your changes
3. Ensure tests pass: `go test ./...`
4. Ensure linting passes: `golangci-lint run`
5. Push (pre-push hook will verify)
6. Open a pull request

## Code Style

- Follow standard Go conventions
- Run `gofmt` or let your editor handle formatting
- Keep functions focused and small
- Add tests for new functionality
