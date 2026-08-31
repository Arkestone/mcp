# Parties prenantes — Teikyō

| Partie prenante | Intérêt | Source |
|---|---|---|
| Développeur individuel (utilisateur d'un assistant IA) | Consomme les serveurs en local (transport stdio) depuis VS Code, Claude Code, Cursor, Windsurf, JetBrains ou Zed | README racine, section « MCP Client Configuration » |
| Équipe plateforme/outillage IA | Maintient les 6 serveurs et l'infrastructure partagée (`pkg/`), publie les images Docker et binaires | README racine, `CONTRIBUTING.md` |
| Équipe consommatrice (tout projet du portefeuille Arkestone) | Fournit son propre contenu (instructions, skills, ADR, prompts) que les serveurs découvrent et servent — ne développe pas son propre serveur de contexte | Inféré de l'architecture par répertoires locaux/dépôts GitHub |
| Administrateur d'un déploiement HTTP partagé | Exploite un ou plusieurs serveurs en mode Streamable HTTP pour une équipe entière plutôt qu'en local par développeur | README racine, section « Transport Mechanisms » |

Aucun RACI formel n'a été trouvé dans le dépôt réel — ce tableau est une inférence à partir de l'architecture et de la documentation d'installation, à confirmer avec l'équipe qui maintient le dépôt.
