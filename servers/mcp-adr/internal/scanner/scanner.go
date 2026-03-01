// Package scanner provides on-demand access to Architecture Decision Records
// from local directories and GitHub repositories.
package scanner

import (
	"context"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"github.com/Arkestone/mcp/pkg/config"
	"github.com/Arkestone/mcp/pkg/github"
	"github.com/Arkestone/mcp/pkg/syncer"
)

// ADR represents a single Architecture Decision Record.
type ADR struct {
	ID      string // filename without .md extension (e.g. "0001-use-postgresql")
	Title   string // from frontmatter "title:" or derived from ID
	Status  string // proposed|accepted|deprecated|superseded
	Date    string // from frontmatter "date:"
	Content string // full markdown content (including frontmatter)
	Source  string // origin (directory basename or owner/repo)
	Path    string // relative path within the source
	URI     string // MCP resource URI: adrs://{source}/{id}
}

// Scanner provides on-demand access to ADR files.
type Scanner struct {
	cfg    *config.Config
	gh     *github.Client
	syncer *syncer.Syncer
}

// New creates a Scanner with its background syncer.
func New(cfg *config.Config, gh *github.Client) *Scanner {
	s := &Scanner{cfg: cfg, gh: gh}
	s.syncer = syncer.New(cfg.Cache.SyncInterval, s.syncAllRepos)
	return s
}

// Start begins background sync. Stop must be called to shut down.
func (s *Scanner) Start(ctx context.Context) { s.syncer.Start(ctx) }

// Stop shuts down the background sync.
func (s *Scanner) Stop() { s.syncer.Stop() }

// ForceSync triggers an immediate sync of all remote repos.
func (s *Scanner) ForceSync() { s.syncer.ForceSync() }

// List returns all ADRs. Local dirs are read from disk; repos from cache.
func (s *Scanner) List() []ADR {
	var out []ADR
	for _, dir := range s.cfg.Sources.Dirs {
		out = append(out, scanDir(dir, filepath.Base(dir))...)
	}
	for _, ref := range s.cfg.ParsedRepos() {
		cacheDir := repoCacheDir(s.cfg.Cache.Dir, ref)
		out = append(out, scanDir(cacheDir, ref.Owner+"/"+ref.Repo)...)
	}
	return out
}

// Get returns a single ADR by URI.
func (s *Scanner) Get(uri string) (ADR, bool) {
	for _, a := range s.List() {
		if a.URI == uri {
			return a, true
		}
	}
	return ADR{}, false
}

// adrDirs are the subdirectory paths (relative to repo root) scanned for ADRs.
var adrDirs = []string{"docs/adr", "docs/decisions", "doc/adr"}

// scanDir scans a directory for ADR markdown files in known ADR subdirectories.
func scanDir(dir, source string) []ADR {
	var out []ADR
	for _, subdir := range adrDirs {
		adrPath := filepath.Join(dir, filepath.FromSlash(subdir))
		entries, err := os.ReadDir(adrPath)
		if err != nil {
			continue
		}
		for _, entry := range entries {
			if entry.IsDir() || !strings.HasSuffix(entry.Name(), ".md") {
				continue
			}
			content, err := os.ReadFile(filepath.Join(adrPath, entry.Name()))
			if err != nil {
				continue
			}
			id := strings.TrimSuffix(entry.Name(), ".md")
			title, status, date := parseFrontmatter(string(content))
			if title == "" {
				title = humanize(id)
			}
			out = append(out, ADR{
				ID:      id,
				Title:   title,
				Status:  status,
				Date:    date,
				Content: string(content),
				Source:  source,
				Path:    subdir + "/" + entry.Name(),
				URI:     fmt.Sprintf("adrs://%s/%s", source, id),
			})
		}
	}
	// Sort deterministically by ID within this source
	sort.Slice(out, func(i, j int) bool { return out[i].ID < out[j].ID })
	return out
}

// parseFrontmatter extracts title, status, and date from YAML frontmatter.
func parseFrontmatter(content string) (title, status, date string) {
	if !strings.HasPrefix(content, "---") {
		return "", "", ""
	}
	rest := content[3:]
	idx := strings.Index(rest, "\n---")
	if idx < 0 {
		return "", "", ""
	}
	for _, line := range strings.Split(rest[:idx], "\n") {
		trimmed := strings.TrimSpace(line)
		if strings.HasPrefix(trimmed, "title:") {
			title = strings.Trim(strings.TrimSpace(strings.TrimPrefix(trimmed, "title:")), "\"'")
		} else if strings.HasPrefix(trimmed, "status:") {
			status = strings.Trim(strings.TrimSpace(strings.TrimPrefix(trimmed, "status:")), "\"'")
		} else if strings.HasPrefix(trimmed, "date:") {
			date = strings.Trim(strings.TrimSpace(strings.TrimPrefix(trimmed, "date:")), "\"'")
		}
	}
	return title, status, date
}

// humanize converts a filename-style ID to a readable title.
// "0001-use-postgresql" → "0001 Use Postgresql"
func humanize(id string) string {
	parts := strings.Split(id, "-")
	for i := 1; i < len(parts); i++ {
		if len(parts[i]) > 0 {
			parts[i] = strings.ToUpper(parts[i][:1]) + parts[i][1:]
		}
	}
	return strings.Join(parts, " ")
}

func (s *Scanner) syncAllRepos() {
	for _, ref := range s.cfg.ParsedRepos() {
		if err := s.syncRepo(ref); err != nil {
			log.Printf("sync %s/%s: %v", ref.Owner, ref.Repo, err)
		}
	}
}

func (s *Scanner) syncRepo(ref config.RepoRef) error {
	cacheDir := repoCacheDir(s.cfg.Cache.Dir, ref)
	ctx := context.Background()

	for _, subdir := range adrDirs {
		entries, err := s.gh.FetchDir(ctx, ref.Owner, ref.Repo, ref.Ref, subdir)
		if err != nil {
			continue
		}
		localDir := filepath.Join(cacheDir, filepath.FromSlash(subdir))
		if err := os.MkdirAll(localDir, 0o755); err != nil {
			continue
		}
		for _, entry := range entries {
			if entry.Type == "dir" || !strings.HasSuffix(entry.Name, ".md") {
				continue
			}
			content, err := s.gh.FetchFile(ctx, ref.Owner, ref.Repo, ref.Ref, entry.Path)
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
