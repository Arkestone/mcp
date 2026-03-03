package loader

import (
	"context"
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"sort"
	"strings"
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
	desc, mode, _, _, body := parseFrontmatter("just content here")
	if desc != "" || mode != "" || body != "just content here" {
		t.Errorf("got desc=%q mode=%q body=%q", desc, mode, body)
	}
}

func TestParseFrontmatter_WithDescriptionAndMode(t *testing.T) {
	content := "---\ndescription: Create component\nmode: agent\n---\nbody text"
	desc, mode, _, _, body := parseFrontmatter(content)
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
	desc, mode, _, _, body := parseFrontmatter(content)
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
	desc, mode, _, _, body := parseFrontmatter(content)
	if desc != "" || mode != "" {
		t.Errorf("expected empty desc/mode, got desc=%q mode=%q", desc, mode)
	}
	if body != "body" {
		t.Errorf("body = %q, want %q", body, "body")
	}
}

func TestParseFrontmatter_UnclosedDelimiter(t *testing.T) {
	content := "---\ndescription: test\nbody without closing"
	desc, mode, _, _, body := parseFrontmatter(content)
	if desc != "" || mode != "" {
		t.Errorf("expected empty desc/mode for unclosed frontmatter")
	}
	if body != content {
		t.Errorf("body should equal original content")
	}
}

// TestParseFrontmatter_NoTrailingNewline tests content where body is truly empty
// (no trailing newline after the closing "---").  This kills the
// CONDITIONALS_BOUNDARY mutation len(body)>0 → len(body)>=0 which would panic
// on an empty-slice access when len==0.
func TestParseFrontmatter_NoTrailingNewline(t *testing.T) {
	content := "---\ndescription: nodesc\n---"
	desc, _, _, _, body := parseFrontmatter(content)
	if desc != "nodesc" {
		t.Errorf("desc = %q, want %q", desc, "nodesc")
	}
	if body != "" {
		t.Errorf("body = %q, want empty", body)
	}
}

