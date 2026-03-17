package smokepod

import (
	"context"
	"strings"
	"testing"
	"time"
)

func TestDockerTarget_Exec(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Minute)
	defer cancel()

	container, err := NewContainer(ctx, ContainerConfig{
		Image: "alpine:latest",
	})
	if err != nil {
		t.Fatalf("NewContainer failed: %v", err)
	}
	defer func() { _ = container.Terminate(ctx) }()

	target := NewDockerTarget(container)

	result, err := target.Exec(ctx, "echo hello")
	if err != nil {
		t.Fatalf("Exec failed: %v", err)
	}

	if result.ExitCode != 0 {
		t.Errorf("ExitCode = %d, want 0", result.ExitCode)
	}

	if !strings.Contains(result.Stdout, "hello") {
		t.Errorf("Stdout = %q, want to contain %q", result.Stdout, "hello")
	}
}

func TestDockerTarget_ExecWithExitCode(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Minute)
	defer cancel()

	container, err := NewContainer(ctx, ContainerConfig{
		Image: "alpine:latest",
	})
	if err != nil {
		t.Fatalf("NewContainer failed: %v", err)
	}
	defer func() { _ = container.Terminate(ctx) }()

	target := NewDockerTarget(container)

	result, err := target.Exec(ctx, "exit 42")
	if err != nil {
		t.Fatalf("Exec failed: %v", err)
	}

	if result.ExitCode != 42 {
		t.Errorf("ExitCode = %d, want 42", result.ExitCode)
	}
}

func TestDockerTarget_ExecWithStderr(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Minute)
	defer cancel()

	container, err := NewContainer(ctx, ContainerConfig{
		Image: "alpine:latest",
	})
	if err != nil {
		t.Fatalf("NewContainer failed: %v", err)
	}
	defer func() { _ = container.Terminate(ctx) }()

	target := NewDockerTarget(container)

	result, err := target.Exec(ctx, "echo error >&2")
	if err != nil {
		t.Fatalf("Exec failed: %v", err)
	}

	if !strings.Contains(result.Stderr, "error") {
		t.Errorf("Stderr = %q, want to contain %q", result.Stderr, "error")
	}
}

func TestDockerTarget_Close(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Minute)
	defer cancel()

	container, err := NewContainer(ctx, ContainerConfig{
		Image: "alpine:latest",
	})
	if err != nil {
		t.Fatalf("NewContainer failed: %v", err)
	}

	target := NewDockerTarget(container)

	if err := target.Close(); err != nil {
		t.Errorf("Close returned error: %v", err)
	}
}
