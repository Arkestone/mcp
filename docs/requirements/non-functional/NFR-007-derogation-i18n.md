# NFR-007 — Dérogation explicite : internationalisation

**Statut** : Confirmé (dérogation)

Teikyō ne présente aucun texte à un utilisateur final humain dans une langue donnée : son contenu est soit transmis tel quel (fichiers d'instructions, skills, ADR déjà rédigés dans la langue choisie par l'équipe qui les a écrits), soit consommé par un LLM (mémoire, graphe). Le standard portefeuille d'internationalisation obligatoire (`Shared/architecture/product-platform-standards.md`) ne s'applique donc pas à une interface utilisateur de Teikyō, qui n'en a pas.

## Nuance à ne pas perdre

Le contenu SERVI par Teikyō (instructions, skills d'une équipe) peut lui-même être rédigé dans n'importe quelle langue — mais c'est la responsabilité de l'équipe consommatrice qui l'a rédigé, pas une exigence i18n interne à Teikyō lui-même, qui reste agnostique de la langue du contenu qu'il transporte.
