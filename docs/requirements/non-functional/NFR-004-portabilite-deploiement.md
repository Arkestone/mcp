# NFR-004 — Portabilité de déploiement (déjà exemplaire)

**Statut** : Confirmé — déjà implémenté et vérifiable dans le dépôt réel

Chaque serveur est distribuable par `go install`, image Docker publiée (`ghcr.io/arkestone/<serveur>`), binaire précompilé multi-OS (Linux/macOS/Windows), et orchestrable via docker-compose (tous les serveurs ensemble) — fonctionnant on-premise, en cloud privé ou public, avec ou sans proxy.

## Relation avec le standard portefeuille

Ce NFR ne fait que constater ce que `Shared/architecture/deployment-portability.md` demande à tous les projets du portefeuille — Teikyō est, à la lecture du dépôt réel, l'un des rares projets à l'avoir déjà pleinement atteint avant même que ce standard soit formalisé. À citer comme référence dans l'ADR d'adoption plutôt qu'à retraiter comme un manque à combler.

## Ce qui manque encore pour Kubernetes

`[QUESTION OUVERTE]` — aucun manifeste Kubernetes/Helm chart n'a été trouvé dans le dépôt (seul docker-compose est documenté) ; l'image Docker existante le permettrait techniquement, mais l'artefact d'orchestration K8s reste à produire si ce besoin se confirme.
