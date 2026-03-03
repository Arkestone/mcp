package scanner

// agent_scenarios_test.go — realistic agent workflow tests for skills.
//
// Each test simulates an AI agent loading skills for a specific context,
// verifying correct filtering and scoring.

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/Arkestone/mcp/pkg/config"
	"github.com/Arkestone/mcp/pkg/github"
)

// createTestDir creates a temp directory populated from a map of relative path → content.
func createTestDir(t *testing.T, files map[string]string) string {
	t.Helper()
	dir := t.TempDir()
	for path, content := range files {
		fullPath := filepath.Join(dir, path)
		if err := os.MkdirAll(filepath.Dir(fullPath), 0o755); err != nil {
			t.Fatal(err)
		}
		if err := os.WriteFile(fullPath, []byte(content), 0o644); err != nil {
			t.Fatal(err)
		}
	}
	return dir
}

// ---------------------------------------------------------------------------
// Polyglot skills fixture — 12 realistic SKILL.md files
// ---------------------------------------------------------------------------

func skillsWorkspace(t *testing.T) string {
	t.Helper()
	return createTestDir(t, map[string]string{
		// Global skills (no files: restriction)
		".github/skills/git-workflow/SKILL.md": `---
name: git-workflow
description: Git branching strategy, commit conventions, and PR workflow
tags: [git, workflow, version-control]
---
# Git Workflow

Use feature branches for all changes.
Commit messages follow Conventional Commits spec.
PR must pass CI before merge.
`,
		".github/skills/error-handling/SKILL.md": `---
name: error-handling
description: Patterns for robust error handling across languages
tags: [errors, resilience, best-practices]
---
# Error Handling

Never silently ignore errors.
Wrap errors with context at each layer boundary.
Return errors up the call stack, don't log and return.
`,
		".github/skills/security-patterns/SKILL.md": `---
name: security-patterns
description: Security patterns for authentication, authorization, and input validation
tags: [security, auth, validation, owasp]
---
# Security Patterns

Validate all inputs at trust boundaries.
Use parameterized queries for SQL.
Never store secrets in code or logs.
`,
		// TypeScript/JS skills
		".github/skills/typescript-patterns/SKILL.md": `---
name: typescript-patterns
description: TypeScript-specific patterns, generics, and type safety
tags: [typescript, types, generics]
files:
  - "**/*.ts"
  - "**/*.tsx"
---
# TypeScript Patterns

Use strict mode in tsconfig.json.
Prefer readonly arrays and objects.
Use discriminated unions for state machines.
`,
		// React skills
		".github/skills/react-patterns/SKILL.md": `---
name: react-patterns
description: React component patterns, hooks, and state management
tags: [react, hooks, components, state]
files:
  - "**/*.tsx"
  - "**/*.jsx"
---
# React Patterns

Use custom hooks to encapsulate stateful logic.
Prefer controlled components over uncontrolled.
Use React.memo and useMemo judiciously.
`,
		// Go skills
		".github/skills/go-idioms/SKILL.md": `---
name: go-idioms
description: Idiomatic Go patterns, interfaces, and concurrency
tags: [go, golang, concurrency, interfaces]
files:
  - "**/*.go"
---
# Go Idioms

Interfaces should be small and focused.
Use channels for communication, mutexes for state.
Table-driven tests are the standard.
`,
		// Python skills
		".github/skills/python-style/SKILL.md": `---
name: python-style
description: Python style, type hints, and Pythonic patterns
tags: [python, pep8, type-hints]
files:
  - "**/*.py"
---
# Python Style

Follow PEP 8 and PEP 257.
Use type hints on all public functions.
Context managers for resource management.
`,
		// SQL skills
		".github/skills/sql-patterns/SKILL.md": `---
name: sql-patterns
description: SQL query patterns, indexing strategies, and migration best practices
tags: [sql, database, indexing, migrations]
files:
  - "**/*.sql"
---
# SQL Patterns

Use CTEs for complex queries.
Every FK should have an index.
Migration scripts must be idempotent.
`,
		// Docker skills
		".github/skills/docker-patterns/SKILL.md": `---
name: docker-patterns
description: Dockerfile best practices, multi-stage builds, and security
tags: [docker, containers, security, multi-stage]
files:
  - "Dockerfile"
  - "Dockerfile.*"
  - "docker-compose*.yml"
---
# Docker Patterns

Use multi-stage builds to minimize image size.
Run as non-root user in final stage.
Pin base image digests for reproducibility.
`,
		// CI/CD skills
		".github/skills/ci-patterns/SKILL.md": `---
name: ci-patterns
description: CI/CD pipeline patterns and GitHub Actions best practices
tags: [ci, cd, github-actions, pipelines]
files:
  - ".github/workflows/**/*.yml"
  - ".github/workflows/**/*.yaml"
---
# CI Patterns

Pin all action versions to a specific SHA.
Cache dependencies to speed up builds.
Use matrix builds for cross-platform testing.
`,
		// Testing skills
		".github/skills/testing-patterns/SKILL.md": `---
name: testing-patterns
description: Testing strategies, test pyramid, and TDD practices
tags: [testing, tdd, unit-tests, integration-tests]
---
# Testing Patterns

Follow the test pyramid: many unit, few integration, fewer E2E.
Tests should be FIRST: Fast, Independent, Repeatable, Self-validating, Timely.
`,
		// Deeply hidden skill
		".hidden/internal/skills/rest-api/SKILL.md": `---
name: rest-api
description: RESTful API design principles and HTTP conventions
tags: [api, rest, http, design]
---
# REST API Design

Use nouns for resource URLs, not verbs.
Proper HTTP status codes for all responses.
API versioning via URL path (/v1/).
`,
	})
}

