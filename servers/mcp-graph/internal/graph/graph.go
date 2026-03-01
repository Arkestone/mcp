// Package graph provides a thread-safe in-memory knowledge graph
// backed by a JSON file for persistence.
// Nodes are entities with a label and name; edges are directed
// relationships between nodes with a relation type.
package graph

import (
	"encoding/json"
	"fmt"
	"math/rand"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"sync"
	"time"
)

// Node represents an entity in the knowledge graph.
type Node struct {
	ID        string            `json:"id"`
	Label     string            `json:"label"`  // e.g. "Person", "Concept", "Technology"
	Name      string            `json:"name"`   // display name
	Props     map[string]string `json:"props,omitempty"`
	CreatedAt time.Time         `json:"created_at"`
}

// Edge represents a directed relationship from one node to another.
type Edge struct {
	ID        string            `json:"id"`
	From      string            `json:"from"`
	To        string            `json:"to"`
	Relation  string            `json:"relation"` // e.g. "knows", "depends_on", "uses"
	Props     map[string]string `json:"props,omitempty"`
	CreatedAt time.Time         `json:"created_at"`
}

// PathResult holds the ordered nodes and edges of a graph path.
type PathResult struct {
	Nodes []Node `json:"nodes"`
	Edges []Edge `json:"edges"`
}

// NeighborResult holds the neighbors of a node and the connecting edges.
type NeighborResult struct {
	Nodes []Node `json:"nodes"`
	Edges []Edge `json:"edges"`
}

// diskGraph is the serialisation format for the JSON file.
type diskGraph struct {
	Nodes map[string]*Node `json:"nodes"`
	Edges map[string]*Edge `json:"edges"`
}

// Graph is a thread-safe knowledge graph backed by a JSON file.
type Graph struct {
	mu    sync.RWMutex
	file  string
	nodes map[string]*Node
	edges map[string]*Edge
	out   map[string][]string // nodeID → outgoing edge IDs
	in    map[string][]string // nodeID → incoming edge IDs
}

// New loads or creates a Graph whose data is stored in dir/graph.json.
func New(dir string) (*Graph, error) {
	if err := os.MkdirAll(dir, 0o755); err != nil {
		return nil, fmt.Errorf("creating graph dir: %w", err)
	}
	g := &Graph{
		file:  filepath.Join(dir, "graph.json"),
		nodes: make(map[string]*Node),
		edges: make(map[string]*Edge),
		out:   make(map[string][]string),
		in:    make(map[string][]string),
	}
	if err := g.load(); err != nil {
		return nil, err
	}
	return g, nil
}

// AddNode creates a new node with the given label, name, and optional properties.
func (g *Graph) AddNode(label, name string, props map[string]string) (Node, error) {
	if strings.TrimSpace(label) == "" {
		return Node{}, fmt.Errorf("label must not be empty")
	}
	if strings.TrimSpace(name) == "" {
		return Node{}, fmt.Errorf("name must not be empty")
	}
	g.mu.Lock()
	defer g.mu.Unlock()
	n := &Node{
		ID:        generateID(),
		Label:     label,
		Name:      name,
		Props:     props,
		CreatedAt: time.Now().UTC(),
	}
	g.nodes[n.ID] = n
	if err := g.save(); err != nil {
		delete(g.nodes, n.ID)
		return Node{}, err
	}
	return *n, nil
}

// AddEdge creates a directed edge from→to with the given relation type.
func (g *Graph) AddEdge(from, to, relation string, props map[string]string) (Edge, error) {
	if strings.TrimSpace(relation) == "" {
		return Edge{}, fmt.Errorf("relation must not be empty")
	}
	g.mu.Lock()
	defer g.mu.Unlock()
	if _, ok := g.nodes[from]; !ok {
		return Edge{}, fmt.Errorf("node not found: %s", from)
	}
	if _, ok := g.nodes[to]; !ok {
		return Edge{}, fmt.Errorf("node not found: %s", to)
	}
	e := &Edge{
		ID:        generateID(),
		From:      from,
		To:        to,
		Relation:  relation,
		Props:     props,
		CreatedAt: time.Now().UTC(),
	}
	g.edges[e.ID] = e
	g.out[from] = append(g.out[from], e.ID)
	g.in[to] = append(g.in[to], e.ID)
	if err := g.save(); err != nil {
		// roll back
		delete(g.edges, e.ID)
		g.out[from] = g.out[from][:len(g.out[from])-1]
		g.in[to] = g.in[to][:len(g.in[to])-1]
		return Edge{}, err
	}
	return *e, nil
}

