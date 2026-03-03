package main

import (
	"context"
	"encoding/json"
	"os"
	"strings"
	"testing"

	"github.com/modelcontextprotocol/go-sdk/mcp"

	"github.com/Arkestone/mcp/servers/mcp-memory/internal/store"
)

// ---------------------------------------------------------------------------
// helpers
// ---------------------------------------------------------------------------

func setupTestStore(t *testing.T) *store.Store {
	t.Helper()
	st, err := store.New(t.TempDir())
	if err != nil {
		t.Fatalf("store.New: %v", err)
	}
	return st
}

func connectClientServer(t *testing.T, ctx context.Context, srv *mcp.Server) *mcp.ClientSession {
	t.Helper()
	st, ct := mcp.NewInMemoryTransports()
	go srv.Run(ctx, st) //nolint:errcheck
	client := mcp.NewClient(&mcp.Implementation{Name: "test-client", Version: "0.0.1"}, nil)
	cs, err := client.Connect(ctx, ct, nil)
	if err != nil {
		t.Fatalf("client.Connect: %v", err)
	}
	t.Cleanup(func() { cs.Close() })
	return cs
}

// textContent extracts the text from the first TextContent in a CallToolResult.
func textContent(t *testing.T, res *mcp.CallToolResult) string {
	t.Helper()
	if len(res.Content) == 0 {
		t.Fatal("no content in tool result")
	}
	tc, ok := res.Content[0].(*mcp.TextContent)
	if !ok {
		t.Fatalf("expected *mcp.TextContent, got %T", res.Content[0])
	}
	return tc.Text
}

// ---------------------------------------------------------------------------
// Unit tests
// ---------------------------------------------------------------------------

func TestExpandHome(t *testing.T) {
	t.Run("absolute path unchanged", func(t *testing.T) {
		got := expandHome("/absolute/path")
		if got != "/absolute/path" {
			t.Errorf("expandHome(%q) = %q, want %q", "/absolute/path", got, "/absolute/path")
		}
	})

	t.Run("tilde expands to home dir", func(t *testing.T) {
		home, _ := os.UserHomeDir()
		got := expandHome("~/foo/bar")
		want := home + "/foo/bar"
		if got != want {
			t.Errorf("expandHome(~/foo/bar) = %q, want %q", got, want)
		}
	})

	t.Run("no tilde prefix unchanged", func(t *testing.T) {
		got := expandHome("relative/path")
		if got != "relative/path" {
			t.Errorf("expandHome(%q) = %q", "relative/path", got)
		}
	})
}

// ---------------------------------------------------------------------------
// Resource tests
// ---------------------------------------------------------------------------

func TestRegisterResources_MemoryAll(t *testing.T) {
	st := setupTestStore(t)
	ctx := context.Background()
	srv := mcp.NewServer(&mcp.Implementation{Name: "test", Version: "0.1.0"}, nil)
	registerResources(srv, st)
	cs := connectClientServer(t, ctx, srv)

	t.Run("empty store", func(t *testing.T) {
		res, err := cs.ReadResource(ctx, &mcp.ReadResourceParams{URI: "memory://all"})
		if err != nil {
			t.Fatalf("ReadResource(memory://all): %v", err)
		}
		if len(res.Contents) == 0 {
			t.Fatal("no contents")
		}
		// Empty store should return empty or whitespace-only text
		if strings.TrimSpace(res.Contents[0].Text) != "" {
			t.Errorf("expected empty text for empty store, got: %q", res.Contents[0].Text)
		}
	})

	t.Run("with memories", func(t *testing.T) {
		st.Remember("Go is great", []string{"programming"})
		st.Remember("Coffee is good", []string{"food", "drinks"})

		res, err := cs.ReadResource(ctx, &mcp.ReadResourceParams{URI: "memory://all"})
		if err != nil {
			t.Fatalf("ReadResource(memory://all): %v", err)
		}
		text := res.Contents[0].Text
		if !strings.Contains(text, "Go is great") {
			t.Errorf("missing first memory: %s", text)
		}
		if !strings.Contains(text, "Coffee is good") {
			t.Errorf("missing second memory: %s", text)
		}
		if !strings.Contains(text, "programming") {
			t.Errorf("missing tag: %s", text)
		}
	})
}

