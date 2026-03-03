# Client Configuration Examples

Ready-to-use MCP client configuration snippets for all Arkestone MCP servers.

## VS Code (`.vscode/mcp.json`)

```json
{
  "servers": {
    "instructions": {
      "command": "mcp-instructions",
      "args": ["-dirs", "/path/to/repo"]
    },
    "skills": {
      "command": "mcp-skills",
      "args": ["-dirs", "/path/to/skills"]
    },
    "prompts": {
      "command": "mcp-prompts",
      "args": ["-dirs", "/path/to/repo"]
    },
    "adrs": {
      "command": "mcp-adr",
      "args": ["-dirs", "/path/to/repo"]
    },
    "memory": {
      "command": "mcp-memory"
    },
    "graph": {
      "command": "mcp-graph"
    }
  }
}
```

## Claude Desktop (`~/Library/Application Support/Claude/claude_desktop_config.json`)

```json
{
  "mcpServers": {
    "instructions": {
      "command": "mcp-instructions",
      "args": ["-dirs", "/path/to/repo"]
    },
    "skills": {
      "command": "mcp-skills",
      "args": ["-dirs", "/path/to/skills"]
    },
    "prompts": {
      "command": "mcp-prompts",
      "args": ["-dirs", "/path/to/repo"]
    },
    "adrs": {
      "command": "mcp-adr",
      "args": ["-dirs", "/path/to/repo"]
    },
    "memory": {
      "command": "mcp-memory",
      "env": {
        "MEMORY_DIR": "~/.local/share/mcp-memory"
      }
    },
    "graph": {
      "command": "mcp-graph",
      "env": {
        "GRAPH_DIR": "~/.local/share/mcp-graph"
      }
    }
  }
}
```

## Remote (HTTP transport)

```json
{
  "mcpServers": {
    "instructions": { "url": "http://localhost:8080" },
    "skills":       { "url": "http://localhost:8081" },
    "prompts":      { "url": "http://localhost:8082" },
    "adrs":         { "url": "http://localhost:8083" },
    "memory":       { "url": "http://localhost:8084" },
    "graph":        { "url": "http://localhost:8085" }
  }
}
```

## Docker Compose

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

  mcp-graph:
    image: ghcr.io/arkestone/mcp-graph:latest
    ports: ["8085:8085"]
    environment:
      GRAPH_TRANSPORT: http
      GRAPH_DIR: /data
    volumes: ["graph-data:/data"]

volumes:
  memory-data:
  graph-data:
```

## Stdio (direct process, for local use)

```bash
# Instructions server
mcp-instructions -dirs /path/to/repo

# Skills server
mcp-skills -dirs /path/to/skills

# Prompts server
mcp-prompts -dirs /path/to/repo

# ADR server
mcp-adr -dirs /path/to/repo

# Memory server (no dirs needed — uses MEMORY_DIR for storage)
mcp-memory

# Graph server (no dirs needed — uses GRAPH_DIR for storage)
mcp-graph
```

