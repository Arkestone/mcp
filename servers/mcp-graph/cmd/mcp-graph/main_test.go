package main

import (
	"context"
	"encoding/json"
	"os"
	"strings"
	"testing"

	"github.com/modelcontextprotocol/go-sdk/mcp"

	"github.com/Arkestone/mcp/servers/mcp-graph/internal/graph"
)

// ---------------------------------------------------------------------------
// helpers
// ---------------------------------------------------------------------------

func setupTestGraph(t *testing.T) *graph.Graph {
	t.Helper()
	g, err := graph.New(t.TempDir())
	if err != nil {
		t.Fatalf("graph.New: %v", err)
	}
	return g
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

func TestJsonStringSlice(t *testing.T) {
	tests := []struct {
		name  string
		input []string
		want  string
	}{
		{"nil", nil, "[]"},
		{"empty", []string{}, "[]"},
		{"single", []string{"a"}, `["a"]`},
		{"multiple", []string{"a", "b", "c"}, `["a","b","c"]`},
		{"with quotes", []string{`a"b`}, `["a\"b"]`},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			got := jsonStringSlice(tc.input)
			if got != tc.want {
				t.Errorf("jsonStringSlice(%v) = %q, want %q", tc.input, got, tc.want)
			}
		})
	}
}

// ---------------------------------------------------------------------------
// Resource tests
// ---------------------------------------------------------------------------

func TestRegisterResources_Stats(t *testing.T) {
	g := setupTestGraph(t)
	ctx := context.Background()
	srv := mcp.NewServer(&mcp.Implementation{Name: "test", Version: "0.1.0"}, nil)
	registerResources(srv, g)
	cs := connectClientServer(t, ctx, srv)

	t.Run("empty graph", func(t *testing.T) {
		res, err := cs.ReadResource(ctx, &mcp.ReadResourceParams{URI: "graph://stats"})
		if err != nil {
			t.Fatalf("ReadResource(graph://stats): %v", err)
		}
		if len(res.Contents) == 0 {
			t.Fatal("no contents")
		}
		var stats struct {
			Nodes int `json:"nodes"`
			Edges int `json:"edges"`
		}
		if err := json.Unmarshal([]byte(res.Contents[0].Text), &stats); err != nil {
			t.Fatalf("unmarshal stats: %v", err)
		}
		if stats.Nodes != 0 || stats.Edges != 0 {
			t.Errorf("empty graph: nodes=%d edges=%d, want 0 0", stats.Nodes, stats.Edges)
		}
	})

	t.Run("after adding nodes and edges", func(t *testing.T) {
		a, _ := g.AddNode("Person", "Alice", nil)
		b, _ := g.AddNode("Person", "Bob", nil)
		g.AddEdge(a.ID, b.ID, "knows", nil)

		res, err := cs.ReadResource(ctx, &mcp.ReadResourceParams{URI: "graph://stats"})
		if err != nil {
			t.Fatalf("ReadResource(graph://stats): %v", err)
		}
		var stats struct {
			Nodes     int      `json:"nodes"`
			Edges     int      `json:"edges"`
			Relations []string `json:"relations"`
		}
		if err := json.Unmarshal([]byte(res.Contents[0].Text), &stats); err != nil {
			t.Fatalf("unmarshal stats: %v", err)
		}
		if stats.Nodes != 2 {
			t.Errorf("nodes = %d, want 2", stats.Nodes)
		}
		if stats.Edges != 1 {
			t.Errorf("edges = %d, want 1", stats.Edges)
		}
		if len(stats.Relations) != 1 || stats.Relations[0] != "knows" {
			t.Errorf("relations = %v, want [knows]", stats.Relations)
		}
	})
}

func TestRegisterResources_NodeTemplate(t *testing.T) {
	g := setupTestGraph(t)
	ctx := context.Background()
	srv := mcp.NewServer(&mcp.Implementation{Name: "test", Version: "0.1.0"}, nil)
	registerResources(srv, g)
	cs := connectClientServer(t, ctx, srv)

	n, _ := g.AddNode("Technology", "Go", map[string]string{"version": "1.24"})

	t.Run("existing node", func(t *testing.T) {
		res, err := cs.ReadResource(ctx, &mcp.ReadResourceParams{URI: "graph://node/" + n.ID})
		if err != nil {
			t.Fatalf("ReadResource(graph://node/%s): %v", n.ID, err)
		}
		if !strings.Contains(res.Contents[0].Text, "Go") {
			t.Errorf("missing node name: %s", res.Contents[0].Text)
		}
		if !strings.Contains(res.Contents[0].Text, "Technology") {
			t.Errorf("missing node label: %s", res.Contents[0].Text)
		}
	})

	t.Run("nonexistent node returns error", func(t *testing.T) {
		_, err := cs.ReadResource(ctx, &mcp.ReadResourceParams{URI: "graph://node/nonexistent"})
		if err == nil {
			t.Error("expected error for nonexistent node")
		}
	})
}

