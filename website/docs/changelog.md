---
title: Changelog
sidebar_position: 4
---


This monorepo uses per-server changelogs. See:

- [`servers/mcp-instructions/CHANGELOG.md`](https://github.com/Arkestone/mcp/blob/main/servers/mcp-instructions/CHANGELOG.md)
- [`servers/mcp-skills/CHANGELOG.md`](https://github.com/Arkestone/mcp/blob/main/servers/mcp-skills/CHANGELOG.md)

## Monorepo-level changes

### [Unreleased]

### [0.1.0] - 2025-07-07

#### Added

- `servers/` layout with `mcp-instructions` and `mcp-skills`
- Shared `pkg/` packages: config, github, httputil, optimizer, server, syncer, testutil
- Monorepo CI/CD: multi-arch Docker builds, GoReleaser, Trivy security scans
- Comprehensive test suite (100+ tests)
- Dev Container for Codespaces (`.devcontainer/`)
- Client configuration examples (`examples/`)
- AGENTS.md for AI coding assistant guidance

[unreleased]: https://github.com/Arkestone/mcp/compare/v0.1.0...HEAD
[0.1.0]: https://github.com/Arkestone/mcp/releases/tag/v0.1.0
