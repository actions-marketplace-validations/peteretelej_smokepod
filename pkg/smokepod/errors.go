package smokepod

import (
	"errors"
	"fmt"
)

var (
	ErrCIGuard         = errors.New("record called in CI without --update")
	ErrFixtureExists   = errors.New("fixture file already exists")
	ErrFixtureNotFound = errors.New("fixture file missing for verify")
	ErrProcessCrash    = errors.New("process died")
	ErrMalformedJSON   = errors.New("invalid JSON received")
)

// ConfigError represents a configuration validation error.
type ConfigError struct {
	Field   string
	Message string
	Hint    string
}

func (e *ConfigError) Error() string {
	if e.Hint != "" {
		return fmt.Sprintf("Invalid config: %s\n  %s", e.Message, e.Hint)
	}
	return fmt.Sprintf("Invalid config: %s", e.Message)
}

// DockerError represents a Docker-related error.
type DockerError struct {
	Op      string
	Message string
	Hint    string
}

func (e *DockerError) Error() string {
	if e.Hint != "" {
		return fmt.Sprintf("%s: %s\n  %s", e.Op, e.Message, e.Hint)
	}
	return fmt.Sprintf("%s: %s", e.Op, e.Message)
}

// Common error constructors for actionable messages.

// ErrDockerNotRunning returns an error indicating Docker is not available.
func ErrDockerNotRunning() error {
	return &DockerError{
		Op:      "Docker",
		Message: "Docker is not running",
		Hint:    "Start Docker Desktop or the Docker daemon and try again.",
	}
}

// ErrConfigNotFound returns an error for missing config file.
func ErrConfigNotFound(path string) error {
	return &ConfigError{
		Field:   "path",
		Message: fmt.Sprintf("config file not found: %s", path),
		Hint:    "Create a smokepod.yaml file or specify a different path.",
	}
}

// ErrImagePullFailed returns an error for failed image pulls.
func ErrImagePullFailed(image string) error {
	return &DockerError{
		Op:      "Image pull",
		Message: fmt.Sprintf("failed to pull image: %s", image),
		Hint:    "Check that the image exists and you have access.",
	}
}

// ErrTestTimeout returns an error for test timeouts.
func ErrTestTimeout(testName string, timeout string) error {
	return &DockerError{
		Op:      "Test execution",
		Message: fmt.Sprintf("timeout after %s: %s", timeout, testName),
		Hint:    "Increase the timeout or check why the test is slow.",
	}
}

// ErrMissingField returns an error for missing required config fields.
func ErrMissingField(testName, field string) error {
	hint := ""
	switch field {
	case "file":
		hint = "CLI tests require a \"file\" field pointing to a .test file."
	case "path":
		hint = "Playwright tests require a \"path\" field pointing to the test project."
	case "image":
		hint = "CLI tests require an \"image\" field specifying the Docker image."
	}
	return &ConfigError{
		Field:   field,
		Message: fmt.Sprintf("test %q missing required field %q", testName, field),
		Hint:    hint,
	}
}
