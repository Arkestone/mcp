# Changelog — mcp-adr

## [Unreleased]

## [0.1.0] - 2026-03-01

### Added

- MCP server serving Architecture Decision Records from `docs/adr/*.md`
- Scans `docs/adr/`, `docs/decisions/`, and `doc/adr/` within configured directories
- YAML frontmatter parsing: title, status, date
- Resources: `adrs://{source}/{id}`, `adrs://optimized`, `adrs://index`
- Prompts: `get-adrs` with source and status filtering
- Tools: `refresh-adrs`, `list-adrs`, `get-adr`, `optimize-adrs`
- Status filtering (proposed, accepted, deprecated, superseded)
- stdio and HTTP transports

[Unreleased]: https://github.com/Arkestone/mcp/compare/mcp-adr/v0.1.0...HEAD
[0.1.0]: https://github.com/Arkestone/mcp/releases/tag/mcp-adr/v0.1.0