// ---------------------------------------------------------------------------
// Helper
// ---------------------------------------------------------------------------

func loadAllSkills(t *testing.T, dir string) []Skill {
	t.Helper()
	cfg := &config.Config{
		Sources: config.Sources{Dirs: []string{dir}},
		Cache:   config.CacheConfig{Dir: t.TempDir(), SyncInterval: time.Hour},
	}
	s := newScanner(cfg)
	return s.List()
}

// ---------------------------------------------------------------------------
// Table-driven agent file-context scenarios
// ---------------------------------------------------------------------------

func TestAgentScenario_Skills_FileContextRouting(t *testing.T) {
	dir := skillsWorkspace(t)

	tests := []struct {
		name        string
		filePath    string
		mustHave    []string
		mustNotHave []string
	}{
		{
			name:        "TypeScript source file — ts-specific + global",
			filePath:    "src/auth/login.ts",
			mustHave:    []string{"typescript-patterns", "error-handling", "security-patterns"},
			mustNotHave: []string{"react-patterns", "go-idioms", "python-style", "sql-patterns", "docker-patterns"},
		},
		{
			name:        "React TSX component — ts + react + global",
			filePath:    "src/components/UserCard.tsx",
			mustHave:    []string{"typescript-patterns", "react-patterns", "error-handling"},
			mustNotHave: []string{"go-idioms", "python-style", "sql-patterns", "docker-patterns"},
		},
		{
			name:        "Go source file — go-specific + global",
			filePath:    "pkg/server/handler.go",
			mustHave:    []string{"go-idioms", "error-handling", "security-patterns"},
			mustNotHave: []string{"typescript-patterns", "react-patterns", "python-style"},
		},
		{
			name:        "Python script — python-specific + global",
			filePath:    "scripts/migrate.py",
			mustHave:    []string{"python-style", "error-handling"},
			mustNotHave: []string{"typescript-patterns", "go-idioms", "sql-patterns"},
		},
		{
			name:        "SQL file — sql-specific + global",
			filePath:    "migrations/0001_init.sql",
			mustHave:    []string{"sql-patterns", "error-handling"},
			mustNotHave: []string{"typescript-patterns", "go-idioms", "docker-patterns"},
		},
		{
			name:        "Dockerfile — docker-specific + global",
			filePath:    "Dockerfile",
			mustHave:    []string{"docker-patterns", "security-patterns"},
			mustNotHave: []string{"typescript-patterns", "go-idioms", "sql-patterns"},
		},
		{
			name:        "GitHub Actions workflow — ci/cd + global",
			filePath:    ".github/workflows/ci.yml",
			mustHave:    []string{"ci-patterns", "security-patterns"},
			mustNotHave: []string{"typescript-patterns", "go-idioms", "sql-patterns"},
		},
		{
			name:        "No file path — all skills returned",
			filePath:    "",
			mustHave:    []string{"git-workflow", "error-handling", "typescript-patterns", "go-idioms", "rest-api"},
			mustNotHave: []string{},
		},
	}

	all := loadAllSkills(t, dir)
	if len(all) == 0 {
		t.Fatal("no skills loaded from polyglot workspace")
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			var got []Skill
			if tc.filePath == "" {
				got = all
			} else {
				got = FilterByFilePath(all, tc.filePath)
			}

			names := make(map[string]bool)
			for _, s := range got {
				names[s.Name] = true
			}

			for _, want := range tc.mustHave {
				if !names[want] {
					t.Errorf("missing %q for file %q (got %v)", want, tc.filePath, skillNameKeys(names))
				}
			}
			for _, notWant := range tc.mustNotHave {
				if names[notWant] {
					t.Errorf("unexpected %q for file %q", notWant, tc.filePath)
				}
			}
		})
	}
}

