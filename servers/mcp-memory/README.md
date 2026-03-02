# mcp-memory

<!-- install-badges -->
| Transport | VS Code | VS Code Insiders |
|-----------|---------|-----------------|
| stdio  | [![Install in VS Code](https://img.shields.io/badge/VS_Code-Install-0098FF?style=flat-square&logo=visualstudiocode&logoColor=white)](https://insiders.vscode.dev/redirect/mcp/install?name=mcp-memory&config=%7B%22settings%22%3A%20%7B%22mcp%22%3A%20%7B%22servers%22%3A%20%7B%22mcp-memory%22%3A%20%7B%22type%22%3A%20%22stdio%22%2C%20%22command%22%3A%20%22mcp-memory%22%7D%7D%2C%20%22inputs%22%3A%20%5B%5D%7D%7D%7D) | [![Install in VS Code Insiders](https://img.shields.io/badge/VS_Code_Insiders-Install-24bfa5?style=flat-square&logo=visualstudiocode&logoColor=white)](https://insiders.vscode.dev/redirect/mcp/install?name=mcp-memory&config=%7B%22settings%22%3A%20%7B%22mcp%22%3A%20%7B%22servers%22%3A%20%7B%22mcp-memory%22%3A%20%7B%22type%22%3A%20%22stdio%22%2C%20%22command%22%3A%20%22mcp-memory%22%7D%7D%2C%20%22inputs%22%3A%20%5B%5D%7D%7D%7D&quality=insiders) |
| HTTP   | [![Install in VS Code](https://img.shields.io/badge/VS_Code-Install-0098FF?style=flat-square&logo=visualstudiocode&logoColor=white)](https://insiders.vscode.dev/redirect/mcp/install?name=mcp-memory&config=%7B%22settings%22%3A%20%7B%22mcp%22%3A%20%7B%22servers%22%3A%20%7B%22mcp-memory%22%3A%20%7B%22type%22%3A%20%22http%22%2C%20%22url%22%3A%20%22http%3A%2F%2Flocalhost%3A8084%2Fmcp%22%7D%7D%2C%20%22inputs%22%3A%20%5B%5D%7D%7D%7D) | [![Install in VS Code Insiders](https://img.shields.io/badge/VS_Code_Insiders-Install-24bfa5?style=flat-square&logo=visualstudiocode&logoColor=white)](https://insiders.vscode.dev/redirect/mcp/install?name=mcp-memory&config=%7B%22settings%22%3A%20%7B%22mcp%22%3A%20%7B%22servers%22%3A%20%7B%22mcp-memory%22%3A%20%7B%22type%22%3A%20%22http%22%2C%20%22url%22%3A%20%22http%3A%2F%2Flocalhost%3A8084%2Fmcp%22%7D%7D%2C%20%22inputs%22%3A%20%5B%5D%7D%7D%7D&quality=insiders) |
| Docker | [![Install in VS Code](https://img.shields.io/badge/VS_Code-Install-0098FF?style=flat-square&logo=visualstudiocode&logoColor=white)](https://insiders.vscode.dev/redirect/mcp/install?name=mcp-memory&config=%7B%22settings%22%3A%20%7B%22mcp%22%3A%20%7B%22servers%22%3A%20%7B%22mcp-memory%22%3A%20%7B%22type%22%3A%20%22stdio%22%2C%20%22command%22%3A%20%22docker%22%2C%20%22args%22%3A%20%5B%22run%22%2C%20%22--rm%22%2C%20%22-i%22%2C%20%22-v%22%2C%20%22mcp-memory%3A%2Fdata%22%2C%20%22ghcr.io%2Farkestone%2Fmcp-memory%3Alatest%22%5D%7D%7D%2C%20%22inputs%22%3A%20%5B%5D%7D%7D%7D) | [![Install in VS Code Insiders](https://img.shields.io/badge/VS_Code_Insiders-Install-24bfa5?style=flat-square&logo=visualstudiocode&logoColor=white)](https://insiders.vscode.dev/redirect/mcp/install?name=mcp-memory&config=%7B%22settings%22%3A%20%7B%22mcp%22%3A%20%7B%22servers%22%3A%20%7B%22mcp-memory%22%3A%20%7B%22type%22%3A%20%22stdio%22%2C%20%22command%22%3A%20%22docker%22%2C%20%22args%22%3A%20%5B%22run%22%2C%20%22--rm%22%2C%20%22-i%22%2C%20%22-v%22%2C%20%22mcp-memory%3A%2Fdata%22%2C%20%22ghcr.io%2Farkestone%2Fmcp-memory%3Alatest%22%5D%7D%7D%2C%20%22inputs%22%3A%20%5B%5D%7D%7D%7D&quality=insiders) |
<!-- /install-badges -->

An MCP server that provides persistent, searchable memory storage for AI assistants. Memories are stored as Markdown files on disk and can be retrieved by text query or tag.

## Features

- **Persistent storage** — memories survive process restarts; stored as plain Markdown files on disk
- **Tag-based organisation** — attach arbitrary tags to any memory and filter by them later
- **Full-text search** — `recall` searches both content and tags
- **Dual transport** — supports both stdio (for local tools) and Streamable HTTP (for remote deployments)
- **MCP primitives** — memories are exposed as Resources and Tools

## Installation

```bash
# go install (requires Go 1.24+)
go install github.com/Arkestone/mcp/servers/mcp-memory/cmd/mcp-memory@latest

# Docker
docker pull ghcr.io/arkestone/mcp-memory:latest

# Pre-built binary — https://github.com/Arkestone/mcp/releases/latest
```

## Getting Started

```bash
# Run with stdio transport (default)
mcp-memory

# Run with HTTP transport
mcp-memory -transport http -addr :8084

# Build from source
make build-memory   # → ./bin/mcp-memory
```

## Configuration

Configuration is loaded in layers (each overrides the previous):

1. **YAML file** — `config.yaml` in the working directory, or specify with `-config path/to/config.yaml`
2. **Environment variables**
3. **CLI flags**

See [`config.example.yaml`](config.example.yaml) for all options.

### Environment Variables

| Variable | Description | Default |
|---|---|---|
| `MEMORY_CONFIG` | Path to YAML config file | — |
| `MEMORY_DIR` | Directory where memories are stored | `~/.local/share/mcp-memory` |
| `MEMORY_TRANSPORT` | `stdio` (default) or `http` | `stdio` |
| `MEMORY_ADDR` | HTTP listen address | `:8084` |

### CLI Flags

```
-config       Path to YAML config file
-memory-dir   Directory where memories are stored
-transport    Transport: stdio (default) or http
-addr         HTTP listen address (default :8084)
```

## MCP API

### Resources

- **`memory://{id}`** — Content of a single memory by its unique ID, including body text, tags, and creation timestamp.
- **`memory://all`** — All stored memories as newline-separated plain text, suitable for bulk context injection.

### Tools

- **`remember`**
  - Store a new memory with optional tags for later filtering and retrieval.
  - Input:
    - `content` (string, required): the text to remember.
    - `tags` (string array, optional): labels to attach to this memory (e.g. `["go", "architecture"]`).

- **`recall`**
  - Search memories by full-text query and/or tags. Returns all memories whose content or tags match.
  - Input:
    - `query` (string, optional): text to search for in memory content and tags.
    - `tags` (string array, optional): filter to memories that carry all of the given tags.

- **`forget`**
  - Permanently delete a memory by its ID.
  - Input:
    - `id` (string, required): the ID of the memory to delete.

- **`list-memories`**
  - List all stored memories, optionally filtered by tags.
  - Input:
    - `tags` (string array, optional): return only memories that carry all of the given tags.

## Memory Format

Each memory is stored as a plain Markdown file with a YAML-like header:

```
id: 01HXYZ...
tags: [go, architecture]
created: 2025-03-01T12:00:00Z

The content of the memory goes here.
It can span multiple lines.
```

## Docker

```bash
# From the repo root
make docker-memory

# stdio mode
docker run -i -v ~/.local/share/mcp-memory:/data \
  -e MEMORY_DIR=/data \
  ghcr.io/arkestone/mcp-memory:latest

# HTTP mode
docker run -p 8084:8084 \
  -v ~/.local/share/mcp-memory:/data \
  -e MEMORY_DIR=/data \
  -e MEMORY_TRANSPORT=http \
  ghcr.io/arkestone/mcp-memory:latest
```

## MCP Client Configuration

### VS Code / GitHub Copilot

`.vscode/mcp.json`:

```json
{
  "servers": {
    "memory": {
      "command": "mcp-memory"
    }
  }
}
```

**Method 1: User Configuration (Recommended)**
Open the Command Palette (`Ctrl+Shift+P`) and run `MCP: Open User Configuration` to open your user `mcp.json` file and add the server configuration.

**Method 2: Workspace Configuration**
Add the configuration to `.vscode/mcp.json` in your workspace to share it with your team.

> See the [VS Code MCP documentation](https://code.visualstudio.com/docs/copilot/model-context-protocol) for more details.

### Claude Desktop

**macOS**: `~/Library/Application Support/Claude/claude_desktop_config.json`

```json
{
  "mcpServers": {
    "memory": {
      "command": "mcp-memory",
      "env": { "MEMORY_DIR": "~/.local/share/mcp-memory" }
    }
  }
}
```

### Cursor

`.cursor/mcp.json` or `~/.cursor/mcp.json`:

```json
{
  "mcpServers": {
    "memory": {
      "command": "mcp-memory",
      "env": { "MEMORY_DIR": "~/.local/share/mcp-memory" }
    }
  }
}
```

### Windsurf

`~/.codeium/windsurf/mcp_config.json`:

```json
{
  "mcpServers": {
    "memory": {
      "command": "mcp-memory",
      "env": { "MEMORY_DIR": "~/.local/share/mcp-memory" }
    }
  }
}
```

### Claude Code

`.mcp.json` or `~/.mcp.json`:

```json
{
  "mcpServers": {
    "memory": {
      "command": "mcp-memory"
    }
  }
}
```

### Remote (HTTP)

```json
{
  "mcpServers": {
    "memory": {
      "type": "http",
      "url": "http://localhost:8084/mcp"
    }
  }
}
```
