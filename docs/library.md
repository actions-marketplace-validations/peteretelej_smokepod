# Go Library Usage

Smokepod can be used as a Go library for programmatic test execution.

## Installation

```bash
go get github.com/peteretelej/smokepod
```

## Quick Start

```go
package main

import (
    "context"
    "fmt"
    "log"

    "github.com/peteretelej/smokepod/pkg/smokepod"
)

func main() {
    result, err := smokepod.RunFile(context.Background(), "smokepod.yaml")
    if err != nil {
        log.Fatal(err)
    }

    fmt.Printf("Tests passed: %v\n", result.Passed)
    fmt.Printf("Summary: %d passed, %d failed\n",
        result.Summary.Passed, result.Summary.Failed)
}
```

## API Reference

### Running Tests

```go
// Run tests from a config file
result, err := smokepod.RunFile(ctx, "smokepod.yaml")

// Run tests from a Config struct
result, err := smokepod.Run(ctx, config)

// Run with options
result, err := smokepod.RunWithOptions(ctx, config,
    smokepod.OptTimeout(10*time.Minute),
    smokepod.OptParallel(false),
    smokepod.OptFailFast(true),
)
```

### Config Parsing

```go
// Parse and validate a config file
config, err := smokepod.ParseConfig("smokepod.yaml")
if err != nil {
    log.Fatal(err)
}

// Validate a Config struct
err := smokepod.ValidateConfig(&config)
```

### Options

| Function | Description |
|----------|-------------|
| `OptTimeout(d time.Duration)` | Set global timeout |
| `OptParallel(bool)` | Enable/disable parallel execution |
| `OptFailFast(bool)` | Stop on first failure |
| `OptBaseDir(string)` | Base directory for relative paths |

## Types

### Config

```go
type Config struct {
    Name     string           `yaml:"name"`
    Version  string           `yaml:"version"`
    Settings Settings         `yaml:"settings"`
    Tests    []TestDefinition `yaml:"tests"`
}

type Settings struct {
    Timeout  time.Duration `yaml:"timeout"`
    Parallel *bool         `yaml:"parallel"`
    FailFast bool          `yaml:"fail_fast"`
}

type TestDefinition struct {
    Name  string   `yaml:"name"`
    Type  string   `yaml:"type"`   // "cli" or "playwright"
    Image string   `yaml:"image"`
    File  string   `yaml:"file"`   // CLI: path to .test file
    Run   []string `yaml:"run"`    // CLI: specific sections
    Path  string   `yaml:"path"`   // Playwright: project path
    Args  []string `yaml:"args"`   // Playwright: pass-through args
}
```

### Result

```go
type Result struct {
    Name      string        `json:"name"`
    Timestamp time.Time     `json:"timestamp"`
    Duration  time.Duration `json:"duration"`
    Passed    bool          `json:"passed"`
    Summary   Summary       `json:"summary"`
    Tests     []TestResult  `json:"tests"`
}

type Summary struct {
    Total   int `json:"total"`
    Passed  int `json:"passed"`
    Failed  int `json:"failed"`
    Skipped int `json:"skipped"`
}

type TestResult struct {
    Name     string          `json:"name"`
    Type     string          `json:"type"`
    Passed   bool            `json:"passed"`
    Duration time.Duration   `json:"duration"`
    Error    string          `json:"error,omitempty"`
    Sections []SectionResult `json:"sections,omitempty"`
}

type SectionResult struct {
    Name     string          `json:"name"`
    Passed   bool            `json:"passed"`
    Commands []CommandResult `json:"commands"`
}

type CommandResult struct {
    Command  string `json:"command"`
    Line     int    `json:"line"`
    Expected string `json:"expected"`
    Actual   string `json:"actual"`
    Passed   bool   `json:"passed"`
    Error    string `json:"error,omitempty"`
}
```

## Examples

### Basic Usage

```go
result, err := smokepod.RunFile(ctx, "smokepod.yaml")
if err != nil {
    log.Fatalf("Test execution failed: %v", err)
}

if !result.Passed {
    for _, test := range result.Tests {
        if !test.Passed {
            fmt.Printf("FAIL: %s - %s\n", test.Name, test.Error)
        }
    }
    os.Exit(1)
}
```

### With Timeout and Context

```go
ctx, cancel := context.WithTimeout(context.Background(), 5*time.Minute)
defer cancel()

result, err := smokepod.RunFile(ctx, "smokepod.yaml")
```

### Using RunOptions

```go
opts := smokepod.RunOptions{
    Timeout:  10 * time.Minute,
    FailFast: true,
    BaseDir:  "/path/to/project",
}

parallel := false
opts.Parallel = &parallel

result, err := smokepod.RunWithOptions(ctx, config, opts.ToOptions()...)
```

### Custom Reporter

```go
result, err := smokepod.RunFile(ctx, "smokepod.yaml")
if err != nil {
    log.Fatal(err)
}

// Write JSON to file
f, _ := os.Create("results.json")
defer f.Close()

reporter := smokepod.NewReporter(f)
reporter.SetPretty(true)
reporter.Report(result)
```

### Building Config Programmatically

```go
config := smokepod.Config{
    Name:    "my-tests",
    Version: "1",
    Settings: smokepod.Settings{
        Timeout: 5 * time.Minute,
    },
    Tests: []smokepod.TestDefinition{
        {
            Name:  "api-health",
            Type:  "cli",
            Image: "curlimages/curl:latest",
            File:  "tests/api.test",
        },
    },
}

if err := smokepod.ValidateConfig(&config); err != nil {
    log.Fatal(err)
}

result, err := smokepod.Run(ctx, config)
```

### Iterating Results

```go
result, _ := smokepod.RunFile(ctx, "smokepod.yaml")

for _, test := range result.Tests {
    fmt.Printf("%s: %v\n", test.Name, test.Passed)

    // CLI tests have section details
    for _, section := range test.Sections {
        fmt.Printf("  Section %s: %v\n", section.Name, section.Passed)

        for _, cmd := range section.Commands {
            if !cmd.Passed {
                fmt.Printf("    Line %d: %s\n", cmd.Line, cmd.Error)
            }
        }
    }
}
```

## Version Information

```go
// Get version string
version := smokepod.VersionString()
// e.g., "1.0.0 (commit: abc123, built: 2024-01-15)"

// Individual version components
fmt.Println(smokepod.Version)    // "1.0.0"
fmt.Println(smokepod.Commit)     // "abc123"
fmt.Println(smokepod.BuildDate)  // "2024-01-15"
```

## Error Handling

```go
result, err := smokepod.RunFile(ctx, "smokepod.yaml")

switch {
case err != nil:
    // Execution error (Docker issues, config errors, etc.)
    log.Fatalf("Execution error: %v", err)

case !result.Passed:
    // Tests ran but some failed
    log.Printf("Tests failed: %d/%d",
        result.Summary.Failed, result.Summary.Total)
    os.Exit(1)

default:
    // All tests passed
    log.Println("All tests passed!")
}
```
