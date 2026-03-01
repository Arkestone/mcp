package loader

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"runtime"
	"sort"
	"strings"
	"sync"
	"testing"
	"time"

	"github.com/Arkestone/mcp/pkg/config"
	"github.com/Arkestone/mcp/pkg/github"
)

// createTestDir creates a temporary directory with the given instruction files.
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

func newLoaderWithGH(cfg *config.Config, baseURL, token string) *Loader {
	return New(cfg, &github.Client{BaseURL: baseURL, Token: token})
}

func TestScanDirCopilotInstructions(t *testing.T) {
	dir := createTestDir(t, map[string]string{
		".github/copilot-instructions.md": "# Build\ngo build ./...",
	})
	instructions := scanDir(dir, "myrepo")

	if len(instructions) != 1 {
		t.Fatalf("got %d instructions, want 1", len(instructions))
	}
	inst := instructions[0]
	if inst.Source != "myrepo" {
		t.Errorf("Source = %q", inst.Source)
	}
	if inst.Path != ".github/copilot-instructions.md" {
		t.Errorf("Path = %q", inst.Path)
	}
	if inst.Content != "# Build\ngo build ./..." {
		t.Errorf("Content = %q", inst.Content)
	}
	if inst.URI != "instructions://myrepo/copilot-instructions" {
		t.Errorf("URI = %q", inst.URI)
	}
}

func TestScanDirPathSpecificInstructions(t *testing.T) {
	dir := createTestDir(t, map[string]string{
		".github/instructions/golang.instructions.md":  "Use gofmt",
		".github/instructions/testing.instructions.md": "Use table-driven tests",
	})
	instructions := scanDir(dir, "myrepo")

	if len(instructions) != 2 {
		t.Fatalf("got %d instructions, want 2", len(instructions))
	}

	sort.Slice(instructions, func(i, j int) bool {
		return instructions[i].URI < instructions[j].URI
	})

	if instructions[0].URI != "instructions://myrepo/golang" {
		t.Errorf("URI[0] = %q", instructions[0].URI)
	}
	if instructions[0].Content != "Use gofmt" {
		t.Errorf("Content[0] = %q", instructions[0].Content)
	}
	if instructions[1].URI != "instructions://myrepo/testing" {
		t.Errorf("URI[1] = %q", instructions[1].URI)
	}
}

func TestScanDirBothTypes(t *testing.T) {
	dir := createTestDir(t, map[string]string{
		".github/copilot-instructions.md":             "repo-wide",
		".github/instructions/python.instructions.md": "use black",
		".github/instructions/docker.instructions.md": "multi-stage builds",
	})
	instructions := scanDir(dir, "fullrepo")

	if len(instructions) != 3 {
		t.Fatalf("got %d instructions, want 3", len(instructions))
	}
}

func TestScanDirIgnoresNonInstructionFiles(t *testing.T) {
	dir := createTestDir(t, map[string]string{
		".github/instructions/golang.instructions.md": "Use gofmt",
		".github/instructions/README.md":              "Not an instruction",
		".github/instructions/notes.txt":              "Not an instruction",
	})
	instructions := scanDir(dir, "myrepo")

	if len(instructions) != 1 {
		t.Fatalf("got %d instructions, want 1 (only .instructions.md)", len(instructions))
	}
}

func TestScanDirEmpty(t *testing.T) {
	dir := t.TempDir()
	instructions := scanDir(dir, "empty")
	if len(instructions) != 0 {
		t.Errorf("got %d instructions from empty dir", len(instructions))
	}
}

func TestScanDirNonexistent(t *testing.T) {
	instructions := scanDir("/nonexistent/path", "gone")
	if len(instructions) != 0 {
		t.Errorf("got %d instructions from nonexistent dir", len(instructions))
	}
}

func TestScanDirNestedInstructions(t *testing.T) {
	dir := createTestDir(t, map[string]string{
		".github/instructions/sub/deep.instructions.md": "deeply nested",
	})
	instructions := scanDir(dir, "nested")

	if len(instructions) != 1 {
		t.Fatalf("got %d instructions, want 1", len(instructions))
	}
	if instructions[0].Content != "deeply nested" {
		t.Errorf("Content = %q", instructions[0].Content)
	}
}

func TestLoaderListLocalDirs(t *testing.T) {
	dir1 := createTestDir(t, map[string]string{
		".github/copilot-instructions.md": "repo1 instructions",
	})
	dir2 := createTestDir(t, map[string]string{
		".github/copilot-instructions.md":            "repo2 instructions",
		".github/instructions/style.instructions.md": "use prettier",
	})

	cfg := &config.Config{
		Sources: config.Sources{
			Dirs: []string{dir1, dir2},
		},
		Cache: config.CacheConfig{
			Dir:          t.TempDir(),
			SyncInterval: time.Hour,
		},
	}
	ldr := newLoader(cfg)
	instructions := ldr.List()

	if len(instructions) != 3 {
		t.Fatalf("got %d instructions, want 3", len(instructions))
	}
}

