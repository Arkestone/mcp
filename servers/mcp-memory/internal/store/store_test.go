package store

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"
)

func newTestStore(t *testing.T) *Store {
	t.Helper()
	st, err := New(t.TempDir())
	if err != nil {
		t.Fatalf("New: %v", err)
	}
	return st
}

func TestNew_createsDir(t *testing.T) {
	dir := filepath.Join(t.TempDir(), "nested", "memory")
	st, err := New(dir)
	if err != nil {
		t.Fatalf("New: %v", err)
	}
	if st == nil {
		t.Fatal("expected non-nil store")
	}
	if _, err := os.Stat(dir); err != nil {
		t.Errorf("directory not created: %v", err)
	}
}

func TestRemember(t *testing.T) {
	st := newTestStore(t)

	m, err := st.Remember("buy milk", []string{"shopping", "errands"})
	if err != nil {
		t.Fatalf("Remember: %v", err)
	}

	if m.ID == "" {
		t.Error("expected non-empty ID")
	}
	if m.Content != "buy milk" {
		t.Errorf("Content = %q, want %q", m.Content, "buy milk")
	}
	if len(m.Tags) != 2 {
		t.Errorf("Tags = %v, want 2 tags", m.Tags)
	}
	if m.CreatedAt.IsZero() {
		t.Error("CreatedAt should not be zero")
	}

	// Verify file exists on disk
	path := filepath.Join(st.dir, m.ID+".md")
	if _, err := os.Stat(path); err != nil {
		t.Errorf("file not created: %v", err)
	}
}

func TestGet(t *testing.T) {
	st := newTestStore(t)

	m, err := st.Remember("hello world", []string{"test"})
	if err != nil {
		t.Fatalf("Remember: %v", err)
	}

	got, ok := st.Get(m.ID)
	if !ok {
		t.Fatal("Get returned false for existing memory")
	}
	if got.ID != m.ID {
		t.Errorf("ID = %q, want %q", got.ID, m.ID)
	}
	if got.Content != "hello world" {
		t.Errorf("Content = %q, want %q", got.Content, "hello world")
	}
	if len(got.Tags) != 1 || got.Tags[0] != "test" {
		t.Errorf("Tags = %v, want [test]", got.Tags)
	}

	_, ok = st.Get("nonexistent")
	if ok {
		t.Error("Get should return false for nonexistent ID")
	}
}

func TestList(t *testing.T) {
	st := newTestStore(t)

	// Create memories with slight time gap to ensure ordering
	m1, _ := st.Remember("first", nil)
	time.Sleep(10 * time.Millisecond)
	m2, _ := st.Remember("second", nil)
	time.Sleep(10 * time.Millisecond)
	m3, _ := st.Remember("third", nil)

	list, err := st.List()
	if err != nil {
		t.Fatalf("List: %v", err)
	}
	if len(list) != 3 {
		t.Fatalf("List returned %d memories, want 3", len(list))
	}

	// Newest first
	if list[0].ID != m3.ID {
		t.Errorf("first = %q, want %q (newest)", list[0].ID, m3.ID)
	}
	if list[1].ID != m2.ID {
		t.Errorf("second = %q, want %q", list[1].ID, m2.ID)
	}
	if list[2].ID != m1.ID {
		t.Errorf("third = %q, want %q (oldest)", list[2].ID, m1.ID)
	}
}

func TestRecall_empty(t *testing.T) {
	st := newTestStore(t)
	st.Remember("alpha", nil)
	st.Remember("beta", []string{"b"})
	st.Remember("gamma", []string{"c"})

	results, err := st.Recall("", nil)
	if err != nil {
		t.Fatalf("Recall: %v", err)
	}
	if len(results) != 3 {
		t.Errorf("Recall empty = %d, want 3", len(results))
	}
}

func TestRecall_byText(t *testing.T) {
	st := newTestStore(t)
	st.Remember("the quick brown fox", nil)
	st.Remember("lazy dog runs", nil)
	st.Remember("quick silver", nil)

	results, err := st.Recall("quick", nil)
	if err != nil {
		t.Fatalf("Recall: %v", err)
	}
	if len(results) != 2 {
		t.Errorf("Recall byText = %d, want 2", len(results))
	}
	for _, r := range results {
		if !strings.Contains(strings.ToLower(r.Content), "quick") {
			t.Errorf("result %q does not contain 'quick'", r.Content)
		}
	}
}

