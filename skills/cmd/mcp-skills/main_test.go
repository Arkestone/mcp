package main

import (
	"context"
	"encoding/json"
	"fmt"
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
	"github.com/Arkestone/mcp/skills/internal/scanner"
	"github.com/modelcontextprotocol/go-sdk/mcp"
)

func newTestOptimizer(t *testing.T) *optimizer.Optimizer {
	t.Helper()
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		resp := struct {
			Choices []struct {
				Message struct {
					Role    string `json:"role"`
					Content string `json:"content"`
				} `json:"message"`
			} `json:"choices"`
		}{
			Choices: []struct {
				Message struct {
					Role    string `json:"role"`
					Content string `json:"content"`
				} `json:"message"`
			}{
				{Message: struct {
					Role    string `json:"role"`
					Content string `json:"content"`
				}{Role: "assistant", Content: "optimized"}},
			},
		}
		json.NewEncoder(w).Encode(resp)
	}))
	t.Cleanup(srv.Close)

	cfg := optimizer.TestLLMConfig()
	cfg.Endpoint = srv.URL
	return optimizer.New(cfg)
}

func TestShouldOptimize_NoOptimizer(t *testing.T) {
	var opt *optimizer.Optimizer
	if server.ShouldOptimize(opt, true, "true") {
		t.Error("nil optimizer should always return false")
	}
	if server.ShouldOptimize(opt, true, "") {
		t.Error("nil optimizer should always return false")
	}
}

func TestShouldOptimize_TrueOverride(t *testing.T) {
	opt := newTestOptimizer(t)
	if !server.ShouldOptimize(opt, false, "true") {
		t.Error("per-request 'true' should override global false")
	}
	if !server.ShouldOptimize(opt, false, "1") {
		t.Error("per-request '1' should override global false")
	}
	if !server.ShouldOptimize(opt, false, "yes") {
		t.Error("per-request 'yes' should override global false")
	}
}

func TestShouldOptimize_FalseOverride(t *testing.T) {
	opt := newTestOptimizer(t)
	if server.ShouldOptimize(opt, true, "false") {
		t.Error("per-request 'false' should override global true")
	}
	if server.ShouldOptimize(opt, true, "0") {
		t.Error("per-request '0' should override global true")
	}
	if server.ShouldOptimize(opt, true, "no") {
		t.Error("per-request 'no' should override global true")
	}
}

func TestShouldOptimize_DefaultTrue(t *testing.T) {
	opt := newTestOptimizer(t)
	if !server.ShouldOptimize(opt, true, "") {
		t.Error("empty override with global true should return true")
	}
}

func TestShouldOptimize_DefaultFalse(t *testing.T) {
	opt := newTestOptimizer(t)
	if server.ShouldOptimize(opt, false, "") {
		t.Error("empty override with global false should return false")
	}
}

func TestShouldOptimize_CaseInsensitive(t *testing.T) {
	opt := newTestOptimizer(t)
	for _, v := range []string{"TRUE", "True", "YES"} {
		if !server.ShouldOptimize(opt, false, v) {
			t.Errorf("per-request %q should be treated as true", v)
		}
	}
}

func TestToOptimizerInputs_Basic(t *testing.T) {
	skills := []scanner.Skill{
		{Source: "local", Path: "docker/SKILL.md", Content: "Use multi-stage builds"},
		{Source: "repo", Path: "python/SKILL.md", Content: "Use type hints"},
	}
	inputs := toOptimizerInputs(skills)
	if len(inputs) != 2 {
		t.Fatalf("got %d inputs, want 2", len(inputs))
	}
	if inputs[0].Source != "local" || inputs[0].Path != "docker/SKILL.md" || inputs[0].Content != "Use multi-stage builds" {
		t.Errorf("inputs[0] = %+v", inputs[0])
	}
	if inputs[1].Source != "repo" || inputs[1].Path != "python/SKILL.md" || inputs[1].Content != "Use type hints" {
		t.Errorf("inputs[1] = %+v", inputs[1])
	}
}

func TestToOptimizerInputs_Empty(t *testing.T) {
	inputs := toOptimizerInputs(nil)
	if len(inputs) != 0 {
		t.Errorf("got %d inputs for nil slice, want 0", len(inputs))
	}
	inputs = toOptimizerInputs([]scanner.Skill{})
	if len(inputs) != 0 {
		t.Errorf("got %d inputs for empty slice, want 0", len(inputs))
	}
}

