package smokepod

import (
	"bufio"
	"context"
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"strings"
	"testing"
	"time"
)

// TestHelperProcess is invoked as a subprocess by tests. It acts as a JSONL
// server: reads {"command":"..."} lines from stdin, executes each command via
// sh -c, and writes {"stdout":"...","stderr":"...","exit_code":N} back.
func TestHelperProcess(t *testing.T) {
	if os.Getenv("SMOKEPOD_TEST_HELPER") != "1" {
		return
	}

	mode := os.Getenv("SMOKEPOD_TEST_MODE")
	switch mode {
	case "crash":
		// Exit immediately without reading stdin
		os.Exit(1)
	case "bad_json":
		// Write invalid JSON for every request
		scanner := bufio.NewScanner(os.Stdin)
		for scanner.Scan() {
			fmt.Println("not valid json {{{")
		}
		os.Exit(0)
	default:
		// Normal JSONL server
		scanner := bufio.NewScanner(os.Stdin)
		for scanner.Scan() {
			var req processRequest
			if err := json.Unmarshal(scanner.Bytes(), &req); err != nil {
				os.Exit(2)
			}
			cmd := exec.Command("sh", "-c", req.Command)
			var stdout, stderr strings.Builder
			cmd.Stdout = &stdout
			cmd.Stderr = &stderr
			exitCode := 0
			if err := cmd.Run(); err != nil {
				if exitErr, ok := err.(*exec.ExitError); ok {
					exitCode = exitErr.ExitCode()
				}
			}
			resp := processResponse{
				Stdout:   stdout.String(),
				Stderr:   stderr.String(),
				ExitCode: exitCode,
			}
			data, _ := json.Marshal(resp)
			fmt.Println(string(data))
		}
		os.Exit(0)
	}
}

func helperCommand() string {
	return fmt.Sprintf("%s -test.run=TestHelperProcess", os.Args[0])
}

func helperEnv(mode string) []string {
	return []string{
		"SMOKEPOD_TEST_HELPER=1",
		"SMOKEPOD_TEST_MODE=" + mode,
	}
}

func newTestProcessTarget(t *testing.T, mode string) *ProcessTarget {
	t.Helper()
	// We need to set the env vars before creating the target, so we create
	// the process manually using the helper command.
	cmd := helperCommand()
	env := helperEnv(mode)

	// Set env vars for the subprocess
	for _, e := range env {
		parts := strings.SplitN(e, "=", 2)
		t.Setenv(parts[0], parts[1])
	}

	ctx := context.Background()
	target, err := NewProcessTarget(ctx, cmd)
	if err != nil {
		t.Fatalf("NewProcessTarget failed: %v", err)
	}
	t.Cleanup(func() { _ = target.Close() })
	return target
}

func TestProcessTarget_Exec(t *testing.T) {
	target := newTestProcessTarget(t, "")

	result, err := target.Exec(context.Background(), "echo hello")
	if err != nil {
		t.Fatalf("Exec failed: %v", err)
	}

	if result.Stdout != "hello\n" {
		t.Errorf("Stdout = %q, want %q", result.Stdout, "hello\n")
	}
	if result.ExitCode != 0 {
		t.Errorf("ExitCode = %d, want 0", result.ExitCode)
	}
}

func TestProcessTarget_ExecExitCode(t *testing.T) {
	target := newTestProcessTarget(t, "")

	result, err := target.Exec(context.Background(), "exit 42")
	if err != nil {
		t.Fatalf("Exec failed: %v", err)
	}

	if result.ExitCode != 42 {
		t.Errorf("ExitCode = %d, want 42", result.ExitCode)
	}
}

func TestProcessTarget_ExecStderr(t *testing.T) {
	target := newTestProcessTarget(t, "")

	result, err := target.Exec(context.Background(), "echo error >&2")
	if err != nil {
		t.Fatalf("Exec failed: %v", err)
	}

	if result.Stderr != "error\n" {
		t.Errorf("Stderr = %q, want %q", result.Stderr, "error\n")
	}
}

func TestProcessTarget_ExecMultipleCommands(t *testing.T) {
	target := newTestProcessTarget(t, "")

	for i := range 3 {
		result, err := target.Exec(context.Background(), fmt.Sprintf("echo cmd%d", i))
		if err != nil {
			t.Fatalf("Exec %d failed: %v", i, err)
		}
		want := fmt.Sprintf("cmd%d\n", i)
		if result.Stdout != want {
			t.Errorf("Exec %d: Stdout = %q, want %q", i, result.Stdout, want)
		}
	}
}

func TestProcessTarget_ExecProcessCrash(t *testing.T) {
	target := newTestProcessTarget(t, "crash")

	_, err := target.Exec(context.Background(), "echo hello")
	if err == nil {
		t.Fatal("expected error for crashed process, got nil")
	}
}

func TestProcessTarget_ExecMalformedJSON(t *testing.T) {
	target := newTestProcessTarget(t, "bad_json")

	_, err := target.Exec(context.Background(), "echo hello")
	if err == nil {
		t.Fatal("expected error for malformed JSON, got nil")
	}
	if !strings.Contains(err.Error(), "parsing response") {
		t.Errorf("error = %q, want it to contain %q", err.Error(), "parsing response")
	}
}

func TestProcessTarget_ExecTimeout(t *testing.T) {
	target := newTestProcessTarget(t, "")

	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Millisecond)
	defer cancel()

	// Sleep long enough that the context expires
	time.Sleep(5 * time.Millisecond)

	_, err := target.Exec(ctx, "echo hello")
	if err == nil {
		t.Fatal("expected error for timed out context, got nil")
	}
}

func TestProcessTarget_Close(t *testing.T) {
	target := newTestProcessTarget(t, "")

	// Execute a command to confirm process is alive
	_, err := target.Exec(context.Background(), "echo alive")
	if err != nil {
		t.Fatalf("Exec failed: %v", err)
	}

	// Close should shut down cleanly
	// (cleanup runs via t.Cleanup, but test explicit close too)
	if err := target.Close(); err != nil {
		t.Errorf("Close returned error: %v", err)
	}
}