func skillNameKeys(m map[string]bool) []string {
	out := make([]string, 0, len(m))
	for k := range m {
		out = append(out, k)
	}
	return out
}

// ---------------------------------------------------------------------------
// Query relevance — scoring ranks the most relevant skill first
// ---------------------------------------------------------------------------

func TestAgentScenario_Skills_QueryScoring(t *testing.T) {
	dir := skillsWorkspace(t)
	all := loadAllSkills(t, dir)

	tests := []struct {
		query     string
		wantFirst string
	}{
		{"git branching and commit conventions", "git-workflow"},
		{"error handling best practices", "error-handling"},
		{"security authentication validation", "security-patterns"},
		{"TypeScript generics type safety", "typescript-patterns"},
		{"React hooks custom state management", "react-patterns"},
		{"Go idiomatic concurrency interfaces", "go-idioms"},
		{"Docker multi-stage build", "docker-patterns"},
		{"SQL indexing query optimization", "sql-patterns"},
		{"REST API design HTTP conventions", "rest-api"},
	}

	for _, tc := range tests {
		t.Run(tc.query, func(t *testing.T) {
			ranked := FilterByQuery(all, tc.query)
			if len(ranked) == 0 {
				t.Fatal("no results returned")
			}
			if ranked[0].Name != tc.wantFirst {
				t.Errorf("query %q: want first=%q, got first=%q (full: %v)",
					tc.query, tc.wantFirst, ranked[0].Name, func() []string {
						out := make([]string, len(ranked))
						for i, s := range ranked {
							out[i] = s.Name
						}
						return out
					}())
			}
		})
	}
}

// ---------------------------------------------------------------------------
// Combined filter: file path + query
// ---------------------------------------------------------------------------

func TestAgentScenario_Skills_CombinedFilter(t *testing.T) {
	dir := skillsWorkspace(t)
	all := loadAllSkills(t, dir)

	// Go file + query about concurrency → go-idioms first
	goSkills := FilterByFilePath(all, "internal/worker/pool.go")
	ranked := FilterByQuery(goSkills, "concurrency channels goroutines")

	if len(ranked) == 0 {
		t.Fatal("no results after combined filter")
	}
	if ranked[0].Name != "go-idioms" {
		t.Errorf("combined filter Go: want go-idioms first, got %q", ranked[0].Name)
	}

	// error-handling (global) must still appear for Go files
	found := false
	for _, s := range goSkills {
		if s.Name == "error-handling" {
			found = true
			break
		}
	}
	if !found {
		t.Error("error-handling (global skill) missing from Go file results")
	}
}

// ---------------------------------------------------------------------------
// Hidden directories — skills in .hidden paths are discoverable
// ---------------------------------------------------------------------------

func TestAgentScenario_Skills_HiddenDirectories(t *testing.T) {
	dir := skillsWorkspace(t)
	all := loadAllSkills(t, dir)

	found := false
	for _, s := range all {
		if s.Name == "rest-api" {
			found = true
			if !strings.Contains(s.Content, "REST API Design") {
				t.Errorf("rest-api content corrupted: %q", s.Content[:50])
			}
			break
		}
	}
	if !found {
		t.Error("rest-api skill not found — hidden dir scan failed")
	}
}

// ---------------------------------------------------------------------------
// Tags as string vs list — both formats must parse correctly
// ---------------------------------------------------------------------------

