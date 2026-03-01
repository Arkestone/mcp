---
slug: /
sidebar_position: 1
title: Introduction
---

# Arkestone MCP Servers

[![CI](https://github.com/Arkestone/mcp/actions/workflows/ci.yml/badge.svg)](https://github.com/Arkestone/mcp/actions/workflows/ci.yml)
[![Go Report Card](https://goreportcard.com/badge/github.com/Arkestone/mcp)](https://goreportcard.com/report/github.com/Arkestone/mcp)
[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://github.com/Arkestone/mcp/blob/main/LICENSE)
[![Release](https://img.shields.io/github/v/release/Arkestone/mcp)](https://github.com/Arkestone/mcp/releases/latest)

A suite of [Model Context Protocol](https://modelcontextprotocol.io) (MCP) servers for [GitHub Copilot](https://github.com/features/copilot) and other AI coding assistants. Each server dynamically serves a different type of context from local directories and GitHub repositories.

## Available Servers

| Server | Description | Port |
|--------|-------------|------|
| [mcp-instructions](/docs/servers/mcp-instructions) | Copilot custom instruction files from local dirs and GitHub repos | `:8080` |
| [mcp-skills](/docs/servers/mcp-skills) | Copilot skills with frontmatter metadata | `:8081` |
| [mcp-prompts](/docs/servers/mcp-prompts) | VS Code Copilot prompt and chat mode files | `:8082` |
| [mcp-adr](/docs/servers/mcp-adr) | Architecture Decision Records | `:8083` |
| [mcp-memory](/docs/servers/mcp-memory) | Persistent memory store across sessions | `:8084` |
| [mcp-graph](/docs/servers/mcp-graph) | Knowledge graph with entity and relationship storage | `:8085` |

## Installation

### go install

Requires [Go 1.24+](https://go.dev/dl/).

```bash
go install github.com/Arkestone/mcp/servers/mcp-instructions/cmd/mcp-instructions@latest
go install github.com/Arkestone/mcp/servers/mcp-skills/cmd/mcp-skills@latest
go install github.com/Arkestone/mcp/servers/mcp-prompts/cmd/mcp-prompts@latest
go install github.com/Arkestone/mcp/servers/mcp-adr/cmd/mcp-adr@latest
go install github.com/Arkestone/mcp/servers/mcp-memory/cmd/mcp-memory@latest
go install github.com/Arkestone/mcp/servers/mcp-graph/cmd/mcp-graph@latest
```

### Docker

```bash
docker pull ghcr.io/arkestone/mcp-instructions:latest
docker pull ghcr.io/arkestone/mcp-skills:latest
docker pull ghcr.io/arkestone/mcp-prompts:latest
docker pull ghcr.io/arkestone/mcp-adr:latest
docker pull ghcr.io/arkestone/mcp-memory:latest
docker pull ghcr.io/arkestone/mcp-graph:latest
```

### Pre-built Binaries

Download from [GitHub Releases](https://github.com/Arkestone/mcp/releases):

```bash
# Example: Linux amd64
curl -L https://github.com/Arkestone/mcp/releases/latest/download/mcp-instructions_linux_amd64.tar.gz | tar xz
mv mcp-instructions /usr/local/bin/
```

### Build from Source

```bash
git clone https://github.com/Arkestone/mcp.git
cd mcp
make build
```

## Quick Start

```bash
# stdio mode (default) — for IDE integrations
mcp-instructions --dirs /path/to/repo

# HTTP mode — for multi-client access
mcp-instructions --transport http --dirs /path/to/repo
```

## Transport Mechanisms

| Transport | When to use |
|-----------|------------|
| `stdio` | IDE integrations (VS Code, Cursor, etc.) |
| `http` | Remote/multi-client access, Docker |

## MCP Client Configuration

### VS Code

Add to `.vscode/mcp.json` in your workspace:

```json
{
  "servers": {
    "mcp-instructions": {
      "type": "stdio",
      "command": "mcp-instructions",
      "args": ["--dirs", "${workspaceFolder}"]
    },
    "mcp-skills": {
      "type": "stdio",
      "command": "mcp-skills",
      "args": ["--dirs", "${workspaceFolder}"]
    }
  }
}
```

### Claude Desktop

Add to `claude_desktop_config.json`:

```json
{
  "mcpServers": {
    "mcp-instructions": {
      "command": "mcp-instructions",
      "args": ["--dirs", "/path/to/your/repo"]
    },
    "mcp-skills": {
      "command": "mcp-skills",
      "args": ["--dirs", "/path/to/your/repo"]
    }
  }
}
```

### Cursor / Windsurf

```json
{
  "mcpServers": {
    "mcp-instructions": {
      "command": "mcp-instructions",
      "args": ["--dirs", "/path/to/repo"]
    }
  }
}
```

## LLM Optimization

All servers support optional LLM optimization to summarize and improve served content:

```bash
export LLM_ENDPOINT=https://api.openai.com/v1
export LLM_API_KEY=sk-...
export LLM_MODEL=gpt-4o-mini

mcp-instructions --dirs /path/to/repo
```

Enable by default in config:

```yaml
llm:
  endpoint: https://api.openai.com/v1
  model: gpt-4o-mini
  enabled: true
```
