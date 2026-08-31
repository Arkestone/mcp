# FR-008 — Optimisation du contenu par LLM (optionnelle)

**Statut** : Draft — décomposé à partir du produit réel (`pkg/optimizer`)

En tant qu'utilisateur d'un serveur de contexte (instructions, skills, prompts ou ADR), je veux pouvoir activer une consolidation du contenu par un modèle de langage avant qu'il soit servi, afin de réduire le volume/bruit transmis à l'assistant appelant lorsque les sources sont nombreuses ou verbeuses.

## Critères d'acceptation

- Given l'option d'optimisation désactivée (comportement par défaut), When un serveur sert du contenu, Then il transmet le contenu brut des fichiers sources, sans appel à un LLM externe.
- Given l'option activée avec un endpoint compatible OpenAI configuré, When le contenu est servi, Then il passe d'abord par l'optimiseur avant transmission au client.
- Given l'endpoint LLM configuré injoignable, When l'optimisation est activée, Then le serveur ne doit pas bloquer indéfiniment le service du contenu — comportement de repli exact `[QUESTION OUVERTE]`, non confirmé à la lecture du README.

## Dépendances

Optionnelle pour FR-001 à FR-004 (serveurs à scan de répertoires) ; non confirmée pour FR-005/FR-006 (voir leurs notes respectives).
