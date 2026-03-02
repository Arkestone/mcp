# Arkestone MCP Servers

[![CI](https://github.com/Arkestone/mcp/actions/workflows/ci.yml/badge.svg)](https://github.com/Arkestone/mcp/actions/workflows/ci.yml)
[![codecov](https://codecov.io/gh/Arkestone/mcp/graph/badge.svg)](https://codecov.io/gh/Arkestone/mcp)
[![Go Report Card](https://goreportcard.com/badge/github.com/Arkestone/mcp)](https://goreportcard.com/report/github.com/Arkestone/mcp)
[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](LICENSE)
[![Go Version](https://img.shields.io/github/go-mod/go-version/Arkestone/mcp)](go.mod)
[![Release](https://img.shields.io/github/v/release/Arkestone/mcp)](https://github.com/Arkestone/mcp/releases/latest)

A suite of [Model Context Protocol](https://modelcontextprotocol.io) (MCP) servers for [GitHub Copilot](https://github.com/features/copilot) and other AI coding assistants. Each server dynamically serves a different type of context — custom instructions, skills, prompts, ADRs, persistent memory, or knowledge graphs — from local directories and GitHub repositories.

> **Built with the [Go MCP SDK](https://github.com/modelcontextprotocol/go-sdk)**. Browse the [MCP Registry](https://registry.modelcontextprotocol.io/) for the full ecosystem of MCP servers.

## Table of Contents

- [Available Servers](#available-servers)
- [Installation](#installation)
  - [go install](#go-install)
  - [Docker](#docker)
  - [Pre-built Binaries](#pre-built-binaries)
  - [Build from Source](#build-from-source)
- [Quick Start](#quick-start)
- [Transport Mechanisms](#transport-mechanisms)
- [MCP Client Configuration](#mcp-client-configuration)
  - [VS Code](#vs-code)
  - [Claude Desktop](#claude-desktop)
  - [Cursor](#cursor)
  - [Windsurf](#windsurf)
  - [Claude Code](#claude-code)
  - [JetBrains](#jetbrains)
  - [Zed](#zed)
  - [Docker Compose (all servers)](#docker-compose--all-servers)
- [LLM Optimization](#llm-optimization)
- [GitHub Authentication](#github-authentication-optional)
- [Network & Proxy](#network--proxy)
- [Architecture](#architecture)
- [Shared Packages](#shared-packages)
- [Development](#development)
- [Contributing](#contributing)

## Available Servers

| Server | Description | Port |
|--------|-------------|------|
| [mcp-instructions](./servers/mcp-instructions/) | Serves Copilot custom instruction files (`.github/copilot-instructions.md`, `.github/instructions/**`) from local dirs and GitHub repos | `:8080` |
| [mcp-skills](./servers/mcp-skills/) | Serves Copilot skills (`skills/*/SKILL.md`) with frontmatter metadata and reference bundles | `:8081` |
| [mcp-prompts](./servers/mcp-prompts/) | Serves VS Code Copilot prompt files (`.github/prompts/*.prompt.md`) and chat mode files | `:8082` |
| [mcp-adr](./servers/mcp-adr/) | Serves Architecture Decision Records from `docs/adr/`, `docs/decisions/`, or `doc/adr/` | `:8083` |
| [mcp-memory](./servers/mcp-memory/) | Persistent memory store — remember, recall, and forget information across sessions | `:8084` |
| [mcp-graph](./servers/mcp-graph/) | Knowledge graph — store entities and relationships, query neighbors and shortest paths | `:8085` |

Each server has its own [README](./servers/mcp-instructions/README.md) and [CHANGELOG](./servers/mcp-instructions/CHANGELOG.md).

## Installation

### go install

The fastest way to install. Requires [Go 1.24+](https://go.dev/dl/).

```bash
# Install the latest version
go install github.com/Arkestone/mcp/servers/mcp-instructions/cmd/mcp-instructions@latest
go install github.com/Arkestone/mcp/servers/mcp-skills/cmd/mcp-skills@latest
go install github.com/Arkestone/mcp/servers/mcp-prompts/cmd/mcp-prompts@latest
go install github.com/Arkestone/mcp/servers/mcp-adr/cmd/mcp-adr@latest
go install github.com/Arkestone/mcp/servers/mcp-memory/cmd/mcp-memory@latest
go install github.com/Arkestone/mcp/servers/mcp-graph/cmd/mcp-graph@latest

# Install a pinned version
go install github.com/Arkestone/mcp/servers/mcp-instructions/cmd/mcp-instructions@v0.0.1
```

### Docker

Pre-built images are published to the GitHub Container Registry after each release:

```bash
# Pull the latest version
docker pull ghcr.io/arkestone/mcp-instructions:latest
docker pull ghcr.io/arkestone/mcp-skills:latest
docker pull ghcr.io/arkestone/mcp-prompts:latest
docker pull ghcr.io/arkestone/mcp-adr:latest
docker pull ghcr.io/arkestone/mcp-memory:latest
docker pull ghcr.io/arkestone/mcp-graph:latest

# Pull a specific version
docker pull ghcr.io/arkestone/mcp-instructions:v0.0.1
```

### Pre-built Binaries

Download pre-built binaries for Linux, macOS, and Windows from [GitHub Releases](https://github.com/Arkestone/mcp/releases):

```bash
# Example: Linux amd64
curl -L https://github.com/Arkestone/mcp/releases/latest/download/mcp-instructions_linux_amd64.tar.gz | tar xz
mv mcp-instructions /usr/local/bin/
```

### Build from Source

```bash
git clone https://github.com/Arkestone/mcp.git
cd mcp
make build    # → ./bin/mcp-instructions, ./bin/mcp-skills, ...
```

## Quick Start

```bash
# stdio — used by VS Code, Claude Desktop, Cursor, etc.
mcp-instructions -dirs /path/to/repo

# HTTP — for shared team deployments
mcp-instructions -transport http -addr :8080 -repos github/awesome-copilot

# With LLM optimization
mcp-instructions -transport http -llm.enabled -llm.endpoint https://api.openai.com/v1
```

See each server's README for the full configuration reference.

## Transport Mechanisms

| Transport | When to Use |
|-----------|-------------|
| **stdio** | Local clients (VS Code, Claude Desktop, Cursor, etc.) — the client spawns the server process |
| **Streamable HTTP** | Remote/shared deployments — the server runs independently, clients connect over HTTP |

All servers default to stdio. Pass `-transport http` to switch to HTTP.

<!-- quick-install -->
## Quick Install in VS Code

Click to install any server directly into VS Code or VS Code Insiders. **stdio** requires the binary installed locally (`go install`); **HTTP** connects to a running local server; **Docker** runs the server in a container (no local install needed).

| Server | Transport | VS Code | VS Code Insiders |
|--------|-----------|---------|------------------|
| `mcp-instructions` | stdio  | [![VS Code](https://img.shields.io/badge/VS_Code-Install-0098FF?style=flat-square&logo=visualstudiocode&logoColor=white)](https://insiders.vscode.dev/redirect/mcp/install?name=mcp-instructions&config=%7B%22type%22%3A%20%22stdio%22%2C%20%22command%22%3A%20%22mcp-instructions%22%2C%20%22args%22%3A%20%5B%22--dirs%22%2C%20%22%24%7BworkspaceFolder%7D%22%5D%7D) | [![VS Code Insiders](https://img.shields.io/badge/VS_Code_Insiders-Install-24bfa5?style=flat-square&logo=visualstudiocode&logoColor=white)](https://insiders.vscode.dev/redirect/mcp/install?name=mcp-instructions&config=%7B%22type%22%3A%20%22stdio%22%2C%20%22command%22%3A%20%22mcp-instructions%22%2C%20%22args%22%3A%20%5B%22--dirs%22%2C%20%22%24%7BworkspaceFolder%7D%22%5D%7D&quality=insiders) |
| | HTTP   | [![VS Code](https://img.shields.io/badge/VS_Code-Install-0098FF?style=flat-square&logo=visualstudiocode&logoColor=white)](https://insiders.vscode.dev/redirect/mcp/install?name=mcp-instructions&config=%7B%22type%22%3A%20%22http%22%2C%20%22url%22%3A%20%22http%3A%2F%2Flocalhost%3A8080%2Fmcp%22%7D) | [![VS Code Insiders](https://img.shields.io/badge/VS_Code_Insiders-Install-24bfa5?style=flat-square&logo=visualstudiocode&logoColor=white)](https://insiders.vscode.dev/redirect/mcp/install?name=mcp-instructions&config=%7B%22type%22%3A%20%22http%22%2C%20%22url%22%3A%20%22http%3A%2F%2Flocalhost%3A8080%2Fmcp%22%7D&quality=insiders) |
| | Docker | [![VS Code](https://img.shields.io/badge/VS_Code-Install-0098FF?style=flat-square&logo=visualstudiocode&logoColor=white)](https://insiders.vscode.dev/redirect/mcp/install?name=mcp-instructions&config=%7B%22type%22%3A%20%22stdio%22%2C%20%22command%22%3A%20%22docker%22%2C%20%22args%22%3A%20%5B%22run%22%2C%20%22--rm%22%2C%20%22-i%22%2C%20%22-v%22%2C%20%22%24%7BworkspaceFolder%7D%3A%2Fworkspace%3Aro%22%2C%20%22ghcr.io%2Farkestone%2Fmcp-instructions%3Alatest%22%2C%20%22--dirs%22%2C%20%22%2Fworkspace%22%5D%7D) | [![VS Code Insiders](https://img.shields.io/badge/VS_Code_Insiders-Install-24bfa5?style=flat-square&logo=visualstudiocode&logoColor=white)](https://insiders.vscode.dev/redirect/mcp/install?name=mcp-instructions&config=%7B%22type%22%3A%20%22stdio%22%2C%20%22command%22%3A%20%22docker%22%2C%20%22args%22%3A%20%5B%22run%22%2C%20%22--rm%22%2C%20%22-i%22%2C%20%22-v%22%2C%20%22%24%7BworkspaceFolder%7D%3A%2Fworkspace%3Aro%22%2C%20%22ghcr.io%2Farkestone%2Fmcp-instructions%3Alatest%22%2C%20%22--dirs%22%2C%20%22%2Fworkspace%22%5D%7D&quality=insiders) |
| `mcp-skills` | stdio  | [![VS Code](https://img.shields.io/badge/VS_Code-Install-0098FF?style=flat-square&logo=visualstudiocode&logoColor=white)](https://insiders.vscode.dev/redirect/mcp/install?name=mcp-skills&config=%7B%22type%22%3A%20%22stdio%22%2C%20%22command%22%3A%20%22mcp-skills%22%2C%20%22args%22%3A%20%5B%22--dirs%22%2C%20%22%24%7BworkspaceFolder%7D%22%5D%7D) | [![VS Code Insiders](https://img.shields.io/badge/VS_Code_Insiders-Install-24bfa5?style=flat-square&logo=visualstudiocode&logoColor=white)](https://insiders.vscode.dev/redirect/mcp/install?name=mcp-skills&config=%7B%22type%22%3A%20%22stdio%22%2C%20%22command%22%3A%20%22mcp-skills%22%2C%20%22args%22%3A%20%5B%22--dirs%22%2C%20%22%24%7BworkspaceFolder%7D%22%5D%7D&quality=insiders) |
| | HTTP   | [![VS Code](https://img.shields.io/badge/VS_Code-Install-0098FF?style=flat-square&logo=visualstudiocode&logoColor=white)](https://insiders.vscode.dev/redirect/mcp/install?name=mcp-skills&config=%7B%22type%22%3A%20%22http%22%2C%20%22url%22%3A%20%22http%3A%2F%2Flocalhost%3A8081%2Fmcp%22%7D) | [![VS Code Insiders](https://img.shields.io/badge/VS_Code_Insiders-Install-24bfa5?style=flat-square&logo=visualstudiocode&logoColor=white)](https://insiders.vscode.dev/redirect/mcp/install?name=mcp-skills&config=%7B%22type%22%3A%20%22http%22%2C%20%22url%22%3A%20%22http%3A%2F%2Flocalhost%3A8081%2Fmcp%22%7D&quality=insiders) |
| | Docker | [![VS Code](https://img.shields.io/badge/VS_Code-Install-0098FF?style=flat-square&logo=visualstudiocode&logoColor=white)](https://insiders.vscode.dev/redirect/mcp/install?name=mcp-skills&config=%7B%22type%22%3A%20%22stdio%22%2C%20%22command%22%3A%20%22docker%22%2C%20%22args%22%3A%20%5B%22run%22%2C%20%22--rm%22%2C%20%22-i%22%2C%20%22-v%22%2C%20%22%24%7BworkspaceFolder%7D%3A%2Fworkspace%3Aro%22%2C%20%22ghcr.io%2Farkestone%2Fmcp-skills%3Alatest%22%2C%20%22--dirs%22%2C%20%22%2Fworkspace%22%5D%7D) | [![VS Code Insiders](https://img.shields.io/badge/VS_Code_Insiders-Install-24bfa5?style=flat-square&logo=visualstudiocode&logoColor=white)](https://insiders.vscode.dev/redirect/mcp/install?name=mcp-skills&config=%7B%22type%22%3A%20%22stdio%22%2C%20%22command%22%3A%20%22docker%22%2C%20%22args%22%3A%20%5B%22run%22%2C%20%22--rm%22%2C%20%22-i%22%2C%20%22-v%22%2C%20%22%24%7BworkspaceFolder%7D%3A%2Fworkspace%3Aro%22%2C%20%22ghcr.io%2Farkestone%2Fmcp-skills%3Alatest%22%2C%20%22--dirs%22%2C%20%22%2Fworkspace%22%5D%7D&quality=insiders) |
| `mcp-prompts` | stdio  | [![VS Code](https://img.shields.io/badge/VS_Code-Install-0098FF?style=flat-square&logo=visualstudiocode&logoColor=white)](https://insiders.vscode.dev/redirect/mcp/install?name=mcp-prompts&config=%7B%22type%22%3A%20%22stdio%22%2C%20%22command%22%3A%20%22mcp-prompts%22%2C%20%22args%22%3A%20%5B%22--dirs%22%2C%20%22%24%7BworkspaceFolder%7D%22%5D%7D) | [![VS Code Insiders](https://img.shields.io/badge/VS_Code_Insiders-Install-24bfa5?style=flat-square&logo=visualstudiocode&logoColor=white)](https://insiders.vscode.dev/redirect/mcp/install?name=mcp-prompts&config=%7B%22type%22%3A%20%22stdio%22%2C%20%22command%22%3A%20%22mcp-prompts%22%2C%20%22args%22%3A%20%5B%22--dirs%22%2C%20%22%24%7BworkspaceFolder%7D%22%5D%7D&quality=insiders) |
| | HTTP   | [![VS Code](https://img.shields.io/badge/VS_Code-Install-0098FF?style=flat-square&logo=visualstudiocode&logoColor=white)](https://insiders.vscode.dev/redirect/mcp/install?name=mcp-prompts&config=%7B%22type%22%3A%20%22http%22%2C%20%22url%22%3A%20%22http%3A%2F%2Flocalhost%3A8082%2Fmcp%22%7D) | [![VS Code Insiders](https://img.shields.io/badge/VS_Code_Insiders-Install-24bfa5?style=flat-square&logo=visualstudiocode&logoColor=white)](https://insiders.vscode.dev/redirect/mcp/install?name=mcp-prompts&config=%7B%22type%22%3A%20%22http%22%2C%20%22url%22%3A%20%22http%3A%2F%2Flocalhost%3A8082%2Fmcp%22%7D&quality=insiders) |
| | Docker | [![VS Code](https://img.shields.io/badge/VS_Code-Install-0098FF?style=flat-square&logo=visualstudiocode&logoColor=white)](https://insiders.vscode.dev/redirect/mcp/install?name=mcp-prompts&config=%7B%22type%22%3A%20%22stdio%22%2C%20%22command%22%3A%20%22docker%22%2C%20%22args%22%3A%20%5B%22run%22%2C%20%22--rm%22%2C%20%22-i%22%2C%20%22-v%22%2C%20%22%24%7BworkspaceFolder%7D%3A%2Fworkspace%3Aro%22%2C%20%22ghcr.io%2Farkestone%2Fmcp-prompts%3Alatest%22%2C%20%22--dirs%22%2C%20%22%2Fworkspace%22%5D%7D) | [![VS Code Insiders](https://img.shields.io/badge/VS_Code_Insiders-Install-24bfa5?style=flat-square&logo=visualstudiocode&logoColor=white)](https://insiders.vscode.dev/redirect/mcp/install?name=mcp-prompts&config=%7B%22type%22%3A%20%22stdio%22%2C%20%22command%22%3A%20%22docker%22%2C%20%22args%22%3A%20%5B%22run%22%2C%20%22--rm%22%2C%20%22-i%22%2C%20%22-v%22%2C%20%22%24%7BworkspaceFolder%7D%3A%2Fworkspace%3Aro%22%2C%20%22ghcr.io%2Farkestone%2Fmcp-prompts%3Alatest%22%2C%20%22--dirs%22%2C%20%22%2Fworkspace%22%5D%7D&quality=insiders) |
| `mcp-adr` | stdio  | [![VS Code](https://img.shields.io/badge/VS_Code-Install-0098FF?style=flat-square&logo=visualstudiocode&logoColor=white)](https://insiders.vscode.dev/redirect/mcp/install?name=mcp-adr&config=%7B%22type%22%3A%20%22stdio%22%2C%20%22command%22%3A%20%22mcp-adr%22%2C%20%22args%22%3A%20%5B%22--dirs%22%2C%20%22%24%7BworkspaceFolder%7D%22%5D%7D) | [![VS Code Insiders](https://img.shields.io/badge/VS_Code_Insiders-Install-24bfa5?style=flat-square&logo=visualstudiocode&logoColor=white)](https://insiders.vscode.dev/redirect/mcp/install?name=mcp-adr&config=%7B%22type%22%3A%20%22stdio%22%2C%20%22command%22%3A%20%22mcp-adr%22%2C%20%22args%22%3A%20%5B%22--dirs%22%2C%20%22%24%7BworkspaceFolder%7D%22%5D%7D&quality=insiders) |
| | HTTP   | [![VS Code](https://img.shields.io/badge/VS_Code-Install-0098FF?style=flat-square&logo=visualstudiocode&logoColor=white)](https://insiders.vscode.dev/redirect/mcp/install?name=mcp-adr&config=%7B%22type%22%3A%20%22http%22%2C%20%22url%22%3A%20%22http%3A%2F%2Flocalhost%3A8083%2Fmcp%22%7D) | [![VS Code Insiders](https://img.shields.io/badge/VS_Code_Insiders-Install-24bfa5?style=flat-square&logo=visualstudiocode&logoColor=white)](https://insiders.vscode.dev/redirect/mcp/install?name=mcp-adr&config=%7B%22type%22%3A%20%22http%22%2C%20%22url%22%3A%20%22http%3A%2F%2Flocalhost%3A8083%2Fmcp%22%7D&quality=insiders) |
| | Docker | [![VS Code](https://img.shields.io/badge/VS_Code-Install-0098FF?style=flat-square&logo=visualstudiocode&logoColor=white)](https://insiders.vscode.dev/redirect/mcp/install?name=mcp-adr&config=%7B%22type%22%3A%20%22stdio%22%2C%20%22command%22%3A%20%22docker%22%2C%20%22args%22%3A%20%5B%22run%22%2C%20%22--rm%22%2C%20%22-i%22%2C%20%22-v%22%2C%20%22%24%7BworkspaceFolder%7D%3A%2Fworkspace%3Aro%22%2C%20%22ghcr.io%2Farkestone%2Fmcp-adr%3Alatest%22%2C%20%22--dirs%22%2C%20%22%2Fworkspace%22%5D%7D) | [![VS Code Insiders](https://img.shields.io/badge/VS_Code_Insiders-Install-24bfa5?style=flat-square&logo=visualstudiocode&logoColor=white)](https://insiders.vscode.dev/redirect/mcp/install?name=mcp-adr&config=%7B%22type%22%3A%20%22stdio%22%2C%20%22command%22%3A%20%22docker%22%2C%20%22args%22%3A%20%5B%22run%22%2C%20%22--rm%22%2C%20%22-i%22%2C%20%22-v%22%2C%20%22%24%7BworkspaceFolder%7D%3A%2Fworkspace%3Aro%22%2C%20%22ghcr.io%2Farkestone%2Fmcp-adr%3Alatest%22%2C%20%22--dirs%22%2C%20%22%2Fworkspace%22%5D%7D&quality=insiders) |
| `mcp-memory` | stdio  | [![VS Code](https://img.shields.io/badge/VS_Code-Install-0098FF?style=flat-square&logo=visualstudiocode&logoColor=white)](https://insiders.vscode.dev/redirect/mcp/install?name=mcp-memory&config=%7B%22type%22%3A%20%22stdio%22%2C%20%22command%22%3A%20%22mcp-memory%22%7D) | [![VS Code Insiders](https://img.shields.io/badge/VS_Code_Insiders-Install-24bfa5?style=flat-square&logo=visualstudiocode&logoColor=white)](https://insiders.vscode.dev/redirect/mcp/install?name=mcp-memory&config=%7B%22type%22%3A%20%22stdio%22%2C%20%22command%22%3A%20%22mcp-memory%22%7D&quality=insiders) |
| | HTTP   | [![VS Code](https://img.shields.io/badge/VS_Code-Install-0098FF?style=flat-square&logo=visualstudiocode&logoColor=white)](https://insiders.vscode.dev/redirect/mcp/install?name=mcp-memory&config=%7B%22type%22%3A%20%22http%22%2C%20%22url%22%3A%20%22http%3A%2F%2Flocalhost%3A8084%2Fmcp%22%7D) | [![VS Code Insiders](https://img.shields.io/badge/VS_Code_Insiders-Install-24bfa5?style=flat-square&logo=visualstudiocode&logoColor=white)](https://insiders.vscode.dev/redirect/mcp/install?name=mcp-memory&config=%7B%22type%22%3A%20%22http%22%2C%20%22url%22%3A%20%22http%3A%2F%2Flocalhost%3A8084%2Fmcp%22%7D&quality=insiders) |
| | Docker | [![VS Code](https://img.shields.io/badge/VS_Code-Install-0098FF?style=flat-square&logo=visualstudiocode&logoColor=white)](https://insiders.vscode.dev/redirect/mcp/install?name=mcp-memory&config=%7B%22type%22%3A%20%22stdio%22%2C%20%22command%22%3A%20%22docker%22%2C%20%22args%22%3A%20%5B%22run%22%2C%20%22--rm%22%2C%20%22-i%22%2C%20%22-v%22%2C%20%22mcp-memory%3A%2Fdata%22%2C%20%22ghcr.io%2Farkestone%2Fmcp-memory%3Alatest%22%5D%7D) | [![VS Code Insiders](https://img.shields.io/badge/VS_Code_Insiders-Install-24bfa5?style=flat-square&logo=visualstudiocode&logoColor=white)](https://insiders.vscode.dev/redirect/mcp/install?name=mcp-memory&config=%7B%22type%22%3A%20%22stdio%22%2C%20%22command%22%3A%20%22docker%22%2C%20%22args%22%3A%20%5B%22run%22%2C%20%22--rm%22%2C%20%22-i%22%2C%20%22-v%22%2C%20%22mcp-memory%3A%2Fdata%22%2C%20%22ghcr.io%2Farkestone%2Fmcp-memory%3Alatest%22%5D%7D&quality=insiders) |
| `mcp-graph` | stdio  | [![VS Code](https://img.shields.io/badge/VS_Code-Install-0098FF?style=flat-square&logo=visualstudiocode&logoColor=white)](https://insiders.vscode.dev/redirect/mcp/install?name=mcp-graph&config=%7B%22type%22%3A%20%22stdio%22%2C%20%22command%22%3A%20%22mcp-graph%22%7D) | [![VS Code Insiders](https://img.shields.io/badge/VS_Code_Insiders-Install-24bfa5?style=flat-square&logo=visualstudiocode&logoColor=white)](https://insiders.vscode.dev/redirect/mcp/install?name=mcp-graph&config=%7B%22type%22%3A%20%22stdio%22%2C%20%22command%22%3A%20%22mcp-graph%22%7D&quality=insiders) |
| | HTTP   | [![VS Code](https://img.shields.io/badge/VS_Code-Install-0098FF?style=flat-square&logo=visualstudiocode&logoColor=white)](https://insiders.vscode.dev/redirect/mcp/install?name=mcp-graph&config=%7B%22type%22%3A%20%22http%22%2C%20%22url%22%3A%20%22http%3A%2F%2Flocalhost%3A8085%2Fmcp%22%7D) | [![VS Code Insiders](https://img.shields.io/badge/VS_Code_Insiders-Install-24bfa5?style=flat-square&logo=visualstudiocode&logoColor=white)](https://insiders.vscode.dev/redirect/mcp/install?name=mcp-graph&config=%7B%22type%22%3A%20%22http%22%2C%20%22url%22%3A%20%22http%3A%2F%2Flocalhost%3A8085%2Fmcp%22%7D&quality=insiders) |
| | Docker | [![VS Code](https://img.shields.io/badge/VS_Code-Install-0098FF?style=flat-square&logo=visualstudiocode&logoColor=white)](https://insiders.vscode.dev/redirect/mcp/install?name=mcp-graph&config=%7B%22type%22%3A%20%22stdio%22%2C%20%22command%22%3A%20%22docker%22%2C%20%22args%22%3A%20%5B%22run%22%2C%20%22--rm%22%2C%20%22-i%22%2C%20%22-v%22%2C%20%22mcp-graph%3A%2Fdata%22%2C%20%22ghcr.io%2Farkestone%2Fmcp-graph%3Alatest%22%5D%7D) | [![VS Code Insiders](https://img.shields.io/badge/VS_Code_Insiders-Install-24bfa5?style=flat-square&logo=visualstudiocode&logoColor=white)](https://insiders.vscode.dev/redirect/mcp/install?name=mcp-graph&config=%7B%22type%22%3A%20%22stdio%22%2C%20%22command%22%3A%20%22docker%22%2C%20%22args%22%3A%20%5B%22run%22%2C%20%22--rm%22%2C%20%22-i%22%2C%20%22-v%22%2C%20%22mcp-graph%3A%2Fdata%22%2C%20%22ghcr.io%2Farkestone%2Fmcp-graph%3Alatest%22%5D%7D&quality=insiders) |
<!-- /quick-install -->

## MCP Client Configuration

### VS Code

> **Quick install**: Use the badges in [Quick Install in VS Code](#quick-install-in-vs-code) above, or open the Command Palette (`Ctrl+Shift+P` / `Cmd+Shift+P`) and run **`MCP: Open User Configuration`** to add a server to your global `mcp.json`. See the [VS Code MCP documentation](https://code.visualstudio.com/docs/copilot/model-context-protocol) for full details.

#### `.vscode/mcp.json` (project-level)

```json
{
  "servers": {
    "instructions": { "command": "mcp-instructions", "args": ["-dirs", "${workspaceFolder}"] },
    "skills":       { "command": "mcp-skills",       "args": ["-dirs", "${workspaceFolder}"] },
    "prompts":      { "command": "mcp-prompts",      "args": ["-dirs", "${workspaceFolder}"] },
    "adrs":         { "command": "mcp-adr",          "args": ["-dirs", "${workspaceFolder}"] },
    "memory":       { "command": "mcp-memory" },
    "graph":        { "command": "mcp-graph" }
  }
}
```

#### Global user settings (`settings.json`)

```json
{
  "github.copilot.chat.mcp.servers": {
    "instructions": {
      "command": "mcp-instructions",
      "args": ["-dirs", "/path/to/repo"]
    }
  }
}
```

### Claude Desktop

**macOS**: `~/Library/Application Support/Claude/claude_desktop_config.json`
**Windows**: `%APPDATA%\Claude\claude_desktop_config.json`

```json
{
  "mcpServers": {
    "instructions": {
      "command": "mcp-instructions",
      "args": ["-dirs", "/path/to/repo"]
    },
    "skills": {
      "command": "mcp-skills",
      "args": ["-dirs", "/path/to/repo"]
    },
    "memory": {
      "command": "mcp-memory",
      "env": { "MEMORY_DIR": "~/.local/share/mcp-memory" }
    },
    "graph": {
      "command": "mcp-graph"
    }
  }
}
```

### Cursor

#### `.cursor/mcp.json` (project-level)

```json
{
  "mcpServers": {
    "instructions": {
      "command": "mcp-instructions",
      "args": ["-dirs", "."]
    },
    "skills": {
      "command": "mcp-skills",
      "args": ["-dirs", "."]
    },
    "prompts": {
      "command": "mcp-prompts",
      "args": ["-dirs", "."]
    },
    "memory": {
      "command": "mcp-memory"
    }
  }
}
```

#### Global (`~/.cursor/mcp.json`)

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

**macOS/Linux**: `~/.codeium/windsurf/mcp_config.json`
**Windows**: `%USERPROFILE%\.codeium\windsurf\mcp_config.json`

```json
{
  "mcpServers": {
    "instructions": {
      "command": "mcp-instructions",
      "args": ["-dirs", "/path/to/repo"]
    },
    "memory": {
      "command": "mcp-memory",
      "env": { "MEMORY_DIR": "~/.local/share/mcp-memory" }
    }
  }
}
```

### Claude Code

`~/.mcp.json` (global) or `.mcp.json` (project-level)

```json
{
  "mcpServers": {
    "instructions": {
      "command": "mcp-instructions",
      "args": ["-dirs", "."]
    },
    "skills": {
      "command": "mcp-skills",
      "args": ["-dirs", "."]
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

### JetBrains

In JetBrains IDEs (IntelliJ IDEA, GoLand, PyCharm, etc.), go to **Settings → Tools → AI Assistant → Model Context Protocol (MCP)**:

```json
{
  "mcpServers": {
    "instructions": {
      "command": "mcp-instructions",
      "args": ["-dirs", "$PROJECT_DIR$"]
    },
    "memory": {
      "command": "mcp-memory"
    }
  }
}
```

### Zed

Add to `~/.config/zed/settings.json`:

```json
{
  "context_servers": {
    "mcp-instructions": {
      "command": {
        "path": "mcp-instructions",
        "args": ["-dirs", "/path/to/repo"]
      }
    },
    "mcp-memory": {
      "command": {
        "path": "mcp-memory",
        "args": []
      }
    }
  }
}
```

### Docker Compose — all servers

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
      GRAPH_DATA_FILE: /data/graph.json
    volumes: ["graph-data:/data"]

volumes:
  memory-data:
  graph-data:
```

Then connect your client over HTTP:

```json
{
  "mcpServers": {
    "instructions": {
      "type": "http",
      "url": "http://localhost:8080/mcp"
    }
  }
}
```

## LLM Optimization

All content servers (instructions, skills, prompts, ADRs) optionally consolidate multiple sources using an OpenAI-compatible LLM endpoint. This deduplicates and logically organizes content from multiple repositories.

```bash
mcp-instructions \
  -dirs . \
  -llm.enabled \
  -llm.endpoint https://api.openai.com/v1 \
  -llm.api-key $LLM_API_KEY \
  -llm.model gpt-4o-mini
```

**Supported providers**: OpenAI, Azure OpenAI, Ollama, LM Studio, and any OpenAI-compatible endpoint.

See [`config.llm.example.yaml`](./servers/mcp-instructions/config.llm.example.yaml) for a complete configuration example.

## GitHub Authentication (Optional)

A GitHub token is **optional**. Public repositories work without authentication. For private repositories, provide a token (highest priority first):

| Method | Example |
|--------|---------|
| CLI flag | `-github-token ghp_xxx` |
| Prefixed env var | `INSTRUCTIONS_GITHUB_TOKEN=ghp_xxx` |
| Global env var | `GITHUB_TOKEN=ghp_xxx` |
| YAML config | `github_token: ghp_xxx` |

## Network & Proxy

All servers work on-premise, in private/public cloud, with direct internet or through HTTP/HTTPS proxies. See the **[Network & Proxy Guide](docs/network.md)** for firewall rules, proxy configuration, and custom CA certificates.

## Architecture

```
.
├── servers/
│   ├── mcp-instructions/   # custom instructions server  (:8080)
│   ├── mcp-skills/         # skills server               (:8081)
│   ├── mcp-prompts/        # prompt files server         (:8082)
│   ├── mcp-adr/            # ADR server                  (:8083)
│   ├── mcp-memory/         # persistent memory server    (:8084)
│   └── mcp-graph/          # knowledge graph server      (:8085)
├── pkg/
│   ├── config/             # shared configuration loading (YAML → env → flags)
│   ├── github/             # GitHub Contents API client
│   ├── httputil/           # proxy, TLS, header propagation
│   ├── optimizer/          # shared LLM optimization layer (OpenAI-compatible)
│   ├── server/             # MCP server bootstrap helpers
│   └── syncer/             # background repo sync
├── docs/
│   └── network.md          # network / proxy / firewall guide
├── examples/               # client configuration examples
└── AGENTS.md               # AI coding assistant guide
```

Each content server follows the same layered design:

1. **Config** — YAML → environment variables → CLI flags (each layer overrides the previous)
2. **Loader / Scanner** — discovers content from local directories and GitHub repositories
3. **Optimizer** — optional LLM-based consolidation via `pkg/optimizer`
4. **MCP Server** — exposes content as Resources, Prompts, and Tools over stdio or Streamable HTTP

## Shared Packages

| Package | Description |
|---------|-------------|
| `pkg/config` | Unified config loading: YAML → env vars → CLI flags |
| `pkg/github` | GitHub Contents API client with proxy and header pass-through |
| `pkg/httputil` | Proxy support, custom TLS/CA certificates, header propagation |
| `pkg/optimizer` | OpenAI-compatible LLM client for optional content consolidation |
| `pkg/server` | MCP server bootstrap and `/healthz` endpoint helpers |
| `pkg/syncer` | Background periodic sync for remote GitHub repositories |

## Development

```bash
make build              # build all servers into ./bin/
make test               # run unit tests
make test-integration   # run integration tests (requires LLM_ENDPOINT)
make docker             # build Docker images for all servers
make lint               # run golangci-lint
make cover              # generate coverage report
```

## Contributing

Contributions are welcome! Please see [CONTRIBUTING.md](CONTRIBUTING.md) for guidelines.

## License

This project is licensed under the MIT License — see the [LICENSE](LICENSE) file for details.
