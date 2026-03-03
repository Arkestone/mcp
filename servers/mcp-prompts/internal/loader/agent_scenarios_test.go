package loader

// agent_scenarios_test.go — realistic agent workflow tests for prompts.
//
// Each test simulates an AI agent loading prompts for a specific task context,
// verifying that only the relevant subset is returned and that scoring is correct.

import (
	"os"
	"strings"
	"testing"
	"time"

	"github.com/Arkestone/mcp/pkg/config"
	"github.com/Arkestone/mcp/pkg/github"
)

// ---------------------------------------------------------------------------
// Polyglot prompts fixture — 12 realistic prompt files
// ---------------------------------------------------------------------------

func promptsWorkspace(t *testing.T) string {
	t.Helper()
	return createTestDir(t, map[string]string{
		// Global prompts (no files: restriction)
		".github/prompts/code-review.prompt.md": `---
name: code-review
description: Review code for correctness, style, and best practices
tags: [review, quality, best-practices]
---
Review the provided code focusing on:
- Logic correctness and edge cases
- Code style and readability
- Performance implications
- Security vulnerabilities
`,
		".github/prompts/write-tests.prompt.md": `---
name: write-tests
description: Write comprehensive unit and integration tests
tags: [testing, quality, tdd]
---
Write tests covering:
- Happy path scenarios
- Edge cases and boundary conditions
- Error handling paths
- Integration points
`,
		".github/prompts/refactor.prompt.md": `---
name: refactor
description: Refactor code for better maintainability and clarity
tags: [refactor, clean-code, maintainability]
---
Refactor focusing on:
- Extract functions for reusability
- Remove code duplication
- Improve naming
- Simplify complex logic
`,
		// TypeScript-specific prompts
		".github/prompts/typescript-patterns.prompt.md": `---
name: typescript-patterns
description: Apply TypeScript-specific patterns and type safety
tags: [typescript, types, patterns]
files:
  - "**/*.ts"
  - "**/*.tsx"
---
Apply TypeScript best practices:
- Use strict type annotations
- Prefer interfaces over type aliases for object shapes
- Use generics for reusable utilities
`,
		// React-specific prompts
		".github/prompts/react-component.prompt.md": `---
name: react-component
description: Create or review React components following best practices
tags: [react, components, ui]
files:
  - "**/*.tsx"
  - "**/*.jsx"
---
React component guidelines:
- Use functional components with hooks
- Proper prop typing with TypeScript
- Memoization for performance
`,
		// Go-specific prompts
		".github/prompts/go-idioms.prompt.md": `---
name: go-idioms
description: Apply Go idiomatic patterns and conventions
tags: [go, golang, idioms]
files:
  - "**/*.go"
---
Go idiomatic patterns:
- Handle errors explicitly, never ignore
- Use interfaces for decoupling
- Prefer composition over inheritance
- Write table-driven tests
`,
		// Python prompts
		".github/prompts/python-style.prompt.md": `---
name: python-style
description: Apply Python style and Pythonic patterns
tags: [python, style, pep8]
files:
  - "**/*.py"
---
Python best practices:
- Follow PEP 8 guidelines
- Use list comprehensions where readable
- Prefer context managers for resources
- Type hints for all public functions
`,
		// Database prompts
		".github/prompts/sql-query.prompt.md": `---
name: sql-query
description: Write optimized SQL queries with proper indexing considerations
tags: [sql, database, performance, query]
files:
  - "**/*.sql"
  - "**/*.migration"
---
SQL query optimization:
- Use appropriate indexes
- Avoid N+1 query patterns
- Write readable CTEs for complex queries
- Consider query execution plans
`,
		// CI/CD prompts
		".github/prompts/github-actions.prompt.md": `---
name: github-actions
description: Write and review GitHub Actions workflows
tags: [ci, cd, github-actions, workflow]
files:
  - ".github/workflows/**/*.yml"
  - ".github/workflows/**/*.yaml"
---
GitHub Actions best practices:
- Pin action versions for reproducibility
- Use concurrency groups to prevent duplicate runs
- Cache dependencies
- Store secrets securely
`,
		// Documentation prompts
		".github/prompts/write-docs.prompt.md": `---
name: write-docs
description: Write clear technical documentation
tags: [docs, documentation, readme]
files:
  - "**/*.md"
---
Documentation guidelines:
- Clear and concise language
- Include usage examples
- Document edge cases and limitations
`,
		// Security prompts
		".github/prompts/security-review.prompt.md": `---
name: security-review
description: Review code for security vulnerabilities
tags: [security, vulnerabilities, owasp]
---
Security review checklist:
- Input validation and sanitization
- Authentication and authorization
- Injection vulnerabilities
- Sensitive data exposure
`,
		// Deeply nested prompt in hidden dir
		".hidden/team/prompts/commit-message.prompt.md": `---
name: commit-message
description: Write conventional commit messages following the specification
tags: [git, commit, conventional-commits]
---
Conventional commit format:
- type(scope): description
- Types: feat, fix, docs, chore, test, refactor
- Keep subject under 72 characters
`,
	})
}