func TestLoaderListNoSources(t *testing.T) {
	cfg := &config.Config{
		Cache: config.CacheConfig{
			Dir:          t.TempDir(),
			SyncInterval: time.Hour,
		},
	}
	ldr := newLoader(cfg)
	instructions := ldr.List()
	if len(instructions) != 0 {
		t.Errorf("got %d instructions from empty config", len(instructions))
	}
}

func TestLoaderGetByURI(t *testing.T) {
	dir := createTestDir(t, map[string]string{
		".github/copilot-instructions.md":            "main instructions",
		".github/instructions/style.instructions.md": "style guide",
	})

	cfg := &config.Config{
		Sources: config.Sources{Dirs: []string{dir}},
		Cache: config.CacheConfig{
			Dir:          t.TempDir(),
			SyncInterval: time.Hour,
		},
	}
	ldr := newLoader(cfg)
	source := filepath.Base(dir)

	inst, ok := ldr.Get("instructions://" + source + "/copilot-instructions")
	if !ok {
		t.Fatal("Get returned false for copilot-instructions")
	}
	if inst.Content != "main instructions" {
		t.Errorf("Content = %q", inst.Content)
	}

	inst, ok = ldr.Get("instructions://" + source + "/style")
	if !ok {
		t.Fatal("Get returned false for style")
	}
	if inst.Content != "style guide" {
		t.Errorf("Content = %q", inst.Content)
	}
}

func TestLoaderGetNotFound(t *testing.T) {
	cfg := &config.Config{
		Cache: config.CacheConfig{
			Dir:          t.TempDir(),
			SyncInterval: time.Hour,
		},
	}
	ldr := newLoader(cfg)
	_, ok := ldr.Get("instructions://nonexistent/nope")
	if ok {
		t.Error("Get returned true for nonexistent URI")
	}
}

func TestLoaderRepoCacheDir(t *testing.T) {
	cacheDir := t.TempDir()

	tests := []struct {
		ref  config.RepoRef
		want string
	}{
		{config.RepoRef{Owner: "owner", Repo: "repo"}, "owner_repo"},
		{config.RepoRef{Owner: "owner", Repo: "repo", Ref: "main"}, "owner_repo_main"},
		{config.RepoRef{Owner: "org", Repo: "lib", Ref: "v1.2"}, "org_lib_v1.2"},
	}
	for _, tt := range tests {
		got := repoCacheDir(cacheDir, tt.ref)
		expected := filepath.Join(cacheDir, tt.want)
		if got != expected {
			t.Errorf("repoCacheDir(%+v) = %q, want %q", tt.ref, got, expected)
		}
	}
}

// Mock GitHub API server for testing sync.
func newMockGitHubServer(t *testing.T) *httptest.Server {
	t.Helper()
	mux := http.NewServeMux()

	mux.HandleFunc("/repos/testowner/testrepo/contents/.github/copilot-instructions.md", func(w http.ResponseWriter, r *http.Request) {
		if r.Header.Get("Accept") != "application/vnd.github.raw+json" {
			t.Errorf("unexpected Accept header: %q", r.Header.Get("Accept"))
		}
		if r.Header.Get("Authorization") != "Bearer test-token" {
			t.Errorf("unexpected Authorization header: %q", r.Header.Get("Authorization"))
		}
		w.Write([]byte("# Test Repo Instructions\nBuild with make"))
	})

	mux.HandleFunc("/repos/testowner/testrepo/contents/.github/instructions", func(w http.ResponseWriter, r *http.Request) {
		entries := []github.ContentEntry{
			{Name: "go.instructions.md", Path: ".github/instructions/go.instructions.md", Type: "file"},
			{Name: "README.md", Path: ".github/instructions/README.md", Type: "file"},
			{Name: "ci.instructions.md", Path: ".github/instructions/ci.instructions.md", Type: "file"},
		}
		json.NewEncoder(w).Encode(entries)
	})

	mux.HandleFunc("/repos/testowner/testrepo/contents/.github/instructions/go.instructions.md", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Use gofmt and go vet"))
	})
	mux.HandleFunc("/repos/testowner/testrepo/contents/.github/instructions/ci.instructions.md", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Run tests in CI"))
	})

	return httptest.NewServer(mux)
}

func TestLoaderSyncRepo(t *testing.T) {
	server := newMockGitHubServer(t)
	defer server.Close()

	cacheDir := t.TempDir()
	cfg := &config.Config{
		Sources: config.Sources{
			Repos: []string{"testowner/testrepo"},
		},
		Cache: config.CacheConfig{
			Dir:          cacheDir,
			SyncInterval: time.Hour,
		},
	}

	ldr := newLoaderWithGH(cfg, server.URL, "test-token")

	ref := config.RepoRef{Owner: "testowner", Repo: "testrepo"}
	err := ldr.syncRepo(ref)
	if err != nil {
		t.Fatalf("syncRepo failed: %v", err)
	}

	// Verify cached files
	cachedCopilot := filepath.Join(repoCacheDir(cacheDir, ref), ".github", "copilot-instructions.md")
	content, err := os.ReadFile(cachedCopilot)
	if err != nil {
		t.Fatalf("cached copilot-instructions.md not found: %v", err)
	}
	if string(content) != "# Test Repo Instructions\nBuild with make" {
		t.Errorf("cached content = %q", string(content))
	}

	// Verify instruction files
	cachedGo := filepath.Join(repoCacheDir(cacheDir, ref), ".github", "instructions", "go.instructions.md")
	content, err = os.ReadFile(cachedGo)
	if err != nil {
		t.Fatalf("cached go.instructions.md not found: %v", err)
	}
	if string(content) != "Use gofmt and go vet" {
		t.Errorf("cached go content = %q", string(content))
	}

	// Verify non-instruction files are NOT cached
	cachedReadme := filepath.Join(repoCacheDir(cacheDir, ref), ".github", "instructions", "README.md")
	if _, err := os.Stat(cachedReadme); err == nil {
		t.Error("README.md should not be cached (not .instructions.md)")
	}
}

