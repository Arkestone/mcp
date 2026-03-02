---
title: mcp-instructions
sidebar_label: mcp-instructions
---


<!-- install-badges -->
| Transport | VS Code | VS Code Insiders |
|-----------|---------|-----------------|
| stdio | [![Install in VS Code](https://img.shields.io/badge/VS_Code-Install-0098FF?logo=visualstudiocode&logoColor=white)](https://vscode.dev/redirect/mcp/install?name=mcp-instructions&config=%7B%22type%22%3A%20%22stdio%22%2C%20%22command%22%3A%20%22mcp-instructions%22%2C%20%22args%22%3A%20%5B%22--dirs%22%2C%20%22%24%7BworkspaceFolder%7D%22%5D%7D) | [![Install in VS Code Insiders](https://img.shields.io/badge/VS_Code_Insiders-Install-24bfa5?logo=visualstudiocode&logoColor=white)](https://insiders.vscode.dev/redirect/mcp/install?name=mcp-instructions&config=%7B%22type%22%3A%20%22stdio%22%2C%20%22command%22%3A%20%22mcp-instructions%22%2C%20%22args%22%3A%20%5B%22--dirs%22%2C%20%22%24%7BworkspaceFolder%7D%22%5D%7D) |
| HTTP  | [![Install in VS Code](https://img.shields.io/badge/VS_Code-Install-0098FF?logo=visualstudiocode&logoColor=white)](https://vscode.dev/redirect/mcp/install?name=mcp-instructions&config=%7B%22type%22%3A%20%22http%22%2C%20%22url%22%3A%20%22http%3A%2F%2Flocalhost%3A8080%2Fmcp%22%7D) | [![Install in VS Code Insiders](https://img.shields.io/badge/VS_Code_Insiders-Install-24bfa5?logo=visualstudiocode&logoColor=white)](https://insiders.vscode.dev/redirect/mcp/install?name=mcp-instructions&config=%7B%22type%22%3A%20%22http%22%2C%20%22url%22%3A%20%22http%3A%2F%2Flocalhost%3A8080%2Fmcp%22%7D) |
<!-- /install-badges -->

An MCP server that dynamically serves [GitHub Copilot custom instructions](https://docs.github.com/en/copilot/customize-github-copilot/adding-custom-instructions-for-github-copilot) from local directories and GitHub repositories.

It discovers `.github/copilot-instructions.md` and `.github/instructions/**/*.instructions.md` files from configured sources and exposes them via the [Model Context Protocol](https://modelcontextprotocol.io).

## Features

- **On-demand loading** — local directories are read live; no eager loading
- **GitHub repo caching** — remote repos are cached locally and synced periodically in the background
- **LLM optimization** — optionally merge/deduplicate instructions from multiple sources via an OpenAI-compatible endpoint
- **Dual transport** — supports both stdio (for Copilot CLI / local tools) and Streamable HTTP (for remote deployments)
- **MCP primitives** — instructions are exposed as Resources, Prompts, and Tools

## Installation

```bash
# go install (requires Go 1.24+)
go install github.com/Arkestone/mcp/servers/mcp-instructions/cmd/mcp-instructions@latest

# Pinned version
go install github.com/Arkestone/mcp/servers/mcp-instructions/cmd/mcp-instructions@v0.0.1

# Docker
docker pull ghcr.io/arkestone/mcp-instructions:latest

# Pre-built binary — https://github.com/Arkestone/mcp/releases/latest
```

## Getting Started

```bash
# Serve instructions from a local directory (stdio)
mcp-instructions -dirs /path/to/repo

# Serve from a GitHub repo with HTTP transport
export GITHUB_TOKEN=ghp_...
mcp-instructions -repos github/awesome-copilot -transport http -addr :8080

# Build from source
make build-instructions   # → ./bin/mcp-instructions
```

## Configuration

Configuration is loaded in layers (each overrides the previous):

1. **YAML file** — `config.yaml` in the working directory, or specify with `-config path/to/config.yaml`
2. **Environment variables**
3. **CLI flags**

See [`config.example.yaml`](https://github.com/Arkestone/mcp/blob/main/servers/mcp-instructions/config.example.yaml) for all options.

### Environment Variables

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

### CLI Flags

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
| `get-instructions` | `source` (optional), `optimize` (optional) | Get instructions as a prompt, optionally filtered and optimized |

### Tools

| Name | Arguments | Description |
|---|---|---|
| `refresh` | — | Force-sync all remote repo caches |
| `list-instructions` | — | List all available instruction files |
| `optimize-instructions` | `source` (optional), `optimize` (optional) | Get consolidated instructions with optional LLM optimization |

## Instruction File Format

The server discovers two kinds of instruction files from each configured source:

| Path | Scope |
|---|---|
| `.github/copilot-instructions.md` | Repository-wide instructions applied to every conversation |
| `.github/instructions/*.instructions.md` | Path-specific instructions (front-matter `applyTo` glob controls scope) |

For details on authoring instruction files, see [Customizing Copilot with custom instructions](https://docs.github.com/en/copilot/customize-github-copilot/adding-custom-instructions-for-github-copilot).

## Docker

```bash
# From the repo root
make docker-instructions

# stdio mode
docker run -i ghcr.io/arkestone/mcp-instructions:latest -dirs /data

# HTTP mode
docker run -p 8080:8080 \
  -e INSTRUCTIONS_REPOS=github/awesome-copilot \
  -e GITHUB_TOKEN=ghp_... \
  -e INSTRUCTIONS_TRANSPORT=http \
  ghcr.io/arkestone/mcp-instructions:latest
```

## MCP Client Configuration

### VS Code / GitHub Copilot

`.vscode/mcp.json`:

```json
{
  "servers": {
    "instructions": {
      "command": "mcp-instructions",
      "args": ["-dirs", "${workspaceFolder}"]
    }
  }
}
```

### Claude Desktop

**macOS**: `~/Library/Application Support/Claude/claude_desktop_config.json`
**Windows**: `%APPDATA%\Claude\claude_desktop_config.json`

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

### Cursor

`.cursor/mcp.json` (project) or `~/.cursor/mcp.json` (global):

```json
{
  "mcpServers": {
    "instructions": {
      "command": "mcp-instructions",
      "args": ["-dirs", "."]
    }
  }
}
```

### Windsurf

`~/.codeium/windsurf/mcp_config.json`:

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

### Claude Code

`.mcp.json` (project) or `~/.mcp.json` (global):

```json
{
  "mcpServers": {
    "instructions": {
      "command": "mcp-instructions",
      "args": ["-dirs", "."]
    }
  }
}
```

### Remote (HTTP)

For shared team deployments, start the server with `-transport http` and connect clients over HTTP:

```json
{
  "mcpServers": {
    "instructions": {
      "type": "http",
      "url": "http://localhost:8080/mcp"
    }
  }
}
```
