# FR-001 — Servir les instructions personnalisées Copilot

**Statut** : Draft — décomposé à partir du produit réel (`servers/mcp-instructions`)

En tant qu'assistant IA de développement, je veux accéder aux fichiers d'instructions personnalisées d'un projet (`.github/copilot-instructions.md`, `.github/instructions/**`), afin d'appliquer les conventions propres à ce projet sans qu'un développeur les copie-colle manuellement dans chaque outil.

## Critères d'acceptation

- Given un répertoire local contenant `.github/copilot-instructions.md`, When le serveur `mcp-instructions` est démarré avec ce répertoire (`-dirs`), Then le contenu du fichier est exposé comme Resource MCP consultable par le client.
- Given un dépôt GitHub public ou privé (avec token, voir NFR-003), When le serveur est configuré avec ce dépôt (`-repos`), Then les fichiers d'instructions y sont découverts et servis au même titre qu'un répertoire local.
- Given un fichier sous `.github/instructions/**`, When le serveur scanne le répertoire, Then ce fichier est découvert au même titre que le fichier `copilot-instructions.md` racine.

## Dépendances

- [FR-007](./FR-007-double-transport.md) — servi via stdio ou HTTP au choix du client.
- [FR-008](./FR-008-optimisation-llm.md) — consolidation optionnelle du contenu avant service.
- [FR-010](./FR-010-configuration-en-cascade.md) — configuration du serveur (répertoires, dépôts, port).
