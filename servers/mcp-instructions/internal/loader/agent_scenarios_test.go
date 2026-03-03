package loader

// agent_scenarios_test.go — realistic agent workflow tests.
//
// Each test simulates an AI agent loading instructions for a specific file-edit
// context, verifying that only the relevant subset is returned and in the
// expected order. The workspace mimics a real polyglot monorepo.

import (
	"os"
	"path/filepath"
	"sort"
	"strings"
	"testing"
	"time"

	"github.com/Arkestone/mcp/pkg/config"
)

// ---------------------------------------------------------------------------
// Polyglot repo fixture
// ---------------------------------------------------------------------------

// polyglotWorkspace builds a directory tree that looks like a real monorepo
// with TypeScript frontend, Go backend, Python scripts, SQL migrations, and CI.
func polyglotWorkspace(t *testing.T) string {
	t.Helper()
	return createTestDir(t, map[string]string{
		// ── global ──────────────────────────────────────────────────────────
		".github/copilot-instructions.md": `---
---
# Project Standards

- Conventional commits: feat|fix|docs|chore|refactor|test(scope): message
- Every PR needs tests
- Never hardcode credentials
- Use dependency injection
`,
		// ── language-specific ───────────────────────────────────────────────
		".github/instructions/typescript.instructions.md": `---
applyTo: "**/*.{ts,tsx,mts,cts}"
---
# TypeScript Standards

- Enable strict mode: noImplicitAny, strictNullChecks
- Prefer interface over type for object shapes
- Explicit return types on public functions
- Use readonly for immutable properties
- Avoid any; use unknown and narrow
`,
		".github/instructions/react.instructions.md": `---
applyTo: "**/*.{tsx,jsx}"
---
# React Standards

- Functional components only (no class components)
- Custom hooks for reusable stateful logic
- Memo/useMemo/useCallback only when profiled
- Co-locate styles, tests, and stories with components
`,
		".github/instructions/go.instructions.md": `---
applyTo: "**/*.go"
---
# Go Standards

- Handle all errors; never use _ for error returns
- Use table-driven tests with t.Run sub-tests
- Prefer composition over inheritance (embed interfaces)
- Use context.Context as first parameter for cancelable ops
- No init() functions; use explicit initialization
`,
		".github/instructions/python.instructions.md": `---
applyTo: "**/*.py"
---
# Python Standards

- Type annotations required on all public functions
- Use dataclasses or pydantic models, not plain dicts
- Prefer f-strings over .format() or %
- Virtual environment: poetry or uv
- Black + ruff for formatting and linting
`,
		// ── domain-specific ─────────────────────────────────────────────────
		".github/instructions/testing.instructions.md": `---
applyTo: "**/*_test.go,**/*.test.ts,**/*.spec.ts,**/*.test.tsx,**/*.spec.tsx"
---
# Testing Standards

- Arrange-Act-Assert (AAA) pattern
- One logical assertion per test (multiple asserts OK if logically linked)
- No sleep() in tests; use channels or polling helpers
- Mock at the boundary; prefer fake implementations over mocks
- Test names: TestUnit_Input_ExpectedOutput
`,
		".github/instructions/database.instructions.md": `---
applyTo: "**/*.sql,**/migrations/**,**/db/**"
---
# Database Standards

- Every table needs created_at, updated_at columns
- Always add indexes for foreign keys and frequently queried columns
- Use transactions for multi-step mutations
- Migration filenames: NNN_descriptive_name.sql (zero-padded)
- Never DROP without a corresponding rollback script
`,
		".github/instructions/cicd.instructions.md": `---
applyTo: ".github/workflows/**,Dockerfile*,docker-compose*"
---
# CI/CD Standards

- Pin action versions to a SHA, not a tag
- Cache node_modules / Go module cache in CI
- Docker images: use distroless or scratch for production
- Never run CI as root
- Fail fast: lint → unit → integration → build
`,
		// ── security (applies to everything) ────────────────────────────────
		".github/instructions/security.instructions.md": `---
applyTo: "**"
---
# Security Standards

- No hardcoded secrets, tokens, or API keys
- Validate and sanitize all user inputs
- Use parameterized queries; never interpolate SQL
- Log security events (failed auth, privilege escalation)
- Dependency audit on every merge
`,
		// ── hidden / non-standard locations ─────────────────────────────────
		".copilot/instructions/performance.instructions.md": `---
applyTo: "**"
---
# Performance Standards

- Avoid N+1 queries; use eager loading or data loaders
- Profile before optimizing; measure twice, cut once
- Cache expensive computations; document cache invalidation
`,
		// Instruction in unexpected deeply-nested hidden dir
		".vscode/.copilot/extra.instructions.md": `---
applyTo: "**/*.json"
---
# JSON File Standards

- Use 2-space indentation
- Sort keys alphabetically in config files
- No trailing commas (JSON strict)
`,
	})
}

