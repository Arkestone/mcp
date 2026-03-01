// Package scanner provides on-demand access to Copilot skill definitions
// from local directories and GitHub repositories.
package scanner

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

// Skill represents a single Copilot skill parsed from a SKILL.md file.
type Skill struct {
	Name        string      // from YAML frontmatter
	Description string      // from YAML frontmatter
	Content     string      // markdown body (after frontmatter)
	Source      string      // origin (directory basename or owner/repo)
	Path        string      // relative path to SKILL.md
	URI         string      // MCP resource URI: skills://{source}/{name}
	References  []Reference // supporting docs from references/ subdirectory
}

// Reference represents a supporting document in a skill's references/ directory.
type Reference struct {
	Name    string
	Content string
	Path    string
}

// Scanner provides on-demand access to skill definitions.
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

// List returns all skills. Local dirs are read from disk; repos from cache.
// When no dirs are configured, the current working directory is used as default
// so the server works out-of-the-box when run from a repository root.
// Each dir is scanned at both root level and inside its .github/ subdirectory.
func (s *Scanner) List() []Skill {
	dirs := s.cfg.Sources.Dirs
	if len(dirs) == 0 {
		dirs = []string{"."}
	}
	var out []Skill
	for _, dir := range dirs {
		src := sourceFor(dir)
		out = append(out, scanDir(dir, src)...)
		out = append(out, scanDir(filepath.Join(dir, ".github"), src)...)
	}
	for _, ref := range s.cfg.ParsedRepos() {
		cacheDir := repoCacheDir(s.cfg.Cache.Dir, ref)
		out = append(out, scanDir(cacheDir, ref.Owner+"/"+ref.Repo)...)
	}
	return out
}

// Get returns a single skill by URI.
func (s *Scanner) Get(uri string) (Skill, bool) {
	for _, sk := range s.List() {
		if sk.URI == uri {
			return sk, true
		}
	}
	return Skill{}, false
}

// scanDir reads skills from a directory. Each subdirectory with a SKILL.md is a skill.
func scanDir(dir, source string) []Skill {
	entries, err := os.ReadDir(dir)
	if err != nil {
		return nil
	}

	var out []Skill
	for _, entry := range entries {
		if !entry.IsDir() {
			continue
		}

		skillDir := filepath.Join(dir, entry.Name())
		data, err := os.ReadFile(filepath.Join(skillDir, "SKILL.md"))
		if err != nil {
			continue
		}

		name, description, body := parseFrontmatter(data)
		if name == "" {
			name = entry.Name()
		}

		sk := Skill{
			Name:        name,
			Description: description,
			Content:     body,
			Source:      source,
			Path:        filepath.Join(entry.Name(), "SKILL.md"),
			URI:         fmt.Sprintf("skills://%s/%s", source, name),
		}

		// Load references/ if present
		refsDir := filepath.Join(skillDir, "references")
		if info, err := os.Stat(refsDir); err == nil && info.IsDir() {
			if refEntries, err := os.ReadDir(refsDir); err == nil {
				for _, re := range refEntries {
					if re.IsDir() {
						continue
					}
					content, err := os.ReadFile(filepath.Join(refsDir, re.Name()))
					if err != nil {
						continue
					}
					sk.References = append(sk.References, Reference{
						Name:    re.Name(),
						Content: string(content),
						Path:    filepath.Join(entry.Name(), "references", re.Name()),
					})
				}
			}
		}

		out = append(out, sk)
	}
	return out
}

// parseFrontmatter extracts YAML frontmatter (between --- delimiters) and returns
// the name, description, and remaining markdown body.
func parseFrontmatter(data []byte) (name, description, body string) {
	content := string(data)
	if !strings.HasPrefix(content, "---") {
		return "", "", content
	}

	rest := content[3:]
	idx := strings.Index(rest, "\n---")
	if idx < 0 {
		return "", "", content
	}

	frontmatter := rest[:idx]
	body = rest[idx+4:]
	if len(body) > 0 && body[0] == '\n' {
		body = body[1:]
	}

	// Parse simple YAML fields
	lines := strings.Split(frontmatter, "\n")
	for i, line := range lines {
		trimmed := strings.TrimSpace(line)
		if strings.HasPrefix(trimmed, "name:") {
			name = strings.Trim(strings.TrimSpace(strings.TrimPrefix(trimmed, "name:")), "\"'")
		} else if strings.HasPrefix(trimmed, "description:") {
			val := strings.TrimSpace(strings.TrimPrefix(trimmed, "description:"))
			if val == "|" || val == ">" {
				// Multiline: collect indented lines
				var parts []string
				for j := i + 1; j < len(lines); j++ {
					l := lines[j]
					if len(l) > 0 && l[0] != ' ' && l[0] != '\t' {
						break
					}
					parts = append(parts, strings.TrimSpace(l))
				}
				description = strings.Join(parts, " ")
			} else {
				description = strings.Trim(val, "\"'")
			}
		}
	}

	return name, description, body
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
	if err := os.MkdirAll(cacheDir, 0o755); err != nil {
		return fmt.Errorf("creating cache dir: %w", err)
	}

	ctx := context.Background()

	entries, err := s.gh.FetchDir(ctx, ref.Owner, ref.Repo, ref.Ref, "skills")
	if err != nil {
		return fmt.Errorf("listing skills/: %w", err)
	}

	for _, entry := range entries {
		if entry.Type != "dir" {
			continue
		}

		skillCacheDir := filepath.Join(cacheDir, entry.Name)
		if err := os.MkdirAll(skillCacheDir, 0o755); err != nil {
			continue
		}

		content, err := s.gh.FetchFile(ctx, ref.Owner, ref.Repo, ref.Ref, entry.Path+"/SKILL.md")
		if err != nil {
			continue
		}
		_ = os.WriteFile(filepath.Join(skillCacheDir, "SKILL.md"), []byte(content), 0o644)

		// Fetch references/ if present
		refEntries, err := s.gh.FetchDir(ctx, ref.Owner, ref.Repo, ref.Ref, entry.Path+"/references")
		if err != nil {
			continue
		}
		refsCacheDir := filepath.Join(skillCacheDir, "references")
		if err := os.MkdirAll(refsCacheDir, 0o755); err != nil {
			continue
		}
		for _, re := range refEntries {
			if re.Type == "dir" {
				continue
			}
			refContent, err := s.gh.FetchFile(ctx, ref.Owner, ref.Repo, ref.Ref, re.Path)
			if err != nil {
				continue
			}
			_ = os.WriteFile(filepath.Join(refsCacheDir, re.Name), []byte(refContent), 0o644)
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
