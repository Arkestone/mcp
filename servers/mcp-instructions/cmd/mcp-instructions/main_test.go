package main

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/Arkestone/mcp/pkg/config"
	"github.com/Arkestone/mcp/pkg/github"
	"github.com/Arkestone/mcp/pkg/optimizer"
	"github.com/Arkestone/mcp/pkg/server"
	"github.com/Arkestone/mcp/pkg/testutil"
	"github.com/Arkestone/mcp/servers/mcp-instructions/internal/loader"
	"github.com/modelcontextprotocol/go-sdk/mcp"
)

// ---------------------------------------------------------------------------
// helpers (same pattern as integration_test.go, safe to redefine here because
// integration tests only build with -tags integration)
// ---------------------------------------------------------------------------

func setupTestEnv(t *testing.T) (*loader.Loader, *optimizer.Optimizer, func()) {
	t.Helper()

	tmpDir := t.TempDir()

	ghDir := filepath.Join(tmpDir, ".github")
	if err := os.MkdirAll(ghDir, 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(ghDir, "copilot-instructions.md"), []byte("Use Go idioms.\n"), 0o644); err != nil {
		t.Fatal(err)
	}

	instrDir := filepath.Join(ghDir, "instructions")
	if err := os.MkdirAll(instrDir, 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(instrDir, "testing.instructions.md"), []byte("Write table-driven tests.\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(instrDir, "style.instructions.md"), []byte("Use gofmt.\n"), 0o644); err != nil {
		t.Fatal(err)
	}

	cfg := &config.Config{
		Sources: config.Sources{
			Dirs: []string{tmpDir},
		},
		Cache: config.CacheConfig{
			Dir:          filepath.Join(tmpDir, "cache"),
			SyncInterval: 5 * time.Minute,
		},
	}

	ldr := loader.New(cfg, &github.Client{})
	ldr.Start(context.Background())
	opt := optimizer.New(cfg.LLM)

	cleanup := func() { ldr.Stop() }
	return ldr, opt, cleanup
}

func connectClientServer(t *testing.T, ctx context.Context, srv *mcp.Server) *mcp.ClientSession {
	t.Helper()

	st, ct := mcp.NewInMemoryTransports()
	if _, err := srv.Connect(ctx, st, nil); err != nil {
		t.Fatalf("server connect: %v", err)
	}

	client := mcp.NewClient(&mcp.Implementation{Name: "test-client", Version: "0.0.1"}, nil)
	cs, err := client.Connect(ctx, ct, nil)
	if err != nil {
		t.Fatalf("client connect: %v", err)
	}
	return cs
}

func source(ldr *loader.Loader) string {
	items := ldr.List()
	if len(items) == 0 {
		return ""
	}
	return items[0].Source
}

// ---------------------------------------------------------------------------
// Existing tests
// ---------------------------------------------------------------------------

func TestShouldOptimize(t *testing.T) {
	opt := optimizer.New(testutil.LLMConfig())
	var nilOpt *optimizer.Optimizer

	tests := []struct {
		name          string
		opt           *optimizer.Optimizer
		globalDefault bool
		perRequest    string
		want          bool
	}{
		{"nil optimizer", nilOpt, true, "", false},
		{"nil optimizer with true override", nilOpt, true, "true", false},
		{"enabled, global true, no override", opt, true, "", true},
		{"enabled, global false, no override", opt, false, "", false},
		{"enabled, global false, override true", opt, false, "true", true},
		{"enabled, global true, override false", opt, true, "false", false},
		{"enabled, override yes", opt, false, "yes", true},
		{"enabled, override no", opt, true, "no", false},
		{"enabled, override 1", opt, false, "1", true},
		{"enabled, override 0", opt, true, "0", false},
		{"enabled, override TRUE (case insensitive)", opt, false, "TRUE", true},
		{"enabled, override FALSE (case insensitive)", opt, true, "FALSE", false},
		{"enabled, override garbage uses default true", opt, true, "maybe", true},
		{"enabled, override garbage uses default false", opt, false, "maybe", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := server.ShouldOptimize(tt.opt, tt.globalDefault, tt.perRequest)
			if got != tt.want {
				t.Errorf("ShouldOptimize(%v, %v, %q) = %v, want %v",
					tt.opt != nil, tt.globalDefault, tt.perRequest, got, tt.want)
			}
		})
	}
}

