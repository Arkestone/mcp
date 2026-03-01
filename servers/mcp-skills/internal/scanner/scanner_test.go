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

func createTestSkill(t *testing.T, dir, skillName, name, desc, body string) {
	t.Helper()
	skillDir := filepath.Join(dir, skillName)
	if err := os.MkdirAll(skillDir, 0o755); err != nil {
		t.Fatal(err)
	}
	content := fmt.Sprintf("---\nname: %s\ndescription: %s\n---\n%s", name, desc, body)
	if err := os.WriteFile(filepath.Join(skillDir, "SKILL.md"), []byte(content), 0o644); err != nil {
		t.Fatal(err)
	}
}

func createTestSkillRaw(t *testing.T, dir, skillName, raw string) {
	t.Helper()
	skillDir := filepath.Join(dir, skillName)
	if err := os.MkdirAll(skillDir, 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(skillDir, "SKILL.md"), []byte(raw), 0o644); err != nil {
		t.Fatal(err)
	}
}

func createTestRef(t *testing.T, dir, skillName, refName, content string) {
	t.Helper()
	refDir := filepath.Join(dir, skillName, "references")
	if err := os.MkdirAll(refDir, 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(refDir, refName), []byte(content), 0o644); err != nil {
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

func newScannerWithGH(cfg *config.Config, baseURL, token string) *Scanner {
	return New(cfg, &github.Client{BaseURL: baseURL, Token: token})
}

// ---------------------------------------------------------------------------
// parseFrontmatter tests
// ---------------------------------------------------------------------------

func TestParseFrontmatter_Basic(t *testing.T) {
	data := []byte("---\nname: Docker Expert\ndescription: Helps with Docker\n---\nBody markdown here\n")
	name, desc, body := parseFrontmatter(data)
	if name != "Docker Expert" {
		t.Errorf("name = %q, want %q", name, "Docker Expert")
	}
	if desc != "Helps with Docker" {
		t.Errorf("description = %q, want %q", desc, "Helps with Docker")
	}
	if body != "Body markdown here\n" {
		t.Errorf("body = %q, want %q", body, "Body markdown here\n")
	}
}

func TestParseFrontmatter_NoFrontmatter(t *testing.T) {
	data := []byte("Just plain markdown\nNo frontmatter here\n")
	name, desc, body := parseFrontmatter(data)
	if name != "" {
		t.Errorf("name = %q, want empty", name)
	}
	if desc != "" {
		t.Errorf("description = %q, want empty", desc)
	}
	if body != "Just plain markdown\nNo frontmatter here\n" {
		t.Errorf("body = %q", body)
	}
}

func TestParseFrontmatter_EmptyFrontmatter(t *testing.T) {
	data := []byte("---\n---\nBody only\n")
	name, desc, body := parseFrontmatter(data)
	if name != "" {
		t.Errorf("name = %q, want empty", name)
	}
	if desc != "" {
		t.Errorf("description = %q, want empty", desc)
	}
	if body != "Body only\n" {
		t.Errorf("body = %q, want %q", body, "Body only\n")
	}
}

func TestParseFrontmatter_MultilineDescription(t *testing.T) {
	data := []byte("---\nname: MySkill\ndescription: |\n  line one\n  line two\n---\nBody\n")
	name, desc, body := parseFrontmatter(data)
	if name != "MySkill" {
		t.Errorf("name = %q, want %q", name, "MySkill")
	}
	if !strings.Contains(desc, "line one") || !strings.Contains(desc, "line two") {
		t.Errorf("description = %q, want multiline content", desc)
	}
	if body != "Body\n" {
		t.Errorf("body = %q", body)
	}
}

func TestParseFrontmatter_DescriptionWithBlockScalar(t *testing.T) {
	data := []byte("---\nname: Folded\ndescription: >\n  folded line one\n  folded line two\n---\nContent\n")
	name, desc, body := parseFrontmatter(data)
	if name != "Folded" {
		t.Errorf("name = %q, want %q", name, "Folded")
	}
	if !strings.Contains(desc, "folded line one") {
		t.Errorf("description = %q, want folded content", desc)
	}
	if body != "Content\n" {
		t.Errorf("body = %q", body)
	}
}

func TestParseFrontmatter_QuotedValues(t *testing.T) {
	data := []byte("---\nname: \"quoted name\"\ndescription: 'quoted desc'\n---\nBody\n")
	name, desc, _ := parseFrontmatter(data)
	if name != "quoted name" {
		t.Errorf("name = %q, want %q", name, "quoted name")
	}
	if desc != "quoted desc" {
		t.Errorf("description = %q, want %q", desc, "quoted desc")
	}
}

func TestParseFrontmatter_OnlyName(t *testing.T) {
	data := []byte("---\nname: JustName\n---\nBody\n")
	name, desc, _ := parseFrontmatter(data)
	if name != "JustName" {
		t.Errorf("name = %q, want %q", name, "JustName")
	}
	if desc != "" {
		t.Errorf("description = %q, want empty", desc)
	}
}

func TestParseFrontmatter_OnlyDescription(t *testing.T) {
	data := []byte("---\ndescription: JustDesc\n---\nBody\n")
	name, desc, _ := parseFrontmatter(data)
	if name != "" {
		t.Errorf("name = %q, want empty", name)
	}
	if desc != "JustDesc" {
		t.Errorf("description = %q, want %q", desc, "JustDesc")
	}
}

func TestParseFrontmatter_EmptyBody(t *testing.T) {
	data := []byte("---\nname: NoBody\ndescription: desc\n---\n")
	name, desc, body := parseFrontmatter(data)
	if name != "NoBody" {
		t.Errorf("name = %q", name)
	}
	if desc != "desc" {
		t.Errorf("description = %q", desc)
	}
	if body != "" {
		t.Errorf("body = %q, want empty", body)
	}
}

// TestParseFrontmatter_NoTrailingNewline tests content with no newline after the
// closing "---", making body truly empty (len==0).  Kills scanner:157
// CONDITIONALS_BOUNDARY len(body)>0 → len(body)>=0, which would panic on
// an empty-slice access.
func TestParseFrontmatter_NoTrailingNewline(t *testing.T) {
	data := []byte("---\nname: NoTrail\ndescription: trail\n---")
	name, desc, body := parseFrontmatter(data)
	if name != "NoTrail" {
		t.Errorf("name = %q, want NoTrail", name)
	}
	if desc != "trail" {
		t.Errorf("description = %q, want trail", desc)
	}
	if body != "" {
		t.Errorf("body = %q, want empty", body)
	}
}

func TestParseFrontmatter_SpecialCharacters(t *testing.T) {
	data := []byte("---\nname: \"Skill: C++ & Go!\"\ndescription: Uses <html> & \"quotes\"\n---\nBody\n")
	name, desc, _ := parseFrontmatter(data)
	if name != "Skill: C++ & Go!" {
		t.Errorf("name = %q, want %q", name, "Skill: C++ & Go!")
	}
	if desc == "" {
		t.Error("description should not be empty")
	}
}

// ---------------------------------------------------------------------------
// scanDir tests
// ---------------------------------------------------------------------------

func TestScanDir_SingleSkill(t *testing.T) {
	dir := t.TempDir()
	createTestSkill(t, dir, "docker", "Docker Expert", "Docker help", "Docker body")

	skills := scanDir(dir, "local")
	if len(skills) != 1 {
		t.Fatalf("got %d skills, want 1", len(skills))
	}
	if skills[0].Name != "Docker Expert" {
		t.Errorf("name = %q", skills[0].Name)
	}
	if skills[0].Content != "Docker body" {
		t.Errorf("content = %q", skills[0].Content)
	}
}

func TestScanDir_MultipleSkills(t *testing.T) {
	dir := t.TempDir()
	createTestSkill(t, dir, "docker", "Docker", "d1", "b1")
	createTestSkill(t, dir, "k8s", "Kubernetes", "d2", "b2")
	createTestSkill(t, dir, "go", "Go", "d3", "b3")

	skills := scanDir(dir, "local")
	if len(skills) != 3 {
		t.Fatalf("got %d skills, want 3", len(skills))
	}
}

func TestScanDir_WithReferences(t *testing.T) {
	dir := t.TempDir()
	createTestSkill(t, dir, "docker", "Docker", "d", "body")
	createTestRef(t, dir, "docker", "compose.md", "compose ref content")

	skills := scanDir(dir, "local")
	if len(skills) != 1 {
		t.Fatalf("got %d skills, want 1", len(skills))
	}
	if len(skills[0].References) != 1 {
		t.Fatalf("got %d refs, want 1", len(skills[0].References))
	}
	ref := skills[0].References[0]
	if ref.Name != "compose.md" {
		t.Errorf("ref name = %q", ref.Name)
	}
	if ref.Content != "compose ref content" {
		t.Errorf("ref content = %q", ref.Content)
	}
	if !strings.Contains(ref.Path, "references") {
		t.Errorf("ref path = %q, want to contain 'references'", ref.Path)
	}
}

func TestScanDir_EmptyDir(t *testing.T) {
	dir := t.TempDir()
	skills := scanDir(dir, "local")
	if len(skills) != 0 {
		t.Fatalf("got %d skills, want 0", len(skills))
	}
}

func TestScanDir_NoSkillMd(t *testing.T) {
	dir := t.TempDir()
	// Create a subdir but no SKILL.md in it
	os.MkdirAll(filepath.Join(dir, "empty-skill"), 0o755)

	skills := scanDir(dir, "local")
	if len(skills) != 0 {
		t.Fatalf("got %d skills, want 0", len(skills))
	}
}

func TestScanDir_MixedContent(t *testing.T) {
	dir := t.TempDir()
	createTestSkill(t, dir, "valid", "Valid", "desc", "body")
	os.MkdirAll(filepath.Join(dir, "no-skill"), 0o755)
	// also create a plain file at the top level (should be ignored)
	os.WriteFile(filepath.Join(dir, "README.md"), []byte("readme"), 0o644)

	skills := scanDir(dir, "local")
	if len(skills) != 1 {
		t.Fatalf("got %d skills, want 1", len(skills))
	}
	if skills[0].Name != "Valid" {
		t.Errorf("name = %q", skills[0].Name)
	}
}

func TestScanDir_NameFromFrontmatter(t *testing.T) {
	dir := t.TempDir()
	createTestSkill(t, dir, "dirname", "Frontmatter Name", "desc", "body")

	skills := scanDir(dir, "src")
	if len(skills) != 1 {
		t.Fatalf("got %d skills", len(skills))
	}
	if skills[0].Name != "Frontmatter Name" {
		t.Errorf("name = %q, want %q", skills[0].Name, "Frontmatter Name")
	}
}

func TestScanDir_NameFallbackToDirname(t *testing.T) {
	dir := t.TempDir()
	// Skill with no name in frontmatter
	createTestSkillRaw(t, dir, "my-skill", "---\ndescription: desc only\n---\nbody\n")

	skills := scanDir(dir, "src")
	if len(skills) != 1 {
		t.Fatalf("got %d skills", len(skills))
	}
	if skills[0].Name != "my-skill" {
		t.Errorf("name = %q, want %q", skills[0].Name, "my-skill")
	}
}

func TestScanDir_URIFormat(t *testing.T) {
	dir := t.TempDir()
	createTestSkill(t, dir, "docker", "Docker Expert", "desc", "body")

	skills := scanDir(dir, "org/repo")
	if len(skills) != 1 {
		t.Fatalf("got %d skills", len(skills))
	}
	want := "skills://org/repo/Docker Expert"
	if skills[0].URI != want {
		t.Errorf("URI = %q, want %q", skills[0].URI, want)
	}
}

func TestScanDir_NonexistentDir(t *testing.T) {
	skills := scanDir("/nonexistent/path/does/not/exist", "missing")
	if skills != nil {
		t.Errorf("got %v, want nil", skills)
	}
}

func TestScanDir_MultipleReferences(t *testing.T) {
	dir := t.TempDir()
	createTestSkill(t, dir, "docker", "Docker", "desc", "body")
	createTestRef(t, dir, "docker", "ref1.md", "content1")
	createTestRef(t, dir, "docker", "ref2.md", "content2")
	createTestRef(t, dir, "docker", "ref3.txt", "content3")

	skills := scanDir(dir, "local")
	if len(skills) != 1 {
		t.Fatalf("got %d skills", len(skills))
	}
	if len(skills[0].References) != 3 {
		t.Fatalf("got %d refs, want 3", len(skills[0].References))
	}
}

// ---------------------------------------------------------------------------
// Scanner integration tests
// ---------------------------------------------------------------------------

func TestScanner_List_LocalDirs(t *testing.T) {
	dir := t.TempDir()
	createTestSkill(t, dir, "docker", "Docker", "desc", "body")
	createTestSkill(t, dir, "go", "Go", "desc", "body")

	cfg := newTestConfig([]string{dir}, nil, t.TempDir())
	s := newScanner(cfg)

	skills := s.List()
	if len(skills) != 2 {
		t.Fatalf("got %d skills, want 2", len(skills))
	}
}

func TestScanner_Get_Found(t *testing.T) {
	dir := t.TempDir()
	createTestSkill(t, dir, "docker", "Docker", "desc", "body")

	cfg := newTestConfig([]string{dir}, nil, t.TempDir())
	s := newScanner(cfg)

	source := filepath.Base(dir)
	uri := fmt.Sprintf("skills://%s/Docker", source)
	sk, ok := s.Get(uri)
	if !ok {
		t.Fatal("Get returned false")
	}
	if sk.Name != "Docker" {
		t.Errorf("name = %q", sk.Name)
	}
}

func TestScanner_Get_NotFound(t *testing.T) {
	dir := t.TempDir()
	createTestSkill(t, dir, "docker", "Docker", "desc", "body")

	cfg := newTestConfig([]string{dir}, nil, t.TempDir())
	s := newScanner(cfg)

	_, ok := s.Get("skills://nonexistent/Missing")
	if ok {
		t.Error("Get returned true for nonexistent URI")
	}
}

func TestScanner_List_EmptyConfig(t *testing.T) {
	cfg := newTestConfig(nil, nil, t.TempDir())
	s := newScanner(cfg)

	skills := s.List()
	if len(skills) != 0 {
		t.Fatalf("got %d skills, want 0", len(skills))
	}
}

func TestScanner_StartStop(t *testing.T) {
	cfg := newTestConfig(nil, nil, t.TempDir())
	s := newScanner(cfg)

	ctx := context.Background()
	s.Start(ctx)
	// Should stop cleanly without hanging
	done := make(chan struct{})
	go func() {
		s.Stop()
		close(done)
	}()
	select {
	case <-done:
		// ok
	case <-time.After(5 * time.Second):
		t.Fatal("Stop did not return within 5 seconds")
	}
}

func TestScanner_ForceSync_NoRepos(t *testing.T) {
	cfg := newTestConfig(nil, nil, t.TempDir())
	s := newScanner(cfg)
	// Should not panic
	s.ForceSync()
}

// ---------------------------------------------------------------------------
// GitHub sync tests (httptest)
// ---------------------------------------------------------------------------

func newGitHubTestServer(t *testing.T, handler http.HandlerFunc) *httptest.Server {
	t.Helper()
	ts := httptest.NewServer(handler)
	t.Cleanup(ts.Close)
	return ts
}

func TestSyncRepo_Basic(t *testing.T) {
	ts := newGitHubTestServer(t, func(w http.ResponseWriter, r *http.Request) {
		path := r.URL.Path
		switch {
		// top-level skills/ listing
		case strings.HasSuffix(path, "/contents/skills"):
			json.NewEncoder(w).Encode([]github.ContentEntry{
				{Name: "docker", Path: "skills/docker", Type: "dir"},
			})
		// SKILL.md raw content
		case strings.HasSuffix(path, "/contents/skills/docker/SKILL.md"):
			w.Write([]byte("---\nname: Docker\ndescription: Docker help\n---\nDocker body\n"))
		// references/ listing
		case strings.HasSuffix(path, "/contents/skills/docker/references"):
			w.WriteHeader(404)
		default:
			w.WriteHeader(404)
		}
	})

	cacheDir := t.TempDir()
	cfg := newTestConfig(nil, []string{"testowner/testrepo"}, cacheDir)
	s := newScannerWithGH(cfg, ts.URL, "")

	ref := config.ParseRepoRef("testowner/testrepo")
	if err := s.syncRepo(ref); err != nil {
		t.Fatalf("syncRepo: %v", err)
	}

	// Verify the cache was populated
	skills := scanDir(repoCacheDir(cacheDir, ref), "testowner/testrepo")
	if len(skills) != 1 {
		t.Fatalf("got %d skills, want 1", len(skills))
	}
	if skills[0].Name != "Docker" {
		t.Errorf("name = %q", skills[0].Name)
	}
}

func TestSyncRepo_WithReferences(t *testing.T) {
	ts := newGitHubTestServer(t, func(w http.ResponseWriter, r *http.Request) {
		path := r.URL.Path
		switch {
		case strings.HasSuffix(path, "/contents/skills"):
			json.NewEncoder(w).Encode([]github.ContentEntry{
				{Name: "docker", Path: "skills/docker", Type: "dir"},
			})
		case strings.HasSuffix(path, "/contents/skills/docker/SKILL.md"):
			w.Write([]byte("---\nname: Docker\ndescription: d\n---\nbody\n"))
		case strings.HasSuffix(path, "/contents/skills/docker/references"):
			json.NewEncoder(w).Encode([]github.ContentEntry{
				{Name: "compose.md", Path: "skills/docker/references/compose.md", Type: "file"},
			})
		case strings.HasSuffix(path, "/contents/skills/docker/references/compose.md"):
			w.Write([]byte("compose reference content"))
		default:
			w.WriteHeader(404)
		}
	})

	cacheDir := t.TempDir()
	cfg := newTestConfig(nil, []string{"owner/repo"}, cacheDir)
	s := newScannerWithGH(cfg, ts.URL, "")

	ref := config.ParseRepoRef("owner/repo")
	if err := s.syncRepo(ref); err != nil {
		t.Fatalf("syncRepo: %v", err)
	}

	skills := scanDir(repoCacheDir(cacheDir, ref), "owner/repo")
	if len(skills) != 1 {
		t.Fatalf("got %d skills", len(skills))
	}
	if len(skills[0].References) != 1 {
		t.Fatalf("got %d refs, want 1", len(skills[0].References))
	}
	if skills[0].References[0].Content != "compose reference content" {
		t.Errorf("ref content = %q", skills[0].References[0].Content)
	}
}

func TestSyncRepo_404(t *testing.T) {
	ts := newGitHubTestServer(t, func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(404)
	})

	cacheDir := t.TempDir()
	cfg := newTestConfig(nil, []string{"owner/missing"}, cacheDir)
	s := newScannerWithGH(cfg, ts.URL, "")

	ref := config.ParseRepoRef("owner/missing")
	err := s.syncRepo(ref)
	if err == nil {
		t.Fatal("expected error for 404, got nil")
	}
	if !strings.Contains(err.Error(), "404") {
		t.Errorf("error = %q, want to contain '404'", err.Error())
	}
}

func TestSyncRepo_InvalidJSON(t *testing.T) {
	ts := newGitHubTestServer(t, func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("not json"))
	})

	cacheDir := t.TempDir()
	cfg := newTestConfig(nil, []string{"owner/bad"}, cacheDir)
	s := newScannerWithGH(cfg, ts.URL, "")

	ref := config.ParseRepoRef("owner/bad")
	err := s.syncRepo(ref)
	// fetchGitHubDir will fail to decode, so we expect an error
	if err == nil {
		t.Fatal("expected error for invalid JSON, got nil")
	}
}