func TestLoaderSyncRepoListViaLoader(t *testing.T) {
	server := newMockGitHubServer(t)
	defer server.Close()

	cacheDir := t.TempDir()
	cfg := &config.Config{
		Sources: config.Sources{
			Repos: []string{"testowner/testrepo"},
		},
		Cache: config.CacheConfig{
			Dir:          cacheDir,
			SyncInterval: time.Hour,
		},
	}

	ldr := newLoaderWithGH(cfg, server.URL, "test-token")
	ldr.syncAllRepos()

	instructions := ldr.List()
	if len(instructions) != 3 {
		t.Fatalf("got %d instructions, want 3 (copilot + go + ci)", len(instructions))
	}

	uris := make(map[string]bool)
	for _, inst := range instructions {
		uris[inst.URI] = true
	}
	want := []string{
		"instructions://testowner/testrepo/copilot-instructions",
		"instructions://testowner/testrepo/go",
		"instructions://testowner/testrepo/ci",
	}
	for _, u := range want {
		if !uris[u] {
			t.Errorf("missing URI: %s", u)
		}
	}
}

func TestLoaderSyncRepoWithRef(t *testing.T) {
	mux := http.NewServeMux()
	mux.HandleFunc("/repos/owner/repo/contents/.github/copilot-instructions.md", func(w http.ResponseWriter, r *http.Request) {
		ref := r.URL.Query().Get("ref")
		if ref != "develop" {
			t.Errorf("expected ref=develop, got %q", ref)
		}
		w.Write([]byte("develop branch instructions"))
	})
	mux.HandleFunc("/repos/owner/repo/contents/.github/instructions", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(404)
	})
	server := httptest.NewServer(mux)
	defer server.Close()

	cacheDir := t.TempDir()
	cfg := &config.Config{
		Sources: config.Sources{Repos: []string{"owner/repo@develop"}},
		Cache:   config.CacheConfig{Dir: cacheDir, SyncInterval: time.Hour},
	}
	ldr := newLoaderWithGH(cfg, server.URL, "")
	ldr.syncAllRepos()

	instructions := ldr.List()
	if len(instructions) != 1 {
		t.Fatalf("got %d instructions, want 1", len(instructions))
	}
	if instructions[0].Content != "develop branch instructions" {
		t.Errorf("Content = %q", instructions[0].Content)
	}
}

func TestLoaderSyncRepoNotFound(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(404)
	}))
	defer server.Close()

	cfg := &config.Config{
		Sources: config.Sources{Repos: []string{"owner/notfound"}},
		Cache:   config.CacheConfig{Dir: t.TempDir(), SyncInterval: time.Hour},
	}
	ldr := newLoaderWithGH(cfg, server.URL, "")
	ldr.syncAllRepos()

	instructions := ldr.List()
	if len(instructions) != 0 {
		t.Errorf("got %d instructions from 404 repo", len(instructions))
	}
}

func TestLoaderStartStop(t *testing.T) {
	dir := createTestDir(t, map[string]string{
		".github/copilot-instructions.md": "test",
	})

	cfg := &config.Config{
		Sources: config.Sources{Dirs: []string{dir}},
		Cache:   config.CacheConfig{Dir: t.TempDir(), SyncInterval: 100 * time.Millisecond},
	}
	ldr := newLoader(cfg)

	ctx, cancel := context.WithCancel(context.Background())
	ldr.Start(ctx)

	instructions := ldr.List()
	if len(instructions) != 1 {
		t.Fatalf("got %d instructions after Start", len(instructions))
	}

	cancel()
	ldr.Stop()
}

func TestLoaderForceSync(t *testing.T) {
	server := newMockGitHubServer(t)
	defer server.Close()

	cfg := &config.Config{
		Sources: config.Sources{Repos: []string{"testowner/testrepo"}},
		Cache:   config.CacheConfig{Dir: t.TempDir(), SyncInterval: time.Hour},
	}
	ldr := newLoaderWithGH(cfg, server.URL, "test-token")

	instructions := ldr.List()
	if len(instructions) != 0 {
		t.Fatalf("expected 0 instructions before sync, got %d", len(instructions))
	}

	ldr.ForceSync()

	instructions = ldr.List()
	if len(instructions) != 3 {
		t.Fatalf("expected 3 instructions after ForceSync, got %d", len(instructions))
	}
}

