package smokepod

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"

	"github.com/peteretelej/smokepod/pkg/smokepod/runners"
)

var knownShells = map[string]bool{
	"sh": true, "bash": true, "zsh": true,
	"dash": true, "ksh": true, "fish": true,
	"cmd": true, "cmd.exe": true,
	"powershell": true, "powershell.exe": true,
	"pwsh": true, "pwsh.exe": true,
}

// defaultShell returns the platform's default shell for running commands.
func defaultShell() string {
	if runtime.GOOS == "windows" {
		return "cmd.exe"
	}
	return "/bin/sh"
}

// shellExecFlag returns the flag used to pass a command string to a shell.
// cmd.exe uses "/c", everything else uses "-c".
func shellExecFlag(shell string) string {
	base := strings.TrimSuffix(filepath.Base(shell), ".exe")
	if base == "cmd" {
		return "/c"
	}
	return "-c"
}

// IsShellTarget returns true if the given path refers to a known shell.
// It handles both Unix and Windows path separators regardless of the
// current platform.
func IsShellTarget(path string) bool {
	base := filepath.Base(path)
	// Handle Windows backslash paths on non-Windows platforms
	if i := strings.LastIndex(base, "\\"); i >= 0 {
		base = base[i+1:]
	}
	return knownShells[base]
}

type LocalTarget struct {
	spec targetLaunchSpec
	env  []string
	mode string
}

func NewLocalTarget(path string, args []string, env []string, mode string) *LocalTarget {
	if path == "" {
		path = defaultShell()
	}
	if mode == "" {
		mode = "shell"
	}
	if len(args) > 0 && (mode == "wrap" || (mode == "shell" && !IsShellTarget(path))) {
		fmt.Fprintf(os.Stderr, "Warning: target-args %v are not passed directly in %s mode; available as $SMOKEPOD_TARGET_ARGS\n", args, mode)
	}
	return &LocalTarget{
		spec: targetLaunchSpec{path: path, args: args},
		env:  env,
		mode: mode,
	}
}

func (t *LocalTarget) Exec(ctx context.Context, command string) (runners.ExecResult, error) {
	var cmd *exec.Cmd
	if t.mode == "wrap" || (t.mode == "shell" && !IsShellTarget(t.spec.path)) {
		shell := defaultShell()
		cmd = exec.CommandContext(ctx, shell, shellExecFlag(shell), command)
	} else {
		cmd = t.spec.cmd(ctx, shellExecFlag(t.spec.path), command)
	}

	env := append(os.Environ(), t.env...)

	// In wrap/non-shell mode, expose the target binary and its args as
	// environment variables so commands can reference them.
	if t.mode == "wrap" || (t.mode == "shell" && !IsShellTarget(t.spec.path)) {
		env = append(env, "SMOKEPOD_TARGET="+t.spec.path)
		if len(t.spec.args) > 0 {
			env = append(env, "SMOKEPOD_TARGET_ARGS="+strings.Join(t.spec.args, " "))
		}
	}
	cmd.Env = env

	var stdout, stderr strings.Builder
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	err := cmd.Run()

	result := runners.ExecResult{
		Stdout:   stdout.String(),
		Stderr:   stderr.String(),
		ExitCode: 0,
	}

	if err != nil {
		if exitErr, ok := err.(*exec.ExitError); ok {
			result.ExitCode = exitErr.ExitCode()
		} else {
			return runners.ExecResult{}, fmt.Errorf("executing command: %w", err)
		}
	}

	return result, nil
}

func (t *LocalTarget) Close() error {
	return nil
}

func (t *LocalTarget) GetVersion(ctx context.Context) string {
	var cmd *exec.Cmd
	if t.mode == "wrap" || (t.mode == "shell" && !IsShellTarget(t.spec.path)) {
		// In wrap/non-shell mode, Exec uses the default shell so report its version
		cmd = exec.CommandContext(ctx, defaultShell(), "--version")
	} else if len(t.spec.args) > 0 {
		// Fixed args exist: run path, args..., "--version"
		cmd = t.spec.cmd(ctx, "--version")
	} else {
		// No fixed args: run path, "--version"
		cmd = exec.CommandContext(ctx, t.spec.path, "--version")
	}

	var stdout strings.Builder
	cmd.Stdout = &stdout

	if err := cmd.Run(); err != nil {
		return ""
	}

	output := strings.TrimSpace(stdout.String())
	if output == "" {
		return ""
	}

	firstLine := strings.Split(output, "\n")[0]
	if len(firstLine) > 100 {
		firstLine = firstLine[:100]
	}
	return firstLine
}
