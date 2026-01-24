package runners

import (
	"context"
	"os/exec"
	"strings"
	"testing"
	"time"

	"github.com/peteretelej/smokepod/internal/testfile"
	"github.com/peteretelej/smokepod/pkg/smokepod"
)

func dockerAvailable() bool {
	cmd := exec.Command("docker", "info")
	return cmd.Run() == nil
}

func setupContainer(t *testing.T) (*smokepod.Container, context.Context, context.CancelFunc) {
	t.Helper()
	if !dockerAvailable() {
		t.Skip("docker not available")
	}

	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Minute)

	c, err := smokepod.NewContainer(ctx, smokepod.ContainerConfig{
		Image: "alpine:latest",
	})
	if err != nil {
		cancel()
		t.Fatalf("NewContainer failed: %v", err)
	}

	return c, ctx, cancel
}

func TestCLIRunner_ExactMatch(t *testing.T) {
	c, ctx, cancel := setupContainer(t)
	defer cancel()
	defer func() { _ = c.Terminate(ctx) }()

	runner := NewCLIRunner(c)
	section := &testfile.Section{
		Name: "test",
		Commands: []testfile.Command{
			{
				Line: 1,
				Cmd:  "echo hello",
				Expected: []testfile.Expect{
					{Line: 2, Text: "hello", IsRegex: false},
				},
				ExitCode: 0,
			},
		},
	}

	result, err := runner.Run(ctx, section)
	if err != nil {
		t.Fatalf("Run failed: %v", err)
	}

	if !result.Passed {
		t.Errorf("result.Passed = false, want true")
	}
	if len(result.Commands) != 1 {
		t.Fatalf("commands = %d, want 1", len(result.Commands))
	}
	if !result.Commands[0].Passed {
		t.Errorf("command[0].Passed = false, want true: %s", result.Commands[0].Error)
	}
}

func TestCLIRunner_ExactMismatch(t *testing.T) {
	c, ctx, cancel := setupContainer(t)
	defer cancel()
	defer func() { _ = c.Terminate(ctx) }()

	runner := NewCLIRunner(c)
	section := &testfile.Section{
		Name: "test",
		Commands: []testfile.Command{
			{
				Line: 1,
				Cmd:  "echo hello",
				Expected: []testfile.Expect{
					{Line: 2, Text: "goodbye", IsRegex: false},
				},
				ExitCode: 0,
			},
		},
	}

	result, err := runner.Run(ctx, section)
	if err != nil {
		t.Fatalf("Run failed: %v", err)
	}

	if result.Passed {
		t.Error("result.Passed = true, want false")
	}
	if result.Commands[0].Passed {
		t.Error("command[0].Passed = true, want false")
	}
	if !strings.Contains(result.Commands[0].Error, "mismatch") {
		t.Errorf("error = %q, want to contain %q", result.Commands[0].Error, "mismatch")
	}
}

func TestCLIRunner_RegexMatch(t *testing.T) {
	c, ctx, cancel := setupContainer(t)
	defer cancel()
	defer func() { _ = c.Terminate(ctx) }()

	runner := NewCLIRunner(c)
	section := &testfile.Section{
		Name: "test",
		Commands: []testfile.Command{
			{
				Line: 1,
				Cmd:  "echo 'version 1.2.3'",
				Expected: []testfile.Expect{
					{Line: 2, Text: `version \d+\.\d+\.\d+`, IsRegex: true},
				},
				ExitCode: 0,
			},
		},
	}

	result, err := runner.Run(ctx, section)
	if err != nil {
		t.Fatalf("Run failed: %v", err)
	}

	if !result.Passed {
		t.Errorf("result.Passed = false, want true: %s", result.Commands[0].Error)
	}
}

func TestCLIRunner_RegexMismatch(t *testing.T) {
	c, ctx, cancel := setupContainer(t)
	defer cancel()
	defer func() { _ = c.Terminate(ctx) }()

	runner := NewCLIRunner(c)
	section := &testfile.Section{
		Name: "test",
		Commands: []testfile.Command{
			{
				Line: 1,
				Cmd:  "echo hello",
				Expected: []testfile.Expect{
					{Line: 2, Text: `^\d+$`, IsRegex: true},
				},
				ExitCode: 0,
			},
		},
	}

	result, err := runner.Run(ctx, section)
	if err != nil {
		t.Fatalf("Run failed: %v", err)
	}

	if result.Passed {
		t.Error("result.Passed = true, want false")
	}
	if !strings.Contains(result.Commands[0].Error, "regex mismatch") {
		t.Errorf("error = %q, want to contain %q", result.Commands[0].Error, "regex mismatch")
	}
}

func TestCLIRunner_ExitCodePass(t *testing.T) {
	c, ctx, cancel := setupContainer(t)
	defer cancel()
	defer func() { _ = c.Terminate(ctx) }()

	runner := NewCLIRunner(c)
	section := &testfile.Section{
		Name: "test",
		Commands: []testfile.Command{
			{
				Line:     1,
				Cmd:      "sh -c 'exit 42'",
				ExitCode: 42,
			},
		},
	}

	result, err := runner.Run(ctx, section)
	if err != nil {
		t.Fatalf("Run failed: %v", err)
	}

	if !result.Passed {
		t.Errorf("result.Passed = false, want true: %s", result.Commands[0].Error)
	}
}