// FindNodes returns nodes matching the optional label and/or name substring.
// Both comparisons are case-insensitive. Empty strings match everything.
func (g *Graph) FindNodes(label, name string) []Node {
	g.mu.RLock()
	defer g.mu.RUnlock()
	label = strings.ToLower(strings.TrimSpace(label))
	name = strings.ToLower(strings.TrimSpace(name))
	var out []Node
	for _, n := range g.nodes {
		if label != "" && strings.ToLower(n.Label) != label {
			continue
		}
		if name != "" && !strings.Contains(strings.ToLower(n.Name), name) {
			continue
		}
		out = append(out, *n)
	}
	sort.Slice(out, func(i, j int) bool { return out[i].ID < out[j].ID })
	return out
}

// GetNode returns a single node by ID.
func (g *Graph) GetNode(id string) (Node, bool) {
	g.mu.RLock()
	defer g.mu.RUnlock()
	n, ok := g.nodes[id]
	if !ok {
		return Node{}, false
	}
	return *n, true
}

// Neighbors returns the direct neighbors of a node and the connecting edges.
// direction is "out" (default outgoing), "in" (incoming), or "both".
// relation optionally filters by edge relation type (case-insensitive).
func (g *Graph) Neighbors(id, direction, relation string) (NeighborResult, error) {
	g.mu.RLock()
	defer g.mu.RUnlock()
	if _, ok := g.nodes[id]; !ok {
		return NeighborResult{}, fmt.Errorf("node not found: %s", id)
	}
	if direction == "" {
		direction = "both"
	}
	relation = strings.ToLower(strings.TrimSpace(relation))

	seen := map[string]bool{}
	var nodes []Node
	var edges []Edge

	addEdge := func(eid string) {
		e := g.edges[eid]
		if relation != "" && strings.ToLower(e.Relation) != relation {
			return
		}
		edges = append(edges, *e)
		// the peer is the other endpoint
		peer := e.To
		if peer == id {
			peer = e.From
		}
		if !seen[peer] {
			seen[peer] = true
			if n, ok := g.nodes[peer]; ok {
				nodes = append(nodes, *n)
			}
		}
	}

	if direction == "out" || direction == "both" {
		for _, eid := range g.out[id] {
			addEdge(eid)
		}
	}
	if direction == "in" || direction == "both" {
		for _, eid := range g.in[id] {
			addEdge(eid)
		}
	}

	sort.Slice(nodes, func(i, j int) bool { return nodes[i].ID < nodes[j].ID })
	sort.Slice(edges, func(i, j int) bool { return edges[i].ID < edges[j].ID })
	return NeighborResult{Nodes: nodes, Edges: edges}, nil
}

// ShortestPath finds the shortest directed path from → to using BFS.
// maxDepth caps the search depth (0 → use default of 10).
func (g *Graph) ShortestPath(from, to string, maxDepth int) (PathResult, error) {
	g.mu.RLock()
	defer g.mu.RUnlock()
	if _, ok := g.nodes[from]; !ok {
		return PathResult{}, fmt.Errorf("node not found: %s", from)
	}
	if _, ok := g.nodes[to]; !ok {
		return PathResult{}, fmt.Errorf("node not found: %s", to)
	}
	if from == to {
		return PathResult{Nodes: []Node{*g.nodes[from]}}, nil
	}
	if maxDepth <= 0 {
		maxDepth = 10
	}

	type state struct {
		nodeID  string
		nodeIDs []string
		edgeIDs []string
	}
	queue := []state{{nodeID: from, nodeIDs: []string{from}, edgeIDs: []string{}}}
	visited := map[string]bool{from: true}

	for len(queue) > 0 {
		cur := queue[0]
		queue = queue[1:]
		if len(cur.nodeIDs)-1 >= maxDepth {
			continue
		}
		for _, eid := range g.out[cur.nodeID] {
			e := g.edges[eid]
			if e.To == to {
				// found — reconstruct path
				nodeIDs := append(append([]string{}, cur.nodeIDs...), to)
				edgeIDs := append(append([]string{}, cur.edgeIDs...), eid)
				result := PathResult{}
				for _, nid := range nodeIDs {
					if n, ok := g.nodes[nid]; ok {
						result.Nodes = append(result.Nodes, *n)
					}
				}
				for _, eid2 := range edgeIDs {
					if e2, ok := g.edges[eid2]; ok {
						result.Edges = append(result.Edges, *e2)
					}
				}
				return result, nil
			}
			if !visited[e.To] {
				visited[e.To] = true
				newNodes := make([]string, len(cur.nodeIDs)+1)
				copy(newNodes, cur.nodeIDs)
				newNodes[len(cur.nodeIDs)] = e.To
				newEdges := make([]string, len(cur.edgeIDs)+1)
				copy(newEdges, cur.edgeIDs)
				newEdges[len(cur.edgeIDs)] = eid
				queue = append(queue, state{nodeID: e.To, nodeIDs: newNodes, edgeIDs: newEdges})
			}
		}
	}
	return PathResult{}, fmt.Errorf("no path found from %s to %s within depth %d", from, to, maxDepth)
}