func TestParseFrontmatter_QuotedValues(t *testing.T) {
	content := "---\ndescription: \"Quoted description\"\nmode: 'ask'\n---\nbody"
	desc, mode, _, _, _ := parseFrontmatter(content)
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

// ---------------------------------------------------------------------------
// Root-level scanning (not just .github/)
// ---------------------------------------------------------------------------

func TestScanDir_RootLevelPrompts(t *testing.T) {
	dir := createTestDir(t, map[string]string{
		"prompts/rootprompt.prompt.md": "---\ndescription: Root prompt\nmode: ask\n---\nbody",
	})
	prompts := scanDir(dir, "rootrepo")
	if len(prompts) != 1 {
		t.Fatalf("got %d prompts, want 1", len(prompts))
	}
	if prompts[0].Name != "rootprompt" {
		t.Errorf("Name = %q", prompts[0].Name)
	}
	if prompts[0].Type != TypePrompt {
		t.Errorf("Type = %q", prompts[0].Type)
	}
	if prompts[0].URI != "prompts://rootrepo/rootprompt" {
		t.Errorf("URI = %q", prompts[0].URI)
	}
}

func TestScanDir_RootLevelChatmodes(t *testing.T) {
	dir := createTestDir(t, map[string]string{
		"chatmodes/rootmode.chatmode.md": "---\ndescription: Root chatmode\n---\nbody",
	})
	prompts := scanDir(dir, "rootrepo")
	if len(prompts) != 1 {
		t.Fatalf("got %d prompts, want 1", len(prompts))
	}
	if prompts[0].Type != TypeChatmode {
		t.Errorf("Type = %q", prompts[0].Type)
	}
}

func TestScanDir_GithubTakesPriorityOverRoot(t *testing.T) {
	// When same prompt exists in .github/ and root, .github/ wins.
	dir := createTestDir(t, map[string]string{
		".github/prompts/foo.prompt.md": "---\ndescription: From github\n---\nbody",
		"prompts/foo.prompt.md":         "---\ndescription: From root\n---\nbody",
	})
	prompts := scanDir(dir, "repo")
	if len(prompts) != 1 {
		t.Fatalf("got %d prompts, want 1 (deduplication)", len(prompts))
	}
	if prompts[0].Description != "From github" {
		t.Errorf("expected .github/ version to win, got %q", prompts[0].Description)
	}
}

func TestScanDir_BothGithubAndRootDifferentFiles(t *testing.T) {
	dir := createTestDir(t, map[string]string{
		".github/prompts/a.prompt.md": "---\ndescription: A\n---\nbody",
		"prompts/b.prompt.md":         "---\ndescription: B\n---\nbody",
	})
	prompts := scanDir(dir, "repo")
	if len(prompts) != 2 {
		t.Fatalf("got %d prompts, want 2", len(prompts))
	}
}

// ---------------------------------------------------------------------------
// FilterByQuery tests
// ---------------------------------------------------------------------------

func TestFilterByQuery_EmptyQuery(t *testing.T) {
prompts := []Prompt{
{Name: "component", Description: "Create a React component"},
{Name: "test", Description: "Write unit tests"},
}
got := FilterByQuery(prompts, "")
if len(got) != len(prompts) {
t.Errorf("empty query should return all %d prompts, got %d", len(prompts), len(got))
}
}

func TestFilterByQuery_Match(t *testing.T) {
prompts := []Prompt{
{Name: "jwt-auth", Description: "JWT authentication setup", Tags: []string{"auth", "security"}},
{Name: "deploy", Description: "Deploy to production"},
{Name: "review", Description: "Code review checklist"},
}
got := FilterByQuery(prompts, "auth")
if len(got) != 1 || got[0].Name != "jwt-auth" {
t.Errorf("'auth' should match only jwt-auth, got %v", got)
}
}

func TestFilterByQuery_NoMatch(t *testing.T) {
prompts := []Prompt{
{Name: "deploy", Description: "Deploy to production"},
}
got := FilterByQuery(prompts, "authentication")
// "authentication" should match "auth" via stem but "deploy" has no auth tokens
if len(got) != 0 {
t.Errorf("'authentication' should not match deploy prompt, got %v", got)
}
}

func TestFilterByQuery_SortsByScore(t *testing.T) {
prompts := []Prompt{
// "build-tool" matches "build" only in name → score=10
{Name: "build-tool", Description: "A hammering utility"},
// "build" matches in name (10) + desc (3) → higher total score
{Name: "build", Description: "Build the project from source"},
}
got := FilterByQuery(prompts, "build")
if len(got) < 2 {
t.Fatalf("expected 2 matching prompts, got %d", len(got))
}
if got[0].Name != "build" {
t.Errorf("higher-scoring item should rank first, got %q", got[0].Name)
}
}

func TestFilterByQuery_StopwordsIgnored(t *testing.T) {
prompts := []Prompt{
{Name: "real-match", Description: "JWT authentication flow"},
{Name: "unrelated", Description: "Unrelated content"},
}
// "how to use" are all stopwords — should act as pass-through
got := FilterByQuery(prompts, "how to use")
if len(got) != 2 {
t.Errorf("all-stopword query should return all prompts, got %d", len(got))
}
}

// ---------------------------------------------------------------------------
// Loader lifecycle tests (Start/Stop/ForceSync)
// ---------------------------------------------------------------------------

func TestLoaderStartStop(t *testing.T) {
cfg := &config.Config{
Cache: config.CacheConfig{Dir: t.TempDir(), SyncInterval: time.Hour},
}
ldr := newLoader(cfg)
ctx, cancel := context.WithCancel(t.Context())
defer cancel()
ldr.Start(ctx)
ldr.Stop()
}

func TestLoaderForceSync(t *testing.T) {
dir := createTestDir(t, map[string]string{
".github/prompts/v1.prompt.md": "---\ndescription: Version 1\n---\nv1 content",
})
cfg := &config.Config{
Sources: config.Sources{Dirs: []string{dir}},
Cache:   config.CacheConfig{Dir: t.TempDir(), SyncInterval: time.Hour},
}
ldr := newLoader(cfg)

prompts := ldr.List()
if len(prompts) == 0 || prompts[0].Content != "v1 content" {
t.Fatalf("initial content wrong: %v", prompts)
}

// Update file and force cache invalidation.
if err := os.WriteFile(
filepath.Join(dir, ".github/prompts/v1.prompt.md"),
[]byte("---\ndescription: Version 2\n---\nv2 content"),
0o644,
); err != nil {
t.Fatal(err)
}

ldr.ForceSync()
prompts = ldr.List()
if len(prompts) == 0 || prompts[0].Content != "v2 content" {
t.Errorf("after ForceSync, content = %q, want v2 content", prompts[0].Content)
}
}

// ---------------------------------------------------------------------------
// repoCacheDir / sourceFor tests
// ---------------------------------------------------------------------------

func TestRepoCacheDir(t *testing.T) {
tests := []struct {
base, owner, repo, ref string
want                   string
}{
{"/tmp/cache", "owner", "repo", "", "owner_repo"},
{"/tmp/cache", "owner", "repo", "main", "owner_repo_main"},
}
for _, tt := range tests {
ref := config.RepoRef{Owner: tt.owner, Repo: tt.repo, Ref: tt.ref}
got := repoCacheDir(tt.base, ref)
if !strings.HasSuffix(got, tt.want) {
t.Errorf("repoCacheDir(%q,%q,%q) = %q, want suffix %q", tt.base, tt.owner, tt.repo, got, tt.want)
}
}
}

func TestSourceFor(t *testing.T) {
// Should return the basename of the resolved absolute path.
got := sourceFor(".")
if got == "" {
t.Error("sourceFor('.') should not be empty")
}
got2 := sourceFor("/tmp/my-dir")
if got2 != "my-dir" {
t.Errorf("sourceFor('/tmp/my-dir') = %q, want %q", got2, "my-dir")
}
}

// ---------------------------------------------------------------------------
// syncRepo / syncAllRepos tests
// ---------------------------------------------------------------------------

func TestSyncRepo_DownloadsPromptFiles(t *testing.T) {
srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
switch {
case strings.HasSuffix(r.URL.Path, "/contents/"):
w.Header().Set("Content-Type", "application/json")
fmt.Fprintf(w, `[{"type":"file","path":".github/prompts/hello.prompt.md","name":"hello.prompt.md","size":40}]`)
case strings.HasSuffix(r.URL.Path, "/contents/.github/prompts/hello.prompt.md"):
w.Header().Set("Content-Type", "application/vnd.github.raw+json")
fmt.Fprint(w, "---\ndescription: Hello prompt\n---\nHello world")
default:
http.NotFound(w, r)
}
}))
t.Cleanup(srv.Close)
gh := &github.Client{BaseURL: srv.URL}
cacheDir := t.TempDir()
cfg := &config.Config{
Cache:   config.CacheConfig{Dir: cacheDir, SyncInterval: time.Hour},
Sources: config.Sources{Repos: []string{"owner/repo"}},
}
ldr := New(cfg, gh)
ref := config.RepoRef{Owner: "owner", Repo: "repo"}
if err := ldr.syncRepo(ref); err != nil {
t.Fatalf("syncRepo: %v", err)
}
target := filepath.Join(repoCacheDir(cacheDir, ref), ".github/prompts/hello.prompt.md")
data, err := os.ReadFile(target)
if err != nil {
t.Fatalf("cached file not found: %v", err)
}
if !strings.Contains(string(data), "Hello world") {
t.Errorf("cached content wrong: %q", data)
}
}