func TestRegisterResources_MemoryByID(t *testing.T) {
	st := setupTestStore(t)
	ctx := context.Background()
	srv := mcp.NewServer(&mcp.Implementation{Name: "test", Version: "0.1.0"}, nil)
	registerResources(srv, st)
	cs := connectClientServer(t, ctx, srv)

	m, _ := st.Remember("important note", nil)

	t.Run("existing memory by ID", func(t *testing.T) {
		res, err := cs.ReadResource(ctx, &mcp.ReadResourceParams{URI: "memory://" + m.ID})
		if err != nil {
			t.Fatalf("ReadResource(memory://%s): %v", m.ID, err)
		}
		if !strings.Contains(res.Contents[0].Text, "important note") {
			t.Errorf("missing content: %s", res.Contents[0].Text)
		}
	})

	t.Run("nonexistent ID returns error", func(t *testing.T) {
		_, err := cs.ReadResource(ctx, &mcp.ReadResourceParams{URI: "memory://nonexistent"})
		if err == nil {
			t.Error("expected error for nonexistent memory")
		}
	})
}

// ---------------------------------------------------------------------------
// Tool tests
// ---------------------------------------------------------------------------

func TestRegisterTools_Remember(t *testing.T) {
	st := setupTestStore(t)
	ctx := context.Background()
	srv := mcp.NewServer(&mcp.Implementation{Name: "test", Version: "0.1.0"}, nil)
	registerTools(srv, st)
	cs := connectClientServer(t, ctx, srv)

	t.Run("remember with tags", func(t *testing.T) {
		res, err := cs.CallTool(ctx, &mcp.CallToolParams{
			Name: "remember",
			Arguments: map[string]any{
				"content": "Go is a systems programming language",
				"tags":    []any{"programming", "go"},
			},
		})
		if err != nil {
			t.Fatalf("CallTool(remember): %v", err)
		}
		text := textContent(t, res)
		var out struct {
			ID   string   `json:"id"`
			Tags []string `json:"tags"`
		}
		if err := json.Unmarshal([]byte(text), &out); err != nil {
			t.Fatalf("unmarshal remember output: %v", err)
		}
		if out.ID == "" {
			t.Error("expected non-empty ID")
		}
		if len(out.Tags) != 2 {
			t.Errorf("expected 2 tags, got %d", len(out.Tags))
		}
	})

	t.Run("remember without tags", func(t *testing.T) {
		res, err := cs.CallTool(ctx, &mcp.CallToolParams{
			Name:      "remember",
			Arguments: map[string]any{"content": "plain note with no tags"},
		})
		if err != nil {
			t.Fatalf("CallTool(remember no tags): %v", err)
		}
		text := textContent(t, res)
		if !strings.Contains(text, `"id"`) {
			t.Errorf("missing id in output: %s", text)
		}
	})

	t.Run("remember with empty content returns error", func(t *testing.T) {
		res, err := cs.CallTool(ctx, &mcp.CallToolParams{
			Name:      "remember",
			Arguments: map[string]any{"content": "   "},
		})
		if err != nil {
			return // protocol error acceptable
		}
		if res != nil && !res.IsError {
			t.Error("expected error for empty content")
		}
	})
}

func TestRegisterTools_Recall(t *testing.T) {
	st := setupTestStore(t)
	ctx := context.Background()
	srv := mcp.NewServer(&mcp.Implementation{Name: "test", Version: "0.1.0"}, nil)
	registerTools(srv, st)
	cs := connectClientServer(t, ctx, srv)

	st.Remember("Go is fast and compiled", []string{"programming", "go"})
	st.Remember("Python is easy to learn", []string{"programming", "python"})
	st.Remember("Coffee helps focus", []string{"food"})

	t.Run("recall all (no filters)", func(t *testing.T) {
		res, err := cs.CallTool(ctx, &mcp.CallToolParams{
			Name:      "recall",
			Arguments: map[string]any{},
		})
		if err != nil {
			t.Fatalf("CallTool(recall all): %v", err)
		}
		text := textContent(t, res)
		if !strings.Contains(text, `"count":3`) {
			t.Errorf("expected count 3: %s", text)
		}
	})

	t.Run("recall by text query", func(t *testing.T) {
		res, err := cs.CallTool(ctx, &mcp.CallToolParams{
			Name:      "recall",
			Arguments: map[string]any{"query": "Go"},
		})
		if err != nil {
			t.Fatalf("CallTool(recall query=Go): %v", err)
		}
		text := textContent(t, res)
		if !strings.Contains(text, "Go is fast") {
			t.Errorf("missing Go memory: %s", text)
		}
		if strings.Contains(text, "Python") {
			t.Errorf("unexpected Python in Go query results: %s", text)
		}
	})

	t.Run("recall by tag", func(t *testing.T) {
		res, err := cs.CallTool(ctx, &mcp.CallToolParams{
			Name:      "recall",
			Arguments: map[string]any{"tags": []any{"programming"}},
		})
		if err != nil {
			t.Fatalf("CallTool(recall tags=[programming]): %v", err)
		}
		text := textContent(t, res)
		if !strings.Contains(text, "Go is fast") || !strings.Contains(text, "Python") {
			t.Errorf("missing programming memories: %s", text)
		}
		if strings.Contains(text, "Coffee") {
			t.Errorf("unexpected food memory in programming results: %s", text)
		}
	})

	t.Run("recall by text and tag", func(t *testing.T) {
		res, err := cs.CallTool(ctx, &mcp.CallToolParams{
			Name: "recall",
			Arguments: map[string]any{
				"query": "fast",
				"tags":  []any{"go"},
			},
		})
		if err != nil {
			t.Fatalf("CallTool(recall text+tags): %v", err)
		}
		text := textContent(t, res)
		if !strings.Contains(text, "Go is fast") {
			t.Errorf("missing matching memory: %s", text)
		}
	})
}

