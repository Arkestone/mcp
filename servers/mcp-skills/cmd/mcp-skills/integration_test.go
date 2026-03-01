//go:build integration

package main

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/Arkestone/mcp/pkg/config"
	"github.com/Arkestone/mcp/pkg/github"
	"github.com/Arkestone/mcp/servers/mcp-skills/internal/scanner"
	"github.com/modelcontextprotocol/go-sdk/mcp"
)

func setupTestSkills(t *testing.T) string {
	t.Helper()
	dir := t.TempDir()
	for _, sk := range []struct{ dir, name, desc, body string }{
		{"docker-expert", "Docker Expert", "Docker skills", "Use multi-stage builds"},
		{"python-dev", "Python Developer", "Python skills", "Use type hints"},
	} {
		if err := os.MkdirAll(filepath.Join(dir, sk.dir), 0o755); err != nil {
			t.Fatal(err)
		}
		content := fmt.Sprintf("---\nname: %s\ndescription: %s\n---\n%s", sk.name, sk.desc, sk.body)
		if err := os.WriteFile(filepath.Join(dir, sk.dir, "SKILL.md"), []byte(content), 0o644); err != nil {
			t.Fatal(err)
		}
	}
	return dir
}

func setupIntegration(t *testing.T) *mcp.ClientSession {
	t.Helper()

	dir := setupTestSkills(t)

	cfg := &config.Config{
		Sources: config.Sources{
			Dirs: []string{dir},
		},
		Cache: config.CacheConfig{
			Dir:          t.TempDir(),
			SyncInterval: 5 * time.Minute,
		},
	}
	scn := scanner.New(cfg, &github.Client{})

	server := mcp.NewServer(
		&mcp.Implementation{Name: "mcp-skills-test", Version: "0.1.0"},
		nil,
	)
	registerResources(server, scn, nil, false)
	registerPrompts(server, scn, nil, false)
	registerTools(server, scn, nil, false)

	serverTransport, clientTransport := mcp.NewInMemoryTransports()

	ctx := context.Background()
	go server.Run(ctx, serverTransport)

	client := mcp.NewClient(&mcp.Implementation{Name: "test", Version: "0.1.0"}, nil)
	session, err := client.Connect(ctx, clientTransport, nil)
	if err != nil {
		t.Fatalf("client.Connect: %v", err)
	}
	t.Cleanup(func() { session.Close() })

	return session
}

func TestIntegration_ListResources(t *testing.T) {
	session := setupIntegration(t)
	ctx := context.Background()

	res, err := session.ListResources(ctx, &mcp.ListResourcesParams{})
	if err != nil {
		t.Fatalf("ListResources: %v", err)
	}
	// We expect at least the two static resources (skills://optimized, skills://index)
	if len(res.Resources) < 2 {
		t.Fatalf("got %d resources, want at least 2", len(res.Resources))
	}

	var uris []string
	for _, r := range res.Resources {
		uris = append(uris, r.URI)
	}
	for _, want := range []string{"skills://optimized", "skills://index"} {
		found := false
		for _, u := range uris {
			if u == want {
				found = true
				break
			}
		}
		if !found {
			t.Errorf("resource %q not found in %v", want, uris)
		}
	}
}

func TestIntegration_ReadResource(t *testing.T) {
	session := setupIntegration(t)
	ctx := context.Background()

	// Read the index to find a skill URI
	indexRes, err := session.ReadResource(ctx, &mcp.ReadResourceParams{URI: "skills://index"})
	if err != nil {
		t.Fatalf("ReadResource(index): %v", err)
	}
	if len(indexRes.Contents) == 0 {
		t.Fatal("index returned no contents")
	}
	indexText := indexRes.Contents[0].Text
	if !strings.Contains(indexText, "Docker Expert") {
		t.Errorf("index missing Docker Expert: %s", indexText)
	}
}

func TestIntegration_ReadIndex(t *testing.T) {
	session := setupIntegration(t)
	ctx := context.Background()

	res, err := session.ReadResource(ctx, &mcp.ReadResourceParams{URI: "skills://index"})
	if err != nil {
		t.Fatalf("ReadResource(index): %v", err)
	}
	if len(res.Contents) == 0 {
		t.Fatal("no contents returned")
	}
	text := res.Contents[0].Text
	if !strings.Contains(text, "Docker Expert") {
		t.Errorf("index missing Docker Expert")
	}
	if !strings.Contains(text, "Python Developer") {
		t.Errorf("index missing Python Developer")
	}
}

