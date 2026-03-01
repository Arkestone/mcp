package main

import (
	"context"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/Arkestone/mcp/pkg/config"
	"github.com/Arkestone/mcp/pkg/github"
	"github.com/Arkestone/mcp/pkg/optimizer"
	"github.com/Arkestone/mcp/pkg/server"
	"github.com/Arkestone/mcp/servers/mcp-adr/internal/scanner"
	"github.com/modelcontextprotocol/go-sdk/mcp"
)

// setupTestADREnv creates a temp dir with two ADR files and returns scanner + optimizer.
func setupTestADREnv(t *testing.T) (*scanner.Scanner, *optimizer.Optimizer) {
	t.Helper()
	dir := t.TempDir()

	adrDir := filepath.Join(dir, "docs", "adr")
	if err := os.MkdirAll(adrDir, 0o755); err != nil {
		t.Fatal(err)
	}

	adr1 := "---\ntitle: Use Go\nstatus: accepted\ndate: 2023-01-01\n---\nWe use Go for all services."
	if err := os.WriteFile(filepath.Join(adrDir, "0001-use-go.md"), []byte(adr1), 0o644); err != nil {
		t.Fatal(err)
	}

	adr2 := "---\ntitle: Use PostgreSQL\nstatus: proposed\ndate: 2023-02-01\n---\nWe use PostgreSQL for persistence."
	if err := os.WriteFile(filepath.Join(adrDir, "0002-use-postgresql.md"), []byte(adr2), 0o644); err != nil {
		t.Fatal(err)
	}

	cfg := &config.Config{
		Sources: config.Sources{Dirs: []string{dir}},
		Cache:   config.CacheConfig{Dir: t.TempDir(), SyncInterval: 5 * time.Minute},
	}
	scn := scanner.New(cfg, &github.Client{})
	scn.Start(context.Background())
	t.Cleanup(scn.Stop)
	return scn, optimizer.New(cfg.LLM)
}

func connectClientServer(t *testing.T, ctx context.Context, srv *mcp.Server) *mcp.ClientSession {
	t.Helper()
	st, ct := mcp.NewInMemoryTransports()
	go srv.Run(ctx, st)
	client := mcp.NewClient(&mcp.Implementation{Name: "test-client", Version: "0.0.1"}, nil)
	cs, err := client.Connect(ctx, ct, nil)
	if err != nil {
		t.Fatalf("client.Connect: %v", err)
	}
	t.Cleanup(func() { cs.Close() })
	return cs
}

// ---------------------------------------------------------------------------
// Unit tests
// ---------------------------------------------------------------------------

func TestToOptimizerInputs(t *testing.T) {
	adrs := []scanner.ADR{
		{Source: "local", Path: "docs/adr/0001-use-go.md", Content: "Use Go"},
		{Source: "repo", Path: "docs/adr/0002-use-psql.md", Content: "Use PostgreSQL"},
	}
	inputs := toOptimizerInputs(adrs)
	if len(inputs) != 2 {
		t.Fatalf("got %d inputs, want 2", len(inputs))
	}
	if inputs[0].Source != "local" || inputs[0].Path != "docs/adr/0001-use-go.md" || inputs[0].Content != "Use Go" {
		t.Errorf("inputs[0] = %+v", inputs[0])
	}
	if inputs[1].Source != "repo" {
		t.Errorf("inputs[1].Source = %q, want repo", inputs[1].Source)
	}
}

func TestToOptimizerInputs_Empty(t *testing.T) {
	if len(toOptimizerInputs(nil)) != 0 {
		t.Error("nil slice should return empty inputs")
	}
	if len(toOptimizerInputs([]scanner.ADR{})) != 0 {
		t.Error("empty slice should return empty inputs")
	}
}

func TestFilterBySource(t *testing.T) {
	adrs := []scanner.ADR{
		{ID: "A", Source: "local"},
		{ID: "B", Source: "remote"},
		{ID: "C", Source: "local"},
	}
	t.Run("empty returns all", func(t *testing.T) {
		if len(filterBySource(adrs, "")) != 3 {
			t.Error("empty source should return all")
		}
	})
	t.Run("filter local", func(t *testing.T) {
		got := filterBySource(adrs, "local")
		if len(got) != 2 {
			t.Fatalf("got %d, want 2", len(got))
		}
		for _, a := range got {
			if a.Source != "local" {
				t.Errorf("unexpected source %q", a.Source)
			}
		}
	})
	t.Run("no match", func(t *testing.T) {
		if len(filterBySource(adrs, "nope")) != 0 {
			t.Error("no-match source should return empty")
		}
	})
}