func TestSyncRepo_AuthHeader(t *testing.T) {
	var gotAuth string
	ts := newGitHubTestServer(t, func(w http.ResponseWriter, r *http.Request) {
		gotAuth = r.Header.Get("Authorization")
		// Return valid response so we can verify the header was sent
		if strings.HasSuffix(r.URL.Path, "/contents/skills") {
			json.NewEncoder(w).Encode([]github.ContentEntry{})
		} else {
			w.WriteHeader(404)
		}
	})

	cacheDir := t.TempDir()
	cfg := newTestConfig(nil, []string{"owner/repo"}, cacheDir)
	s := newScannerWithGH(cfg, ts.URL, "test-token-12345")

	ref := config.ParseRepoRef("owner/repo")
	_ = s.syncRepo(ref)

	if gotAuth != "Bearer test-token-12345" {
		t.Errorf("Authorization = %q, want %q", gotAuth, "Bearer test-token-12345")
	}
}

// ---------------------------------------------------------------------------
// Additional edge-case tests
// ---------------------------------------------------------------------------

func TestScanDir_SkillPathField(t *testing.T) {
	dir := t.TempDir()
	createTestSkill(t, dir, "myskill", "MySkill", "desc", "body")

	skills := scanDir(dir, "local")
	if len(skills) != 1 {
		t.Fatalf("got %d skills", len(skills))
	}
	want := filepath.Join("myskill", "SKILL.md")
	if skills[0].Path != want {
		t.Errorf("Path = %q, want %q", skills[0].Path, want)
	}
}

