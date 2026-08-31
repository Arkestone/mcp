# ADR-001 — Adoption des standards transverses du portefeuille Arkestone

**Statut** : Accepted

## Contexte

Le portefeuille Arkestone a formalisé 6 standards d'architecture transverses dans `Shared/architecture/` du mono-dépôt de documentation : DDD/BDD/TDD, carte de contextes, patterns d'architecture applicative, portail/API/micro-frontends/i18n/a11y/design system, NFR transverses, et portabilité de déploiement.

Teikyō a une nature différente des autres projets du portefeuille : c'est un outillage backend/CLI pour développeurs (six serveurs MCP), pas une application métier avec portail utilisateur final. Chaque standard doit donc être évalué pour sa pertinence réelle plutôt qu'adopté mécaniquement.

## Décision

- [DDD/BDD/TDD](https://github.com/Arkestone/projects/blob/main/Shared/architecture/ddd-bdd-tdd.md) — **adopté** : chaque serveur (`mcp-instructions`, `mcp-skills`, `mcp-prompts`, `mcp-adr`, `mcp-memory`, `mcp-graph`) constitue un bounded context indépendant, déjà organisé en modules Go distincts (`servers/*/`). Les critères Given/When/Then des FR de ce document sont directement transposables en scénarios Gherkin.
- [Patterns d'architecture applicative](https://github.com/Arkestone/projects/blob/main/Shared/architecture/application-architecture-patterns.md) — **cohérent avec l'existant** : le design en couches déjà en place (Config → Loader/Scanner → Optimiseur → Serveur MCP) suit le même principe que Ports & Adapters (l'optimiseur LLM et le client GitHub sont des adaptateurs remplaçables) sans qu'une réécriture soit nécessaire.
- [Portail/API/micro-frontends/i18n/a11y/design system](https://github.com/Arkestone/projects/blob/main/Shared/architecture/product-platform-standards.md) — **partiellement adopté** : l'interopérabilité via API versionnée est assurée nativement par le protocole MCP lui-même (voir NFR-001) ; micro-frontends, i18n et a11y ne s'appliquent pas, Teikyō n'exposant aucune UI (dérogations explicites : NFR-006, NFR-007).
- [NFR transverses](https://github.com/Arkestone/projects/blob/main/Shared/architecture/nfr-transverses.md) — **partiellement adopté** : pas de NFR de disponibilité chiffrée ni de multi-fuseau horaire pertinents pour un outil sans notion de délai métier ou d'échéance humaine.
- [Portabilité de déploiement](https://github.com/Arkestone/projects/blob/main/Shared/architecture/deployment-portability.md) — **déjà pleinement atteint** (voir NFR-004) : Teikyō est une référence à citer pour ce standard plutôt qu'un projet ayant un rattrapage à faire.
- [Carte de contextes](https://github.com/Arkestone/projects/blob/main/Shared/architecture/context-map.md) — Teikyō n'a, à ce jour, aucune relation de dépendance documentée avec un autre projet du portefeuille métier (Sekisho, Meikan, etc.) : c'est un outil d'outillage IA consommé par des développeurs, pas un service métier consommé par une autre application. À mettre à jour si une intégration réelle apparaît (ex. un projet du portefeuille qui utiliserait `mcp-adr` pour exposer ses propres ADR à un assistant IA).

## Conséquences

- Les futures FR/NFR de Teikyō doivent continuer à distinguer explicitement ce qui s'applique (DDD, architecture en couches, déploiement) de ce qui ne s'applique pas (a11y, i18n, micro-frontends) plutôt que d'appliquer les standards portefeuille mécaniquement.
- Toute UI d'administration future (non identifiée dans le dépôt réel actuel) devrait déclencher une révision de cette ADR et des NFR-006/NFR-007.