// ---------------------------------------------------------------------------
// 1. TestToOptimizerInputs — calls the real toOptimizerInputs()
// ---------------------------------------------------------------------------

func TestToOptimizerInputs(t *testing.T) {
	instructions := []loader.Instruction{
		{Source: "repo1", Path: "path1", Content: "content1", URI: "instructions://repo1/a"},
		{Source: "repo2", Path: "path2", Content: "content2", URI: "instructions://repo2/b"},
	}

	inputs := toOptimizerInputs(instructions)

	if len(inputs) != 2 {
		t.Fatalf("got %d inputs, want 2", len(inputs))
	}
	if inputs[0].Source != "repo1" || inputs[0].Path != "path1" || inputs[0].Content != "content1" {
		t.Errorf("inputs[0] = %+v", inputs[0])
	}
	if inputs[1].Source != "repo2" || inputs[1].Path != "path2" || inputs[1].Content != "content2" {
		t.Errorf("inputs[1] = %+v", inputs[1])
	}

	// empty slice
	empty := toOptimizerInputs(nil)
	if len(empty) != 0 {
		t.Errorf("expected 0 inputs for nil, got %d", len(empty))
	}
}

// ---------------------------------------------------------------------------
// 2. TestFilterBySource
// ---------------------------------------------------------------------------

func TestFilterBySource(t *testing.T) {
	instructions := []loader.Instruction{
		{Source: "repo1", Path: "p1", Content: "c1"},
		{Source: "repo2", Path: "p2", Content: "c2"},
		{Source: "repo1", Path: "p3", Content: "c3"},
	}

	t.Run("empty source returns all", func(t *testing.T) {
		got := filterBySource(instructions, "")
		if len(got) != 3 {
			t.Errorf("expected 3, got %d", len(got))
		}
	})

	t.Run("specific source returns matching", func(t *testing.T) {
		got := filterBySource(instructions, "repo1")
		if len(got) != 2 {
			t.Fatalf("expected 2, got %d", len(got))
		}
		for _, inst := range got {
			if inst.Source != "repo1" {
				t.Errorf("unexpected source %q", inst.Source)
			}
		}
	})

	t.Run("non-matching source returns empty", func(t *testing.T) {
		got := filterBySource(instructions, "nonexistent")
		if len(got) != 0 {
			t.Errorf("expected 0, got %d", len(got))
		}
	})

	t.Run("empty list returns empty", func(t *testing.T) {
		got := filterBySource(nil, "repo1")
		if len(got) != 0 {
			t.Errorf("expected 0, got %d", len(got))
		}
	})
}

// ---------------------------------------------------------------------------
// 3. TestPromptResult
// ---------------------------------------------------------------------------

func TestPromptResult(t *testing.T) {
	res := promptResult("my description", "some text content")

	if res.Description != "my description" {
		t.Errorf("Description = %q, want %q", res.Description, "my description")
	}
	if len(res.Messages) != 1 {
		t.Fatalf("expected 1 message, got %d", len(res.Messages))
	}
	if res.Messages[0].Role != "user" {
		t.Errorf("Role = %q, want %q", res.Messages[0].Role, "user")
	}
	tc, ok := res.Messages[0].Content.(*mcp.TextContent)
	if !ok {
		t.Fatalf("expected *mcp.TextContent, got %T", res.Messages[0].Content)
	}
	if tc.Text != "some text content" {
		t.Errorf("Text = %q, want %q", tc.Text, "some text content")
	}
}

