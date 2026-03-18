package whitespace

import "testing"

func TestIsWhitespaceDiff(t *testing.T) {
	tests := []struct {
		name string
		a, b string
		want bool
	}{
		{"trailing space", "a ", "a", true},
		{"tab vs spaces", "\thello", "  hello", true},
		{"content differs", "a", "b", false},
		{"identical", "a", "a", false},
		{"leading space", " a", "a", true},
		{"both empty", "", "", false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := IsWhitespaceDiff(tt.a, tt.b)
			if got != tt.want {
				t.Errorf("IsWhitespaceDiff(%q, %q) = %v, want %v", tt.a, tt.b, got, tt.want)
			}
		})
	}
}

func TestRenderWhitespace(t *testing.T) {
	tests := []struct {
		name  string
		input string
		want  string
	}{
		{"space becomes dot", "a b", "a·b"},
		{"tab becomes arrow", "a\tb", "a→b"},
		{"cr becomes not", "a\rb", "a¬b"},
		{"no whitespace", "abc", "abc"},
		{"mixed content", "hello world\t", "hello·world→"},
		{"multiple spaces", "  a  ", "··a··"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := RenderWhitespace(tt.input)
			if got != tt.want {
				t.Errorf("RenderWhitespace(%q) = %q, want %q", tt.input, got, tt.want)
			}
		})
	}
}