func TestLoaderMixedLocalAndRemote(t *testing.T) {
	localDir := createTestDir(t, map[string]string{
		".github/copilot-instructions.md": "local instructions",
	})

	server := newMockGitHubServer(t)
	defer server.Close()

	cfg := &config.Config{
		Sources: config.Sources{
			Dirs:  []string{localDir},
			Repos: []string{"testowner/testrepo"},
		},
		Cache: config.CacheConfig{Dir: t.TempDir(), SyncInterval: time.Hour},
	}
	ldr := newLoaderWithGH(cfg, server.URL, "test-token")
	ldr.ForceSync()

	instructions := ldr.List()
	if len(instructions) != 4 {
		t.Fatalf("got %d instructions, want 4", len(instructions))
	}

	sources := make(map[string]int)
	for _, inst := range instructions {
		sources[inst.Source]++
	}
	localSource := filepath.Base(localDir)
	if sources[localSource] != 1 {
		t.Errorf("local source count = %d, want 1", sources[localSource])
	}
	if sources["testowner/testrepo"] != 3 {
		t.Errorf("remote source count = %d, want 3", sources["testowner/testrepo"])
	}
}

func TestLoaderLocalReadsLive(t *testing.T) {
	dir := createTestDir(t, map[string]string{
		".github/copilot-instructions.md": "version 1",
	})

	cfg := &config.Config{
		Sources: config.Sources{Dirs: []string{dir}},
		Cache:   config.CacheConfig{Dir: t.TempDir(), SyncInterval: time.Hour},
	}
	ldr := newLoader(cfg)

	instructions := ldr.List()
	if instructions[0].Content != "version 1" {
		t.Fatalf("initial content = %q", instructions[0].Content)
	}

	os.WriteFile(filepath.Join(dir, ".github", "copilot-instructions.md"), []byte("version 2"), 0o644)

	instructions = ldr.List()
	if instructions[0].Content != "version 2" {
		t.Errorf("after update, content = %q, want 'version 2'", instructions[0].Content)
	}
}

func TestScanDirUnreadableFile(t *testing.T) {
	if runtime.GOOS == "windows" {
		t.Skip("file permissions not enforced the same way on Windows")
	}
	dir := createTestDir(t, map[string]string{
		".github/instructions/secret.instructions.md": "secret content",
	})
	os.Chmod(filepath.Join(dir, ".github", "instructions", "secret.instructions.md"), 0o000)
	t.Cleanup(func() {
		os.Chmod(filepath.Join(dir, ".github", "instructions", "secret.instructions.md"), 0o644)
	})

	instructions := scanDir(dir, "unreadable")
	if len(instructions) != 0 {
		t.Errorf("got %d instructions, want 0 (file unreadable)", len(instructions))
	}
}

func TestScanDirSymlink(t *testing.T) {
	if runtime.GOOS == "windows" {
		t.Skip("symlinks may require special privileges on Windows")
	}
	dir := t.TempDir()
	actualFile := filepath.Join(dir, "actual.instructions.md")
	os.WriteFile(actualFile, []byte("symlinked content"), 0o644)

	instrDir := filepath.Join(dir, ".github", "instructions")
	os.MkdirAll(instrDir, 0o755)
	os.Symlink(actualFile, filepath.Join(instrDir, "linked.instructions.md"))

	instructions := scanDir(dir, "symrepo")
	if len(instructions) != 1 {
		t.Fatalf("got %d instructions, want 1", len(instructions))
	}
	if instructions[0].Content != "symlinked content" {
		t.Errorf("Content = %q, want %q", instructions[0].Content, "symlinked content")
	}
	if instructions[0].URI != "instructions://symrepo/linked" {
		t.Errorf("URI = %q", instructions[0].URI)
	}
}

func TestLoaderConcurrentList(t *testing.T) {
	dir := createTestDir(t, map[string]string{
		".github/copilot-instructions.md":        "concurrent",
		".github/instructions/a.instructions.md": "aaa",
		".github/instructions/b.instructions.md": "bbb",
	})
	cfg := &config.Config{
		Sources: config.Sources{Dirs: []string{dir}},
		Cache:   config.CacheConfig{Dir: t.TempDir(), SyncInterval: time.Hour},
	}
	ldr := newLoader(cfg)

	var wg sync.WaitGroup
	for i := 0; i < 10; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			list := ldr.List()
			if len(list) != 3 {
				t.Errorf("got %d instructions, want 3", len(list))
			}
		}()
	}
	wg.Wait()
}

func TestLoaderConcurrentSyncAndList(t *testing.T) {
	server := newMockGitHubServer(t)
	defer server.Close()

	cfg := &config.Config{
		Sources: config.Sources{
			Repos: []string{"testowner/testrepo"},
		},
		Cache: config.CacheConfig{Dir: t.TempDir(), SyncInterval: time.Hour},
	}
	ldr := newLoaderWithGH(cfg, server.URL, "test-token")

	var wg sync.WaitGroup
	for i := 0; i < 5; i++ {
		wg.Add(2)
		go func() {
			defer wg.Done()
			ldr.ForceSync()
		}()
		go func() {
			defer wg.Done()
			_ = ldr.List()
		}()
	}
	wg.Wait()
}