// ---------------------------------------------------------------------------
// 4. TestOptimizeContent
// ---------------------------------------------------------------------------

func TestOptimizeContent(t *testing.T) {
	ctx := context.Background()
	inputs := []optimizer.ContentInput{
		{Source: "s1", Path: "p1", Content: "c1"},
		{Source: "s2", Path: "p2", Content: "c2"},
	}
	fallback := optimizer.ConcatRaw(inputs)

	t.Run("nil optimizer falls back to ConcatRaw", func(t *testing.T) {
		got := optimizeContent(ctx, nil, true, "", inputs)
		if got != fallback {
			t.Errorf("got %q, want ConcatRaw output", got)
		}
	})

	t.Run("enabled optimizer succeeds", func(t *testing.T) {
		ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			json.NewEncoder(w).Encode(map[string]any{
				"choices": []map[string]any{
					{"message": map[string]string{"role": "assistant", "content": "optimized result"}},
				},
			})
		}))
		defer ts.Close()

		opt := optimizer.New(optimizer.LLMConfig{
			Endpoint: ts.URL, APIKey: "key", Model: "m", Enabled: true,
		})
		got := optimizeContent(ctx, opt, true, "", inputs)
		if got != "optimized result" {
			t.Errorf("got %q, want %q", got, "optimized result")
		}
	})

	t.Run("enabled optimizer override false falls back", func(t *testing.T) {
		opt := optimizer.New(testutil.LLMConfig())
		got := optimizeContent(ctx, opt, true, "false", inputs)
		if got != fallback {
			t.Errorf("got %q, want ConcatRaw output", got)
		}
	})

	t.Run("enabled optimizer error falls back to ConcatRaw", func(t *testing.T) {
		ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(500)
			w.Write([]byte("server error"))
		}))
		defer ts.Close()

		opt := optimizer.New(optimizer.LLMConfig{
			Endpoint: ts.URL, APIKey: "key", Model: "m", Enabled: true,
		})
		got := optimizeContent(ctx, opt, true, "", inputs)
		if got != fallback {
			t.Errorf("got %q, want ConcatRaw fallback", got)
		}
	})
}

// ---------------------------------------------------------------------------
// 5. TestRegisterResources_ViaProtocol
// ---------------------------------------------------------------------------

func TestRegisterResources_ViaProtocol(t *testing.T) {
	ctx := context.Background()
	ldr, opt, cleanup := setupTestEnv(t)
	defer cleanup()

	srv := mcp.NewServer(&mcp.Implementation{Name: "mcp-instructions", Version: "test"}, nil)
	registerResources(srv, ldr, opt, false)

	cs := connectClientServer(t, ctx, srv)
	defer cs.Close()

	src := source(ldr)
	if src == "" {
		t.Fatal("no instructions found")
	}

	t.Run("instructions://index", func(t *testing.T) {
		res, err := cs.ReadResource(ctx, &mcp.ReadResourceParams{URI: "instructions://index"})
		if err != nil {
			t.Fatalf("ReadResource index: %v", err)
		}
		if len(res.Contents) != 1 {
			t.Fatalf("expected 1 content, got %d", len(res.Contents))
		}
		text := res.Contents[0].Text
		for _, want := range []string{"copilot-instructions", "testing", "style"} {
			if !strings.Contains(text, want) {
				t.Errorf("index missing %q", want)
			}
		}
	})

	t.Run("instructions://optimized", func(t *testing.T) {
		res, err := cs.ReadResource(ctx, &mcp.ReadResourceParams{URI: "instructions://optimized"})
		if err != nil {
			t.Fatalf("ReadResource optimized: %v", err)
		}
		text := res.Contents[0].Text
		for _, want := range []string{"Go idioms", "table-driven", "gofmt"} {
			if !strings.Contains(text, want) {
				t.Errorf("optimized missing %q", want)
			}
		}
	})

	t.Run("individual instruction via template URI", func(t *testing.T) {
		res, err := cs.ReadResource(ctx, &mcp.ReadResourceParams{
			URI: "instructions://" + src + "/copilot-instructions",
		})
		if err != nil {
			t.Fatalf("ReadResource individual: %v", err)
		}
		if !strings.Contains(res.Contents[0].Text, "Go idioms") {
			t.Errorf("unexpected content: %q", res.Contents[0].Text)
		}
	})

	t.Run("individual instruction not found", func(t *testing.T) {
		_, err := cs.ReadResource(ctx, &mcp.ReadResourceParams{
			URI: "instructions://" + src + "/nonexistent",
		})
		if err == nil {
			t.Error("expected error for nonexistent instruction")
		}
	})
}

