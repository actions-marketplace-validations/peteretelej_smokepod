package smokepod

import (
	"bytes"
	"encoding/json"
	"testing"
	"time"
)

func TestReporter_JSON(t *testing.T) {
	result := &Result{
		Name:      "test-suite",
		Timestamp: time.Date(2024, 1, 1, 12, 0, 0, 0, time.UTC),
		Duration:  5 * time.Second,
		Passed:    true,
		Summary: Summary{
			Total:  2,
			Passed: 2,
		},
		Tests: []TestResult{
			{Name: "test1", Type: "cli", Passed: true},
			{Name: "test2", Type: "playwright", Passed: true},
		},
	}

	var buf bytes.Buffer
	reporter := NewReporter(&buf)

	if err := reporter.Report(result); err != nil {
		t.Fatalf("Report() error = %v", err)
	}

	// Should be compact JSON (no newlines inside)
	output := buf.String()
	if output == "" {
		t.Error("Report() produced empty output")
	}

	// Verify it's valid JSON by parsing it back
	var parsed Result
	if err := json.Unmarshal(buf.Bytes(), &parsed); err != nil {
		t.Errorf("Report() produced invalid JSON: %v", err)
	}

	if parsed.Name != "test-suite" {
		t.Errorf("parsed.Name = %q, want %q", parsed.Name, "test-suite")
	}
}

func TestReporter_Pretty(t *testing.T) {
	result := &Result{
		Name:   "pretty-test",
		Passed: true,
	}

	var buf bytes.Buffer
	reporter := NewReporter(&buf)
	reporter.SetPretty(true)

	if err := reporter.Report(result); err != nil {
		t.Fatalf("Report() error = %v", err)
	}

	output := buf.String()

	// Pretty output should have indentation
	if !bytes.Contains(buf.Bytes(), []byte("  ")) {
		t.Error("Pretty output should contain indentation")
	}

	// Should end with newline
	if output[len(output)-1] != '\n' {
		t.Error("Pretty output should end with newline")
	}
}

func TestReporter_ValidJSON(t *testing.T) {
	// Test with various result states
	testCases := []struct {
		name   string
		result *Result
	}{
		{
			name: "passed",
			result: &Result{
				Name:   "passed-suite",
				Passed: true,
				Summary: Summary{
					Total:  1,
					Passed: 1,
				},
			},
		},
		{
			name: "failed",
			result: &Result{
				Name:   "failed-suite",
				Passed: false,
				Summary: Summary{
					Total:  2,
					Passed: 1,
					Failed: 1,
				},
				Tests: []TestResult{
					{Name: "ok", Passed: true},
					{Name: "fail", Passed: false, Error: "something went wrong"},
				},
			},
		},
		{
			name: "empty",
			result: &Result{
				Name:   "empty-suite",
				Passed: true,
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			var buf bytes.Buffer
			reporter := NewReporter(&buf)

			if err := reporter.Report(tc.result); err != nil {
				t.Fatalf("Report() error = %v", err)
			}

			// Parse it back and verify round-trip
			var parsed Result
			if err := json.Unmarshal(buf.Bytes(), &parsed); err != nil {
				t.Errorf("Report() produced invalid JSON: %v\nOutput: %s", err, buf.String())
			}

			if parsed.Name != tc.result.Name {
				t.Errorf("round-trip Name = %q, want %q", parsed.Name, tc.result.Name)
			}

			if parsed.Passed != tc.result.Passed {
				t.Errorf("round-trip Passed = %v, want %v", parsed.Passed, tc.result.Passed)
			}
		})
	}
}
