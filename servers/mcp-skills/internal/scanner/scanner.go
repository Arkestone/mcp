// Package scanner provides on-demand access to Copilot skill definitions
// from local directories and GitHub repositories.
package scanner

import (
"bytes"
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

const maxFileSize = 1 << 20 // 1 MiB

type Skill struct {
Name        string
Description string
Tags        []string
Files       []string // glob patterns from frontmatter files: — restricts which file paths this skill applies to
Content     string
Source      string
Path        string
URI         string
References  []Reference
}

func FilterByQuery(skills []Skill, query string) []Skill {
return filter.SortByScore(skills, func(s Skill) int {
return filter.Score(query, s.Name, s.Description, s.Tags)
})
}

// FilterByFilePath returns skills applicable to the given file path.
// Skills without Files patterns are always included (global scope).
// Skills with Files patterns are included only when at least one pattern matches filePath.
// If filePath is empty, all skills are returned unchanged.
func FilterByFilePath(skills []Skill, filePath string) []Skill {
if filePath == "" {
return skills
}
fp := filepath.ToSlash(filePath)
var out []Skill
for _, s := range skills {
if len(s.Files) == 0 || glob.MatchAny(s.Files, fp) {
out = append(out, s)
}
}
return out
}

type Reference struct {
Name    string
Content string
Path    string
}

type Scanner struct {
cfg    *config.Config
gh     *github.Client
syncer *syncer.Syncer
cache  cache.List[Skill]
}

func New(cfg *config.Config, gh *github.Client) *Scanner {
s := &Scanner{cfg: cfg, gh: gh}
s.syncer = syncer.New(cfg.Cache.SyncInterval, s.syncAllRepos)
return s
}

func (s *Scanner) Start(ctx context.Context) { s.syncer.Start(ctx) }
func (s *Scanner) Stop()                     { s.syncer.Stop() }

func (s *Scanner) ForceSync() {
s.cache.Invalidate()
s.syncer.ForceSync()
}

func (s *Scanner) List() []Skill {
return s.cache.Get(s.scan)
}

func (s *Scanner) scan() []Skill {
dirs := s.cfg.Sources.Dirs
if len(dirs) == 0 {
dirs = []string{"."}
}
var out []Skill
for _, dir := range dirs {
out = append(out, scanDir(dir, sourceFor(dir))...)
}
for _, ref := range s.cfg.ParsedRepos() {
cacheDir := repoCacheDir(s.cfg.Cache.Dir, ref)
out = append(out, scanDir(cacheDir, ref.Owner+"/"+ref.Repo)...)
}
return out
}

func (s *Scanner) Get(uri string) (Skill, bool) {
for _, sk := range s.List() {
if sk.URI == uri {
return sk, true
}
}
return Skill{}, false
}

var skipDirs = map[string]bool{
".git":         true,
"node_modules": true,
}

func scanDir(dir, source string) []Skill {
var out []Skill
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
if info.Name() != "SKILL.md" || info.Size() > maxFileSize {
return nil
}
data, err := os.ReadFile(fpath)
if err != nil {
return nil
}
name, description, tags, files, body := parseFrontmatter(data)
skillDir := filepath.Dir(fpath)
if name == "" {
name = filepath.Base(skillDir)
}
rel, _ := filepath.Rel(dir, fpath)
sk := Skill{
Name: name, Description: description, Tags: tags, Files: files,
Content: body, Source: source, Path: filepath.ToSlash(rel),
URI: fmt.Sprintf("skills://%s/%s", source, name),
}
refsDir := filepath.Join(skillDir, "references")
if rinfo, err := os.Stat(refsDir); err == nil && rinfo.IsDir() {
if refEntries, err := os.ReadDir(refsDir); err == nil {
for _, re := range refEntries {
if re.IsDir() {
continue
}
ri, err := re.Info()
if err != nil || ri.Size() > maxFileSize {
continue
}
content, err := os.ReadFile(filepath.Join(refsDir, re.Name()))
if err != nil {
continue
}
refRel, _ := filepath.Rel(dir, filepath.Join(refsDir, re.Name()))
sk.References = append(sk.References, Reference{
Name: re.Name(), Content: string(content), Path: filepath.ToSlash(refRel),
})
}
}
}
out = append(out, sk)
return nil
})
return out
}

type skillMeta struct {
Name        string      `yaml:"name"`
Description string      `yaml:"description"`
Tags        interface{} `yaml:"tags"`
Files       interface{} `yaml:"files"` // string or []string glob patterns
}

func parseFrontmatter(data []byte) (name, description string, tags, files []string, body string) {
var meta skillMeta
rest, _ := frontmatter.Parse(bytes.NewReader(data), &meta)
return meta.Name, meta.Description, toStringSlice(meta.Tags), toStringSlice(meta.Files), string(rest)
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

func (s *Scanner) syncAllRepos() {
s.cache.Invalidate()
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
entries, err := s.gh.FetchDirRecursive(ctx, ref.Owner, ref.Repo, ref.Ref, "")
if err != nil {
if github.IsRateLimitError(err) || (s.gh.Token == "" && isAuthError(err)) {
log.Printf("sync %s/%s: API failed (%v), falling back to ZIP download", ref.Owner, ref.Repo, err)
return s.gh.FetchZipAndExtract(ctx, ref.Owner, ref.Repo, ref.Ref, cacheDir)
}
return fmt.Errorf("listing repo: %w", err)
}
skillDirs := map[string]bool{}
for _, e := range entries {
if path.Base(e.Path) == "SKILL.md" {
skillDirs[path.Dir(e.Path)] = true
}
}
for _, entry := range entries {
base := path.Base(entry.Path)
parentDir := path.Dir(entry.Path)
parentBase := path.Base(parentDir)
if base != "SKILL.md" && !(parentBase == "references" && skillDirs[path.Dir(parentDir)]) {
continue
}
content, err := s.gh.FetchFile(ctx, ref.Owner, ref.Repo, ref.Ref, entry.Path)
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

func repoCacheDir(cacheBase string, ref config.RepoRef) string {
key := ref.Owner + "_" + ref.Repo
if ref.Ref != "" {
key += "_" + ref.Ref
}
return filepath.Join(cacheBase, key)
}

// isAuthError reports whether err is a GitHub authentication / access error.
func isAuthError(err error) bool {
if err == nil {
return false
}
s := err.Error()
return strings.Contains(s, "HTTP 401") || strings.Contains(s, "HTTP 403") || strings.Contains(s, "HTTP 404")
}

func sourceFor(dir string) string {
abs, err := filepath.Abs(dir)
if err != nil {
return filepath.Base(dir)
}
return filepath.Base(abs)
}
