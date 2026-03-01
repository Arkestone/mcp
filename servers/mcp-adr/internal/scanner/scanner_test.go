package scanner

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/Arkestone/mcp/pkg/config"
	"github.com/Arkestone/mcp/pkg/github"
)

// ---------------------------------------------------------------------------
// helpers
// ---------------------------------------------------------------------------

func createTestADR(t *testing.T, dir, subdir, filename, title, status, dateVal, body string) {
	t.Helper()
	adrDir := filepath.Join(dir, filepath.FromSlash(subdir))
	if err := os.MkdirAll(adrDir, 0o755); err != nil {
		t.Fatal(err)
	}
	var content string
	if title != "" || status != "" || dateVal != "" {
		content = "---\n"
		if title != "" {
			content += fmt.Sprintf("title: %s\n", title)
		}
		if status != "" {
			content += fmt.Sprintf("status: %s\n", status)
		}
		if dateVal != "" {
			content += fmt.Sprintf("date: %s\n", dateVal)
		}
		content += "---\n" + body
	} else {
		content = body
	}
	if err := os.WriteFile(filepath.Join(adrDir, filename), []byte(content), 0o644); err != nil {
		t.Fatal(err)
	}
}

func newTestConfig(dirs []string, repos []string, cacheDir string) *config.Config {
	return &config.Config{
		Sources: config.Sources{
			Dirs:  dirs,
			Repos: repos,
		},
		Cache: config.CacheConfig{
			Dir:          cacheDir,
			SyncInterval: time.Hour,
		},
	}
}

func newScanner(cfg *config.Config) *Scanner {
	return New(cfg, &github.Client{})
}

// ---------------------------------------------------------------------------
// parseFrontmatter tests
// ---------------------------------------------------------------------------

func TestParseFrontmatter_AllFields(t *testing.T) {
	content := "---\ntitle: Use PostgreSQL\nstatus: accepted\ndate: 2023-01-15\n---\nBody text"
	title, status, date := parseFrontmatter(content)
	if title != "Use PostgreSQL" {
		t.Errorf("title = %q, want %q", title, "Use PostgreSQL")
	}
	if status != "accepted" {
		t.Errorf("status = %q, want %q", status, "accepted")
	}
	if date != "2023-01-15" {
		t.Errorf("date = %q, want %q", date, "2023-01-15")
	}
}

func TestParseFrontmatter_PartialFields(t *testing.T) {
	content := "---\ntitle: Only Title\n---\nBody"
	title, status, date := parseFrontmatter(content)
	if title != "Only Title" {
		t.Errorf("title = %q, want %q", title, "Only Title")
	}
	if status != "" {
		t.Errorf("status = %q, want empty", status)
	}
	if date != "" {
		t.Errorf("date = %q, want empty", date)
	}
}

func TestParseFrontmatter_NoFrontmatter(t *testing.T) {
	content := "# ADR Title\n\nSome content"
	title, status, date := parseFrontmatter(content)
	if title != "" || status != "" || date != "" {
		t.Errorf("expected all empty, got title=%q status=%q date=%q", title, status, date)
	}
}

func TestParseFrontmatter_OnlyStatus(t *testing.T) {
	content := "---\nstatus: proposed\n---\nBody"
	title, status, date := parseFrontmatter(content)
	if title != "" {
		t.Errorf("title = %q, want empty", title)
	}
	if status != "proposed" {
		t.Errorf("status = %q, want proposed", status)
	}
	if date != "" {
		t.Errorf("date = %q, want empty", date)
	}
}

func TestParseFrontmatter_QuotedValues(t *testing.T) {
	content := "---\ntitle: \"Quoted Title\"\nstatus: 'accepted'\ndate: \"2024-06-01\"\n---\n"
	title, status, date := parseFrontmatter(content)
	if title != "Quoted Title" {
		t.Errorf("title = %q, want %q", title, "Quoted Title")
	}
	if status != "accepted" {
		t.Errorf("status = %q, want %q", status, "accepted")
	}
	if date != "2024-06-01" {
		t.Errorf("date = %q, want %q", date, "2024-06-01")
	}
}

