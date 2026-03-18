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
	cmd.Env = append(os.Environ(), t.env...)

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
	if len(t.spec.args) > 0 {
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
