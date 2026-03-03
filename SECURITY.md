# Security Policy

## Reporting a Vulnerability

**Please do not report security vulnerabilities through public GitHub issues.**

Instead, please use [GitHub Security Advisories](https://github.com/github/github-mcp-server/security/advisories/new)
to report vulnerabilities privately. This ensures that security issues are handled
responsibly and disclosed only after a fix is available.

When reporting, please include as much of the following as possible:

- A description of the vulnerability and its potential impact
- Steps to reproduce the issue
- Affected version(s)
- Any relevant logs, screenshots, or proof-of-concept code

## Scope

The following areas are in scope for security reports:

- **MCP server code** — the core server implementation and transport layers
- **Configuration handling** — parsing, validation, and storage of server configuration
- **GitHub API interactions** — authentication, authorization, and data handling with the GitHub API
- **LLM proxy** — request forwarding, prompt injection defenses, and API key management

Out of scope:

- Issues in third-party dependencies (report these upstream)
- Social engineering attacks
- Denial-of-service attacks that require excessive resource consumption

## Response Timeline

| Action                          | Timeframe       |
| ------------------------------- | --------------- |
| Acknowledge receipt of report   | Within 48 hours |
| Provide an initial assessment   | Within 7 days   |
| Release a fix (critical issues) | Best effort     |

We will keep you informed of our progress throughout the process.

## Supported Versions

Only the **latest release** is actively supported with security updates.
We strongly recommend that users always run the most recent version.

| Version        | Supported |
| -------------- | --------- |
| Latest release | ✅         |
| Older releases | ❌         |

## Security Best Practices

When using this project, follow these guidelines to protect your environment:

- **Never commit secrets** — Do not commit tokens, API keys, or credentials to source control.
- **Use environment variables** — Pass sensitive values via `GITHUB_TOKEN` and `LLM_API_KEY` environment variables rather than command-line arguments or config files.
- **Restrict token scopes** — Grant only the minimum required permissions to any access tokens.
- **Rotate credentials regularly** — Periodically rotate tokens and API keys.
- **Review configuration** — Audit your `config.yaml` to ensure no sensitive data is stored in plain text.
- **Pin dependencies** — Use dependency pinning and verify checksums to prevent supply-chain attacks.

## Preferred Languages

We prefer all communications to be in English.
