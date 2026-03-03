# awesome-copilot preset

Load the community instruction library from [github/awesome-copilot](https://github.com/github/awesome-copilot) — 100+ curated Copilot instruction files — alongside your own workspace instructions.

## VS Code — `.vscode/mcp.json`

### Local binary (stdio)

```json
{
  "servers": {
    "instructions": {
      "command": "mcp-instructions",
      "args": [
        "-dirs", "${workspaceFolder}",
        "-repos", "github/awesome-copilot"
      ]
    }
  }
}
```

### Docker

```json
{
  "servers": {
    "instructions": {
      "command": "docker",
      "args": [
        "run", "--rm", "-i",
        "-v", "${workspaceFolder}:/workspace:ro",
        "ghcr.io/arkestone/mcp-instructions:latest",
        "-dirs", "/workspace",
        "-repos", "github/awesome-copilot"
      ]
    }
  }
}
```

### HTTP (already running server)

```json
{
  "servers": {
    "instructions": {
      "type": "http",
      "url": "http://localhost:3000"
    }
  }
}
```

---

## Claude Desktop — `claude_desktop_config.json`

> **macOS:** `~/Library/Application Support/Claude/claude_desktop_config.json`  
> **Windows:** `%APPDATA%\Claude\claude_desktop_config.json`

```json
{
  "mcpServers": {
    "instructions": {
      "command": "mcp-instructions",
      "args": [
        "-dirs", "/path/to/your/project",
        "-repos", "github/awesome-copilot"
      ]
    }
  }
}
```

---

## With a GitHub token (rate-limit protection)

The server fetches `github/awesome-copilot` via the public GitHub API. For heavy usage or CI, supply a token to avoid rate limits:

```json
{
  "servers": {
    "instructions": {
      "command": "mcp-instructions",
      "args": [
        "-dirs", "${workspaceFolder}",
        "-repos", "github/awesome-copilot",
        "-github-token", "${env:GITHUB_TOKEN}"
      ]
    }
  }
}
```

---

## Pin to a specific commit

```json
{
  "servers": {
    "instructions": {
      "command": "mcp-instructions",
      "args": [
        "-dirs", "${workspaceFolder}",
        "-repos", "github/awesome-copilot@beb33e6"
      ]
    }
  }
}
```

---

## What you get

Once connected, Copilot can access all instructions from `github/awesome-copilot` as MCP resources, for example:

| URI | Description |
|-----|-------------|
| `instructions://github/awesome-copilot/golang` | Go best practices |
| `instructions://github/awesome-copilot/typescript` | TypeScript conventions |
| `instructions://github/awesome-copilot/react` | React patterns |
| `instructions://github/awesome-copilot/python` | Python style guide |
| `instructions://github/awesome-copilot/agents` | Copilot agent authoring |
| `instructions://github/awesome-copilot/security` | Security guidelines |
| …and 100+ more | |

> The repository is cached locally and re-synced every **30 minutes** by default.  
> Override with `-sync-interval 5m` for fresher content or `-sync-interval 24h` to reduce network calls.
