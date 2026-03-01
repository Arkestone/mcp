# Arkestone MCP Servers

[![CI](https://github.com/Arkestone/mcp/actions/workflows/ci.yml/badge.svg)](https://github.com/Arkestone/mcp/actions/workflows/ci.yml)
[![codecov](https://codecov.io/gh/Arkestone/mcp/graph/badge.svg)](https://codecov.io/gh/Arkestone/mcp)
[![Go Report Card](https://goreportcard.com/badge/github.com/Arkestone/mcp)](https://goreportcard.com/report/github.com/Arkestone/mcp)
[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](LICENSE)

A suite of [Model Context Protocol](https://modelcontextprotocol.io) (MCP) servers for [GitHub Copilot](https://github.com/features/copilot) customization. Each server dynamically serves a different type of Copilot configuration — custom instructions, skills, prompts, ADRs, or persistent memory — from local directories and GitHub repositories.

## Available Servers

| Server | Description | Port | Docs |
|--------|-------------|------|------|
| [mcp-instructions](./servers/mcp-instructions/) | Serves Copilot custom instruction files (`.github/copilot-instructions.md`, `.github/instructions/**/*.instructions.md`) | `:8080` | [README](./servers/mcp-instructions/README.md) · [CHANGELOG](./servers/mcp-instructions/CHANGELOG.md) |
| [mcp-skills](./servers/mcp-skills/) | Serves Copilot skills (`SKILL.md`) with frontmatter metadata | `:8081` | [README](./servers/mcp-skills/README.md) · [CHANGELOG](./servers/mcp-skills/CHANGELOG.md) |
| [mcp-prompts](./servers/mcp-prompts/) | Serves VS Code Copilot prompt files (`.github/prompts/*.prompt.md`) and chat mode files | `:8082` | [README](./servers/mcp-prompts/README.md) · [CHANGELOG](./servers/mcp-prompts/CHANGELOG.md) |
| [mcp-graph](./servers/mcp-graph/) | Knowledge graph MCP server — store entities and relationships, query neighbors and shortest paths | `:8085` | [README](./servers/mcp-graph/README.md) · [CHANGELOG](./servers/mcp-graph/CHANGELOG.md) |
| [mcp-adr](./servers/mcp-adr/) | Serves Architecture Decision Records from `docs/adr/`, `docs/decisions/`, or `doc/adr/` | `:8083` | [README](./servers/mcp-adr/README.md) · [CHANGELOG](./servers/mcp-adr/CHANGELOG.md) |
| [mcp-memory](./servers/mcp-memory/) | Persistent memory store — remember, recall, and forget information across sessions | `:8084` | [README](./servers/mcp-memory/README.md) · [CHANGELOG](./servers/mcp-memory/CHANGELOG.md) |

## Quick Start

```bash
# Build all servers
make build

# Or build individually
make build-instructions   # → ./bin/mcp-instructions
make build-skills         # → ./bin/mcp-skills
make build-prompts        # → ./bin/mcp-prompts
make build-adr            # → ./bin/mcp-adr
make build-memory         # → ./bin/mcp-memory
```

Each server supports stdio (default) and HTTP transports:

```bash
# stdio — used by Copilot CLI and desktop clients
./bin/mcp-instructions -dirs /path/to/repo

# HTTP — used for remote/shared deployments
./bin/mcp-instructions -transport http -addr :8080 -repos github/awesome-copilot
```

See each server's README for the full configuration reference.

## MCP Client Configuration

### VS Code (`.vscode/mcp.json`)

```json
{
  "servers": {
    "instructions": { "command": "mcp-instructions", "args": ["-dirs", "/path/to/repo"] },
    "skills":       { "command": "mcp-skills",       "args": ["-dirs", "/path/to/skills"] },
    "prompts":      { "command": "mcp-prompts",      "args": ["-dirs", "/path/to/repo"] },
    "adrs":         { "command": "mcp-adr",          "args": ["-dirs", "/path/to/repo"] },
    "memory":       { "command": "mcp-memory" }
  }
}
```

### Claude Desktop (`~/Library/Application Support/Claude/claude_desktop_config.json`)

```json
{
  "mcpServers": {
    "instructions": {
      "command": "mcp-instructions",
      "args": ["-dirs", "/path/to/repo"]
    },
    "memory": {
      "command": "mcp-memory",
      "env": { "MEMORY_DIR": "~/.local/share/mcp-memory" }
    }
  }
}
```

### Docker Compose — all servers

```yaml
services:
  mcp-instructions:
    image: ghcr.io/arkestone/mcp-instructions:latest
    ports: ["8080:8080"]
    environment:
      INSTRUCTIONS_TRANSPORT: http
      INSTRUCTIONS_DIRS: /data
    volumes: ["./instructions:/data:ro"]

  mcp-skills:
    image: ghcr.io/arkestone/mcp-skills:latest
    ports: ["8081:8081"]
    environment:
      SKILLS_TRANSPORT: http
      SKILLS_DIRS: /data
    volumes: ["./skills:/data:ro"]

  mcp-prompts:
    image: ghcr.io/arkestone/mcp-prompts:latest
    ports: ["8082:8082"]
    environment:
      PROMPTS_TRANSPORT: http
      PROMPTS_SOURCES_DIRS: /data
    volumes: ["./prompts:/data:ro"]

  mcp-adr:
    image: ghcr.io/arkestone/mcp-adr:latest
    ports: ["8083:8083"]
    environment:
      ADR_TRANSPORT: http
      ADR_DIRS: /data
    volumes: ["./docs:/data:ro"]

  mcp-memory:
    image: ghcr.io/arkestone/mcp-memory:latest
    ports: ["8084:8084"]
    environment:
      MEMORY_TRANSPORT: http
      MEMORY_DIR: /data
    volumes: ["memory-data:/data"]

volumes:
  memory-data:
```

## Architecture

```
.
├── servers/
│   ├── mcp-instructions/   # custom instructions server  (:8080)
│   ├── mcp-skills/         # skills server               (:8081)
│   ├── mcp-prompts/        # prompt files server         (:8082)
│   └── mcp-graph/          # knowledge graph server      (:8085)
│   ├── mcp-adr/            # ADR server                  (:8083)
│   └── mcp-memory/         # persistent memory server    (:8084)
├── pkg/
│   ├── config/             # shared configuration loading (YAML → env → flags)
│   ├── github/             # GitHub Contents API client
│   ├── httputil/           # proxy, TLS, header propagation
│   ├── optimizer/          # shared LLM optimization layer (OpenAI-compatible)
│   ├── server/             # MCP server bootstrap helpers
│   └── syncer/             # background repo sync
├── docs/
│   └── network.md          # network / proxy / firewall guide
├── examples/               # client configuration examples
├── AGENTS.md               # AI coding assistant guide
├── Makefile
├── go.mod
└── go.sum
```

Each content server (instructions, skills, prompts, ADRs) follows the same layered design:

1. **Config** — YAML → environment variables → CLI flags (each layer overrides the previous)
2. **Loader / Scanner** — discovers content from local directories and GitHub repositories
3. **Optimizer** — optional LLM-based consolidation via `pkg/optimizer`
4. **MCP Server** — exposes content as Resources, Prompts, and Tools over stdio or Streamable HTTP

`mcp-memory` is a standalone persistent store and does not use the loader or optimizer layers.

## Shared Packages

| Package | Description |
|---------|-------------|
| `pkg/config` | Unified config loading: YAML → env vars → CLI flags |
| `pkg/github` | GitHub Contents API client with proxy and header pass-through |
| `pkg/httputil` | Proxy support, custom TLS/CA certificates, header propagation |
| `pkg/optimizer` | OpenAI-compatible LLM client for optional content consolidation |
| `pkg/server` | MCP server bootstrap and `/healthz` endpoint helpers |
| `pkg/syncer` | Background periodic sync for remote GitHub repositories |

## Development

```bash
make build              # build all servers into ./bin/
make test               # run unit tests
make test-integration   # run integration tests
make docker             # build Docker images for all servers
make lint               # run golangci-lint
make cover              # generate coverage report
```

## GitHub Authentication (Optional)

A GitHub token is **optional**. Public repositories work without authentication. For private repositories, provide a token (highest priority first):

| Method | Example |
|--------|---------|
| CLI flag | `-github-token ghp_xxx` |
| Prefixed env var | `INSTRUCTIONS_GITHUB_TOKEN=ghp_xxx` |
| Global env var | `GITHUB_TOKEN=ghp_xxx` |
| YAML config | `github_token: ghp_xxx` |

## Network & Proxy

All servers work on-premise, in private/public cloud, with direct internet or through HTTP/HTTPS proxies. See the **[Network & Proxy Guide](docs/network.md)** for firewall rules, proxy configuration, and custom CA certificates.

## Contributing

Contributions are welcome! Please see [CONTRIBUTING.md](CONTRIBUTING.md) for guidelines.

## License

This project is licensed under the MIT License — see the [LICENSE](LICENSE) file for details.