func TestAgentScenario_Skills_TagFormats(t *testing.T) {
	dir := createTestDir(t, map[string]string{
		".github/skills/string-tags/SKILL.md": `---
name: string-tags
description: Tags as comma-separated string
tags: "typescript, go, python"
---
# String Tags Skill
Content.
`,
		".github/skills/list-tags/SKILL.md": `---
name: list-tags
description: Tags as YAML list
tags:
  - typescript
  - go
  - python
---
# List Tags Skill
Content.
`,
	})

	all := loadAllSkills(t, dir)
	byName := make(map[string]Skill)
	for _, s := range all {
		byName[s.Name] = s
	}

	for _, name := range []string{"string-tags", "list-tags"} {
		s, ok := byName[name]
		if !ok {
			t.Fatalf("skill %q not found", name)
		}
		if len(s.Tags) != 3 {
			t.Errorf("%q: want 3 tags, got %d: %v", name, len(s.Tags), s.Tags)
		}
		tagSet := make(map[string]bool)
		for _, tag := range s.Tags {
			tagSet[strings.TrimSpace(tag)] = true
		}
		for _, expected := range []string{"typescript", "go", "python"} {
			if !tagSet[expected] {
				t.Errorf("%q: missing tag %q, got %v", name, expected, s.Tags)
			}
		}
	}
}

// ---------------------------------------------------------------------------
// Brace expansion in files: patterns
// ---------------------------------------------------------------------------

func TestAgentScenario_Skills_BraceExpansionGlob(t *testing.T) {
	dir := createTestDir(t, map[string]string{
		".github/skills/ts-skill/SKILL.md": `---
name: ts-skill
description: TypeScript and React skill
files: "**/*.{ts,tsx}"
---
# TS Skill
`,
		".github/skills/go-skill/SKILL.md": `---
name: go-skill
description: Go skill
files: "**/*.go"
---
# Go Skill
`,
	})

	all := loadAllSkills(t, dir)

	tsResults := FilterByFilePath(all, "src/components/App.tsx")
	tsNames := make(map[string]bool)
	for _, s := range tsResults {
		tsNames[s.Name] = true
	}
	if !tsNames["ts-skill"] {
		t.Error("ts-skill not matched for .tsx file — brace expansion broken")
	}
	if tsNames["go-skill"] {
		t.Error("go-skill should not match .tsx file")
	}

	goResults := FilterByFilePath(all, "pkg/api/server.go")
	goNames := make(map[string]bool)
	for _, s := range goResults {
		goNames[s.Name] = true
	}
	if !goNames["go-skill"] {
		t.Error("go-skill not matched for .go file")
	}
	if goNames["ts-skill"] {
		t.Error("ts-skill should not match .go file")
	}
}

// ---------------------------------------------------------------------------
// Multiple directories — skills from all sources merged
// ---------------------------------------------------------------------------

func TestAgentScenario_Skills_MultipleDirectories(t *testing.T) {
	dir1 := createTestDir(t, map[string]string{
		".github/skills/skill-a/SKILL.md": "---\nname: skill-a\ndescription: from dir1\ntags: [dir1]\n---\n# Skill A\n",
	})
	dir2 := createTestDir(t, map[string]string{
		".github/skills/skill-b/SKILL.md": "---\nname: skill-b\ndescription: from dir2\ntags: [dir2]\n---\n# Skill B\n",
	})

	cfg := newTestConfig([]string{dir1, dir2}, nil, t.TempDir())
	s := newScanner(cfg)
	all := s.List()

	names := make(map[string]bool)
	for _, skill := range all {
		names[skill.Name] = true
	}
	if !names["skill-a"] {
		t.Error("missing skill-a from dir1")
	}
	if !names["skill-b"] {
		t.Error("missing skill-b from dir2")
	}
}

// ---------------------------------------------------------------------------
// Content integrity — no truncation, URI stable
// ---------------------------------------------------------------------------

func TestAgentScenario_Skills_ContentIntegrity(t *testing.T) {
	marker := "UNIQUE_SKILL_MARKER_ABC_67890"
	dir := createTestDir(t, map[string]string{
		".github/skills/special/SKILL.md": "---\nname: special-skill\ndescription: special\ntags: [test]\n---\n# Special\n" + marker + "\n",
	})

	all := loadAllSkills(t, dir)
	if len(all) == 0 {
		t.Fatal("no skills loaded")
	}

	found := false
	for _, s := range all {
		if s.Name == "special-skill" {
			found = true
			if !strings.Contains(s.Content, marker) {
				t.Errorf("content marker missing: %q", s.Content)
			}
			if s.URI == "" {
				t.Error("URI must not be empty")
			}
			break
		}
	}
	if !found {
		t.Error("special-skill not found")
	}
}