// shortNames extracts just the file-stem names from instruction URIs
// (e.g. "instructions://workspace/typescript" → "typescript").
func shortNames(entries []Instruction) []string {
	names := make([]string, 0, len(entries))
	for _, e := range entries {
		parts := strings.Split(e.URI, "/")
		names = append(names, parts[len(parts)-1])
	}
	sort.Strings(names)
	return names
}

func loadAll(t *testing.T, dir string) []Instruction {
	t.Helper()
	cfg := &config.Config{Sources: config.Sources{Dirs: []string{dir}}, Cache: config.CacheConfig{Dir: t.TempDir(), SyncInterval: time.Hour}}
	l := newLoader(cfg)
	inst := l.List()
	return inst
}

// ---------------------------------------------------------------------------
// Table-driven agent file-context scenarios
// ---------------------------------------------------------------------------

func TestAgentScenario_Instructions_FileContextRouting(t *testing.T) {
	dir := polyglotWorkspace(t)
	all := loadAll(t, dir)

	// Verify we loaded all expected instructions.
	allNames := shortNames(all)
	wantAll := []string{
		"copilot-instructions",
		"cicd",
		"database",
		"extra",
		"go",
		"performance",
		"python",
		"react",
		"security",
		"testing",
		"typescript",
	}
	if len(allNames) != len(wantAll) {
		t.Fatalf("loaded %d instructions, want %d: got %v, want %v", len(allNames), len(wantAll), allNames, wantAll)
	}

	tests := []struct {
		name        string
		filePath    string
		wantIn      []string // must be present
		wantOut     []string // must be absent
		wantAtLeast int
	}{
		{
			name:     "TypeScript source file",
			filePath: "src/auth/jwt.ts",
			wantIn:   []string{"copilot-instructions", "typescript", "security", "performance"},
			wantOut:  []string{"react", "go", "python", "database", "cicd", "extra"},
		},
		{
			name:     "React component",
			filePath: "src/components/Button.tsx",
			wantIn:   []string{"copilot-instructions", "typescript", "react", "security", "performance"},
			wantOut:  []string{"go", "python", "database", "cicd", "extra"},
		},
		{
			name:     "Go source file",
			filePath: "internal/api/user_handler.go",
			wantIn:   []string{"copilot-instructions", "go", "security", "performance"},
			wantOut:  []string{"typescript", "react", "python", "database", "cicd", "extra"},
		},
		{
			name:     "Go test file",
			filePath: "internal/api/user_handler_test.go",
			wantIn:   []string{"copilot-instructions", "go", "testing", "security", "performance"},
			wantOut:  []string{"typescript", "react", "python", "database", "cicd", "extra"},
		},
		{
			name:     "Python script",
			filePath: "scripts/data_import.py",
			wantIn:   []string{"copilot-instructions", "python", "security", "performance"},
			wantOut:  []string{"typescript", "react", "go", "database", "cicd", "extra"},
		},
		{
			name:     "SQL migration file",
			filePath: "db/migrations/003_add_user_indexes.sql",
			wantIn:   []string{"copilot-instructions", "database", "security", "performance"},
			wantOut:  []string{"typescript", "react", "go", "python", "cicd", "extra"},
		},
		{
			name:     "File inside db/ directory",
			filePath: "db/queries/list_users.sql",
			wantIn:   []string{"database", "security"},
			wantOut:  []string{"typescript", "react", "go", "python", "cicd"},
		},
		{
			name:     "GitHub Actions workflow",
			filePath: ".github/workflows/ci.yml",
			wantIn:   []string{"copilot-instructions", "cicd", "security", "performance"},
			wantOut:  []string{"typescript", "react", "go", "python", "database", "extra"},
		},
		{
			name:     "Dockerfile",
			filePath: "Dockerfile",
			wantIn:   []string{"copilot-instructions", "cicd", "security", "performance"},
			wantOut:  []string{"typescript", "react", "go", "python", "database", "extra"},
		},
		{
			name:     "Docker Compose file",
			filePath: "docker-compose.yml",
			wantIn:   []string{"copilot-instructions", "cicd", "security", "performance"},
			wantOut:  []string{"typescript", "react", "go", "python", "database"},
		},
		{
			name:     "JSON config file",
			filePath: "config/settings.json",
			wantIn:   []string{"copilot-instructions", "security", "performance", "extra"},
			wantOut:  []string{"typescript", "react", "go", "python", "database", "cicd"},
		},
		{
			name:     "Markdown docs file",
			filePath: "docs/api-reference.md",
			wantIn:   []string{"copilot-instructions", "security", "performance"},
			wantOut:  []string{"typescript", "react", "go", "python", "database", "cicd", "extra"},
		},
		{
			name:        "No file path — all instructions returned",
			filePath:    "",
			wantAtLeast: len(wantAll),
		},
		{
			name:     "TypeScript test file — gets both ts and testing",
			filePath: "src/__tests__/auth.test.ts",
			wantIn:   []string{"copilot-instructions", "typescript", "testing", "security", "performance"},
			wantOut:  []string{"react", "go", "python", "database", "cicd", "extra"},
		},
		{
			name:     "React component test",
			filePath: "src/components/Button.test.tsx",
			wantIn:   []string{"copilot-instructions", "typescript", "react", "testing", "security", "performance"},
			wantOut:  []string{"go", "python", "database", "cicd", "extra"},
		},
		{
			name:    "Deeply nested Go file",
			filePath: "pkg/internal/auth/middleware/jwt/validator.go",
			wantIn:  []string{"go", "security"},
			wantOut: []string{"typescript", "react", "python", "database", "cicd"},
		},
		{
			name:     "Python test file",
			filePath: "tests/test_migration.py",
			// Python matches *.py; testing only matches *_test.go/*.test.ts/*.spec.ts
			wantIn:  []string{"python", "security", "performance"},
			wantOut: []string{"typescript", "react", "go", "database", "cicd"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := FilterByFilePath(all, tt.filePath)
			names := shortNames(got)

			if tt.wantAtLeast > 0 {
				if len(got) < tt.wantAtLeast {
					t.Errorf("got %d instructions, want at least %d", len(got), tt.wantAtLeast)
				}
				return
			}

			nameSet := make(map[string]bool, len(names))
			for _, n := range names {
				nameSet[n] = true
			}

			for _, want := range tt.wantIn {
				if !nameSet[want] {
					t.Errorf("file=%q: want %q in result, got %v", tt.filePath, want, names)
				}
			}
			for _, notWant := range tt.wantOut {
				if nameSet[notWant] {
					t.Errorf("file=%q: want %q NOT in result, got %v", tt.filePath, notWant, names)
				}
			}
		})
	}
}

