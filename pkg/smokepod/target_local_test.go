package smokepod

import (
	"context"
	"strings"
	"testing"
)

func TestLocalTarget_Exec(t *testing.T) {
	target := NewLocalTarget("", nil, nil, "")

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
	target := NewLocalTarget("", nil, nil, "")

	result, err := target.Exec(context.Background(), "exit 42")
	if err != nil {
		t.Fatalf("Exec failed: %v", err)
	}

	if result.ExitCode != 42 {
		t.Errorf("ExitCode = %d, want 42", result.ExitCode)
	}
}

func TestLocalTarget_ExecWithStderr(t *testing.T) {
	target := NewLocalTarget("", nil, nil, "")

	result, err := target.Exec(context.Background(), "echo error >&2")
	if err != nil {
		t.Fatalf("Exec failed: %v", err)
	}

	if result.Stderr != "error\n" {
		t.Errorf("Stderr = %q, want %q", result.Stderr, "error\n")
	}
}

func TestLocalTarget_ExecWithEnv(t *testing.T) {
	target := NewLocalTarget("", nil, []string{"MY_VAR=testvalue"}, "")

	result, err := target.Exec(context.Background(), "echo $MY_VAR")
	if err != nil {
		t.Fatalf("Exec failed: %v", err)
	}

	if result.Stdout != "testvalue\n" {
		t.Errorf("Stdout = %q, want %q", result.Stdout, "testvalue\n")
	}
}

func TestLocalTarget_Close(t *testing.T) {
	target := NewLocalTarget("", nil, nil, "")

	if err := target.Close(); err != nil {
		t.Errorf("Close returned error: %v", err)
	}
}

func TestLocalTarget_CustomShell(t *testing.T) {
	target := NewLocalTarget("/bin/bash", nil, nil, "")

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
	target := NewLocalTarget("/bin/bash", []string{"--norc"}, nil, "")

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
	target := NewLocalTarget("/bin/bash", []string{"--norc", "--noprofile"}, nil, "")

	result, err := target.Exec(context.Background(), "echo works")
	if err != nil {
		t.Fatalf("Exec failed: %v", err)
	}

	if result.Stdout != "works\n" {
		t.Errorf("Stdout = %q, want %q", result.Stdout, "works\n")
	}
}

func TestLocalTarget_GetVersion_NoArgs(t *testing.T) {
	target := NewLocalTarget("/bin/bash", nil, nil, "")

	version := target.GetVersion(context.Background())
	if version == "" {
		t.Error("GetVersion returned empty string for /bin/bash")
	}
}

func TestLocalTarget_GetVersion_WithArgs(t *testing.T) {
	// With fixed args, GetVersion should run: path, args..., "--version"
	target := NewLocalTarget("/bin/bash", []string{"--norc"}, nil, "")

	version := target.GetVersion(context.Background())
	// bash --norc --version should still return version info
	if version == "" {
		t.Error("GetVersion returned empty string for /bin/bash with --norc")
	}
	if !strings.Contains(version, "bash") {
		t.Errorf("GetVersion = %q, expected it to contain 'bash'", version)
	}
}

func TestIsShellTarget(t *testing.T) {
	tests := []struct {
		path string
		want bool
	}{
		{"/bin/bash", true},
		{"/usr/bin/zsh", true},
		{"bash", true},
		{"/usr/bin/jq", false},
		{"jq", false},
		{"python3", false},
		{"/bin/sh", true},
		{"fish", true},
		{"dash", true},
		{"ksh", true},
		{"cmd.exe", true},
		{"cmd", true},
		{"C:\\Windows\\System32\\cmd.exe", true},
		{"powershell.exe", true},
		{"pwsh", true},
	}
	for _, tt := range tests {
		got := IsShellTarget(tt.path)
		if got != tt.want {
			t.Errorf("IsShellTarget(%q) = %v, want %v", tt.path, got, tt.want)
		}
	}
}

func TestLocalTarget_NonShellExec(t *testing.T) {
	// Non-shell target with mode "shell" should wrap in /bin/sh
	target := NewLocalTarget("/usr/bin/jq", nil, nil, "shell")
	result, err := target.Exec(context.Background(), "echo hello")
	if err != nil {
		t.Fatalf("Exec failed: %v", err)
	}
	if result.Stdout != "hello\n" {
		t.Errorf("Stdout = %q, want %q", result.Stdout, "hello\n")
	}
}

func TestLocalTarget_WrapMode(t *testing.T) {
	// Wrap mode should use /bin/sh even for shell targets
	target := NewLocalTarget("/bin/bash", nil, nil, "wrap")
	result, err := target.Exec(context.Background(), "echo hello")
	if err != nil {
		t.Fatalf("Exec failed: %v", err)
	}
	if result.Stdout != "hello\n" {
		t.Errorf("Stdout = %q, want %q", result.Stdout, "hello\n")
	}
}

func TestLocalTarget_WrapMode_ExposesTargetEnv(t *testing.T) {
	// Wrap mode should expose SMOKEPOD_TARGET and SMOKEPOD_TARGET_ARGS
	target := NewLocalTarget("/usr/bin/node", []string{"--experimental-vm-modules"}, nil, "wrap")
	result, err := target.Exec(context.Background(), "echo $SMOKEPOD_TARGET")
	if err != nil {
		t.Fatalf("Exec failed: %v", err)
	}
	if strings.TrimSpace(result.Stdout) != "/usr/bin/node" {
		t.Errorf("SMOKEPOD_TARGET = %q, want /usr/bin/node", strings.TrimSpace(result.Stdout))
	}

	result, err = target.Exec(context.Background(), "echo $SMOKEPOD_TARGET_ARGS")
	if err != nil {
		t.Fatalf("Exec failed: %v", err)
	}
	if strings.TrimSpace(result.Stdout) != "--experimental-vm-modules" {
		t.Errorf("SMOKEPOD_TARGET_ARGS = %q, want --experimental-vm-modules", strings.TrimSpace(result.Stdout))
	}
}

func TestLocalTarget_NonShellExec_ExposesTargetEnv(t *testing.T) {
	// Non-shell target in shell mode should also expose env vars
	target := NewLocalTarget("/usr/bin/jq", []string{"--tab"}, nil, "shell")
	result, err := target.Exec(context.Background(), "echo $SMOKEPOD_TARGET")
	if err != nil {
		t.Fatalf("Exec failed: %v", err)
	}
	if strings.TrimSpace(result.Stdout) != "/usr/bin/jq" {
		t.Errorf("SMOKEPOD_TARGET = %q, want /usr/bin/jq", strings.TrimSpace(result.Stdout))
	}
}
