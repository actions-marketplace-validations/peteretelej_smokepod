package smokepod

import (
	"context"
	"os/exec"
)

// targetLaunchSpec holds the executable path and fixed arguments for a target.
type targetLaunchSpec struct {
	path string
	args []string
}

// shellArgs returns the argument list for shell-mode execution: [args..., "-c", command].
func (s targetLaunchSpec) shellArgs(command string) []string {
	out := make([]string, 0, len(s.args)+2)
	out = append(out, s.args...)
	out = append(out, "-c", command)
	return out
}

// processArgs returns the argument list for process-mode execution: [args...].
func (s targetLaunchSpec) processArgs() []string {
	return s.args
}

// cmd builds an exec.Cmd for this launch spec with the given extra args.
// It copies s.args before appending to avoid mutating the backing array.
func (s targetLaunchSpec) cmd(ctx context.Context, extraArgs ...string) *exec.Cmd {
	allArgs := make([]string, 0, len(s.args)+len(extraArgs))
	allArgs = append(allArgs, s.args...)
	allArgs = append(allArgs, extraArgs...)
	return exec.CommandContext(ctx, s.path, allArgs...)
}
