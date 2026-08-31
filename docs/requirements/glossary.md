# Glossaire — Teikyō

**Serveur MCP** : processus exposant du contexte à un assistant IA via le protocole [Model Context Protocol](https://modelcontextprotocol.io). Teikyō en fournit six, un par type de contexte.

**Transport** : mécanisme de communication entre un client (assistant IA) et un serveur MCP. Deux transports supportés par tous les serveurs : `stdio` (le client démarre le processus serveur localement) et `Streamable HTTP` (le serveur tourne en continu, plusieurs clients s'y connectent à distance).

**Resource** (ressource, primitive MCP) : contenu adressable exposé par un serveur, consultable par le client (ex. `graph://stats`, `graph://node/{id}`).

**Tool** (outil, primitive MCP) : action invocable par le client sur un serveur (ex. `add-node`, `add-edge` sur `mcp-graph`).

**Prompt** (primitive MCP) : gabarit de message réutilisable exposé par un serveur (utilisé par `mcp-prompts`).

**Optimiseur LLM** (`pkg/optimizer`) : couche optionnelle de consolidation du contenu par un modèle de langage compatible OpenAI, avant de le servir — réduit le volume/bruit du contexte transmis à l'assistant appelant.

**Syncer** (`pkg/syncer`) : synchronisation périodique en arrière-plan du contenu provenant de dépôts GitHub distants, pour ne pas requérir un appel réseau à chaque requête client.

**Config en cascade** (`pkg/config`) : résolution de la configuration d'un serveur dans l'ordre YAML → variables d'environnement → flags CLI, chaque niveau pouvant surcharger le précédent.

**Nœud / Arête** (`mcp-graph`) : un nœud est une entité (label + nom + propriétés) ; une arête est une relation dirigée entre deux nœuds (type de relation + propriétés optionnelles).

**Mémoire** (`mcp-memory`) : unité de contenu persistant stockée en fichier Markdown, retrouvable par recherche plein texte ou par tag.
