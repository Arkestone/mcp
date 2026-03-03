// Package loader provides on-demand access to Copilot custom instruction files
// from local directories and GitHub repositories.
package loader

import (
	"bytes"
	"context"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/adrg/frontmatter"

	"github.com/Arkestone/mcp/pkg/cache"
	"github.com/Arkestone/mcp/pkg/config"
	"github.com/Arkestone/mcp/pkg/glob"
	"github.com/Arkestone/mcp/pkg/github"
	"github.com/Arkestone/mcp/pkg/syncer"
)

const maxFileSize = 1 << 20 // 1 MiB

// Instruction represents a single instruction file.
type Instruction struct {
	Source  string
	Path    string
	Content string
	URI     string
	ApplyTo []string // VS Code glob patterns; empty = applies everywhere
}

// Loader provides on-demand access to instruction files.
type Loader struct {
	cfg    *config.Config
	gh     *github.Client
	syncer *syncer.Syncer
	cache  cache.List[Instruction]
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

// List returns all instructions, cached for cache.DefaultTTL.
func (l *Loader) List() []Instruction {
	return l.cache.Get(l.scan)
}

func (l *Loader) scan() []Instruction {
	dirs := l.cfg.Sources.Dirs
	if len(dirs) == 0 {
		dirs = []string{"."}
	}
	var out []Instruction
	for _, dir := range dirs {
		out = append(out, scanDir(dir, sourceFor(dir))...)
	}
	for _, ref := range l.cfg.ParsedRepos() {
		cacheDir := repoCacheDir(l.cfg.Cache.Dir, ref)
		out = append(out, scanDir(cacheDir, ref.Owner+"/"+ref.Repo)...)
	}
	return out
}

func (l *Loader) Get(uri string) (Instruction, bool) {
	for _, inst := range l.List() {
		if inst.URI == uri {
			return inst, true
		}
	}
	return Instruction{}, false
}

// FilterByFilePath returns instructions applicable to the given file path.
// Instructions with an empty ApplyTo are always included (global rules).
// Reproducible: the order of returned instructions is stable (same as List order).
// If filePath is empty, all instructions are returned unchanged.
func FilterByFilePath(instructions []Instruction, filePath string) []Instruction {
	if filePath == "" {
		return instructions
	}
	fp := filepath.ToSlash(filePath)
	var out []Instruction
	for _, inst := range instructions {
		if len(inst.ApplyTo) == 0 || glob.MatchAny(inst.ApplyTo, fp) {
			out = append(out, inst)
		}
	}
	return out
}

var skipDirs = map[string]bool{
	".git":         true,
	"node_modules": true,
}

// scanDir walks the entire directory tree and collects *.instructions.md and
// copilot-instructions.md. .github/ is visited first (lexically) so it takes
// URI deduplication priority over root-level files with the same stem name.
func scanDir(dir, source string) []Instruction {
	seen := map[string]bool{}
	var out []Instruction

	add := func(inst Instruction) {
		if !seen[inst.URI] {
			seen[inst.URI] = true
			out = append(out, inst)
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

		if name == "copilot-instructions.md" {
			content, err := os.ReadFile(fpath)
			if err != nil {
				return nil
			}
			applyTo, body := parseFrontmatter(content)
			uriName := strings.TrimSuffix(rel, ".md")
			if rel == ".github/copilot-instructions.md" || rel == "copilot-instructions.md" {
				uriName = "copilot-instructions"
			}
			add(Instruction{
				Source: source, Path: rel, Content: body, ApplyTo: applyTo,
				URI: fmt.Sprintf("instructions://%s/%s", source, uriName),
			})
			return nil
		}

		if strings.HasSuffix(name, ".instructions.md") {
			content, err := os.ReadFile(fpath)
			if err != nil {
				return nil
			}
			applyTo, body := parseFrontmatter(content)
			instrName := strings.TrimSuffix(name, ".instructions.md")
			add(Instruction{
				Source: source, Path: rel, Content: body, ApplyTo: applyTo,
				URI: fmt.Sprintf("instructions://%s/%s", source, instrName),
			})
			return nil
		}

		return nil
	})
	return out
}

type instructionMeta struct {
	ApplyTo interface{} `yaml:"applyTo"`
}

func parseFrontmatter(data []byte) (applyTo []string, body string) {
	var meta instructionMeta
	rest, _ := frontmatter.Parse(bytes.NewReader(data), &meta)
	return toStringSlice(meta.ApplyTo), string(rest)
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
		name := filepath.Base(entry.Path)
		if name != "copilot-instructions.md" && !strings.HasSuffix(name, ".instructions.md") {
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