func TestScanDir_SourceField(t *testing.T) {
	dir := t.TempDir()
	createTestSkill(t, dir, "myskill", "MySkill", "desc", "body")

	skills := scanDir(dir, "test-source")
	if len(skills) != 1 {
		t.Fatalf("got %d skills", len(skills))
	}
	if skills[0].Source != "test-source" {
		t.Errorf("Source = %q, want %q", skills[0].Source, "test-source")
	}
}

func TestScanDir_DescriptionField(t *testing.T) {
	dir := t.TempDir()
	createTestSkill(t, dir, "myskill", "MySkill", "My Description", "body")

	skills := scanDir(dir, "local")
	if len(skills) != 1 {
		t.Fatalf("got %d skills", len(skills))
	}
	if skills[0].Description != "My Description" {
		t.Errorf("Description = %q, want %q", skills[0].Description, "My Description")
	}
}

func TestScanner_List_MultipleLocalDirs(t *testing.T) {
	dir1 := t.TempDir()
	dir2 := t.TempDir()
	createTestSkill(t, dir1, "skill1", "Skill1", "d1", "b1")
	createTestSkill(t, dir2, "skill2", "Skill2", "d2", "b2")

	cfg := newTestConfig([]string{dir1, dir2}, nil, t.TempDir())
	s := newScanner(cfg)

	skills := s.List()
	if len(skills) != 2 {
		t.Fatalf("got %d skills, want 2", len(skills))
	}
}