// ---------------------------------------------------------------------------
// 6. TestRegisterTools_ViaProtocol
// ---------------------------------------------------------------------------

func TestRegisterTools_ViaProtocol(t *testing.T) {
	ctx := context.Background()
	ldr, opt, cleanup := setupTestEnv(t)
	defer cleanup()

	srv := mcp.NewServer(&mcp.Implementation{Name: "mcp-instructions", Version: "test"}, nil)
	registerTools(srv, ldr, opt, false)

	cs := connectClientServer(t, ctx, srv)
	defer cs.Close()

	src := source(ldr)

	t.Run("refresh-instructions", func(t *testing.T) {
		res, err := cs.CallTool(ctx, &mcp.CallToolParams{Name: "refresh-instructions"})
		if err != nil {
			t.Fatalf("CallTool: %v", err)
		}
		if res.IsError {
			t.Fatal("returned error")
		}
		raw, _ := json.Marshal(res.StructuredContent)
		var out struct {
			Message string   `json:"message"`
			Count   int      `json:"count"`
			Sources []string `json:"sources"`
		}
		if err := json.Unmarshal(raw, &out); err != nil {
			t.Fatalf("unmarshal: %v", err)
		}
		if out.Count != 3 {
			t.Errorf("count = %d, want 3", out.Count)
		}
		if len(out.Sources) == 0 {
			t.Error("expected at least one source")
		}
	})

	t.Run("list-instructions", func(t *testing.T) {
		res, err := cs.CallTool(ctx, &mcp.CallToolParams{Name: "list-instructions"})
		if err != nil {
			t.Fatalf("CallTool: %v", err)
		}
		raw, _ := json.Marshal(res.StructuredContent)
		var out struct {
			Entries []struct {
				URI    string `json:"uri"`
				Source string `json:"source"`
				Path   string `json:"path"`
			} `json:"entries"`
		}
		if err := json.Unmarshal(raw, &out); err != nil {
			t.Fatalf("unmarshal: %v", err)
		}
		if len(out.Entries) != 3 {
			t.Fatalf("expected 3 entries, got %d", len(out.Entries))
		}
		uris := map[string]bool{}
		for _, e := range out.Entries {
			uris[e.URI] = true
		}
		for _, name := range []string{"copilot-instructions", "testing", "style"} {
			if !uris["instructions://"+src+"/"+name] {
				t.Errorf("missing URI for %s", name)
			}
		}
	})

	t.Run("optimize-instructions", func(t *testing.T) {
		res, err := cs.CallTool(ctx, &mcp.CallToolParams{
			Name:      "optimize-instructions",
			Arguments: map[string]any{},
		})
		if err != nil {
			t.Fatalf("CallTool: %v", err)
		}
		raw, _ := json.Marshal(res.StructuredContent)
		var out struct {
			Content string `json:"content"`
			Sources int    `json:"sources"`
		}
		if err := json.Unmarshal(raw, &out); err != nil {
			t.Fatalf("unmarshal: %v", err)
		}
		if out.Sources != 3 {
			t.Errorf("sources = %d, want 3", out.Sources)
		}
		for _, want := range []string{"Go idioms", "table-driven", "gofmt"} {
			if !strings.Contains(out.Content, want) {
				t.Errorf("content missing %q", want)
			}
		}
	})

	t.Run("optimize-instructions with source filter", func(t *testing.T) {
		res, err := cs.CallTool(ctx, &mcp.CallToolParams{
			Name:      "optimize-instructions",
			Arguments: map[string]any{"source": "nonexistent"},
		})
		if err != nil {
			t.Fatalf("CallTool: %v", err)
		}
		raw, _ := json.Marshal(res.StructuredContent)
		var out struct {
			Sources int `json:"sources"`
		}
		json.Unmarshal(raw, &out)
		if out.Sources != 0 {
			t.Errorf("sources = %d, want 0", out.Sources)
		}
	})
}

