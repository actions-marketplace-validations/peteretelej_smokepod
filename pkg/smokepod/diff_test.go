package smokepod

import (
	"strings"
	"testing"
)

func TestFormatDiff_WhitespaceOnly(t *testing.T) {
	expected := []string{"hello "}
	actual := []string{"hello"}

	diff, wsDiff := formatDiff(expected, actual)

	if !wsDiff {
		t.Error("Expected wsDiff=true for trailing space difference")
	}

	if !strings.Contains(diff, "·") {
		t.Errorf("Expected · marker in diff output, got:\n%s", diff)
	}
}

func TestFormatDiff_TabVsSpace(t *testing.T) {
	expected := []string{"\thello"}
	actual := []string{"  hello"}

	diff, wsDiff := formatDiff(expected, actual)

	if !wsDiff {
		t.Error("Expected wsDiff=true for tab vs space difference")
	}

	if !strings.Contains(diff, "→") {
		t.Errorf("Expected → marker for tab, got:\n%s", diff)
	}

	if !strings.Contains(diff, "·") {
		t.Errorf("Expected · marker for spaces, got:\n%s", diff)
	}
}

func TestFormatDiff_ContentDiffOnly(t *testing.T) {
	expected := []string{"hello"}
	actual := []string{"world"}

	diff, wsDiff := formatDiff(expected, actual)

	if wsDiff {
		t.Error("Expected wsDiff=false for content-only difference")
	}

	if strings.Contains(diff, "·") || strings.Contains(diff, "→") || strings.Contains(diff, "¬") {
		t.Errorf("Expected no whitespace markers for content diff, got:\n%s", diff)
	}
}

func TestFormatDiff_Mixed(t *testing.T) {
	expected := []string{"same", "hello ", "different"}
	actual := []string{"same", "hello", "changed"}

	diff, wsDiff := formatDiff(expected, actual)

	if !wsDiff {
		t.Error("Expected wsDiff=true when at least one pair differs only by whitespace")
	}

	// The whitespace pair should have markers
	if !strings.Contains(diff, "·") {
		t.Errorf("Expected · marker for whitespace pair, got:\n%s", diff)
	}

	// The content pair should appear raw
	if !strings.Contains(diff, "-different") {
		t.Errorf("Expected raw '-different' for content pair, got:\n%s", diff)
	}

	if !strings.Contains(diff, "+changed") {
		t.Errorf("Expected raw '+changed' for content pair, got:\n%s", diff)
	}
}

func TestFormatDiff_CarriageReturn(t *testing.T) {
	expected := []string{"hello\r"}
	actual := []string{"hello"}

	diff, wsDiff := formatDiff(expected, actual)

	if !wsDiff {
		t.Error("Expected wsDiff=true for carriage return difference")
	}

	if !strings.Contains(diff, "¬") {
		t.Errorf("Expected ¬ marker for CR, got:\n%s", diff)
	}
}

func TestFormatDiff_UnpairedLines(t *testing.T) {
	expected := []string{"line1", "line2", "extra"}
	actual := []string{"line1", "line2"}

	diff, _ := formatDiff(expected, actual)

	// Unpaired lines should not have whitespace markers
	if strings.Contains(diff, "·") || strings.Contains(diff, "→") || strings.Contains(diff, "¬") {
		t.Errorf("Expected no whitespace markers on unpaired lines, got:\n%s", diff)
	}

	if !strings.Contains(diff, "-extra") {
		t.Errorf("Expected '-extra' for unpaired line, got:\n%s", diff)
	}
}