func TestScanner_List_WithCachedRepo(t *testing.T) {
	// Prepopulate a cache dir to simulate a previously synced repo
	localDir := t.TempDir()
	cacheDir := t.TempDir()
	createTestSkill(t, localDir, "local-skill", "Local", "d", "b")

	// Simulate cached repo: owner_repo directory in cache
	repoCache := filepath.Join(cacheDir, "owner_repo")
	createTestSkill(t, repoCache, "cached-skill", "Cached", "d", "b")

	cfg := newTestConfig([]string{localDir}, []string{"owner/repo"}, cacheDir)
	s := newScanner(cfg)

	skills := s.List()
	if len(skills) != 2 {
		t.Fatalf("got %d skills, want 2", len(skills))
	}

	found := map[string]bool{}
	for _, sk := range skills {
		found[sk.Name] = true
	}
	if !found["Local"] {
		t.Error("missing Local skill")
	}
	if !found["Cached"] {
		t.Error("missing Cached skill")
	}
}

func TestScanner_RepoCacheDir(t *testing.T) {
	cacheDir := t.TempDir()

	ref := config.RepoRef{Owner: "org", Repo: "repo", Ref: "main"}
	got := repoCacheDir(cacheDir, ref)
	want := filepath.Join(cacheDir, "org_repo_main")
	if got != want {
		t.Errorf("repoCacheDir = %q, want %q", got, want)
	}
}

