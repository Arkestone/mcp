package loader

import (
	"os"
	"path/filepath"
	"sort"
	"testing"
	"time"

	"github.com/Arkestone/mcp/pkg/config"
	"github.com/Arkestone/mcp/pkg/github"
)

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

func newLoader(cfg *config.Config) *Loader {
	return New(cfg, &github.Client{})
}

// ---------------------------------------------------------------------------
// parseFrontmatter tests
// ---------------------------------------------------------------------------

func TestParseFrontmatter_NoFrontmatter(t *testing.T) {
	desc, mode, body := parseFrontmatter("just content here")
	if desc != "" || mode != "" || body != "just content here" {
		t.Errorf("got desc=%q mode=%q body=%q", desc, mode, body)
	}
}

func TestParseFrontmatter_WithDescriptionAndMode(t *testing.T) {
	content := "---\ndescription: Create component\nmode: agent\n---\nbody text"
	desc, mode, body := parseFrontmatter(content)
	if desc != "Create component" {
		t.Errorf("desc = %q, want %q", desc, "Create component")
	}
	if mode != "agent" {
		t.Errorf("mode = %q, want %q", mode, "agent")
	}
	if body != "body text" {
		t.Errorf("body = %q, want %q", body, "body text")
	}
}

func TestParseFrontmatter_DescriptionOnly(t *testing.T) {
	content := "---\ndescription: Code reviewer\n---\nreviewer body"
	desc, mode, body := parseFrontmatter(content)
	if desc != "Code reviewer" {
		t.Errorf("desc = %q, want %q", desc, "Code reviewer")
	}
	if mode != "" {
		t.Errorf("mode = %q, want empty", mode)
	}
	if body != "reviewer body" {
		t.Errorf("body = %q, want %q", body, "reviewer body")
	}
}

func TestParseFrontmatter_EmptyFields(t *testing.T) {
	content := "---\n---\nbody"
	desc, mode, body := parseFrontmatter(content)
	if desc != "" || mode != "" {
		t.Errorf("expected empty desc/mode, got desc=%q mode=%q", desc, mode)
	}
	if body != "body" {
		t.Errorf("body = %q, want %q", body, "body")
	}
}

func TestParseFrontmatter_UnclosedDelimiter(t *testing.T) {
	content := "---\ndescription: test\nbody without closing"
	desc, mode, body := parseFrontmatter(content)
	if desc != "" || mode != "" {
		t.Errorf("expected empty desc/mode for unclosed frontmatter")
	}
	if body != content {
		t.Errorf("body should equal original content")
	}
}

func TestParseFrontmatter_QuotedValues(t *testing.T) {
	content := "---\ndescription: \"Quoted description\"\nmode: 'ask'\n---\nbody"
	desc, mode, _ := parseFrontmatter(content)
	if desc != "Quoted description" {
		t.Errorf("desc = %q, want %q", desc, "Quoted description")
	}
	if mode != "ask" {
		t.Errorf("mode = %q, want %q", mode, "ask")
	}
}

// ---------------------------------------------------------------------------
// scanDir tests
// ---------------------------------------------------------------------------

func TestScanDir_PromptFiles(t *testing.T) {
	dir := createTestDir(t, map[string]string{
		".github/prompts/component.prompt.md": "---\ndescription: Create component\nmode: agent\n---\nCreate a component.",
	})
	prompts := scanDir(dir, "myrepo")

	if len(prompts) != 1 {
		t.Fatalf("got %d prompts, want 1", len(prompts))
	}
	p := prompts[0]
	if p.Name != "component" {
		t.Errorf("Name = %q, want %q", p.Name, "component")
	}
	if p.Description != "Create component" {
		t.Errorf("Description = %q", p.Description)
	}
	if p.Mode != "agent" {
		t.Errorf("Mode = %q", p.Mode)
	}
	if p.Type != TypePrompt {
		t.Errorf("Type = %q, want %q", p.Type, TypePrompt)
	}
	if p.Source != "myrepo" {
		t.Errorf("Source = %q", p.Source)
	}
	if p.Path != ".github/prompts/component.prompt.md" {
		t.Errorf("Path = %q", p.Path)
	}
	if p.URI != "prompts://myrepo/component" {
		t.Errorf("URI = %q", p.URI)
	}
	if p.Content != "Create a component." {
		t.Errorf("Content = %q", p.Content)
	}
}

func TestScanDir_ChatmodeFiles(t *testing.T) {
	dir := createTestDir(t, map[string]string{
		".github/chatmodes/reviewer.chatmode.md": "---\ndescription: Code reviewer\n---\nReview code carefully.",
	})
	prompts := scanDir(dir, "myrepo")

	if len(prompts) != 1 {
		t.Fatalf("got %d prompts, want 1", len(prompts))
	}
	p := prompts[0]
	if p.Name != "reviewer" {
		t.Errorf("Name = %q, want %q", p.Name, "reviewer")
	}
	if p.Description != "Code reviewer" {
		t.Errorf("Description = %q", p.Description)
	}
	if p.Mode != "" {
		t.Errorf("Mode should be empty for chatmode, got %q", p.Mode)
	}
	if p.Type != TypeChatmode {
		t.Errorf("Type = %q, want %q", p.Type, TypeChatmode)
	}
	if p.URI != "prompts://myrepo/reviewer" {
		t.Errorf("URI = %q", p.URI)
	}
}

