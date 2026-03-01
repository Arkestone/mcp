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
	"github.com/Arkestone/mcp/skills/internal/scanner"
	"github.com/modelcontextprotocol/go-sdk/mcp"
)

var version = "dev"  // set at build time via -ldflags
var commit = "none"  // set at build time via -ldflags
var date = "unknown" // set at build time via -ldflags

func main() {
	log.Printf("mcp-skills version=%s commit=%s date=%s", version, commit, date)
	cfg := config.Load(config.Options{
		EnvPrefix:        "SKILLS",
		DefaultAddr:      ":8081",
		DefaultCacheName: "mcp-skills",
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
		&mcp.Implementation{Name: "mcp-skills", Version: "0.1.0"},
		nil,
	)

	registerResources(srv, scn, opt, cfg.LLM.Enabled)
	registerPrompts(srv, scn, opt, cfg.LLM.Enabled)
	registerTools(srv, scn, opt, cfg.LLM.Enabled)

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

func toOptimizerInputs(skills []scanner.Skill) []optimizer.ContentInput {
	inputs := make([]optimizer.ContentInput, len(skills))
	for i, s := range skills {
		inputs[i] = optimizer.ContentInput{
			Source: s.Source, Path: s.Path, Content: s.Content,
		}
	}
	return inputs
}

func registerResources(srv *mcp.Server, scn *scanner.Scanner, opt *optimizer.Optimizer, optimizeDefault bool) {
	srv.AddResourceTemplate(
		&mcp.ResourceTemplate{
			URITemplate: "skills://{source}/{name}",
			Name:        "Skill",
			Description: "Individual skill content with metadata and references",
			MIMEType:    "text/markdown",
		},
		func(ctx context.Context, req *mcp.ReadResourceRequest) (*mcp.ReadResourceResult, error) {
			skill, ok := scn.Get(req.Params.URI)
			if !ok {
				return nil, fmt.Errorf("skill not found: %s", req.Params.URI)
			}
			return &mcp.ReadResourceResult{
				Contents: []*mcp.ResourceContents{{
					URI: req.Params.URI, MIMEType: "text/markdown",
					Text: renderSkill(skill),
				}},
			}, nil
		},
	)

	srv.AddResource(
		&mcp.Resource{
			URI: "skills://optimized", Name: "Optimized Skills",
			Description: "All skills merged via LLM (or concatenated)", MIMEType: "text/markdown",
		},
		func(ctx context.Context, req *mcp.ReadResourceRequest) (*mcp.ReadResourceResult, error) {
			content := optimizeContent(ctx, opt, optimizeDefault, "", toOptimizerInputs(scn.List()))
			return &mcp.ReadResourceResult{
				Contents: []*mcp.ResourceContents{{
					URI: "skills://optimized", MIMEType: "text/markdown", Text: content,
				}},
			}, nil
		},
	)

	srv.AddResource(
		&mcp.Resource{
			URI: "skills://index", Name: "Skills Index",
			Description: "List of all available skills", MIMEType: "text/plain",
		},
		func(ctx context.Context, req *mcp.ReadResourceRequest) (*mcp.ReadResourceResult, error) {
			var sb strings.Builder
			for _, s := range scn.List() {
				fmt.Fprintf(&sb, "%s  [%s] %s — %s\n", s.URI, s.Source, s.Name, s.Description)
			}
			return &mcp.ReadResourceResult{
				Contents: []*mcp.ResourceContents{{
					URI: "skills://index", MIMEType: "text/plain", Text: sb.String(),
				}},
			}, nil
		},
	)
}

func registerPrompts(srv *mcp.Server, scn *scanner.Scanner, opt *optimizer.Optimizer, optimizeDefault bool) {
	srv.AddPrompt(
		&mcp.Prompt{
			Name:        "get-skills",
			Description: "Get all skills, optionally optimized via LLM",
			Arguments: []*mcp.PromptArgument{
				{Name: "source", Description: "Filter by source (optional)", Required: false},
				{Name: "optimize", Description: "Override LLM optimization (true/false)", Required: false},
			},
		},
		func(ctx context.Context, req *mcp.GetPromptRequest) (*mcp.GetPromptResult, error) {
			filtered := filterBySource(scn.List(), req.Params.Arguments["source"])
			if len(filtered) == 0 {
				return promptResult("Skills for AI assistants", "No skills found."), nil
			}
			content := optimizeContent(ctx, opt, optimizeDefault,
				req.Params.Arguments["optimize"], toOptimizerInputs(filtered))
			return promptResult("Skills for AI assistants", content), nil
		},
	)

	srv.AddPrompt(
		&mcp.Prompt{
			Name:        "get-skill",
			Description: "Get a single skill's full content including references",
			Arguments: []*mcp.PromptArgument{
				{Name: "name", Description: "The skill name", Required: true},
			},
		},
		func(ctx context.Context, req *mcp.GetPromptRequest) (*mcp.GetPromptResult, error) {
			name := req.Params.Arguments["name"]
			for _, s := range scn.List() {
				if s.Name == name {
					return promptResult(fmt.Sprintf("Skill: %s", s.Name), renderSkill(s)), nil
				}
			}
			return promptResult("Skill not found", fmt.Sprintf("Skill %q not found.", name)), nil
		},
	)
}

func registerTools(srv *mcp.Server, scn *scanner.Scanner, opt *optimizer.Optimizer, optimizeDefault bool) {
	type RefreshInput struct {
		Source string `json:"source,omitempty" jsonschema:"Optional source filter. When set only the matching source is returned in the result. All sources are always refreshed."`
	}
	type RefreshOutput struct {
		Message string   `json:"message"`
		Count   int      `json:"count"`
		Sources []string `json:"sources"`
	}
	mcp.AddTool(srv, &mcp.Tool{
		Name: "refresh-skills",
		Description: "Force an immediate re-sync of all skill sources. " +
			"This fetches the latest SKILL.md files and references from every configured GitHub repository " +
			"into the local cache and re-reads local directories. Call this tool when skills may have changed " +
			"on disk or in remote repositories and you need the most up-to-date content. " +
			"Returns the total count of skills found after the refresh and the list of distinct sources.",
	},
		func(ctx context.Context, req *mcp.CallToolRequest, input RefreshInput) (*mcp.CallToolResult, RefreshOutput, error) {
			scn.ForceSync()
			skills := scn.List()
			sourceSet := map[string]struct{}{}
			for _, s := range skills {
				sourceSet[s.Source] = struct{}{}
			}
			sources := make([]string, 0, len(sourceSet))
			for s := range sourceSet {
				sources = append(sources, s)
			}
			sort.Strings(sources)
			return nil, RefreshOutput{
				Message: "All skill sources refreshed successfully",
				Count:   len(skills),
				Sources: sources,
			}, nil
		},
	)

	type ListInput struct{}
	type ListEntry struct {
		URI         string `json:"uri"`
		Name        string `json:"name"`
		Description string `json:"description"`
		Source      string `json:"source"`
	}
	type ListOutput struct {
		Entries []ListEntry `json:"entries"`
	}
	mcp.AddTool(srv, &mcp.Tool{Name: "list-skills", Description: "List all available skills"},
		func(ctx context.Context, req *mcp.CallToolRequest, input ListInput) (*mcp.CallToolResult, ListOutput, error) {
			skills := scn.List()
			entries := make([]ListEntry, len(skills))
			for i, s := range skills {
				entries[i] = ListEntry{URI: s.URI, Name: s.Name, Description: s.Description, Source: s.Source}
			}
			return nil, ListOutput{Entries: entries}, nil
		},
	)

	type GetSkillInput struct {
		Name string `json:"name"`
	}
	type GetSkillOutput struct {
		Content    string              `json:"content"`
		References []scanner.Reference `json:"references,omitempty"`
	}
	mcp.AddTool(srv, &mcp.Tool{Name: "get-skill", Description: "Get a skill by name"},
		func(ctx context.Context, req *mcp.CallToolRequest, input GetSkillInput) (*mcp.CallToolResult, GetSkillOutput, error) {
			for _, s := range scn.List() {
				if s.Name == input.Name {
					return nil, GetSkillOutput{Content: s.Content, References: s.References}, nil
				}
			}
			return nil, GetSkillOutput{}, fmt.Errorf("skill %q not found", input.Name)
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
	mcp.AddTool(srv, &mcp.Tool{Name: "optimize-skills", Description: "Get consolidated skills"},
		func(ctx context.Context, req *mcp.CallToolRequest, input OptimizeInput) (*mcp.CallToolResult, OptimizeOutput, error) {
			filtered := filterBySource(scn.List(), input.Source)
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

func renderSkill(s scanner.Skill) string {
	var sb strings.Builder
	fmt.Fprintf(&sb, "---\nname: %s\ndescription: %s\n---\n\n", s.Name, s.Description)
	sb.WriteString(s.Content)
	if len(s.References) > 0 {
		sb.WriteString("\n\n## References\n\n")
		for _, ref := range s.References {
			fmt.Fprintf(&sb, "### %s\n\n%s\n\n", ref.Name, ref.Content)
		}
	}
	return sb.String()
}

func filterBySource(skills []scanner.Skill, source string) []scanner.Skill {
	if source == "" {
		return skills
	}
	var out []scanner.Skill
	for _, s := range skills {
		if s.Source == source {
			out = append(out, s)
		}
	}
	return out
}
