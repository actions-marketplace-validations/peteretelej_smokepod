package smokepod

import (
	"context"
	"path/filepath"
	"testing"
	"time"
)

func TestRun(t *testing.T) {
	// Test that Run works with a minimal config
	cfg := Config{
		Name:    "test-suite",
		Version: "1",
		Settings: Settings{
			Timeout: 30 * time.Second,
		},
		Tests: []TestDefinition{}, // Empty tests should succeed
	}

	ctx := context.Background()
	result, err := Run(ctx, cfg)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if result == nil {
		t.Fatal("result is nil")
		return
	}

	if result.Name != "test-suite" {
		t.Errorf("result.Name = %q, want %q", result.Name, "test-suite")
	}

	if !result.Passed {
		t.Error("result.Passed should be true for empty test suite")
	}
}

func TestRunFile_Valid(t *testing.T) {
	ctx := context.Background()
	configPath := filepath.Join("..", "..", "testdata", "fixtures", "minimal.yaml")

	result, err := RunFile(ctx, configPath)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if result == nil {
		t.Fatal("result is nil")
		return
	}

	// minimal.yaml has 1 test defined
	if result.Summary.Total != 1 {
		t.Errorf("expected 1 test, got %d", result.Summary.Total)
	}
}

func TestRunFile_NotFound(t *testing.T) {
	ctx := context.Background()
	_, err := RunFile(ctx, "nonexistent.yaml")
	if err == nil {
		t.Fatal("expected error, got nil")
	}
}

func TestRunWithOptions_Timeout(t *testing.T) {
	cfg := Config{
		Name:    "test-suite",
		Version: "1",
		Settings: Settings{
			Timeout: 1 * time.Hour, // Original timeout
		},
		Tests: []TestDefinition{},
	}

	ctx := context.Background()

	// Override timeout to be much shorter
	result, err := RunWithOptions(ctx, cfg, OptTimeout(10*time.Second))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if result == nil {
		t.Fatal("result is nil")
	}
}

func TestRunWithOptions_FailFast(t *testing.T) {
	cfg := Config{
		Name:    "test-suite",
		Version: "1",
		Settings: Settings{
			FailFast: false, // Original setting
		},
		Tests: []TestDefinition{},
	}

	ctx := context.Background()

	// Override fail-fast
	result, err := RunWithOptions(ctx, cfg, OptFailFast(true))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if result == nil {
		t.Fatal("result is nil")
	}
}

func TestRunWithOptions_Parallel(t *testing.T) {
	parallel := true
	cfg := Config{
		Name:    "test-suite",
		Version: "1",
		Settings: Settings{
			Parallel: &parallel, // Original setting is parallel
		},
		Tests: []TestDefinition{},
	}

	ctx := context.Background()

	// Override to sequential
	result, err := RunWithOptions(ctx, cfg, OptParallel(false))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if result == nil {
		t.Fatal("result is nil")
	}
}

func TestValidateConfig_Valid(t *testing.T) {
	cfg := &Config{
		Name:    "test-suite",
		Version: "1",
		Tests: []TestDefinition{
			{
				Name:  "cli-test",
				Type:  "cli",
				Image: "alpine:latest",
				File:  "test.test",
			},
		},
	}

	err := ValidateConfig(cfg)
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}
}

func TestValidateConfig_MissingName(t *testing.T) {
	cfg := &Config{
		Version: "1",
		Tests:   []TestDefinition{},
	}

	err := ValidateConfig(cfg)
	if err == nil {
		t.Error("expected error for missing name")
	}
}

func TestValidateConfig_InvalidVersion(t *testing.T) {
	cfg := &Config{
		Name:    "test-suite",
		Version: "2", // Only "1" is valid
		Tests:   []TestDefinition{},
	}

	err := ValidateConfig(cfg)
	if err == nil {
		t.Error("expected error for invalid version")
	}
}

func TestRunOptions_ToOptions(t *testing.T) {
	parallel := false
	opts := RunOptions{
		Timeout:  5 * time.Minute,
		Parallel: &parallel,
		FailFast: true,
		BaseDir:  "/some/path",
	}

	result := opts.ToOptions()

	// Should have 4 options
	if len(result) != 4 {
		t.Errorf("expected 4 options, got %d", len(result))
	}

	// Test that we can apply them without panic
	cfg := Config{
		Name:    "test-suite",
		Version: "1",
	}
	executor := NewExecutor(&cfg, result...)
	if executor == nil {
		t.Error("executor should not be nil")
	}
}

func TestOption_WithTimeout(t *testing.T) {
	cfg := Config{
		Name:    "test-suite",
		Version: "1",
	}

	executor := NewExecutor(&cfg, WithTimeout(10*time.Minute))
	if executor.timeout != 10*time.Minute {
		t.Errorf("timeout = %v, want %v", executor.timeout, 10*time.Minute)
	}
}

func TestOption_WithParallel(t *testing.T) {
	cfg := Config{
		Name:    "test-suite",
		Version: "1",
	}

	executor := NewExecutor(&cfg, WithParallel(false))
	if executor.parallel != false {
		t.Errorf("parallel = %v, want false", executor.parallel)
	}
}

func TestOption_WithFailFast(t *testing.T) {
	cfg := Config{
		Name:    "test-suite",
		Version: "1",
	}

	executor := NewExecutor(&cfg, WithFailFast(true))
	if executor.failFast != true {
		t.Errorf("failFast = %v, want true", executor.failFast)
	}
}

func TestOption_WithBaseDir(t *testing.T) {
	cfg := Config{
		Name:    "test-suite",
		Version: "1",
	}

	executor := NewExecutor(&cfg, WithBaseDir("/custom/path"))
	if executor.baseDir != "/custom/path" {
		t.Errorf("baseDir = %v, want /custom/path", executor.baseDir)
	}
}
