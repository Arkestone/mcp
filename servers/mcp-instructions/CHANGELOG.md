# Changelog — mcp-instructions

All notable changes to the **mcp-instructions** server are documented here.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.1.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [Unreleased]

## [1.2.0] - 2026-03-03

### Added

- ZIP fallback: when the GitHub Contents API is rate-limited or returns a 403/auth error, the server automatically downloads the full repository as a single ZIP archive instead. This allows syncing large repos like `github/awesome-copilot` (175+ files) without a `GITHUB_TOKEN`.

## [1.1.0] - 2026-03-02

### Added

- Context-aware filtering: `get-context` tool returns instructions dynamically ranked by relevance to the current file path and query
- Efficient result caching with configurable TTL for filtered results

### Fixed

- Brace expansion in `applyTo:` glob patterns (`**/*.{ts,tsx}`) no longer split incorrectly at the comma inside braces
- `tags:` field as a plain comma-separated string now parses correctly
- Invalid glob pattern in `applyTo:` uses inclusive fallback (item shown rather than hidden due to a typo)

### Tests

- Comprehensive agent scenario integration tests: `TestAgentScenario_Instructions_*` covering 17 file-context cases, query scoring, hidden directories, non-standard layouts, brace expansion, live reload, determinism, and adversarial frontmatter

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