func TestLoaderGetAllListedURIs(t *testing.T) {
	dir := createTestDir(t, map[string]string{
		".github/copilot-instructions.md":              "main",
		".github/instructions/style.instructions.md":   "style guide",
		".github/instructions/testing.instructions.md": "testing guide",
	})
	cfg := &config.Config{
		Sources: config.Sources{Dirs: []string{dir}},
		Cache:   config.CacheConfig{Dir: t.TempDir(), SyncInterval: time.Hour},
	}
	ldr := newLoader(cfg)

	listed := ldr.List()
	if len(listed) == 0 {
		t.Fatal("no instructions listed")
	}
	for _, inst := range listed {
		got, ok := ldr.Get(inst.URI)
		if !ok {
			t.Errorf("Get(%q) returned false", inst.URI)
			continue
		}
		if got.Content != inst.Content {
			t.Errorf("Get(%q).Content = %q, want %q", inst.URI, got.Content, inst.Content)
		}
	}
}

func TestScanDirLargeFile(t *testing.T) {
	dir := t.TempDir()
	instrDir := filepath.Join(dir, ".github", "instructions")
	os.MkdirAll(instrDir, 0o755)

	largeContent := strings.Repeat("A", 1024*1024)
	os.WriteFile(filepath.Join(instrDir, "large.instructions.md"), []byte(largeContent), 0o644)

	instructions := scanDir(dir, "largerepo")
	if len(instructions) != 1 {
		t.Fatalf("got %d instructions, want 1", len(instructions))
	}
	if len(instructions[0].Content) != 1024*1024 {
		t.Errorf("content length = %d, want %d", len(instructions[0].Content), 1024*1024)
	}
}

func TestLoaderEmptyInstructionFile(t *testing.T) {
	dir := createTestDir(t, map[string]string{
		".github/instructions/empty.instructions.md": "",
	})
	cfg := &config.Config{
		Sources: config.Sources{Dirs: []string{dir}},
		Cache:   config.CacheConfig{Dir: t.TempDir(), SyncInterval: time.Hour},
	}
	ldr := newLoader(cfg)

	instructions := ldr.List()
	if len(instructions) != 1 {
		t.Fatalf("got %d instructions, want 1", len(instructions))
	}
	if instructions[0].Content != "" {
		t.Errorf("Content = %q, want empty string", instructions[0].Content)
	}
	if instructions[0].URI != "instructions://"+filepath.Base(dir)+"/empty" {
		t.Errorf("URI = %q", instructions[0].URI)
	}
}

func TestScanDirSpecialCharsInFilename(t *testing.T) {
	dir := t.TempDir()
	instrDir := filepath.Join(dir, ".github", "instructions")
	os.MkdirAll(instrDir, 0o755)

	names := []string{
		"my file.instructions.md",
		"my-file.instructions.md",
		"my.extra.dots.instructions.md",
	}
	for _, name := range names {
		os.WriteFile(filepath.Join(instrDir, name), []byte("content of "+name), 0o644)
	}

	instructions := scanDir(dir, "special")
	if len(instructions) != len(names) {
		t.Fatalf("got %d instructions, want %d", len(instructions), len(names))
	}
	sort.Slice(instructions, func(i, j int) bool {
		return instructions[i].URI < instructions[j].URI
	})
	for _, inst := range instructions {
		if !strings.HasPrefix(inst.Content, "content of ") {
			t.Errorf("unexpected content: %q", inst.Content)
		}
	}
}

func TestLoaderRepoCacheDirIdempotent(t *testing.T) {
	cacheDir := t.TempDir()

	ref := config.RepoRef{Owner: "owner", Repo: "repo", Ref: "main"}
	first := repoCacheDir(cacheDir, ref)
	second := repoCacheDir(cacheDir, ref)
	third := repoCacheDir(cacheDir, ref)

	if first != second || second != third {
		t.Errorf("repoCacheDir not idempotent: %q, %q, %q", first, second, third)
	}
}

func TestLoaderSyncRepoUpdatesCache(t *testing.T) {
	var mu sync.Mutex
	version := "version 1"

	mux := http.NewServeMux()
	mux.HandleFunc("/repos/owner/cacherepo/contents/.github/copilot-instructions.md", func(w http.ResponseWriter, r *http.Request) {
		mu.Lock()
		defer mu.Unlock()
		w.Write([]byte(version))
	})
	mux.HandleFunc("/repos/owner/cacherepo/contents/.github/instructions", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(404)
	})
	server := httptest.NewServer(mux)
	defer server.Close()

	cfg := &config.Config{
		Sources: config.Sources{Repos: []string{"owner/cacherepo"}},
		Cache:   config.CacheConfig{Dir: t.TempDir(), SyncInterval: time.Hour},
	}
	ldr := newLoaderWithGH(cfg, server.URL, "")

	// First sync
	ldr.ForceSync()
	instructions := ldr.List()
	if len(instructions) != 1 {
		t.Fatalf("got %d instructions, want 1", len(instructions))
	}
	if instructions[0].Content != "version 1" {
		t.Errorf("Content = %q, want %q", instructions[0].Content, "version 1")
	}

	// Update mock response
	mu.Lock()
	version = "version 2"
	mu.Unlock()

	// Sync again
	ldr.ForceSync()
	instructions = ldr.List()
	if len(instructions) != 1 {
		t.Fatalf("got %d instructions after second sync, want 1", len(instructions))
	}
	if instructions[0].Content != "version 2" {
		t.Errorf("Content after update = %q, want %q", instructions[0].Content, "version 2")
	}
}