// ---------------------------------------------------------------------------
// Live reload — modifications on disk reflected after ForceSync
// ---------------------------------------------------------------------------

func TestAgentScenario_Skills_LiveReload(t *testing.T) {
	dir := t.TempDir()
	skillDir := dir + "/.github/skills/evolving"
	if err := os.MkdirAll(skillDir, 0o755); err != nil {
		t.Fatal(err)
	}
	skillPath := skillDir + "/SKILL.md"
	if err := os.WriteFile(skillPath, []byte("---\nname: evolving\ndescription: v1\ntags: [v1]\n---\n# V1\nVersion 1 content\n"), 0o644); err != nil {
		t.Fatal(err)
	}

	cfg := &config.Config{
		Sources: config.Sources{Dirs: []string{dir}},
		Cache:   config.CacheConfig{Dir: t.TempDir(), SyncInterval: time.Hour},
	}
	s := newScanner(cfg)
	all1 := s.List()
	if len(all1) == 0 {
		t.Fatal("no skills loaded initially")
	}
	if !strings.Contains(all1[0].Content, "Version 1") {
		t.Errorf("initial: want Version 1, got %q", all1[0].Content)
	}

	if err := os.WriteFile(skillPath, []byte("---\nname: evolving\ndescription: v2\ntags: [v2]\n---\n# V2\nVersion 2 content\n"), 0o644); err != nil {
		t.Fatal(err)
	}

	s.ForceSync()
	all2 := s.List()
	if len(all2) == 0 {
		t.Fatal("no skills after update")
	}
	if !strings.Contains(all2[0].Content, "Version 2") {
		t.Errorf("updated: want Version 2, got %q", all2[0].Content)
	}
}

// ---------------------------------------------------------------------------
// Determinism — List() order is stable across multiple calls
// ---------------------------------------------------------------------------

func TestAgentScenario_Skills_Determinism(t *testing.T) {
	dir := skillsWorkspace(t)
	cfg := &config.Config{
		Sources: config.Sources{Dirs: []string{dir}},
		Cache:   config.CacheConfig{Dir: t.TempDir(), SyncInterval: time.Hour},
	}
	s := newScanner(cfg)

	all1 := s.List()
	s.ForceSync()
	all2 := s.List()

	if len(all1) != len(all2) {
		t.Fatalf("non-deterministic count: %d vs %d", len(all1), len(all2))
	}
	for i := range all1 {
		if all1[i].URI != all2[i].URI {
			t.Errorf("position %d: URI changed (%q → %q)", i, all1[i].URI, all2[i].URI)
		}
	}
}

// ---------------------------------------------------------------------------
// Adversarial: malformed SKILL.md files — must not crash
// ---------------------------------------------------------------------------

func TestAgentScenario_Skills_AdversarialFrontmatter(t *testing.T) {
	dir := createTestDir(t, map[string]string{
		".github/skills/good/SKILL.md": "---\nname: good-skill\ndescription: Valid skill\ntags: [valid]\n---\n# Good Skill\nContent.\n",
		".github/skills/no-fm/SKILL.md": "No frontmatter at all.\nRaw content only.\n",
		".github/skills/empty-fm/SKILL.md": "---\n---\nBody without fields.\n",
		".github/skills/malformed/SKILL.md": "---\nname: bad\n  broken: indentation\ntags: [unclosed\n---\nContent.\n",
		".github/skills/wrong-type/SKILL.md": "---\nname: wrong-type\ndescription: files is integer\nfiles: 42\n---\n# Wrong Type\n",
		".github/skills/invalid-glob/SKILL.md": "---\nname: invalid-glob\ndescription: Bad glob\nfiles:\n  - \"[invalid[glob\"\n---\n# Invalid Glob\n",
	})

	all := loadAllSkills(t, dir)

	found := false
	for _, s := range all {
		if s.Name == "good-skill" {
			found = true
			break
		}
	}
	if !found {
		t.Error("good-skill missing — valid skill dropped due to bad neighbors")
	}

	// invalid-glob (if loaded) should appear in file-filtered results via inclusive fallback
	filtered := FilterByFilePath(all, "src/anything.ts")
	_ = filtered // just ensure it doesn't panic
}

// force github import usage
var _ = (*github.Client)(nil)