// RemoveNode deletes a node and all its incident edges.
// Returns the number of edges removed.
func (g *Graph) RemoveNode(id string) (int, error) {
	g.mu.Lock()
	defer g.mu.Unlock()
	if _, ok := g.nodes[id]; !ok {
		return 0, fmt.Errorf("node not found: %s", id)
	}
	// collect all incident edge IDs
	incident := map[string]bool{}
	for _, eid := range g.out[id] {
		incident[eid] = true
	}
	for _, eid := range g.in[id] {
		incident[eid] = true
	}
	// remove edges from adjacency index and edge map
	for eid := range incident {
		e := g.edges[eid]
		g.out[e.From] = removeStr(g.out[e.From], eid)
		g.in[e.To] = removeStr(g.in[e.To], eid)
		delete(g.edges, eid)
	}
	delete(g.out, id)
	delete(g.in, id)
	delete(g.nodes, id)
	if err := g.save(); err != nil {
		return 0, err
	}
	return len(incident), nil
}

// RemoveEdge deletes an edge by ID.
func (g *Graph) RemoveEdge(id string) error {
	g.mu.Lock()
	defer g.mu.Unlock()
	e, ok := g.edges[id]
	if !ok {
		return fmt.Errorf("edge not found: %s", id)
	}
	g.out[e.From] = removeStr(g.out[e.From], id)
	g.in[e.To] = removeStr(g.in[e.To], id)
	delete(g.edges, id)
	return g.save()
}

// ListRelations returns all unique relation types in the graph, sorted.
func (g *Graph) ListRelations() []string {
	g.mu.RLock()
	defer g.mu.RUnlock()
	seen := map[string]bool{}
	for _, e := range g.edges {
		seen[e.Relation] = true
	}
	out := make([]string, 0, len(seen))
	for r := range seen {
		out = append(out, r)
	}
	sort.Strings(out)
	return out
}

// Stats returns the total number of nodes and edges.
func (g *Graph) Stats() (nodes, edges int) {
	g.mu.RLock()
	defer g.mu.RUnlock()
	return len(g.nodes), len(g.edges)
}

// AllNodes returns all nodes sorted by ID.
func (g *Graph) AllNodes() []Node {
	g.mu.RLock()
	defer g.mu.RUnlock()
	out := make([]Node, 0, len(g.nodes))
	for _, n := range g.nodes {
		out = append(out, *n)
	}
	sort.Slice(out, func(i, j int) bool { return out[i].ID < out[j].ID })
	return out
}

// AllEdges returns all edges sorted by ID.
func (g *Graph) AllEdges() []Edge {
	g.mu.RLock()
	defer g.mu.RUnlock()
	out := make([]Edge, 0, len(g.edges))
	for _, e := range g.edges {
		out = append(out, *e)
	}
	sort.Slice(out, func(i, j int) bool { return out[i].ID < out[j].ID })
	return out
}

// load reads the graph from the JSON file. A missing file starts an empty graph.
func (g *Graph) load() error {
	data, err := os.ReadFile(g.file)
	if os.IsNotExist(err) {
		return nil
	}
	if err != nil {
		return fmt.Errorf("reading graph file: %w", err)
	}
	var d diskGraph
	if err := json.Unmarshal(data, &d); err != nil {
		return fmt.Errorf("parsing graph file: %w", err)
	}
	if d.Nodes != nil {
		g.nodes = d.Nodes
	}
	if d.Edges != nil {
		g.edges = d.Edges
	}
	// rebuild adjacency index
	for _, e := range g.edges {
		g.out[e.From] = append(g.out[e.From], e.ID)
		g.in[e.To] = append(g.in[e.To], e.ID)
	}
	return nil
}

// save atomically persists the graph to the JSON file.
func (g *Graph) save() error {
	d := diskGraph{Nodes: g.nodes, Edges: g.edges}
	data, err := json.MarshalIndent(d, "", "  ")
	if err != nil {
		return fmt.Errorf("serializing graph: %w", err)
	}
	tmp := g.file + ".tmp"
	if err := os.WriteFile(tmp, data, 0o644); err != nil {
		return fmt.Errorf("writing graph: %w", err)
	}
	return os.Rename(tmp, g.file)
}

func removeStr(s []string, v string) []string {
	out := s[:0]
	for _, x := range s {
		if x != v {
			out = append(out, x)
		}
	}
	return out
}

func generateID() string {
	const chars = "abcdefghijklmnopqrstuvwxyz0123456789"
	b := make([]byte, 8)
	for i := range b {
		b[i] = chars[rand.Intn(len(chars))]
	}
	return string(b)
}
