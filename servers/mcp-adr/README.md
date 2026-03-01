# mcp-adr

An MCP server that serves Architecture Decision Records (ADRs) from local directories and GitHub repositories. It scans `docs/adr/`, `docs/decisions/`, and `doc/adr/` within each configured source and exposes ADRs via the Model Context Protocol.

## Quick Start

```bash
# Build
go build -o mcp-adr ./servers/mcp-adr/cmd/mcp-adr

# Run with stdio transport (default)
./mcp-adr

# Run with HTTP transport
ADR_TRANSPORT=http ADR_ADDR=:8083 ./mcp-adr
```

## Configuration

Configuration is loaded from a YAML file, environment variables, and CLI flags (later overrides earlier). Copy `config.example.yaml` and adjust:

```yaml
sources:
  dirs:
    - /path/to/local/repo       # scanned for docs/adr/*.md

  repos:
    - myorg/my-project          # GitHub repos (owner/repo or owner/repo@ref)

cache:
  dir: ~/.cache/mcp-adr
  sync_interval: 5m

llm:
  endpoint: ""                  # OpenAI-compatible endpoint
  model: gpt-4o-mini
  enabled: false

transport: stdio                # stdio | http
addr: ":8083"
```

### Environment variables

All variables are prefixed with `ADR_`:

| Variable | Description |
|---|---|
| `ADR_CONFIG` | Path to config YAML file |
| `ADR_TRANSPORT` | `stdio` or `http` |
| `ADR_ADDR` | HTTP listen address |
| `ADR_GITHUB_TOKEN` | GitHub personal access token |
| `ADR_LLM_ENDPOINT` | LLM API endpoint |
| `ADR_LLM_MODEL` | LLM model name |
| `ADR_LLM_ENABLED` | Enable LLM optimization (`true`/`false`) |

## MCP Primitives

### Resources

| URI | Description |
|---|---|
| `adrs://{source}/{id}` | Individual ADR content (Markdown + frontmatter) |
| `adrs://optimized` | All ADRs merged via LLM (or concatenated) |
| `adrs://index` | Plain-text list of all available ADRs |

### Prompts

| Name | Arguments | Description |
|---|---|---|
| `get-adrs` | `source` (opt), `status` (opt), `optimize` (opt) | Get ADRs, optionally filtered by source or status |

### Tools

| Name | Description |
|---|---|
| `refresh-adrs` | Force re-sync of all ADR sources |
| `list-adrs` | List all ADRs with optional `source` and `status` filters |
| `get-adr` | Get a single ADR by `uri` |
| `optimize-adrs` | Get consolidated ADR content with optional LLM optimization |

## ADR Format

Each ADR is a Markdown file with optional YAML frontmatter:

```markdown
---
title: Use PostgreSQL
status: accepted
date: 2023-06-01
---

## Context

We need a relational database...

## Decision

We will use PostgreSQL.
```

Supported statuses: `proposed`, `accepted`, `deprecated`, `superseded`.

If `title` is not present in the frontmatter, it is derived from the filename (e.g. `0001-use-postgresql` → `0001 Use Postgresql`).

## Scanned Directories

Within each configured source directory or GitHub repository root, the server looks for ADRs in:

- `docs/adr/`
- `docs/decisions/`
- `doc/adr/`

## License

MIT