func TestCLIRunner_ExitCodeFail(t *testing.T) {
	c, ctx, cancel := setupContainer(t)
	defer cancel()
	defer func() { _ = c.Terminate(ctx) }()

	runner := NewCLIRunner(c)
	section := &testfile.Section{
		Name: "test",
		Commands: []testfile.Command{
			{
				Line:     1,
				Cmd:      "sh -c 'exit 1'",
				ExitCode: 0, // expecting 0 but command exits 1
			},
		},
	}

	result, err := runner.Run(ctx, section)
	if err != nil {
		t.Fatalf("Run failed: %v", err)
	}

	if result.Passed {
		t.Error("result.Passed = true, want false")
	}
	if !strings.Contains(result.Commands[0].Error, "exit code") {
		t.Errorf("error = %q, want to contain %q", result.Commands[0].Error, "exit code")
	}
}

func TestCLIRunner_MultipleCommands(t *testing.T) {
	c, ctx, cancel := setupContainer(t)
	defer cancel()
	defer func() { _ = c.Terminate(ctx) }()

	runner := NewCLIRunner(c)
	section := &testfile.Section{
		Name: "test",
		Commands: []testfile.Command{
			{
				Line: 1,
				Cmd:  "echo one",
				Expected: []testfile.Expect{
					{Line: 2, Text: "one", IsRegex: false},
				},
				ExitCode: 0,
			},
			{
				Line: 4,
				Cmd:  "echo two",
				Expected: []testfile.Expect{
					{Line: 5, Text: "two", IsRegex: false},
				},
				ExitCode: 0,
			},
			{
				Line: 7,
				Cmd:  "echo three",
				Expected: []testfile.Expect{
					{Line: 8, Text: "three", IsRegex: false},
				},
				ExitCode: 0,
			},
		},
	}

	result, err := runner.Run(ctx, section)
	if err != nil {
		t.Fatalf("Run failed: %v", err)
	}

	if !result.Passed {
		t.Error("result.Passed = false, want true")
		for i, cmd := range result.Commands {
			if !cmd.Passed {
				t.Errorf("command[%d] failed: %s", i, cmd.Error)
			}
		}
	}
	if len(result.Commands) != 3 {
		t.Errorf("commands = %d, want 3", len(result.Commands))
	}
}

func TestCLIRunner_MultilineOutput(t *testing.T) {
	c, ctx, cancel := setupContainer(t)
	defer cancel()
	defer func() { _ = c.Terminate(ctx) }()

	runner := NewCLIRunner(c)
	section := &testfile.Section{
		Name: "test",
		Commands: []testfile.Command{
			{
				Line: 1,
				Cmd:  "printf 'line1\\nline2\\nline3'",
				Expected: []testfile.Expect{
					{Line: 2, Text: "line1", IsRegex: false},
					{Line: 3, Text: "line2", IsRegex: false},
					{Line: 4, Text: "line3", IsRegex: false},
				},
				ExitCode: 0,
			},
		},
	}

	result, err := runner.Run(ctx, section)
	if err != nil {
		t.Fatalf("Run failed: %v", err)
	}

	if !result.Passed {
		t.Errorf("result.Passed = false, want true: %s", result.Commands[0].Error)
	}
}

func TestCLIRunner_LineMismatchCount(t *testing.T) {
	c, ctx, cancel := setupContainer(t)
	defer cancel()
	defer func() { _ = c.Terminate(ctx) }()

	runner := NewCLIRunner(c)
	section := &testfile.Section{
		Name: "test",
		Commands: []testfile.Command{
			{
				Line: 1,
				Cmd:  "printf 'line1\\nline2'",
				Expected: []testfile.Expect{
					{Line: 2, Text: "line1", IsRegex: false},
				},
				ExitCode: 0,
			},
		},
	}

	result, err := runner.Run(ctx, section)
	if err != nil {
		t.Fatalf("Run failed: %v", err)
	}

	if result.Passed {
		t.Error("result.Passed = true, want false")
	}
	if !strings.Contains(result.Commands[0].Error, "line count") {
		t.Errorf("error = %q, want to contain %q", result.Commands[0].Error, "line count")
	}
}

func TestCLIRunner_NoExpectedOutput(t *testing.T) {
	c, ctx, cancel := setupContainer(t)
	defer cancel()
	defer func() { _ = c.Terminate(ctx) }()

	runner := NewCLIRunner(c)
	section := &testfile.Section{
		Name: "test",
		Commands: []testfile.Command{
			{
				Line:     1,
				Cmd:      "true",
				ExitCode: 0,
			},
		},
	}

	result, err := runner.Run(ctx, section)
	if err != nil {
		t.Fatalf("Run failed: %v", err)
	}

	if !result.Passed {
		t.Errorf("result.Passed = false, want true: %s", result.Commands[0].Error)
	}
}

func TestCLIRunner_SectionName(t *testing.T) {
	c, ctx, cancel := setupContainer(t)
	defer cancel()
	defer func() { _ = c.Terminate(ctx) }()

	runner := NewCLIRunner(c)
	section := &testfile.Section{
		Name: "my-section",
		Commands: []testfile.Command{
			{
				Line:     1,
				Cmd:      "true",
				ExitCode: 0,
			},
		},
	}

	result, err := runner.Run(ctx, section)
	if err != nil {
		t.Fatalf("Run failed: %v", err)
	}

	if result.Name != "my-section" {
		t.Errorf("result.Name = %q, want %q", result.Name, "my-section")
	}
}