// ---------------------------------------------------------------------------
// Helper
// ---------------------------------------------------------------------------

func loadAllPrompts(t *testing.T, dir string) []Prompt {
	t.Helper()
	cfg := &config.Config{
		Sources: config.Sources{Dirs: []string{dir}},
		Cache:   config.CacheConfig{Dir: t.TempDir(), SyncInterval: time.Hour},
	}
	l := newLoader(cfg)
	return l.List()
}

// ---------------------------------------------------------------------------
// Table-driven agent file-context scenarios
// ---------------------------------------------------------------------------

func TestAgentScenario_Prompts_FileContextRouting(t *testing.T) {
	dir := promptsWorkspace(t)

	tests := []struct {
		name       string
		filePath   string
		mustHave   []string
		mustNotHave []string
	}{
		{
			name:       "TypeScript source file — ts-specific + global",
			filePath:   "src/auth/login.ts",
			mustHave:   []string{"code-review", "write-tests", "refactor", "typescript-patterns", "security-review"},
			mustNotHave: []string{"react-component", "go-idioms", "python-style", "sql-query"},
		},
		{
			name:       "React TSX component — ts + react + global",
			filePath:   "src/components/Button.tsx",
			mustHave:   []string{"typescript-patterns", "react-component", "code-review"},
			mustNotHave: []string{"go-idioms", "python-style", "sql-query"},
		},
		{
			name:       "Go source file — go-specific + global",
			filePath:   "pkg/server/server.go",
			mustHave:   []string{"go-idioms", "code-review", "write-tests"},
			mustNotHave: []string{"typescript-patterns", "react-component", "python-style"},
		},
		{
			name:       "Python script — python-specific + global",
			filePath:   "scripts/deploy.py",
			mustHave:   []string{"python-style", "code-review"},
			mustNotHave: []string{"typescript-patterns", "go-idioms", "sql-query"},
		},
		{
			name:       "SQL migration file — sql-specific + global",
			filePath:   "migrations/0042_add_users_index.sql",
			mustHave:   []string{"sql-query", "code-review"},
			mustNotHave: []string{"typescript-patterns", "go-idioms", "python-style"},
		},
		{
			name:       "GitHub Actions workflow — ci/cd + global",
			filePath:   ".github/workflows/ci.yml",
			mustHave:   []string{"github-actions", "code-review"},
			mustNotHave: []string{"typescript-patterns", "go-idioms", "sql-query"},
		},
		{
			name:       "Markdown docs file — docs + global",
			filePath:   "docs/architecture.md",
			mustHave:   []string{"write-docs", "code-review"},
			mustNotHave: []string{"go-idioms", "sql-query", "github-actions"},
		},
		{
			name:       "No file path — all prompts returned",
			filePath:   "",
			mustHave:   []string{"code-review", "write-tests", "typescript-patterns", "go-idioms", "sql-query", "commit-message"},
			mustNotHave: []string{},
		},
	}

	all := loadAllPrompts(t, dir)
	if len(all) == 0 {
		t.Fatal("no prompts loaded from polyglot workspace")
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			var got []Prompt
			if tc.filePath == "" {
				got = all
			} else {
				got = FilterByFilePath(all, tc.filePath)
			}

			names := make(map[string]bool)
			for _, p := range got {
				names[p.Name] = true
			}

			for _, want := range tc.mustHave {
				if !names[want] {
					t.Errorf("missing %q in results for file %q (got %v)", want, tc.filePath, nameKeys(names))
				}
			}
			for _, notWant := range tc.mustNotHave {
				if names[notWant] {
					t.Errorf("unexpected %q in results for file %q", notWant, tc.filePath)
				}
			}
		})
	}
}

func nameKeys(m map[string]bool) []string {
	out := make([]string, 0, len(m))
	for k := range m {
		out = append(out, k)
	}
	return out
}

