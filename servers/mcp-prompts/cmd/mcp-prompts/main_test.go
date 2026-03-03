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

	"github.com/modelcontextprotocol/go-sdk/mcp"

	"github.com/Arkestone/mcp/pkg/config"
	"github.com/Arkestone/mcp/pkg/github"
	"github.com/Arkestone/mcp/pkg/optimizer"
	"github.com/Arkestone/mcp/pkg/server"
	"github.com/Arkestone/mcp/pkg/testutil"
	"github.com/Arkestone/mcp/servers/mcp-prompts/internal/loader"
)

// ---------------------------------------------------------------------------
// helpers
// ---------------------------------------------------------------------------

func setupTestEnv(t *testing.T) (*loader.Loader, *optimizer.Optimizer, func()) {
	t.Helper()

	tmpDir := t.TempDir()

	promptsDir := filepath.Join(tmpDir, ".github", "prompts")
	if err := os.MkdirAll(promptsDir, 0o755); err != nil {
		t.Fatal(err)
	}
	compContent := "---\ndescription: Create component\nmode: agent\n---\nCreate a React component.\n"
	if err := os.WriteFile(filepath.Join(promptsDir, "component.prompt.md"), []byte(compContent), 0o644); err != nil {
		t.Fatal(err)
	}

	chatmodesDir := filepath.Join(tmpDir, ".github", "chatmodes")
	if err := os.MkdirAll(chatmodesDir, 0o755); err != nil {
		t.Fatal(err)
	}
	reviewContent := "---\ndescription: Code reviewer\n---\nReview code carefully.\n"
	if err := os.WriteFile(filepath.Join(chatmodesDir, "reviewer.chatmode.md"), []byte(reviewContent), 0o644); err != nil {
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

func sourceOf(ldr *loader.Loader) string {
	items := ldr.List()
	if len(items) == 0 {
		return ""
	}
	return items[0].Source
}

// ---------------------------------------------------------------------------
// 1. TestToOptimizerInputs
// ---------------------------------------------------------------------------

func TestToOptimizerInputs(t *testing.T) {
	prompts := []loader.Prompt{
		{Source: "repo1", Path: "path1", Content: "content1", URI: "prompts://repo1/a"},
		{Source: "repo2", Path: "path2", Content: "content2", URI: "prompts://repo2/b"},
	}

	inputs := toOptimizerInputs(prompts)

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
	prompts := []loader.Prompt{
		{Source: "repo1", Path: "p1", Type: loader.TypePrompt},
		{Source: "repo2", Path: "p2", Type: loader.TypeChatmode},
		{Source: "repo1", Path: "p3", Type: loader.TypeChatmode},
	}

	t.Run("empty source returns all", func(t *testing.T) {
		got := filterBySource(prompts, "")
		if len(got) != 3 {
			t.Errorf("expected 3, got %d", len(got))
		}
	})

	t.Run("specific source returns matching", func(t *testing.T) {
		got := filterBySource(prompts, "repo1")
		if len(got) != 2 {
			t.Fatalf("expected 2, got %d", len(got))
		}
		for _, p := range got {
			if p.Source != "repo1" {
				t.Errorf("unexpected source %q", p.Source)
			}
		}
	})

	t.Run("non-matching source returns empty", func(t *testing.T) {
		got := filterBySource(prompts, "nonexistent")
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
// 5. TestShouldOptimize (re-export from server package)
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
		{"enabled, global true, no override", opt, true, "", true},
		{"enabled, global false, no override", opt, false, "", false},
		{"enabled, global false, override true", opt, false, "true", true},
		{"enabled, global true, override false", opt, true, "false", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := server.ShouldOptimize(tt.opt, tt.globalDefault, tt.perRequest)
			if got != tt.want {
				t.Errorf("ShouldOptimize = %v, want %v", got, tt.want)
			}
		})
	}
}

// ---------------------------------------------------------------------------
// 6. TestRegisterResources_ViaProtocol
// ---------------------------------------------------------------------------

func TestRegisterResources_ViaProtocol(t *testing.T) {
	ctx := context.Background()
	ldr, opt, cleanup := setupTestEnv(t)
	defer cleanup()

	srv := mcp.NewServer(&mcp.Implementation{Name: "mcp-prompts", Version: "test"}, nil)
	registerResources(srv, ldr, opt, false)

	cs := connectClientServer(t, ctx, srv)
	defer cs.Close()

	src := sourceOf(ldr)
	if src == "" {
		t.Fatal("no prompts found")
	}

	t.Run("prompts://index", func(t *testing.T) {
		res, err := cs.ReadResource(ctx, &mcp.ReadResourceParams{URI: "prompts://index"})
		if err != nil {
			t.Fatalf("ReadResource index: %v", err)
		}
		if len(res.Contents) != 1 {
			t.Fatalf("expected 1 content, got %d", len(res.Contents))
		}
		text := res.Contents[0].Text
		for _, want := range []string{"component", "reviewer", "type="} {
			if !strings.Contains(text, want) {
				t.Errorf("index missing %q; got: %q", want, text)
			}
		}
	})

	t.Run("prompts://optimized", func(t *testing.T) {
		res, err := cs.ReadResource(ctx, &mcp.ReadResourceParams{URI: "prompts://optimized"})
		if err != nil {
			t.Fatalf("ReadResource optimized: %v", err)
		}
		text := res.Contents[0].Text
		for _, want := range []string{"React component", "Review code"} {
			if !strings.Contains(text, want) {
				t.Errorf("optimized missing %q", want)
			}
		}
	})

	t.Run("individual prompt via template URI", func(t *testing.T) {
		res, err := cs.ReadResource(ctx, &mcp.ReadResourceParams{
			URI: "prompts://" + src + "/component",
		})
		if err != nil {
			t.Fatalf("ReadResource individual: %v", err)
		}
		if !strings.Contains(res.Contents[0].Text, "React component") {
			t.Errorf("unexpected content: %q", res.Contents[0].Text)
		}
	})

	t.Run("individual prompt not found", func(t *testing.T) {
		_, err := cs.ReadResource(ctx, &mcp.ReadResourceParams{
			URI: "prompts://" + src + "/nonexistent",
		})
		if err == nil {
			t.Error("expected error for nonexistent prompt")
		}
	})
}

// ---------------------------------------------------------------------------
// 7. TestRegisterTools_ViaProtocol
// ---------------------------------------------------------------------------

func TestRegisterTools_ViaProtocol(t *testing.T) {
	ctx := context.Background()
	ldr, opt, cleanup := setupTestEnv(t)
	defer cleanup()

	srv := mcp.NewServer(&mcp.Implementation{Name: "mcp-prompts", Version: "test"}, nil)
	registerTools(srv, ldr, opt, false)

	cs := connectClientServer(t, ctx, srv)
	defer cs.Close()

	src := sourceOf(ldr)

	t.Run("refresh-prompts", func(t *testing.T) {
		res, err := cs.CallTool(ctx, &mcp.CallToolParams{Name: "refresh-prompts"})
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
		if out.Count != 2 {
			t.Errorf("count = %d, want 2", out.Count)
		}
		if len(out.Sources) == 0 {
			t.Error("expected at least one source")
		}
	})

	t.Run("list-prompts", func(t *testing.T) {
		res, err := cs.CallTool(ctx, &mcp.CallToolParams{Name: "list-prompts"})
		if err != nil {
			t.Fatalf("CallTool: %v", err)
		}
		raw, _ := json.Marshal(res.StructuredContent)
		var out struct {
			Entries []struct {
				URI         string `json:"uri"`
				Source      string `json:"source"`
				Path        string `json:"path"`
				Type        string `json:"type"`
				Description string `json:"description"`
				Mode        string `json:"mode,omitempty"`
			} `json:"entries"`
		}
		if err := json.Unmarshal(raw, &out); err != nil {
			t.Fatalf("unmarshal: %v", err)
		}
		if len(out.Entries) != 2 {
			t.Fatalf("expected 2 entries, got %d", len(out.Entries))
		}
		uris := map[string]bool{}
		types := map[string]string{}
		for _, e := range out.Entries {
			uris[e.URI] = true
			types[e.URI] = e.Type
		}
		compURI := "prompts://" + src + "/component"
		reviewURI := "prompts://" + src + "/reviewer"
		if !uris[compURI] {
			t.Errorf("missing URI %s", compURI)
		}
		if !uris[reviewURI] {
			t.Errorf("missing URI %s", reviewURI)
		}
		if types[compURI] != loader.TypePrompt {
			t.Errorf("component type = %q, want %q", types[compURI], loader.TypePrompt)
		}
		if types[reviewURI] != loader.TypeChatmode {
			t.Errorf("reviewer type = %q, want %q", types[reviewURI], loader.TypeChatmode)
		}
	})

	t.Run("get-prompt found", func(t *testing.T) {
		uri := "prompts://" + src + "/component"
		res, err := cs.CallTool(ctx, &mcp.CallToolParams{
			Name:      "get-prompt",
			Arguments: map[string]any{"uri": uri},
		})
		if err != nil {
			t.Fatalf("CallTool: %v", err)
		}
		if res.IsError {
			t.Fatal("returned error")
		}
		raw, _ := json.Marshal(res.StructuredContent)
		var out struct {
			URI         string `json:"uri"`
			Name        string `json:"name"`
			Description string `json:"description"`
			Mode        string `json:"mode"`
			Type        string `json:"type"`
			Content     string `json:"content"`
		}
		if err := json.Unmarshal(raw, &out); err != nil {
			t.Fatalf("unmarshal: %v", err)
		}
		if out.URI != uri {
			t.Errorf("URI = %q, want %q", out.URI, uri)
		}
		if out.Name != "component" {
			t.Errorf("Name = %q", out.Name)
		}
		if out.Mode != "agent" {
			t.Errorf("Mode = %q, want agent", out.Mode)
		}
		if out.Type != loader.TypePrompt {
			t.Errorf("Type = %q", out.Type)
		}
		if !strings.Contains(out.Content, "React component") {
			t.Errorf("Content missing expected text: %q", out.Content)
		}
	})

	t.Run("get-prompt not found", func(t *testing.T) {
		res, err := cs.CallTool(ctx, &mcp.CallToolParams{
			Name:      "get-prompt",
			Arguments: map[string]any{"uri": "prompts://nonexistent/nothere"},
		})
		if err != nil {
			t.Fatalf("CallTool: %v", err)
		}
		if !res.IsError {
			t.Error("expected IsError=true for missing prompt")
		}
	})

	t.Run("optimize-prompts", func(t *testing.T) {
		res, err := cs.CallTool(ctx, &mcp.CallToolParams{
			Name:      "optimize-prompts",
			Arguments: map[string]any{},
		})
		if err != nil {
			t.Fatalf("CallTool: %v", err)
		}
		raw, _ := json.Marshal(res.StructuredContent)
		var out struct {
			Content string `json:"content"`
			Matched int    `json:"matched"`
		}
		if err := json.Unmarshal(raw, &out); err != nil {
			t.Fatalf("unmarshal: %v", err)
		}
		if out.Matched != 2 {
			t.Errorf("matched = %d, want 2", out.Matched)
		}
		for _, want := range []string{"React component", "Review code"} {
			if !strings.Contains(out.Content, want) {
				t.Errorf("content missing %q", want)
			}
		}
	})

	t.Run("optimize-prompts with source filter", func(t *testing.T) {
		res, err := cs.CallTool(ctx, &mcp.CallToolParams{
			Name:      "optimize-prompts",
			Arguments: map[string]any{"source": "nonexistent"},
		})
		if err != nil {
			t.Fatalf("CallTool: %v", err)
		}
		raw, _ := json.Marshal(res.StructuredContent)
		var out struct {
			Matched int    `json:"matched"`
		}
		json.Unmarshal(raw, &out)
		if out.Matched != 0 {
			t.Errorf("matched = %d, want 0", out.Matched)
		}
	})
}

// ---------------------------------------------------------------------------
// 8. TestRegisterPrompts_ViaProtocol
// ---------------------------------------------------------------------------

func TestRegisterPrompts_ViaProtocol(t *testing.T) {
	ctx := context.Background()
	ldr, opt, cleanup := setupTestEnv(t)
	defer cleanup()

	srv := mcp.NewServer(&mcp.Implementation{Name: "mcp-prompts", Version: "test"}, nil)
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
			if p.Name == "get-prompts" {
				found = true
			}
		}
		if !found {
			t.Error("get-prompts prompt not found")
		}
	})

	t.Run("get-prompts without filter", func(t *testing.T) {
		res, err := cs.GetPrompt(ctx, &mcp.GetPromptParams{Name: "get-prompts"})
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
		for _, want := range []string{"React component", "Review code"} {
			if !strings.Contains(tc.Text, want) {
				t.Errorf("prompt missing %q", want)
			}
		}
	})

	t.Run("get-prompts with non-matching source", func(t *testing.T) {
		res, err := cs.GetPrompt(ctx, &mcp.GetPromptParams{
			Name:      "get-prompts",
			Arguments: map[string]string{"source": "nonexistent"},
		})
		if err != nil {
			t.Fatalf("GetPrompt: %v", err)
		}
		tc, ok := res.Messages[0].Content.(*mcp.TextContent)
		if !ok {
			t.Fatalf("expected TextContent, got %T", res.Messages[0].Content)
		}
		if tc.Text != "No prompts found." {
			t.Errorf("expected 'No prompts found.', got %q", tc.Text)
		}
	})

	t.Run("get-prompts with matching source", func(t *testing.T) {
		src := sourceOf(ldr)
		res, err := cs.GetPrompt(ctx, &mcp.GetPromptParams{
			Name:      "get-prompts",
			Arguments: map[string]string{"source": src},
		})
		if err != nil {
			t.Fatalf("GetPrompt: %v", err)
		}
		tc, ok := res.Messages[0].Content.(*mcp.TextContent)
		if !ok {
			t.Fatalf("expected TextContent, got %T", res.Messages[0].Content)
		}
		if !strings.Contains(tc.Text, "React component") {
			t.Errorf("expected prompt content, got: %q", tc.Text)
		}
	})
}