// ---------------------------------------------------------------------------
// Tool tests
// ---------------------------------------------------------------------------

func TestRegisterTools_AddNode(t *testing.T) {
	g := setupTestGraph(t)
	ctx := context.Background()
	srv := mcp.NewServer(&mcp.Implementation{Name: "test", Version: "0.1.0"}, nil)
	registerTools(srv, g)
	cs := connectClientServer(t, ctx, srv)

	res, err := cs.CallTool(ctx, &mcp.CallToolParams{
		Name:      "add-node",
		Arguments: map[string]any{"label": "Technology", "name": "Go"},
	})
	if err != nil {
		t.Fatalf("CallTool(add-node): %v", err)
	}
	text := textContent(t, res)
	if !strings.Contains(text, "Go") || !strings.Contains(text, "Technology") {
		t.Errorf("unexpected output: %s", text)
	}

	nodes, _ := g.Stats()
	if nodes != 1 {
		t.Errorf("graph has %d nodes after add-node, want 1", nodes)
	}
}

func TestRegisterTools_AddEdge(t *testing.T) {
	g := setupTestGraph(t)
	ctx := context.Background()
	srv := mcp.NewServer(&mcp.Implementation{Name: "test", Version: "0.1.0"}, nil)
	registerTools(srv, g)
	cs := connectClientServer(t, ctx, srv)

	alice, _ := g.AddNode("Person", "Alice", nil)
	go_, _ := g.AddNode("Technology", "Go", nil)

	t.Run("valid edge", func(t *testing.T) {
		res, err := cs.CallTool(ctx, &mcp.CallToolParams{
			Name:      "add-edge",
			Arguments: map[string]any{"from": alice.ID, "to": go_.ID, "relation": "uses"},
		})
		if err != nil {
			t.Fatalf("CallTool(add-edge): %v", err)
		}
		text := textContent(t, res)
		if !strings.Contains(text, "uses") {
			t.Errorf("relation missing from output: %s", text)
		}
	})

	t.Run("nonexistent source node", func(t *testing.T) {
		res, err := cs.CallTool(ctx, &mcp.CallToolParams{
			Name:      "add-edge",
			Arguments: map[string]any{"from": "nonexistent", "to": go_.ID, "relation": "uses"},
		})
		if err != nil {
			return // protocol error is acceptable
		}
		if res != nil && !res.IsError {
			t.Error("expected error for nonexistent source node")
		}
	})
}

func TestRegisterTools_FindNodes(t *testing.T) {
	g := setupTestGraph(t)
	ctx := context.Background()
	srv := mcp.NewServer(&mcp.Implementation{Name: "test", Version: "0.1.0"}, nil)
	registerTools(srv, g)
	cs := connectClientServer(t, ctx, srv)

	g.AddNode("Person", "Alice", nil)
	g.AddNode("Person", "Bob", nil)
	g.AddNode("Technology", "Go", nil)

	t.Run("filter by label", func(t *testing.T) {
		res, err := cs.CallTool(ctx, &mcp.CallToolParams{
			Name:      "find-nodes",
			Arguments: map[string]any{"label": "Person"},
		})
		if err != nil {
			t.Fatalf("CallTool(find-nodes label=Person): %v", err)
		}
		text := textContent(t, res)
		if !strings.Contains(text, "Alice") || !strings.Contains(text, "Bob") {
			t.Errorf("missing persons: %s", text)
		}
		if strings.Contains(text, "\"count\":3") {
			t.Errorf("Technology node should be filtered out: %s", text)
		}
	})

	t.Run("filter by name substring", func(t *testing.T) {
		res, err := cs.CallTool(ctx, &mcp.CallToolParams{
			Name:      "find-nodes",
			Arguments: map[string]any{"name": "alice"},
		})
		if err != nil {
			t.Fatalf("CallTool(find-nodes name=alice): %v", err)
		}
		text := textContent(t, res)
		if !strings.Contains(text, "Alice") {
			t.Errorf("missing Alice (case-insensitive): %s", text)
		}
	})

	t.Run("empty filters return all", func(t *testing.T) {
		res, err := cs.CallTool(ctx, &mcp.CallToolParams{
			Name:      "find-nodes",
			Arguments: map[string]any{},
		})
		if err != nil {
			t.Fatalf("CallTool(find-nodes empty): %v", err)
		}
		text := textContent(t, res)
		if !strings.Contains(text, `"count":3`) {
			t.Errorf("expected count 3 for all nodes: %s", text)
		}
	})
}

