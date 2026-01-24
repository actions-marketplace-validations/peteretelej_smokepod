package smokepod

import (
	"context"
	"testing"
	"time"
)

func TestNewExecutor(t *testing.T) {
	cfg := &Config{
		Name:    "test",
		Version: "1",
		Settings: Settings{
			Timeout:  10 * time.Minute,
			FailFast: true,
		},
	}

	e := NewExecutor(cfg)

	if !e.parallel {
		t.Error("parallel should default to true")
	}

	if !e.failFast {
		t.Error("failFast should be true from config")
	}

	if e.timeout != 10*time.Minute {
		t.Errorf("timeout = %v, want 10m", e.timeout)
	}
}

func TestNewExecutor_WithOptions(t *testing.T) {
	cfg := &Config{
		Name:    "test",
		Version: "1",
	}

	e := NewExecutor(cfg,
		WithParallel(false),
		WithFailFast(true),
		WithTimeout(30*time.Second),
		WithBaseDir("/custom/path"),
	)

	if e.parallel {
		t.Error("parallel should be false from option")
	}

	if !e.failFast {
		t.Error("failFast should be true from option")
	}

	if e.timeout != 30*time.Second {
		t.Errorf("timeout = %v, want 30s", e.timeout)
	}

	if e.baseDir != "/custom/path" {
		t.Errorf("baseDir = %q, want /custom/path", e.baseDir)
	}
}

func TestExecutor_Aggregate(t *testing.T) {
	cfg := &Config{
		Name:    "aggregate-test",
		Version: "1",
	}

	e := NewExecutor(cfg)

	testCases := []struct {
		name          string
		results       []TestResult
		expectPassed  bool
		expectTotal   int
		expectPassed2 int
		expectFailed  int
		expectSkipped int
	}{
		{
			name: "all passed",
			results: []TestResult{
				{Name: "t1", Passed: true},
				{Name: "t2", Passed: true},
			},
			expectPassed:  true,
			expectTotal:   2,
			expectPassed2: 2,
			expectFailed:  0,
			expectSkipped: 0,
		},
		{
			name: "one failed",
			results: []TestResult{
				{Name: "t1", Passed: true},
				{Name: "t2", Passed: false, Error: "something went wrong"},
			},
			expectPassed:  false,
			expectTotal:   2,
			expectPassed2: 1,
			expectFailed:  1,
			expectSkipped: 0,
		},
		{
			name: "with skipped",
			results: []TestResult{
				{Name: "t1", Passed: false, Error: "failed"},
				{Name: "t2", Passed: false, Error: "skipped (fail-fast)"},
				{Name: "t3", Passed: false, Error: "cancelled"},
			},
			expectPassed:  false,
			expectTotal:   3,
			expectPassed2: 0,
			expectFailed:  1,
			expectSkipped: 2,
		},
		{
			name:          "empty",
			results:       []TestResult{},
			expectPassed:  true,
			expectTotal:   0,
			expectPassed2: 0,
			expectFailed:  0,
			expectSkipped: 0,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			result := e.aggregate(tc.results, time.Second)

			if result.Passed != tc.expectPassed {
				t.Errorf("Passed = %v, want %v", result.Passed, tc.expectPassed)
			}

			if result.Summary.Total != tc.expectTotal {
				t.Errorf("Summary.Total = %d, want %d", result.Summary.Total, tc.expectTotal)
			}

			if result.Summary.Passed != tc.expectPassed2 {
				t.Errorf("Summary.Passed = %d, want %d", result.Summary.Passed, tc.expectPassed2)
			}

			if result.Summary.Failed != tc.expectFailed {
				t.Errorf("Summary.Failed = %d, want %d", result.Summary.Failed, tc.expectFailed)
			}

			if result.Summary.Skipped != tc.expectSkipped {
				t.Errorf("Summary.Skipped = %d, want %d", result.Summary.Skipped, tc.expectSkipped)
			}

			if result.Name != "aggregate-test" {
				t.Errorf("Name = %q, want aggregate-test", result.Name)
			}
		})
	}
}

func TestExecutor_ResolvePath(t *testing.T) {
	cfg := &Config{Name: "test", Version: "1"}

	e := NewExecutor(cfg, WithBaseDir("/base/dir"))

	testCases := []struct {
		input    string
		expected string
	}{
		{
			input:    "relative/path.test",
			expected: "/base/dir/relative/path.test",
		},
		{
			input:    "/absolute/path.test",
			expected: "/absolute/path.test",
		},
	}

	for _, tc := range testCases {
		result := e.resolvePath(tc.input)
		if result != tc.expected {
			t.Errorf("resolvePath(%q) = %q, want %q", tc.input, result, tc.expected)
		}
	}
}

func TestExecutor_ExecuteSequential_CancelledContext(t *testing.T) {
	cfg := &Config{
		Name:    "cancel-test",
		Version: "1",
		Tests: []TestDefinition{
			{Name: "t1", Type: "cli", Image: "alpine", File: "test.test"},
			{Name: "t2", Type: "cli", Image: "alpine", File: "test.test"},
		},
	}

	e := NewExecutor(cfg, WithParallel(false))

	// Create an already-cancelled context
	ctx, cancel := context.WithCancel(context.Background())
	cancel()

	results := e.executeSequential(ctx)

	if len(results) != 2 {
		t.Fatalf("len(results) = %d, want 2", len(results))
	}

	for i, r := range results {
		if r.Passed {
			t.Errorf("results[%d].Passed = true, want false", i)
		}
		if r.Error != "cancelled" {
			t.Errorf("results[%d].Error = %q, want cancelled", i, r.Error)
		}
	}
}

func TestExecutor_GlobalTimeout(t *testing.T) {
	cfg := &Config{
		Name:    "timeout-test",
		Version: "1",
	}

	// Very short timeout
	e := NewExecutor(cfg, WithTimeout(1*time.Nanosecond))

	ctx := context.Background()

	// Execute should apply timeout
	result, err := e.Execute(ctx)
	if err != nil {
		t.Fatalf("Execute() error = %v", err)
	}

	// With no tests, should still complete
	if !result.Passed {
		t.Errorf("Passed = false, want true (no tests)")
	}
}
