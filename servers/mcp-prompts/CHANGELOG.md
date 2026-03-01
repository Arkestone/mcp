# Changelog — mcp-prompts

## [Unreleased]

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
