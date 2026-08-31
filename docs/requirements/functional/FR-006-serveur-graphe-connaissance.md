# FR-006 — Graphe de connaissance pour assistant IA

**Statut** : Draft — décomposé à partir du produit réel (`servers/mcp-graph`)

En tant qu'assistant IA de développement, je veux construire et interroger un graphe d'entités et de relations propre à un domaine, afin de raisonner sur des associations (qui dépend de quoi, qui est lié à qui) plutôt que de ne manipuler que des faits isolés.

## Critères d'acceptation

- Given un nœud identifié par un label, un nom et des propriétés optionnelles, When l'outil `add-node` est invoqué, Then le nœud est créé et devient consultable via `graph://node/{id}`.
- Given deux nœuds existants, When l'outil `add-edge` crée une relation dirigée typée entre eux, Then cette relation est prise en compte dans les traversées ultérieures du graphe.
- Given deux nœuds distants dans le graphe, When une recherche de plus court chemin est demandée, Then le résultat est calculé par parcours en largeur (BFS) plutôt que par une heuristique non déterministe.
- Given un redémarrage du serveur, When le graphe est rechargé depuis son fichier JSON, Then l'ensemble des index (recherche par label/nom, relations) est reconstruit intégralement, sans perte de nœud ni d'arête.

## Dépendances

- [FR-007](./FR-007-double-transport.md), [FR-010](./FR-010-configuration-en-cascade.md).

## Note

Comme FR-005, ce serveur gère sa propre persistance locale (écriture JSON atomique) plutôt que de scanner des répertoires externes.
