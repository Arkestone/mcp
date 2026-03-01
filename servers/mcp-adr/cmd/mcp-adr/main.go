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
	"github.com/Arkestone/mcp/servers/mcp-adr/internal/scanner"
)

var version = "dev"
var commit = "none"
var date = "unknown"

func main() {
	log.Printf("mcp-adr version=%s commit=%s date=%s", version, commit, date)
	cfg := config.Load(config.Options{
		EnvPrefix:        "ADR",
		DefaultAddr:      ":8083",
		DefaultCacheName: "mcp-adr",
	})

	httpClient, err := httputil.NewClient(cfg.Proxy, 30*time.Second)
	if err != nil {
		log.Fatalf("creating HTTP client: %v", err)
	}

	gh := &github.Client{Token: cfg.GitHubToken, HTTPClient: httpClient}
	scn := scanner.New(cfg, gh)

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	scn.Start(ctx)
	defer scn.Stop()

	opt := optimizer.New(cfg.LLM, httpClient)
	if opt.Enabled() {
		log.Println("LLM optimization enabled")
	}

	srv := mcp.NewServer(
		&mcp.Implementation{Name: "mcp-adr", Version: "0.1.0"},
		nil,
	)

	registerResources(srv, scn, opt, cfg.LLM.Enabled)
	registerPrompts(srv, scn, opt, cfg.LLM.Enabled)
	registerTools(srv, scn, opt, cfg.LLM.Enabled)

	switch cfg.Transport {
	case "stdio":
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

func toOptimizerInputs(adrs []scanner.ADR) []optimizer.ContentInput {
	inputs := make([]optimizer.ContentInput, len(adrs))
	for i, a := range adrs {
		inputs[i] = optimizer.ContentInput{Source: a.Source, Path: a.Path, Content: a.Content}
	}
	return inputs
}

func registerResources(srv *mcp.Server, scn *scanner.Scanner, opt *optimizer.Optimizer, optimizeDefault bool) {
	srv.AddResourceTemplate(
		&mcp.ResourceTemplate{
			URITemplate: "adrs://{source}/{id}",
			Name:        "Architecture Decision Record",
			Description: "An ADR from configured sources",
			MIMEType:    "text/markdown",
		},
		func(ctx context.Context, req *mcp.ReadResourceRequest) (*mcp.ReadResourceResult, error) {
			a, ok := scn.Get(req.Params.URI)
			if !ok {
				return nil, fmt.Errorf("ADR not found: %s", req.Params.URI)
			}
			return &mcp.ReadResourceResult{
				Contents: []*mcp.ResourceContents{{URI: req.Params.URI, MIMEType: "text/markdown", Text: a.Content}},
			}, nil
		},
	)

	srv.AddResource(
		&mcp.Resource{URI: "adrs://optimized", Name: "Optimized ADRs",
			Description: "All ADRs merged via LLM (or concatenated)", MIMEType: "text/markdown"},
		func(ctx context.Context, req *mcp.ReadResourceRequest) (*mcp.ReadResourceResult, error) {
			content := optimizeContent(ctx, opt, optimizeDefault, "", toOptimizerInputs(scn.List()))
			return &mcp.ReadResourceResult{
				Contents: []*mcp.ResourceContents{{URI: "adrs://optimized", MIMEType: "text/markdown", Text: content}},
			}, nil
		},
	)

	srv.AddResource(
		&mcp.Resource{URI: "adrs://index", Name: "ADRs Index",
			Description: "List of all available ADRs", MIMEType: "text/plain"},
		func(ctx context.Context, req *mcp.ReadResourceRequest) (*mcp.ReadResourceResult, error) {
			var sb strings.Builder
			for _, a := range scn.List() {
				fmt.Fprintf(&sb, "%s  [%s] %s  status=%s\n", a.URI, a.Source, a.Title, a.Status)
			}
			return &mcp.ReadResourceResult{
				Contents: []*mcp.ResourceContents{{URI: "adrs://index", MIMEType: "text/plain", Text: sb.String()}},
			}, nil
		},
	)
}

func registerPrompts(srv *mcp.Server, scn *scanner.Scanner, opt *optimizer.Optimizer, optimizeDefault bool) {
	srv.AddPrompt(
		&mcp.Prompt{
			Name:        "get-adrs",
			Description: "Get Architecture Decision Records, optionally filtered by source or status",
			Arguments: []*mcp.PromptArgument{
				{Name: "source", Description: "Filter by source (optional)", Required: false},
				{Name: "status", Description: "Filter by status: proposed, accepted, deprecated (optional)", Required: false},
				{Name: "optimize", Description: "Override LLM optimization (true/false)", Required: false},
			},
		},
		func(ctx context.Context, req *mcp.GetPromptRequest) (*mcp.GetPromptResult, error) {
			filtered := filterBySource(scn.List(), req.Params.Arguments["source"])
			filtered = filterByStatus(filtered, req.Params.Arguments["status"])
			if len(filtered) == 0 {
				return promptResult("Architecture Decision Records", "No ADRs found."), nil
			}
			content := optimizeContent(ctx, opt, optimizeDefault,
				req.Params.Arguments["optimize"], toOptimizerInputs(filtered))
			return promptResult("Architecture Decision Records", content), nil
		},
	)
}

func registerTools(srv *mcp.Server, scn *scanner.Scanner, opt *optimizer.Optimizer, optimizeDefault bool) {
	type RefreshInput struct{}
	type RefreshOutput struct {
		Message string   `json:"message"`
		Count   int      `json:"count"`
		Sources []string `json:"sources"`
	}
	mcp.AddTool(srv, &mcp.Tool{Name: "refresh-adrs", Description: "Force re-sync of all ADR sources"},
		func(ctx context.Context, req *mcp.CallToolRequest, input RefreshInput) (*mcp.CallToolResult, RefreshOutput, error) {
			scn.ForceSync()
			adrs := scn.List()
			sourceSet := map[string]struct{}{}
			for _, a := range adrs {
				sourceSet[a.Source] = struct{}{}
			}
			sources := make([]string, 0, len(sourceSet))
			for s := range sourceSet {
				sources = append(sources, s)
			}
			sort.Strings(sources)
			return nil, RefreshOutput{Message: "All ADR sources refreshed successfully", Count: len(adrs), Sources: sources}, nil
		},
	)

	type ListInput struct {
		Source string `json:"source,omitempty"`
		Status string `json:"status,omitempty"`
	}
	type ListEntry struct {
		URI    string `json:"uri"`
		Source string `json:"source"`
		ID     string `json:"id"`
		Title  string `json:"title"`
		Status string `json:"status"`
		Date   string `json:"date,omitempty"`
	}
	type ListOutput struct {
		Entries []ListEntry `json:"entries"`
	}
	mcp.AddTool(srv, &mcp.Tool{Name: "list-adrs", Description: "List all Architecture Decision Records"},
		func(ctx context.Context, req *mcp.CallToolRequest, input ListInput) (*mcp.CallToolResult, ListOutput, error) {
			adrs := filterByStatus(filterBySource(scn.List(), input.Source), input.Status)
			entries := make([]ListEntry, len(adrs))
			for i, a := range adrs {
				entries[i] = ListEntry{URI: a.URI, Source: a.Source, ID: a.ID, Title: a.Title, Status: a.Status, Date: a.Date}
			}
			return nil, ListOutput{Entries: entries}, nil
		},
	)

	type GetADRInput struct {
		URI string `json:"uri" jsonschema:"required"`
	}
	type GetADROutput struct {
		URI     string `json:"uri"`
		ID      string `json:"id"`
		Title   string `json:"title"`
		Status  string `json:"status"`
		Date    string `json:"date,omitempty"`
		Content string `json:"content"`
	}
	mcp.AddTool(srv, &mcp.Tool{Name: "get-adr", Description: "Get a single ADR by URI"},
		func(ctx context.Context, req *mcp.CallToolRequest, input GetADRInput) (*mcp.CallToolResult, GetADROutput, error) {
			a, ok := scn.Get(input.URI)
			if !ok {
				return nil, GetADROutput{}, fmt.Errorf("ADR not found: %s", input.URI)
			}
			return nil, GetADROutput{URI: a.URI, ID: a.ID, Title: a.Title, Status: a.Status, Date: a.Date, Content: a.Content}, nil
		},
	)

	type OptimizeInput struct {
		Source   string `json:"source,omitempty"`
		Status   string `json:"status,omitempty"`
		Optimize string `json:"optimize,omitempty"`
	}
	type OptimizeOutput struct {
		Content string `json:"content"`
		Sources int    `json:"sources"`
	}
	mcp.AddTool(srv, &mcp.Tool{Name: "optimize-adrs", Description: "Get consolidated ADRs"},
		func(ctx context.Context, req *mcp.CallToolRequest, input OptimizeInput) (*mcp.CallToolResult, OptimizeOutput, error) {
			filtered := filterByStatus(filterBySource(scn.List(), input.Source), input.Status)
			content := optimizeContent(ctx, opt, optimizeDefault, input.Optimize, toOptimizerInputs(filtered))
			return nil, OptimizeOutput{Content: content, Sources: len(filtered)}, nil
		},
	)
}

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

func filterBySource(adrs []scanner.ADR, source string) []scanner.ADR {
	if source == "" {
		return adrs
	}
	var out []scanner.ADR
	for _, a := range adrs {
		if a.Source == source {
			out = append(out, a)
		}
	}
	return out
}

func filterByStatus(adrs []scanner.ADR, status string) []scanner.ADR {
	if status == "" {
		return adrs
	}
	var out []scanner.ADR
	for _, a := range adrs {
		if strings.EqualFold(a.Status, status) {
			out = append(out, a)
		}
	}
	return out
}
