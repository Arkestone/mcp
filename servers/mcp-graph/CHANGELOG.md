# Changelog — mcp-graph

All notable changes to the `mcp-graph` MCP server are documented here.

## [Unreleased]

## [0.1.0] — 2026-03-01

### Added
- Initial implementation of `mcp-graph` knowledge graph MCP server
- In-memory graph with atomic JSON persistence (`graph.json`)
- Node model: label, name, optional key-value properties
- Edge model: directed relationship with type and optional properties
- Tools: `add-node`, `add-edge`, `find-nodes`, `get-node`, `neighbors`,
  `shortest-path`, `remove-node`, `remove-edge`, `list-relations`
- Resources: `graph://stats`, `graph://node/{id}`
- BFS-based `shortest-path` with configurable `max_depth`
- `neighbors` with `direction` (`out`/`in`/`both`) and `relation` filter
- `stdio` and HTTP (Streamable HTTP) transports
- Dockerfile and Dockerfile.goreleaser for multi-arch container images
- Port `:8085`
