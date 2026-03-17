package smokepod

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"strings"

	"github.com/peteretelej/smokepod/pkg/smokepod/runners"
)

type LocalTarget struct {
	shell string
	env   []string
}

func NewLocalTarget(shell string, env []string) *LocalTarget {
	if shell == "" {
		shell = "/bin/sh"
	}
	return &LocalTarget{
		shell: shell,
		env:   env,
	}
}

func (t *LocalTarget) Exec(ctx context.Context, command string) (runners.ExecResult, error) {
	cmd := exec.CommandContext(ctx, t.shell, "-c", command)
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
	cmd := exec.CommandContext(ctx, t.shell, "--version")
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
