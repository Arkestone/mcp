package graph

import (
	"os"
	"path/filepath"
	"testing"
)

func tmpGraph(t *testing.T) *Graph {
	t.Helper()
	dir := t.TempDir()
	g, err := New(dir)
	if err != nil {
		t.Fatalf("New: %v", err)
	}
	return g
}

// --- Node operations ---

func TestAddNode_basic(t *testing.T) {
	g := tmpGraph(t)
	n, err := g.AddNode("Person", "Alice", nil)
	if err != nil {
		t.Fatalf("AddNode: %v", err)
	}
	if n.ID == "" {
		t.Error("expected non-empty ID")
	}
	if n.Label != "Person" {
		t.Errorf("label: got %q, want %q", n.Label, "Person")
	}
	if n.Name != "Alice" {
		t.Errorf("name: got %q, want %q", n.Name, "Alice")
	}
}

func TestAddNode_withProps(t *testing.T) {
	g := tmpGraph(t)
	props := map[string]string{"role": "engineer", "city": "Paris"}
	n, err := g.AddNode("Person", "Bob", props)
	if err != nil {
		t.Fatalf("AddNode: %v", err)
	}
	if n.Props["role"] != "engineer" {
		t.Errorf("prop role: got %q", n.Props["role"])
	}
}

func TestAddNode_emptyLabel(t *testing.T) {
	g := tmpGraph(t)
	_, err := g.AddNode("", "Alice", nil)
	if err == nil {
		t.Error("expected error for empty label")
	}
}

func TestAddNode_emptyName(t *testing.T) {
	g := tmpGraph(t)
	_, err := g.AddNode("Person", "", nil)
	if err == nil {
		t.Error("expected error for empty name")
	}
}

// --- GetNode ---

func TestGetNode(t *testing.T) {
	g := tmpGraph(t)
	n, _ := g.AddNode("Tech", "Go", nil)
	got, ok := g.GetNode(n.ID)
	if !ok {
		t.Fatal("GetNode: not found")
	}
	if got.Name != "Go" {
		t.Errorf("name: got %q", got.Name)
	}
}

func TestGetNode_notFound(t *testing.T) {
	g := tmpGraph(t)
	_, ok := g.GetNode("nonexistent")
	if ok {
		t.Error("expected not found")
	}
}

// --- Edge operations ---

func TestAddEdge_basic(t *testing.T) {
	g := tmpGraph(t)
	a, _ := g.AddNode("Person", "Alice", nil)
	b, _ := g.AddNode("Person", "Bob", nil)
	e, err := g.AddEdge(a.ID, b.ID, "knows", nil)
	if err != nil {
		t.Fatalf("AddEdge: %v", err)
	}
	if e.From != a.ID || e.To != b.ID {
		t.Errorf("edge endpoints: from=%s to=%s", e.From, e.To)
	}
	if e.Relation != "knows" {
		t.Errorf("relation: got %q", e.Relation)
	}
}

func TestAddEdge_emptyRelation(t *testing.T) {
	g := tmpGraph(t)
	a, _ := g.AddNode("X", "A", nil)
	b, _ := g.AddNode("X", "B", nil)
	_, err := g.AddEdge(a.ID, b.ID, "", nil)
	if err == nil {
		t.Error("expected error for empty relation")
	}
}

func TestAddEdge_unknownFrom(t *testing.T) {
	g := tmpGraph(t)
	b, _ := g.AddNode("X", "B", nil)
	_, err := g.AddEdge("nope", b.ID, "rel", nil)
	if err == nil {
		t.Error("expected error for unknown from-node")
	}
}

func TestAddEdge_unknownTo(t *testing.T) {
	g := tmpGraph(t)
	a, _ := g.AddNode("X", "A", nil)
	_, err := g.AddEdge(a.ID, "nope", "rel", nil)
	if err == nil {
		t.Error("expected error for unknown to-node")
	}
}

// --- FindNodes ---

func TestFindNodes_byLabel(t *testing.T) {
	g := tmpGraph(t)
	g.AddNode("Person", "Alice", nil)
	g.AddNode("Person", "Bob", nil)
	g.AddNode("Tech", "Go", nil)
	persons := g.FindNodes("person", "")
	if len(persons) != 2 {
		t.Errorf("FindNodes by label: got %d, want 2", len(persons))
	}
}

