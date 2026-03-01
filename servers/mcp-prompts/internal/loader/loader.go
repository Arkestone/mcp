// Package loader provides on-demand access to Copilot prompt files
// from local directories and GitHub repositories.
package loader

import (
	"context"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/adrg/frontmatter"

	"github.com/Arkestone/mcp/pkg/config"
	"github.com/Arkestone/mcp/pkg/github"
	"github.com/Arkestone/mcp/pkg/syncer"
)

const (
	TypePrompt   = "prompt"
	TypeChatmode = "chatmode"
)

// Prompt represents a single Copilot prompt or chat mode file.
type Prompt struct {
	Name        string // from frontmatter or filename
	Description string // from frontmatter
	Mode        string // ask|edit|agent (prompt files only, from frontmatter)
	Type        string // "prompt" or "chatmode"
	Content     string // markdown body (after frontmatter)
	Source      string // origin (directory basename or owner/repo)
	Path        string // relative path within the source
	URI         string // MCP resource URI: prompts://{source}/{name}
}

// Loader provides on-demand access to Copilot prompt files.
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

// List returns all prompts and chat modes. Local dirs are read from disk; repos from cache.
// When no dirs are configured, the current working directory is used as default
// so the server works out-of-the-box when run from a repository root.
func (l *Loader) List() []Prompt {
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

// Get returns a single prompt by URI.
func (l *Loader) Get(uri string) (Prompt, bool) {
	for _, p := range l.List() {
		if p.URI == uri {
			return p, true
		}
	}
	return Prompt{}, false
}

// scanDir reads prompt and chatmode files from a directory.
// It searches in both the .github/ subdirectory (canonical location) and the
// root of dir (alternative location), preferring .github/ when both exist.
func scanDir(dir, source string) []Prompt {
	seen := map[string]bool{}
	var out []Prompt

	add := func(p Prompt) {
		if !seen[p.URI] {
			seen[p.URI] = true
			out = append(out, p)
		}
	}

	// Search both .github/ prefix and root level; .github/ takes priority.
	for _, base := range []string{filepath.Join(dir, ".github"), dir} {
		// prompts/*.prompt.md
		promptsDir := filepath.Join(base, "prompts")
		if info, err := os.Stat(promptsDir); err == nil && info.IsDir() {
			entries, _ := os.ReadDir(promptsDir)
			for _, e := range entries {
				if e.IsDir() || !strings.HasSuffix(e.Name(), ".prompt.md") {
					continue
				}
				content, err := os.ReadFile(filepath.Join(promptsDir, e.Name()))
				if err != nil {
					continue
				}
				name := strings.TrimSuffix(e.Name(), ".prompt.md")
				desc, mode, body := parseFrontmatter(string(content))
				if desc == "" {
					desc = name
				}
				relPath, _ := filepath.Rel(dir, filepath.Join(promptsDir, e.Name()))
				add(Prompt{
					Name:        name,
					Description: desc,
					Mode:        mode,
					Type:        TypePrompt,
					Content:     body,
					Source:      source,
					Path:        relPath,
					URI:         fmt.Sprintf("prompts://%s/%s", source, name),
				})
			}
		}

		// chatmodes/*.chatmode.md
		chatmodesDir := filepath.Join(base, "chatmodes")
		if info, err := os.Stat(chatmodesDir); err == nil && info.IsDir() {
			entries, _ := os.ReadDir(chatmodesDir)
			for _, e := range entries {
				if e.IsDir() || !strings.HasSuffix(e.Name(), ".chatmode.md") {
					continue
				}
				content, err := os.ReadFile(filepath.Join(chatmodesDir, e.Name()))
				if err != nil {
					continue
				}
				name := strings.TrimSuffix(e.Name(), ".chatmode.md")
				desc, _, body := parseFrontmatter(string(content))
				if desc == "" {
					desc = name
				}
				relPath, _ := filepath.Rel(dir, filepath.Join(chatmodesDir, e.Name()))
				add(Prompt{
					Name:        name,
					Description: desc,
					Type:        TypeChatmode,
					Content:     body,
					Source:      source,
					Path:        relPath,
					URI:         fmt.Sprintf("prompts://%s/%s", source, name),
				})
			}
		}
	}

	return out
}

type promptMeta struct {
	Description string `yaml:"description"`
	Mode        string `yaml:"mode"`
}

// parseFrontmatter extracts description and mode from YAML frontmatter.
// Returns description, mode, and the remaining body.
func parseFrontmatter(content string) (description, mode, body string) {
	var meta promptMeta
	rest, _ := frontmatter.Parse(strings.NewReader(content), &meta)
	return meta.Description, meta.Mode, string(rest)
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

	for _, subdir := range []string{"prompts", "chatmodes"} {
		remotePath := ".github/" + subdir
		entries, err := l.gh.FetchDir(ctx, ref.Owner, ref.Repo, ref.Ref, remotePath)
		if err != nil {
			continue
		}
		localDir := filepath.Join(ghDir, subdir)
		if err := os.MkdirAll(localDir, 0o755); err != nil {
			continue
		}
		for _, entry := range entries {
			suffix := ".prompt.md"
			if subdir == "chatmodes" {
				suffix = ".chatmode.md"
			}
			if !strings.HasSuffix(entry.Name, suffix) {
				continue
			}
			content, err := l.gh.FetchFile(ctx, ref.Owner, ref.Repo, ref.Ref, entry.Path)
			if err != nil {
				continue
			}
			_ = os.WriteFile(filepath.Join(localDir, entry.Name), []byte(content), 0o644)
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

// sourceFor returns the source label for a directory.
// For relative paths (e.g. ".") it resolves to the absolute path first so the
// label is a meaningful directory name rather than ".".
func sourceFor(dir string) string {
	abs, err := filepath.Abs(dir)
	if err != nil {
		return filepath.Base(dir)
	}
	return filepath.Base(abs)
}
