package glob_test

import (
	"testing"

	"github.com/Arkestone/mcp/pkg/glob"
)

func TestMatch(t *testing.T) {
	tests := []struct {
		pattern string
		name    string
		want    bool
	}{
		{"**/*.go", "main.go", true},
		{"**/*.go", "pkg/server/main.go", true},
		{"**/*.go", "pkg/server/main.ts", false},
		{"**/*.{go,ts}", "main.go", true},
		{"**/*.{go,ts}", "main.ts", true},
		{"**/*.{go,ts}", "main.py", false},
		{"src/**/*.go", "src/pkg/main.go", true},
		{"src/**/*.go", "other/pkg/main.go", false},
		{"*.md", "README.md", true},
		{"*.md", "dir/README.md", false},
		{"**", "anything/at/all", true},
		{"**/*.instructions.md", ".github/instructions/go.instructions.md", true},
		{"**/*.instructions.md", ".github/instructions/go.instructions.ts", false},
	}
	for _, tt := range tests {
		got := glob.Match(tt.pattern, tt.name)
		if got != tt.want {
			t.Errorf("Match(%q, %q) = %v, want %v", tt.pattern, tt.name, got, tt.want)
		}
	}
}

func TestMatchAny(t *testing.T) {
	patterns := []string{"**/*.go", "**/*.ts"}
	if !glob.MatchAny(patterns, "main.go") {
		t.Error("expected match for main.go")
	}
	if !glob.MatchAny(patterns, "src/app.ts") {
		t.Error("expected match for src/app.ts")
	}
	if glob.MatchAny(patterns, "README.md") {
		t.Error("expected no match for README.md")
	}
}

func TestMatch_InvalidPattern_ReturnsTrue(t *testing.T) {
	// An invalid glob (unclosed bracket) should match everything rather than silently hide items.
	got := glob.Match("[invalid[glob", "anything.ts")
	if !got {
		t.Error("Match with invalid pattern should return true (inclusive fallback)")
	}
}

func TestMatchAny_AllInvalidPatterns_ReturnsTrue(t *testing.T) {
	// If all patterns are invalid, the item should not be hidden.
	got := glob.MatchAny([]string{"[bad[", "[also[bad["}, "src/file.ts")
	if !got {
		t.Error("all-invalid patterns should return true (inclusive fallback)")
	}
}