// ---------------------------------------------------------------------------
// Non-standard file structures
// ---------------------------------------------------------------------------

func TestAgentScenario_Instructions_NonStandardLayouts(t *testing.T) {
	t.Run("root level copilot-instructions only", func(t *testing.T) {
		dir := createTestDir(t, map[string]string{
			"copilot-instructions.md": "# Root level global instruction\n",
		})
		all := loadAll(t, dir)
		if len(all) != 1 {
			t.Errorf("expected 1 instruction, got %d", len(all))
		}
		// Applies globally
		filtered := FilterByFilePath(all, "any/file.ts")
		if len(filtered) != 1 {
			t.Errorf("global root instruction should match any file, got %d", len(filtered))
		}
	})

	t.Run("comma-separated applyTo pattern", func(t *testing.T) {
		dir := createTestDir(t, map[string]string{
			".github/instructions/multi.instructions.md": `---
applyTo: "**/*.ts,**/*.tsx,**/*.mts"
---
# Multi-extension TypeScript
`,
		})
		all := loadAll(t, dir)

		tsMatch := FilterByFilePath(all, "src/app.ts")
		if len(tsMatch) != 1 {
			t.Errorf(".ts: want 1, got %d", len(tsMatch))
		}
		tsxMatch := FilterByFilePath(all, "src/App.tsx")
		if len(tsxMatch) != 1 {
			t.Errorf(".tsx: want 1, got %d", len(tsxMatch))
		}
		mtsMatch := FilterByFilePath(all, "src/worker.mts")
		if len(mtsMatch) != 1 {
			t.Errorf(".mts: want 1, got %d", len(mtsMatch))
		}
		goNoMatch := FilterByFilePath(all, "main.go")
		if len(goNoMatch) != 0 {
			t.Errorf(".go: want 0, got %d", len(goNoMatch))
		}
	})

	t.Run("brace expansion applyTo", func(t *testing.T) {
		dir := createTestDir(t, map[string]string{
			".github/instructions/frontend.instructions.md": `---
applyTo: "src/**/*.{ts,tsx,js,jsx,css,scss}"
---
# Frontend standards
`,
		})
		all := loadAll(t, dir)

		cases := []struct{ path string; want bool }{
			{"src/components/Button.tsx", true},
			{"src/utils/helpers.ts", true},
			{"src/styles/main.css", true},
			{"src/styles/theme.scss", true},
			{"src/app.js", true},
			{"src/index.jsx", true},
			{"internal/server.go", false},
			{"tests/utils.py", false},
			// Outside src/ — should NOT match
			{"components/Button.tsx", false},
		}

		for _, c := range cases {
			filtered := FilterByFilePath(all, c.path)
			got := len(filtered) == 1
			if got != c.want {
				t.Errorf("path=%q: want match=%v, got match=%v", c.path, c.want, got)
			}
		}
	})

	t.Run("glob with path separator in pattern", func(t *testing.T) {
		dir := createTestDir(t, map[string]string{
			".github/instructions/api.instructions.md": `---
applyTo: "internal/api/**,pkg/handlers/**"
---
# API handler standards
`,
		})
		all := loadAll(t, dir)

		cases := []struct{ path string; want bool }{
			{"internal/api/user.go", true},
			{"internal/api/v2/auth.go", true},
			{"pkg/handlers/http.go", true},
			{"pkg/handlers/grpc/server.go", true},
			{"cmd/main.go", false},
			{"internal/service/user.go", false},
		}

		for _, c := range cases {
			filtered := FilterByFilePath(all, c.path)
			got := len(filtered) == 1
			if got != c.want {
				t.Errorf("path=%q: want match=%v, got match=%v", c.path, c.want, got)
			}
		}
	})

	t.Run("multiple directories — sources merged", func(t *testing.T) {
		dir1 := createTestDir(t, map[string]string{
			".github/instructions/ts.instructions.md": `---
applyTo: "**/*.ts"
---
# TypeScript from dir1
`,
		})
		dir2 := createTestDir(t, map[string]string{
			".github/instructions/go.instructions.md": `---
applyTo: "**/*.go"
---
# Go from dir2
`,
		})

		cfg := &config.Config{Sources: config.Sources{Dirs: []string{dir1, dir2}}, Cache: config.CacheConfig{Dir: t.TempDir(), SyncInterval: time.Hour}}
		l := newLoader(cfg)
		all := l.List()
		if len(all) != 2 {
			t.Errorf("want 2 instructions from 2 dirs, got %d", len(all))
		}

		tsOnly := FilterByFilePath(all, "auth.ts")
		if len(tsOnly) != 1 {
			t.Errorf("ts: want 1, got %d", len(tsOnly))
		}
		goOnly := FilterByFilePath(all, "main.go")
		if len(goOnly) != 1 {
			t.Errorf("go: want 1, got %d", len(goOnly))
		}
	})
}