func TestParseFrontmatter_NoClosingDashes(t *testing.T) {
	content := "---\ntitle: No Close\nstatus: accepted\n"
	title, status, date := parseFrontmatter(content)
	if title != "" || status != "" || date != "" {
		t.Errorf("expected all empty for unclosed frontmatter, got title=%q status=%q date=%q", title, status, date)
	}
}

// ---------------------------------------------------------------------------
// humanize tests
// ---------------------------------------------------------------------------

func TestHumanize_Basic(t *testing.T) {
	tests := []struct {
		input string
		want  string
	}{
		{"0001-use-postgresql", "0001 Use Postgresql"},
		{"0002-adopt-go", "0002 Adopt Go"},
		{"no-dashes", "no Dashes"},
		{"single", "single"},
		{"", ""},
	}
	for _, tc := range tests {
		got := humanize(tc.input)
		if got != tc.want {
			t.Errorf("humanize(%q) = %q, want %q", tc.input, got, tc.want)
		}
	}
}

// ---------------------------------------------------------------------------
// scanDir tests
// ---------------------------------------------------------------------------

func TestScanDir_BasicDocsADR(t *testing.T) {
	dir := t.TempDir()
	createTestADR(t, dir, "docs/adr", "0001-use-go.md", "Use Go", "accepted", "2023-01-01", "We use Go.")

	adrs := scanDir(dir, "myrepo")
	if len(adrs) != 1 {
		t.Fatalf("got %d ADRs, want 1", len(adrs))
	}
	a := adrs[0]
	if a.ID != "0001-use-go" {
		t.Errorf("ID = %q, want %q", a.ID, "0001-use-go")
	}
	if a.Title != "Use Go" {
		t.Errorf("Title = %q, want %q", a.Title, "Use Go")
	}
	if a.Status != "accepted" {
		t.Errorf("Status = %q, want %q", a.Status, "accepted")
	}
	if a.Date != "2023-01-01" {
		t.Errorf("Date = %q, want %q", a.Date, "2023-01-01")
	}
	if a.Source != "myrepo" {
		t.Errorf("Source = %q, want %q", a.Source, "myrepo")
	}
	if a.URI != "adrs://myrepo/0001-use-go" {
		t.Errorf("URI = %q, want %q", a.URI, "adrs://myrepo/0001-use-go")
	}
	if a.Path != "docs/adr/0001-use-go.md" {
		t.Errorf("Path = %q, want %q", a.Path, "docs/adr/0001-use-go.md")
	}
}

func TestScanDir_MultipleADRs(t *testing.T) {
	dir := t.TempDir()
	createTestADR(t, dir, "docs/adr", "0001-use-go.md", "Use Go", "accepted", "2023-01-01", "")
	createTestADR(t, dir, "docs/adr", "0002-use-postgresql.md", "Use PostgreSQL", "proposed", "2023-02-01", "")

	adrs := scanDir(dir, "src")
	if len(adrs) != 2 {
		t.Fatalf("got %d ADRs, want 2", len(adrs))
	}
	// Should be sorted by ID
	if adrs[0].ID != "0001-use-go" {
		t.Errorf("adrs[0].ID = %q, want %q", adrs[0].ID, "0001-use-go")
	}
	if adrs[1].ID != "0002-use-postgresql" {
		t.Errorf("adrs[1].ID = %q, want %q", adrs[1].ID, "0002-use-postgresql")
	}
}

func TestScanDir_DocsDecisions(t *testing.T) {
	dir := t.TempDir()
	createTestADR(t, dir, "docs/decisions", "0001-use-kubernetes.md", "Use Kubernetes", "accepted", "", "")

	adrs := scanDir(dir, "myproject")
	if len(adrs) != 1 {
		t.Fatalf("got %d ADRs, want 1", len(adrs))
	}
	if adrs[0].Path != "docs/decisions/0001-use-kubernetes.md" {
		t.Errorf("Path = %q, want docs/decisions/0001-use-kubernetes.md", adrs[0].Path)
	}
}

