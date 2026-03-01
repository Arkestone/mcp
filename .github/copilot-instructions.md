# Build & run

This is a Go monorepo containing multiple MCP servers and shared packages.

```bash
# Instructions MCP server
go build -o mcp-instructions ./instructions/cmd/mcp-instructions

# Skills MCP server
go build -o mcp-skills ./skills/cmd/mcp-skills

# Run all tests
go test ./...
```

# Architecture

This monorepo hosts multiple MCP servers built with the official Go SDK (`github.com/modelcontextprotocol/go-sdk/mcp`). Each server serves different Copilot data via the Model Context Protocol.

**Monorepo layout:**
- `instructions/` — Instructions MCP server (serves Copilot custom instruction files)
- `skills/` — Skills MCP server (`mcp-skills`, serves Copilot skills)
- `pkg/` — Shared packages used across servers
  - `pkg/optimizer/` — shared LLM optimization (OpenAI-compatible client, LLMConfig)
- `cmd/mcp-instructions/` — entry point for the instructions server
- `cmd/mcp-skills/` — entry point for the skills server

## Instructions server

- `instructions/internal/config/` — YAML + env + CLI flag config loading (layered, each overrides the previous)
- `instructions/internal/loader/` — on-demand instruction file discovery from local dirs and cached GitHub repos with periodic background sync
- `instructions/internal/optimizer/` — instruction-specific optimization wiring

**Key data flow:** MCP client request → server handler → `loader.List()` or `loader.Get(uri)` reads from disk on-demand → optionally passed through `optimizer.Optimize()` → returned as MCP resource/prompt/tool result.

Local directories are always read live from disk. Remote GitHub repos are cached locally under `cache.dir` and synced in a background goroutine at `cache.sync_interval`.

## Skills server

- `skills/` — skills scanner with SKILL.md frontmatter parsing
- Skills-specific MCP resources, prompts, and tools

# Conventions

- Config precedence: YAML file < environment variables < CLI flags
- Instruction URIs follow the pattern `instructions://{source}/{name}` where source is a directory basename or `owner/repo`
- The `optimize` argument on prompts and tools accepts `true`/`false` to override the global `llm.enabled` default per-request
- HTTP transport uses `mcp.NewStreamableHTTPHandler` (stateless, no session management)
- All GitHub API calls use `application/vnd.github.raw+json` for file content and `application/vnd.github+json` for directory listings
- Shared code lives in `pkg/`; server-specific code stays under each server's directory