// ---------------------------------------------------------------------------
// Full workspace content tests — verify content integrity
// ---------------------------------------------------------------------------

func TestAgentScenario_Instructions_ContentIntegrity(t *testing.T) {
	dir := createTestDir(t, map[string]string{
		".github/instructions/go.instructions.md": `---
applyTo: "**/*.go"
---
# Go Standards

Handle all errors explicitly.
Use context.Context as first parameter.
`,
	})

	cfg := &config.Config{Sources: config.Sources{Dirs: []string{dir}}, Cache: config.CacheConfig{Dir: t.TempDir(), SyncInterval: time.Hour}}
	l := newLoader(cfg)
	all := l.List()
	if len(all) != 1 {
		t.Fatalf("want 1, got %d", len(all))
	}

	inst := all[0]
	if !strings.Contains(inst.Content, "Handle all errors") {
		t.Errorf("content missing expected text, got: %q", inst.Content[:min(100, len(inst.Content))])
	}
	if !strings.Contains(inst.Content, "context.Context") {
		t.Errorf("content missing context.Context, got: %q", inst.Content[:min(100, len(inst.Content))])
	}
	if inst.ApplyTo == nil || len(inst.ApplyTo) == 0 {
		t.Error("ApplyTo should be populated")
	}
	if inst.ApplyTo[0] != "**/*.go" {
		t.Errorf("ApplyTo[0] = %q, want %q", inst.ApplyTo[0], "**/*.go")
	}
}