func TestRegisterTools_GetNode(t *testing.T) {
	g := setupTestGraph(t)
	ctx := context.Background()
	srv := mcp.NewServer(&mcp.Implementation{Name: "test", Version: "0.1.0"}, nil)
	registerTools(srv, g)
	cs := connectClientServer(t, ctx, srv)

	n, _ := g.AddNode("Technology", "Go", map[string]string{"version": "1.24"})

	t.Run("existing node", func(t *testing.T) {
		res, err := cs.CallTool(ctx, &mcp.CallToolParams{
			Name:      "get-node",
			Arguments: map[string]any{"id": n.ID},
		})
		if err != nil {
			t.Fatalf("CallTool(get-node): %v", err)
		}
		text := textContent(t, res)
		if !strings.Contains(text, "Go") || !strings.Contains(text, "Technology") {
			t.Errorf("missing node fields: %s", text)
		}
	})

	t.Run("nonexistent node returns error", func(t *testing.T) {
		res, err := cs.CallTool(ctx, &mcp.CallToolParams{
			Name:      "get-node",
			Arguments: map[string]any{"id": "nonexistent"},
		})
		if err != nil {
			return
		}
		if res != nil && !res.IsError {
			t.Error("expected error for nonexistent node")
		}
	})
}

func TestRegisterTools_Neighbors(t *testing.T) {
	g := setupTestGraph(t)
	ctx := context.Background()
	srv := mcp.NewServer(&mcp.Implementation{Name: "test", Version: "0.1.0"}, nil)
	registerTools(srv, g)
	cs := connectClientServer(t, ctx, srv)

	alice, _ := g.AddNode("Person", "Alice", nil)
	bob, _ := g.AddNode("Person", "Bob", nil)
	go_, _ := g.AddNode("Technology", "Go", nil)
	g.AddEdge(alice.ID, bob.ID, "knows", nil)
	g.AddEdge(alice.ID, go_.ID, "uses", nil)

	t.Run("all outbound neighbors", func(t *testing.T) {
		res, err := cs.CallTool(ctx, &mcp.CallToolParams{
			Name:      "neighbors",
			Arguments: map[string]any{"id": alice.ID, "direction": "out"},
		})
		if err != nil {
			t.Fatalf("CallTool(neighbors direction=out): %v", err)
		}
		text := textContent(t, res)
		if !strings.Contains(text, "Bob") || !strings.Contains(text, "Go") {
			t.Errorf("missing neighbors: %s", text)
		}
	})

	t.Run("filter by relation", func(t *testing.T) {
		res, err := cs.CallTool(ctx, &mcp.CallToolParams{
			Name:      "neighbors",
			Arguments: map[string]any{"id": alice.ID, "relation": "uses"},
		})
		if err != nil {
			t.Fatalf("CallTool(neighbors relation=uses): %v", err)
		}
		text := textContent(t, res)
		if !strings.Contains(text, "Go") {
			t.Errorf("missing Go: %s", text)
		}
		if strings.Contains(text, "Bob") {
			t.Errorf("unexpected Bob in uses relation: %s", text)
		}
	})

	t.Run("inbound neighbors", func(t *testing.T) {
		res, err := cs.CallTool(ctx, &mcp.CallToolParams{
			Name:      "neighbors",
			Arguments: map[string]any{"id": bob.ID, "direction": "in"},
		})
		if err != nil {
			t.Fatalf("CallTool(neighbors direction=in): %v", err)
		}
		text := textContent(t, res)
		if !strings.Contains(text, "Alice") {
			t.Errorf("missing inbound neighbor Alice: %s", text)
		}
	})
}