// --------------- helpers ---------------

func setupTestSkillEnv(t *testing.T) (*scanner.Scanner, *optimizer.Optimizer) {
	t.Helper()
	dir := t.TempDir()

	// Skill 1: docker-expert (no references)
	dockerDir := filepath.Join(dir, "docker-expert")
	if err := os.MkdirAll(dockerDir, 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(dockerDir, "SKILL.md"),
		[]byte("---\nname: Docker Expert\ndescription: Docker skills\n---\nUse multi-stage builds"), 0o644); err != nil {
		t.Fatal(err)
	}

	// Skill 2: python-dev (with references)
	pythonDir := filepath.Join(dir, "python-dev")
	refsDir := filepath.Join(pythonDir, "references")
	if err := os.MkdirAll(refsDir, 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(pythonDir, "SKILL.md"),
		[]byte("---\nname: Python Developer\ndescription: Python skills\n---\nUse type hints"), 0o644); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(refsDir, "pep8.md"),
		[]byte("PEP 8 style guide content"), 0o644); err != nil {
		t.Fatal(err)
	}

	// Skill 3: go-expert (no spaces in name, for template URI matching)
	goDir := filepath.Join(dir, "go-expert")
	if err := os.MkdirAll(goDir, 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(goDir, "SKILL.md"),
		[]byte("---\nname: GoExpert\ndescription: Go skills\n---\nUse gofmt"), 0o644); err != nil {
		t.Fatal(err)
	}

	cfg := &config.Config{
		Sources: config.Sources{Dirs: []string{dir}},
		Cache:   config.CacheConfig{Dir: t.TempDir(), SyncInterval: 5 * time.Minute},
	}
	scn := scanner.New(cfg, &github.Client{})
	scn.Start(context.Background())
	t.Cleanup(scn.Stop)
	return scn, optimizer.New(cfg.LLM) // nil optimizer since LLM not configured
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

// --------------- TestFilterBySource ---------------

func TestFilterBySource(t *testing.T) {
	skills := []scanner.Skill{
		{Name: "A", Source: "local"},
		{Name: "B", Source: "remote"},
		{Name: "C", Source: "local"},
	}

	t.Run("empty source returns all", func(t *testing.T) {
		got := filterBySource(skills, "")
		if len(got) != 3 {
			t.Fatalf("got %d, want 3", len(got))
		}
	})
	t.Run("specific source", func(t *testing.T) {
		got := filterBySource(skills, "local")
		if len(got) != 2 {
			t.Fatalf("got %d, want 2", len(got))
		}
		for _, s := range got {
			if s.Source != "local" {
				t.Errorf("got source %q, want local", s.Source)
			}
		}
	})
	t.Run("non-matching source", func(t *testing.T) {
		got := filterBySource(skills, "nope")
		if len(got) != 0 {
			t.Fatalf("got %d, want 0", len(got))
		}
	})
}

// --------------- TestPromptResult ---------------

func TestPromptResult(t *testing.T) {
	res := promptResult("my description", "hello world")
	if res.Description != "my description" {
		t.Errorf("Description = %q, want %q", res.Description, "my description")
	}
	if len(res.Messages) != 1 {
		t.Fatalf("got %d messages, want 1", len(res.Messages))
	}
	tc, ok := res.Messages[0].Content.(*mcp.TextContent)
	if !ok {
		t.Fatalf("expected *mcp.TextContent, got %T", res.Messages[0].Content)
	}
	if tc.Text != "hello world" {
		t.Errorf("Text = %q, want %q", tc.Text, "hello world")
	}
	if res.Messages[0].Role != "user" {
		t.Errorf("Role = %q, want user", res.Messages[0].Role)
	}
}

// --------------- TestRenderSkill ---------------

