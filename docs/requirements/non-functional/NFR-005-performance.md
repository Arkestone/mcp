# NFR-005 — Performance de service du contexte

**Statut** : Draft — question ouverte, aucun chiffre trouvé dans le dépôt réel

Un serveur Teikyō doit répondre à une requête MCP (Resource, Tool, Prompt) sans latence perceptible pour l'utilisateur d'un assistant IA interactif, y compris lorsque le contenu provient d'un dépôt GitHub distant.

## Comment cette exigence est déjà partiellement satisfaite

La synchronisation périodique en arrière-plan ([FR-009](../functional/FR-009-synchronisation-arriere-plan.md)) évite qu'une requête client déclenche un appel réseau synchrone vers GitHub — la latence perçue dépend donc principalement du scan local et, si activée, du temps de réponse de l'optimiseur LLM ([FR-008](../functional/FR-008-optimisation-llm.md)).

## Cible

`[QUESTION OUVERTE]` — aucun seuil p95 chiffré n'a été trouvé dans le README ou les tests d'intégration (`make test-integration`) ; à établir avec l'équipe qui maintient le dépôt réel si un besoin de SLA se confirme.
