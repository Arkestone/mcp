// Package store provides a file-based persistent memory store.
// Each memory is a markdown file with YAML frontmatter in a configured directory.
package store

import (
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"time"

	"github.com/adrg/frontmatter"
	"github.com/google/uuid"
)

// Memory represents a single stored memory.
type Memory struct {
	ID        string
	Tags      []string
	Content   string // the text to remember
	CreatedAt time.Time
	UpdatedAt time.Time
}

// Store manages memories in a local directory.
type Store struct {
	dir string
}

// New creates a Store backed by the given directory.
func New(dir string) (*Store, error) {
	if err := os.MkdirAll(dir, 0o755); err != nil {
		return nil, fmt.Errorf("creating memory store dir: %w", err)
	}
	return &Store{dir: dir}, nil
}

// Remember stores a new memory and returns its ID.
func (s *Store) Remember(content string, tags []string) (Memory, error) {
	now := time.Now().UTC()
	id := generateID()
	m := Memory{
		ID:        id,
		Tags:      tags,
		Content:   content,
		CreatedAt: now,
		UpdatedAt: now,
	}
	if err := s.write(m); err != nil {
		return Memory{}, err
	}
	return m, nil
}

// Recall searches memories by tags and/or text. Empty query returns all.
func (s *Store) Recall(query string, tags []string) ([]Memory, error) {
	all, err := s.List()
	if err != nil {
		return nil, err
	}

	query = strings.ToLower(strings.TrimSpace(query))

	var out []Memory
	for _, m := range all {
		if !matchesTags(m, tags) {
			continue
		}
		if query != "" && !strings.Contains(strings.ToLower(m.Content), query) {
			continue
		}
		out = append(out, m)
	}
	return out, nil
}

// Forget deletes a memory by ID. Returns error if not found.
func (s *Store) Forget(id string) error {
	path := filepath.Join(s.dir, id+".md")
	if err := os.Remove(path); err != nil {
		if os.IsNotExist(err) {
			return fmt.Errorf("memory not found: %s", id)
		}
		return err
	}
	return nil
}

// Get returns a single memory by ID.
func (s *Store) Get(id string) (Memory, bool) {
	path := filepath.Join(s.dir, id+".md")
	data, err := os.ReadFile(path)
	if err != nil {
		return Memory{}, false
	}
	m, err := parse(id, string(data))
	if err != nil {
		return Memory{}, false
	}
	return m, true
}

// List returns all memories sorted by created_at descending (newest first).
func (s *Store) List() ([]Memory, error) {
	entries, err := os.ReadDir(s.dir)
	if err != nil {
		return nil, err
	}

	var out []Memory
	for _, entry := range entries {
		if entry.IsDir() || !strings.HasSuffix(entry.Name(), ".md") {
			continue
		}
		id := strings.TrimSuffix(entry.Name(), ".md")
		data, err := os.ReadFile(filepath.Join(s.dir, entry.Name()))
		if err != nil {
			continue
		}
		m, err := parse(id, string(data))
		if err != nil {
			continue
		}
		out = append(out, m)
	}

	// Sort newest first
	sort.Slice(out, func(i, j int) bool {
		return out[i].CreatedAt.After(out[j].CreatedAt)
	})
	return out, nil
}

func (s *Store) write(m Memory) error {
	var sb strings.Builder
	sb.WriteString("---\n")
	fmt.Fprintf(&sb, "id: %s\n", m.ID)
	if len(m.Tags) > 0 {
		sb.WriteString("tags: [")
		for i, tag := range m.Tags {
			if i > 0 {
				sb.WriteString(", ")
			}
			sb.WriteString(tag)
		}
		sb.WriteString("]\n")
	}
	fmt.Fprintf(&sb, "created_at: %s\n", m.CreatedAt.Format(time.RFC3339Nano))
	fmt.Fprintf(&sb, "updated_at: %s\n", m.UpdatedAt.Format(time.RFC3339Nano))
	sb.WriteString("---\n")
	sb.WriteString(m.Content)
	if !strings.HasSuffix(m.Content, "\n") {
		sb.WriteString("\n")
	}

	return os.WriteFile(filepath.Join(s.dir, m.ID+".md"), []byte(sb.String()), 0o644)
}

type memoryFrontmatter struct {
	Tags      []string `yaml:"tags"`
	CreatedAt string   `yaml:"created_at"`
	UpdatedAt string   `yaml:"updated_at"`
}

// parse reads a Memory from its file contents.
func parse(id, content string) (Memory, error) {
	m := Memory{ID: id}

	if !strings.HasPrefix(content, "---") {
		m.Content = content
		return m, nil
	}

	var meta memoryFrontmatter
	rest, _ := frontmatter.Parse(strings.NewReader(content), &meta)

	m.Tags = meta.Tags
	body := strings.TrimPrefix(string(rest), "\n")
	m.Content = strings.TrimRight(body, "\n")

	if t, err := time.Parse(time.RFC3339Nano, meta.CreatedAt); err == nil {
		m.CreatedAt = t
	}
	if t, err := time.Parse(time.RFC3339Nano, meta.UpdatedAt); err == nil {
		m.UpdatedAt = t
	}

	return m, nil
}

func matchesTags(m Memory, tags []string) bool {
	if len(tags) == 0 {
		return true
	}
	tagSet := make(map[string]struct{}, len(m.Tags))
	for _, t := range m.Tags {
		tagSet[strings.ToLower(t)] = struct{}{}
	}
	for _, want := range tags {
		if _, ok := tagSet[strings.ToLower(want)]; !ok {
			return false
		}
	}
	return true
}

// generateID returns a new random UUID string.
func generateID() string {
	return uuid.New().String()
}
