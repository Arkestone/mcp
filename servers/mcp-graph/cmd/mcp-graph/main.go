package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	"github.com/Arkestone/mcp/pkg/config"
	"github.com/Arkestone/mcp/pkg/server"
	"github.com/Arkestone/mcp/servers/mcp-graph/internal/graph"
	"github.com/modelcontextprotocol/go-sdk/mcp"
)

var version = "dev"
var commit = "none"
var date = "unknown"

const defaultGraphDir = "~/.local/share/mcp-graph"

func main() {
	log.Printf("mcp-graph version=%s commit=%s date=%s", version, commit, date)
	cfg := config.Load(config.Options{
		EnvPrefix:        "GRAPH",
		DefaultAddr:      ":8085",
		DefaultCacheName: "mcp-graph",
	})

	dataDir := graphDir(cfg)
	log.Printf("graph store: %s", dataDir)

	g, err := graph.New(dataDir)
	if err != nil {
		log.Fatalf("creating graph store: %v", err)
	}

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	srv := mcp.NewServer(
		&mcp.Implementation{Name: "mcp-graph", Version: version},
		nil,
	)

	registerResources(srv, g)
	registerTools(srv, g)

	switch cfg.Transport {
	case "stdio":
		log.Println("starting MCP server on stdio")
		if err := srv.Run(ctx, &mcp.StdioTransport{}); err != nil {
			log.Fatalf("server error: %v", err)
		}
	case "http":
		log.Printf("starting MCP server on %s", cfg.Addr)
		if err := server.RunHTTP(ctx, srv, cfg.Addr, cfg.Proxy.HeaderPassthrough); err != nil {
			log.Fatalf("server error: %v", err)
		}
	default:
		log.Fatalf("unknown transport: %s", cfg.Transport)
	}
}

func graphDir(cfg *config.Config) string {
	if len(cfg.Sources.Dirs) > 0 && cfg.Sources.Dirs[0] != "" {
		return expandHome(cfg.Sources.Dirs[0])
	}
	return expandHome(defaultGraphDir)
}

func expandHome(path string) string {
	if strings.HasPrefix(path, "~/") {
		if home, err := os.UserHomeDir(); err == nil {
			return home + path[1:]
		}
	}
	return path
}

func registerResources(srv *mcp.Server, g *graph.Graph) {
	srv.AddResource(
		&mcp.Resource{
			URI:         "graph://stats",
			Name:        "Graph Stats",
			Description: "Total node and edge counts and all relation types in the graph",
			MIMEType:    "application/json",
		},
		func(ctx context.Context, req *mcp.ReadResourceRequest) (*mcp.ReadResourceResult, error) {
			nodes, edges := g.Stats()
			rels := g.ListRelations()
			text := fmt.Sprintf(`{"nodes":%d,"edges":%d,"relations":%s}`,
				nodes, edges, jsonStringSlice(rels))
			return &mcp.ReadResourceResult{
				Contents: []*mcp.ResourceContents{{
					URI: "graph://stats", MIMEType: "application/json", Text: text,
				}},
			}, nil
		},
	)

	srv.AddResourceTemplate(
		&mcp.ResourceTemplate{
			URITemplate: "graph://node/{id}",
			Name:        "Graph Node",
			Description: "A single node in the knowledge graph",
			MIMEType:    "application/json",
		},
		func(ctx context.Context, req *mcp.ReadResourceRequest) (*mcp.ReadResourceResult, error) {
			id := strings.TrimPrefix(req.Params.URI, "graph://node/")
			n, ok := g.GetNode(id)
			if !ok {
				return nil, fmt.Errorf("node not found: %s", id)
			}
			text := fmt.Sprintf(`{"id":%q,"label":%q,"name":%q,"created_at":%q}`,
				n.ID, n.Label, n.Name, n.CreatedAt.Format(time.RFC3339))
			return &mcp.ReadResourceResult{
				Contents: []*mcp.ResourceContents{{
					URI: req.Params.URI, MIMEType: "application/json", Text: text,
				}},
			}, nil
		},
	)
}