func TestScanDir_DocADR(t *testing.T) {
	dir := t.TempDir()
	createTestADR(t, dir, "doc/adr", "0001-use-redis.md", "Use Redis", "proposed", "", "")

	adrs := scanDir(dir, "svc")
	if len(adrs) != 1 {
		t.Fatalf("got %d ADRs, want 1", len(adrs))
	}
	if adrs[0].Path != "doc/adr/0001-use-redis.md" {
		t.Errorf("Path = %q, want doc/adr/0001-use-redis.md", adrs[0].Path)
	}
}

func TestScanDir_TitleFallbackToHumanize(t *testing.T) {
	dir := t.TempDir()
	createTestADR(t, dir, "docs/adr", "0001-use-postgresql.md", "", "", "", "No frontmatter body")

	adrs := scanDir(dir, "src")
	if len(adrs) != 1 {
		t.Fatalf("got %d ADRs, want 1", len(adrs))
	}
	if adrs[0].Title != "0001 Use Postgresql" {
		t.Errorf("Title = %q, want %q", adrs[0].Title, "0001 Use Postgresql")
	}
}

func TestScanDir_SkipsNonMD(t *testing.T) {
	dir := t.TempDir()
	createTestADR(t, dir, "docs/adr", "0001-use-go.md", "Use Go", "accepted", "", "")
	// Create a non-md file
	if err := os.WriteFile(filepath.Join(dir, "docs/adr", "README.txt"), []byte("readme"), 0o644); err != nil {
		t.Fatal(err)
	}

	adrs := scanDir(dir, "src")
	if len(adrs) != 1 {
		t.Fatalf("got %d ADRs, want 1 (should skip .txt)", len(adrs))
	}
}

func TestScanDir_SkipsSubdirs(t *testing.T) {
	dir := t.TempDir()
	createTestADR(t, dir, "docs/adr", "0001-use-go.md", "Use Go", "accepted", "", "")
	// Create a subdirectory inside docs/adr
	if err := os.MkdirAll(filepath.Join(dir, "docs/adr/subdir"), 0o755); err != nil {
		t.Fatal(err)
	}

	adrs := scanDir(dir, "src")
	if len(adrs) != 1 {
		t.Fatalf("got %d ADRs, want 1 (should skip subdirs)", len(adrs))
	}
}

func TestScanDir_EmptyDir(t *testing.T) {
	dir := t.TempDir()
	adrs := scanDir(dir, "src")
	if len(adrs) != 0 {
		t.Errorf("got %d ADRs, want 0", len(adrs))
	}
}

func TestScanDir_NonexistentDir(t *testing.T) {
	adrs := scanDir("/nonexistent/path/xyz", "src")
	if len(adrs) != 0 {
		t.Errorf("got %d ADRs, want 0", len(adrs))
	}
}

func TestScanDir_ContentIncludesFrontmatter(t *testing.T) {
	dir := t.TempDir()
	raw := "---\ntitle: Test ADR\nstatus: accepted\n---\nSome content"
	adrDir := filepath.Join(dir, "docs/adr")
	if err := os.MkdirAll(adrDir, 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(adrDir, "0001-test.md"), []byte(raw), 0o644); err != nil {
		t.Fatal(err)
	}

	adrs := scanDir(dir, "src")
	if len(adrs) != 1 {
		t.Fatalf("got %d ADRs, want 1", len(adrs))
	}
	if adrs[0].Content != raw {
		t.Errorf("Content = %q, want %q", adrs[0].Content, raw)
	}
}

// ---------------------------------------------------------------------------
// Scanner List / Get tests
// ---------------------------------------------------------------------------

