# mcp-memory

An MCP server that provides persistent, searchable memory storage for AI assistants. Memories are stored as Markdown files on disk and can be retrieved by text query or tag.

## Features

- **Persistent storage** — memories survive process restarts; stored as plain Markdown files on disk
- **Tag-based organisation** — attach arbitrary tags to any memory and filter by them later
- **Full-text search** — `recall` searches both content and tags
- **Dual transport** — supports both stdio (for local tools) and Streamable HTTP (for remote deployments)
- **MCP primitives** — memories are exposed as Resources and Tools

## Getting Started

```bash
# From the repo root
make build-memory

# Run with stdio transport (default)
./bin/mcp-memory

# Run with HTTP transport
./bin/mcp-memory -transport http -addr :8084
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

## MCP Primitives

### Resources

| URI | Description |
|---|---|
| `memory://{id}` | Individual memory by ID |
| `memory://all` | All stored memories (newline-separated) |

### Tools

| Name | Arguments | Description |
|---|---|---|
| `remember` | `content` (string), `tags` ([]string, optional) | Store a new memory with optional tags |
| `recall` | `query` (string, optional), `tags` ([]string, optional) | Search memories by text and/or tags |
| `forget` | `id` (string) | Delete a memory by ID |
| `list-memories` | `tags` ([]string, optional) | List all memories, optionally filtered by tags |

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

### VS Code (`.vscode/mcp.json`)

```json
{
  "servers": {
    "memory": {
      "command": "mcp-memory",
      "args": []
    }
  }
}
```

### Claude Desktop

```json
{
  "mcpServers": {
    "memory": {
      "command": "mcp-memory",
      "args": [],
      "env": {
        "MEMORY_DIR": "/path/to/memories"
      }
    }
  }
}
```

### Remote (HTTP)

```json
{
  "mcpServers": {
    "memory": {
      "url": "http://localhost:8084"
    }
  }
}
```