func registerTools(srv *mcp.Server, g *graph.Graph) {
	// add-node
	type AddNodeInput struct {
		Label string            `json:"label" jsonschema:"required,description=Node type/category (e.g. 'Person' 'Technology' 'Concept')"`
		Name  string            `json:"name"  jsonschema:"required,description=Display name of the entity"`
		Props map[string]string `json:"props,omitempty" jsonschema:"description=Optional key-value properties"`
	}
	type NodeOutput struct {
		ID        string            `json:"id"`
		Label     string            `json:"label"`
		Name      string            `json:"name"`
		Props     map[string]string `json:"props,omitempty"`
		CreatedAt string            `json:"created_at"`
	}
	mcp.AddTool(srv, &mcp.Tool{
		Name:        "add-node",
		Description: "Add a new entity (node) to the knowledge graph",
	}, func(ctx context.Context, req *mcp.CallToolRequest, in AddNodeInput) (*mcp.CallToolResult, NodeOutput, error) {
		n, err := g.AddNode(in.Label, in.Name, in.Props)
		if err != nil {
			return nil, NodeOutput{}, err
		}
		return nil, NodeOutput{ID: n.ID, Label: n.Label, Name: n.Name, Props: n.Props, CreatedAt: n.CreatedAt.Format(time.RFC3339)}, nil
	})

	// add-edge
	type AddEdgeInput struct {
		From     string            `json:"from"     jsonschema:"required,description=ID of the source node"`
		To       string            `json:"to"       jsonschema:"required,description=ID of the target node"`
		Relation string            `json:"relation" jsonschema:"required,description=Relationship type (e.g. 'knows' 'depends_on' 'uses')"`
		Props    map[string]string `json:"props,omitempty" jsonschema:"description=Optional key-value properties on the edge"`
	}
	type EdgeOutput struct {
		ID        string            `json:"id"`
		From      string            `json:"from"`
		To        string            `json:"to"`
		Relation  string            `json:"relation"`
		Props     map[string]string `json:"props,omitempty"`
		CreatedAt string            `json:"created_at"`
	}
	mcp.AddTool(srv, &mcp.Tool{
		Name:        "add-edge",
		Description: "Create a directed relationship (edge) between two existing nodes",
	}, func(ctx context.Context, req *mcp.CallToolRequest, in AddEdgeInput) (*mcp.CallToolResult, EdgeOutput, error) {
		e, err := g.AddEdge(in.From, in.To, in.Relation, in.Props)
		if err != nil {
			return nil, EdgeOutput{}, err
		}
		return nil, EdgeOutput{ID: e.ID, From: e.From, To: e.To, Relation: e.Relation, Props: e.Props, CreatedAt: e.CreatedAt.Format(time.RFC3339)}, nil
	})

	// find-nodes
	type FindNodesInput struct {
		Label string `json:"label,omitempty" jsonschema:"description=Filter by node label (case-insensitive exact match)"`
		Name  string `json:"name,omitempty"  jsonschema:"description=Filter by name substring (case-insensitive)"`
	}
	type FindNodesOutput struct {
		Nodes []NodeOutput `json:"nodes"`
		Count int          `json:"count"`
	}
	mcp.AddTool(srv, &mcp.Tool{
		Name:        "find-nodes",
		Description: "Search nodes by label and/or name substring",
	}, func(ctx context.Context, req *mcp.CallToolRequest, in FindNodesInput) (*mcp.CallToolResult, FindNodesOutput, error) {
		nodes := g.FindNodes(in.Label, in.Name)
		out := make([]NodeOutput, len(nodes))
		for i, n := range nodes {
			out[i] = NodeOutput{ID: n.ID, Label: n.Label, Name: n.Name, Props: n.Props, CreatedAt: n.CreatedAt.Format(time.RFC3339)}
		}
		return nil, FindNodesOutput{Nodes: out, Count: len(out)}, nil
	})

	// get-node
	type GetNodeInput struct {
		ID string `json:"id" jsonschema:"required,description=Node ID"`
	}
	mcp.AddTool(srv, &mcp.Tool{
		Name:        "get-node",
		Description: "Get a node by ID including its properties",
	}, func(ctx context.Context, req *mcp.CallToolRequest, in GetNodeInput) (*mcp.CallToolResult, NodeOutput, error) {
		n, ok := g.GetNode(in.ID)
		if !ok {
			return nil, NodeOutput{}, fmt.Errorf("node not found: %s", in.ID)
		}
		return nil, NodeOutput{ID: n.ID, Label: n.Label, Name: n.Name, Props: n.Props, CreatedAt: n.CreatedAt.Format(time.RFC3339)}, nil
	})

	// neighbors
	type NeighborsInput struct {
		ID        string `json:"id"                  jsonschema:"required,description=Node ID"`
		Direction string `json:"direction,omitempty" jsonschema:"description=Edge direction: 'out' 'in' or 'both' (default)"`
		Relation  string `json:"relation,omitempty"  jsonschema:"description=Optional relation type filter (case-insensitive)"`
	}
	type NeighborsOutput struct {
		Nodes []NodeOutput `json:"nodes"`
		Edges []EdgeOutput `json:"edges"`
		Count int          `json:"count"`
	}
	mcp.AddTool(srv, &mcp.Tool{
		Name:        "neighbors",
		Description: "List direct neighbors of a node with connecting edges. Use direction='out'/'in'/'both' and optionally filter by relation type.",
	}, func(ctx context.Context, req *mcp.CallToolRequest, in NeighborsInput) (*mcp.CallToolResult, NeighborsOutput, error) {
		res, err := g.Neighbors(in.ID, in.Direction, in.Relation)
		if err != nil {
			return nil, NeighborsOutput{}, err
		}
		outNodes := make([]NodeOutput, len(res.Nodes))
		for i, n := range res.Nodes {
			outNodes[i] = NodeOutput{ID: n.ID, Label: n.Label, Name: n.Name, Props: n.Props, CreatedAt: n.CreatedAt.Format(time.RFC3339)}
		}
		outEdges := make([]EdgeOutput, len(res.Edges))
		for i, e := range res.Edges {
			outEdges[i] = EdgeOutput{ID: e.ID, From: e.From, To: e.To, Relation: e.Relation, Props: e.Props, CreatedAt: e.CreatedAt.Format(time.RFC3339)}
		}
		return nil, NeighborsOutput{Nodes: outNodes, Edges: outEdges, Count: len(outNodes)}, nil
	})

	// shortest-path
	type ShortestPathInput struct {
		From     string `json:"from"               jsonschema:"required,description=Source node ID"`
		To       string `json:"to"                 jsonschema:"required,description=Target node ID"`
		MaxDepth int    `json:"max_depth,omitempty" jsonschema:"description=Maximum path length (default 10)"`
	}
	type PathOutput struct {
		Nodes []NodeOutput `json:"nodes"`
		Edges []EdgeOutput `json:"edges"`
		Hops  int          `json:"hops"`
	}
	mcp.AddTool(srv, &mcp.Tool{
		Name:        "shortest-path",
		Description: "Find the shortest directed path between two nodes using BFS",
	}, func(ctx context.Context, req *mcp.CallToolRequest, in ShortestPathInput) (*mcp.CallToolResult, PathOutput, error) {
		res, err := g.ShortestPath(in.From, in.To, in.MaxDepth)
		if err != nil {
			return nil, PathOutput{}, err
		}
		outNodes := make([]NodeOutput, len(res.Nodes))
		for i, n := range res.Nodes {
			outNodes[i] = NodeOutput{ID: n.ID, Label: n.Label, Name: n.Name, Props: n.Props, CreatedAt: n.CreatedAt.Format(time.RFC3339)}
		}
		outEdges := make([]EdgeOutput, len(res.Edges))
		for i, e := range res.Edges {
			outEdges[i] = EdgeOutput{ID: e.ID, From: e.From, To: e.To, Relation: e.Relation, Props: e.Props, CreatedAt: e.CreatedAt.Format(time.RFC3339)}
		}
		return nil, PathOutput{Nodes: outNodes, Edges: outEdges, Hops: len(outEdges)}, nil
	})

	// remove-node
	type RemoveNodeInput struct {
		ID string `json:"id" jsonschema:"required,description=ID of the node to delete"`
	}
	type RemoveNodeOutput struct {
		Message      string `json:"message"`
		DeletedEdges int    `json:"deleted_edges"`
	}
	mcp.AddTool(srv, &mcp.Tool{
		Name:        "remove-node",
		Description: "Delete a node and all its incident edges",
	}, func(ctx context.Context, req *mcp.CallToolRequest, in RemoveNodeInput) (*mcp.CallToolResult, RemoveNodeOutput, error) {
		n, err := g.RemoveNode(in.ID)
		if err != nil {
			return nil, RemoveNodeOutput{}, err
		}
		return nil, RemoveNodeOutput{
			Message:      fmt.Sprintf("node %s deleted", in.ID),
			DeletedEdges: n,
		}, nil
	})

	// remove-edge
	type RemoveEdgeInput struct {
		ID string `json:"id" jsonschema:"required,description=ID of the edge to delete"`
	}
	type RemoveEdgeOutput struct {
		Message string `json:"message"`
	}
	mcp.AddTool(srv, &mcp.Tool{
		Name:        "remove-edge",
		Description: "Delete a relationship (edge) by ID",
	}, func(ctx context.Context, req *mcp.CallToolRequest, in RemoveEdgeInput) (*mcp.CallToolResult, RemoveEdgeOutput, error) {
		if err := g.RemoveEdge(in.ID); err != nil {
			return nil, RemoveEdgeOutput{}, err
		}
		return nil, RemoveEdgeOutput{Message: fmt.Sprintf("edge %s deleted", in.ID)}, nil
	})

	// list-relations
	type ListRelationsInput struct{}
	type ListRelationsOutput struct {
		Relations []string `json:"relations"`
		Count     int      `json:"count"`
	}
	mcp.AddTool(srv, &mcp.Tool{
		Name:        "list-relations",
		Description: "List all unique relationship types present in the graph",
	}, func(ctx context.Context, req *mcp.CallToolRequest, in ListRelationsInput) (*mcp.CallToolResult, ListRelationsOutput, error) {
		rels := g.ListRelations()
		return nil, ListRelationsOutput{Relations: rels, Count: len(rels)}, nil
	})
}

// jsonStringSlice encodes a []string as a JSON array without importing encoding/json in a hot path.
func jsonStringSlice(ss []string) string {
	if len(ss) == 0 {
		return "[]"
	}
	var sb strings.Builder
	sb.WriteByte('[')
	for i, s := range ss {
		if i > 0 {
			sb.WriteByte(',')
		}
		sb.WriteByte('"')
		sb.WriteString(strings.ReplaceAll(s, `"`, `\"`))
		sb.WriteByte('"')
	}
	sb.WriteByte(']')
	return sb.String()
}
