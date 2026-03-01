# mcp-prompts

An MCP server that serves VS Code Copilot prompt files (`.github/prompts/*.prompt.md`) and chat mode files (`.github/chatmodes/*.chatmode.md`) from local directories and GitHub repositories.

## What it does

`mcp-prompts` scans configured local directories and GitHub repos for Copilot prompt/chat mode files and exposes them as MCP resources, prompts, and tools. Frontmatter fields (`description`, `mode`) are parsed and surfaced in tool responses.

## Installation

```bash
# go install (requires Go 1.24+)
go install github.com/Arkestone/mcp/servers/mcp-prompts/cmd/mcp-prompts@latest

# Docker
docker pull ghcr.io/arkestone/mcp-prompts:latest

# Pre-built binary — https://github.com/Arkestone/mcp/releases/latest
```

## Quick start

```bash
# Run with a local directory (stdio transport)
mcp-prompts -dirs ./my-repo

# Run as HTTP server
mcp-prompts -transport http -addr :8082 -dirs ./my-repo

# Build from source
go build -o mcp-prompts ./servers/mcp-prompts/cmd/mcp-prompts
```

## MCP Client Configuration

### VS Code / GitHub Copilot

`.vscode/mcp.json`:

```json
{
  "servers": {
    "prompts": {
      "command": "mcp-prompts",
      "args": ["-dirs", "${workspaceFolder}"]
    }
  }
}
```

### Claude Desktop

```json
{
  "mcpServers": {
    "prompts": {
      "command": "mcp-prompts",
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
    "prompts": {
      "command": "mcp-prompts",
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
    "prompts": {
      "command": "mcp-prompts",
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
    "prompts": {
      "command": "mcp-prompts",
      "args": ["-dirs", "."]
    }
  }
}
```

### Remote (HTTP)

```json
{
  "mcpServers": {
    "prompts": {
      "type": "http",
      "url": "http://localhost:8082/mcp"
    }
  }
}
```

## Configuration

Copy `config.example.yaml` and adjust:

```yaml
sources:
  dirs:
    - /path/to/local/repo
  repos:
    - owner/repo          # latest default branch
    - owner/repo@main     # specific branch/tag/SHA

cache:
  dir: ~/.cache/mcp-prompts
  sync_interval: 5m

llm:
  endpoint: https://api.openai.com/v1
  model: gpt-4o-mini
  enabled: false

transport: stdio   # or http
addr: ":8082"
```

### Environment variables

All config keys are available as `PROMPTS_` prefixed env vars, e.g.:

| Variable | Description |
|---|---|
| `PROMPTS_TRANSPORT` | `stdio` or `http` |
| `PROMPTS_ADDR` | HTTP listen address (default `:8082`) |
| `PROMPTS_SOURCES_DIRS` | Comma-separated local directories |
| `PROMPTS_SOURCES_REPOS` | Comma-separated GitHub repos |
| `PROMPTS_GITHUB_TOKEN` | GitHub personal access token |
| `PROMPTS_CACHE_DIR` | Cache directory |
| `PROMPTS_LLM_ENDPOINT` | OpenAI-compatible API endpoint |
| `PROMPTS_LLM_APIKEY` | LLM API key |
| `PROMPTS_LLM_MODEL` | LLM model name |
| `PROMPTS_LLM_ENABLED` | Enable LLM optimization (`true`/`false`) |

## MCP primitives

### Resources

| URI | Description |
|---|---|
| `prompts://{source}/{name}` | Single prompt or chat mode file |
| `prompts://optimized` | All files merged via LLM (or concatenated) |
| `prompts://index` | Plain-text index of all files |

### Prompts

| Name | Arguments | Description |
|---|---|---|
| `get-prompts` | `source` (optional), `optimize` (optional) | Get all prompt/chat mode content |

### Tools

| Name | Description |
|---|---|
| `refresh-prompts` | Force immediate re-sync from all sources |
| `list-prompts` | List all files with URI, source, path, type, description, mode |
| `get-prompt` | Get a single prompt by URI |
| `optimize-prompts` | Get consolidated content via LLM or concatenation |

## Examples

```jsonc
// list-prompts
{"tool": "list-prompts"}
// → {"entries": [{"uri": "prompts://myrepo/component", "type": "prompt", "mode": "agent", ...}]}

// get-prompt
{"tool": "get-prompt", "arguments": {"uri": "prompts://myrepo/component"}}
// → {"uri": "...", "name": "component", "mode": "agent", "type": "prompt", "content": "..."}
```