func TestSyncAllRepos_ContinuesOnError(t *testing.T) {
	mux := http.NewServeMux()

	// fail-repo always returns 500
	mux.HandleFunc("/repos/owner/fail-repo/", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(500)
	})

	// ok-repo has valid content
	mux.HandleFunc("/repos/owner/ok-repo/contents/.github/copilot-instructions.md", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("OK repo instructions"))
	})
	mux.HandleFunc("/repos/owner/ok-repo/contents/.github/instructions", func(w http.ResponseWriter, r *http.Request) {
		json.NewEncoder(w).Encode([]github.ContentEntry{})
	})

	server := httptest.NewServer(mux)
	defer server.Close()

	cacheDir := t.TempDir()
	cfg := &config.Config{
		Sources: config.Sources{
			Repos: []string{"owner/fail-repo", "owner/ok-repo"},
		},
		Cache: config.CacheConfig{
			Dir:          cacheDir,
			SyncInterval: time.Hour,
		},
	}

	ldr := newLoaderWithGH(cfg, server.URL, "")
	ldr.syncAllRepos()

	// Verify ok-repo was still synced despite fail-repo erroring
	instructions := ldr.List()
	found := false
	for _, instr := range instructions {
		if strings.Contains(instr.Source, "ok-repo") || strings.Contains(instr.Content, "OK repo") {
			found = true
		}
	}
	if !found {
		t.Error("syncAllRepos should continue to ok-repo after fail-repo errors")
	}
}

// ---------------------------------------------------------------------------
// Additional nominal / error / limit tests
// ---------------------------------------------------------------------------

func TestLoaderListEmptyConfig(t *testing.T) {
	cfg := &config.Config{}
	ldr := newLoader(cfg)
	instructions := ldr.List()
	if len(instructions) != 0 {
		t.Errorf("expected 0 instructions, got %d", len(instructions))
	}
}

func TestLoaderGetNotFoundEmpty(t *testing.T) {
	cfg := &config.Config{}
	ldr := newLoader(cfg)
	_, found := ldr.Get("instructions://nonexistent/foo")
	if found {
		t.Error("expected not found on empty loader")
	}
}

func TestSyncRepo_NoCopilotInstructionsButHasDir(t *testing.T) {
	// Repo has no copilot-instructions.md but has .github/instructions/ dir
	mux := http.NewServeMux()
	mux.HandleFunc("/repos/o/r/contents/.github/copilot-instructions.md", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(404)
	})
	mux.HandleFunc("/repos/o/r/contents/.github/instructions", func(w http.ResponseWriter, r *http.Request) {
		json.NewEncoder(w).Encode([]github.ContentEntry{
			{Name: "go.instructions.md", Path: ".github/instructions/go.instructions.md", Type: "file"},
		})
	})
	mux.HandleFunc("/repos/o/r/contents/.github/instructions/go.instructions.md", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Use Go"))
	})
	server := httptest.NewServer(mux)
	defer server.Close()

	cacheDir := t.TempDir()
	cfg := &config.Config{
		Sources: config.Sources{Repos: []string{"o/r"}},
		Cache:   config.CacheConfig{Dir: cacheDir, SyncInterval: time.Hour},
	}
	ldr := newLoaderWithGH(cfg, server.URL, "")

	ref := config.ParseRepoRef("o/r")
	if err := ldr.syncRepo(ref); err != nil {
		t.Fatalf("syncRepo: %v", err)
	}

	instructions := scanDir(repoCacheDir(cacheDir, ref), "o/r")
	if len(instructions) != 1 {
		t.Fatalf("got %d instructions, want 1", len(instructions))
	}
	if instructions[0].Content != "Use Go" {
		t.Errorf("content = %q", instructions[0].Content)
	}
}

