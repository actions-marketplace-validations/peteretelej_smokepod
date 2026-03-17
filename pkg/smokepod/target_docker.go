package smokepod

import (
	"bytes"
	"context"
	"fmt"
	"io"

	"github.com/peteretelej/smokepod/pkg/smokepod/runners"
)

type DockerTarget struct {
	container *Container
}

func NewDockerTarget(container *Container) *DockerTarget {
	return &DockerTarget{container: container}
}

func (t *DockerTarget) Exec(ctx context.Context, command string) (runners.ExecResult, error) {
	exitCode, reader, err := t.container.container.Exec(ctx, []string{"sh", "-c", command})
	if err != nil {
		return runners.ExecResult{}, fmt.Errorf("executing command: %w", err)
	}

	output, err := io.ReadAll(reader)
	if err != nil {
		return runners.ExecResult{}, fmt.Errorf("reading output: %w", err)
	}

	stdout, stderr := demultiplexDockerStream(output)

	return runners.ExecResult{
		ExitCode: exitCode,
		Stdout:   stdout,
		Stderr:   stderr,
	}, nil
}

func (t *DockerTarget) Close() error {
	return t.container.Terminate(context.Background())
}

func demultiplexDockerStream(data []byte) (stdout, stderr string) {
	var stdoutBuf, stderrBuf bytes.Buffer

	for len(data) >= 8 {
		streamType := data[0]
		size := int(data[4])<<24 | int(data[5])<<16 | int(data[6])<<8 | int(data[7])
		data = data[8:]

		if len(data) < size {
			break
		}

		payload := data[:size]
		data = data[size:]

		switch streamType {
		case 1:
			stdoutBuf.Write(payload)
		case 2:
			stderrBuf.Write(payload)
		}
	}

	if stdoutBuf.Len() == 0 && stderrBuf.Len() == 0 {
		return string(data), ""
	}

	return stdoutBuf.String(), stderrBuf.String()
}
