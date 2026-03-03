# Changelog

This monorepo uses per-server changelogs. See:

- [`servers/mcp-instructions/CHANGELOG.md`](servers/mcp-instructions/CHANGELOG.md)
- [`servers/mcp-skills/CHANGELOG.md`](servers/mcp-skills/CHANGELOG.md)

## Monorepo-level changes

### [Unreleased]

### [1.2.0] - 2026-03-03

#### Added

- ZIP fallback for GitHub repository downloads: when the Contents API returns a rate-limit or auth error, all 3 servers (instructions, prompts, skills) automatically fall back to downloading the full repo as a ZIP archive — a single request that bypasses the 60 req/hr unauthenticated API limit
- `pkg/github`: `FetchZipAndExtract(ctx, owner, repo, ref, targetDir)` — downloads the GitHub zipball and extracts it locally, stripping the top-level directory prefix GitHub adds
- `pkg/github`: `IsRateLimitError(err)` — helper to detect rate-limit errors from GitHub API responses

#### Fixed

- Without `GITHUB_TOKEN`, large repos (e.g. `github/awesome-copilot` with 175+ files) would exhaust the 60 req/hr rate limit and fail silently; the ZIP fallback recovers automatically

### [1.1.0] - 2026-03-02

#### Added

- Context-aware filtering: `get-context` tool returns instructions/prompts/skills dynamically ranked by relevance to the current file and query context
- `files:` frontmatter field for prompts and skills — restricts items to matching file globs (same as `applyTo:` for instructions)
- Filter engine improvements: phrase adjacency bonus, short-tag exact match, name precision ratio bonus, stopwords, suffix-stripping stemmer, co-occurrence bonus

#### Fixed

- Brace expansion in glob patterns (`**/*.{ts,tsx}`) no longer split incorrectly at the comma inside braces
- `tags:` field as a plain comma-separated string now parses correctly (was silently dropped when YAML expected `[]string`)
- Invalid glob pattern in `files:` / `applyTo:` now uses inclusive fallback (item shown rather than hidden)
- Missing `[]interface{}` case in `toStringSlice` for YAML list-style `files:` in prompts loader

#### Tests

- Comprehensive agent scenario integration tests for all three servers (instructions, prompts, skills) covering file-context routing, query scoring, hidden directories, brace expansion, live reload, determinism, and adversarial frontmatter

### [1.0.1] - 2026-03-02

#### Changed

- CI: merged 20 workflow files into 9 consolidated workflows (39 → 29 files)
- CI: `release.yml` — GHCR visibility now set with `GITHUB_TOKEN` (`packages: write`); `PACKAGES_TOKEN` used as fallback
- CI: added `ghcr-publish.yml` — standalone dispatchable workflow to publish GHCR packages

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

[unreleased]: https://github.com/Arkestone/mcp/compare/v1.0.1...HEAD
[1.0.1]: https://github.com/Arkestone/mcp/compare/v1.0.0...v1.0.1
[1.0.0]: https://github.com/Arkestone/mcp/releases/tag/v1.0.0
[0.1.0]: https://github.com/Arkestone/mcp/releases/tag/v0.1.0
