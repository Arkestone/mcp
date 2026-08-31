# NFR-006 — Dérogation explicite : accessibilité et micro-frontends

**Statut** : Confirmé (dérogation)

Teikyō est un ensemble de serveurs backend/CLI consommés par des clients MCP (extensions IDE, applications de bureau) — il n'expose aucune interface utilisateur graphique propre. Les standards portefeuille d'accessibilité (`Shared/architecture/product-platform-standards.md`, WCAG 2.2 AAA/RGAA 4.1.2) et de micro-frontends intégrables ne s'appliquent donc pas ici.

## Justification de la dérogation

Contrairement à un projet comme Meikan ou Kiseki qui expose un portail web consulté par un humain, l'interface de Teikyō est le protocole MCP lui-même, consommé par un logiciel client (VS Code, Claude Desktop, etc.) qui porte lui-même sa propre responsabilité d'accessibilité. Si un tableau de bord web d'administration était ajouté un jour à Teikyō (non identifié dans le dépôt réel actuel), ce NFR devrait être révisé et les standards a11y/micro-frontends s'appliqueraient à cette interface spécifiquement.