func TestRenderSkill(t *testing.T) {
	t.Run("no references", func(t *testing.T) {
		s := scanner.Skill{Name: "Go", Description: "Go skills", Content: "Use gofmt"}
		got := renderSkill(s)
		if !strings.Contains(got, "---\nname: Go\ndescription: Go skills\n---") {
			t.Errorf("missing frontmatter: %s", got)
		}
		if !strings.Contains(got, "Use gofmt") {
			t.Errorf("missing content: %s", got)
		}
		if strings.Contains(got, "## References") {
			t.Error("should not contain References section")
		}
	})

	t.Run("with references", func(t *testing.T) {
		s := scanner.Skill{
			Name: "Go", Description: "Go skills", Content: "Use gofmt",
			References: []scanner.Reference{
				{Name: "style.md", Content: "style guide"},
			},
		}
		got := renderSkill(s)
		if !strings.Contains(got, "## References") {
			t.Error("missing References section")
		}
		if !strings.Contains(got, "### style.md") {
			t.Error("missing reference heading")
		}
		if !strings.Contains(got, "style guide") {
			t.Error("missing reference content")
		}
	})
}

// --------------- TestOptimizeContent ---------------

func TestOptimizeContent(t *testing.T) {
	inputs := []optimizer.ContentInput{
		{Source: "local", Path: "a/SKILL.md", Content: "skill A"},
		{Source: "remote", Path: "b/SKILL.md", Content: "skill B"},
	}

	t.Run("nil optimizer falls back to ConcatRaw", func(t *testing.T) {
		got := optimizeContent(context.Background(), nil, true, "true", inputs)
		want := optimizer.ConcatRaw(inputs)
		if got != want {
			t.Errorf("got %q, want %q", got, want)
		}
	})

	t.Run("enabled optimizer override false falls back", func(t *testing.T) {
		opt := newTestOptimizer(t)
		got := optimizeContent(context.Background(), opt, true, "false", inputs)
		want := optimizer.ConcatRaw(inputs)
		if got != want {
			t.Errorf("got %q, want %q", got, want)
		}
	})

	t.Run("optimizer error falls back to ConcatRaw", func(t *testing.T) {
		errSrv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusInternalServerError)
		}))
		t.Cleanup(errSrv.Close)
		cfg := optimizer.TestLLMConfig()
		cfg.Endpoint = errSrv.URL
		opt := optimizer.New(cfg)
		got := optimizeContent(context.Background(), opt, true, "true", inputs)
		want := optimizer.ConcatRaw(inputs)
		if got != want {
			t.Errorf("got %q, want %q", got, want)
		}
	})

	t.Run("optimizer success", func(t *testing.T) {
		opt := newTestOptimizer(t)
		got := optimizeContent(context.Background(), opt, true, "true", inputs)
		if got != "optimized" {
			t.Errorf("got %q, want %q", got, "optimized")
		}
	})
}

// --------------- TestRegisterResources_ViaProtocol ---------------

func TestRegisterResources_ViaProtocol(t *testing.T) {
	scn, opt := setupTestSkillEnv(t)
	ctx := context.Background()

	srv := mcp.NewServer(&mcp.Implementation{Name: "test-res", Version: "0.1.0"}, nil)
	registerResources(srv, scn, opt, false)
	cs := connectClientServer(t, ctx, srv)

	t.Run("read skills://index", func(t *testing.T) {
		res, err := cs.ReadResource(ctx, &mcp.ReadResourceParams{URI: "skills://index"})
		if err != nil {
			t.Fatalf("ReadResource(index): %v", err)
		}
		if len(res.Contents) == 0 {
			t.Fatal("no contents")
		}
		text := res.Contents[0].Text
		if !strings.Contains(text, "Docker Expert") {
			t.Errorf("missing Docker Expert in index: %s", text)
		}
		if !strings.Contains(text, "Python Developer") {
			t.Errorf("missing Python Developer in index: %s", text)
		}
	})

	t.Run("read skills://optimized", func(t *testing.T) {
		res, err := cs.ReadResource(ctx, &mcp.ReadResourceParams{URI: "skills://optimized"})
		if err != nil {
			t.Fatalf("ReadResource(optimized): %v", err)
		}
		if len(res.Contents) == 0 {
			t.Fatal("no contents")
		}
		text := res.Contents[0].Text
		if !strings.Contains(text, "multi-stage builds") {
			t.Errorf("missing docker content: %s", text)
		}
	})

	t.Run("read nonexistent skill via template URI", func(t *testing.T) {
		_, err := cs.ReadResource(ctx, &mcp.ReadResourceParams{URI: "skills://nosource/NoSkill"})
		if err == nil {
			t.Error("expected error for nonexistent skill")
		}
	})

	t.Run("read individual skill via template URI", func(t *testing.T) {
		// Use GoExpert (no spaces) which matches the URI template regex
		skills := scn.List()
		if len(skills) == 0 {
			t.Fatal("no skills loaded")
		}
		var uri string
		for _, s := range skills {
			if s.Name == "GoExpert" {
				uri = s.URI
				break
			}
		}
		if uri == "" {
			t.Fatal("GoExpert not found")
		}
		res, err := cs.ReadResource(ctx, &mcp.ReadResourceParams{URI: uri})
		if err != nil {
			t.Fatalf("ReadResource(%s): %v", uri, err)
		}
		if len(res.Contents) == 0 {
			t.Fatal("no contents")
		}
		text := res.Contents[0].Text
		if !strings.Contains(text, "Use gofmt") {
			t.Errorf("missing go content: %s", text)
		}
		if !strings.Contains(text, "name: GoExpert") {
			t.Errorf("missing frontmatter: %s", text)
		}
	})
}

