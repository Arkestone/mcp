# Client Configuration Examples

This directory contains ready-to-use MCP client configuration snippets for both servers.

## VS Code (`.vscode/mcp.json`)

```json
{
  "servers": {
    "mcp-instructions": {
      "type": "http",
      "url": "http://localhost:8080/mcp"
    },
    "mcp-skills": {
      "type": "http",
      "url": "http://localhost:8081/mcp"
    }
  }
}
```

## Claude Desktop (`~/Library/Application Support/Claude/claude_desktop_config.json`)

```json
{
  "mcpServers": {
    "mcp-instructions": {
      "command": "mcp-instructions",
      "args": ["--http", "--addr", "127.0.0.1:8080"],
      "env": {
        "INSTRUCTIONS_DIRS": "/path/to/instructions"
      }
    },
    "mcp-skills": {
      "command": "mcp-skills",
      "args": ["--http", "--addr", "127.0.0.1:8081"],
      "env": {
        "SKILLS_DIRS": "/path/to/skills"
      }
    }
  }
}
```

## Docker Compose

```yaml
services:
  mcp-instructions:
    image: ghcr.io/arkestone/mcp-instructions:latest
    ports:
      - "8080:8080"
    command: ["--http", "--addr", "0.0.0.0:8080"]
    volumes:
      - ./instructions:/instructions:ro
    environment:
      INSTRUCTIONS_DIRS: /instructions

  mcp-skills:
    image: ghcr.io/arkestone/mcp-skills:latest
    ports:
      - "8081:8081"
    command: ["--http", "--addr", "0.0.0.0:8081"]
    volumes:
      - ./skills:/skills:ro
    environment:
      SKILLS_DIRS: /skills
```

## Stdio (direct process, for local use)

```bash
# Instructions server via stdio
mcp-instructions --config config.yaml

# Skills server via stdio
mcp-skills --config config.yaml
```
