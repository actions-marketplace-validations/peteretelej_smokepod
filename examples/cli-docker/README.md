# CLI Docker Example

A minimal smokepod example running CLI tests inside an Alpine Linux Docker container.

## Run

```bash
npx smokepod run smokepod.yaml
```

Or with the Go-installed binary:

```bash
smokepod run smokepod.yaml
```

## Files

- `smokepod.yaml` - Configuration file
- `tests/basic.test` - Test file with example commands

## What It Tests

- Basic echo commands
- Multi-line output
- Exit code assertions
- Regex matching
- Environment variables