func setupScannerDir(t *testing.T) (string, *Scanner) {
	t.Helper()
	dir := t.TempDir()
	createTestADR(t, dir, "docs/adr", "0001-use-go.md", "Use Go", "accepted", "2023-01-01", "We use Go.")
	createTestADR(t, dir, "docs/adr", "0002-use-postgresql.md", "Use PostgreSQL", "proposed", "2023-02-01", "We use PostgreSQL.")

	cfg := newTestConfig([]string{dir}, nil, t.TempDir())
	scn := newScanner(cfg)
	scn.Start(context.Background())
	t.Cleanup(scn.Stop)
	return dir, scn
}

func TestScanner_List_ReturnsSorted(t *testing.T) {
	_, scn := setupScannerDir(t)
	adrs := scn.List()
	if len(adrs) != 2 {
		t.Fatalf("got %d ADRs, want 2", len(adrs))
	}
	if adrs[0].ID != "0001-use-go" {
		t.Errorf("adrs[0].ID = %q, want 0001-use-go", adrs[0].ID)
	}
	if adrs[1].ID != "0002-use-postgresql" {
		t.Errorf("adrs[1].ID = %q, want 0002-use-postgresql", adrs[1].ID)
	}
}

func TestScanner_Get_Found(t *testing.T) {
	dir := t.TempDir()
	createTestADR(t, dir, "docs/adr", "0001-use-go.md", "Use Go", "accepted", "2023-01-01", "We use Go.")

	cfg := newTestConfig([]string{dir}, nil, t.TempDir())
	scn := newScanner(cfg)

	source := filepath.Base(dir)
	uri := fmt.Sprintf("adrs://%s/0001-use-go", source)
	a, ok := scn.Get(uri)
	if !ok {
		t.Fatalf("Get(%q) returned false", uri)
	}
	if a.ID != "0001-use-go" {
		t.Errorf("ID = %q, want 0001-use-go", a.ID)
	}
	if a.Title != "Use Go" {
		t.Errorf("Title = %q, want Use Go", a.Title)
	}
}

func TestScanner_Get_NotFound(t *testing.T) {
	dir := t.TempDir()
	cfg := newTestConfig([]string{dir}, nil, t.TempDir())
	scn := newScanner(cfg)

	_, ok := scn.Get("adrs://src/nonexistent")
	if ok {
		t.Error("expected Get to return false for nonexistent ADR")
	}
}

func TestScanner_List_EmptyConfig(t *testing.T) {
	cfg := newTestConfig(nil, nil, t.TempDir())
	scn := newScanner(cfg)
	adrs := scn.List()
	if len(adrs) != 0 {
		t.Errorf("got %d ADRs, want 0", len(adrs))
	}
}

func TestScanner_List_MultipleLocalDirs(t *testing.T) {
	dir1 := t.TempDir()
	dir2 := t.TempDir()
	createTestADR(t, dir1, "docs/adr", "0001-adr-one.md", "ADR One", "accepted", "", "")
	createTestADR(t, dir2, "docs/adr", "0001-adr-two.md", "ADR Two", "proposed", "", "")

	cfg := newTestConfig([]string{dir1, dir2}, nil, t.TempDir())
	scn := newScanner(cfg)
	adrs := scn.List()
	if len(adrs) != 2 {
		t.Fatalf("got %d ADRs, want 2", len(adrs))
	}
}

func TestScanner_StartStop(t *testing.T) {
	cfg := newTestConfig(nil, nil, t.TempDir())
	scn := newScanner(cfg)
	ctx := context.Background()
	scn.Start(ctx)
	scn.Stop() // should not panic
}

func TestScanner_ForceSync_NoRepos(t *testing.T) {
	cfg := newTestConfig(nil, nil, t.TempDir())
	scn := newScanner(cfg)
	scn.ForceSync() // should not panic with no repos
}

// ---------------------------------------------------------------------------
// repoCacheDir tests
// ---------------------------------------------------------------------------

func TestRepoCacheDir_WithRef(t *testing.T) {
	ref := config.RepoRef{Owner: "myorg", Repo: "myrepo", Ref: "main"}
	got := repoCacheDir("/cache", ref)
	want := "/cache/myorg_myrepo_main"
	if got != want {
		t.Errorf("repoCacheDir = %q, want %q", got, want)
	}
}

