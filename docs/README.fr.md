<div align="center">

<img src="../.github/assets/banner.svg?v=1.0.2" alt="Agent-First Feature Flags" width="100%"/>

<br/>

[![Statut](https://img.shields.io/badge/🚀%20Statut-v1.0.2%20Prêt-10b981?style=flat-square&labelColor=1c2128)](#)
[![Licence](https://img.shields.io/badge/Licence-MIT-79c0ff?style=flat-square&labelColor=1c2128)](../LICENSE)
[![Go](https://img.shields.io/badge/Go-1.27+-00ADD8?style=flat-square&logo=go&labelColor=1c2128)](#)

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

## Installation

### Depuis les Releases GitHub

Téléchargez le dernier binaire correspondant à votre plateforme depuis les [Releases](https://github.com/legends-deadlines/Agent-First-Markdown-Feature-Flag/releases) :

```bash
# Linux/macOS
curl -L https://github.com/legends-deadlines/Agent-First-Markdown-Feature-Flag/releases/latest/download/mdflag-linux-amd64 -o mdflag
chmod +x mdflag
sudo mv mdflag /usr/local/bin/

# Vérifier l'installation
mdflag --help
```

### Depuis le Code Source

```bash
git clone https://github.com/legends-deadlines/Agent-First-Markdown-Feature-Flag.git
cd Agent-First-Markdown-Feature-Flag
go build -o mdflag ./cmd/mdflag/
```

---

## Démarrage Rapide

### 1. Créer un flag

```bash
mdflag create \
  --name new-checkout-flow \
  --percentage 0 \
  --hypothesis "Augmenter la conversion du paiement de 5%" \
  --metrics "conversion_rate,checkout_time_seconds" \
  --author "claude-code-agent" \
  --description "Nouveau flux de paiement sur une seule page remplaçant le processus à 3 étapes"
```

Ceci crée le fichier `.mdflag/new-checkout-flow.md` :

```markdown
---
name: new-checkout-flow
percentage: 0
targeting: user_id
status: active
created: 2026-09-17T10:00:00Z
author: claude-code-agent
hypothesis: Augmenter la conversion du paiement de 5%
metrics:
    - conversion_rate
    - checkout_time_seconds
agent_section_hash: "sha256:a1b2c3..."
human_section_hash: "sha256:d4e5f6..."
---

## Description
Nouveau flux de paiement sur une seule page remplaçant le processus à 3 étapes
```

### 2. Utiliser le flag dans votre code

```go
import "github.com/legends-deadlines/Agent-First-Markdown-Feature-Flag/pkg/mdflag"

// Initialiser le client au démarrage de l'application
client, err := mdflag.New(".mdflag")
if err != nil {
    log.Fatal(err)
}

// Évaluer le flag dans votre gestionnaire de requêtes
func handleCheckout(w http.ResponseWriter, r *http.Request) {
    userID := getUserID(r)

    if client.Enabled("new-checkout-flow", userID) {
        // Nouvelle branche de code
        renderOnePageCheckout(w, r)
    } else {
        // Ancienne branche de code stable
        renderThreeStepCheckout(w, r)
    }
}
```

### 3. Déploiement progressif

```bash
# Activer pour 5% des utilisateurs
mdflag rollout --name new-checkout-flow --percentage 5

# Surveiller les métriques puis augmenter à 25%
mdflag rollout --name new-checkout-flow --percentage 25

# Si tout est correct, activer pour tout le monde (100%)
mdflag rollout --name new-checkout-flow --percentage 100

# En cas d'anomalie, désactiver instantanément (0%)
mdflag rollout --name new-checkout-flow --percentage 0
```

### 4. Vérifier l'intégrité

```bash
mdflag verify
```

---

## Intégration avec les Agents IA

MDFLAG s'interface avec les agents via MCP (Model Context Protocol).

### Pour Cursor

Dans `.cursor/mcp.json` à la racine de votre projet :

```json
{
  "mcpServers": {
    "mdflag": {
      "command": "/path/to/mdflag",
      "args": ["serve", "--dir", ".mdflag"]
    }
  }
}
```

### Pour Claude Code

Dans `.claude/mcp.json` :

```json
{
  "mcpServers": {
    "mdflag": {
      "command": "/path/to/mdflag",
      "args": ["serve", "--dir", ".mdflag"]
    }
  }
}
```

---

## Outils MCP Disponibles

| Outil | Accès Agent | Description |
|-------|-------------|-------------|
| `mdflag_create` | Oui | Créer un nouveau feature flag |
| `mdflag_list` | Oui | Lister tous les flags et leur statut |
| `mdflag_verify` | Oui | Vérifier l'intégrité de tous les flags |
| `mdflag rollout` | Non | Modifier le pourcentage du flag (CLI uniquement) |

---

## Modèle de Sécurité

MDFLAG intègre un système de protection à trois niveaux :
1. **Intégrité basée sur les hachages** : SHA-256 pour vérifier le contenu
2. **Séparation des outils** : Les agents créent des flags mais ne peuvent pas les activer
3. **Git hooks (optionnel)** : Blocage au niveau du commit en cas d'altération

---

## Licence

Licence MIT — voir [LICENSE](../LICENSE) pour plus de détails.

---

<div align="center">

<img src="../.github/assets/seal.svg" alt="MDFLAG Seal Mascot" width="600"/>

</div>
