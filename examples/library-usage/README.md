# Library Usage Example

Example of using smokepod as a Go library.

## Run

```bash
go run . ../cli-docker/smokepod.yaml
```

Or build and run separately:

```bash
go build -o library-usage .
./library-usage ../cli-docker/smokepod.yaml
```

## Code Overview

The example demonstrates:

1. Using `smokepod.RunFile()` to run tests from a config file
2. Setting up context with timeout
3. Accessing `Result` struct fields
4. Iterating over test results
5. JSON serialization of results

## Alternative API Usage

```go
// Using RunWithOptions for more control
result, err := smokepod.RunWithOptions(ctx, config,
    smokepod.OptTimeout(5*time.Minute),
    smokepod.OptParallel(true),
    smokepod.OptFailFast(false),
)

// Parsing config separately
config, err := smokepod.ParseConfig("config.yaml")
if err != nil {
    log.Fatal(err)
}
result, err := smokepod.Run(ctx, *config)
```