func TestScanDir_BothTypes(t *testing.T) {
	dir := createTestDir(t, map[string]string{
		".github/prompts/component.prompt.md":    "---\ndescription: Create component\nmode: agent\n---\nbody",
		".github/chatmodes/reviewer.chatmode.md": "---\ndescription: Code reviewer\n---\nbody",
	})
	prompts := scanDir(dir, "myrepo")
	if len(prompts) != 2 {
		t.Fatalf("got %d prompts, want 2", len(prompts))
	}
	sort.Slice(prompts, func(i, j int) bool { return prompts[i].Type < prompts[j].Type })
	// chatmode < prompt alphabetically
	if prompts[0].Type != TypeChatmode {
		t.Errorf("expected first to be chatmode")
	}
	if prompts[1].Type != TypePrompt {
		t.Errorf("expected second to be prompt")
	}
}

func TestScanDir_IgnoresNonMatchingFiles(t *testing.T) {
	dir := createTestDir(t, map[string]string{
		".github/prompts/component.prompt.md":    "content",
		".github/prompts/README.md":              "not a prompt",
		".github/prompts/notes.txt":              "not a prompt",
		".github/chatmodes/reviewer.chatmode.md": "content",
		".github/chatmodes/README.md":            "not a chatmode",
	})
	prompts := scanDir(dir, "myrepo")
	if len(prompts) != 2 {
		t.Fatalf("got %d prompts, want 2 (only .prompt.md and .chatmode.md)", len(prompts))
	}
}

func TestScanDir_NoFrontmatterUsesFilename(t *testing.T) {
	dir := createTestDir(t, map[string]string{
		".github/prompts/myfile.prompt.md": "just content",
	})
	prompts := scanDir(dir, "myrepo")
	if len(prompts) != 1 {
		t.Fatalf("got %d prompts, want 1", len(prompts))
	}
	if prompts[0].Description != "myfile" {
		t.Errorf("Description = %q, want %q (filename fallback)", prompts[0].Description, "myfile")
	}
	if prompts[0].Content != "just content" {
		t.Errorf("Content = %q", prompts[0].Content)
	}
}

func TestScanDir_Empty(t *testing.T) {
	dir := t.TempDir()
	prompts := scanDir(dir, "empty")
	if len(prompts) != 0 {
		t.Errorf("got %d prompts from empty dir", len(prompts))
	}
}

func TestScanDir_Nonexistent(t *testing.T) {
	prompts := scanDir("/nonexistent/path/xyz", "gone")
	if len(prompts) != 0 {
		t.Errorf("got %d prompts from nonexistent dir", len(prompts))
	}
}

// ---------------------------------------------------------------------------
// Loader.List tests
// ---------------------------------------------------------------------------

func TestLoaderList_LocalDirs(t *testing.T) {
	dir1 := createTestDir(t, map[string]string{
		".github/prompts/comp.prompt.md": "---\ndescription: Component\n---\nbody",
	})
	dir2 := createTestDir(t, map[string]string{
		".github/prompts/test.prompt.md":       "---\ndescription: Test\n---\nbody",
		".github/chatmodes/review.chatmode.md": "---\ndescription: Reviewer\n---\nbody",
	})

	cfg := &config.Config{
		Sources: config.Sources{Dirs: []string{dir1, dir2}},
		Cache:   config.CacheConfig{Dir: t.TempDir(), SyncInterval: time.Hour},
	}
	ldr := newLoader(cfg)
	prompts := ldr.List()
	if len(prompts) != 3 {
		t.Fatalf("got %d prompts, want 3", len(prompts))
	}
}

func TestLoaderList_NoSources(t *testing.T) {
	cfg := &config.Config{
		Cache: config.CacheConfig{Dir: t.TempDir(), SyncInterval: time.Hour},
	}
	ldr := newLoader(cfg)
	if got := ldr.List(); len(got) != 0 {
		t.Errorf("got %d prompts from empty config", len(got))
	}
}

// ---------------------------------------------------------------------------
// Loader.Get tests
// ---------------------------------------------------------------------------

func TestLoaderGet_ByURI(t *testing.T) {
	dir := createTestDir(t, map[string]string{
		".github/prompts/component.prompt.md":    "---\ndescription: Create component\nmode: agent\n---\nContent here.",
		".github/chatmodes/reviewer.chatmode.md": "---\ndescription: Code reviewer\n---\nReview body.",
	})

	cfg := &config.Config{
		Sources: config.Sources{Dirs: []string{dir}},
		Cache:   config.CacheConfig{Dir: t.TempDir(), SyncInterval: time.Hour},
	}
	ldr := newLoader(cfg)

	src := filepath.Base(dir)

	t.Run("found prompt", func(t *testing.T) {
		p, ok := ldr.Get("prompts://" + src + "/component")
		if !ok {
			t.Fatal("expected to find prompt")
		}
		if p.Name != "component" {
			t.Errorf("Name = %q", p.Name)
		}
		if p.Mode != "agent" {
			t.Errorf("Mode = %q", p.Mode)
		}
		if p.Type != TypePrompt {
			t.Errorf("Type = %q", p.Type)
		}
	})

	t.Run("found chatmode", func(t *testing.T) {
		p, ok := ldr.Get("prompts://" + src + "/reviewer")
		if !ok {
			t.Fatal("expected to find chatmode")
		}
		if p.Type != TypeChatmode {
			t.Errorf("Type = %q", p.Type)
		}
	})

	t.Run("not found", func(t *testing.T) {
		_, ok := ldr.Get("prompts://" + src + "/nonexistent")
		if ok {
			t.Error("expected not found")
		}
	})

	t.Run("missing URI returns false", func(t *testing.T) {
		_, ok := ldr.Get("")
		if ok {
			t.Error("expected false for empty URI")
		}
	})
}