func TestSyncRepo_APIError(t *testing.T) {
srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
http.Error(w, "not found", http.StatusNotFound)
}))
t.Cleanup(srv.Close)
gh := &github.Client{BaseURL: srv.URL}
cfg := &config.Config{
Cache: config.CacheConfig{Dir: t.TempDir(), SyncInterval: time.Hour},
}
ldr := New(cfg, gh)
ref := config.RepoRef{Owner: "owner", Repo: "repo"}
err := ldr.syncRepo(ref)
if err == nil {
t.Error("expected error when API returns 404")
}
}

func TestSyncAllRepos_NoRepos(t *testing.T) {
cfg := &config.Config{
Cache: config.CacheConfig{Dir: t.TempDir(), SyncInterval: time.Hour},
}
ldr := newLoader(cfg)
ldr.syncAllRepos() // should not panic
}

// ---------------------------------------------------------------------------
// FilterByFilePath tests
// ---------------------------------------------------------------------------

func TestFilterByFilePath_EmptyPath(t *testing.T) {
prompts := []Prompt{
{Name: "ts-prompt", Files: []string{"**/*.ts"}},
{Name: "global-prompt"},
}
got := FilterByFilePath(prompts, "")
if len(got) != 2 {
t.Errorf("empty filePath should return all prompts, got %d", len(got))
}
}