func TestScanner_RepoCacheDir_NoRef(t *testing.T) {
	cacheDir := t.TempDir()

	ref := config.RepoRef{Owner: "org", Repo: "repo"}
	got := repoCacheDir(cacheDir, ref)
	want := filepath.Join(cacheDir, "org_repo")
	if got != want {
		t.Errorf("repoCacheDir = %q, want %q", got, want)
	}
}

func TestParseFrontmatter_NoClosingDashes(t *testing.T) {
	// --- at start but no closing ---, treat as no frontmatter
	data := []byte("---\nname: Test\nno closing dashes\n")
	name, desc, body := parseFrontmatter(data)
	if name != "" {
		t.Errorf("name = %q, want empty", name)
	}
	if desc != "" {
		t.Errorf("desc = %q, want empty", desc)
	}
	if body != string(data) {
		t.Errorf("body should be original content")
	}
}

func TestParseFrontmatter_ExtraFieldsIgnored(t *testing.T) {
	data := []byte("---\nname: Test\ndescription: Desc\ntags: [a, b]\nversion: 1.0\n---\nBody\n")
	name, desc, body := parseFrontmatter(data)
	if name != "Test" {
		t.Errorf("name = %q", name)
	}
	if desc != "Desc" {
		t.Errorf("desc = %q", desc)
	}
	if body != "Body\n" {
		t.Errorf("body = %q", body)
	}
}

func TestScanDir_SkillNoFrontmatter(t *testing.T) {
	dir := t.TempDir()
	createTestSkillRaw(t, dir, "plain", "Just plain markdown\nNo frontmatter\n")

	skills := scanDir(dir, "local")
	if len(skills) != 1 {
		t.Fatalf("got %d skills", len(skills))
	}
	// Name should fall back to dirname
	if skills[0].Name != "plain" {
		t.Errorf("name = %q, want %q", skills[0].Name, "plain")
	}
	if skills[0].Content != "Just plain markdown\nNo frontmatter\n" {
		t.Errorf("content = %q", skills[0].Content)
	}
}

func TestSyncRepo_MultipleSkills(t *testing.T) {
	ts := newGitHubTestServer(t, func(w http.ResponseWriter, r *http.Request) {
		path := r.URL.Path
		switch {
		case strings.HasSuffix(path, "/contents/skills"):
			json.NewEncoder(w).Encode([]github.ContentEntry{
				{Name: "docker", Path: "skills/docker", Type: "dir"},
				{Name: "go", Path: "skills/go", Type: "dir"},
				{Name: "readme.md", Path: "skills/readme.md", Type: "file"}, // should be skipped
			})
		case strings.HasSuffix(path, "/contents/skills/docker/SKILL.md"):
			w.Write([]byte("---\nname: Docker\ndescription: d\n---\nbody\n"))
		case strings.HasSuffix(path, "/contents/skills/go/SKILL.md"):
			w.Write([]byte("---\nname: Go\ndescription: g\n---\ngo body\n"))
		case strings.Contains(path, "/references"):
			w.WriteHeader(404)
		default:
			w.WriteHeader(404)
		}
	})

	cacheDir := t.TempDir()
	cfg := newTestConfig(nil, []string{"owner/repo"}, cacheDir)
	s := newScannerWithGH(cfg, ts.URL, "")

	ref := config.ParseRepoRef("owner/repo")
	if err := s.syncRepo(ref); err != nil {
		t.Fatalf("syncRepo: %v", err)
	}

	skills := scanDir(repoCacheDir(cacheDir, ref), "owner/repo")
	if len(skills) != 2 {
		t.Fatalf("got %d skills, want 2", len(skills))
	}
}

func TestScanner_StartStop_WithRepos(t *testing.T) {
	ts := newGitHubTestServer(t, func(w http.ResponseWriter, r *http.Request) {
		if strings.HasSuffix(r.URL.Path, "/contents/skills") {
			json.NewEncoder(w).Encode([]github.ContentEntry{})
		} else {
			w.WriteHeader(404)
		}
	})

	cacheDir := t.TempDir()
	cfg := newTestConfig(nil, []string{"owner/repo"}, cacheDir)
	cfg.Cache.SyncInterval = time.Hour
	s := newScannerWithGH(cfg, ts.URL, "")

	ctx := context.Background()
	s.Start(ctx)

	done := make(chan struct{})
	go func() {
		s.Stop()
		close(done)
	}()
	select {
	case <-done:
	case <-time.After(5 * time.Second):
		t.Fatal("Stop did not return within 5 seconds")
	}
}

func TestParseFrontmatter_MultilineStopsAtNonIndented(t *testing.T) {
	// The multiline parser should stop collecting when it hits a non-indented,
	// non-empty line. This tests the condition: len(l) > 0 && l[0] != ' ' && l[0] != '\t'
	data := []byte("---\nname: Test\ndescription: |\n  first line\n  second line\nname: Override\n---\nBody\n")
	name, desc, body := parseFrontmatter(data)
	if !strings.Contains(desc, "first line") || !strings.Contains(desc, "second line") {
		t.Errorf("description = %q, should contain both indented lines", desc)
	}
	// "name: Override" is non-indented → should have broken the multiline block
	if strings.Contains(desc, "Override") {
		t.Errorf("description = %q, should NOT contain the non-indented line", desc)
	}
	// The re-parsed name key after multiline block isn't re-read by our simple parser,
	// so name stays "Test"
	if name != "Override" && name != "Test" {
		t.Errorf("name = %q", name)
	}
	if body != "Body\n" {
		t.Errorf("body = %q", body)
	}
}

