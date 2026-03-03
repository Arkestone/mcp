# Contributing to mcp-instructions

Thank you for your interest in contributing to **mcp-instructions** — an MCP server that dynamically serves GitHub Copilot custom instructions from local directories and GitHub repositories.

We welcome contributions of all kinds: bug fixes, new features, documentation improvements, and more. Check out our [open issues](../../issues) for things to work on — issues labeled **good first issue** are a great starting point.

## Prerequisites

Before you begin, make sure you have the following installed:

- **Go 1.24+** — [download](https://go.dev/dl/)
- **Git** — [download](https://git-scm.com/)
- **Docker** (optional) — only needed if you want to build or test the container image

Verify your Go version:

```bash
go version
```

## Getting Started

1. **Clone the repository:**

   ```bash
   git clone https://github.com/Arkestone/mcp.git
   cd instructions
   ```

2. **Build the project:**

   ```bash
   go build ./...
   ```

3. **Run the tests:**

   ```bash
   go test ./...
   ```

4. **Run the server locally:**

   ```bash
   go run ./cmd/mcp-instructions -dirs /path/to/repo
   ```

See the [README](README.md) for full configuration options including environment variables and CLI flags.

## Project Structure

```
├── cmd/
│   └── mcp-instructions/    # Application entrypoint
├── internal/
│   ├── config/              # Configuration loading (YAML, env vars, CLI flags)
│   ├── loader/              # Instruction discovery and loading (local dirs, GitHub repos)
│   └── optimizer/           # LLM-based instruction merging and deduplication
├── config.example.yaml      # Example configuration file
├── Dockerfile               # Container image definition
└── go.mod                   # Go module definition
```

- **`cmd/mcp-instructions`** — Main entry point. Parses flags, initializes the config, and starts the MCP server.
- **`internal/config`** — Handles layered configuration: YAML files, environment variables, and CLI flags.
- **`internal/loader`** — Discovers `.github/copilot-instructions.md` and `.github/instructions/**/*.instructions.md` files from local directories and GitHub repositories. Manages caching and background sync for remote sources.
- **`internal/optimizer`** — Optionally consolidates instructions from multiple sources using an OpenAI-compatible LLM endpoint.

## Development Workflow

1. **Fork** the repository on GitHub.
2. **Create a branch** from `main` for your change:

   ```bash
   git checkout -b my-feature
   ```

3. **Make your changes** — keep commits focused and atomic.
4. **Run tests and linters** before committing (see [Testing](#testing) and [Code Style](#code-style)).
5. **Push** your branch and open a **Pull Request** against `main`.

## Code Style

We follow standard Go conventions:

- Run `go fmt ./...` before committing — all code must be formatted.
- Run `go vet ./...` to catch common issues.
- Use clear, descriptive names for functions, types, and variables.
- Keep exported API surfaces small; prefer `internal/` packages for implementation details.
- Aim for meaningful test coverage on new code. If you add a feature, add tests for it.

## Testing

### Unit Tests

```bash
go test ./...
```

### Race Detection

Run tests with the race detector to catch concurrency issues:

```bash
go test -race ./...
```

### Integration Tests

Integration tests (e.g., tests that call external services or require Docker) are gated behind a build tag:

```bash
go test -tags integration ./...
```

### Coverage

To generate a coverage report:

```bash
go test -coverprofile=coverage.out ./...
go tool cover -html=coverage.out
```

## Submitting Changes

When you open a pull request, please ensure:

1. **All tests pass** — `go test ./...` completes without failures.
2. **Code is formatted** — `go fmt ./...` produces no changes.
3. **No vet warnings** — `go vet ./...` is clean.
4. **PR description is clear** — describe *what* changed and *why*. Link to related issues if applicable.
5. **Commits are clean** — squash fixup commits; each commit should represent a logical change.
6. **Code has been reviewed** — at least one maintainer approval is required before merging.

## Reporting Issues

Found a bug or have an idea? [Open an issue](../../issues/new).

- **Bug reports** — include steps to reproduce, expected behavior, actual behavior, and your Go version / OS.
- **Feature requests** — describe the use case and the behavior you'd like to see.

Please search existing issues before opening a new one to avoid duplicates.

## License

This project is licensed under the [MIT License](LICENSE). By contributing, you agree that your contributions will be licensed under the same license.

No Contributor License Agreement (CLA) is required.
