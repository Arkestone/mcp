# Changelog — mcp-memory

All notable changes to the **mcp-memory** server are documented here.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.1.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [Unreleased]

## [0.1.0] - 2025-03-01

### Added

- MCP server with stdio and HTTP (Streamable HTTP) transports
- Persistent on-disk memory storage (plain Markdown files)
- `remember` tool — store a new memory with optional tags
- `recall` tool — full-text and tag-based memory search
- `forget` tool — delete a memory by ID
- `list-memories` tool — list memories with optional tag filter
- MCP Resources: `memory://{id}` (single memory) and `memory://all` (all memories)
- YAML configuration with environment variable and CLI flag overrides
- Docker support with multi-stage builds
- `/healthz` endpoint for liveness checks in HTTP mode
- Graceful HTTP shutdown (5 s drain window)

[Unreleased]: https://github.com/Arkestone/mcp/compare/v0.0.1...HEAD
[0.1.0]: https://github.com/Arkestone/mcp/releases/tag/v0.0.1
