// Package loader provides on-demand access to Copilot prompt files
// from local directories and GitHub repositories.
package loader

import (
	"context"
	"fmt"
	"log"
	"os"
	"path"
	"path/filepath"
	"strings"

	"github.com/adrg/frontmatter"

	"github.com/Arkestone/mcp/pkg/cache"
	"github.com/Arkestone/mcp/pkg/config"
	"github.com/Arkestone/mcp/pkg/filter"
	"github.com/Arkestone/mcp/pkg/glob"
	"github.com/Arkestone/mcp/pkg/github"
	"github.com/Arkestone/mcp/pkg/syncer"
)

const (
	TypePrompt   = "prompt"
	TypeChatmode = "chatmode"

	maxFileSize = 1 << 20 // 1 MiB
)

// Prompt represents a single Copilot prompt or chat mode file.
type Prompt struct {
	Name        string
	Description string
	Mode        string
	Type        string
	Tags        []string
	Files       []string // glob patterns from frontmatter files: — restricts which file paths this prompt applies to
	Content     string
	Source      string
	Path        string
	URI         string
}

// FilterByQuery returns prompts scored and sorted by relevance to query using
// word-boundary tokenization for precision and a stable sort for reproducibility.
// Prompts scoring 0 are excluded. If query is empty, all prompts are returned unchanged.
func FilterByQuery(prompts []Prompt, query string) []Prompt {
	return filter.SortByScore(prompts, func(p Prompt) int {
		return filter.Score(query, p.Name, p.Description, p.Tags)
	})
}

// FilterByFilePath returns prompts applicable to the given file path.
// Prompts without Files patterns are always included (global scope).
// Prompts with Files patterns are included only when at least one pattern matches filePath.
// If filePath is empty, all prompts are returned unchanged.
func FilterByFilePath(prompts []Prompt, filePath string) []Prompt {
	if filePath == "" {
		return prompts
	}
	fp := filepath.ToSlash(filePath)
	var out []Prompt
	for _, p := range prompts {
		if len(p.Files) == 0 || glob.MatchAny(p.Files, fp) {
			out = append(out, p)
		}
	}
	return out
}

// Loader provides on-demand access to Copilot prompt files.
type Loader struct {
	cfg    *config.Config
	gh     *github.Client
	syncer *syncer.Syncer
	cache  cache.List[Prompt]
}

func New(cfg *config.Config, gh *github.Client) *Loader {
	l := &Loader{cfg: cfg, gh: gh}
	l.syncer = syncer.New(cfg.Cache.SyncInterval, l.syncAllRepos)
	return l
}

func (l *Loader) Start(ctx context.Context) { l.syncer.Start(ctx) }
func (l *Loader) Stop()                     { l.syncer.Stop() }

func (l *Loader) ForceSync() {
	l.cache.Invalidate()
	l.syncer.ForceSync()
}

func (l *Loader) List() []Prompt {
	return l.cache.Get(l.scan)
}

func (l *Loader) scan() []Prompt {
	dirs := l.cfg.Sources.Dirs
	if len(dirs) == 0 {
		dirs = []string{"."}
	}
	var out []Prompt
	for _, dir := range dirs {
		out = append(out, scanDir(dir, sourceFor(dir))...)
	}
	for _, ref := range l.cfg.ParsedRepos() {
		cacheDir := repoCacheDir(l.cfg.Cache.Dir, ref)
		out = append(out, scanDir(cacheDir, ref.Owner+"/"+ref.Repo)...)
	}
	return out
}

func (l *Loader) Get(uri string) (Prompt, bool) {
	for _, p := range l.List() {
		if p.URI == uri {
			return p, true
		}
	}
	return Prompt{}, false
}

var skipDirs = map[string]bool{
	".git":         true,
	"node_modules": true,
}

func scanDir(dir, source string) []Prompt {
	seen := map[string]bool{}
	var out []Prompt

	add := func(p Prompt) {
		if !seen[p.URI] {
			seen[p.URI] = true
			out = append(out, p)
		}
	}

	_ = filepath.Walk(dir, func(fpath string, info os.FileInfo, err error) error {
		if err != nil {
			return nil
		}
		if info.IsDir() {
			if skipDirs[info.Name()] {
				return filepath.SkipDir
			}
			return nil
		}
		if info.Size() > maxFileSize {
			return nil
		}

		name := info.Name()
		rel, _ := filepath.Rel(dir, fpath)
		rel = filepath.ToSlash(rel)

		if strings.HasSuffix(name, ".prompt.md") {
			content, err := os.ReadFile(fpath)
			if err != nil {
				return nil
			}
			pname := strings.TrimSuffix(name, ".prompt.md")
			desc, mode, tags, files, body := parseFrontmatter(string(content))
			if desc == "" {
				desc = pname
			}
			add(Prompt{
				Name: pname, Description: desc, Mode: mode, Tags: tags, Files: files,
				Type: TypePrompt, Content: body, Source: source, Path: rel,
				URI: fmt.Sprintf("prompts://%s/%s", source, pname),
			})
			return nil
		}

		if strings.HasSuffix(name, ".chatmode.md") {
			content, err := os.ReadFile(fpath)
			if err != nil {
				return nil
			}
			pname := strings.TrimSuffix(name, ".chatmode.md")
			desc, _, tags, files, body := parseFrontmatter(string(content))
			if desc == "" {
				desc = pname
			}
			add(Prompt{
				Name: pname, Description: desc, Tags: tags, Files: files,
				Type: TypeChatmode, Content: body, Source: source, Path: rel,
				URI: fmt.Sprintf("prompts://%s/%s", source, pname),
			})
			return nil
		}

		return nil
	})
	return out
}