func TestFilterByStatus(t *testing.T) {
	adrs := []scanner.ADR{
		{ID: "A", Status: "accepted"},
		{ID: "B", Status: "proposed"},
		{ID: "C", Status: "Accepted"}, // mixed case
	}
	t.Run("empty returns all", func(t *testing.T) {
		if len(filterByStatus(adrs, "")) != 3 {
			t.Error("empty status should return all")
		}
	})
	t.Run("filter accepted case-insensitive", func(t *testing.T) {
		got := filterByStatus(adrs, "accepted")
		if len(got) != 2 {
			t.Fatalf("got %d, want 2 (case-insensitive match)", len(got))
		}
	})
	t.Run("filter proposed", func(t *testing.T) {
		got := filterByStatus(adrs, "proposed")
		if len(got) != 1 {
			t.Fatalf("got %d, want 1", len(got))
		}
		if got[0].ID != "B" {
			t.Errorf("ID = %q, want B", got[0].ID)
		}
	})
	t.Run("no match", func(t *testing.T) {
		if len(filterByStatus(adrs, "deprecated")) != 0 {
			t.Error("no-match status should return empty")
		}
	})
}

func TestPromptResult(t *testing.T) {
	res := promptResult("Architecture Decision Records", "some content")
	if res.Description != "Architecture Decision Records" {
		t.Errorf("Description = %q", res.Description)
	}
	if len(res.Messages) != 1 {
		t.Fatalf("got %d messages, want 1", len(res.Messages))
	}
	tc, ok := res.Messages[0].Content.(*mcp.TextContent)
	if !ok {
		t.Fatalf("expected *mcp.TextContent, got %T", res.Messages[0].Content)
	}
	if tc.Text != "some content" {
		t.Errorf("Text = %q, want %q", tc.Text, "some content")
	}
	if res.Messages[0].Role != "user" {
		t.Errorf("Role = %q, want user", res.Messages[0].Role)
	}
}

func TestShouldOptimize_NoOptimizer(t *testing.T) {
	var opt *optimizer.Optimizer
	if server.ShouldOptimize(opt, true, "true") {
		t.Error("nil optimizer should always return false")
	}
}

// ---------------------------------------------------------------------------
// Protocol tests
// ---------------------------------------------------------------------------

func TestRegisterResources_ViaProtocol(t *testing.T) {
	scn, opt := setupTestADREnv(t)
	ctx := context.Background()

	srv := mcp.NewServer(&mcp.Implementation{Name: "test-res", Version: "0.1.0"}, nil)
	registerResources(srv, scn, opt, false)
	cs := connectClientServer(t, ctx, srv)

	t.Run("read adrs://index", func(t *testing.T) {
		res, err := cs.ReadResource(ctx, &mcp.ReadResourceParams{URI: "adrs://index"})
		if err != nil {
			t.Fatalf("ReadResource(index): %v", err)
		}
		if len(res.Contents) == 0 {
			t.Fatal("no contents")
		}
		text := res.Contents[0].Text
		if !strings.Contains(text, "Use Go") {
			t.Errorf("missing Use Go in index: %s", text)
		}
		if !strings.Contains(text, "Use PostgreSQL") {
			t.Errorf("missing Use PostgreSQL in index: %s", text)
		}
	})

	t.Run("read adrs://optimized", func(t *testing.T) {
		res, err := cs.ReadResource(ctx, &mcp.ReadResourceParams{URI: "adrs://optimized"})
		if err != nil {
			t.Fatalf("ReadResource(optimized): %v", err)
		}
		if len(res.Contents) == 0 {
			t.Fatal("no contents")
		}
		text := res.Contents[0].Text
		if !strings.Contains(text, "We use Go") {
			t.Errorf("missing Go content: %s", text)
		}
	})

	t.Run("read individual ADR via template URI", func(t *testing.T) {
		adrs := scn.List()
		if len(adrs) == 0 {
			t.Fatal("no ADRs loaded")
		}
		uri := adrs[0].URI
		res, err := cs.ReadResource(ctx, &mcp.ReadResourceParams{URI: uri})
		if err != nil {
			t.Fatalf("ReadResource(%s): %v", uri, err)
		}
		if len(res.Contents) == 0 {
			t.Fatal("no contents")
		}
		text := res.Contents[0].Text
		if !strings.Contains(text, "We use Go") {
			t.Errorf("missing content: %s", text)
		}
	})

	t.Run("read nonexistent ADR", func(t *testing.T) {
		_, err := cs.ReadResource(ctx, &mcp.ReadResourceParams{URI: "adrs://nosource/nonexistent"})
		if err == nil {
			t.Error("expected error for nonexistent ADR")
		}
	})
}

