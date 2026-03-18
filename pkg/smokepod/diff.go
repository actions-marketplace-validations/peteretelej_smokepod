package smokepod

import (
	"fmt"
	"strings"

	"github.com/peteretelej/smokepod/internal/whitespace"
)

func formatDiff(expected, actual []string) (string, bool) {
	var b strings.Builder
	hasWSDiff := false

	b.WriteString("--- expected\n")
	b.WriteString("+++ actual\n")

	hdr := formatHunkHeader(1, len(expected), 1, len(actual))
	b.WriteString(hdr)
	b.WriteString("\n")

	maxLen := len(expected)
	if len(actual) > maxLen {
		maxLen = len(actual)
	}

	for i := 0; i < maxLen; i++ {
		if i < len(expected) && i < len(actual) {
			if expected[i] == actual[i] {
				b.WriteString(" ")
				b.WriteString(expected[i])
				b.WriteString("\n")
			} else if whitespace.IsWhitespaceDiff(expected[i], actual[i]) {
				hasWSDiff = true
				b.WriteString("-")
				b.WriteString(whitespace.RenderWhitespace(expected[i]))
				b.WriteString("\n")
				b.WriteString("+")
				b.WriteString(whitespace.RenderWhitespace(actual[i]))
				b.WriteString("\n")
			} else {
				b.WriteString("-")
				b.WriteString(expected[i])
				b.WriteString("\n")
				b.WriteString("+")
				b.WriteString(actual[i])
				b.WriteString("\n")
			}
		} else if i < len(expected) {
			b.WriteString("-")
			b.WriteString(expected[i])
			b.WriteString("\n")
		} else {
			b.WriteString("+")
			b.WriteString(actual[i])
			b.WriteString("\n")
		}
	}

	return b.String(), hasWSDiff
}

func formatHunkHeader(expStart, expLen, actStart, actLen int) string {
	expRange := formatRange(expStart, expLen)
	actRange := formatRange(actStart, actLen)
	return fmt.Sprintf("@@ -%s +%s @@", expRange, actRange)
}

func formatRange(start, length int) string {
	if length == 1 {
		return fmt.Sprintf("%d", start)
	}
	return fmt.Sprintf("%d,%d", start, length)
}
