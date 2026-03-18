# Record/Verify Example

The simplest smokepod workflow: record expected outputs from a reference shell, then verify that another target produces the same results.

## Record

Run commands from `.test` files against a reference target and save results to fixtures:

```bash
npx smokepod record --target /bin/bash --tests tests/ --fixtures fixtures/
```

This creates JSON fixture files in `fixtures/` with the recorded output.

## Verify

Re-run the same commands against a different target and compare:

```bash
npx smokepod verify --target ./my-shell --tests tests/ --fixtures fixtures/
```

Verify fails if any output differs from the recorded fixtures.

## Files

- `tests/basics.test` - Example test file with shell commands
- `fixtures/` - Recorded fixture output (created by `record`)

## How It Works

1. Write `.test` files with commands to execute
2. Record fixtures using a known-good target (e.g. `/bin/bash`)
3. Verify any other target produces identical output
4. Commit fixtures to version control so CI can verify against them
