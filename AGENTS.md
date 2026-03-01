# AGENTS.md — AI Coding Assistant Guide

This file provides structured guidance for AI coding assistants (GitHub Copilot, Cursor, Claude, etc.) working in this repository.

## Repository Overview

This is a Go monorepo hosting five MCP (Model Context Protocol) servers:

| Server | Directory | Binary | Default Port |
|--------|-----------|--------|-------------|
| Instructions | `servers/mcp-instructions/` | `mcp-instructions` | `:8080` |
| Skills | `servers/mcp-skills/` | `mcp-skills` | `:8081` |
| Prompts | `servers/mcp-prompts/` | `mcp-prompts` | `:8082` |
| ADR | `servers/mcp-adr/` | `mcp-adr` | `:8083` |
| Memory | `servers/mcp-memory/` | `mcp-memory` | `:8084` |

Shared packages live in `pkg/`. Never put server-specific code in `pkg/`.

## Build Commands

```bash
# Build all servers (requires Go 1.24+)
export PATH=$PATH:~/go-install/go/bin   # if go not in PATH
make build

# Build individual servers
go build -buildvcs=false -o mcp-instructions ./servers/mcp-instructions/cmd/mcp-instructions
go build -buildvcs=false -o mcp-skills       ./servers/mcp-skills/cmd/mcp-skills
go build -buildvcs=false -o mcp-prompts      ./servers/mcp-prompts/cmd/mcp-prompts
go build -buildvcs=false -o mcp-adr          ./servers/mcp-adr/cmd/mcp-adr
go build -buildvcs=false -o mcp-memory       ./servers/mcp-memory/cmd/mcp-memory

# Run all tests
go test -buildvcs=false ./...

# Format check
gofmt -l .
```

## Directory Layout

```
mcp/
├── servers/
│   ├── mcp-instructions/          # Instructions MCP server      (:8080)
│   │   ├── cmd/mcp-instructions/  # Entry point (main.go)
│   │   ├── internal/
│   │   │   ├── config/            # Config loading (YAML + env + flags)
│   │   │   ├── loader/            # Instruction file discovery + GitHub cache
│   │   │   └── optimizer/         # Instruction-specific LLM optimization
│   │   ├── Dockerfile
│   │   ├── README.md
│   │   └── CHANGELOG.md
│   ├── mcp-skills/                # Skills MCP server            (:8081)
│   │   ├── cmd/mcp-skills/        # Entry point (main.go)
│   │   ├── internal/
│   │   │   └── scanner/           # SKILL.md frontmatter scanner
│   │   ├── Dockerfile
│   │   ├── README.md
│   │   └── CHANGELOG.md
│   ├── mcp-prompts/               # Prompts MCP server           (:8082)
│   │   ├── cmd/mcp-prompts/
│   │   ├── internal/
│   │   ├── Dockerfile
│   │   ├── README.md
│   │   └── CHANGELOG.md
│   ├── mcp-adr/                   # ADR MCP server               (:8083)
│   │   ├── cmd/mcp-adr/
│   │   ├── internal/
│   │   ├── Dockerfile
│   │   ├── README.md
│   │   └── CHANGELOG.md
│   └── mcp-memory/                # Memory MCP server            (:8084)
│       ├── cmd/mcp-memory/
│       ├── internal/
│       │   └── store/             # On-disk memory store
│       ├── Dockerfile
│       ├── README.md
│       └── CHANGELOG.md
├── pkg/
│   ├── config/                    # Shared config primitives
│   ├── github/                    # GitHub API client (raw file + dir listing)
│   ├── httputil/                  # HTTP helpers (retry, timeouts)
│   ├── optimizer/                 # OpenAI-compatible LLM client + LLMConfig
│   ├── server/                    # HTTP server wrapper (RunHTTP, /healthz)
│   ├── syncer/                    # Background sync goroutine
│   └── testutil/                  # Shared test helpers (not production code)
├── examples/                      # Client configuration examples
├── docs/                          # Architecture and reference docs
└── .devcontainer/                 # Dev Container for Codespaces
```

## Conventions

- **Config precedence**: YAML file < environment variables < CLI flags
- **Instruction URIs**: `instructions://{source}/{name}` where source is a directory basename or `owner/repo`
- **`optimize` parameter**: accepts `"true"`/`"false"` to override the global `llm.enabled` per-request
- **HTTP transport**: `mcp.NewStreamableHTTPHandler` — stateless, no session management
- **GitHub API**: `application/vnd.github.raw+json` for file content; `application/vnd.github+json` for directory listings
- **Determinism**: always sort any map-derived slices before returning them from MCP tools
- **Test helpers**: add to `pkg/testutil/` not to production packages

## Adding a New MCP Server

1. Create `servers/<server-name>/cmd/<server-name>/main.go` (copy structure from `servers/mcp-skills/`)
2. Add internal packages under `servers/<server-name>/internal/`
3. Add `servers/<server-name>/Dockerfile` (copy and adjust from an existing one)
4. Add `servers/<server-name>/README.md` and `servers/<server-name>/CHANGELOG.md`
5. Add build targets to `Makefile`, goreleaser config, and CI workflow
6. Add a client config example to `examples/`

## Important Notes

- **`-buildvcs=false`**: required in this environment (no VCS info available at build time)
- **CGO_ENABLED=0**: always disabled; do not use `-race` flag
- **`pkg/testutil`**: has no `_test.go` files (it's a helper for other packages' tests); `[no test files]` is expected
- **`go.mod` module path**: `github.com/Arkestone/mcp` (capital A in Arkestone)
