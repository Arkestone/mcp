# Business Case — Teikyō

**Énoncé du problème** : les assistants IA de développement (GitHub Copilot, Claude Code, Cursor, Windsurf, JetBrains, Zed) n'ont, par défaut, aucun accès structuré au contexte propre à une équipe ou une organisation — ses instructions personnalisées, ses procédures réutilisables (skills), ses prompts standardisés, ses décisions d'architecture passées, ce qu'un assistant a appris d'une session à l'autre, ou les relations entre entités de son domaine. Ce contexte existe déjà sous forme de fichiers (`.github/copilot-instructions.md`, `skills/*/SKILL.md`, `.github/prompts/*.prompt.md`, `docs/adr/`) dispersés dans des dépôts locaux ou GitHub, sans mécanisme standard pour qu'un assistant les découvre à la demande.

Le vrai dépôt applicatif Arkestone MCP Servers (30 commits, actif) a déjà répondu à ce besoin en construisant six serveurs MCP indépendants — `mcp-instructions`, `mcp-skills`, `mcp-prompts`, `mcp-adr`, `mcp-memory`, `mcp-graph` — partageant une infrastructure commune (`pkg/config`, `pkg/github`, `pkg/httputil`, `pkg/optimizer`, `pkg/server`, `pkg/syncer`). Ce document décompose ce produit réel en exigences atomiques suivant la convention du portefeuille de documentation Arkestone, sans inventer de périmètre au-delà de ce qui est réellement implémenté.

## Objectifs métier

| Objectif | Baseline (état actuel) | Cible | Statut |
|---|---|---|---|
| Donner à tout assistant IA de développement un accès uniforme au contexte organisationnel | 6 serveurs MCP indépendants déjà en production, chacun avec son propre mécanisme de découverte avant ce projet | Les 6 serveurs partagent la même infrastructure (config, découverte GitHub, transport, optimisation LLM) — déjà atteint côté code (`pkg/`) | Confirmé, déjà implémenté |
| Réduire l'effort d'ajout d'un nouveau type de contexte à servir | Chaque nouveau serveur réutilise `pkg/server`, `pkg/config`, `pkg/github` plutôt que d'être développé isolément | `[QUESTION OUVERTE]` — aucun 7e serveur n'a encore été ajouté pour valider ce gain en pratique | Architecture en place, gain non mesuré |
| Rendre chaque serveur déployable dans n'importe quel environnement sans adaptation | go install, image Docker par serveur, binaires précompilés, docker-compose multi-serveurs, on-premise/cloud avec proxy — déjà documenté et fonctionnel | Cible déjà atteinte, à formaliser comme référence pour le standard `deployment-portability.md` du portefeuille | Confirmé, déjà implémenté |

## Valeur

Réduction du travail de configuration manuelle et de copier-coller de contexte entre projets et outils IA ; cohérence garantie entre les six types de contexte servis (même transport, même mécanisme de configuration, même politique d'authentification GitHub) ; et référence interne réutilisable pour toute équipe du portefeuille Arkestone souhaitant équiper ses assistants IA sans réinventer un serveur de contexte spécifique.

## Ce que ce business case ne couvre pas

- Aucune mesure d'adoption réelle (nombre d'équipes utilisant effectivement ces serveurs) n'a été trouvée dans le dépôt — à confirmer séparément.
- Le contenu servi par chaque serveur (quelles instructions, quels skills, quelle mémoire) reste la responsabilité de chaque équipe consommatrice — Teikyō sert le mécanisme, pas le contenu métier.