func TestSyncRepo_InstructionFetchFails(t *testing.T) {
	// Repo lists instructions dir but individual file fetch fails
	mux := http.NewServeMux()
	mux.HandleFunc("/repos/o/r/contents/.github/copilot-instructions.md", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(404)
	})
	mux.HandleFunc("/repos/o/r/contents/.github/instructions", func(w http.ResponseWriter, r *http.Request) {
		json.NewEncoder(w).Encode([]github.ContentEntry{
			{Name: "good.instructions.md", Path: ".github/instructions/good.instructions.md", Type: "file"},
			{Name: "fail.instructions.md", Path: ".github/instructions/fail.instructions.md", Type: "file"},
		})
	})
	mux.HandleFunc("/repos/o/r/contents/.github/instructions/good.instructions.md", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Good content"))
	})
	mux.HandleFunc("/repos/o/r/contents/.github/instructions/fail.instructions.md", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(500)
	})
	server := httptest.NewServer(mux)
	defer server.Close()

	cacheDir := t.TempDir()
	cfg := &config.Config{
		Sources: config.Sources{Repos: []string{"o/r"}},
		Cache:   config.CacheConfig{Dir: cacheDir, SyncInterval: time.Hour},
	}
	ldr := newLoaderWithGH(cfg, server.URL, "")

	ref := config.ParseRepoRef("o/r")
	err := ldr.syncRepo(ref)
	// syncRepo should NOT return error — it continues on individual file failures
	if err != nil {
		t.Fatalf("syncRepo: %v", err)
	}

	instructions := scanDir(repoCacheDir(cacheDir, ref), "o/r")
	if len(instructions) != 1 {
		t.Fatalf("got %d instructions, want 1 (only good.instructions.md)", len(instructions))
	}
}

func TestSyncRepo_NonInstructionFilesSkipped(t *testing.T) {
	mux := http.NewServeMux()
	mux.HandleFunc("/repos/o/r/contents/.github/copilot-instructions.md", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(404)
	})
	mux.HandleFunc("/repos/o/r/contents/.github/instructions", func(w http.ResponseWriter, r *http.Request) {
		json.NewEncoder(w).Encode([]github.ContentEntry{
			{Name: "README.md", Path: ".github/instructions/README.md", Type: "file"},
			{Name: "notes.txt", Path: ".github/instructions/notes.txt", Type: "file"},
			{Name: "go.instructions.md", Path: ".github/instructions/go.instructions.md", Type: "file"},
		})
	})
	mux.HandleFunc("/repos/o/r/contents/.github/instructions/go.instructions.md", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Go instructions"))
	})
	server := httptest.NewServer(mux)
	defer server.Close()

	cacheDir := t.TempDir()
	cfg := &config.Config{
		Sources: config.Sources{Repos: []string{"o/r"}},
		Cache:   config.CacheConfig{Dir: cacheDir, SyncInterval: time.Hour},
	}
	ldr := newLoaderWithGH(cfg, server.URL, "")

	ref := config.ParseRepoRef("o/r")
	ldr.syncRepo(ref)

	instructions := scanDir(repoCacheDir(cacheDir, ref), "o/r")
	if len(instructions) != 1 {
		t.Fatalf("got %d instructions, want 1", len(instructions))
	}
	if instructions[0].Content != "Go instructions" {
		t.Errorf("content = %q", instructions[0].Content)
	}
}

func TestLoaderGetMultipleSources(t *testing.T) {
	dir1 := createTestDir(t, map[string]string{
		".github/copilot-instructions.md": "From dir1",
	})
	dir2 := createTestDir(t, map[string]string{
		".github/copilot-instructions.md": "From dir2",
	})

	cfg := &config.Config{
		Sources: config.Sources{Dirs: []string{dir1, dir2}},
	}
	ldr := newLoader(cfg)

	instructions := ldr.List()
	if len(instructions) != 2 {
		t.Fatalf("got %d instructions, want 2", len(instructions))
	}

	// Get should find the first one
	got, found := ldr.Get(instructions[0].URI)
	if !found {
		t.Fatal("expected to find instruction")
	}
	if got.Content != instructions[0].Content {
		t.Errorf("content = %q, want %q", got.Content, instructions[0].Content)
	}
}

func TestScanDirPermissionDenied(t *testing.T) {
	if runtime.GOOS == "windows" {
		t.Skip("permission test not reliable on Windows")
	}

	dir := t.TempDir()
	ghDir := filepath.Join(dir, ".github", "instructions")
	if err := os.MkdirAll(ghDir, 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(ghDir, "test.instructions.md"), []byte("content"), 0o000); err != nil {
		t.Fatal(err)
	}

	instructions := scanDir(dir, "test")
	// Should gracefully skip unreadable file
	if len(instructions) != 0 {
		t.Errorf("got %d instructions, want 0 (file unreadable)", len(instructions))
	}
}

func TestScanDir_OnlyCopilotInstructions(t *testing.T) {
	dir := createTestDir(t, map[string]string{
		".github/copilot-instructions.md": "Only copilot",
	})
	instructions := scanDir(dir, "test")
	if len(instructions) != 1 {
		t.Fatalf("got %d, want 1", len(instructions))
	}
	if instructions[0].URI != "instructions://test/copilot-instructions" {
		t.Errorf("URI = %q", instructions[0].URI)
	}
}

func TestScanDir_OnlyPathSpecific(t *testing.T) {
	dir := createTestDir(t, map[string]string{
		".github/instructions/go.instructions.md": "Go only",
	})
	instructions := scanDir(dir, "test")
	if len(instructions) != 1 {
		t.Fatalf("got %d, want 1", len(instructions))
	}
	if instructions[0].URI != "instructions://test/go" {
		t.Errorf("URI = %q", instructions[0].URI)
	}
}

