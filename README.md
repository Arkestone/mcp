# Arkestone MCP Servers

[![CI](https://github.com/Arkestone/mcp/actions/workflows/ci.yml/badge.svg)](https://github.com/Arkestone/mcp/actions/workflows/ci.yml)
[![codecov](https://codecov.io/gh/Arkestone/mcp/graph/badge.svg)](https://codecov.io/gh/Arkestone/mcp)
[![Go Report Card](https://goreportcard.com/badge/github.com/Arkestone/mcp)](https://goreportcard.com/report/github.com/Arkestone/mcp)
[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](LICENSE)

## Overview

A monorepo hosting [Model Context Protocol](https://modelcontextprotocol.io) (MCP) servers for [GitHub Copilot](https://github.com/features/copilot) customization. Each server implements the MCP specification to dynamically serve different types of Copilot configuration — custom instructions, skills, and more — from local directories and GitHub repositories.

## Servers

| Server | Description | Docs |
|--------|-------------|------|
| [mcp-instructions](./instructions/) | Serves Copilot custom instructions | [README](./instructions/README.md) |
| [mcp-skills](./skills/) | Serves Copilot skills | [README](./skills/README.md) |

## Shared Components

- **`pkg/optimizer`** — OpenAI-compatible LLM optimization layer shared by all servers. Optionally merges, deduplicates, and consolidates content from multiple sources via any OpenAI-compatible endpoint.
- **`pkg/config`** — Unified configuration loading (YAML → env vars → CLI flags).
- **`pkg/github`** — GitHub Contents API client with proxy and header pass-through.
- **`pkg/httputil`** — Shared HTTP infrastructure: proxy support, custom TLS/CA, header propagation.
- **`pkg/syncer`** — Background periodic sync for remote repositories.
- **`pkg/server`** — MCP server bootstrap helpers.

## Quick Start

### mcp-instructions

```bash
make build-instructions
./bin/mcp-instructions -dirs /path/to/repo
```

### mcp-skills

```bash
make build-skills
./bin/mcp-skills -repos github/awesome-copilot
```

See each server's README for full configuration and usage details.

## Architecture

```
.
├── instructions/           # mcp-instructions server
│   ├── cmd/mcp-instructions/
│   └── internal/loader/
├── skills/                 # mcp-skills server
│   ├── cmd/mcp-skills/
│   └── internal/scanner/
├── pkg/
│   ├── config/             # shared configuration loading
│   ├── github/             # GitHub Contents API client
│   ├── httputil/           # proxy, TLS, header propagation
│   ├── optimizer/          # shared LLM optimization layer
│   ├── server/             # MCP server helpers
│   └── syncer/             # background sync
├── docs/
│   └── network.md          # network / proxy / firewall guide
├── Makefile
├── go.mod
└── go.sum
```

Each server follows the same layered design:

1. **Config** — YAML → environment variables → CLI flags (each layer overrides the previous)
2. **Loader / Scanner** — discovers content from local directories and GitHub repositories
3. **Optimizer** — optional LLM-based consolidation via `pkg/optimizer`
4. **MCP Server** — exposes content as Resources, Prompts, and Tools over stdio or Streamable HTTP

## Development

```bash
make build              # build all servers
make test               # run unit tests
make test-integration   # run integration tests
make docker             # build Docker images for all servers
make lint               # run golangci-lint
make cover              # generate coverage report
```

## Configuration

Each server supports configuration via YAML files, environment variables, and CLI flags. See the individual server READMEs for details:

- [mcp-instructions configuration](./instructions/README.md#configuration)
- [mcp-skills configuration](./skills/README.md#configuration)

### GitHub Authentication (Optional)

A GitHub token is **optional**. Public repositories work without authentication. For private repositories, provide a token via any of these methods (highest priority first):

| Method | Example |
|--------|---------|
| CLI flag | `--github-token ghp_xxx` |
| Prefixed env var | `INSTRUCTIONS_GITHUB_TOKEN=ghp_xxx` / `SKILLS_GITHUB_TOKEN=ghp_xxx` |
| Global env var | `GITHUB_TOKEN=ghp_xxx` |
| YAML config | `github_token: ghp_xxx` |

When a repository returns HTTP 401/403/404 without a token configured, the error message will hint that authentication may be required.

## Network & Proxy

All servers work on-premise, in private/public cloud, with direct internet or through HTTP/HTTPS proxies. See the **[Network & Proxy Guide](docs/network.md)** for:

- Firewall rules to open
- Proxy and custom CA certificate configuration
- HTTP header pass-through from incoming requests to outbound calls
- Deployment scenarios (direct, proxy, air-gapped, Kubernetes)

## Contributing

Contributions are welcome! Please see [CONTRIBUTING.md](CONTRIBUTING.md) for guidelines.

## License

This project is licensed under the MIT License — see the [LICENSE](LICENSE) file for details.
