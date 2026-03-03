# mcp-prompts

<!-- install-badges -->
| Transport | VS Code | VS Code Insiders |
|-----------|---------|-----------------|
| stdio  | [![Install in VS Code](https://img.shields.io/badge/VS_Code-Install-0098FF?style=flat-square&logo=visualstudiocode&logoColor=white)](https://insiders.vscode.dev/redirect/mcp/install?name=mcp-prompts&config=%7B%22command%22%3A%20%22mcp-prompts%22%2C%20%22args%22%3A%20%5B%22--dirs%22%2C%20%22%24%7BworkspaceFolder%7D%22%5D%7D) | [![Install in VS Code Insiders](https://img.shields.io/badge/VS_Code_Insiders-Install-24bfa5?style=flat-square&logo=visualstudiocode&logoColor=white)](https://insiders.vscode.dev/redirect/mcp/install?name=mcp-prompts&config=%7B%22command%22%3A%20%22mcp-prompts%22%2C%20%22args%22%3A%20%5B%22--dirs%22%2C%20%22%24%7BworkspaceFolder%7D%22%5D%7D&quality=insiders) |
| HTTP   | [![Install in VS Code](https://img.shields.io/badge/VS_Code-Install-0098FF?style=flat-square&logo=visualstudiocode&logoColor=white)](https://insiders.vscode.dev/redirect/mcp/install?name=mcp-prompts&config=%7B%22type%22%3A%20%22http%22%2C%20%22url%22%3A%20%22http%3A%2F%2Flocalhost%3A8082%2Fmcp%22%7D) | [![Install in VS Code Insiders](https://img.shields.io/badge/VS_Code_Insiders-Install-24bfa5?style=flat-square&logo=visualstudiocode&logoColor=white)](https://insiders.vscode.dev/redirect/mcp/install?name=mcp-prompts&config=%7B%22type%22%3A%20%22http%22%2C%20%22url%22%3A%20%22http%3A%2F%2Flocalhost%3A8082%2Fmcp%22%7D&quality=insiders) |
| Docker | [![Install in VS Code](https://img.shields.io/badge/VS_Code-Install-0098FF?style=flat-square&logo=visualstudiocode&logoColor=white)](https://insiders.vscode.dev/redirect/mcp/install?name=mcp-prompts&config=%7B%22command%22%3A%20%22docker%22%2C%20%22args%22%3A%20%5B%22run%22%2C%20%22--rm%22%2C%20%22-i%22%2C%20%22-v%22%2C%20%22%24%7BworkspaceFolder%7D%3A%2Fworkspace%3Aro%22%2C%20%22ghcr.io%2Farkestone%2Fmcp-prompts%3Alatest%22%2C%20%22--dirs%22%2C%20%22%2Fworkspace%22%5D%7D) | [![Install in VS Code Insiders](https://img.shields.io/badge/VS_Code_Insiders-Install-24bfa5?style=flat-square&logo=visualstudiocode&logoColor=white)](https://insiders.vscode.dev/redirect/mcp/install?name=mcp-prompts&config=%7B%22command%22%3A%20%22docker%22%2C%20%22args%22%3A%20%5B%22run%22%2C%20%22--rm%22%2C%20%22-i%22%2C%20%22-v%22%2C%20%22%24%7BworkspaceFolder%7D%3A%2Fworkspace%3Aro%22%2C%20%22ghcr.io%2Farkestone%2Fmcp-prompts%3Alatest%22%2C%20%22--dirs%22%2C%20%22%2Fworkspace%22%5D%7D&quality=insiders) |
<!-- /install-badges -->

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

**Method 1: User Configuration (Recommended)**
Open the Command Palette (`Ctrl+Shift+P`) and run `MCP: Open User Configuration` to open your user `mcp.json` file and add the server configuration.

**Method 2: Workspace Configuration**
Add the configuration to `.vscode/mcp.json` in your workspace to share it with your team.

> See the [VS Code MCP documentation](https://code.visualstudio.com/docs/copilot/chat/mcp-servers) for more details.

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

## MCP API

### Resources

- **`prompts://{source}/{name}`** — Content of a single prompt or chat mode file. `{source}` is the directory basename or `owner/repo`; `{name}` is the file stem (e.g. `component` for `component.prompt.md`).
- **`prompts://optimized`** — All prompt files from every configured source merged and deduplicated (via LLM if configured, otherwise concatenated).
- **`prompts://index`** — Plain-text index listing all available prompt URIs, sources, types, and descriptions.

### Tools

- **`refresh-prompts`**
  - Force an immediate re-sync of all configured sources. Local directories are re-scanned; GitHub repo caches are fetched without waiting for the background interval.
  - No input required.

- **`list-prompts`**
  - Return metadata for every discovered file: URI, source, path, type (`prompt` or `chatmode`), description, and mode.
  - No input required.

- **`get-prompt`**
  - Return the full content and metadata of a single prompt or chat mode file.
  - Input:
    - `uri` (string, required): the resource URI in `prompts://{source}/{name}` form.

- **`optimize-prompts`**
  - Return consolidated prompt content, optionally passed through an LLM for merging and deduplication.
  - Input:
    - `optimize` (boolean, optional): override the global `llm.enabled` setting for this request.

### Prompts

- **`get-prompts`**
  - Inject prompt/chat mode content into the conversation as a prompt message, optionally filtered to a single source and/or optimized.
  - Parameters:
    - `source` (string, optional): restrict output to one source (directory basename or `owner/repo`).
    - `optimize` (string `"true"` / `"false"`, optional): override the global LLM optimization setting for this request.

## Examples

```jsonc
// list-prompts
{"tool": "list-prompts"}
// → {"entries": [{"uri": "prompts://myrepo/component", "type": "prompt", "mode": "agent", ...}]}

// get-prompt
{"tool": "get-prompt", "arguments": {"uri": "prompts://myrepo/component"}}
// → {"uri": "...", "name": "component", "mode": "agent", "type": "prompt", "content": "..."}
```
