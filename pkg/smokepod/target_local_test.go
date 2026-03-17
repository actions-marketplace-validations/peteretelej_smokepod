package smokepod

import (
	"context"
	"strings"
	"testing"
)

func TestLocalTarget_Exec(t *testing.T) {
	target := NewLocalTarget("", nil, nil)

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
	target := NewLocalTarget("", nil, nil)

	result, err := target.Exec(context.Background(), "exit 42")
	if err != nil {
		t.Fatalf("Exec failed: %v", err)
	}

	if result.ExitCode != 42 {
		t.Errorf("ExitCode = %d, want 42", result.ExitCode)
	}
}

func TestLocalTarget_ExecWithStderr(t *testing.T) {
	target := NewLocalTarget("", nil, nil)

	result, err := target.Exec(context.Background(), "echo error >&2")
	if err != nil {
		t.Fatalf("Exec failed: %v", err)
	}

	if result.Stderr != "error\n" {
		t.Errorf("Stderr = %q, want %q", result.Stderr, "error\n")
	}
}

func TestLocalTarget_ExecWithEnv(t *testing.T) {
	target := NewLocalTarget("", nil, []string{"MY_VAR=testvalue"})

	result, err := target.Exec(context.Background(), "echo $MY_VAR")
	if err != nil {
		t.Fatalf("Exec failed: %v", err)
	}

	if result.Stdout != "testvalue\n" {
		t.Errorf("Stdout = %q, want %q", result.Stdout, "testvalue\n")
	}
}

func TestLocalTarget_Close(t *testing.T) {
	target := NewLocalTarget("", nil, nil)

	if err := target.Close(); err != nil {
		t.Errorf("Close returned error: %v", err)
	}
}

func TestLocalTarget_CustomShell(t *testing.T) {
	target := NewLocalTarget("/bin/bash", nil, nil)

	result, err := target.Exec(context.Background(), "echo $0")
	if err != nil {
		t.Fatalf("Exec failed: %v", err)
	}

	if result.ExitCode != 0 {
		t.Errorf("ExitCode = %d, want 0", result.ExitCode)
	}
}

func TestLocalTarget_ExecWithFixedArgs(t *testing.T) {
	// Use bash with --norc to verify fixed args are passed before -c
	target := NewLocalTarget("/bin/bash", []string{"--norc"}, nil)

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

func TestLocalTarget_ExecWithMultipleFixedArgs(t *testing.T) {
	target := NewLocalTarget("/bin/bash", []string{"--norc", "--noprofile"}, nil)

	result, err := target.Exec(context.Background(), "echo works")
	if err != nil {
		t.Fatalf("Exec failed: %v", err)
	}

	if result.Stdout != "works\n" {
		t.Errorf("Stdout = %q, want %q", result.Stdout, "works\n")
	}
}

func TestLocalTarget_GetVersion_NoArgs(t *testing.T) {
	target := NewLocalTarget("/bin/bash", nil, nil)

	version := target.GetVersion(context.Background())
	if version == "" {
		t.Error("GetVersion returned empty string for /bin/bash")
	}
}

func TestLocalTarget_GetVersion_WithArgs(t *testing.T) {
	// With fixed args, GetVersion should run: path, args..., "--version"
	target := NewLocalTarget("/bin/bash", []string{"--norc"}, nil)

	version := target.GetVersion(context.Background())
	// bash --norc --version should still return version info
	if version == "" {
		t.Error("GetVersion returned empty string for /bin/bash with --norc")
	}
	if !strings.Contains(version, "bash") {
		t.Errorf("GetVersion = %q, expected it to contain 'bash'", version)
	}
}