func TestRepoCacheDir_NoRef(t *testing.T) {
	ref := config.RepoRef{Owner: "myorg", Repo: "myrepo"}
	got := repoCacheDir("/cache", ref)
	want := "/cache/myorg_myrepo"
	if got != want {
		t.Errorf("repoCacheDir = %q, want %q", got, want)
	}
}

// ---------------------------------------------------------------------------
// HTTP sync tests
// ---------------------------------------------------------------------------

func newGitHubTestServer(t *testing.T, handler http.HandlerFunc) *httptest.Server {
	t.Helper()
	ts := httptest.NewServer(handler)
	t.Cleanup(ts.Close)
	return ts
}

func TestSyncRepo_Basic(t *testing.T) {
	fileContent := "---\ntitle: Use Go\nstatus: accepted\ndate: 2023-01-01\n---\nWe use Go."

	ts := newGitHubTestServer(t, func(w http.ResponseWriter, r *http.Request) {
		path := r.URL.Path
		switch {
		case strings.HasSuffix(path, "/contents/docs/adr"):
			json.NewEncoder(w).Encode([]github.ContentEntry{
				{Name: "0001-use-go.md", Path: "docs/adr/0001-use-go.md", Type: "file"},
			})
		case strings.HasSuffix(path, "/contents/docs/adr/0001-use-go.md"):
			w.Write([]byte(fileContent))
		default:
			http.NotFound(w, r)
		}
	})

	cacheDir := t.TempDir()
	cfg := newTestConfig(nil, []string{"myorg/myrepo"}, cacheDir)
	scn := New(cfg, &github.Client{BaseURL: ts.URL, Token: "test"})

	ref := config.ParseRepoRef("myorg/myrepo")
	if err := scn.syncRepo(ref); err != nil {
		t.Fatalf("syncRepo: %v", err)
	}

	adrPath := filepath.Join(cacheDir, "myorg_myrepo", "docs", "adr", "0001-use-go.md")
	content, err := os.ReadFile(adrPath)
	if err != nil {
		t.Fatalf("cached file not found: %v", err)
	}
	if string(content) != fileContent {
		t.Errorf("cached content = %q, want %q", string(content), fileContent)
	}
}

func TestSyncRepo_404(t *testing.T) {
	ts := newGitHubTestServer(t, func(w http.ResponseWriter, r *http.Request) {
		http.NotFound(w, r)
	})

	cfg := newTestConfig(nil, []string{"myorg/myrepo"}, t.TempDir())
	scn := New(cfg, &github.Client{BaseURL: ts.URL, Token: "test"})
	// Should not panic, just log
	scn.syncAllRepos()
}

func TestScanner_List_WithCachedRepo(t *testing.T) {
	cacheDir := t.TempDir()
	// Pre-populate cache
	adrDir := filepath.Join(cacheDir, "myorg_myrepo", "docs", "adr")
	if err := os.MkdirAll(adrDir, 0o755); err != nil {
		t.Fatal(err)
	}
	content := "---\ntitle: Use Go\nstatus: accepted\ndate: 2023-01-01\n---\nWe use Go."
	if err := os.WriteFile(filepath.Join(adrDir, "0001-use-go.md"), []byte(content), 0o644); err != nil {
		t.Fatal(err)
	}

	cfg := newTestConfig(nil, []string{"myorg/myrepo"}, cacheDir)
	scn := newScanner(cfg)
	adrs := scn.List()
	if len(adrs) != 1 {
		t.Fatalf("got %d ADRs, want 1", len(adrs))
	}
	if adrs[0].Source != "myorg/myrepo" {
		t.Errorf("Source = %q, want myorg/myrepo", adrs[0].Source)
	}
	if adrs[0].Title != "Use Go" {
		t.Errorf("Title = %q, want Use Go", adrs[0].Title)
	}
}