func TestRegisterTools_ShortestPath(t *testing.T) {
	g := setupTestGraph(t)
	ctx := context.Background()
	srv := mcp.NewServer(&mcp.Implementation{Name: "test", Version: "0.1.0"}, nil)
	registerTools(srv, g)
	cs := connectClientServer(t, ctx, srv)

	a, _ := g.AddNode("Node", "A", nil)
	b, _ := g.AddNode("Node", "B", nil)
	c, _ := g.AddNode("Node", "C", nil)
	g.AddEdge(a.ID, b.ID, "connects", nil)
	g.AddEdge(b.ID, c.ID, "connects", nil)

	t.Run("path exists", func(t *testing.T) {
		res, err := cs.CallTool(ctx, &mcp.CallToolParams{
			Name:      "shortest-path",
			Arguments: map[string]any{"from": a.ID, "to": c.ID},
		})
		if err != nil {
			t.Fatalf("CallTool(shortest-path): %v", err)
		}
		text := textContent(t, res)
		if !strings.Contains(text, `"hops":2`) {
			t.Errorf("expected 2 hops: %s", text)
		}
	})

	t.Run("no path between nodes", func(t *testing.T) {
		d, _ := g.AddNode("Node", "D", nil)
		res, err := cs.CallTool(ctx, &mcp.CallToolParams{
			Name:      "shortest-path",
			Arguments: map[string]any{"from": a.ID, "to": d.ID},
		})
		if err != nil {
			return
		}
		if res != nil && !res.IsError {
			t.Error("expected error for unreachable node")
		}
	})
}

func TestRegisterTools_RemoveEdge(t *testing.T) {
	g := setupTestGraph(t)
	ctx := context.Background()
	srv := mcp.NewServer(&mcp.Implementation{Name: "test", Version: "0.1.0"}, nil)
	registerTools(srv, g)
	cs := connectClientServer(t, ctx, srv)

	a, _ := g.AddNode("Node", "A", nil)
	b, _ := g.AddNode("Node", "B", nil)
	e, _ := g.AddEdge(a.ID, b.ID, "connects", nil)

	t.Run("remove existing edge", func(t *testing.T) {
		res, err := cs.CallTool(ctx, &mcp.CallToolParams{
			Name:      "remove-edge",
			Arguments: map[string]any{"id": e.ID},
		})
		if err != nil {
			t.Fatalf("CallTool(remove-edge): %v", err)
		}
		if !strings.Contains(textContent(t, res), "deleted") {
			t.Error("expected 'deleted' in response")
		}
	})

	t.Run("remove nonexistent edge", func(t *testing.T) {
		res, err := cs.CallTool(ctx, &mcp.CallToolParams{
			Name:      "remove-edge",
			Arguments: map[string]any{"id": "nonexistent"},
		})
		if err != nil {
			return
		}
		if res != nil && !res.IsError {
			t.Error("expected error for nonexistent edge")
		}
	})
}

func TestRegisterTools_RemoveNode(t *testing.T) {
	g := setupTestGraph(t)
	ctx := context.Background()
	srv := mcp.NewServer(&mcp.Implementation{Name: "test", Version: "0.1.0"}, nil)
	registerTools(srv, g)
	cs := connectClientServer(t, ctx, srv)

	a, _ := g.AddNode("Node", "A", nil)
	b, _ := g.AddNode("Node", "B", nil)
	g.AddEdge(a.ID, b.ID, "connects", nil)

	res, err := cs.CallTool(ctx, &mcp.CallToolParams{
		Name:      "remove-node",
		Arguments: map[string]any{"id": a.ID},
	})
	if err != nil {
		t.Fatalf("CallTool(remove-node): %v", err)
	}
	text := textContent(t, res)
	if !strings.Contains(text, "deleted") {
		t.Errorf("expected 'deleted': %s", text)
	}
	if !strings.Contains(text, `"deleted_edges":1`) {
		t.Errorf("expected 1 deleted edge: %s", text)
	}
}

func TestRegisterTools_ListRelations(t *testing.T) {
	g := setupTestGraph(t)
	ctx := context.Background()
	srv := mcp.NewServer(&mcp.Implementation{Name: "test", Version: "0.1.0"}, nil)
	registerTools(srv, g)
	cs := connectClientServer(t, ctx, srv)

	t.Run("empty graph", func(t *testing.T) {
		res, err := cs.CallTool(ctx, &mcp.CallToolParams{Name: "list-relations"})
		if err != nil {
			t.Fatalf("CallTool(list-relations): %v", err)
		}
		if !strings.Contains(textContent(t, res), `"count":0`) {
			t.Errorf("expected count 0 in empty graph")
		}
	})

	t.Run("with relations", func(t *testing.T) {
		a, _ := g.AddNode("A", "a", nil)
		b, _ := g.AddNode("B", "b", nil)
		g.AddEdge(a.ID, b.ID, "knows", nil)
		g.AddEdge(a.ID, b.ID, "uses", nil)

		res, err := cs.CallTool(ctx, &mcp.CallToolParams{Name: "list-relations"})
		if err != nil {
			t.Fatalf("CallTool(list-relations): %v", err)
		}
		text := textContent(t, res)
		if !strings.Contains(text, "knows") || !strings.Contains(text, "uses") {
			t.Errorf("missing relations: %s", text)
		}
	})
}
