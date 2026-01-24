package smokepod

import (
	"context"
	"io"
	"path/filepath"
	"time"
)

// Option is an alias for ExecutorOption for the public API.
type Option = ExecutorOption

// Run executes tests with the given configuration.
func Run(ctx context.Context, config Config) (*Result, error) {
	return RunWithOptions(ctx, config)
}

// RunFile loads config from path and executes tests.
func RunFile(ctx context.Context, path string) (*Result, error) {
	config, err := ParseConfig(path)
	if err != nil {
		return nil, err
	}

	// Set base directory to config file's directory for relative path resolution
	baseDir := filepath.Dir(path)
	return RunWithOptions(ctx, *config, WithBaseDir(baseDir))
}

// RunWithOptions executes tests with additional options.
func RunWithOptions(ctx context.Context, config Config, opts ...Option) (*Result, error) {
	executor := NewExecutor(&config, opts...)
	return executor.Execute(ctx)
}

// ValidateConfig validates a config without running tests.
func ValidateConfig(config *Config) error {
	return validateConfig(config)
}

// Convenience option constructors - these are aliases to the With* functions
// for API consistency. Users can use either form.

// OptTimeout sets the global timeout for all tests.
func OptTimeout(d time.Duration) Option {
	return WithTimeout(d)
}

// OptParallel sets whether to run tests in parallel.
func OptParallel(enabled bool) Option {
	return WithParallel(enabled)
}

// OptFailFast sets whether to stop on first failure.
func OptFailFast(enabled bool) Option {
	return WithFailFast(enabled)
}

// OptBaseDir sets the base directory for resolving relative paths.
func OptBaseDir(dir string) Option {
	return WithBaseDir(dir)
}

// RunOptions holds options for running tests (alternative to functional options).
type RunOptions struct {
	Timeout  time.Duration
	Parallel *bool
	FailFast bool
	Output   io.Writer
	Pretty   bool
	BaseDir  string
}

// ToOptions converts RunOptions to a slice of Option.
func (o RunOptions) ToOptions() []Option {
	var opts []Option
	if o.Timeout > 0 {
		opts = append(opts, OptTimeout(o.Timeout))
	}
	if o.Parallel != nil {
		opts = append(opts, OptParallel(*o.Parallel))
	}
	if o.FailFast {
		opts = append(opts, OptFailFast(true))
	}
	if o.BaseDir != "" {
		opts = append(opts, OptBaseDir(o.BaseDir))
	}
	return opts
}
