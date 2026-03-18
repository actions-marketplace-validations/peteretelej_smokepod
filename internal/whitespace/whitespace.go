package whitespace

import "strings"

// IsWhitespaceDiff returns true when a and b have the same trimmed content
// but differ in their original form (trailing spaces, tabs vs spaces, etc.).
func IsWhitespaceDiff(a, b string) bool {
	return strings.TrimSpace(a) == strings.TrimSpace(b) && a != b
}

// RenderWhitespace replaces invisible whitespace characters with visible
// Unicode markers: space -> · (U+00B7), tab -> → (U+2192), \r -> ¬ (U+00AC).
func RenderWhitespace(s string) string {
	r := strings.NewReplacer(" ", "·", "\t", "→", "\r", "¬")
	return r.Replace(s)
}