func TestLoaderStartStopMultipleTimes(t *testing.T) {
	cfg := &config.Config{
		Cache: config.CacheConfig{SyncInterval: time.Hour},
	}
	ldr := newLoader(cfg)

	ctx := context.Background()
	ldr.Start(ctx)
	ldr.Stop()
	// Stop should be idempotent
	ldr.Stop()
}

// TestSyncAllRepos_MkdirFailLogsAndContinues covers the log.Printf branch in
// syncAllRepos when syncRepo returns an error (os.MkdirAll fails because a
// regular file already occupies the expected cache directory path).
func TestSyncAllRepos_MkdirFailLogsAndContinues(t *testing.T) {
	cacheDir := t.TempDir()
	// repoCacheDir("x", "y") == cacheDir/"x_y"; create it as a file so that
	// os.MkdirAll(cacheDir/x_y/.github, 0755) fails.
	if err := os.WriteFile(filepath.Join(cacheDir, "x_y"), []byte("block"), 0o644); err != nil {
		t.Fatal(err)
	}

	cfg := &config.Config{
		Sources: config.Sources{Repos: []string{"x/y"}},
		Cache:   config.CacheConfig{Dir: cacheDir, SyncInterval: time.Hour},
	}
	ldr := newLoader(cfg)
	// Must not panic; the error is logged and execution continues.
	ldr.syncAllRepos()
}

func TestLoaderForceSyncWithRemote(t *testing.T) {
	callCount := 0
	mux := http.NewServeMux()
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		callCount++
		if strings.HasSuffix(r.URL.Path, "/copilot-instructions.md") {
			w.Write([]byte("synced"))
		} else if strings.HasSuffix(r.URL.Path, "/instructions") {
			json.NewEncoder(w).Encode([]github.ContentEntry{})
		} else {
			w.WriteHeader(404)
		}
	})
	server := httptest.NewServer(mux)
	defer server.Close()

	cacheDir := t.TempDir()
	cfg := &config.Config{
		Sources: config.Sources{Repos: []string{"o/r"}},
		Cache:   config.CacheConfig{Dir: cacheDir, SyncInterval: time.Hour},
	}
	ldr := newLoaderWithGH(cfg, server.URL, "")

	ldr.ForceSync()
	if callCount == 0 {
		t.Error("ForceSync should have triggered API calls")
	}

	first := callCount
	ldr.ForceSync()
	if callCount <= first {
		t.Error("second ForceSync should make additional API calls")
	}
}

// ---------------------------------------------------------------------------
// Root-level scanning (not just .github/)
// ---------------------------------------------------------------------------

func TestScanDir_RootLevelCopilotInstructions(t *testing.T) {
	// copilot-instructions.md at root (no .github/ prefix)
	dir := createTestDir(t, map[string]string{
		"copilot-instructions.md": "root instructions",
	})
	instructions := scanDir(dir, "rootrepo")
	if len(instructions) != 1 {
		t.Fatalf("got %d instructions, want 1", len(instructions))
	}
	if instructions[0].Content != "root instructions" {
		t.Errorf("Content = %q", instructions[0].Content)
	}
	if instructions[0].URI != "instructions://rootrepo/copilot-instructions" {
		t.Errorf("URI = %q", instructions[0].URI)
	}
}

func TestScanDir_RootLevelInstructionsDir(t *testing.T) {
	// instructions/*.instructions.md at root (no .github/ prefix)
	dir := createTestDir(t, map[string]string{
		"instructions/style.instructions.md": "root style guide",
	})
	instructions := scanDir(dir, "rootrepo")
	if len(instructions) != 1 {
		t.Fatalf("got %d instructions, want 1", len(instructions))
	}
	if instructions[0].Content != "root style guide" {
		t.Errorf("Content = %q", instructions[0].Content)
	}
	if instructions[0].URI != "instructions://rootrepo/style" {
		t.Errorf("URI = %q", instructions[0].URI)
	}
}

func TestScanDir_GithubTakesPriorityOverRoot(t *testing.T) {
	// When the same file exists in both .github/ and root, .github/ wins.
	dir := createTestDir(t, map[string]string{
		".github/copilot-instructions.md": "from github",
		"copilot-instructions.md":         "from root",
	})
	instructions := scanDir(dir, "repo")
	if len(instructions) != 1 {
		t.Fatalf("got %d instructions, want 1 (deduplication)", len(instructions))
	}
	if instructions[0].Content != "from github" {
		t.Errorf("expected .github/ version to win, got %q", instructions[0].Content)
	}
}

func TestScanDir_BothGithubAndRootDifferentFiles(t *testing.T) {
	// Files in .github/ and different files at root are both included.
	dir := createTestDir(t, map[string]string{
		".github/copilot-instructions.md":    "main instructions",
		"instructions/extra.instructions.md": "extra root guide",
	})
	instructions := scanDir(dir, "repo")
	if len(instructions) != 2 {
		t.Fatalf("got %d instructions, want 2", len(instructions))
	}
}
