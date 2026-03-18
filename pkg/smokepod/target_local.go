package smokepod

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/peteretelej/smokepod/pkg/smokepod/runners"
)

var knownShells = map[string]bool{
	"sh": true, "bash": true, "zsh": true,
	"dash": true, "ksh": true, "fish": true,
}

// IsShellTarget returns true if the given path refers to a known shell.
func IsShellTarget(path string) bool {
	return knownShells[filepath.Base(path)]
}

type LocalTarget struct {
	spec targetLaunchSpec
	env  []string
	mode string
}

func NewLocalTarget(path string, args []string, env []string, mode string) *LocalTarget {
	if path == "" {
		path = "/bin/sh"
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
		cmd = exec.CommandContext(ctx, "/bin/sh", "-c", command)
	} else {
		cmd = t.spec.cmd(ctx, "-c", command)
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
		// In wrap/non-shell mode, Exec uses /bin/sh so report the shell version
		cmd = exec.CommandContext(ctx, "/bin/sh", "--version")
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
