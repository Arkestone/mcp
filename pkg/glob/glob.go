// Package glob provides glob pattern matching with ** support and brace expansion.
// It is used to match file paths against VS Code-style glob patterns (e.g. applyTo).
package glob

import (
	"path"
	"strings"
)

// Match reports whether name matches the shell glob pattern.
// Unlike path.Match, it supports ** as a wildcard spanning multiple path
// segments and {a,b} brace expansion (e.g. "**/*.{go,ts}").
// name and pattern should use forward slashes.
func Match(pattern, name string) bool {
	pattern = forwardSlash(pattern)
	name = forwardSlash(name)
	for _, p := range expandBraces(pattern) {
		if matchSingle(p, name) {
			return true
		}
	}
	return false
}

// MatchAny reports whether name matches any of the patterns.
func MatchAny(patterns []string, name string) bool {
	for _, p := range patterns {
		if Match(p, name) {
			return true
		}
	}
	return false
}

func forwardSlash(s string) string { return strings.ReplaceAll(s, "\\", "/") }

// matchSingle matches one (already brace-expanded) pattern against name.
func matchSingle(pattern, name string) bool {
	if !strings.Contains(pattern, "**") {
		m, err := path.Match(pattern, name)
		// Invalid pattern (e.g. unclosed bracket) → treat as wildcard to avoid silently hiding items.
		if err != nil {
			return true
		}
		return m
	}
	i := strings.Index(pattern, "**")
	prefix := pattern[:i]
	rest := strings.TrimPrefix(pattern[i+2:], "/")
	if prefix != "" {
		if !strings.HasPrefix(name, prefix) {
			return false
		}
		name = strings.TrimPrefix(name[len(prefix):], "/")
	}
	if rest == "" {
		return true
	}
	// Try matching rest against every possible suffix of name.
	segments := strings.Split(name, "/")
	for i := range segments {
		sub := strings.Join(segments[i:], "/")
		if matchSingle(rest, sub) {
			return true
		}
	}
	return false
}

// expandBraces expands brace alternatives: "*.{go,ts}" → ["*.go", "*.ts"].
func expandBraces(pattern string) []string {
	start := strings.Index(pattern, "{")
	if start < 0 {
		return []string{pattern}
	}
	end := -1
	depth := 0
	for i := start; i < len(pattern); i++ {
		switch pattern[i] {
		case '{':
			depth++
		case '}':
			depth--
			if depth == 0 {
				end = i
			}
		}
		if end >= 0 {
			break
		}
	}
	if end < 0 {
		return []string{pattern}
	}
	pre, suf := pattern[:start], pattern[end+1:]
	alts := strings.Split(pattern[start+1:end], ",")
	var out []string
	for _, alt := range alts {
		for _, x := range expandBraces(pre + strings.TrimSpace(alt) + suf) {
			out = append(out, x)
		}
	}
	return out
}
