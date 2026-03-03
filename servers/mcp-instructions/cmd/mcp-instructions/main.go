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

	"github.com/modelcontextprotocol/go-sdk/mcp"

	"github.com/Arkestone/mcp/pkg/config"
	"github.com/Arkestone/mcp/pkg/github"
	"github.com/Arkestone/mcp/pkg/httputil"
	"github.com/Arkestone/mcp/pkg/optimizer"
	"github.com/Arkestone/mcp/pkg/server"
	"github.com/Arkestone/mcp/servers/mcp-instructions/internal/loader"
)

var version = "dev"  // set at build time via -ldflags
var commit = "none"  // set at build time via -ldflags
var date = "unknown" // set at build time via -ldflags

func main() {
	log.Printf("mcp-instructions version=%s commit=%s date=%s", version, commit, date)
	cfg := config.Load(config.Options{
		EnvPrefix:        "INSTRUCTIONS",
		DefaultAddr:      ":8080",
		DefaultCacheName: "mcp-instructions",
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
		&mcp.Implementation{Name: "mcp-instructions", Version: "0.1.0"},
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

func toOptimizerInputs(instructions []loader.Instruction) []optimizer.ContentInput {
	inputs := make([]optimizer.ContentInput, len(instructions))
	for i, inst := range instructions {
		inputs[i] = optimizer.ContentInput{
			Source:  inst.Source,
			Path:    inst.Path,
			Content: inst.Content,
		}
	}
	return inputs
}

func registerResources(srv *mcp.Server, ldr *loader.Loader, opt *optimizer.Optimizer, optimizeDefault bool) {
	srv.AddResourceTemplate(
		&mcp.ResourceTemplate{
			URITemplate: "instructions://{source}/{name}",
			Name:        "Copilot Instructions",
			Description: "Custom instruction files from configured sources",
			MIMEType:    "text/markdown",
		},
		func(ctx context.Context, req *mcp.ReadResourceRequest) (*mcp.ReadResourceResult, error) {
			inst, ok := ldr.Get(req.Params.URI)
			if !ok {
				return nil, fmt.Errorf("instruction not found: %s", req.Params.URI)
			}
			return &mcp.ReadResourceResult{
				Contents: []*mcp.ResourceContents{{
					URI: req.Params.URI, MIMEType: "text/markdown", Text: inst.Content,
				}},
			}, nil
		},
	)

	srv.AddResource(
		&mcp.Resource{
			URI: "instructions://optimized", Name: "Optimized Instructions",
			Description: "All instructions merged via LLM (or concatenated)", MIMEType: "text/markdown",
		},
		func(ctx context.Context, req *mcp.ReadResourceRequest) (*mcp.ReadResourceResult, error) {
			content := optimizeContent(ctx, opt, optimizeDefault, "", toOptimizerInputs(ldr.List()))
			return &mcp.ReadResourceResult{
				Contents: []*mcp.ResourceContents{{
					URI: "instructions://optimized", MIMEType: "text/markdown", Text: content,
				}},
			}, nil
		},
	)

	srv.AddResource(
		&mcp.Resource{
			URI: "instructions://index", Name: "Instructions Index",
			Description: "List of all available instruction files", MIMEType: "text/plain",
		},
		func(ctx context.Context, req *mcp.ReadResourceRequest) (*mcp.ReadResourceResult, error) {
			var sb strings.Builder
			for _, inst := range ldr.List() {
				fmt.Fprintf(&sb, "%s  [%s] %s\n", inst.URI, inst.Source, inst.Path)
			}
			return &mcp.ReadResourceResult{
				Contents: []*mcp.ResourceContents{{
					URI: "instructions://index", MIMEType: "text/plain", Text: sb.String(),
				}},
			}, nil
		},
	)
}

func registerPrompts(srv *mcp.Server, ldr *loader.Loader, opt *optimizer.Optimizer, optimizeDefault bool) {
	srv.AddPrompt(
		&mcp.Prompt{
			Name:        "get-instructions",
			Description: "Get all custom instructions, optionally optimized via LLM",
			Arguments: []*mcp.PromptArgument{
				{Name: "source", Description: "Filter by source (optional)", Required: false},
				{Name: "optimize", Description: "Override LLM optimization (true/false)", Required: false},
				{Name: "file_path", Description: "Return only instructions applicable to this file path (uses applyTo glob matching)", Required: false},
			},
		},
		func(ctx context.Context, req *mcp.GetPromptRequest) (*mcp.GetPromptResult, error) {
			filtered := filterBySource(ldr.List(), req.Params.Arguments["source"])
			filtered = loader.FilterByFilePath(filtered, req.Params.Arguments["file_path"])
			if len(filtered) == 0 {
				return promptResult("Custom instructions for AI assistants", "No instructions found."), nil
			}
			content := optimizeContent(ctx, opt, optimizeDefault,
				req.Params.Arguments["optimize"], toOptimizerInputs(filtered))
			return promptResult("Custom instructions for AI assistants", content), nil
		},
	)
}

func registerTools(srv *mcp.Server, ldr *loader.Loader, opt *optimizer.Optimizer, optimizeDefault bool) {
	type RefreshInput struct {
		Source string `json:"source,omitempty" jsonschema:"Optional source filter. When set only the matching source is returned in the result. All sources are always refreshed."`
	}
	type RefreshOutput struct {
		Message string   `json:"message"`
		Count   int      `json:"count"`
		Sources []string `json:"sources"`
	}
	mcp.AddTool(srv, &mcp.Tool{
		Name: "refresh-instructions",
		Description: "Force an immediate re-sync of all instruction sources. " +
			"This fetches the latest content from every configured GitHub repository into the local cache " +
			"and re-reads local directories. Call this tool when instructions may have changed on disk or " +
			"in remote repositories and you need the most up-to-date content. " +
			"Returns the total count of instruction files found after the refresh and the list of distinct sources.",
	},
		func(ctx context.Context, req *mcp.CallToolRequest, input RefreshInput) (*mcp.CallToolResult, RefreshOutput, error) {
			ldr.ForceSync()
			instructions := ldr.List()
			sourceSet := map[string]struct{}{}
			for _, inst := range instructions {
				sourceSet[inst.Source] = struct{}{}
			}
			sources := make([]string, 0, len(sourceSet))
			for s := range sourceSet {
				sources = append(sources, s)
			}
			sort.Strings(sources)
			return nil, RefreshOutput{
				Message: "All instruction sources refreshed successfully",
				Count:   len(instructions),
				Sources: sources,
			}, nil
		},
	)

	type ListInput struct {
		FilePath string `json:"file_path,omitempty" jsonschema:"Filter by file path using applyTo glob patterns. Instructions with no applyTo always match."`
		Source   string `json:"source,omitempty"   jsonschema:"Filter by source name."`
	}
	type ListEntry struct {
		URI     string   `json:"uri"`
		Source  string   `json:"source"`
		Path    string   `json:"path"`
		ApplyTo []string `json:"apply_to,omitempty"`
	}
	type ListOutput struct {
		Total   int         `json:"total"`
		Matched int         `json:"matched"`
		Entries []ListEntry `json:"entries"`
	}
	mcp.AddTool(srv, &mcp.Tool{
		Name:        "list-instructions",
		Description: "List available instruction files with their URIs and applyTo patterns. Use file_path to filter by relevance before calling get-context.",
	},
		func(ctx context.Context, req *mcp.CallToolRequest, input ListInput) (*mcp.CallToolResult, ListOutput, error) {
			all := ldr.List()
			filtered := loader.FilterByFilePath(filterBySource(all, input.Source), input.FilePath)
			entries := make([]ListEntry, len(filtered))
			for i, inst := range filtered {
				entries[i] = ListEntry{URI: inst.URI, Source: inst.Source, Path: inst.Path, ApplyTo: inst.ApplyTo}
			}
			return nil, ListOutput{Total: len(all), Matched: len(filtered), Entries: entries}, nil
		},
	)

	// get-context is the primary agent tool: returns all applicable instructions
	// merged into a single ready-to-inject block. Filter by file_path to keep
	// only the instructions relevant to the current context.
	type GetContextInput struct {
		FilePath string `json:"file_path,omitempty" jsonschema:"Current file path being edited or generated. Used to match applyTo glob patterns. Omit to load all instructions."`
		Source   string `json:"source,omitempty"   jsonschema:"Restrict to a specific source."`
		Optimize string `json:"optimize,omitempty" jsonschema:"Override LLM optimization: true or false."`
	}
	type GetContextOutput struct {
		FilePath string `json:"file_path,omitempty"`
		Matched  int    `json:"matched"`
		Total    int    `json:"total"`
		Content  string `json:"content"`
	}
	mcp.AddTool(srv, &mcp.Tool{
		Name: "get-context",
		Description: "PRIMARY AGENT TOOL. Returns all instructions applicable to the current context as a single merged block ready to inject into the agent system prompt. " +
			"Pass file_path (e.g. 'src/auth/jwt.go') to receive only instructions whose applyTo glob matches — " +
			"instructions with no applyTo are always included as global rules. " +
			"This is the most efficient way to load relevant instructions without bloating the context window.",
	},
		func(ctx context.Context, req *mcp.CallToolRequest, input GetContextInput) (*mcp.CallToolResult, GetContextOutput, error) {
			all := ldr.List()
			filtered := loader.FilterByFilePath(filterBySource(all, input.Source), input.FilePath)
			if len(filtered) == 0 {
				return nil, GetContextOutput{FilePath: input.FilePath, Matched: 0, Total: len(all), Content: ""}, nil
			}
			content := optimizeContent(ctx, opt, optimizeDefault, input.Optimize, toOptimizerInputs(filtered))
			return nil, GetContextOutput{
				FilePath: input.FilePath,
				Matched:  len(filtered),
				Total:    len(all),
				Content:  content,
			}, nil
		},
	)

	type OptimizeInput struct {
		Source   string `json:"source,omitempty"`
		FilePath string `json:"file_path,omitempty" jsonschema:"Filter by file path before optimizing."`
		Optimize string `json:"optimize,omitempty"`
	}
	type OptimizeOutput struct {
		Content string `json:"content"`
		Matched int    `json:"matched"`
	}
	mcp.AddTool(srv, &mcp.Tool{Name: "optimize-instructions", Description: "Get instructions merged via LLM (or concatenated), optionally filtered by file path and source."},
		func(ctx context.Context, req *mcp.CallToolRequest, input OptimizeInput) (*mcp.CallToolResult, OptimizeOutput, error) {
			filtered := loader.FilterByFilePath(filterBySource(ldr.List(), input.Source), input.FilePath)
			content := optimizeContent(ctx, opt, optimizeDefault, input.Optimize, toOptimizerInputs(filtered))
			return nil, OptimizeOutput{Content: content, Matched: len(filtered)}, nil
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

func filterBySource(instructions []loader.Instruction, source string) []loader.Instruction {
	if source == "" {
		return instructions
	}
	var out []loader.Instruction
	for _, inst := range instructions {
		if inst.Source == source {
			out = append(out, inst)
		}
	}
	return out
}
