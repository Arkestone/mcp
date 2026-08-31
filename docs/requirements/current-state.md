# État actuel (as-is) — Teikyō

À la différence des autres projets du portefeuille dont l'as-is décrit un processus métier externe non outillé, l'as-is de Teikyō **est déjà le produit cible** : six serveurs MCP réellement implémentés et actifs (dépôt Go, 30 commits), ce document en fait l'inventaire factuel plutôt que de constater un manque.

## Inventaire des six serveurs réels

| Serveur | Port par défaut | Contexte servi | Mécanisme de découverte |
|---|---|---|---|
| `mcp-instructions` | 8080 | Instructions personnalisées Copilot (`.github/copilot-instructions.md`, `.github/instructions/**`) | Répertoires locaux + dépôts GitHub |
| `mcp-skills` | 8081 | Skills Copilot (`skills/*/SKILL.md`) avec métadonnées frontmatter et bundles de référence | Répertoires locaux + dépôts GitHub |
| `mcp-prompts` | 8082 | Fichiers de prompt VS Code Copilot (`.github/prompts/*.prompt.md`) et fichiers de chat mode | Répertoires locaux + dépôts GitHub |
| `mcp-adr` | 8083 | Architecture Decision Records (`docs/adr/`, `docs/decisions/`, `doc/adr/`) | Répertoires locaux + dépôts GitHub |
| `mcp-memory` | 8084 | Mémoire persistante (remember/recall/forget), stockage Markdown sur disque | Fichiers locaux gérés par le serveur lui-même |
| `mcp-graph` | 8085 | Graphe de connaissance (nœuds/arêtes), persistance JSON atomique | Fichiers locaux gérés par le serveur lui-même |

## Infrastructure partagée déjà en place

- `pkg/config` — configuration en cascade YAML → env → flags, commune aux 6 serveurs.
- `pkg/github` — client GitHub Contents API avec support proxy et propagation d'en-têtes.
- `pkg/httputil` — proxy, TLS/CA personnalisés, propagation d'en-têtes.
- `pkg/optimizer` — client LLM compatible OpenAI pour consolidation optionnelle du contenu.
- `pkg/server` — bootstrap MCP commun + endpoint `/healthz`.
- `pkg/syncer` — synchronisation périodique en arrière-plan des dépôts GitHub distants.

Les quatre premiers serveurs (instructions, skills, prompts, adr) suivent un design en couches identique : Config → Loader/Scanner → Optimiseur (optionnel) → Serveur MCP. Les deux derniers (memory, graph) sont stateful et gèrent leur propre persistance locale plutôt que de scanner des répertoires externes — une différence architecturale réelle à ne pas gommer dans les FR.

## Contraintes et hypothèses

- **Contrainte héritée** : ce dépôt de documentation ne contient aucun code applicatif (`CLAUDE.md` racine) — Teikyō documente un produit dont l'implémentation vit exclusivement dans le dépôt applicatif réel.
- **Authentification GitHub optionnelle** : un token GitHub n'est nécessaire que pour accéder à des dépôts privés (4 sources de configuration par ordre de priorité : flag CLI, variable d'environnement préfixée, variable d'environnement globale, YAML).
- **Portabilité de déploiement déjà exemplaire** : go install, image Docker par serveur, binaires précompilés multi-OS, docker-compose (tous les serveurs), fonctionnement on-premise/cloud privé/public avec proxy — déjà conforme, voire en avance, sur le standard `deployment-portability.md` du portefeuille.
- **Hors périmètre déclaré** : Teikyō ne stocke ni ne décide du contenu métier servi (quelles instructions, quels ADR) — cette responsabilité reste celle de chaque équipe consommatrice.