func TestFindNodes_byName(t *testing.T) {
	g := tmpGraph(t)
	g.AddNode("Person", "Alice", nil)
	g.AddNode("Person", "Bob", nil)
	got := g.FindNodes("", "ali")
	if len(got) != 1 || got[0].Name != "Alice" {
		t.Errorf("FindNodes by name: got %+v", got)
	}
}

func TestFindNodes_all(t *testing.T) {
	g := tmpGraph(t)
	g.AddNode("X", "A", nil)
	g.AddNode("X", "B", nil)
	all := g.FindNodes("", "")
	if len(all) != 2 {
		t.Errorf("FindNodes all: got %d, want 2", len(all))
	}
}

func TestFindNodes_none(t *testing.T) {
	g := tmpGraph(t)
	g.AddNode("X", "A", nil)
	got := g.FindNodes("nonexistent", "")
	if len(got) != 0 {
		t.Errorf("expected empty, got %d", len(got))
	}
}

// --- Neighbors ---

func TestNeighbors_out(t *testing.T) {
	g := tmpGraph(t)
	a, _ := g.AddNode("X", "A", nil)
	b, _ := g.AddNode("X", "B", nil)
	g.AddEdge(a.ID, b.ID, "links", nil)

	res, err := g.Neighbors(a.ID, "out", "")
	if err != nil {
		t.Fatalf("Neighbors: %v", err)
	}
	if len(res.Nodes) != 1 || res.Nodes[0].ID != b.ID {
		t.Errorf("out neighbors: %+v", res.Nodes)
	}
}

func TestNeighbors_in(t *testing.T) {
	g := tmpGraph(t)
	a, _ := g.AddNode("X", "A", nil)
	b, _ := g.AddNode("X", "B", nil)
	g.AddEdge(a.ID, b.ID, "links", nil)

	res, err := g.Neighbors(b.ID, "in", "")
	if err != nil {
		t.Fatalf("Neighbors: %v", err)
	}
	if len(res.Nodes) != 1 || res.Nodes[0].ID != a.ID {
		t.Errorf("in neighbors: %+v", res.Nodes)
	}
}

func TestNeighbors_both(t *testing.T) {
	g := tmpGraph(t)
	a, _ := g.AddNode("X", "A", nil)
	b, _ := g.AddNode("X", "B", nil)
	c, _ := g.AddNode("X", "C", nil)
	g.AddEdge(a.ID, b.ID, "rel", nil)
	g.AddEdge(c.ID, a.ID, "rel", nil)

	res, err := g.Neighbors(a.ID, "both", "")
	if err != nil {
		t.Fatalf("Neighbors: %v", err)
	}
	if len(res.Nodes) != 2 {
		t.Errorf("both neighbors: got %d nodes, want 2", len(res.Nodes))
	}
}

func TestNeighbors_filterRelation(t *testing.T) {
	g := tmpGraph(t)
	a, _ := g.AddNode("X", "A", nil)
	b, _ := g.AddNode("X", "B", nil)
	c, _ := g.AddNode("X", "C", nil)
	g.AddEdge(a.ID, b.ID, "knows", nil)
	g.AddEdge(a.ID, c.ID, "uses", nil)

	res, err := g.Neighbors(a.ID, "out", "knows")
	if err != nil {
		t.Fatalf("Neighbors: %v", err)
	}
	if len(res.Nodes) != 1 || res.Nodes[0].ID != b.ID {
		t.Errorf("filtered neighbors: %+v", res.Nodes)
	}
}

func TestNeighbors_notFound(t *testing.T) {
	g := tmpGraph(t)
	_, err := g.Neighbors("nonexistent", "both", "")
	if err == nil {
		t.Error("expected error for unknown node")
	}
}

// --- ShortestPath ---

func TestShortestPath_direct(t *testing.T) {
	g := tmpGraph(t)
	a, _ := g.AddNode("X", "A", nil)
	b, _ := g.AddNode("X", "B", nil)
	g.AddEdge(a.ID, b.ID, "rel", nil)

	res, err := g.ShortestPath(a.ID, b.ID, 0)
	if err != nil {
		t.Fatalf("ShortestPath: %v", err)
	}
	if len(res.Nodes) != 2 {
		t.Errorf("nodes: got %d, want 2", len(res.Nodes))
	}
	if len(res.Edges) != 1 {
		t.Errorf("edges: got %d, want 1", len(res.Edges))
	}
}

