package runners

import "context"

// ContainerExecutor can execute commands in a container.
type ContainerExecutor interface {
	Exec(ctx context.Context, cmd []string) (ExecResult, error)
}

// ExecResult holds the result of a command execution.
type ExecResult struct {
	ExitCode int
	Stdout   string
	Stderr   string
}

// SectionResult contains results for a test file section.
type SectionResult struct {
	Name     string          `json:"name"`
	Passed   bool            `json:"passed"`
	Commands []CommandResult `json:"commands"`
}

// CommandResult contains the result for a single command.
type CommandResult struct {
	Command  string `json:"command"`
	Line     int    `json:"line"`
	Expected string `json:"expected"`
	Actual   string `json:"actual"`
	Passed   bool   `json:"passed"`
	Error    string `json:"error,omitempty"`
}