func TestIntegration_GetPrompt(t *testing.T) {
	session := setupIntegration(t)
	ctx := context.Background()

	res, err := session.GetPrompt(ctx, &mcp.GetPromptParams{
		Name:      "get-skills",
		Arguments: map[string]string{},
	})
	if err != nil {
		t.Fatalf("GetPrompt(get-skills): %v", err)
	}
	if len(res.Messages) == 0 {
		t.Fatal("no messages returned")
	}
	tc, ok := res.Messages[0].Content.(*mcp.TextContent)
	if !ok {
		t.Fatalf("expected TextContent, got %T", res.Messages[0].Content)
	}
	if !strings.Contains(tc.Text, "multi-stage builds") {
		t.Errorf("prompt missing docker content: %s", tc.Text)
	}
}

func TestIntegration_GetSingleSkillPrompt(t *testing.T) {
	session := setupIntegration(t)
	ctx := context.Background()

	res, err := session.GetPrompt(ctx, &mcp.GetPromptParams{
		Name:      "get-skill",
		Arguments: map[string]string{"name": "Docker Expert"},
	})
	if err != nil {
		t.Fatalf("GetPrompt(get-skill): %v", err)
	}
	if len(res.Messages) == 0 {
		t.Fatal("no messages returned")
	}
	tc, ok := res.Messages[0].Content.(*mcp.TextContent)
	if !ok {
		t.Fatalf("expected TextContent, got %T", res.Messages[0].Content)
	}
	if !strings.Contains(tc.Text, "multi-stage builds") {
		t.Errorf("prompt missing docker content: %s", tc.Text)
	}
}

func TestIntegration_CallRefreshTool(t *testing.T) {
	session := setupIntegration(t)
	ctx := context.Background()

	res, err := session.CallTool(ctx, &mcp.CallToolParams{
		Name: "refresh-skills",
	})
	if err != nil {
		t.Fatalf("CallTool(refresh-skills): %v", err)
	}
	if len(res.Content) == 0 {
		t.Fatal("no content returned")
	}
	tc, ok := res.Content[0].(*mcp.TextContent)
	if !ok {
		t.Fatalf("expected TextContent, got %T", res.Content[0])
	}
	if !strings.Contains(tc.Text, "refreshed") {
		t.Errorf("unexpected result: %s", tc.Text)
	}
}

func TestIntegration_CallListTool(t *testing.T) {
	session := setupIntegration(t)
	ctx := context.Background()

	res, err := session.CallTool(ctx, &mcp.CallToolParams{
		Name: "list-skills",
	})
	if err != nil {
		t.Fatalf("CallTool(list-skills): %v", err)
	}
	if len(res.Content) == 0 {
		t.Fatal("no content returned")
	}
	tc, ok := res.Content[0].(*mcp.TextContent)
	if !ok {
		t.Fatalf("expected TextContent, got %T", res.Content[0])
	}
	if !strings.Contains(tc.Text, "Docker Expert") {
		t.Errorf("list missing Docker Expert: %s", tc.Text)
	}
	if !strings.Contains(tc.Text, "Python Developer") {
		t.Errorf("list missing Python Developer: %s", tc.Text)
	}
}

func TestIntegration_CallGetSkillTool(t *testing.T) {
	session := setupIntegration(t)
	ctx := context.Background()

	res, err := session.CallTool(ctx, &mcp.CallToolParams{
		Name:      "get-skill",
		Arguments: map[string]any{"name": "Docker Expert"},
	})
	if err != nil {
		t.Fatalf("CallTool(get-skill): %v", err)
	}
	if len(res.Content) == 0 {
		t.Fatal("no content returned")
	}
	tc, ok := res.Content[0].(*mcp.TextContent)
	if !ok {
		t.Fatalf("expected TextContent, got %T", res.Content[0])
	}
	if !strings.Contains(tc.Text, "multi-stage builds") {
		t.Errorf("get-skill missing docker content: %s", tc.Text)
	}
}