func TestRegisterTools_ViaProtocol(t *testing.T) {
	scn, opt := setupTestADREnv(t)
	ctx := context.Background()

	srv := mcp.NewServer(&mcp.Implementation{Name: "test-tools", Version: "0.1.0"}, nil)
	registerTools(srv, scn, opt, false)
	cs := connectClientServer(t, ctx, srv)

	t.Run("refresh-adrs", func(t *testing.T) {
		res, err := cs.CallTool(ctx, &mcp.CallToolParams{Name: "refresh-adrs"})
		if err != nil {
			t.Fatalf("CallTool(refresh-adrs): %v", err)
		}
		if len(res.Content) == 0 {
			t.Fatal("no content")
		}
		tc, ok := res.Content[0].(*mcp.TextContent)
		if !ok {
			t.Fatalf("expected TextContent, got %T", res.Content[0])
		}
		if !strings.Contains(tc.Text, "refreshed") {
			t.Errorf("unexpected result: %s", tc.Text)
		}
	})

	t.Run("list-adrs", func(t *testing.T) {
		res, err := cs.CallTool(ctx, &mcp.CallToolParams{Name: "list-adrs"})
		if err != nil {
			t.Fatalf("CallTool(list-adrs): %v", err)
		}
		if len(res.Content) == 0 {
			t.Fatal("no content")
		}
		tc, ok := res.Content[0].(*mcp.TextContent)
		if !ok {
			t.Fatalf("expected TextContent, got %T", res.Content[0])
		}
		if !strings.Contains(tc.Text, "Use Go") {
			t.Errorf("missing Use Go: %s", tc.Text)
		}
		if !strings.Contains(tc.Text, "Use PostgreSQL") {
			t.Errorf("missing Use PostgreSQL: %s", tc.Text)
		}
	})

	t.Run("list-adrs with status filter", func(t *testing.T) {
		res, err := cs.CallTool(ctx, &mcp.CallToolParams{
			Name:      "list-adrs",
			Arguments: map[string]any{"status": "accepted"},
		})
		if err != nil {
			t.Fatalf("CallTool(list-adrs status=accepted): %v", err)
		}
		if len(res.Content) == 0 {
			t.Fatal("no content")
		}
		tc, ok := res.Content[0].(*mcp.TextContent)
		if !ok {
			t.Fatalf("expected TextContent, got %T", res.Content[0])
		}
		if !strings.Contains(tc.Text, "Use Go") {
			t.Errorf("missing Use Go (accepted): %s", tc.Text)
		}
		if strings.Contains(tc.Text, "Use PostgreSQL") {
			t.Errorf("should not contain proposed ADR: %s", tc.Text)
		}
	})

	t.Run("get-adr", func(t *testing.T) {
		adrs := scn.List()
		if len(adrs) == 0 {
			t.Fatal("no ADRs")
		}
		uri := adrs[0].URI
		res, err := cs.CallTool(ctx, &mcp.CallToolParams{
			Name:      "get-adr",
			Arguments: map[string]any{"uri": uri},
		})
		if err != nil {
			t.Fatalf("CallTool(get-adr): %v", err)
		}
		if len(res.Content) == 0 {
			t.Fatal("no content")
		}
		tc, ok := res.Content[0].(*mcp.TextContent)
		if !ok {
			t.Fatalf("expected TextContent, got %T", res.Content[0])
		}
		if !strings.Contains(tc.Text, "We use Go") {
			t.Errorf("missing content: %s", tc.Text)
		}
	})

	t.Run("get-adr not found", func(t *testing.T) {
		res, err := cs.CallTool(ctx, &mcp.CallToolParams{
			Name:      "get-adr",
			Arguments: map[string]any{"uri": "adrs://nosource/nonexistent"},
		})
		if err != nil {
			return // error path is fine
		}
		if res == nil || !res.IsError {
			t.Error("expected error result for nonexistent ADR")
		}
	})
}

