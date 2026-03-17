package smokepod

import (
	"context"
	"testing"
)

func TestLocalTarget_Exec(t *testing.T) {
	target := NewLocalTarget("", nil)

	result, err := target.Exec(context.Background(), "echo hello")
	if err != nil {
		t.Fatalf("Exec failed: %v", err)
	}

	if result.ExitCode != 0 {
		t.Errorf("ExitCode = %d, want 0", result.ExitCode)
	}

	if result.Stdout != "hello\n" {
		t.Errorf("Stdout = %q, want %q", result.Stdout, "hello\n")
	}
}

func TestLocalTarget_ExecWithExitCode(t *testing.T) {
	target := NewLocalTarget("", nil)

	result, err := target.Exec(context.Background(), "exit 42")
	if err != nil {
		t.Fatalf("Exec failed: %v", err)
	}

	if result.ExitCode != 42 {
		t.Errorf("ExitCode = %d, want 42", result.ExitCode)
	}
}

func TestLocalTarget_ExecWithStderr(t *testing.T) {
	target := NewLocalTarget("", nil)

	result, err := target.Exec(context.Background(), "echo error >&2")
	if err != nil {
		t.Fatalf("Exec failed: %v", err)
	}

	if result.Stderr != "error\n" {
		t.Errorf("Stderr = %q, want %q", result.Stderr, "error\n")
	}
}

func TestLocalTarget_ExecWithEnv(t *testing.T) {
	target := NewLocalTarget("", []string{"MY_VAR=testvalue"})

	result, err := target.Exec(context.Background(), "echo $MY_VAR")
	if err != nil {
		t.Fatalf("Exec failed: %v", err)
	}

	if result.Stdout != "testvalue\n" {
		t.Errorf("Stdout = %q, want %q", result.Stdout, "testvalue\n")
	}
}

func TestLocalTarget_Close(t *testing.T) {
	target := NewLocalTarget("", nil)

	if err := target.Close(); err != nil {
		t.Errorf("Close returned error: %v", err)
	}
}

func TestLocalTarget_CustomShell(t *testing.T) {
	target := NewLocalTarget("/bin/bash", nil)

	result, err := target.Exec(context.Background(), "echo $0")
	if err != nil {
		t.Fatalf("Exec failed: %v", err)
	}

	if result.ExitCode != 0 {
		t.Errorf("ExitCode = %d, want 0", result.ExitCode)
	}
}