func TestShortestPath_multiHop(t *testing.T) {
	g := tmpGraph(t)
	a, _ := g.AddNode("X", "A", nil)
	b, _ := g.AddNode("X", "B", nil)
	c, _ := g.AddNode("X", "C", nil)
	g.AddEdge(a.ID, b.ID, "rel", nil)
	g.AddEdge(b.ID, c.ID, "rel", nil)

	res, err := g.ShortestPath(a.ID, c.ID, 0)
	if err != nil {
		t.Fatalf("ShortestPath: %v", err)
	}
	if len(res.Nodes) != 3 {
		t.Errorf("nodes: got %d, want 3", len(res.Nodes))
	}
}

func TestShortestPath_noPath(t *testing.T) {
	g := tmpGraph(t)
	a, _ := g.AddNode("X", "A", nil)
	b, _ := g.AddNode("X", "B", nil)

	_, err := g.ShortestPath(a.ID, b.ID, 0)
	if err == nil {
		t.Error("expected error for unreachable target")
	}
}

func TestShortestPath_sameNode(t *testing.T) {
	g := tmpGraph(t)
	a, _ := g.AddNode("X", "A", nil)
	res, err := g.ShortestPath(a.ID, a.ID, 0)
	if err != nil {
		t.Fatalf("ShortestPath same node: %v", err)
	}
	if len(res.Nodes) != 1 || res.Nodes[0].ID != a.ID {
		t.Errorf("same node path: %+v", res.Nodes)
	}
}

func TestShortestPath_maxDepth(t *testing.T) {
	g := tmpGraph(t)
	a, _ := g.AddNode("X", "A", nil)
	b, _ := g.AddNode("X", "B", nil)
	c, _ := g.AddNode("X", "C", nil)
	g.AddEdge(a.ID, b.ID, "rel", nil)
	g.AddEdge(b.ID, c.ID, "rel", nil)

	_, err := g.ShortestPath(a.ID, c.ID, 1) // depth 1 insufficient
	if err == nil {
		t.Error("expected error when max_depth too small")
	}
}

func TestShortestPath_unknownFrom(t *testing.T) {
	g := tmpGraph(t)
	b, _ := g.AddNode("X", "B", nil)
	_, err := g.ShortestPath("nope", b.ID, 0)
	if err == nil {
		t.Error("expected error for unknown from-node")
	}
}

func TestShortestPath_unknownTo(t *testing.T) {
	g := tmpGraph(t)
	a, _ := g.AddNode("X", "A", nil)
	_, err := g.ShortestPath(a.ID, "nope", 0)
	if err == nil {
		t.Error("expected error for unknown to-node")
	}
}

// --- RemoveNode ---

func TestRemoveNode(t *testing.T) {
	g := tmpGraph(t)
	a, _ := g.AddNode("X", "A", nil)
	b, _ := g.AddNode("X", "B", nil)
	g.AddEdge(a.ID, b.ID, "rel", nil)

	deleted, err := g.RemoveNode(a.ID)
	if err != nil {
		t.Fatalf("RemoveNode: %v", err)
	}
	if deleted != 1 {
		t.Errorf("deleted edges: got %d, want 1", deleted)
	}
	if _, ok := g.GetNode(a.ID); ok {
		t.Error("node still present after removal")
	}
	if len(g.AllEdges()) != 0 {
		t.Error("edge not removed with node")
	}
}

func TestRemoveNode_notFound(t *testing.T) {
	g := tmpGraph(t)
	_, err := g.RemoveNode("nonexistent")
	if err == nil {
		t.Error("expected error for unknown node")
	}
}

// --- RemoveEdge ---

func TestRemoveEdge(t *testing.T) {
	g := tmpGraph(t)
	a, _ := g.AddNode("X", "A", nil)
	b, _ := g.AddNode("X", "B", nil)
	e, _ := g.AddEdge(a.ID, b.ID, "rel", nil)

	if err := g.RemoveEdge(e.ID); err != nil {
		t.Fatalf("RemoveEdge: %v", err)
	}
	if len(g.AllEdges()) != 0 {
		t.Error("edge still present after removal")
	}
}

