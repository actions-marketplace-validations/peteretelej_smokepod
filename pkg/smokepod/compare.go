package smokepod

import (
	"strings"
)

type CompareResult struct {
	Matched        bool
	Diff           string
	ExitCode       int
	ExitMatched    bool
	WhitespaceDiff bool
}

func CompareOutput(expected, actual string) CompareResult {
	expectedLines := splitLines(expected)
	actualLines := splitLines(actual)

	if len(expectedLines) != len(actualLines) {
		diff, wsDiff := formatDiff(expectedLines, actualLines)
		return CompareResult{
			Matched:        false,
			Diff:           diff,
			WhitespaceDiff: wsDiff,
		}
	}

	for i, exp := range expectedLines {
		if exp != actualLines[i] {
			diff, wsDiff := formatDiff(expectedLines, actualLines)
			return CompareResult{
				Matched:        false,
				Diff:           diff,
				WhitespaceDiff: wsDiff,
			}
		}
	}

	return CompareResult{
		Matched: true,
		Diff:    "",
	}
}

func CompareExitCode(expected, actual int) bool {
	return expected == actual
}

func splitLines(s string) []string {
	if s == "" {
		return nil
	}

	lines := strings.Split(s, "\n")

	if len(lines) > 0 && lines[len(lines)-1] == "" {
		lines = lines[:len(lines)-1]
	}

	return lines
}