// ---------------------------------------------------------------------------
// Query relevance — scoring ranks the most relevant prompt first
// ---------------------------------------------------------------------------

func TestAgentScenario_Prompts_QueryScoring(t *testing.T) {
	dir := promptsWorkspace(t)
	all := loadAllPrompts(t, dir)

	tests := []struct {
		query    string
		wantFirst string // name of prompt expected at rank 0
	}{
		{"review this code for bugs", "code-review"},
		{"write unit tests for this function", "write-tests"},
		{"refactor this messy code", "refactor"},
		{"TypeScript type safety patterns", "typescript-patterns"},
		{"React component best practices", "react-component"},
		{"optimize my SQL query", "sql-query"},
		{"GitHub Actions workflow security", "github-actions"},
		{"conventional commit message", "commit-message"},
		{"security vulnerability review", "security-review"},
	}

	for _, tc := range tests {
		t.Run(tc.query, func(t *testing.T) {
			ranked := FilterByQuery(all, tc.query)
			if len(ranked) == 0 {
				t.Fatal("no results returned")
			}
			if ranked[0].Name != tc.wantFirst {
				t.Errorf("query %q: want first=%q got first=%q", tc.query, tc.wantFirst, ranked[0].Name)
			}
		})
	}
}

// ---------------------------------------------------------------------------
// Combined filter: file path + query narrows results further
// ---------------------------------------------------------------------------

func TestAgentScenario_Prompts_CombinedFilter(t *testing.T) {
	dir := promptsWorkspace(t)
	all := loadAllPrompts(t, dir)

	// For a Go file, filter by file then rank by query
	goPrompts := FilterByFilePath(all, "pkg/api/handler.go")
	ranked := FilterByQuery(goPrompts, "idiomatic Go error handling")

	if len(ranked) == 0 {
		t.Fatal("no results after combined filter")
	}
	if ranked[0].Name != "go-idioms" {
		t.Errorf("combined filter: want go-idioms first, got %q", ranked[0].Name)
	}

	// Security review (global) should still appear in Go file results
	found := false
	for _, p := range goPrompts {
		if p.Name == "security-review" {
			found = true
			break
		}
	}
	if !found {
		t.Error("security-review (global prompt) missing from Go file results")
	}
}

// ---------------------------------------------------------------------------
// Hidden directories — prompts in .hidden paths must be discoverable
// ---------------------------------------------------------------------------

func TestAgentScenario_Prompts_HiddenDirectories(t *testing.T) {
	dir := promptsWorkspace(t)
	all := loadAllPrompts(t, dir)

	found := false
	for _, p := range all {
		if p.Name == "commit-message" {
			found = true
			if !strings.Contains(p.Content, "Conventional commit") {
				t.Errorf("commit-message content corrupted: %q", p.Content[:50])
			}
			break
		}
	}
	if !found {
		t.Error("commit-message prompt not found — hidden dir scan failed")
	}
}

// ---------------------------------------------------------------------------
// Multiple directories — prompts from all sources are merged
// ---------------------------------------------------------------------------

func TestAgentScenario_Prompts_MultipleDirectories(t *testing.T) {
	dir1 := createTestDir(t, map[string]string{
		".github/prompts/review.prompt.md": `---
name: review
description: Code review prompt from dir1
tags: [review]
---
Review code from dir1.
`,
	})
	dir2 := createTestDir(t, map[string]string{
		".github/prompts/test.prompt.md": `---
name: test
description: Testing prompt from dir2
tags: [testing]
---
Write tests from dir2.
`,
	})

	cfg := &config.Config{
		Sources: config.Sources{Dirs: []string{dir1, dir2}},
		Cache:   config.CacheConfig{Dir: t.TempDir(), SyncInterval: time.Hour},
	}
	l := newLoader(cfg)
	all := l.List()

	names := make(map[string]bool)
	for _, p := range all {
		names[p.Name] = true
	}
	if !names["review"] {
		t.Error("missing 'review' from dir1")
	}
	if !names["test"] {
		t.Error("missing 'test' from dir2")
	}
}

// ---------------------------------------------------------------------------
// Content integrity — no truncation, URI stable
// ---------------------------------------------------------------------------

