# Gouvernance documentaire — Teikyō

Document amorcé pour combler l'absence de critères de sortie de brouillon et de cadence de revue, sur le même principe que les autres projets du portefeuille — adapté à la particularité de Teikyō : son corpus documente un produit **déjà implémenté et actif**, pas une intention à cadrer.

## Critères de sortie du statut Draft

1. **Confirmation par l'équipe qui maintient le dépôt réel** — chaque FR/NFR doit être relue par l'équipe plateforme/outillage IA (voir `stakeholders.md`) avant de passer en Confirmed, pour vérifier qu'aucune capacité n'a été mal interprétée depuis la seule lecture du code et de la documentation.
2. **Absence de contradiction avec le code source** — une exigence qui décrit un comportement non vérifiable dans le dépôt réel (`servers/*/`, `pkg/*/`) reste en Draft.

## Cadence de revue

À réévaluer à chaque changement notable du dépôt réel (nouveau serveur, changement d'architecture partagée) plutôt qu'à date fixe — ce corpus documente un produit qui évolue par son propre cycle de développement, non piloté par ce dépôt de documentation.
