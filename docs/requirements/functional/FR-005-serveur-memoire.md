# FR-005 — Mémoire persistante pour assistant IA

**Statut** : Draft — décomposé à partir du produit réel (`servers/mcp-memory`)

En tant qu'assistant IA de développement, je veux mémoriser une information au fil d'une session et la retrouver lors d'une session ultérieure, afin de ne pas redemander à l'utilisateur un contexte déjà fourni précédemment.

## Critères d'acceptation

- Given une information soumise via l'outil `remember` avec un ou plusieurs tags, When elle est enregistrée, Then elle persiste sous forme de fichier Markdown sur disque, survivant à un redémarrage du serveur.
- Given une recherche via l'outil `recall` par mot-clé ou par tag, When la recherche est exécutée, Then elle porte à la fois sur le contenu et sur les tags de chaque mémoire (recherche plein texte).
- Given une mémoire devenue obsolète, When l'outil `forget` est invoqué sur son identifiant, Then elle n'est plus retournée par `recall`.

## Dépendances

- [FR-007](./FR-007-double-transport.md), [FR-010](./FR-010-configuration-en-cascade.md).

## Note

Contrairement à FR-001 à FR-004, ce serveur ne scanne pas de répertoires externes : il gère sa propre persistance locale. Ne pas lui appliquer FR-008 (optimisation LLM) sans vérifier que le produit réel l'implémente effectivement pour ce serveur — non confirmé à la lecture du README.