func TestParseFrontmatter_MultilineIncludesEmptyLines(t *testing.T) {
	// Empty lines (len == 0) should NOT break multiline collection
	data := []byte("---\nname: Test\ndescription: |\n  first\n\n  third\n---\nBody\n")
	_, desc, _ := parseFrontmatter(data)
	if !strings.Contains(desc, "first") {
		t.Errorf("description = %q, should contain 'first'", desc)
	}
	if !strings.Contains(desc, "third") {
		t.Errorf("description = %q, should contain 'third'", desc)
	}
}

func TestParseFrontmatter_MultilineTabIndent(t *testing.T) {
	data := []byte("---\nname: TabTest\ndescription: |\n\ttab indented\n\tsecond tab\n---\nBody\n")
	_, desc, _ := parseFrontmatter(data)
	if !strings.Contains(desc, "tab indented") || !strings.Contains(desc, "second tab") {
		t.Errorf("description = %q, should contain tab-indented lines", desc)
	}
}

func TestSyncAllRepos_ContinuesOnError(t *testing.T) {
	ts := newGitHubTestServer(t, func(w http.ResponseWriter, r *http.Request) {
		path := r.URL.Path
		switch {
		case strings.Contains(path, "fail-repo"):
			w.WriteHeader(500)
		case strings.HasSuffix(path, "ok-repo/contents/skills"):
			json.NewEncoder(w).Encode([]github.ContentEntry{
				{Name: "my-skill", Path: "skills/my-skill", Type: "dir"},
			})
		case strings.HasSuffix(path, "ok-repo/contents/skills/my-skill/SKILL.md"):
			w.Write([]byte("---\nname: Synced\ndescription: From OK repo\n---\nBody"))
		case strings.HasSuffix(path, "ok-repo/contents/skills/my-skill/references"):
			w.WriteHeader(404)
		default:
			w.WriteHeader(404)
		}
	})

	cacheDir := t.TempDir()
	cfg := newTestConfig(nil, []string{"owner/fail-repo", "owner/ok-repo"}, cacheDir)
	s := newScannerWithGH(cfg, ts.URL, "")

	// Trigger syncAllRepos
	s.syncAllRepos()

	// Verify the second repo was still synced despite the first failing
	skills := s.List()
	found := false
	for _, sk := range skills {
		if strings.Contains(sk.Name, "Synced") || strings.Contains(sk.Source, "ok-repo") {
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

func TestScanner_ListEmpty(t *testing.T) {
	cfg := newTestConfig(nil, nil, t.TempDir())
	s := newScanner(cfg)
	skills := s.List()
	if len(skills) != 0 {
		t.Errorf("expected 0 skills, got %d", len(skills))
	}
}

func TestScanner_GetNotFoundEmpty(t *testing.T) {
	cfg := newTestConfig(nil, nil, t.TempDir())
	s := newScanner(cfg)
	_, found := s.Get("skills://nonexistent/foo")
	if found {
		t.Error("expected not found on empty scanner")
	}
}

func TestScanner_GetFirstMatch(t *testing.T) {
	dir := t.TempDir()
	createTestSkill(t, dir, "skill-a", "Alpha", "First skill", "Body A")
	createTestSkill(t, dir, "skill-b", "Beta", "Second skill", "Body B")

	cfg := newTestConfig([]string{dir}, nil, t.TempDir())
	s := newScanner(cfg)

	sk, found := s.Get(fmt.Sprintf("skills://%s/Alpha", filepath.Base(dir)))
	if !found {
		t.Fatal("expected to find Alpha")
	}
	if sk.Name != "Alpha" {
		t.Errorf("Name = %q, want Alpha", sk.Name)
	}
}

func TestSyncRepo_SkillMdFetchFails(t *testing.T) {
	ts := newGitHubTestServer(t, func(w http.ResponseWriter, r *http.Request) {
		path := r.URL.Path
		switch {
		case strings.HasSuffix(path, "/contents/skills"):
			json.NewEncoder(w).Encode([]github.ContentEntry{
				{Name: "good-skill", Path: "skills/good-skill", Type: "dir"},
				{Name: "bad-skill", Path: "skills/bad-skill", Type: "dir"},
			})
		case strings.HasSuffix(path, "/contents/skills/good-skill/SKILL.md"):
			w.Write([]byte("---\nname: Good\ndescription: works\n---\nBody"))
		case strings.HasSuffix(path, "/contents/skills/bad-skill/SKILL.md"):
			w.WriteHeader(500) // Fails
		case strings.HasSuffix(path, "/contents/skills/good-skill/references"):
			w.WriteHeader(404)
		default:
			w.WriteHeader(404)
		}
	})

	cacheDir := t.TempDir()
	cfg := newTestConfig(nil, []string{"o/r"}, cacheDir)
	s := newScannerWithGH(cfg, ts.URL, "")

	ref := config.ParseRepoRef("o/r")
	err := s.syncRepo(ref)
	if err != nil {
		t.Fatalf("syncRepo: %v", err)
	}

	// good-skill should be cached, bad-skill should not
	skills := scanDir(repoCacheDir(cacheDir, ref), "o/r")
	if len(skills) != 1 {
		t.Fatalf("got %d skills, want 1", len(skills))
	}
	if skills[0].Name != "Good" {
		t.Errorf("name = %q, want Good", skills[0].Name)
	}
}

func TestSyncRepo_FileInSkillsDir(t *testing.T) {
	// Files directly in skills/ should be skipped (only dirs are skills)
	ts := newGitHubTestServer(t, func(w http.ResponseWriter, r *http.Request) {
		path := r.URL.Path
		switch {
		case strings.HasSuffix(path, "/contents/skills"):
			json.NewEncoder(w).Encode([]github.ContentEntry{
				{Name: "README.md", Path: "skills/README.md", Type: "file"},
				{Name: "my-skill", Path: "skills/my-skill", Type: "dir"},
			})
		case strings.HasSuffix(path, "/contents/skills/my-skill/SKILL.md"):
			w.Write([]byte("---\nname: Test\n---\nBody"))
		case strings.HasSuffix(path, "/contents/skills/my-skill/references"):
			w.WriteHeader(404)
		default:
			w.WriteHeader(404)
		}
	})

	cacheDir := t.TempDir()
	cfg := newTestConfig(nil, []string{"o/r"}, cacheDir)
	s := newScannerWithGH(cfg, ts.URL, "")

	ref := config.ParseRepoRef("o/r")
	err := s.syncRepo(ref)
	if err != nil {
		t.Fatalf("syncRepo: %v", err)
	}

	skills := scanDir(repoCacheDir(cacheDir, ref), "o/r")
	if len(skills) != 1 {
		t.Fatalf("got %d skills, want 1", len(skills))
	}
}

func TestSyncRepo_DirInReferences(t *testing.T) {
	// Dirs in references/ should be skipped
	ts := newGitHubTestServer(t, func(w http.ResponseWriter, r *http.Request) {
		path := r.URL.Path
		switch {
		case strings.HasSuffix(path, "/contents/skills"):
			json.NewEncoder(w).Encode([]github.ContentEntry{
				{Name: "my-skill", Path: "skills/my-skill", Type: "dir"},
			})
		case strings.HasSuffix(path, "/contents/skills/my-skill/SKILL.md"):
			w.Write([]byte("---\nname: RefTest\n---\nBody"))
		case strings.HasSuffix(path, "/contents/skills/my-skill/references"):
			json.NewEncoder(w).Encode([]github.ContentEntry{
				{Name: "subdir", Path: "skills/my-skill/references/subdir", Type: "dir"},
				{Name: "doc.md", Path: "skills/my-skill/references/doc.md", Type: "file"},
			})
		case strings.HasSuffix(path, "/contents/skills/my-skill/references/doc.md"):
			w.Write([]byte("Reference content"))
		default:
			w.WriteHeader(404)
		}
	})

	cacheDir := t.TempDir()
	cfg := newTestConfig(nil, []string{"o/r"}, cacheDir)
	s := newScannerWithGH(cfg, ts.URL, "")

	ref := config.ParseRepoRef("o/r")
	s.syncRepo(ref)

	skills := scanDir(repoCacheDir(cacheDir, ref), "o/r")
	if len(skills) != 1 {
		t.Fatalf("got %d skills, want 1", len(skills))
	}
	if len(skills[0].References) != 1 {
		t.Fatalf("got %d refs, want 1", len(skills[0].References))
	}
	if skills[0].References[0].Name != "doc.md" {
		t.Errorf("ref name = %q, want doc.md", skills[0].References[0].Name)
	}
}

func TestSyncRepo_ReferenceFetchFails(t *testing.T) {
	ts := newGitHubTestServer(t, func(w http.ResponseWriter, r *http.Request) {
		path := r.URL.Path
		switch {
		case strings.HasSuffix(path, "/contents/skills"):
			json.NewEncoder(w).Encode([]github.ContentEntry{
				{Name: "my-skill", Path: "skills/my-skill", Type: "dir"},
			})
		case strings.HasSuffix(path, "/contents/skills/my-skill/SKILL.md"):
			w.Write([]byte("---\nname: RefFail\n---\nBody"))
		case strings.HasSuffix(path, "/contents/skills/my-skill/references"):
			json.NewEncoder(w).Encode([]github.ContentEntry{
				{Name: "good.md", Path: "skills/my-skill/references/good.md", Type: "file"},
				{Name: "bad.md", Path: "skills/my-skill/references/bad.md", Type: "file"},
			})
		case strings.HasSuffix(path, "/contents/skills/my-skill/references/good.md"):
			w.Write([]byte("Good ref"))
		case strings.HasSuffix(path, "/contents/skills/my-skill/references/bad.md"):
			w.WriteHeader(500)
		default:
			w.WriteHeader(404)
		}
	})

	cacheDir := t.TempDir()
	cfg := newTestConfig(nil, []string{"o/r"}, cacheDir)
	s := newScannerWithGH(cfg, ts.URL, "")

	ref := config.ParseRepoRef("o/r")
	s.syncRepo(ref)

	skills := scanDir(repoCacheDir(cacheDir, ref), "o/r")
	if len(skills) != 1 {
		t.Fatalf("got %d skills, want 1", len(skills))
	}
	// Only good.md should be cached
	if len(skills[0].References) != 1 {
		t.Fatalf("got %d refs, want 1 (bad.md failed)", len(skills[0].References))
	}
}

func TestScanDir_UnreadableSkillMd(t *testing.T) {
	dir := t.TempDir()
	skillDir := filepath.Join(dir, "my-skill")
	os.MkdirAll(skillDir, 0o755)
	os.WriteFile(filepath.Join(skillDir, "SKILL.md"), []byte("content"), 0o000)

	skills := scanDir(dir, "test")
	// Unreadable SKILL.md should be skipped
	if len(skills) != 0 {
		t.Errorf("got %d skills, want 0", len(skills))
	}
}

func TestScanDir_UnreadableReference(t *testing.T) {
	dir := t.TempDir()
	createTestSkill(t, dir, "my-skill", "Test", "desc", "body")
	refDir := filepath.Join(dir, "my-skill", "references")
	os.MkdirAll(refDir, 0o755)
	os.WriteFile(filepath.Join(refDir, "good.md"), []byte("good"), 0o644)
	os.WriteFile(filepath.Join(refDir, "bad.md"), []byte("bad"), 0o000)

	skills := scanDir(dir, "test")
	if len(skills) != 1 {
		t.Fatalf("got %d skills, want 1", len(skills))
	}
	// Only good.md should be readable
	if len(skills[0].References) != 1 {
		t.Errorf("got %d refs, want 1 (bad.md unreadable)", len(skills[0].References))
	}
}

func TestParseFrontmatter_OnlyDashes(t *testing.T) {
	data := []byte("---\n---\nBody\n")
	name, desc, body := parseFrontmatter(data)
	if name != "" || desc != "" {
		t.Errorf("name=%q desc=%q, want empty", name, desc)
	}
	if body != "Body\n" {
		t.Errorf("body = %q", body)
	}
}

func TestParseFrontmatter_UnicodeContent(t *testing.T) {
	data := []byte("---\nname: 日本語スキル\ndescription: スキルの説明\n---\n本文\n")
	name, desc, body := parseFrontmatter(data)
	if name != "日本語スキル" {
		t.Errorf("name = %q, want Japanese", name)
	}
	if desc != "スキルの説明" {
		t.Errorf("desc = %q, want Japanese", desc)
	}
	if body != "本文\n" {
		t.Errorf("body = %q", body)
	}
}

func TestParseFrontmatter_LargeBody(t *testing.T) {
	large := strings.Repeat("Content line\n", 10000)
	data := []byte("---\nname: Big\n---\n" + large)
	name, _, body := parseFrontmatter(data)
	if name != "Big" {
		t.Errorf("name = %q", name)
	}
	if body != large {
		t.Errorf("body length = %d, want %d", len(body), len(large))
	}
}

func TestParseFrontmatter_EmptyInput(t *testing.T) {
	name, desc, body := parseFrontmatter([]byte(""))
	if name != "" || desc != "" || body != "" {
		t.Errorf("expected all empty, got name=%q desc=%q body=%q", name, desc, body)
	}
}

func TestParseFrontmatter_DashesOnlyNoNewline(t *testing.T) {
	name, desc, body := parseFrontmatter([]byte("---"))
	if name != "" || desc != "" {
		t.Errorf("name=%q desc=%q, want empty", name, desc)
	}
	if body != "---" {
		t.Errorf("body = %q, want raw content returned", body)
	}
}

func TestScanner_StartStopIdempotent(t *testing.T) {
	cfg := newTestConfig(nil, nil, t.TempDir())
	s := newScanner(cfg)
	ctx := context.Background()
	s.Start(ctx)
	s.Stop()
	s.Stop() // idempotent
}

func TestScanner_ForceSyncMultiple(t *testing.T) {
	dir := t.TempDir()
	createTestSkill(t, dir, "my-skill", "Test", "desc", "body")

	cfg := newTestConfig([]string{dir}, nil, t.TempDir())
	s := newScanner(cfg)

	s.ForceSync()
	skills := s.List()
	if len(skills) != 1 {
		t.Fatalf("got %d skills after first sync", len(skills))
	}

	s.ForceSync()
	skills = s.List()
	if len(skills) != 1 {
		t.Fatalf("got %d skills after second sync", len(skills))
	}
}

func TestScanDir_EmptySkillMd(t *testing.T) {
	dir := t.TempDir()
	skillDir := filepath.Join(dir, "empty-skill")
	os.MkdirAll(skillDir, 0o755)
	os.WriteFile(filepath.Join(skillDir, "SKILL.md"), []byte(""), 0o644)

	skills := scanDir(dir, "test")
	if len(skills) != 1 {
		t.Fatalf("got %d skills, want 1", len(skills))
	}
	// Empty file means no frontmatter, name falls back to dir name
	if skills[0].Name != "empty-skill" {
		t.Errorf("name = %q, want dirname fallback", skills[0].Name)
	}
}

func TestScanDir_ManySkills(t *testing.T) {
	dir := t.TempDir()
	for i := 0; i < 50; i++ {
		name := fmt.Sprintf("skill-%03d", i)
		createTestSkill(t, dir, name, name, "Description", "Body")
	}

	skills := scanDir(dir, "test")
	if len(skills) != 50 {
		t.Errorf("got %d skills, want 50", len(skills))
	}
}

// ---------------------------------------------------------------------------
// .github/ subdirectory scanning
// ---------------------------------------------------------------------------

func TestScanner_List_GithubSubdir(t *testing.T) {
// Skills placed under dir/.github/<skill-name>/SKILL.md are discovered.
dir := t.TempDir()
createTestSkill(t, filepath.Join(dir, ".github"), "myskill", "My Skill", "A skill in .github", "Skill body")

cfg := newTestConfig([]string{dir}, nil, t.TempDir())
s := newScanner(cfg)
skills := s.List()
if len(skills) != 1 {
t.Fatalf("got %d skills, want 1", len(skills))
}
if skills[0].Name != "My Skill" {
t.Errorf("Name = %q", skills[0].Name)
}
}

func TestScanner_List_GithubAndRootSkills(t *testing.T) {
// Skills at root level and in .github/ are both included.
dir := t.TempDir()
createTestSkill(t, dir, "rootskill", "Root Skill", "Root level", "Body")
createTestSkill(t, filepath.Join(dir, ".github"), "githubskill", "Github Skill", "In .github", "Body")

cfg := newTestConfig([]string{dir}, nil, t.TempDir())
s := newScanner(cfg)
skills := s.List()
if len(skills) != 2 {
t.Fatalf("got %d skills, want 2", len(skills))
}
}

func TestScanner_List_DefaultCWD(t *testing.T) {
// With no dirs configured, scanner uses CWD. In test environment (package dir),
// there are no SKILL.md subdirectories, so result is empty — but it must not panic.
cfg := newTestConfig(nil, nil, t.TempDir())
s := newScanner(cfg)
_ = s.List() // must not panic
}
