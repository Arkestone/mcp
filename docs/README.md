# Teikyō (提供)

Suite de six serveurs Model Context Protocol (MCP) partageant une infrastructure commune, servant dynamiquement du contexte organisationnel — instructions personnalisées, skills, prompts, ADR, mémoire persistante, graphe de connaissance — à GitHub Copilot et aux autres assistants IA de développement.

## Pitch

Un assistant IA de développement n'a, par défaut, aucun accès structuré au contexte propre à une équipe : ses instructions personnalisées, ses procédures réutilisables, ses prompts standardisés, ses décisions d'architecture passées, ce qu'il a appris d'une session à l'autre, ou les relations entre les entités de son domaine. Ce contexte existe déjà — dispersé dans des fichiers locaux ou des dépôts GitHub (`.github/copilot-instructions.md`, `skills/*/SKILL.md`, `docs/adr/`) — mais sans mécanisme standard pour qu'un assistant le découvre à la demande plutôt que de le faire recopier manuellement à chaque nouveau projet.

Teikyō répond à ce besoin avec six serveurs MCP indépendants — `mcp-instructions`, `mcp-skills`, `mcp-prompts`, `mcp-adr`, `mcp-memory` (mémoire persistante remember/recall/forget) et `mcp-graph` (graphe de connaissance avec recherche de plus court chemin) — partageant la même infrastructure : configuration en cascade YAML→variables d'environnement→CLI, double transport stdio/HTTP, découverte GitHub, optimisation LLM optionnelle, synchronisation en arrière-plan. Chaque nouveau serveur réutilise cette base plutôt que d'être redéveloppé isolément.

Le nom porte cette fonction sans détour : *teikyō* (提供) signifie « fourniture, mise à disposition » — exactement ce que fait un serveur qui sert du contexte à la demande.

## Vision

L'ambition de Teikyō est qu'un assistant IA, sur n'importe quelle machine et pour n'importe quel projet d'une équipe, retrouve instantanément la même mémoire organisationnelle — sans copier-coller de fichiers de contexte entre outils ou entre postes.

Teikyō est déjà, avant même la formalisation des standards transverses du portefeuille cette session, un exemple abouti de portabilité de déploiement : `go install`, image Docker par serveur, binaires précompilés, `docker-compose` multi-serveurs, hébergement on-premise ou cloud avec proxy — tout ce que `deployment-portability.md` demande aux autres projets est déjà en production ici. C'est une référence à citer pour les autres projets du portefeuille plutôt qu'un standard à leur imposer depuis Teikyō.

À terme, l'infrastructure déjà partagée entre les six serveurs devrait rendre l'ajout d'un septième serveur — pourquoi pas un jour un serveur exposant les FR/NFR de ce portefeuille de documentation lui-même à un assistant IA — significativement moins coûteux qu'un développement isolé, sans qu'aucun septième serveur n'ait encore été construit pour le vérifier en pratique.

## En bref

- **Statut** : produit réel actif (Go, 30+ commits). Documentation de portefeuille amorcée le 2026-08-31 : 10 exigences fonctionnelles (une par serveur/capacité réelle) et 7 non fonctionnelles, décomposées à partir du code source réel (`pkg/`, `servers/`), pas d'une recherche externe générique.
- **Dérogations explicites** : accessibilité et micro-frontends non pertinents (outil CLI/serveur sans interface utilisateur) — documentées comme telles plutôt que forcées.
- **Points ouverts signalés** : comportement de repli si l'endpoint LLM d'optimisation est injoignable non documenté ; aucun SLA de disponibilité ni seuil de performance chiffré ; aucun manifeste Kubernetes/Helm malgré une image Docker déjà prête.
- **Par où commencer** : [`requirements/business-case.md`](./requirements/business-case.md), puis le README racine du dépôt pour la liste complète des serveurs disponibles.

## Documentation

| Dossier | Contenu | Statut |
|---|---|---|
| `network.md` | Documentation technique réseau du produit réel | Source de vérité produit |
| [`requirements/`](./requirements/README.md) | Business case, parties prenantes, glossaire, état actuel, FR/NFR atomiques — décomposition du produit réel selon la convention du portefeuille Arkestone | Issu du mono-dépôt de documentation [`Arkestone/projects`](https://github.com/Arkestone/projects) (sous-projet Teikyō) |
| [`adr/`](./adr/README.md) | Adoption des standards transverses du portefeuille pertinents pour un outil CLI/serveur | Idem |
