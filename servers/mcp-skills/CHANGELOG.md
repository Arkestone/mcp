# Changelog — mcp-skills

All notable changes to the **mcp-skills** server are documented here.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.1.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [Unreleased]

## [0.1.0] - 2025-07-07

### Added

- MCP server with stdio and HTTP (Streamable HTTP) transports
- Skills scanner with SKILL.md frontmatter parsing
- Skills-specific MCP resources, prompts, and tools
- YAML configuration with environment variable and CLI flag overrides
- Docker support with multi-stage builds
- `/healthz` endpoint for liveness checks in HTTP mode
- Graceful HTTP shutdown (5 s drain window)
- Deterministic `Sources []string` ordering in refresh output

[Unreleased]: https://github.com/Arkestone/mcp/compare/mcp-skills/v0.1.0...HEAD
[0.1.0]: https://github.com/Arkestone/mcp/releases/tag/mcp-skills/v0.1.0