func TestRegisterPrompts_ViaProtocol(t *testing.T) {
	scn, opt := setupTestADREnv(t)
	ctx := context.Background()

	srv := mcp.NewServer(&mcp.Implementation{Name: "test-prompts", Version: "0.1.0"}, nil)
	registerPrompts(srv, scn, opt, false)
	cs := connectClientServer(t, ctx, srv)

	t.Run("get-adrs prompt", func(t *testing.T) {
		res, err := cs.GetPrompt(ctx, &mcp.GetPromptParams{
			Name:      "get-adrs",
			Arguments: map[string]string{},
		})
		if err != nil {
			t.Fatalf("GetPrompt(get-adrs): %v", err)
		}
		if res.Description != "Architecture Decision Records" {
			t.Errorf("Description = %q", res.Description)
		}
		if len(res.Messages) == 0 {
			t.Fatal("no messages")
		}
		tc, ok := res.Messages[0].Content.(*mcp.TextContent)
		if !ok {
			t.Fatalf("expected TextContent, got %T", res.Messages[0].Content)
		}
		if !strings.Contains(tc.Text, "We use Go") {
			t.Errorf("missing Go content: %s", tc.Text)
		}
		if !strings.Contains(tc.Text, "We use PostgreSQL") {
			t.Errorf("missing PostgreSQL content: %s", tc.Text)
		}
	})

	t.Run("get-adrs prompt with source filter", func(t *testing.T) {
		adrs := scn.List()
		if len(adrs) == 0 {
			t.Fatal("no ADRs")
		}
		source := adrs[0].Source
		res, err := cs.GetPrompt(ctx, &mcp.GetPromptParams{
			Name:      "get-adrs",
			Arguments: map[string]string{"source": source},
		})
		if err != nil {
			t.Fatalf("GetPrompt(get-adrs with source): %v", err)
		}
		if len(res.Messages) == 0 {
			t.Fatal("no messages")
		}
	})

	t.Run("get-adrs prompt with status filter", func(t *testing.T) {
		res, err := cs.GetPrompt(ctx, &mcp.GetPromptParams{
			Name:      "get-adrs",
			Arguments: map[string]string{"status": "accepted"},
		})
		if err != nil {
			t.Fatalf("GetPrompt(get-adrs status=accepted): %v", err)
		}
		if len(res.Messages) == 0 {
			t.Fatal("no messages")
		}
		tc, ok := res.Messages[0].Content.(*mcp.TextContent)
		if !ok {
			t.Fatalf("expected TextContent, got %T", res.Messages[0].Content)
		}
		if !strings.Contains(tc.Text, "We use Go") {
			t.Errorf("missing accepted ADR content: %s", tc.Text)
		}
	})

	t.Run("get-adrs prompt no match returns no ADRs found", func(t *testing.T) {
		res, err := cs.GetPrompt(ctx, &mcp.GetPromptParams{
			Name:      "get-adrs",
			Arguments: map[string]string{"source": "nonexistent-source"},
		})
		if err != nil {
			t.Fatalf("GetPrompt: %v", err)
		}
		if len(res.Messages) == 0 {
			t.Fatal("no messages")
		}
		tc, ok := res.Messages[0].Content.(*mcp.TextContent)
		if !ok {
			t.Fatalf("expected TextContent, got %T", res.Messages[0].Content)
		}
		if !strings.Contains(tc.Text, "No ADRs found") {
			t.Errorf("expected 'No ADRs found', got: %s", tc.Text)
		}
	})
}