func TestFilterByFilePath_Matches(t *testing.T) {
prompts := []Prompt{
{Name: "ts-prompt", Files: []string{"**/*.ts"}},
{Name: "go-prompt", Files: []string{"**/*.go"}},
{Name: "global-prompt"},
}
got := FilterByFilePath(prompts, "src/auth.ts")
if len(got) != 2 {
t.Errorf("expected 2 prompts (ts + global), got %d", len(got))
}
for _, p := range got {
if p.Name == "go-prompt" {
t.Error("go-prompt should be excluded for .ts file")
}
}
}

func TestFilterByFilePath_NoPattern_AlwaysIncluded(t *testing.T) {
prompts := []Prompt{{Name: "global", Files: nil}}
got := FilterByFilePath(prompts, "anything.py")
if len(got) != 1 {
t.Errorf("prompt with no Files should always be included, got %d", len(got))
}
}

func TestFilterByFilePath_MultiplePatterns(t *testing.T) {
p := Prompt{Name: "web", Files: []string{"**/*.ts", "**/*.tsx", "**/*.js"}}
if got := FilterByFilePath([]Prompt{p}, "components/Button.tsx"); len(got) != 1 {
t.Error("Button.tsx should match **/*.tsx")
}
if got := FilterByFilePath([]Prompt{p}, "main.go"); len(got) != 0 {
t.Error("main.go should not match any web pattern")
}
}

func TestFilterByFilePath_ParsedFromFrontmatter(t *testing.T) {
dir := createTestDir(t, map[string]string{
".github/prompts/ts-review.prompt.md": "---\ndescription: TypeScript reviewer\nfiles: \"**/*.ts\"\n---\nbody",
".github/prompts/global.prompt.md":    "---\ndescription: Global prompt\n---\nbody",
})
cfg := &config.Config{
Sources: config.Sources{Dirs: []string{dir}},
Cache:   config.CacheConfig{Dir: t.TempDir(), SyncInterval: time.Hour},
}
ldr := newLoader(cfg)
all := ldr.List()
if len(all) != 2 {
t.Fatalf("expected 2 prompts, got %d", len(all))
}
filtered := FilterByFilePath(all, "src/auth.ts")
if len(filtered) != 2 {
t.Errorf("both prompts should match .ts file (ts-review matches, global always matches), got %d", len(filtered))
}
filtered2 := FilterByFilePath(all, "main.go")
if len(filtered2) != 1 {
t.Errorf("only global prompt should match .go file, got %d", len(filtered2))
}
}

func TestParseFrontmatter_TagsAsString(t *testing.T) {
content := "---\ndescription: test\ntags: \"typescript, go, python\"\n---\nbody"
desc, _, tags, _, body := parseFrontmatter(content)
if desc != "test" {
t.Errorf("description = %q", desc)
}
if len(tags) != 3 || tags[0] != "typescript" || tags[1] != "go" || tags[2] != "python" {
t.Errorf("tags = %v, want [typescript go python]", tags)
}
if body != "body" {
t.Errorf("body = %q", body)
}
}

func TestParseFrontmatter_TagsAsList(t *testing.T) {
content := "---\ndescription: test\ntags:\n  - ts\n  - go\n---\nbody"
_, _, tags, _, _ := parseFrontmatter(content)
if len(tags) != 2 || tags[0] != "ts" || tags[1] != "go" {
t.Errorf("tags from list = %v, want [ts go]", tags)
}
}
