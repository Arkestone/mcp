# NFR-001 — Interopérabilité via le protocole MCP standard

**Statut** : Draft

Chaque serveur Teikyō expose son contenu exclusivement via les primitives standard du [Model Context Protocol](https://modelcontextprotocol.io) (Resources, Tools, Prompts), construites avec le Go MCP SDK — jamais via un protocole propriétaire Arkestone.

## Justification

À la différence des autres projets du portefeuille, qui versionnent leur propre API REST (voir `Shared/architecture/product-platform-standards.md`), l'interopérabilité de Teikyō repose sur la conformité à un standard externe déjà versionné par l'écosystème MCP lui-même — c'est cette conformité, pas un schéma interne, qui garantit qu'un serveur Teikyō reste utilisable par n'importe quel client MCP présent ou futur (VS Code, Claude Desktop, Cursor, Windsurf, JetBrains, Zed, et tout client à venir).

## Critère de vérification

Given un nouveau client MCP conforme au protocole standard, When il se connecte à un serveur Teikyō, Then il peut consommer ses Resources/Tools/Prompts sans adaptation spécifique à Arkestone.
