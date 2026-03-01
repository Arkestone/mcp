//go:build integration

package main

import (
	"context"
	"encoding/json"
	"strings"
	"testing"

	"github.com/modelcontextprotocol/go-sdk/mcp"
)

func TestIntegration_ReadResource_Individual(t *testing.T) {
	ctx := context.Background()
	ldr, opt, cleanup := setupTestEnv(t)
	defer cleanup()

	server := mcp.NewServer(&mcp.Implementation{Name: "mcp-instructions", Version: "test"}, nil)
	registerResources(server, ldr, opt, false)
	registerPrompts(server, ldr, opt, false)
	registerTools(server, ldr, opt, false)

	cs := connectClientServer(t, ctx, server)
	defer cs.Close()

	src := source(ldr)
	if src == "" {
		t.Fatal("no instructions found")
	}

	// Read copilot-instructions
	res, err := cs.ReadResource(ctx, &mcp.ReadResourceParams{
		URI: "instructions://" + src + "/copilot-instructions",
	})
	if err != nil {
		t.Fatalf("ReadResource copilot-instructions: %v", err)
	}
	if len(res.Contents) != 1 {
		t.Fatalf("expected 1 content, got %d", len(res.Contents))
	}
	if !strings.Contains(res.Contents[0].Text, "Go idioms") {
		t.Errorf("unexpected content: %q", res.Contents[0].Text)
	}

	// Read testing instructions
	res, err = cs.ReadResource(ctx, &mcp.ReadResourceParams{
		URI: "instructions://" + src + "/testing",
	})
	if err != nil {
		t.Fatalf("ReadResource testing: %v", err)
	}
	if !strings.Contains(res.Contents[0].Text, "table-driven") {
		t.Errorf("unexpected content: %q", res.Contents[0].Text)
	}
}

func TestIntegration_ReadResource_Index(t *testing.T) {
	ctx := context.Background()
	ldr, opt, cleanup := setupTestEnv(t)
	defer cleanup()

	server := mcp.NewServer(&mcp.Implementation{Name: "mcp-instructions", Version: "test"}, nil)
	registerResources(server, ldr, opt, false)
	registerPrompts(server, ldr, opt, false)
	registerTools(server, ldr, opt, false)

	cs := connectClientServer(t, ctx, server)
	defer cs.Close()

	res, err := cs.ReadResource(ctx, &mcp.ReadResourceParams{URI: "instructions://index"})
	if err != nil {
		t.Fatalf("ReadResource index: %v", err)
	}
	if len(res.Contents) != 1 {
		t.Fatalf("expected 1 content, got %d", len(res.Contents))
	}
	text := res.Contents[0].Text
	// Should list all three instruction files
	if !strings.Contains(text, "copilot-instructions") {
		t.Errorf("index missing copilot-instructions: %q", text)
	}
	if !strings.Contains(text, "testing") {
		t.Errorf("index missing testing: %q", text)
	}
	if !strings.Contains(text, "style") {
		t.Errorf("index missing style: %q", text)
	}
}

func TestIntegration_ReadResource_Optimized(t *testing.T) {
	ctx := context.Background()
	ldr, opt, cleanup := setupTestEnv(t)
	defer cleanup()

	server := mcp.NewServer(&mcp.Implementation{Name: "mcp-instructions", Version: "test"}, nil)
	registerResources(server, ldr, opt, false)
	registerPrompts(server, ldr, opt, false)
	registerTools(server, ldr, opt, false)

	cs := connectClientServer(t, ctx, server)
	defer cs.Close()

	res, err := cs.ReadResource(ctx, &mcp.ReadResourceParams{URI: "instructions://optimized"})
	if err != nil {
		t.Fatalf("ReadResource optimized: %v", err)
	}
	if len(res.Contents) != 1 {
		t.Fatalf("expected 1 content, got %d", len(res.Contents))
	}
	text := res.Contents[0].Text
	// Without LLM, should concatenate all instructions
	if !strings.Contains(text, "Go idioms") {
		t.Errorf("optimized missing 'Go idioms': %q", text)
	}
	if !strings.Contains(text, "table-driven") {
		t.Errorf("optimized missing 'table-driven': %q", text)
	}
	if !strings.Contains(text, "gofmt") {
		t.Errorf("optimized missing 'gofmt': %q", text)
	}
}