// ---------------------------------------------------------------------------
// 7. TestRegisterPrompts_ViaProtocol
// ---------------------------------------------------------------------------

func TestRegisterPrompts_ViaProtocol(t *testing.T) {
	ctx := context.Background()
	ldr, opt, cleanup := setupTestEnv(t)
	defer cleanup()

	srv := mcp.NewServer(&mcp.Implementation{Name: "mcp-instructions", Version: "test"}, nil)
	registerPrompts(srv, ldr, opt, false)

	cs := connectClientServer(t, ctx, srv)
	defer cs.Close()

	t.Run("list prompts", func(t *testing.T) {
		listRes, err := cs.ListPrompts(ctx, nil)
		if err != nil {
			t.Fatalf("ListPrompts: %v", err)
		}
		found := false
		for _, p := range listRes.Prompts {
			if p.Name == "get-instructions" {
				found = true
			}
		}
		if !found {
			t.Error("get-instructions prompt not found")
		}
	})

	t.Run("get-instructions without filter", func(t *testing.T) {
		res, err := cs.GetPrompt(ctx, &mcp.GetPromptParams{Name: "get-instructions"})
		if err != nil {
			t.Fatalf("GetPrompt: %v", err)
		}
		if len(res.Messages) != 1 {
			t.Fatalf("expected 1 message, got %d", len(res.Messages))
		}
		tc, ok := res.Messages[0].Content.(*mcp.TextContent)
		if !ok {
			t.Fatalf("expected TextContent, got %T", res.Messages[0].Content)
		}
		for _, want := range []string{"Go idioms", "table-driven"} {
			if !strings.Contains(tc.Text, want) {
				t.Errorf("prompt missing %q", want)
			}
		}
	})

	t.Run("get-instructions with non-matching source", func(t *testing.T) {
		res, err := cs.GetPrompt(ctx, &mcp.GetPromptParams{
			Name:      "get-instructions",
			Arguments: map[string]string{"source": "nonexistent"},
		})
		if err != nil {
			t.Fatalf("GetPrompt: %v", err)
		}
		tc, ok := res.Messages[0].Content.(*mcp.TextContent)
		if !ok {
			t.Fatalf("expected TextContent, got %T", res.Messages[0].Content)
		}
		if tc.Text != "No instructions found." {
			t.Errorf("expected 'No instructions found.', got %q", tc.Text)
		}
	})

	t.Run("get-instructions with matching source", func(t *testing.T) {
		src := source(ldr)
		res, err := cs.GetPrompt(ctx, &mcp.GetPromptParams{
			Name:      "get-instructions",
			Arguments: map[string]string{"source": src},
		})
		if err != nil {
			t.Fatalf("GetPrompt: %v", err)
		}
		tc, ok := res.Messages[0].Content.(*mcp.TextContent)
		if !ok {
			t.Fatalf("expected TextContent, got %T", res.Messages[0].Content)
		}
		if !strings.Contains(tc.Text, "Go idioms") {
			t.Errorf("expected instructions content, got: %q", tc.Text)
		}
	})
}