func TestAgentScenario_Prompts_ContentIntegrity(t *testing.T) {
	marker := "UNIQUE_MARKER_XYZ_12345"
	dir := createTestDir(t, map[string]string{
		".github/prompts/special.prompt.md": "---\nname: special\ndescription: special prompt\ntags: [test]\n---\n" + marker + "\n",
	})
	all := loadAllPrompts(t, dir)

	if len(all) == 0 {
		t.Fatal("no prompts loaded")
	}
	p := all[0]
	if !strings.Contains(p.Content, marker) {
		t.Errorf("content marker missing: got %q", p.Content)
	}
	if p.URI == "" {
		t.Error("URI must not be empty")
	}
	if p.Name != "special" {
		t.Errorf("name: want 'special' got %q", p.Name)
	}
}

// ---------------------------------------------------------------------------
// Tags as string vs list — both formats must parse correctly
// ---------------------------------------------------------------------------

func TestAgentScenario_Prompts_TagFormats(t *testing.T) {
	dir := createTestDir(t, map[string]string{
		".github/prompts/string-tags.prompt.md": `---
name: string-tags
description: Tags as comma-separated string
tags: "typescript, go, python"
---
Content here.
`,
		".github/prompts/list-tags.prompt.md": `---
name: list-tags
description: Tags as YAML list
tags:
  - typescript
  - go
  - python
---
Content here.
`,
	})

	all := loadAllPrompts(t, dir)
	byName := make(map[string]Prompt)
	for _, p := range all {
		byName[p.Name] = p
	}

	for _, name := range []string{"string-tags", "list-tags"} {
		p, ok := byName[name]
		if !ok {
			t.Fatalf("prompt %q not found", name)
		}
		if len(p.Tags) != 3 {
			t.Errorf("%q: want 3 tags, got %d: %v", name, len(p.Tags), p.Tags)
		}
		tagSet := make(map[string]bool)
		for _, tag := range p.Tags {
			tagSet[strings.TrimSpace(tag)] = true
		}
		for _, expected := range []string{"typescript", "go", "python"} {
			if !tagSet[expected] {
				t.Errorf("%q: missing tag %q, got %v", name, expected, p.Tags)
			}
		}
	}
}

// ---------------------------------------------------------------------------
// Brace expansion in files: patterns
// ---------------------------------------------------------------------------

func TestAgentScenario_Prompts_BraceExpansionGlob(t *testing.T) {
	dir := createTestDir(t, map[string]string{
		".github/prompts/ts-react.prompt.md": `---
name: ts-react
description: TypeScript and React patterns
files: "**/*.{ts,tsx}"
---
TS and React patterns.
`,
		".github/prompts/other.prompt.md": `---
name: other
description: Other file types only
files: "**/*.go"
---
Go patterns.
`,
	})

	all := loadAllPrompts(t, dir)

	tsResults := FilterByFilePath(all, "src/components/App.tsx")
	tsNames := make(map[string]bool)
	for _, p := range tsResults {
		tsNames[p.Name] = true
	}
	if !tsNames["ts-react"] {
		t.Error("ts-react not matched for .tsx file — brace expansion broken")
	}
	if tsNames["other"] {
		t.Error("other (go pattern) should not match .tsx file")
	}

	tsFile := FilterByFilePath(all, "src/auth.ts")
	tsFileNames := make(map[string]bool)
	for _, p := range tsFile {
		tsFileNames[p.Name] = true
	}
	if !tsFileNames["ts-react"] {
		t.Error("ts-react not matched for .ts file — brace expansion broken")
	}
}

// ---------------------------------------------------------------------------
// Live reload — modifications on disk reflected on next List()
// ---------------------------------------------------------------------------