// --------------- TestRegisterTools_ViaProtocol ---------------

func TestRegisterTools_ViaProtocol(t *testing.T) {
	scn, opt := setupTestSkillEnv(t)
	ctx := context.Background()

	srv := mcp.NewServer(&mcp.Implementation{Name: "test-tools", Version: "0.1.0"}, nil)
	registerTools(srv, scn, opt, false)
	cs := connectClientServer(t, ctx, srv)

	t.Run("refresh-skills", func(t *testing.T) {
		res, err := cs.CallTool(ctx, &mcp.CallToolParams{Name: "refresh-skills"})
		if err != nil {
			t.Fatalf("CallTool(refresh-skills): %v", err)
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

	t.Run("list-skills", func(t *testing.T) {
		res, err := cs.CallTool(ctx, &mcp.CallToolParams{Name: "list-skills"})
		if err != nil {
			t.Fatalf("CallTool(list-skills): %v", err)
		}
		if len(res.Content) == 0 {
			t.Fatal("no content")
		}
		tc, ok := res.Content[0].(*mcp.TextContent)
		if !ok {
			t.Fatalf("expected TextContent, got %T", res.Content[0])
		}
		if !strings.Contains(tc.Text, "Docker Expert") {
			t.Errorf("missing Docker Expert: %s", tc.Text)
		}
		if !strings.Contains(tc.Text, "Python Developer") {
			t.Errorf("missing Python Developer: %s", tc.Text)
		}
	})

	t.Run("get-skill", func(t *testing.T) {
		res, err := cs.CallTool(ctx, &mcp.CallToolParams{
			Name:      "get-skill",
			Arguments: map[string]any{"name": "Docker Expert"},
		})
		if err != nil {
			t.Fatalf("CallTool(get-skill): %v", err)
		}
		if len(res.Content) == 0 {
			t.Fatal("no content")
		}
		tc, ok := res.Content[0].(*mcp.TextContent)
		if !ok {
			t.Fatalf("expected TextContent, got %T", res.Content[0])
		}
		if !strings.Contains(tc.Text, "multi-stage builds") {
			t.Errorf("missing docker content: %s", tc.Text)
		}
	})

	t.Run("get-skill not found", func(t *testing.T) {
		res, err := cs.CallTool(ctx, &mcp.CallToolParams{
			Name:      "get-skill",
			Arguments: map[string]any{"name": "Nonexistent"},
		})
		// The SDK may return the error as isError in the result or as a Go error
		if err != nil {
			return // error path is fine
		}
		if res == nil || !res.IsError {
			t.Error("expected error result for nonexistent skill")
		}
	})

	t.Run("optimize-skills", func(t *testing.T) {
		res, err := cs.CallTool(ctx, &mcp.CallToolParams{
			Name:      "optimize-skills",
			Arguments: map[string]any{},
		})
		if err != nil {
			t.Fatalf("CallTool(optimize-skills): %v", err)
		}
		if len(res.Content) == 0 {
			t.Fatal("no content")
		}
		tc, ok := res.Content[0].(*mcp.TextContent)
		if !ok {
			t.Fatalf("expected TextContent, got %T", res.Content[0])
		}
		if !strings.Contains(tc.Text, "multi-stage builds") {
			t.Errorf("missing docker content: %s", tc.Text)
		}
	})
}

// --------------- TestRegisterPrompts_ViaProtocol ---------------

func TestRegisterPrompts_ViaProtocol(t *testing.T) {
	scn, opt := setupTestSkillEnv(t)
	ctx := context.Background()

	srv := mcp.NewServer(&mcp.Implementation{Name: "test-prompts", Version: "0.1.0"}, nil)
	registerPrompts(srv, scn, opt, false)
	cs := connectClientServer(t, ctx, srv)

	t.Run("get-skills prompt", func(t *testing.T) {
		res, err := cs.GetPrompt(ctx, &mcp.GetPromptParams{
			Name:      "get-skills",
			Arguments: map[string]string{},
		})
		if err != nil {
			t.Fatalf("GetPrompt(get-skills): %v", err)
		}
		if res.Description != "Skills for AI assistants" {
			t.Errorf("Description = %q", res.Description)
		}
		if len(res.Messages) == 0 {
			t.Fatal("no messages")
		}
		tc, ok := res.Messages[0].Content.(*mcp.TextContent)
		if !ok {
			t.Fatalf("expected TextContent, got %T", res.Messages[0].Content)
		}
		if !strings.Contains(tc.Text, "multi-stage builds") {
			t.Errorf("missing docker content: %s", tc.Text)
		}
		if !strings.Contains(tc.Text, "type hints") {
			t.Errorf("missing python content: %s", tc.Text)
		}
	})

	t.Run("get-skills prompt with source filter", func(t *testing.T) {
		// The source is the basename of the temp dir — find it
		skills := scn.List()
		if len(skills) == 0 {
			t.Fatal("no skills")
		}
		source := skills[0].Source
		res, err := cs.GetPrompt(ctx, &mcp.GetPromptParams{
			Name:      "get-skills",
			Arguments: map[string]string{"source": source},
		})
		if err != nil {
			t.Fatalf("GetPrompt(get-skills with source): %v", err)
		}
		if len(res.Messages) == 0 {
			t.Fatal("no messages")
		}
		tc, ok := res.Messages[0].Content.(*mcp.TextContent)
		if !ok {
			t.Fatalf("expected TextContent, got %T", res.Messages[0].Content)
		}
		if !strings.Contains(tc.Text, "multi-stage builds") && !strings.Contains(tc.Text, "type hints") {
			t.Errorf("expected skill content: %s", tc.Text)
		}
	})

	t.Run("get-skills prompt no match returns no skills found", func(t *testing.T) {
		res, err := cs.GetPrompt(ctx, &mcp.GetPromptParams{
			Name:      "get-skills",
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
		if !strings.Contains(tc.Text, "No skills found") {
			t.Errorf("expected 'No skills found', got: %s", tc.Text)
		}
	})

	t.Run("get-skill prompt found", func(t *testing.T) {
		res, err := cs.GetPrompt(ctx, &mcp.GetPromptParams{
			Name:      "get-skill",
			Arguments: map[string]string{"name": "Python Developer"},
		})
		if err != nil {
			t.Fatalf("GetPrompt(get-skill): %v", err)
		}
		if !strings.Contains(res.Description, "Python Developer") {
			t.Errorf("Description = %q", res.Description)
		}
		if len(res.Messages) == 0 {
			t.Fatal("no messages")
		}
		tc, ok := res.Messages[0].Content.(*mcp.TextContent)
		if !ok {
			t.Fatalf("expected TextContent, got %T", res.Messages[0].Content)
		}
		if !strings.Contains(tc.Text, "type hints") {
			t.Errorf("missing python content: %s", tc.Text)
		}
		// Should include reference content since renderSkill is used
		if !strings.Contains(tc.Text, "pep8.md") {
			t.Errorf("missing reference: %s", tc.Text)
		}
	})

	t.Run("get-skill prompt not found", func(t *testing.T) {
		res, err := cs.GetPrompt(ctx, &mcp.GetPromptParams{
			Name:      "get-skill",
			Arguments: map[string]string{"name": "Nonexistent"},
		})
		if err != nil {
			t.Fatalf("GetPrompt(get-skill): %v", err)
		}
		if len(res.Messages) == 0 {
			t.Fatal("no messages")
		}
		tc, ok := res.Messages[0].Content.(*mcp.TextContent)
		if !ok {
			t.Fatalf("expected TextContent, got %T", res.Messages[0].Content)
		}
		if !strings.Contains(tc.Text, "not found") {
			t.Errorf("expected not found message: %s", tc.Text)
		}
	})
}

// suppress unused import warning
var _ = fmt.Sprintf