func TestRegisterTools_Forget(t *testing.T) {
	st := setupTestStore(t)
	ctx := context.Background()
	srv := mcp.NewServer(&mcp.Implementation{Name: "test", Version: "0.1.0"}, nil)
	registerTools(srv, st)
	cs := connectClientServer(t, ctx, srv)

	m, _ := st.Remember("to be deleted", nil)

	t.Run("forget existing memory", func(t *testing.T) {
		res, err := cs.CallTool(ctx, &mcp.CallToolParams{
			Name:      "forget",
			Arguments: map[string]any{"id": m.ID},
		})
		if err != nil {
			t.Fatalf("CallTool(forget): %v", err)
		}
		if !strings.Contains(textContent(t, res), "deleted") {
			t.Error("expected 'deleted' in response")
		}
		// Verify it's gone
		_, ok := st.Get(m.ID)
		if ok {
			t.Error("memory should be deleted from store")
		}
	})

	t.Run("forget nonexistent memory returns error", func(t *testing.T) {
		res, err := cs.CallTool(ctx, &mcp.CallToolParams{
			Name:      "forget",
			Arguments: map[string]any{"id": "nonexistent"},
		})
		if err != nil {
			return
		}
		if res != nil && !res.IsError {
			t.Error("expected error for nonexistent memory")
		}
	})
}

func TestRegisterTools_ListMemories(t *testing.T) {
	st := setupTestStore(t)
	ctx := context.Background()
	srv := mcp.NewServer(&mcp.Implementation{Name: "test", Version: "0.1.0"}, nil)
	registerTools(srv, st)
	cs := connectClientServer(t, ctx, srv)

	t.Run("empty store", func(t *testing.T) {
		res, err := cs.CallTool(ctx, &mcp.CallToolParams{Name: "list-memories"})
		if err != nil {
			t.Fatalf("CallTool(list-memories empty): %v", err)
		}
		if !strings.Contains(textContent(t, res), `"count":0`) {
			t.Error("expected count 0 for empty store")
		}
	})

	t.Run("list all", func(t *testing.T) {
		st.Remember("memory one", []string{"tag1"})
		st.Remember("memory two", []string{"tag2"})

		res, err := cs.CallTool(ctx, &mcp.CallToolParams{Name: "list-memories"})
		if err != nil {
			t.Fatalf("CallTool(list-memories): %v", err)
		}
		text := textContent(t, res)
		if !strings.Contains(text, "memory one") || !strings.Contains(text, "memory two") {
			t.Errorf("missing memories: %s", text)
		}
	})

	t.Run("list filtered by tags", func(t *testing.T) {
		res, err := cs.CallTool(ctx, &mcp.CallToolParams{
			Name:      "list-memories",
			Arguments: map[string]any{"tags": []any{"tag1"}},
		})
		if err != nil {
			t.Fatalf("CallTool(list-memories tag1): %v", err)
		}
		text := textContent(t, res)
		if !strings.Contains(text, "memory one") {
			t.Errorf("missing tag1 memory: %s", text)
		}
		if strings.Contains(text, "memory two") {
			t.Errorf("unexpected tag2 memory in tag1 results: %s", text)
		}
	})
}