func TestIntegration_Prompt_GetInstructions(t *testing.T) {
	ctx := context.Background()
	ldr, opt, cleanup := setupTestEnv(t)
	defer cleanup()

	server := mcp.NewServer(&mcp.Implementation{Name: "mcp-instructions", Version: "test"}, nil)
	registerResources(server, ldr, opt, false)
	registerPrompts(server, ldr, opt, false)
	registerTools(server, ldr, opt, false)

	cs := connectClientServer(t, ctx, server)
	defer cs.Close()

	// List prompts first
	listRes, err := cs.ListPrompts(ctx, nil)
	if err != nil {
		t.Fatalf("ListPrompts: %v", err)
	}
	found := false
	for _, p := range listRes.Prompts {
		if p.Name == "get-instructions" {
			found = true
			break
		}
	}
	if !found {
		t.Fatal("get-instructions prompt not found")
	}

	// Get prompt without filters
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
	if !strings.Contains(tc.Text, "Go idioms") {
		t.Errorf("prompt missing 'Go idioms': %q", tc.Text)
	}
	if !strings.Contains(tc.Text, "table-driven") {
		t.Errorf("prompt missing 'table-driven': %q", tc.Text)
	}
}

func TestIntegration_Prompt_GetInstructions_SourceFilter(t *testing.T) {
	ctx := context.Background()
	ldr, opt, cleanup := setupTestEnv(t)
	defer cleanup()

	server := mcp.NewServer(&mcp.Implementation{Name: "mcp-instructions", Version: "test"}, nil)
	registerResources(server, ldr, opt, false)
	registerPrompts(server, ldr, opt, false)
	registerTools(server, ldr, opt, false)

	cs := connectClientServer(t, ctx, server)
	defer cs.Close()

	src := source(ldr)

	// With matching source filter — should return instructions
	res, err := cs.GetPrompt(ctx, &mcp.GetPromptParams{
		Name:      "get-instructions",
		Arguments: map[string]string{"source": src},
	})
	if err != nil {
		t.Fatalf("GetPrompt with source: %v", err)
	}
	tc, ok := res.Messages[0].Content.(*mcp.TextContent)
	if !ok {
		t.Fatal("expected TextContent")
	}
	if !strings.Contains(tc.Text, "Go idioms") {
		t.Errorf("expected instructions content, got: %q", tc.Text)
	}

	// With non-matching source filter — should return "No instructions found."
	res, err = cs.GetPrompt(ctx, &mcp.GetPromptParams{
		Name:      "get-instructions",
		Arguments: map[string]string{"source": "nonexistent-source"},
	})
	if err != nil {
		t.Fatalf("GetPrompt with bad source: %v", err)
	}
	tc, ok = res.Messages[0].Content.(*mcp.TextContent)
	if !ok {
		t.Fatal("expected TextContent")
	}
	if tc.Text != "No instructions found." {
		t.Errorf("expected 'No instructions found.', got: %q", tc.Text)
	}
}

func TestIntegration_Tool_Refresh(t *testing.T) {
	ctx := context.Background()
	ldr, opt, cleanup := setupTestEnv(t)
	defer cleanup()

	server := mcp.NewServer(&mcp.Implementation{Name: "mcp-instructions", Version: "test"}, nil)
	registerResources(server, ldr, opt, false)
	registerPrompts(server, ldr, opt, false)
	registerTools(server, ldr, opt, false)

	cs := connectClientServer(t, ctx, server)
	defer cs.Close()

	res, err := cs.CallTool(ctx, &mcp.CallToolParams{Name: "refresh-instructions"})
	if err != nil {
		t.Fatalf("CallTool refresh: %v", err)
	}
	if res.IsError {
		t.Fatal("refresh returned error")
	}

	// The structured content should contain message and count
	raw, err := json.Marshal(res.StructuredContent)
	if err != nil {
		t.Fatalf("marshal structured content: %v", err)
	}
	var result struct {
		Message string `json:"message"`
		Count   int    `json:"count"`
	}
	if err := json.Unmarshal(raw, &result); err != nil {
		t.Fatalf("unmarshal: %v", err)
	}
	if result.Message != "All instruction sources refreshed successfully" {
		t.Errorf("expected 'All instruction sources refreshed successfully', got %q", result.Message)
	}
	if result.Count != 3 {
		t.Errorf("expected count=3, got %d", result.Count)
	}
}