func TestRecall_byTags(t *testing.T) {
	st := newTestStore(t)
	st.Remember("item A", []string{"work", "urgent"})
	st.Remember("item B", []string{"work"})
	st.Remember("item C", []string{"personal"})

	results, err := st.Recall("", []string{"work"})
	if err != nil {
		t.Fatalf("Recall: %v", err)
	}
	if len(results) != 2 {
		t.Errorf("Recall byTags(work) = %d, want 2", len(results))
	}

	results, err = st.Recall("", []string{"work", "urgent"})
	if err != nil {
		t.Fatalf("Recall: %v", err)
	}
	if len(results) != 1 {
		t.Errorf("Recall byTags(work,urgent) = %d, want 1", len(results))
	}
}

func TestRecall_textAndTags(t *testing.T) {
	st := newTestStore(t)
	st.Remember("buy groceries", []string{"shopping"})
	st.Remember("buy coffee beans", []string{"shopping", "food"})
	st.Remember("read a book", []string{"leisure"})

	results, err := st.Recall("buy", []string{"shopping"})
	if err != nil {
		t.Fatalf("Recall: %v", err)
	}
	if len(results) != 2 {
		t.Errorf("Recall textAndTags = %d, want 2", len(results))
	}

	results, err = st.Recall("coffee", []string{"shopping", "food"})
	if err != nil {
		t.Fatalf("Recall: %v", err)
	}
	if len(results) != 1 {
		t.Errorf("Recall textAndTags precise = %d, want 1", len(results))
	}
	if results[0].Content != "buy coffee beans" {
		t.Errorf("Content = %q", results[0].Content)
	}
}

func TestForget(t *testing.T) {
	st := newTestStore(t)

	m, err := st.Remember("to be deleted", nil)
	if err != nil {
		t.Fatalf("Remember: %v", err)
	}

	if err := st.Forget(m.ID); err != nil {
		t.Fatalf("Forget: %v", err)
	}

	_, ok := st.Get(m.ID)
	if ok {
		t.Error("memory should be gone after Forget")
	}

	// Forget nonexistent
	err = st.Forget("doesnotexist")
	if err == nil {
		t.Error("expected error for nonexistent ID")
	}
	if !strings.Contains(err.Error(), "not found") {
		t.Errorf("error = %q, want 'not found'", err.Error())
	}
}

func TestParse_noFrontmatter(t *testing.T) {
	m, err := parse("myid", "plain content without frontmatter")
	if err != nil {
		t.Fatalf("parse: %v", err)
	}
	if m.ID != "myid" {
		t.Errorf("ID = %q, want %q", m.ID, "myid")
	}
	if m.Content != "plain content without frontmatter" {
		t.Errorf("Content = %q", m.Content)
	}
	if len(m.Tags) != 0 {
		t.Errorf("Tags = %v, want empty", m.Tags)
	}
	if !m.CreatedAt.IsZero() {
		t.Error("CreatedAt should be zero for no frontmatter")
	}
}

// TestRemember_noTags_writeFormat verifies the serialised file does NOT contain
// a "tags:" line when no tags are provided (kills store:137 CONDITIONALS_BOUNDARY).
func TestRemember_noTags_writeFormat(t *testing.T) {
	st := newTestStore(t)
	m, err := st.Remember("no tags here", nil)
	if err != nil {
		t.Fatalf("Remember: %v", err)
	}
	raw, err := os.ReadFile(filepath.Join(st.dir, m.ID+".md"))
	if err != nil {
		t.Fatalf("ReadFile: %v", err)
	}
	if strings.Contains(string(raw), "tags:") {
		t.Errorf("file should not contain 'tags:' for a no-tags memory, got:\n%s", raw)
	}
}

