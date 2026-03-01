# mcp-instructions

An MCP server that dynamically serves [GitHub Copilot custom instructions](https://docs.github.com/en/copilot/customize-github-copilot/adding-custom-instructions-for-github-copilot) from local directories and GitHub repositories.

It discovers `.github/copilot-instructions.md` and `.github/instructions/**/*.instructions.md` files from configured sources and exposes them via the [Model Context Protocol](https://modelcontextprotocol.io).

## Features

- **On-demand loading** — local directories are read live; no eager loading
- **GitHub repo caching** — remote repos are cached locally and synced periodically in the background
- **LLM optimization** — optionally merge/deduplicate instructions from multiple sources via an OpenAI-compatible endpoint
- **Dual transport** — supports both stdio (for Copilot CLI / local tools) and Streamable HTTP (for remote deployments)
- **MCP primitives** — instructions are exposed as Resources, Prompts, and Tools

## Getting Started

```bash
# From the repo root
make build-instructions

# Serve instructions from a local directory (stdio)
./bin/mcp-instructions -dirs /path/to/repo

# Serve from a GitHub repo with HTTP transport
export GITHUB_TOKEN=ghp_...
./bin/mcp-instructions -repos github/awesome-copilot -transport http -addr :8080
```

## Configuration

Configuration is loaded in layers (each overrides the previous):

1. **YAML file** — `config.yaml` in the working directory, or specify with `-config path/to/config.yaml`
2. **Environment variables**
3. **CLI flags**

See [`config.example.yaml`](config.example.yaml) for all options.

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

### Copilot CLI / VS Code

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
