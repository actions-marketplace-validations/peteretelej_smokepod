package smokepod

import (
	"context"
	"os/exec"
	"strings"
	"testing"
	"time"
)

func dockerAvailable() bool {
	cmd := exec.Command("docker", "info")
	return cmd.Run() == nil
}

func TestNewContainer(t *testing.T) {
	if !dockerAvailable() {
		t.Skip("docker not available")
	}

	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Minute)
	defer cancel()

	c, err := NewContainer(ctx, ContainerConfig{
		Image: "alpine:latest",
	})
	if err != nil {
		t.Fatalf("NewContainer failed: %v", err)
	}
	defer func() { _ = c.Terminate(ctx) }()

	if c.container == nil {
		t.Error("container is nil")
	}
}

func TestContainer_Exec(t *testing.T) {
	if !dockerAvailable() {
		t.Skip("docker not available")
	}

	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Minute)
	defer cancel()

	c, err := NewContainer(ctx, ContainerConfig{
		Image: "alpine:latest",
	})
	if err != nil {
		t.Fatalf("NewContainer failed: %v", err)
	}
	defer func() { _ = c.Terminate(ctx) }()

	result, err := c.Exec(ctx, []string{"echo", "hello"})
	if err != nil {
		t.Fatalf("Exec failed: %v", err)
	}

	if result.ExitCode != 0 {
		t.Errorf("exit code = %d, want 0", result.ExitCode)
	}

	if !strings.Contains(result.Stdout, "hello") {
		t.Errorf("stdout = %q, want to contain %q", result.Stdout, "hello")
	}
}

func TestContainer_ExecExitCode(t *testing.T) {
	if !dockerAvailable() {
		t.Skip("docker not available")
	}

	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Minute)
	defer cancel()

	c, err := NewContainer(ctx, ContainerConfig{
		Image: "alpine:latest",
	})
	if err != nil {
		t.Fatalf("NewContainer failed: %v", err)
	}
	defer func() { _ = c.Terminate(ctx) }()

	result, err := c.Exec(ctx, []string{"sh", "-c", "exit 1"})
	if err != nil {
		t.Fatalf("Exec failed: %v", err)
	}

	if result.ExitCode != 1 {
		t.Errorf("exit code = %d, want 1", result.ExitCode)
	}
}

func TestContainer_Terminate(t *testing.T) {
	if !dockerAvailable() {
		t.Skip("docker not available")
	}

	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Minute)
	defer cancel()

	c, err := NewContainer(ctx, ContainerConfig{
		Image: "alpine:latest",
	})
	if err != nil {
		t.Fatalf("NewContainer failed: %v", err)
	}

	err = c.Terminate(ctx)
	if err != nil {
		t.Errorf("Terminate failed: %v", err)
	}
}

func TestContainer_ExecWithEnv(t *testing.T) {
	if !dockerAvailable() {
		t.Skip("docker not available")
	}

	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Minute)
	defer cancel()

	c, err := NewContainer(ctx, ContainerConfig{
		Image: "alpine:latest",
		Env: map[string]string{
			"MY_VAR": "testvalue",
		},
	})
	if err != nil {
		t.Fatalf("NewContainer failed: %v", err)
	}
	defer func() { _ = c.Terminate(ctx) }()

	result, err := c.Exec(ctx, []string{"sh", "-c", "echo $MY_VAR"})
	if err != nil {
		t.Fatalf("Exec failed: %v", err)
	}

	if !strings.Contains(result.Stdout, "testvalue") {
		t.Errorf("stdout = %q, want to contain %q", result.Stdout, "testvalue")
	}
}