func TestIntegration_Tool_ListInstructions(t *testing.T) {
	ctx := context.Background()
	ldr, opt, cleanup := setupTestEnv(t)
	defer cleanup()

	server := mcp.NewServer(&mcp.Implementation{Name: "mcp-instructions", Version: "test"}, nil)
	registerResources(server, ldr, opt, false)
	registerPrompts(server, ldr, opt, false)
	registerTools(server, ldr, opt, false)

	cs := connectClientServer(t, ctx, server)
	defer cs.Close()

	res, err := cs.CallTool(ctx, &mcp.CallToolParams{Name: "list-instructions"})
	if err != nil {
		t.Fatalf("CallTool list-instructions: %v", err)
	}
	if res.IsError {
		t.Fatal("list-instructions returned error")
	}

	raw, err := json.Marshal(res.StructuredContent)
	if err != nil {
		t.Fatalf("marshal: %v", err)
	}
	var output struct {
		Entries []struct {
			URI    string `json:"uri"`
			Source string `json:"source"`
			Path   string `json:"path"`
		} `json:"entries"`
	}
	if err := json.Unmarshal(raw, &output); err != nil {
		t.Fatalf("unmarshal: %v", err)
	}
	if len(output.Entries) != 3 {
		t.Fatalf("expected 3 entries, got %d", len(output.Entries))
	}

	uris := make(map[string]bool)
	for _, e := range output.Entries {
		uris[e.URI] = true
	}
	src := source(ldr)
	for _, name := range []string{"copilot-instructions", "testing", "style"} {
		uri := "instructions://" + src + "/" + name
		if !uris[uri] {
			t.Errorf("missing URI: %s", uri)
		}
	}
}

func TestIntegration_Tool_OptimizeInstructions(t *testing.T) {
	ctx := context.Background()
	ldr, opt, cleanup := setupTestEnv(t)
	defer cleanup()

	server := mcp.NewServer(&mcp.Implementation{Name: "mcp-instructions", Version: "test"}, nil)
	registerResources(server, ldr, opt, false)
	registerPrompts(server, ldr, opt, false)
	registerTools(server, ldr, opt, false)

	cs := connectClientServer(t, ctx, server)
	defer cs.Close()

	res, err := cs.CallTool(ctx, &mcp.CallToolParams{
		Name:      "optimize-instructions",
		Arguments: map[string]any{},
	})
	if err != nil {
		t.Fatalf("CallTool optimize-instructions: %v", err)
	}
	if res.IsError {
		t.Fatal("optimize-instructions returned error")
	}

	raw, err := json.Marshal(res.StructuredContent)
	if err != nil {
		t.Fatalf("marshal: %v", err)
	}
	var result struct {
		Content string `json:"content"`
		Sources int    `json:"sources"`
	}
	if err := json.Unmarshal(raw, &result); err != nil {
		t.Fatalf("unmarshal: %v", err)
	}
	if result.Sources != 3 {
		t.Errorf("expected sources=3, got %d", result.Sources)
	}
	// Without LLM, content is concatenation
	if !strings.Contains(result.Content, "Go idioms") {
		t.Errorf("content missing 'Go idioms'")
	}
	if !strings.Contains(result.Content, "table-driven") {
		t.Errorf("content missing 'table-driven'")
	}
	if !strings.Contains(result.Content, "gofmt") {
		t.Errorf("content missing 'gofmt'")
	}
}
