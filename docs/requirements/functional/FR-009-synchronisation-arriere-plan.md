# FR-009 — Synchronisation périodique des dépôts distants

**Statut** : Draft — décomposé à partir du produit réel (`pkg/syncer`)

En tant qu'utilisateur d'un serveur configuré avec un ou plusieurs dépôts GitHub distants, je veux que leur contenu soit rafraîchi périodiquement en arrière-plan, afin qu'une requête client n'attende jamais un appel réseau complet vers GitHub à chaque consultation.

## Critères d'acceptation

- Given un serveur configuré avec `-repos`, When il démarre, Then une synchronisation initiale du contenu distant a lieu avant de commencer à répondre aux requêtes clients.
- Given une synchronisation déjà effectuée, When l'intervalle périodique de resynchronisation s'écoule, Then le contenu distant est rafraîchi en arrière-plan sans interrompre le service des requêtes en cours.

## Dépendances

- [FR-001](./FR-001-serveur-instructions.md) à [FR-004](./FR-004-serveur-adr.md) — les 4 serveurs à scan de répertoires en dépendent lorsqu'un dépôt distant est configuré.
