# NFR-003 — Sécurité de l'authentification et du réseau

**Statut** : Draft

L'accès à des dépôts GitHub privés se fait via un token, résolu par ordre de priorité décroissant : flag CLI, variable d'environnement préfixée au serveur, variable d'environnement globale `GITHUB_TOKEN`, valeur YAML. Aucun token n'est requis pour les dépôts publics. Le trafic réseau supporte les proxys HTTP/HTTPS et les certificats TLS/CA personnalisés (`pkg/httputil`), pour un fonctionnement derrière un pare-feu d'entreprise.

## Critères de vérification

- Given un token défini à la fois en flag CLI et en variable d'environnement, When le serveur résout sa configuration, Then le flag CLI prévaut (priorité la plus haute).
- Given un environnement réseau derrière un proxy d'entreprise avec certificat CA personnalisé, When un serveur Teikyō accède à GitHub, Then la connexion aboutit sans désactiver la vérification TLS.

## Hors périmètre

Le contenu servi (mémoire, instructions, ADR) peut lui-même être sensible selon l'équipe consommatrice — la classification de cette sensibilité relève de chaque équipe, pas de Teikyō, qui ne fait que transporter le contenu qu'on lui pointe.
