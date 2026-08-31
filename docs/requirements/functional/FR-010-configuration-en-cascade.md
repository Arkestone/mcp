# FR-010 — Configuration en cascade YAML → environnement → CLI

**Statut** : Draft — décomposé à partir du produit réel (`pkg/config`)

En tant qu'exploitant d'un serveur Teikyō, je veux pouvoir définir sa configuration à plusieurs niveaux (fichier YAML, variables d'environnement, flags de ligne de commande) avec une priorité claire, afin d'adapter un même déploiement à différents environnements sans dupliquer les fichiers de configuration.

## Critères d'acceptation

- Given une valeur définie à la fois en YAML et en variable d'environnement, When le serveur démarre, Then la variable d'environnement prévaut sur la valeur YAML.
- Given une valeur définie à la fois en variable d'environnement et en flag CLI, When le serveur démarre, Then le flag CLI prévaut sur la variable d'environnement.
- Given un token GitHub configurable par plusieurs sources (flag, variable préfixée, variable globale, YAML), When plusieurs sources sont renseignées simultanément, Then l'ordre de priorité documenté dans `current-state.md` est respecté sans ambiguïté.

## Dépendances

Transversal — condition préalable de FR-001 à FR-006.