func TestRemoveEdge_notFound(t *testing.T) {
	g := tmpGraph(t)
	if err := g.RemoveEdge("nonexistent"); err == nil {
		t.Error("expected error for unknown edge")
	}
}

// --- ListRelations ---

func TestListRelations(t *testing.T) {
	g := tmpGraph(t)
	a, _ := g.AddNode("X", "A", nil)
	b, _ := g.AddNode("X", "B", nil)
	c, _ := g.AddNode("X", "C", nil)
	g.AddEdge(a.ID, b.ID, "knows", nil)
	g.AddEdge(a.ID, c.ID, "uses", nil)
	g.AddEdge(b.ID, c.ID, "knows", nil)

	rels := g.ListRelations()
	if len(rels) != 2 || rels[0] != "knows" || rels[1] != "uses" {
		t.Errorf("ListRelations: got %v", rels)
	}
}

func TestListRelations_empty(t *testing.T) {
	g := tmpGraph(t)
	if len(g.ListRelations()) != 0 {
		t.Error("expected empty relations list")
	}
}

// --- Stats ---

func TestStats(t *testing.T) {
	g := tmpGraph(t)
	a, _ := g.AddNode("X", "A", nil)
	b, _ := g.AddNode("X", "B", nil)
	g.AddEdge(a.ID, b.ID, "rel", nil)

	nodes, edges := g.Stats()
	if nodes != 2 || edges != 1 {
		t.Errorf("Stats: nodes=%d edges=%d", nodes, edges)
	}
}

// --- Persistence ---

func TestPersistence(t *testing.T) {
	dir := t.TempDir()
	g, _ := New(dir)
	a, _ := g.AddNode("Person", "Alice", map[string]string{"role": "dev"})
	b, _ := g.AddNode("Person", "Bob", nil)
	g.AddEdge(a.ID, b.ID, "knows", nil)

	// reload from disk
	g2, err := New(dir)
	if err != nil {
		t.Fatalf("reload: %v", err)
	}
	n, ok := g2.GetNode(a.ID)
	if !ok {
		t.Fatal("Alice not found after reload")
	}
	if n.Props["role"] != "dev" {
		t.Errorf("props not persisted: %+v", n.Props)
	}
	nNodes, nEdges := g2.Stats()
	if nNodes != 2 || nEdges != 1 {
		t.Errorf("stats after reload: nodes=%d edges=%d", nNodes, nEdges)
	}
}

func TestNew_existingDir(t *testing.T) {
	dir := t.TempDir()
	// New should succeed on already existing dir
	_, err := New(dir)
	if err != nil {
		t.Fatalf("New on existing dir: %v", err)
	}
}

func TestNew_emptyFile(t *testing.T) {
	dir := t.TempDir()
	// An empty but present graph.json should not cause an error
	// (handled by os.IsNotExist check)
	_ = os.WriteFile(filepath.Join(dir, "graph.json"), []byte("{}"), 0o644)
	_, err := New(dir)
	if err != nil {
		t.Fatalf("New with empty JSON: %v", err)
	}
}

// --- AllNodes / AllEdges sorted ---

func TestAllNodes_sorted(t *testing.T) {
	g := tmpGraph(t)
	g.AddNode("X", "Z", nil)
	g.AddNode("X", "A", nil)
	nodes := g.AllNodes()
	for i := 1; i < len(nodes); i++ {
		if nodes[i-1].ID > nodes[i].ID {
			t.Error("AllNodes not sorted")
		}
	}
}

func TestAllEdges_sorted(t *testing.T) {
	g := tmpGraph(t)
	a, _ := g.AddNode("X", "A", nil)
	b, _ := g.AddNode("X", "B", nil)
	c, _ := g.AddNode("X", "C", nil)
	g.AddEdge(a.ID, b.ID, "r", nil)
	g.AddEdge(b.ID, c.ID, "r", nil)
	edges := g.AllEdges()
	for i := 1; i < len(edges); i++ {
		if edges[i-1].ID > edges[i].ID {
			t.Error("AllEdges not sorted")
		}
	}
}
