# Changelog — mcp-instructions

All notable changes to the **mcp-instructions** server are documented here.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.1.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [Unreleased]

## [0.1.0] - 2025-07-07

### Added

- MCP server with stdio and HTTP (Streamable HTTP) transports
- On-demand instruction loading from local directories
- GitHub repository caching with configurable sync interval
- Optional LLM-based instruction optimization (OpenAI-compatible)
- MCP Resources, Prompts, and Tools primitives
- YAML configuration with environment variable and CLI flag overrides
- Docker support with multi-stage builds
- `/healthz` endpoint for liveness checks in HTTP mode
- Graceful HTTP shutdown (5 s drain window)
- Deterministic `Sources []string` ordering in refresh output

[Unreleased]: https://github.com/Arkestone/mcp/compare/mcp-instructions/v0.1.0...HEAD
[0.1.0]: https://github.com/Arkestone/mcp/releases/tag/mcp-instructions/v0.1.0
