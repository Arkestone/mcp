# NFR-002 — Disponibilité en déploiement partagé

**Statut** : Draft

Tout serveur déployé en mode Streamable HTTP (déploiement partagé pour une équipe) expose un endpoint `/healthz` (déjà implémenté dans `pkg/server`) permettant à un orchestrateur de superviser sa disponibilité indépendamment de l'état d'une requête MCP en cours.

## Justification

Condition nécessaire à l'usage de Teikyō sous Kubernetes/docker-compose avec probes de santé (voir `Shared/architecture/deployment-portability.md`) — sans cet endpoint, un orchestrateur ne peut pas distinguer un serveur bloqué d'un serveur simplement occupé.

## Cible

`[QUESTION OUVERTE]` — aucun SLA chiffré (taux de disponibilité mensuel) n'a été trouvé dans le dépôt réel ; à établir si un déploiement HTTP partagé devient critique pour une équipe consommatrice.
