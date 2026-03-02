# mcp-graph — Knowledge Graph MCP Server

<!-- install-badges -->
| Transport | VS Code | VS Code Insiders |
|-----------|---------|-----------------|
| stdio | [![Install in VS Code](https://img.shields.io/badge/VS_Code-Install-0098FF?style=flat-square&logo=visualstudiocode&logoColor=white)](https://insiders.vscode.dev/redirect/mcp/install?name=mcp-graph&config=%7B%22type%22%3A%20%22stdio%22%2C%20%22command%22%3A%20%22mcp-graph%22%7D) | [![Install in VS Code Insiders](https://img.shields.io/badge/VS_Code_Insiders-Install-24bfa5?style=flat-square&logo=visualstudiocode&logoColor=white)](https://insiders.vscode.dev/redirect/mcp/install?name=mcp-graph&config=%7B%22type%22%3A%20%22stdio%22%2C%20%22command%22%3A%20%22mcp-graph%22%7D&quality=insiders) |
| HTTP  | [![Install in VS Code](https://img.shields.io/badge/VS_Code-Install-0098FF?style=flat-square&logo=visualstudiocode&logoColor=white)](https://insiders.vscode.dev/redirect/mcp/install?name=mcp-graph&config=%7B%22type%22%3A%20%22http%22%2C%20%22url%22%3A%20%22http%3A%2F%2Flocalhost%3A8085%2Fmcp%22%7D) | [![Install in VS Code Insiders](https://img.shields.io/badge/VS_Code_Insiders-Install-24bfa5?style=flat-square&logo=visualstudiocode&logoColor=white)](https://insiders.vscode.dev/redirect/mcp/install?name=mcp-graph&config=%7B%22type%22%3A%20%22http%22%2C%20%22url%22%3A%20%22http%3A%2F%2Flocalhost%3A8085%2Fmcp%22%7D&quality=insiders) |
<!-- /install-badges -->

`mcp-graph` is an MCP server that stores entities (nodes) and relationships (edges) in an in-memory knowledge graph, persisted to disk as JSON. It enables an LLM to build, query, and traverse arbitrary association networks.

## Features

- **Nodes** — entities with a label (type), a name, and optional key-value properties
- **Edges** — directed relationships between nodes with a relation type and optional properties
- **Graph traversal** — BFS-based shortest-path search, N-hop neighbor exploration
- **Filters** — find nodes by label/name substring, filter neighbors by relation type
- **Persistence** — atomic JSON file writes; reloads on restart with full index rebuild
- **Transports** — `stdio` (default) and HTTP (Streamable HTTP)

## MCP API

### Resources

- **`graph://stats`** — Node and edge counts plus the full list of unique relation types currently in the graph.
- **`graph://node/{id}`** — Full details of a single node by ID, including label, name, properties, and creation timestamp.

### Tools

- **`add-node`**
  - Create a new entity node in the graph.
  - Input:
    - `label` (string, required): the node type / category (e.g. `"Person"`, `"Technology"`).
    - `name` (string, required): a human-readable name for the entity.
    - `properties` (object, optional): arbitrary key-value metadata to attach to the node.

- **`add-edge`**
  - Create a directed relationship between two existing nodes.
  - Input:
    - `from` (string, required): ID of the source node.
    - `to` (string, required): ID of the target node.
    - `relation` (string, required): relationship type in active voice (e.g. `"knows"`, `"depends_on"`).
    - `properties` (object, optional): arbitrary key-value metadata to attach to the edge.

- **`find-nodes`**
  - Search for nodes by label and/or name substring. Returns all matching nodes.
  - Input:
    - `label` (string, optional): filter by node label (exact match).
    - `name` (string, optional): filter by node name (case-insensitive substring match).

- **`get-node`**
  - Retrieve a single node by its ID, including all properties and metadata.
  - Input:
    - `id` (string, required): the node ID.

- **`neighbors`**
  - List the immediate neighbors of a node, with optional direction and relation-type filters.
  - Input:
    - `id` (string, required): the node ID to expand.
    - `direction` (string `"in"` / `"out"` / `"both"`, optional): which edge direction to follow (default: `"both"`).
    - `relation` (string, optional): restrict results to edges of this relation type.

- **`shortest-path`**
  - Find the shortest path between two nodes using breadth-first search.
  - Input:
    - `from` (string, required): ID of the start node.
    - `to` (string, required): ID of the end node.

- **`remove-node`**
  - Delete a node and all edges incident to it.
  - Input:
    - `id` (string, required): the node ID to delete.

- **`remove-edge`**
  - Delete a single edge by its ID.
  - Input:
    - `id` (string, required): the edge ID to delete.

- **`list-relations`**
  - Return all unique relation type strings currently present in the graph.
  - No input required.

## Installation

```bash
# go install (requires Go 1.24+)
go install github.com/Arkestone/mcp/servers/mcp-graph/cmd/mcp-graph@latest

# Docker
docker pull ghcr.io/arkestone/mcp-graph:latest

# Pre-built binary — https://github.com/Arkestone/mcp/releases/latest
```

## Quick Start

### VS Code / GitHub Copilot

`.vscode/mcp.json`:

```json
{
  "servers": {
    "graph": {
      "command": "mcp-graph",
      "env": {
        "GRAPH_DIRS": "~/.local/share/mcp-graph"
      }
    }
  }
}
```

**Method 1: User Configuration (Recommended)**
Open the Command Palette (`Ctrl+Shift+P`) and run `MCP: Open User Configuration` to open your user `mcp.json` file and add the server configuration.

**Method 2: Workspace Configuration**
Add the configuration to `.vscode/mcp.json` in your workspace to share it with your team.

> See the [VS Code MCP documentation](https://code.visualstudio.com/docs/copilot/model-context-protocol) for more details.

### Cursor

`.cursor/mcp.json`:

```json
{
  "mcpServers": {
    "graph": {
      "command": "mcp-graph"
    }
  }
}
```

### Windsurf

`~/.codeium/windsurf/mcp_config.json`:

```json
{
  "mcpServers": {
    "graph": {
      "command": "mcp-graph"
    }
  }
}
```

### Claude Code

`.mcp.json`:

```json
{
  "mcpServers": {
    "graph": {
      "command": "mcp-graph"
    }
  }
}
```

### HTTP

```bash
mcp-graph --transport http --addr :8085
```

```json
{
  "mcpServers": {
    "graph": {
      "type": "http",
      "url": "http://localhost:8085/mcp"
    }
  }
}
```

### Docker

```bash
docker run -d \
  -v ~/.local/share/mcp-graph:/home/mcp/.local/share/mcp-graph \
  -p 8085:8085 \
  ghcr.io/arkestone/mcp-graph:latest \
  --transport http --addr :8085
```

## Configuration

| Source | Format |
|--------|--------|
| YAML file | `--config config.yaml` or auto-detected `config.yaml` |
| Environment variables | `GRAPH_ADDR`, `GRAPH_TRANSPORT`, `GRAPH_DIRS` |
| CLI flags | `--addr`, `--transport`, `--dirs` |

See [`config.example.yaml`](config.example.yaml) for all options.

## Data Model

### Node
```json
{
  "id": "a1b2c3d4",
  "label": "Person",
  "name": "Alice",
  "props": { "role": "engineer" },
  "created_at": "2026-03-01T22:45:32Z"
}
```

### Edge
```json
{
  "id": "x9y8z7w6",
  "from": "a1b2c3d4",
  "to": "e5f6g7h8",
  "relation": "knows",
  "props": { "since": "2024" },
  "created_at": "2026-03-01T22:45:32Z"
}
```

## Example Usage

```
# Build a knowledge graph about a tech stack
add-node label="Technology" name="Go"
add-node label="Technology" name="PostgreSQL"
add-node label="Project"    name="MyApp"

add-edge from=<go-id>   to=<myapp-id>    relation="used_by"
add-edge from=<pg-id>   to=<myapp-id>    relation="used_by"
add-edge from=<myapp-id> to=<go-id>      relation="depends_on"

neighbors id=<myapp-id> direction=in  relation=used_by
shortest-path from=<go-id> to=<myapp-id>
```
