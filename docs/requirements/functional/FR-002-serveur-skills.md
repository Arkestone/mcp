# FR-002 — Servir les skills Copilot

**Statut** : Draft — décomposé à partir du produit réel (`servers/mcp-skills`)

En tant qu'assistant IA de développement, je veux accéder aux skills définis pour un projet (`skills/*/SKILL.md`, avec leurs métadonnées frontmatter et bundles de référence), afin de réutiliser des procédures déjà validées plutôt que de réinventer une méthode à chaque session.

## Critères d'acceptation

- Given un répertoire `skills/<nom>/SKILL.md` avec un en-tête frontmatter (nom, description), When le serveur `mcp-skills` scanne ce répertoire, Then le skill est exposé avec ses métadonnées consultables séparément de son contenu complet.
- Given un skill accompagné d'un bundle de fichiers de référence additionnels, When le client demande le skill, Then ces fichiers de référence restent accessibles en complément du fichier `SKILL.md` principal.
- Given plusieurs répertoires `skills/` provenant de sources différentes (local + GitHub), When le serveur les agrège, Then aucun conflit de nom de skill n'est résolu silencieusement — un skill dupliqué doit être signalé plutôt que masqué.

## Dépendances

- [FR-007](./FR-007-double-transport.md), [FR-008](./FR-008-optimisation-llm.md), [FR-010](./FR-010-configuration-en-cascade.md) — mêmes mécanismes transverses que FR-001.