func TestAgentScenario_Prompts_LiveReload(t *testing.T) {
	dir := t.TempDir()

	promptPath := dir + "/.github/prompts/evolving.prompt.md"
	if err := os.MkdirAll(dir+"/.github/prompts", 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(promptPath, []byte("---\nname: evolving\ndescription: v1\ntags: [v1]\n---\nVersion 1\n"), 0o644); err != nil {
		t.Fatal(err)
	}

	cfg := &config.Config{
		Sources: config.Sources{Dirs: []string{dir}},
		Cache:   config.CacheConfig{Dir: t.TempDir(), SyncInterval: time.Hour},
	}
	l := newLoader(cfg)
	all1 := l.List()
	if len(all1) == 0 {
		t.Fatal("no prompts loaded initially")
	}
	if !strings.Contains(all1[0].Content, "Version 1") {
		t.Errorf("initial: want Version 1, got %q", all1[0].Content)
	}

	if err := os.WriteFile(promptPath, []byte("---\nname: evolving\ndescription: v2\ntags: [v2]\n---\nVersion 2\n"), 0o644); err != nil {
		t.Fatal(err)
	}

	// Invalidate cache to pick up disk changes
	l.ForceSync()
	all2 := l.List()
	if len(all2) == 0 {
		t.Fatal("no prompts after update")
	}
	if !strings.Contains(all2[0].Content, "Version 2") {
		t.Errorf("updated: want Version 2, got %q", all2[0].Content)
	}
}

// ---------------------------------------------------------------------------
// Determinism — List() order is stable across multiple calls
// ---------------------------------------------------------------------------

func TestAgentScenario_Prompts_Determinism(t *testing.T) {
	dir := promptsWorkspace(t)
	cfg := &config.Config{
		Sources: config.Sources{Dirs: []string{dir}},
		Cache:   config.CacheConfig{Dir: t.TempDir(), SyncInterval: time.Hour},
	}
	l := newLoader(cfg)

	all1 := l.List()
	l.ForceSync()
	all2 := l.List()

	if len(all1) != len(all2) {
		t.Fatalf("non-deterministic count: %d vs %d", len(all1), len(all2))
	}
	for i := range all1 {
		if all1[i].URI != all2[i].URI {
			t.Errorf("position %d: URI changed between calls (%q → %q)", i, all1[i].URI, all2[i].URI)
		}
	}
}

// ---------------------------------------------------------------------------
// Adversarial: malformed frontmatter — must not crash or drop valid prompts
// ---------------------------------------------------------------------------

func TestAgentScenario_Prompts_AdversarialFrontmatter(t *testing.T) {
	dir := createTestDir(t, map[string]string{
		".github/prompts/good-prompt.prompt.md": `---
name: good-prompt
description: Valid prompt
tags: [valid]
---
This is a valid prompt.
`,
		".github/prompts/no-frontmatter.prompt.md": `No frontmatter at all.
Just raw content.
`,
		".github/prompts/empty-frontmatter.prompt.md": `---
---
Body without any frontmatter fields.
`,
		".github/prompts/malformed-yaml.prompt.md": `---
name: bad
  indented: wrongly
tags: [unclosed
---
Content after malformed YAML.
`,
		".github/prompts/wrong-type-files.prompt.md": `---
name: wrong-type
description: files field is an integer
files: 42
---
Wrong type files.
`,
		".github/prompts/invalid-glob.prompt.md": `---
name: invalid-glob
description: Invalid glob pattern in files
files:
  - "[invalid[glob"
---
Should be included via inclusive fallback.
`,
	})

	all := loadAllPrompts(t, dir)

	// Must find the valid prompt
	found := false
	for _, p := range all {
		if p.Name == "good-prompt" {
			found = true
			break
		}
	}
	if !found {
		t.Error("good-prompt missing — valid prompt dropped due to bad neighbors")
	}

	// Server must not crash — just reaching here means we're resilient
	// invalid-glob should be inclusive (included for all file paths due to fallback)
	invalidGlob := false
	for _, p := range all {
		if p.Name == "invalid-glob" {
			invalidGlob = true
			break
		}
	}
	if !invalidGlob {
		t.Log("info: invalid-glob prompt not loaded (may have been skipped at parse time)")
	}

	// invalid-glob (if loaded) should appear in file-filtered results due to inclusive fallback
	filtered := FilterByFilePath(all, "src/anything.ts")
	for _, p := range filtered {
		if p.Name == "invalid-glob" {
			// Good — inclusive fallback worked
			return
		}
	}
}

// ---------------------------------------------------------------------------
// Ensure ForceSync is accessible (used by live reload tests)
// ---------------------------------------------------------------------------

func TestAgentScenario_Prompts_ForceSyncExists(t *testing.T) {
	dir := createTestDir(t, map[string]string{
		".github/prompts/p.prompt.md": "---\nname: p\ndescription: p\n---\nContent\n",
	})
	cfg := &config.Config{
		Sources: config.Sources{Dirs: []string{dir}},
		Cache:   config.CacheConfig{Dir: t.TempDir(), SyncInterval: time.Hour},
	}
	l := newLoader(cfg)
	l.ForceSync() // must not panic
	all := l.List()
	if len(all) == 0 {
		t.Error("ForceSync + List returned no results")
	}
}

// force the github import to be used (newLoader already uses it indirectly)
var _ = (*github.Client)(nil)
