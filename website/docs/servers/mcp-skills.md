---
title: mcp-skills
sidebar_label: mcp-skills
---


An MCP server that dynamically serves [GitHub Copilot skills](https://docs.github.com/en/copilot) from local directories and GitHub repositories.

It scans configured sources for skills following the `skills/{skill-name}/SKILL.md` convention and exposes them via the [Model Context Protocol](https://modelcontextprotocol.io).

## Features

- **Skill discovery** — automatically scans `skills/` directories for `SKILL.md` files with YAML frontmatter
- **Reference bundling** — includes optional `references/` alongside each skill for supporting documentation
- **GitHub repo caching** — remote repos are cached locally and synced periodically in the background
- **LLM optimization** — optionally consolidate multiple skills via an OpenAI-compatible endpoint
- **Dual transport** — supports both stdio and Streamable HTTP
- **MCP primitives** — skills are exposed as Resources, Prompts, and Tools

## Skills Format

Each skill lives in its own directory under `skills/` and contains a `SKILL.md` file with YAML frontmatter:

```
skills/
├── code-review/
│   ├── SKILL.md
│   └── references/
│       └── best-practices.md
├── testing/
│   ├── SKILL.md
│   └── references/
│       ├── unit-testing.md
│       └── integration-testing.md
└── documentation/
    └── SKILL.md
```

### SKILL.md

```markdown
---
name: code-review
description: Guidelines for performing thorough code reviews
---

Review code for correctness, readability, and maintainability...
```

The YAML frontmatter must include `name` and `description`. The Markdown body contains the skill content that will be served to Copilot.

## Installation

```bash
# go install (requires Go 1.24+)
go install github.com/Arkestone/mcp/servers/mcp-skills/cmd/mcp-skills@latest

# Docker
docker pull ghcr.io/arkestone/mcp-skills:latest

# Pre-built binary — https://github.com/Arkestone/mcp/releases/latest
```

## Getting Started

```bash
# Serve skills from a local directory (stdio)
mcp-skills -dirs /path/to/skills-repo

# Serve from a GitHub repo with HTTP transport
export GITHUB_TOKEN=ghp_...
mcp-skills -repos github/awesome-copilot -transport http -addr :8081

# Build from source
make build-skills   # → ./bin/mcp-skills
```

## Configuration

Configuration is loaded in layers (each overrides the previous):

1. **YAML file** — `config.yaml` in the working directory, or specify with `-config path/to/config.yaml`
2. **Environment variables**
3. **CLI flags**

See [`config.example.yaml`](https://github.com/Arkestone/mcp/blob/main/servers/mcp-skills/config.example.yaml) for all options.

### Environment Variables

| Variable | Description |
|---|---|
| `SKILLS_CONFIG` | Path to YAML config file |
| `SKILLS_DIRS` | Comma-separated local directories |
| `SKILLS_REPOS` | Comma-separated GitHub repos (`owner/repo` or `owner/repo@ref`) |
| `SKILLS_TRANSPORT` | `stdio` (default) or `http` |
| `SKILLS_ADDR` | HTTP listen address (default `:8081`) |
| `SKILLS_CACHE_DIR` | Local cache directory (default `~/.cache/mcp-skills`) |
| `SKILLS_SYNC_INTERVAL` | Sync interval for remote repos (default `5m`) |
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
-addr            HTTP listen address (default :8081)
-cache-dir       Local cache directory
-sync-interval   Sync interval (e.g. 5m, 1h)
-llm-endpoint    OpenAI-compatible endpoint URL
-llm-model       LLM model name
```

## MCP Primitives

### Resources

| URI | Description |
|---|---|
| `skills://{source}/{skill-name}` | Individual skill content (SKILL.md + references) |
| `skills://optimized` | All skills merged via LLM (or concatenated) |
| `skills://index` | List of all available skill URIs |

### Prompts

| Name | Arguments | Description |
|---|---|---|
| `get-skills` | `source` (optional), `optimize` (optional) | Get skills as a prompt, optionally filtered and optimized |

### Tools

| Name | Arguments | Description |
|---|---|---|
| `refresh` | — | Force-sync all remote repo caches |
| `list-skills` | — | List all available skills |
| `optimize-skills` | `source` (optional), `optimize` (optional) | Get consolidated skills with optional LLM optimization |

## Docker

```bash
# From the repo root
make docker-skills

# stdio mode
docker run -i ghcr.io/arkestone/mcp-skills:latest -dirs /data

# HTTP mode
docker run -p 8081:8081 \
  -e SKILLS_REPOS=github/awesome-copilot \
  -e GITHUB_TOKEN=ghp_... \
  -e SKILLS_TRANSPORT=http \
  ghcr.io/arkestone/mcp-skills:latest
```

## Example Configuration

```yaml
sources:
  dirs:
    - ./my-skills
  repos:
    - github/awesome-copilot

cache:
  dir: ~/.cache/mcp-skills
  sync_interval: 5m

llm:
  endpoint: https://api.openai.com/v1
  model: gpt-4o-mini
  enabled: false

transport: stdio
addr: ":8081"
```

## MCP Client Configuration

### VS Code / GitHub Copilot

`.vscode/mcp.json`:

```json
{
  "servers": {
    "skills": {
      "command": "mcp-skills",
      "args": ["-dirs", "${workspaceFolder}"]
    }
  }
}
```

### Claude Desktop

**macOS**: `~/Library/Application Support/Claude/claude_desktop_config.json`

```json
{
  "mcpServers": {
    "skills": {
      "command": "mcp-skills",
      "args": ["-repos", "github/awesome-copilot"]
    }
  }
}
```

### Cursor

`.cursor/mcp.json`:

```json
{
  "mcpServers": {
    "skills": {
      "command": "mcp-skills",
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
    "skills": {
      "command": "mcp-skills",
      "args": ["-dirs", "/path/to/skills-repo"]
    }
  }
}
```

### Claude Code

`.mcp.json`:

```json
{
  "mcpServers": {
    "skills": {
      "command": "mcp-skills",
      "args": ["-dirs", "."]
    }
  }
}
```

### Remote (HTTP)

```json
{
  "mcpServers": {
    "skills": {
      "type": "http",
      "url": "http://localhost:8081/mcp"
    }
  }
}
```
