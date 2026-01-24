// Package runners provides test runners for different test types.
package runners

import (
	"context"
	"fmt"
	"regexp"
	"strings"

	"github.com/peteretelej/smokepod/internal/testfile"
	"github.com/peteretelej/smokepod/pkg/smokepod"
)

// CLIRunner executes CLI tests in a container.
type CLIRunner struct {
	container *smokepod.Container
}

// NewCLIRunner creates a new CLI test runner.
func NewCLIRunner(container *smokepod.Container) *CLIRunner {
	return &CLIRunner{container: container}
}

// Run executes all commands in a section and returns results.
func (r *CLIRunner) Run(ctx context.Context, section *testfile.Section) (*smokepod.SectionResult, error) {
	result := &smokepod.SectionResult{
		Name:   section.Name,
		Passed: true,
	}

	for _, cmd := range section.Commands {
		cmdResult := r.runCommand(ctx, cmd)
		result.Commands = append(result.Commands, cmdResult)
		if !cmdResult.Passed {
			result.Passed = false
		}
	}

	return result, nil
}

func (r *CLIRunner) runCommand(ctx context.Context, cmd testfile.Command) smokepod.CommandResult {
	result := smokepod.CommandResult{
		Command: cmd.Cmd,
		Line:    cmd.Line,
		Passed:  true,
	}

	// Build expected output string for reporting
	var expectedLines []string
	for _, exp := range cmd.Expected {
		expectedLines = append(expectedLines, exp.Text)
	}
	result.Expected = strings.Join(expectedLines, "\n")

	// Execute command in container
	execResult, err := r.container.Exec(ctx, []string{"sh", "-c", cmd.Cmd})
	if err != nil {
		result.Passed = false
		result.Error = fmt.Sprintf("execution error: %v", err)
		return result
	}

	result.Actual = strings.TrimRight(execResult.Stdout, "\n")

	// Check exit code
	if execResult.ExitCode != cmd.ExitCode {
		result.Passed = false
		result.Error = fmt.Sprintf("exit code: got %d, want %d", execResult.ExitCode, cmd.ExitCode)
		return result
	}

	// Compare output if we have expected lines
	if len(cmd.Expected) > 0 {
		if err := r.compareOutput(cmd.Expected, result.Actual); err != nil {
			result.Passed = false
			result.Error = err.Error()
			return result
		}
	}

	return result
}

func (r *CLIRunner) compareOutput(expected []testfile.Expect, actual string) error {
	actualLines := strings.Split(actual, "\n")

	// Handle empty actual output
	if actual == "" {
		actualLines = []string{}
	}

	if len(actualLines) != len(expected) {
		return fmt.Errorf("line count: got %d, want %d\n%s",
			len(actualLines), len(expected), formatDiff(expected, actualLines))
	}

	for i, exp := range expected {
		actualLine := actualLines[i]
		if exp.IsRegex {
			matched, err := regexp.MatchString(exp.Text, actualLine)
			if err != nil {
				return fmt.Errorf("line %d: invalid regex %q: %v", exp.Line, exp.Text, err)
			}
			if !matched {
				return fmt.Errorf("line %d: regex mismatch\n  pattern: %s\n  actual:  %s",
					exp.Line, exp.Text, actualLine)
			}
		} else {
			if actualLine != exp.Text {
				return fmt.Errorf("line %d: mismatch\n  want: %s\n  got:  %s",
					exp.Line, exp.Text, actualLine)
			}
		}
	}

	return nil
}

func formatDiff(expected []testfile.Expect, actual []string) string {
	var b strings.Builder
	b.WriteString("--- expected\n+++ actual\n")

	maxLen := len(expected)
	if len(actual) > maxLen {
		maxLen = len(actual)
	}

	for i := 0; i < maxLen; i++ {
		if i < len(expected) {
			suffix := ""
			if expected[i].IsRegex {
				suffix = " (re)"
			}
			b.WriteString(fmt.Sprintf("- %s%s\n", expected[i].Text, suffix))
		}
		if i < len(actual) {
			b.WriteString(fmt.Sprintf("+ %s\n", actual[i]))
		}
	}

	return b.String()
}
