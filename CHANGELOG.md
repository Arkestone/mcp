# Changelog

This monorepo uses per-server changelogs. See:

- [`servers/mcp-instructions/CHANGELOG.md`](servers/mcp-instructions/CHANGELOG.md)
- [`servers/mcp-skills/CHANGELOG.md`](servers/mcp-skills/CHANGELOG.md)

## Monorepo-level changes

### [Unreleased]

### [1.0.0] - 2026-03-02

#### Added

- 6 MCP servers: `mcp-instructions`, `mcp-skills`, `mcp-prompts`, `mcp-adr`, `mcp-memory`, `mcp-graph`
- VS Code one-click install badges (stdio, HTTP, Docker) for all servers
- Preset: [Awesome Copilot](https://github.com/github/awesome-copilot) integration for `mcp-instructions`
- Docusaurus documentation site at <https://arkestone.github.io/mcp/>
- Multi-arch Docker images published to GHCR (`linux/amd64`, `linux/arm64`)
- Binary releases for Linux, macOS, Windows (amd64 + arm64)
- Shared packages: `pkg/config`, `pkg/github`, `pkg/httputil`, `pkg/optimizer`, `pkg/server`, `pkg/syncer`, `pkg/testutil`
- Mutation testing with gremlins, security scanning with Trivy + CodeQL

### [0.1.0] - 2025-07-07

#### Added

- `servers/` layout with `mcp-instructions` and `mcp-skills`
- Shared `pkg/` packages: config, github, httputil, optimizer, server, syncer, testutil
- Monorepo CI/CD: multi-arch Docker builds, GoReleaser, Trivy security scans
- Comprehensive test suite (100+ tests)
- Dev Container for Codespaces (`.devcontainer/`)
- Client configuration examples (`examples/`)
- AGENTS.md for AI coding assistant guidance

[unreleased]: https://github.com/Arkestone/mcp/compare/v1.0.0...HEAD
[1.0.0]: https://github.com/Arkestone/mcp/releases/tag/v1.0.0
[0.1.0]: https://github.com/Arkestone/mcp/releases/tag/v0.1.0
