// Package loader provides on-demand access to Copilot custom instruction files
// from local directories and GitHub repositories.
package loader

import (
	"context"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/Arkestone/mcp/pkg/config"
	"github.com/Arkestone/mcp/pkg/github"
	"github.com/Arkestone/mcp/pkg/syncer"
)

// Instruction represents a single instruction file.
type Instruction struct {
	Source  string // origin identifier (directory basename or owner/repo)
	Path    string // relative path within the source
	Content string // raw markdown content
	URI     string // MCP resource URI
}

// Loader provides on-demand access to instruction files.
type Loader struct {
	cfg    *config.Config
	gh     *github.Client
	syncer *syncer.Syncer
}

// New creates a Loader with its background syncer.
func New(cfg *config.Config, gh *github.Client) *Loader {
	l := &Loader{cfg: cfg, gh: gh}
	l.syncer = syncer.New(cfg.Cache.SyncInterval, l.syncAllRepos)
	return l
}

// Start begins background sync. Stop must be called to shut down.
func (l *Loader) Start(ctx context.Context) { l.syncer.Start(ctx) }

// Stop shuts down the background sync.
func (l *Loader) Stop() { l.syncer.Stop() }

// ForceSync triggers an immediate sync of all remote repos.
func (l *Loader) ForceSync() { l.syncer.ForceSync() }

// List returns all instructions. Local dirs are read from disk; repos from cache.
func (l *Loader) List() []Instruction {
	var out []Instruction
	for _, dir := range l.cfg.Sources.Dirs {
		out = append(out, scanDir(dir, filepath.Base(dir))...)
	}
	for _, ref := range l.cfg.ParsedRepos() {
		cacheDir := repoCacheDir(l.cfg.Cache.Dir, ref)
		out = append(out, scanDir(cacheDir, ref.Owner+"/"+ref.Repo)...)
	}
	return out
}

// Get returns a single instruction by URI.
func (l *Loader) Get(uri string) (Instruction, bool) {
	for _, inst := range l.List() {
		if inst.URI == uri {
			return inst, true
		}
	}
	return Instruction{}, false
}

// scanDir reads instruction files from a directory.
// Looks for .github/copilot-instructions.md and .github/instructions/**/*.instructions.md
func scanDir(dir, source string) []Instruction {
	var out []Instruction

	// .github/copilot-instructions.md
	ciPath := filepath.Join(dir, ".github", "copilot-instructions.md")
	if content, err := os.ReadFile(ciPath); err == nil {
		out = append(out, Instruction{
			Source:  source,
			Path:    ".github/copilot-instructions.md",
			Content: string(content),
			URI:     fmt.Sprintf("instructions://%s/copilot-instructions", source),
		})
	}

	// .github/instructions/**/*.instructions.md
	instrDir := filepath.Join(dir, ".github", "instructions")
	if info, err := os.Stat(instrDir); err == nil && info.IsDir() {
		_ = filepath.Walk(instrDir, func(path string, info os.FileInfo, err error) error {
			if err != nil || info.IsDir() {
				return err
			}
			if !strings.HasSuffix(info.Name(), ".instructions.md") {
				return nil
			}
			content, err := os.ReadFile(path)
			if err != nil {
				return nil
			}
			relPath, _ := filepath.Rel(dir, path)
			name := strings.TrimSuffix(info.Name(), ".instructions.md")
			out = append(out, Instruction{
				Source:  source,
				Path:    relPath,
				Content: string(content),
				URI:     fmt.Sprintf("instructions://%s/%s", source, name),
			})
			return nil
		})
	}

	return out
}

func (l *Loader) syncAllRepos() {
	for _, ref := range l.cfg.ParsedRepos() {
		if err := l.syncRepo(ref); err != nil {
			log.Printf("sync %s/%s: %v", ref.Owner, ref.Repo, err)
		}
	}
}

func (l *Loader) syncRepo(ref config.RepoRef) error {
	cacheDir := repoCacheDir(l.cfg.Cache.Dir, ref)
	ghDir := filepath.Join(cacheDir, ".github")

	if err := os.MkdirAll(ghDir, 0o755); err != nil {
		return fmt.Errorf("creating cache dir: %w", err)
	}

	ctx := context.Background()

	// Sync .github/copilot-instructions.md
	if content, err := l.gh.FetchFile(ctx, ref.Owner, ref.Repo, ref.Ref, ".github/copilot-instructions.md"); err == nil {
		_ = os.WriteFile(filepath.Join(ghDir, "copilot-instructions.md"), []byte(content), 0o644)
	}

	// Sync .github/instructions/
	instrCacheDir := filepath.Join(ghDir, "instructions")
	entries, err := l.gh.FetchDir(ctx, ref.Owner, ref.Repo, ref.Ref, ".github/instructions")
	if err == nil {
		_ = os.MkdirAll(instrCacheDir, 0o755)
		for _, entry := range entries {
			if !strings.HasSuffix(entry.Name, ".instructions.md") {
				continue
			}
			content, err := l.gh.FetchFile(ctx, ref.Owner, ref.Repo, ref.Ref, entry.Path)
			if err != nil {
				continue
			}
			_ = os.WriteFile(filepath.Join(instrCacheDir, entry.Name), []byte(content), 0o644)
		}
	}

	return nil
}

func repoCacheDir(cacheBase string, ref config.RepoRef) string {
	key := ref.Owner + "_" + ref.Repo
	if ref.Ref != "" {
		key += "_" + ref.Ref
	}
	return filepath.Join(cacheBase, key)
}
