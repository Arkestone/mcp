# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.1.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [Unreleased]

## [0.1.0] - 2025-07-07

### Added

- Skills MCP server (`mcp-skills`) for serving Copilot skills
- Monorepo structure supporting multiple MCP servers
- Shared `pkg/optimizer` package for LLM optimization
- Skills scanner with SKILL.md frontmatter parsing
- Skills-specific MCP resources, prompts, and tools
- MCP server with stdio and HTTP (Streamable HTTP) transports
- On-demand instruction loading from local directories
- GitHub repository caching with configurable sync interval
- Optional LLM-based instruction optimization (OpenAI-compatible)
- MCP Resources, Prompts, and Tools primitives
- YAML configuration with environment variable and CLI flag overrides
- Docker support with multi-stage builds
- Comprehensive test suite (100+ tests, race-safe)

### Changed

- Moved instructions server to `instructions/` subdirectory
- Extracted LLMConfig to shared `pkg/optimizer` package

[unreleased]: https://github.com/Arkestone/mcp/compare/v0.1.0...HEAD
[0.1.0]: https://github.com/Arkestone/mcp/releases/tag/v0.1.0
