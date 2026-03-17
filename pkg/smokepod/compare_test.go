package smokepod

import (
	"strings"
	"testing"
)

func TestCompareOutput_Match(t *testing.T) {
	expected := "hello\nworld\n"
	actual := "hello\nworld\n"

	result := CompareOutput(expected, actual)

	if !result.Matched {
		t.Errorf("Expected match, got diff:\n%s", result.Diff)
	}
}

func TestCompareOutput_Mismatch(t *testing.T) {
	expected := "hello\nworld\n"
	actual := "hello\nuniverse\n"

	result := CompareOutput(expected, actual)

	if result.Matched {
		t.Error("Expected mismatch, got match")
	}

	if !strings.Contains(result.Diff, "-world") {
		t.Errorf("Diff should contain '-world', got:\n%s", result.Diff)
	}

	if !strings.Contains(result.Diff, "+universe") {
		t.Errorf("Diff should contain '+universe', got:\n%s", result.Diff)
	}
}

func TestCompareOutput_LineCountMismatch(t *testing.T) {
	expected := "line1\nline2\n"
	actual := "line1\n"

	result := CompareOutput(expected, actual)

	if result.Matched {
		t.Error("Expected mismatch due to line count")
	}
}

func TestCompareOutput_Empty(t *testing.T) {
	result := CompareOutput("", "")

	if !result.Matched {
		t.Error("Empty strings should match")
	}
}

func TestCompareExitCode_Match(t *testing.T) {
	if !CompareExitCode(0, 0) {
		t.Error("Exit code 0 should match 0")
	}

	if !CompareExitCode(42, 42) {
		t.Error("Exit code 42 should match 42")
	}
}

func TestCompareExitCode_Mismatch(t *testing.T) {
	if CompareExitCode(0, 1) {
		t.Error("Exit code 0 should not match 1")
	}

	if CompareExitCode(42, 0) {
		t.Error("Exit code 42 should not match 0")
	}
}

func TestFormatDiff(t *testing.T) {
	expected := []string{"line1", "line2", "line3"}
	actual := []string{"line1", "different", "line3"}

	diff := formatDiff(expected, actual)

	if !strings.Contains(diff, "--- expected") {
		t.Error("Diff should contain expected header")
	}

	if !strings.Contains(diff, "+++ actual") {
		t.Error("Diff should contain actual header")
	}

	if !strings.Contains(diff, "-line2") {
		t.Error("Diff should show removed line")
	}

	if !strings.Contains(diff, "+different") {
		t.Error("Diff should show added line")
	}
}
