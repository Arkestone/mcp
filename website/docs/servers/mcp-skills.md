---
title: mcp-skills
sidebar_label: mcp-skills
---


<!-- install-badges -->
| Transport | VS Code | VS Code Insiders |
|-----------|---------|-----------------|
| stdio  | [![Install in VS Code](https://img.shields.io/badge/VS_Code-Install-0098FF?style=flat-square&logo=visualstudiocode&logoColor=white)](https://insiders.vscode.dev/redirect/mcp/install?name=mcp-skills&config=%7B%22type%22%3A%20%22stdio%22%2C%20%22command%22%3A%20%22mcp-skills%22%2C%20%22args%22%3A%20%5B%22--dirs%22%2C%20%22%24%7BworkspaceFolder%7D%22%5D%7D) | [![Install in VS Code Insiders](https://img.shields.io/badge/VS_Code_Insiders-Install-24bfa5?style=flat-square&logo=visualstudiocode&logoColor=white)](https://insiders.vscode.dev/redirect/mcp/install?name=mcp-skills&config=%7B%22type%22%3A%20%22stdio%22%2C%20%22command%22%3A%20%22mcp-skills%22%2C%20%22args%22%3A%20%5B%22--dirs%22%2C%20%22%24%7BworkspaceFolder%7D%22%5D%7D&quality=insiders) |
| HTTP   | [![Install in VS Code](https://img.shields.io/badge/VS_Code-Install-0098FF?style=flat-square&logo=visualstudiocode&logoColor=white)](https://insiders.vscode.dev/redirect/mcp/install?name=mcp-skills&config=%7B%22type%22%3A%20%22http%22%2C%20%22url%22%3A%20%22http%3A%2F%2Flocalhost%3A8081%2Fmcp%22%7D) | [![Install in VS Code Insiders](https://img.shields.io/badge/VS_Code_Insiders-Install-24bfa5?style=flat-square&logo=visualstudiocode&logoColor=white)](https://insiders.vscode.dev/redirect/mcp/install?name=mcp-skills&config=%7B%22type%22%3A%20%22http%22%2C%20%22url%22%3A%20%22http%3A%2F%2Flocalhost%3A8081%2Fmcp%22%7D&quality=insiders) |
| Docker | [![Install in VS Code](https://img.shields.io/badge/VS_Code-Install-0098FF?style=flat-square&logo=visualstudiocode&logoColor=white)](https://insiders.vscode.dev/redirect/mcp/install?name=mcp-skills&config=%7B%22type%22%3A%20%22stdio%22%2C%20%22command%22%3A%20%22docker%22%2C%20%22args%22%3A%20%5B%22run%22%2C%20%22--rm%22%2C%20%22-i%22%2C%20%22-v%22%2C%20%22%24%7BworkspaceFolder%7D%3A%2Fworkspace%3Aro%22%2C%20%22ghcr.io%2Farkestone%2Fmcp-skills%3Alatest%22%2C%20%22--dirs%22%2C%20%22%2Fworkspace%22%5D%7D) | [![Install in VS Code Insiders](https://img.shields.io/badge/VS_Code_Insiders-Install-24bfa5?style=flat-square&logo=visualstudiocode&logoColor=white)](https://insiders.vscode.dev/redirect/mcp/install?name=mcp-skills&config=%7B%22type%22%3A%20%22stdio%22%2C%20%22command%22%3A%20%22docker%22%2C%20%22args%22%3A%20%5B%22run%22%2C%20%22--rm%22%2C%20%22-i%22%2C%20%22-v%22%2C%20%22%24%7BworkspaceFolder%7D%3A%2Fworkspace%3Aro%22%2C%20%22ghcr.io%2Farkestone%2Fmcp-skills%3Alatest%22%2C%20%22--dirs%22%2C%20%22%2Fworkspace%22%5D%7D&quality=insiders) |
<!-- /install-badges -->

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

## MCP API

### Resources

- **`skills://{source}/{name}`** — Full content of a single skill (SKILL.md body plus any bundled reference files). `{source}` is the directory basename or `owner/repo`; `{name}` is the skill's `name` frontmatter value.
- **`skills://optimized`** — All skills from every configured source merged and deduplicated (via LLM if configured, otherwise concatenated).
- **`skills://index`** — Plain-text index listing all available skill URIs and their sources.

### Tools

- **`refresh-skills`**
  - Force an immediate re-sync of all configured sources. Local directories are re-scanned; GitHub repo caches are fetched from the API without waiting for the background sync interval.
  - No input required.

- **`list-skills`**
  - Return metadata for every discovered skill: URI, source, name, description, and path.
  - No input required.

- **`get-skill`**
  - Return the full content (SKILL.md + references) for a single skill by name.
  - Input:
    - `name` (string, required): the skill name as declared in SKILL.md frontmatter.

- **`optimize-skills`**
  - Return consolidated skill content, optionally passed through an LLM for merging and deduplication.
  - Input:
    - `optimize` (boolean, optional): override the global `llm.enabled` setting for this request.

### Prompts

- **`get-skills`**
  - Inject skill content into the conversation as a prompt message, optionally filtered to a single source and/or optimized.
  - Parameters:
    - `source` (string, optional): restrict output to one source (directory basename or `owner/repo`).
    - `optimize` (string `"true"` / `"false"`, optional): override the global LLM optimization setting for this request.

- **`get-skill`**
  - Inject a single skill and its references into the conversation as a prompt message.
  - Parameters:
    - `name` (string, required): the skill name as declared in SKILL.md frontmatter.

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
