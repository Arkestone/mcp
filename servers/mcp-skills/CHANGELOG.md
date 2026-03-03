# Changelog — mcp-skills

All notable changes to the **mcp-skills** server are documented here.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.1.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [Unreleased]

## [1.2.0] - 2026-03-03

### Added

- ZIP fallback: when the GitHub Contents API is rate-limited or returns a 403/auth error, the server automatically downloads the full repository as a single ZIP archive.

## [1.1.0] - 2026-03-02

### Added

- `files:` frontmatter field — restricts a skill to matching file globs (e.g. `**/*.go`, `**/*.{ts,tsx}`)
- Context-aware filtering: `get-context` tool ranks skills by relevance to current file and query
- Filter engine: phrase adjacency bonus, short-tag exact match, name precision ratio bonus, stopwords, suffix-stripping stemmer

### Fixed

- Brace expansion in `files:` glob patterns (`**/*.{ts,tsx}`) no longer split incorrectly at the comma
- `tags:` field as a plain comma-separated string now parses correctly (was silently dropped)
- Invalid glob in `files:` uses inclusive fallback

### Tests

- Comprehensive agent scenario integration tests: `TestAgentScenario_Skills_*` covering file-context routing, query scoring, combined filters, hidden directories, tag formats, brace expansion, live reload, determinism, and adversarial frontmatter

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
