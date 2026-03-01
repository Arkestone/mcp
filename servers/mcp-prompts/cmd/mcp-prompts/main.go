package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"
	"sort"
	"strings"
	"syscall"
	"time"

	"github.com/Arkestone/mcp/pkg/config"
	"github.com/Arkestone/mcp/pkg/github"
	"github.com/Arkestone/mcp/pkg/httputil"
	"github.com/Arkestone/mcp/pkg/optimizer"
	"github.com/Arkestone/mcp/pkg/server"
	"github.com/Arkestone/mcp/servers/mcp-prompts/internal/loader"
	"github.com/modelcontextprotocol/go-sdk/mcp"
)

var version = "dev"
var commit = "none"
var date = "unknown"

func main() {
	log.Printf("mcp-prompts version=%s commit=%s date=%s", version, commit, date)
	cfg := config.Load(config.Options{
		EnvPrefix:        "PROMPTS",
		DefaultAddr:      ":8082",
		DefaultCacheName: "mcp-prompts",
	})

	httpClient, err := httputil.NewClient(cfg.Proxy, 30*time.Second)
	if err != nil {
		log.Fatalf("creating HTTP client: %v", err)
	}

	gh := &github.Client{Token: cfg.GitHubToken, HTTPClient: httpClient}
	ldr := loader.New(cfg, gh)

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	ldr.Start(ctx)
	defer ldr.Stop()

	opt := optimizer.New(cfg.LLM, httpClient)
	if opt.Enabled() {
		log.Println("LLM optimization enabled")
	}

	srv := mcp.NewServer(
		&mcp.Implementation{Name: "mcp-prompts", Version: "0.1.0"},
		nil,
	)

	registerResources(srv, ldr, opt, cfg.LLM.Enabled)
	registerPrompts(srv, ldr, opt, cfg.LLM.Enabled)
	registerTools(srv, ldr, opt, cfg.LLM.Enabled)

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

func toOptimizerInputs(prompts []loader.Prompt) []optimizer.ContentInput {
	inputs := make([]optimizer.ContentInput, len(prompts))
	for i, p := range prompts {
		inputs[i] = optimizer.ContentInput{
			Source:  p.Source,
			Path:    p.Path,
			Content: p.Content,
		}
	}
	return inputs
}

func registerResources(srv *mcp.Server, ldr *loader.Loader, opt *optimizer.Optimizer, optimizeDefault bool) {
	srv.AddResourceTemplate(
		&mcp.ResourceTemplate{
			URITemplate: "prompts://{source}/{name}",
			Name:        "Copilot Prompt",
			Description: "A Copilot prompt or chat mode file from configured sources",
			MIMEType:    "text/markdown",
		},
		func(ctx context.Context, req *mcp.ReadResourceRequest) (*mcp.ReadResourceResult, error) {
			p, ok := ldr.Get(req.Params.URI)
			if !ok {
				return nil, fmt.Errorf("prompt not found: %s", req.Params.URI)
			}
			return &mcp.ReadResourceResult{
				Contents: []*mcp.ResourceContents{{
					URI: req.Params.URI, MIMEType: "text/markdown", Text: p.Content,
				}},
			}, nil
		},
	)

	srv.AddResource(
		&mcp.Resource{
			URI: "prompts://optimized", Name: "Optimized Prompts",
			Description: "All prompts merged via LLM (or concatenated)", MIMEType: "text/markdown",
		},
		func(ctx context.Context, req *mcp.ReadResourceRequest) (*mcp.ReadResourceResult, error) {
			content := optimizeContent(ctx, opt, optimizeDefault, "", toOptimizerInputs(ldr.List()))
			return &mcp.ReadResourceResult{
				Contents: []*mcp.ResourceContents{{
					URI: "prompts://optimized", MIMEType: "text/markdown", Text: content,
				}},
			}, nil
		},
	)

	srv.AddResource(
		&mcp.Resource{
			URI: "prompts://index", Name: "Prompts Index",
			Description: "List of all available prompt and chat mode files", MIMEType: "text/plain",
		},
		func(ctx context.Context, req *mcp.ReadResourceRequest) (*mcp.ReadResourceResult, error) {
			var sb strings.Builder
			for _, p := range ldr.List() {
				fmt.Fprintf(&sb, "%s  [%s] %s  type=%s\n", p.URI, p.Source, p.Path, p.Type)
			}
			return &mcp.ReadResourceResult{
				Contents: []*mcp.ResourceContents{{
					URI: "prompts://index", MIMEType: "text/plain", Text: sb.String(),
				}},
			}, nil
		},
	)
}

func registerPrompts(srv *mcp.Server, ldr *loader.Loader, opt *optimizer.Optimizer, optimizeDefault bool) {
	srv.AddPrompt(
		&mcp.Prompt{
			Name:        "get-prompts",
			Description: "Get all Copilot prompt and chat mode files, optionally optimized via LLM",
			Arguments: []*mcp.PromptArgument{
				{Name: "source", Description: "Filter by source (optional)", Required: false},
				{Name: "optimize", Description: "Override LLM optimization (true/false)", Required: false},
			},
		},
		func(ctx context.Context, req *mcp.GetPromptRequest) (*mcp.GetPromptResult, error) {
			filtered := filterBySource(ldr.List(), req.Params.Arguments["source"])
			if len(filtered) == 0 {
				return promptResult("Copilot prompts and chat modes", "No prompts found."), nil
			}
			content := optimizeContent(ctx, opt, optimizeDefault,
				req.Params.Arguments["optimize"], toOptimizerInputs(filtered))
			return promptResult("Copilot prompts and chat modes", content), nil
		},
	)
}

func registerTools(srv *mcp.Server, ldr *loader.Loader, opt *optimizer.Optimizer, optimizeDefault bool) {
	type RefreshInput struct {
		Source string `json:"source,omitempty"`
	}
	type RefreshOutput struct {
		Message string   `json:"message"`
		Count   int      `json:"count"`
		Sources []string `json:"sources"`
	}
	mcp.AddTool(srv, &mcp.Tool{
		Name:        "refresh-prompts",
		Description: "Force an immediate re-sync of all prompt sources from GitHub repositories and local directories.",
	},
		func(ctx context.Context, req *mcp.CallToolRequest, input RefreshInput) (*mcp.CallToolResult, RefreshOutput, error) {
			ldr.ForceSync()
			prompts := ldr.List()
			sourceSet := map[string]struct{}{}
			for _, p := range prompts {
				sourceSet[p.Source] = struct{}{}
			}
			sources := make([]string, 0, len(sourceSet))
			for s := range sourceSet {
				sources = append(sources, s)
			}
			sort.Strings(sources)
			return nil, RefreshOutput{
				Message: "All prompt sources refreshed successfully",
				Count:   len(prompts),
				Sources: sources,
			}, nil
		},
	)

	type ListInput struct{}
	type ListEntry struct {
		URI         string `json:"uri"`
		Source      string `json:"source"`
		Path        string `json:"path"`
		Type        string `json:"type"`
		Description string `json:"description"`
		Mode        string `json:"mode,omitempty"`
	}
	type ListOutput struct {
		Entries []ListEntry `json:"entries"`
	}
	mcp.AddTool(srv, &mcp.Tool{Name: "list-prompts", Description: "List all prompt and chat mode files"},
		func(ctx context.Context, req *mcp.CallToolRequest, input ListInput) (*mcp.CallToolResult, ListOutput, error) {
			prompts := ldr.List()
			entries := make([]ListEntry, len(prompts))
			for i, p := range prompts {
				entries[i] = ListEntry{URI: p.URI, Source: p.Source, Path: p.Path, Type: p.Type, Description: p.Description, Mode: p.Mode}
			}
			return nil, ListOutput{Entries: entries}, nil
		},
	)

	type GetPromptInput struct {
		URI string `json:"uri"`
	}
	type GetPromptOutput struct {
		URI         string `json:"uri"`
		Name        string `json:"name"`
		Description string `json:"description"`
		Mode        string `json:"mode,omitempty"`
		Type        string `json:"type"`
		Content     string `json:"content"`
	}
	mcp.AddTool(srv, &mcp.Tool{Name: "get-prompt", Description: "Get a single prompt or chat mode by URI"},
		func(ctx context.Context, req *mcp.CallToolRequest, input GetPromptInput) (*mcp.CallToolResult, GetPromptOutput, error) {
			p, ok := ldr.Get(input.URI)
			if !ok {
				return nil, GetPromptOutput{}, fmt.Errorf("prompt not found: %s", input.URI)
			}
			return nil, GetPromptOutput{
				URI: p.URI, Name: p.Name, Description: p.Description,
				Mode: p.Mode, Type: p.Type, Content: p.Content,
			}, nil
		},
	)

	type OptimizeInput struct {
		Source   string `json:"source,omitempty"`
		Optimize string `json:"optimize,omitempty"`
	}
	type OptimizeOutput struct {
		Content string `json:"content"`
		Sources int    `json:"sources"`
	}
	mcp.AddTool(srv, &mcp.Tool{Name: "optimize-prompts", Description: "Get consolidated prompts"},
		func(ctx context.Context, req *mcp.CallToolRequest, input OptimizeInput) (*mcp.CallToolResult, OptimizeOutput, error) {
			filtered := filterBySource(ldr.List(), input.Source)
			content := optimizeContent(ctx, opt, optimizeDefault, input.Optimize, toOptimizerInputs(filtered))
			return nil, OptimizeOutput{Content: content, Sources: len(filtered)}, nil
		},
	)
}

// optimizeContent runs LLM optimization or falls back to concatenation.
func optimizeContent(ctx context.Context, opt *optimizer.Optimizer, defaultOn bool, override string, inputs []optimizer.ContentInput) string {
	if server.ShouldOptimize(opt, defaultOn, override) {
		content, err := opt.Optimize(ctx, inputs)
		if err != nil {
			log.Printf("LLM optimization failed, falling back: %v", err)
			return optimizer.ConcatRaw(inputs)
		}
		return content
	}
	return optimizer.ConcatRaw(inputs)
}

func promptResult(desc, text string) *mcp.GetPromptResult {
	return &mcp.GetPromptResult{
		Description: desc,
		Messages:    []*mcp.PromptMessage{{Role: "user", Content: &mcp.TextContent{Text: text}}},
	}
}

func filterBySource(prompts []loader.Prompt, source string) []loader.Prompt {
	if source == "" {
		return prompts
	}
	var out []loader.Prompt
	for _, p := range prompts {
		if p.Source == source {
			out = append(out, p)
		}
	}
	return out
}