type promptMeta struct {
	Description string      `yaml:"description"`
	Mode        string      `yaml:"mode"`
	Tags        interface{} `yaml:"tags"`
	Files       interface{} `yaml:"files"` // string or []string glob patterns
}

func parseFrontmatter(content string) (description, mode string, tags, files []string, body string) {
	var meta promptMeta
	rest, _ := frontmatter.Parse(strings.NewReader(content), &meta)
	return meta.Description, meta.Mode, toStringSlice(meta.Tags), toStringSlice(meta.Files), string(rest)
}

func toStringSlice(v interface{}) []string {
	switch t := v.(type) {
	case string:
		if t == "" {
			return nil
		}
		// Support comma-separated patterns (e.g. "**/*.ts,**/*.tsx")
		// but only split on commas that are NOT inside brace expansions like {ts,tsx}.
		return splitGlobPatterns(t)
	case []interface{}:
		var out []string
		for _, item := range t {
			if s, ok := item.(string); ok {
				out = append(out, s)
			}
		}
		return out
	}
	return nil
}

// splitGlobPatterns splits a comma-separated glob pattern string into individual
// patterns, ignoring commas that appear inside brace expansions (e.g. {ts,tsx}).
func splitGlobPatterns(s string) []string {
	var out []string
	depth := 0
	start := 0
	for i := 0; i < len(s); i++ {
		switch s[i] {
		case '{':
			depth++
		case '}':
			if depth > 0 {
				depth--
			}
		case ',':
			if depth == 0 {
				if p := strings.TrimSpace(s[start:i]); p != "" {
					out = append(out, p)
				}
				start = i + 1
			}
		}
	}
	if p := strings.TrimSpace(s[start:]); p != "" {
		out = append(out, p)
	}
	return out
}

func (l *Loader) syncAllRepos() {
	l.cache.Invalidate()
	for _, ref := range l.cfg.ParsedRepos() {
		if err := l.syncRepo(ref); err != nil {
			log.Printf("sync %s/%s: %v", ref.Owner, ref.Repo, err)
		}
	}
}

func (l *Loader) syncRepo(ref config.RepoRef) error {
	cacheDir := repoCacheDir(l.cfg.Cache.Dir, ref)
	if err := os.MkdirAll(cacheDir, 0o755); err != nil {
		return fmt.Errorf("creating cache dir: %w", err)
	}
	ctx := context.Background()
	entries, err := l.gh.FetchDirRecursive(ctx, ref.Owner, ref.Repo, ref.Ref, "")
	if err != nil {
		if github.IsRateLimitError(err) || (l.gh.Token == "" && isAuthError(err)) {
			log.Printf("sync %s/%s: API failed (%v), falling back to ZIP download", ref.Owner, ref.Repo, err)
			return l.gh.FetchZipAndExtract(ctx, ref.Owner, ref.Repo, ref.Ref, cacheDir)
		}
		return fmt.Errorf("listing repo: %w", err)
	}
	for _, entry := range entries {
		base := path.Base(entry.Path)
		if !strings.HasSuffix(base, ".prompt.md") && !strings.HasSuffix(base, ".chatmode.md") {
			continue
		}
		content, err := l.gh.FetchFile(ctx, ref.Owner, ref.Repo, ref.Ref, entry.Path)
		if err != nil {
			continue
		}
		localPath := filepath.Join(cacheDir, filepath.FromSlash(entry.Path))
		if err := os.MkdirAll(filepath.Dir(localPath), 0o755); err != nil {
			continue
		}
		_ = os.WriteFile(localPath, []byte(content), 0o644)
	}
	return nil
}

// isAuthError reports whether err is a GitHub authentication / access error.
func isAuthError(err error) bool {
	if err == nil {
		return false
	}
	s := err.Error()
	return strings.Contains(s, "HTTP 401") || strings.Contains(s, "HTTP 403") || strings.Contains(s, "HTTP 404")
}

func repoCacheDir(cacheBase string, ref config.RepoRef) string {
	key := ref.Owner + "_" + ref.Repo
	if ref.Ref != "" {
		key += "_" + ref.Ref
	}
	return filepath.Join(cacheBase, key)
}

func sourceFor(dir string) string {
	abs, err := filepath.Abs(dir)
	if err != nil {
		return filepath.Base(dir)
	}
	return filepath.Base(abs)
}
