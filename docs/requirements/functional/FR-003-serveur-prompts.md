# FR-003 — Servir les fichiers de prompt et chat modes

**Statut** : Draft — décomposé à partir du produit réel (`servers/mcp-prompts`)

En tant qu'assistant IA de développement, je veux accéder aux fichiers de prompt standardisés d'un projet (`.github/prompts/*.prompt.md`) et à ses fichiers de chat mode, afin de proposer les mêmes invites réutilisables que celles déjà définies par l'équipe pour VS Code Copilot.

## Critères d'acceptation

- Given un fichier `.github/prompts/nom.prompt.md`, When le serveur `mcp-prompts` scanne le répertoire, Then ce prompt est exposé comme primitive MCP Prompt, invocable par le client avec ses paramètres éventuels.
- Given un fichier de chat mode associé, When le client le demande, Then son contenu est servi distinctement d'un prompt standard.

## Dépendances

- [FR-007](./FR-007-double-transport.md), [FR-008](./FR-008-optimisation-llm.md), [FR-010](./FR-010-configuration-en-cascade.md) — mêmes mécanismes transverses que FR-001.
