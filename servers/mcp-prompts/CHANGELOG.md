# Changelog — mcp-prompts

## [Unreleased]

## [1.2.0] - 2026-03-03

### Added

- ZIP fallback: when the GitHub Contents API is rate-limited or returns a 403/auth error, the server automatically downloads the full repository as a single ZIP archive.

## [1.1.0] - 2026-03-02

### Added

- `files:` frontmatter field — restricts a prompt to matching file globs (e.g. `**/*.ts`, `**/*.{ts,tsx}`)
- Context-aware filtering: `get-context` tool ranks prompts by relevance to current file and query
- Filter engine: phrase adjacency bonus, short-tag exact match, name precision ratio bonus, stopwords, suffix-stripping stemmer

### Fixed

- Brace expansion in `files:` glob patterns (`**/*.{ts,tsx}`) no longer split incorrectly at the comma
- `tags:` field as a plain comma-separated string now parses correctly (was silently dropped)
- Invalid glob in `files:` uses inclusive fallback
- Missing `[]interface{}` case in YAML list-style `files:` field

### Tests

- Comprehensive agent scenario integration tests: `TestAgentScenario_Prompts_*` covering file-context routing, query scoring, combined filters, hidden directories, tag formats, brace expansion, live reload, determinism, and adversarial frontmatter

## [0.1.0] - 2026-03-01

### Added

- MCP server serving `.github/prompts/*.prompt.md` and `.github/chatmodes/*.chatmode.md`
- Resources: `prompts://{source}/{name}`, `prompts://optimized`, `prompts://index`
- Prompts: `get-prompts`
- Tools: `refresh-prompts`, `list-prompts`, `get-prompt`, `optimize-prompts`
- stdio and HTTP (Streamable HTTP) transports
- Optional LLM-based optimization
- `/healthz` endpoint in HTTP mode
- Graceful HTTP shutdown

[Unreleased]: https://github.com/Arkestone/mcp/compare/mcp-prompts/v0.1.0...HEAD
[0.1.0]: https://github.com/Arkestone/mcp/releases/tag/mcp-prompts/v0.1.0
