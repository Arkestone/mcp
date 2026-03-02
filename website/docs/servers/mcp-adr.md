---
title: mcp-adr
sidebar_label: mcp-adr
---


<!-- install-badges -->
| Transport | VS Code | VS Code Insiders |
|-----------|---------|-----------------|
| stdio | [![Install in VS Code](https://img.shields.io/badge/VS_Code-Install-0098FF?logo=visualstudiocode&logoColor=white)](https://vscode.dev/redirect/mcp/install?name=mcp-adr&config=%7B%22type%22%3A%20%22stdio%22%2C%20%22command%22%3A%20%22mcp-adr%22%2C%20%22args%22%3A%20%5B%22--dirs%22%2C%20%22%24%7BworkspaceFolder%7D%22%5D%7D) | [![Install in VS Code Insiders](https://img.shields.io/badge/VS_Code_Insiders-Install-24bfa5?logo=visualstudiocode&logoColor=white)](https://insiders.vscode.dev/redirect/mcp/install?name=mcp-adr&config=%7B%22type%22%3A%20%22stdio%22%2C%20%22command%22%3A%20%22mcp-adr%22%2C%20%22args%22%3A%20%5B%22--dirs%22%2C%20%22%24%7BworkspaceFolder%7D%22%5D%7D) |
| HTTP  | [![Install in VS Code](https://img.shields.io/badge/VS_Code-Install-0098FF?logo=visualstudiocode&logoColor=white)](https://vscode.dev/redirect/mcp/install?name=mcp-adr&config=%7B%22type%22%3A%20%22http%22%2C%20%22url%22%3A%20%22http%3A%2F%2Flocalhost%3A8083%2Fmcp%22%7D) | [![Install in VS Code Insiders](https://img.shields.io/badge/VS_Code_Insiders-Install-24bfa5?logo=visualstudiocode&logoColor=white)](https://insiders.vscode.dev/redirect/mcp/install?name=mcp-adr&config=%7B%22type%22%3A%20%22http%22%2C%20%22url%22%3A%20%22http%3A%2F%2Flocalhost%3A8083%2Fmcp%22%7D) |
<!-- /install-badges -->

An MCP server that serves Architecture Decision Records (ADRs) from local directories and GitHub repositories. It scans `docs/adr/`, `docs/decisions/`, and `doc/adr/` within each configured source and exposes ADRs via the Model Context Protocol.

## Installation

```bash
# go install (requires Go 1.24+)
go install github.com/Arkestone/mcp/servers/mcp-adr/cmd/mcp-adr@latest

# Docker
docker pull ghcr.io/arkestone/mcp-adr:latest

# Pre-built binary — https://github.com/Arkestone/mcp/releases/latest
```

## Quick Start

```bash
# Run with stdio transport (default)
mcp-adr

# Run with HTTP transport
mcp-adr -transport http -addr :8083

# Build from source
go build -o mcp-adr ./servers/mcp-adr/cmd/mcp-adr
```

## MCP Client Configuration

### VS Code / GitHub Copilot

`.vscode/mcp.json`:

```json
{
  "servers": {
    "adrs": {
      "command": "mcp-adr",
      "args": ["-dirs", "${workspaceFolder}"]
    }
  }
}
```

### Claude Desktop

```json
{
  "mcpServers": {
    "adrs": {
      "command": "mcp-adr",
      "args": ["-dirs", "/path/to/repo"]
    }
  }
}
```

### Cursor

`.cursor/mcp.json`:

```json
{
  "mcpServers": {
    "adrs": {
      "command": "mcp-adr",
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
    "adrs": {
      "command": "mcp-adr",
      "args": ["-dirs", "/path/to/repo"]
    }
  }
}
```

### Claude Code

`.mcp.json`:

```json
{
  "mcpServers": {
    "adrs": {
      "command": "mcp-adr",
      "args": ["-dirs", "."]
    }
  }
}
```

### Remote (HTTP)

```json
{
  "mcpServers": {
    "adrs": {
      "type": "http",
      "url": "http://localhost:8083/mcp"
    }
  }
}
```



## Configuration

Configuration is loaded from a YAML file, environment variables, and CLI flags (later overrides earlier). Copy `config.example.yaml` and adjust:

```yaml
sources:
  dirs:
    - /path/to/local/repo       # scanned for docs/adr/*.md

  repos:
    - myorg/my-project          # GitHub repos (owner/repo or owner/repo@ref)

cache:
  dir: ~/.cache/mcp-adr
  sync_interval: 5m

llm:
  endpoint: ""                  # OpenAI-compatible endpoint
  model: gpt-4o-mini
  enabled: false

transport: stdio                # stdio | http
addr: ":8083"
```

### Environment variables

All variables are prefixed with `ADR_`:

| Variable | Description |
|---|---|
| `ADR_CONFIG` | Path to config YAML file |
| `ADR_TRANSPORT` | `stdio` or `http` |
| `ADR_ADDR` | HTTP listen address |
| `ADR_GITHUB_TOKEN` | GitHub personal access token |
| `ADR_LLM_ENDPOINT` | LLM API endpoint |
| `ADR_LLM_MODEL` | LLM model name |
| `ADR_LLM_ENABLED` | Enable LLM optimization (`true`/`false`) |

## MCP Primitives

### Resources

| URI | Description |
|---|---|
| `adrs://{source}/{id}` | Individual ADR content (Markdown + frontmatter) |
| `adrs://optimized` | All ADRs merged via LLM (or concatenated) |
| `adrs://index` | Plain-text list of all available ADRs |

### Prompts

| Name | Arguments | Description |
|---|---|---|
| `get-adrs` | `source` (opt), `status` (opt), `optimize` (opt) | Get ADRs, optionally filtered by source or status |

### Tools

| Name | Description |
|---|---|
| `refresh-adrs` | Force re-sync of all ADR sources |
| `list-adrs` | List all ADRs with optional `source` and `status` filters |
| `get-adr` | Get a single ADR by `uri` |
| `optimize-adrs` | Get consolidated ADR content with optional LLM optimization |

## ADR Format

Each ADR is a Markdown file with optional YAML frontmatter:

```markdown
---
title: Use PostgreSQL
status: accepted
date: 2023-06-01
---

## Context

We need a relational database...

## Decision

We will use PostgreSQL.
```

Supported statuses: `proposed`, `accepted`, `deprecated`, `superseded`.

If `title` is not present in the frontmatter, it is derived from the filename (e.g. `0001-use-postgresql` → `0001 Use Postgresql`).

## Scanned Directories

Within each configured source directory or GitHub repository root, the server looks for ADRs in:

- `docs/adr/`
- `docs/decisions/`
- `doc/adr/`

## License

MIT