// ---------------------------------------------------------------------------
// Agent query (FilterByQuery) scenarios
// ---------------------------------------------------------------------------

func TestAgentScenario_Instructions_QueryFiltering(t *testing.T) {
	dir := polyglotWorkspace(t)
	all := loadAll(t, dir)

	t.Run("query narrows list", func(t *testing.T) {
		// Instructions don't use query scoring normally (they use file path).
		// But calling FilterByFilePath with an empty query should return everything.
		filtered := FilterByFilePath(all, "")
		if len(filtered) != len(all) {
			t.Errorf("empty file_path: want all %d, got %d", len(all), len(filtered))
		}
	})

	t.Run("combined: query context + file path", func(t *testing.T) {
		// For a Go test file, we want go + testing instructions
		filtered := FilterByFilePath(all, "internal/api/handler_test.go")
		names := shortNames(filtered)
		nameSet := make(map[string]bool)
		for _, n := range names {
			nameSet[n] = true
		}
		if !nameSet["go"] {
			t.Errorf("go instruction missing for _test.go: %v", names)
		}
		if !nameSet["testing"] {
			t.Errorf("testing instruction missing for _test.go: %v", names)
		}
		if nameSet["react"] || nameSet["python"] || nameSet["database"] {
			t.Errorf("unexpected instructions for Go test file: %v", names)
		}
	})
}

// ---------------------------------------------------------------------------
// Performance — large workspace doesn't degrade
// ---------------------------------------------------------------------------

