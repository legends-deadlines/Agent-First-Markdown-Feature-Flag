<div align="center">

<img src="../.github/assets/banner.svg" alt="Agent-First Feature Flags" width="100%"/>

<br/>

[English](../README.md) | [Русский](README.ru.md) | [Español](README.es.md) | [中文](README.zh.md) | [日本語](README.ja.md) | [Français](README.fr.md) | [Deutsch](README.de.md)

</div>

---

# MDFLAG

**Feature flags pour agents d'IA. Sans SaaS. Sans clés API. Juste des fichiers Markdown dans votre dépôt Git.**

Les agents de codage IA (Cursor, Claude Code, Copilot, Windsurf) génèrent et insèrent directement du code dans votre projet. Si le code contient des bugs, tout s'effondre. Annuler les modifications est coûteux et lent.

Les outils de feature flags traditionnels (LaunchDarkly, Statsig, Flipt) sont inaccessibles aux agents autonomes : ils nécessitent une inscription sur un site web, des clés API et la configuration d'un SDK. Les agents n'ont ni navigateurs, ni adresses e-mail, ni cartes de crédit.

MDFLAG résout ce problème en transformant les feature flags en simples fichiers Markdown stockés dans le répertoire de votre projet. Les agents d'IA sachant déjà manipuler des fichiers, ils peuvent ainsi utiliser les flags de manière native.

---

## Idée Clé

**L'IA peut créer des fichiers dans Git, mais ne peut pas s'inscrire à des services SaaS.**

Par conséquent :
- Flags traditionnels (LaunchDarkly) = inaccessibles aux agents.
- Flags basés sur Git (MDFLAG) = accessibles aux agents.

MDFLAG est la première solution qui s'appuie sur les capacités existantes de l'IA, plutôt que d'essayer de lui enseigner des intégrations externes complexes.

---

## Comment Ça Marche

**Sans feature flag :**
L'IA construit une nouvelle pièce en démolissant l'ancienne. Si la nouvelle pièce est bancale, vous vous retrouvez sans toit.

**Avec un feature flag :**
L'IA construit la nouvelle pièce à côté de l'ancienne. Vous décidez dans laquelle vivre. En cas de problème, il vous suffit de fermer la porte et de rester dans l'ancienne.

---

## Fonctionnalités

- **Agent-native** : Fonctionne via MCP (Model Context Protocol) — sans inscription ni clés API
- **Basé sur Git** : Les flags sont des fichiers Markdown dans votre dépôt, versionnés avec votre code
- **Protection de l'intégrité** : Détection des altérations basée sur SHA-256 pour empêcher les modifications non autorisées
- **Déploiement progressif** : Activez les flags pour 5% → 25% → 100% des utilisateurs
- **Annulation instantanée** : Modifiez un seul nombre dans un fichier pour désactiver immédiatement la fonctionnalité
- **Contrôle humain** : Les agents créent les flags, les humains gèrent le déploiement

---

## Démarrage Rapide (CLI)

```bash
# Créer un flag (déploiement à 0%)
mdflag create --name new-checkout --hypothesis "Augmenter la conversion de 5%"

# Modifier le déploiement (25%)
mdflag rollout --name new-checkout --percentage 25

# Vérifier l'intégrité des hachages
mdflag verify

# Lister tous les flags
mdflag list
```

---

## Utilisation dans le Code (Go)

```go
import "github.com/legends-deadlines/mdflag/pkg/mdflag"

client, err := mdflag.New(".mdflag")
if err != nil {
    log.Fatal(err)
}

if client.Enabled("new-checkout", userID) {
    // Nouvelle fonctionnalité
} else {
    // Code stable existant
}
```

---

## Licence

Licence MIT — voir [LICENSE](../LICENSE) pour plus de détails.

---

<div align="center">

<img src="../.github/assets/seal.svg" alt="MDFLAG Seal Mascot" width="600"/>

</div>
