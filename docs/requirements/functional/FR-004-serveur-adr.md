# FR-004 — Servir les décisions d'architecture (ADR)

**Statut** : Draft — décomposé à partir du produit réel (`servers/mcp-adr`)

En tant qu'assistant IA de développement, je veux accéder aux ADR d'un projet (`docs/adr/`, `docs/decisions/`, ou `doc/adr/`), afin de proposer des solutions cohérentes avec les décisions d'architecture déjà prises plutôt que de les ignorer ou les contredire.

## Critères d'acceptation

- Given un projet utilisant l'une des trois conventions de répertoire supportées (`docs/adr/`, `docs/decisions/`, `doc/adr/`), When le serveur `mcp-adr` scanne le projet, Then les ADR sont découvertes quelle que soit la convention utilisée, sans configuration supplémentaire.
- Given plusieurs ADR dont une remplace explicitement une autre (statut « Remplacé par »), When le client consulte l'historique, Then la relation de remplacement reste visible plutôt que de faire disparaître l'ADR obsolète sans trace.

## Dépendances

- [FR-007](./FR-007-double-transport.md), [FR-008](./FR-008-optimisation-llm.md), [FR-010](./FR-010-configuration-en-cascade.md) — mêmes mécanismes transverses que FR-001.

## Note

Le portefeuille de documentation Arkestone (`Shared/architecture/`, et l'ADR de chaque projet) est lui-même un consommateur potentiel de ce serveur une fois les projets extraits en dépôts séparés avec leur propre `docs/adr/`.