func TestAgentScenario_Instructions_LargeWorkspace(t *testing.T) {
	files := map[string]string{
		".github/copilot-instructions.md": "---\n---\n# Global",
	}
	// Add 50 language-specific instructions
	extensions := []string{
		"ts", "tsx", "go", "py", "rb", "java", "kt", "rs",
		"cpp", "c", "cs", "php", "swift", "scala", "clj",
	}
	for i, ext := range extensions {
		files[filepath.Join(".github/instructions", ext+".instructions.md")] =
			"---\napplyTo: \"**/*." + ext + "\"\n---\n# " + ext + " standards " + string(rune('A'+i))
	}

	dir := createTestDir(t, files)
	all := loadAll(t, dir)

	if len(all) != len(files) {
		t.Errorf("want %d instructions, got %d", len(files), len(all))
	}

	// For a .ts file, should get: global + ts + tsx (tsx contains .ts files?  no) = global + ts
	filtered := FilterByFilePath(all, "src/auth.ts")
	names := shortNames(filtered)
	found := false
	for _, n := range names {
		if n == "ts" {
			found = true
		}
	}
	if !found {
		t.Errorf("ts instruction not found for .ts file: %v", names)
	}
	// Should NOT include java, go, py, etc.
	for _, n := range names {
		for _, ext := range []string{"go", "py", "java", "rb", "rs", "cpp"} {
			if n == ext {
				t.Errorf("unexpected %q instruction for .ts file: %v", n, names)
			}
		}
	}
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

// ---------------------------------------------------------------------------
// Edge: instructions loaded from hidden directories (.copilot, .vscode)
// ---------------------------------------------------------------------------

func TestAgentScenario_Instructions_HiddenDirectories(t *testing.T) {
	dir := createTestDir(t, map[string]string{
		// Standard location
		".github/instructions/standard.instructions.md": `---
applyTo: "**/*.go"
---
Standard Go instruction
`,
		// Non-standard hidden locations
		".copilot/instructions/hidden.instructions.md": `---
applyTo: "**"
---
Hidden copilot global instruction
`,
		".vscode/.copilot/nested.instructions.md": `---
applyTo: "**/*.json"
---
Hidden VSCode nested instruction
`,
		// Deep nesting
		".tools/.config/.copilot/deep.instructions.md": `---
applyTo: "**/*.yml"
---
Deep nested YAML instruction
`,
	})

	all := loadAll(t, dir)
	names := shortNames(all)
	nameSet := make(map[string]bool)
	for _, n := range names {
		nameSet[n] = true
	}

	// All should be discovered (recursive scan)
	if !nameSet["standard"] {
		t.Errorf("standard instruction not found: %v", names)
	}
	if !nameSet["hidden"] {
		t.Errorf("hidden .copilot instruction not found: %v", names)
	}
	if !nameSet["nested"] {
		t.Errorf("nested .vscode instruction not found: %v", names)
	}
	if !nameSet["deep"] {
		t.Errorf("deep nested instruction not found: %v", names)
	}

	// Verify content routing still works
	goFiltered := FilterByFilePath(all, "main.go")
	goNames := shortNames(goFiltered)
	goSet := make(map[string]bool)
	for _, n := range goNames {
		goSet[n] = true
	}
	if !goSet["standard"] {
		t.Errorf("standard (applyTo:**/*.go) should match main.go: %v", goNames)
	}
	if !goSet["hidden"] {
		t.Errorf("hidden (applyTo:**) should match main.go: %v", goNames)
	}
	if goSet["nested"] {
		t.Errorf("nested (applyTo:**/*.json) should NOT match main.go: %v", goNames)
	}
	if goSet["deep"] {
		t.Errorf("deep (applyTo:**/*.yml) should NOT match main.go: %v", goNames)
	}
}

// ---------------------------------------------------------------------------
// Determinism — multiple calls return same order
// ---------------------------------------------------------------------------

func TestAgentScenario_Instructions_Determinism(t *testing.T) {
	dir := polyglotWorkspace(t)
	cfg := &config.Config{Sources: config.Sources{Dirs: []string{dir}}, Cache: config.CacheConfig{Dir: t.TempDir(), SyncInterval: time.Hour}}
	l := newLoader(cfg)

	var results [][]string
	for i := 0; i < 5; i++ {
		all := l.List()
		names := shortNames(all)
		results = append(results, names)
	}

	for i := 1; i < len(results); i++ {
		if strings.Join(results[0], ",") != strings.Join(results[i], ",") {
			t.Errorf("non-deterministic order: run 0 = %v, run %d = %v", results[0], i, results[i])
		}
	}
}

// ---------------------------------------------------------------------------
// File with only frontmatter separator, no content
// ---------------------------------------------------------------------------

func TestAgentScenario_Instructions_VariousFrontmatterStyles(t *testing.T) {
	dir := createTestDir(t, map[string]string{
		// Style 1: YAML array for applyTo
		".github/instructions/array-style.instructions.md": `---
applyTo:
  - "**/*.ts"
  - "**/*.tsx"
---
Array-style applyTo instruction
`,
		// Style 2: Single string
		".github/instructions/string-style.instructions.md": `---
applyTo: "**/*.go"
---
String-style applyTo instruction
`,
		// Style 3: Comma-separated string
		".github/instructions/comma-style.instructions.md": `---
applyTo: "**/*.py,**/*.pyi"
---
Comma-separated applyTo instruction
`,
		// Style 4: Brace expansion
		".github/instructions/brace-style.instructions.md": `---
applyTo: "**/*.{rs,toml}"
---
Brace expansion applyTo instruction
`,
		// Style 5: Wildcard only
		".github/instructions/global-style.instructions.md": `---
applyTo: "**"
---
Wildcard applyTo instruction
`,
		// Style 6: No frontmatter (global)
		".github/instructions/no-fm-style.instructions.md": `
# Just content, no frontmatter

This instruction has no YAML block.
`,
	})

	all := loadAll(t, dir)
	if len(all) != 6 {
		t.Errorf("want 6 instructions, got %d: %v", len(all), shortNames(all))
	}

	type testCase struct {
		file      string
		wantNames []string
		wantCount int
	}

	cases := []testCase{
		{
			file:      "src/auth.ts",
			wantNames: []string{"array-style", "global-style", "no-fm-style"},
			wantCount: 3,
		},
		{
			file:      "main.go",
			wantNames: []string{"string-style", "global-style", "no-fm-style"},
			wantCount: 3,
		},
		{
			file:      "scripts/migrate.py",
			wantNames: []string{"comma-style", "global-style", "no-fm-style"},
			wantCount: 3,
		},
		{
			file:      "src/lib.rs",
			wantNames: []string{"brace-style", "global-style", "no-fm-style"},
			wantCount: 3,
		},
		{
			file:      "Cargo.toml",
			wantNames: []string{"brace-style", "global-style", "no-fm-style"},
			wantCount: 3,
		},
		{
			file:      "any/random/file.xyz",
			wantNames: []string{"global-style", "no-fm-style"},
			wantCount: 2,
		},
	}

	for _, c := range cases {
		t.Run("file="+c.file, func(t *testing.T) {
			filtered := FilterByFilePath(all, c.file)
			if len(filtered) != c.wantCount {
				t.Errorf("file=%q: want %d instructions, got %d: %v",
					c.file, c.wantCount, len(filtered), shortNames(filtered))
			}
			names := shortNames(filtered)
			nameSet := make(map[string]bool)
			for _, n := range names {
				nameSet[n] = true
			}
			for _, want := range c.wantNames {
				if !nameSet[want] {
					t.Errorf("file=%q: want %q present, got %v", c.file, want, names)
				}
			}
		})
	}
}

// ---------------------------------------------------------------------------
// Validate ApplyTo field is properly exposed
// ---------------------------------------------------------------------------

func TestAgentScenario_Instructions_ApplyToExposed(t *testing.T) {
	dir := createTestDir(t, map[string]string{
		".github/instructions/typed.instructions.md": `---
applyTo: "**/*.go,**/*.ts"
---
Multi-lang instruction
`,
	})
	all := loadAll(t, dir)
	if len(all) != 1 {
		t.Fatalf("want 1, got %d", len(all))
	}
	applyTo := all[0].ApplyTo
	if len(applyTo) != 2 {
		t.Errorf("ApplyTo len = %d, want 2: %v", len(applyTo), applyTo)
	}
	wantPatterns := map[string]bool{"**/*.go": true, "**/*.ts": true}
	for _, p := range applyTo {
		if !wantPatterns[p] {
			t.Errorf("unexpected ApplyTo pattern: %q", p)
		}
	}
}

// ---------------------------------------------------------------------------
// Loader.Get returns instruction by URI with full content
// ---------------------------------------------------------------------------

func TestAgentScenario_Instructions_GetByURI(t *testing.T) {
	dir := createTestDir(t, map[string]string{
		".github/instructions/go.instructions.md": `---
applyTo: "**/*.go"
---
# Go Standards

Use idiomatic error handling.
`,
	})
	cfg := &config.Config{Sources: config.Sources{Dirs: []string{dir}}, Cache: config.CacheConfig{Dir: t.TempDir(), SyncInterval: time.Hour}}
	l := newLoader(cfg)

	all := l.List()
	if len(all) == 0 {
		t.Fatal("no instructions loaded")
	}

	uri := all[0].URI
	inst, ok := l.Get(uri)
	if !ok {
		t.Fatalf("Get(%q) not found", uri)
	}
	if !strings.Contains(inst.Content, "idiomatic error handling") {
		t.Errorf("Get returned wrong content: %q", inst.Content)
	}
}

// ---------------------------------------------------------------------------
// Instruction file modified on disk — live reload
// ---------------------------------------------------------------------------

func TestAgentScenario_Instructions_LiveReload(t *testing.T) {
	dir := t.TempDir()
	instDir := filepath.Join(dir, ".github", "instructions")
	if err := os.MkdirAll(instDir, 0o755); err != nil {
		t.Fatal(err)
	}

	path := filepath.Join(instDir, "go.instructions.md")
	if err := os.WriteFile(path, []byte("---\napplyTo: \"**/*.go\"\n---\nVersion 1\n"), 0o644); err != nil {
		t.Fatal(err)
	}

	cfg := &config.Config{Sources: config.Sources{Dirs: []string{dir}}, Cache: config.CacheConfig{Dir: t.TempDir(), SyncInterval: time.Hour}}
	l := newLoader(cfg)

	all1 := l.List()
	if len(all1) == 0 {
		t.Fatal("no instructions")
	}
	if !strings.Contains(all1[0].Content, "Version 1") {
		t.Errorf("initial content: want Version 1, got %q", all1[0].Content)
	}

	// Update the file
	if err := os.WriteFile(path, []byte("---\napplyTo: \"**/*.go\"\n---\nVersion 2 updated\n"), 0o644); err != nil {
		t.Fatal(err)
	}

	// Invalidate cache before second read to pick up the updated file.
	l.ForceSync()

	all2 := l.List()
	if len(all2) == 0 {
		t.Fatal("no instructions after update")
	}
	if !strings.Contains(all2[0].Content, "Version 2 updated") {
		t.Errorf("updated content: want Version 2, got %q", all2[0].Content)
	}
}
