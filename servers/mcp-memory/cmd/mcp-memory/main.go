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

	"github.com/modelcontextprotocol/go-sdk/mcp"

	"github.com/Arkestone/mcp/pkg/config"
	"github.com/Arkestone/mcp/pkg/httputil"
	"github.com/Arkestone/mcp/pkg/server"
	"github.com/Arkestone/mcp/servers/mcp-memory/internal/store"
)

var version = "dev"
var commit = "none"
var date = "unknown"

const defaultMemoryDir = "~/.local/share/mcp-memory"

func main() {
	log.Printf("mcp-memory version=%s commit=%s date=%s", version, commit, date)
	cfg := config.Load(config.Options{
		EnvPrefix:        "MEMORY",
		DefaultAddr:      ":8084",
		DefaultCacheName: "mcp-memory",
	})

	memDir := memoryDir(cfg)
	log.Printf("memory store: %s", memDir)

	st, err := store.New(memDir)
	if err != nil {
		log.Fatalf("creating memory store: %v", err)
	}

	httpClient, err := httputil.NewClient(cfg.Proxy, 30*time.Second)
	if err != nil {
		log.Fatalf("creating HTTP client: %v", err)
	}
	_ = httpClient // available for future use

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	srv := mcp.NewServer(
		&mcp.Implementation{Name: "mcp-memory", Version: version},
		nil,
	)

	registerResources(srv, st)
	registerTools(srv, st)

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

// memoryDir returns the memory store directory from config or default.
func memoryDir(cfg *config.Config) string {
	if len(cfg.Sources.Dirs) > 0 && cfg.Sources.Dirs[0] != "" {
		return expandHome(cfg.Sources.Dirs[0])
	}
	return expandHome(defaultMemoryDir)
}

func expandHome(path string) string {
	if strings.HasPrefix(path, "~/") {
		if home, err := os.UserHomeDir(); err == nil {
			return home + path[1:]
		}
	}
	return path
}

func registerResources(srv *mcp.Server, st *store.Store) {
	srv.AddResourceTemplate(
		&mcp.ResourceTemplate{
			URITemplate: "memory://{id}",
			Name:        "Memory",
			Description: "A single stored memory by ID",
			MIMEType:    "text/markdown",
		},
		func(ctx context.Context, req *mcp.ReadResourceRequest) (*mcp.ReadResourceResult, error) {
			// Extract ID from URI "memory://{id}"
			id := strings.TrimPrefix(req.Params.URI, "memory://")
			m, ok := st.Get(id)
			if !ok {
				return nil, fmt.Errorf("memory not found: %s", id)
			}
			return &mcp.ReadResourceResult{
				Contents: []*mcp.ResourceContents{{
					URI: req.Params.URI, MIMEType: "text/markdown", Text: m.Content,
				}},
			}, nil
		},
	)

	srv.AddResource(
		&mcp.Resource{
			URI: "memory://all", Name: "All Memories",
			Description: "All stored memories", MIMEType: "text/plain",
		},
		func(ctx context.Context, req *mcp.ReadResourceRequest) (*mcp.ReadResourceResult, error) {
			memories, err := st.List()
			if err != nil {
				return nil, err
			}
			var sb strings.Builder
			for _, m := range memories {
				tags := strings.Join(m.Tags, ", ")
				fmt.Fprintf(&sb, "[%s] tags=[%s] created=%s\n%s\n\n",
					m.ID, tags, m.CreatedAt.Format("2006-01-02"), m.Content)
			}
			return &mcp.ReadResourceResult{
				Contents: []*mcp.ResourceContents{{
					URI: "memory://all", MIMEType: "text/plain", Text: sb.String(),
				}},
			}, nil
		},
	)
}

func registerTools(srv *mcp.Server, st *store.Store) {
	type RememberInput struct {
		Content string   `json:"content" jsonschema:"The text to remember"`
		Tags    []string `json:"tags,omitempty" jsonschema:"Optional tags to categorize the memory"`
	}
	type RememberOutput struct {
		ID        string   `json:"id"`
		Tags      []string `json:"tags"`
		CreatedAt string   `json:"created_at"`
	}
	mcp.AddTool(srv, &mcp.Tool{
		Name:        "remember",
		Description: "Store a new memory with optional tags for later retrieval",
	},
		func(ctx context.Context, req *mcp.CallToolRequest, input RememberInput) (*mcp.CallToolResult, RememberOutput, error) {
			if strings.TrimSpace(input.Content) == "" {
				return nil, RememberOutput{}, fmt.Errorf("content must not be empty")
			}
			m, err := st.Remember(input.Content, input.Tags)
			if err != nil {
				return nil, RememberOutput{}, err
			}
			return nil, RememberOutput{
				ID: m.ID, Tags: m.Tags, CreatedAt: m.CreatedAt.Format(time.RFC3339),
			}, nil
		},
	)

	type RecallInput struct {
		Query string   `json:"query,omitempty" jsonschema:"Text to search for in memory content"`
		Tags  []string `json:"tags,omitempty" jsonschema:"Filter by tags all must match"`
	}
	type RecallEntry struct {
		ID        string   `json:"id"`
		Tags      []string `json:"tags"`
		Content   string   `json:"content"`
		CreatedAt string   `json:"created_at"`
	}
	type RecallOutput struct {
		Memories []RecallEntry `json:"memories"`
		Count    int           `json:"count"`
	}
	mcp.AddTool(srv, &mcp.Tool{
		Name:        "recall",
		Description: "Search memories by text query and/or tags",
	},
		func(ctx context.Context, req *mcp.CallToolRequest, input RecallInput) (*mcp.CallToolResult, RecallOutput, error) {
			memories, err := st.Recall(input.Query, input.Tags)
			if err != nil {
				return nil, RecallOutput{}, err
			}
			entries := make([]RecallEntry, len(memories))
			for i, m := range memories {
				entries[i] = RecallEntry{
					ID: m.ID, Tags: m.Tags, Content: m.Content,
					CreatedAt: m.CreatedAt.Format(time.RFC3339),
				}
			}
			return nil, RecallOutput{Memories: entries, Count: len(entries)}, nil
		},
	)

	type ForgetInput struct {
		ID string `json:"id" jsonschema:"The ID of the memory to delete"`
	}
	type ForgetOutput struct {
		Message string `json:"message"`
	}
	mcp.AddTool(srv, &mcp.Tool{
		Name:        "forget",
		Description: "Delete a memory by ID",
	},
		func(ctx context.Context, req *mcp.CallToolRequest, input ForgetInput) (*mcp.CallToolResult, ForgetOutput, error) {
			if err := st.Forget(input.ID); err != nil {
				return nil, ForgetOutput{}, err
			}
			return nil, ForgetOutput{Message: fmt.Sprintf("Memory %s deleted", input.ID)}, nil
		},
	)

	type ListMemoriesInput struct {
		Tags []string `json:"tags,omitempty"`
	}
	type ListEntry struct {
		ID        string   `json:"id"`
		Tags      []string `json:"tags"`
		Content   string   `json:"content"`
		CreatedAt string   `json:"created_at"`
	}
	type ListOutput struct {
		Memories []ListEntry `json:"memories"`
		Count    int         `json:"count"`
	}
	mcp.AddTool(srv, &mcp.Tool{
		Name:        "list-memories",
		Description: "List all stored memories, optionally filtered by tags",
	},
		func(ctx context.Context, req *mcp.CallToolRequest, input ListMemoriesInput) (*mcp.CallToolResult, ListOutput, error) {
			memories, err := st.Recall("", input.Tags)
			if err != nil {
				return nil, ListOutput{}, err
			}
			entries := make([]ListEntry, len(memories))
			for i, m := range memories {
				entries[i] = ListEntry{
					ID: m.ID, Tags: m.Tags, Content: m.Content,
					CreatedAt: m.CreatedAt.Format(time.RFC3339),
				}
			}
			return nil, ListOutput{Memories: entries, Count: len(entries)}, nil
		},
	)
}