// TestRemember_singleTag_format verifies a single tag is written as "tags: [mytag]"
// with no leading comma (kills store:140 CONDITIONALS_BOUNDARY on i>0).
func TestRemember_singleTag_format(t *testing.T) {
	st := newTestStore(t)
	m, err := st.Remember("one tag", []string{"mytag"})
	if err != nil {
		t.Fatalf("Remember: %v", err)
	}
	raw, err := os.ReadFile(filepath.Join(st.dir, m.ID+".md"))
	if err != nil {
		t.Fatalf("ReadFile: %v", err)
	}
	if !strings.Contains(string(raw), "tags: [mytag]") {
		t.Errorf("expected 'tags: [mytag]' in file, got:\n%s", raw)
	}
	// roundtrip
	got, ok := st.Get(m.ID)
	if !ok {
		t.Fatal("Get returned false")
	}
	if len(got.Tags) != 1 || got.Tags[0] != "mytag" {
		t.Errorf("Tags = %v, want [mytag]", got.Tags)
	}
}

// TestParse_emptyFrontmatter verifies parse handles content where the frontmatter
// separator "\n---\n" is at position 0 in rest (idx==0), i.e. "---\n\n---\ncontent".
// Kills store:169 CONDITIONALS_BOUNDARY (idx < 0  →  idx <= 0 would skip parsing).
func TestParse_emptyFrontmatter(t *testing.T) {
	m, err := parse("myid", "---\n\n---\ncontent here")
	if err != nil {
		t.Fatalf("parse: %v", err)
	}
	if m.ID != "myid" {
		t.Errorf("ID = %q", m.ID)
	}
	if m.Content != "content here" {
		t.Errorf("Content = %q, want %q", m.Content, "content here")
	}
}

// TestParse_timestamps verifies that valid created_at / updated_at values are
// parsed correctly (kills store:195 CONDITIONALS_NEGATION: err==nil → err!=nil).
func TestParse_timestamps(t *testing.T) {
	now := time.Now().UTC().Truncate(time.Millisecond)
	content := "---\n" +
		"id: tid\n" +
		"created_at: " + now.Format(time.RFC3339Nano) + "\n" +
		"updated_at: " + now.Add(time.Second).Format(time.RFC3339Nano) + "\n" +
		"---\nsome content"
	m, err := parse("tid", content)
	if err != nil {
		t.Fatalf("parse: %v", err)
	}
	if m.CreatedAt.IsZero() {
		t.Error("CreatedAt should not be zero when valid created_at is present")
	}
	if !m.CreatedAt.Equal(now) {
		t.Errorf("CreatedAt = %v, want %v", m.CreatedAt, now)
	}
	if !m.UpdatedAt.Equal(now.Add(time.Second)) {
		t.Errorf("UpdatedAt = %v, want %v", m.UpdatedAt, now.Add(time.Second))
	}
}

// TestRemember_timestampRoundtrip verifies that CreatedAt survives a write+read
// cycle (cross-validates store:195 via the public API).
func TestRemember_timestampRoundtrip(t *testing.T) {
	st := newTestStore(t)
	original, err := st.Remember("roundtrip", nil)
	if err != nil {
		t.Fatalf("Remember: %v", err)
	}
	got, ok := st.Get(original.ID)
	if !ok {
		t.Fatal("Get returned false")
	}
	if !got.CreatedAt.Equal(original.CreatedAt) {
		t.Errorf("CreatedAt mismatch: got %v, want %v", got.CreatedAt, original.CreatedAt)
	}
	if !got.UpdatedAt.Equal(original.UpdatedAt) {
		t.Errorf("UpdatedAt mismatch: got %v, want %v", got.UpdatedAt, original.UpdatedAt)
	}
}

func TestMatchesTags(t *testing.T) {
	tests := []struct {
		name   string
		mTags  []string
		filter []string
		want   bool
	}{
		{"empty filter matches all", []string{"a", "b"}, nil, true},
		{"empty filter empty tags", nil, nil, true},
		{"single tag match", []string{"work"}, []string{"work"}, true},
		{"single tag no match", []string{"work"}, []string{"home"}, false},
		{"all tags match", []string{"work", "urgent"}, []string{"work", "urgent"}, true},
		{"partial tags match fails", []string{"work"}, []string{"work", "urgent"}, false},
		{"case insensitive", []string{"Work"}, []string{"work"}, true},
		{"no memory tags, filter set", nil, []string{"work"}, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := Memory{Tags: tt.mTags}
			got := matchesTags(m, tt.filter)
			if got != tt.want {
				t.Errorf("matchesTags(%v, %v) = %v, want %v", tt.mTags, tt.filter, got, tt.want)
			}
		})
	}
}
