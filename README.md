# mcp-instructions

[![CI](https://github.com/YOUR_ORG/mcp-instructions/actions/workflows/ci.yml/badge.svg)](https://github.com/YOUR_ORG/mcp-instructions/actions/workflows/ci.yml)
[![Go Report Card](https://goreportcard.com/badge/github.com/YOUR_ORG/mcp-instructions)](https://goreportcard.com/report/github.com/YOUR_ORG/mcp-instructions)
[![Go Reference](https://pkg.go.dev/badge/github.com/YOUR_ORG/mcp-instructions.svg)](https://pkg.go.dev/github.com/YOUR_ORG/mcp-instructions)
[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](LICENSE)

## Overview

An MCP server that dynamically serves [GitHub Copilot custom instructions](https://docs.github.com/en/copilot/how-tos/copilot-cli/customize-copilot/add-custom-instructions) from local directories and GitHub repositories.

It discovers `.github/copilot-instructions.md` and `.github/instructions/**/*.instructions.md` files from configured sources and exposes them via the [Model Context Protocol](https://modelcontextprotocol.io).

## Features

- **On-demand loading** — local directories are read live; no eager loading
- **GitHub repo caching** — remote repos are cached locally and synced periodically in the background
- **LLM optimization** — optionally merge/deduplicate instructions from multiple sources via an OpenAI-compatible endpoint
- **Dual transport** — supports both stdio (for Copilot CLI / local tools) and Streamable HTTP (for remote deployments)
- **MCP primitives** — instructions are exposed as Resources, Prompts, and Tools

## Architecture

The server is organized into four layers:

1. **Config** — loads settings from YAML → environment variables → CLI flags (each layer overrides the previous).
2. **Loader** — discovers instruction files. Local directories are read on every request; GitHub repos are cached locally and synced in the background on a configurable interval.
3. **Optimizer** — optionally consolidates multiple instruction files into a single output via an OpenAI-compatible LLM endpoint, deduplicating and merging content.
4. **MCP Server** — exposes the loaded (and optionally optimized) instructions as MCP Resources, Prompts, and Tools over stdio or Streamable HTTP.

## Getting Started

```bash
# Build
go build -o mcp-instructions ./cmd/mcp-instructions

# Serve instructions from a local directory (stdio)
mcp-instructions -dirs /path/to/repo

# Serve from a GitHub repo with HTTP transport
export GITHUB_TOKEN=ghp_...
mcp-instructions -repos github/awesome-copilot -transport http -addr :8080
```

## Configuration

Configuration is loaded in layers (each overrides the previous):

1. **YAML file** — `config.yaml` in the working directory, or specify with `-config path/to/config.yaml`
2. **Environment variables**
3. **CLI flags**

See [`config.example.yaml`](config.example.yaml) for all options.

### Environment variables

| Variable | Description |
|---|---|
| `INSTRUCTIONS_CONFIG` | Path to YAML config file |
| `INSTRUCTIONS_DIRS` | Comma-separated local directories |
| `INSTRUCTIONS_REPOS` | Comma-separated GitHub repos (`owner/repo` or `owner/repo@ref`) |
| `INSTRUCTIONS_TRANSPORT` | `stdio` (default) or `http` |
| `INSTRUCTIONS_ADDR` | HTTP listen address (default `:8080`) |
| `INSTRUCTIONS_CACHE_DIR` | Local cache directory (default `~/.cache/mcp-instructions`) |
| `INSTRUCTIONS_SYNC_INTERVAL` | Sync interval for remote repos (default `5m`) |
| `GITHUB_TOKEN` | GitHub API token for private repos |
| `LLM_ENDPOINT` | OpenAI-compatible API endpoint |
| `LLM_MODEL` | Model name (e.g. `gpt-4o-mini`) |
| `LLM_API_KEY` | LLM API key |
| `LLM_ENABLED` | `true` to enable LLM optimization by default |

### CLI flags

```
-config          Path to YAML config file
-dirs            Comma-separated local directories
-repos           Comma-separated GitHub repos (owner/repo[@ref])
-transport       Transport: stdio (default) or http
-addr            HTTP listen address (default :8080)
-cache-dir       Local cache directory
-sync-interval   Sync interval (e.g. 5m, 1h)
-llm-endpoint    OpenAI-compatible endpoint URL
-llm-model       LLM model name
```

## MCP Primitives

### Resources

| URI | Description |
|---|---|
| `instructions://{source}/{name}` | Individual instruction file content |
| `instructions://optimized` | All instructions merged via LLM (or concatenated) |
| `instructions://index` | List of all available instruction URIs |

### Prompts

| Name | Arguments | Description |
|---|---|---|
| `get-instructions` | `source` (optional), `optimize` (optional: `true`/`false`) | Get instructions as a prompt, optionally filtered and optimized |

### Tools

| Name | Arguments | Description |
|---|---|---|
| `refresh` | — | Force-sync all remote repo caches |
| `list-instructions` | — | List all available instruction files |
| `optimize-instructions` | `source` (optional), `optimize` (optional: `true`/`false`) | Get consolidated instructions with optional LLM optimization |

## Instruction File Format

The server discovers two kinds of instruction files from each configured source:

| Path | Scope |
|---|---|
| `.github/copilot-instructions.md` | Repository-wide instructions applied to every conversation |
| `.github/instructions/*.instructions.md` | Path-specific instructions (the file's front-matter `applyTo` glob controls which files they apply to) |

For details on authoring instruction files, see [Customizing Copilot with custom instructions](https://docs.github.com/en/copilot/how-tos/copilot-cli/customize-copilot/add-custom-instructions).

## Sample Repository

The [`github/awesome-copilot`](https://github.com/github/awesome-copilot) repository is a curated collection of Copilot instructions. Point `mcp-instructions` at it to get started quickly:

```bash
mcp-instructions -repos github/awesome-copilot
```

## Docker

```bash
docker build -t mcp-instructions .

# stdio mode (pipe to MCP client)
docker run -i mcp-instructions -dirs /data

# HTTP mode
docker run -p 8080:8080 \
  -e INSTRUCTIONS_REPOS=github/awesome-copilot \
  -e GITHUB_TOKEN=ghp_... \
  -e INSTRUCTIONS_TRANSPORT=http \
  mcp-instructions
```

## MCP Client Configuration

### Copilot CLI / VS Code

Add to your MCP client settings:

```json
{
  "mcpServers": {
    "instructions": {
      "command": "mcp-instructions",
      "args": ["-dirs", "/path/to/repo"]
    }
  }
}
```

### Remote (HTTP)

```json
{
  "mcpServers": {
    "instructions": {
      "url": "http://localhost:8080"
    }
  }
}
```

## Development

```bash
go build ./...          # build all packages
go test ./...           # run tests
go vet ./...            # lint
```

If a `Makefile` is present, you can also use `make build`, `make test`, and `make lint`.

See [CONTRIBUTING.md](CONTRIBUTING.md) for contribution guidelines (if available).

## License

This project is licensed under the MIT License — see the [LICENSE](LICENSE) file for details.

## Contributing

Contributions are welcome! Please open an issue or submit a pull request. If a [CONTRIBUTING.md](CONTRIBUTING.md) file is present, review it before submitting.
