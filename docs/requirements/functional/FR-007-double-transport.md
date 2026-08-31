# FR-007 — Double transport stdio / Streamable HTTP

**Statut** : Draft — décomposé à partir du produit réel (`pkg/server`, commun aux 6 serveurs)

En tant qu'utilisateur d'un serveur Teikyō, je veux pouvoir le faire tourner soit en mode local (le client spawn le processus), soit en mode partagé (le serveur tourne en continu, plusieurs clients s'y connectent), afin d'adapter le déploiement à mon contexte sans changer d'outil.

## Critères d'acceptation

- Given aucun flag de transport spécifié, When un serveur démarre, Then il utilise `stdio` par défaut.
- Given le flag `-transport http` et une adresse d'écoute, When le serveur démarre, Then il expose le protocole MCP en Streamable HTTP sur cette adresse, consultable par plusieurs clients simultanément.
- Given un déploiement HTTP, When une requête de supervision est faite sur `/healthz`, Then le serveur répond sans dépendre de l'état d'une requête MCP en cours.

## Dépendances

Transversal — condition préalable de FR-001 à FR-006.
