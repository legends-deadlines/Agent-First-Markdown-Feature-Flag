<div align="center">

<img src="../.github/assets/banner.svg" alt="Agent-First Feature Flags" width="100%"/>

<br/>

[English](../README.md) | [Русский](README.ru.md) | [Español](README.es.md) | [中文](README.zh.md) | [日本語](README.ja.md) | [Français](README.fr.md) | [Deutsch](README.de.md)

</div>

---

# MDFLAG

**Feature-Flags für KI-Agenten. Kein SaaS. Keine API-Schlüssel. Nur einfache Markdown-Dateien in Ihrem Git-Repository.**

KI-Coding-Agenten (Cursor, Claude Code, Copilot, Windsurf) generieren und fügen Code direkt in Ihr Projekt ein. Wenn der Code Fehler enthält, bricht die gesamte Anwendung zusammen. Ein Rollback ist teuer und zeitaufwendig.

Herkömmliche Feature-Flag-Tools (LaunchDarkly, Statsig, Flipt) sind für autonome Agenten unerreichbar: Sie erfordern eine Website-Registrierung, API-Schlüssel und SDK-Konfigurationen. KI-Agenten besitzen weder Browser noch E-Mail-Adressen oder Kreditkarten.

MDFLAG löst dieses Problem, indem es Feature-Flags als einfache Markdown-Dateien im Projektverzeichnis bereitstellt. Da KI-Agenten von Natur aus mit Dateien umgehen können, können sie Flags direkt und nativ nutzen.

---

## Kernidee

**KI kann Dateien in Git erstellen, sich aber nicht bei SaaS-Diensten registrieren.**

Daher gilt:
- Traditionelle Flags (LaunchDarkly) = für Agenten unzugänglich.
- Git-basierte Flags (MDFLAG) = für Agenten direkt zugänglich.

MDFLAG ist die erste Lösung, die auf den vorhandenen Fähigkeiten der KI aufbaut, anstatt zu versuchen, ihr komplexe externe Integrationen beizubringen.

---

## Funktionsweise

**Ohne Feature-Flag:**
Die KI baut einen neuen Raum, indem sie den alten abreißt. Wenn der neue Raum schief wird, stehen Sie ohne Dach da.

**Mit Feature-Flag:**
Die KI baut den neuen Raum neben dem alten. Sie entscheiden, in welchem Sie wohnen möchten. Wenn mit dem neuen Raum etwas nicht stimmt, schließen Sie einfach die Tür und bleiben im alten.

---

## Funktionen

- **Agent-native**: Funktioniert über MCP (Model Context Protocol) — keine Registrierung, keine API-Schlüssel
- **Git-basiert**: Flags sind Markdown-Dateien in Ihrem Repository und werden zusammen mit Ihrem Code versioniert
- **Integritätsschutz**: SHA-256-basierte Manipulationserkennung verhindert unbefugte Änderungen
- **Schrittweise Einführung**: Flags schrittweise für 5% → 25% → 100% der Benutzer aktivieren
- **Sofortiges Rollback**: Ändern Sie eine einzelne Zahl in einer Datei, um die Funktion sofort zu deaktivieren
- **Menschliche Kontrolle**: Agenten erstellen Flags, Menschen verwalten die Einführung

---

## Installation

### Aus GitHub Releases

Laden Sie das neueste Binary für Ihre Plattform von [Releases](https://github.com/legends-deadlines/Agent-First-Markdown-Feature-Flag/releases) herunter:

```bash
# Linux/macOS
curl -L https://github.com/legends-deadlines/Agent-First-Markdown-Feature-Flag/releases/latest/download/mdflag-linux-amd64 -o mdflag
chmod +x mdflag
sudo mv mdflag /usr/local/bin/

# Installation überprüfen
mdflag --help
```

### Aus Quellcode bauen

```bash
git clone https://github.com/legends-deadlines/Agent-First-Markdown-Feature-Flag.git
cd Agent-First-Markdown-Feature-Flag
go build -o mdflag ./cmd/mdflag/
```

---

## Schnelleinstieg

### 1. Flag erstellen

```bash
mdflag create \
  --name new-checkout-flow \
  --percentage 0 \
  --hypothesis "Konvertierung des Checkouts um 5% steigern" \
  --metrics "conversion_rate,checkout_time_seconds" \
  --author "claude-code-agent" \
  --description "Neuer einseitiger Checkout-Prozess als Ersatz für den 3-Schritt-Prozess"
```

Dies erstellt die Datei `.mdflag/new-checkout-flow.md`:

```markdown
---
name: new-checkout-flow
percentage: 0
targeting: user_id
status: active
created: 2026-09-17T10:00:00Z
author: claude-code-agent
hypothesis: Konvertierung des Checkouts um 5% steigern
metrics:
    - conversion_rate
    - checkout_time_seconds
agent_section_hash: "sha256:a1b2c3..."
human_section_hash: "sha256:d4e5f6..."
---

## Description
Neuer einseitiger Checkout-Prozess als Ersatz für den 3-Schritt-Prozess
```

### 2. Flag im Code verwenden

```go
import "github.com/legends-deadlines/Agent-First-Markdown-Feature-Flag/pkg/mdflag"

// Client einmalig beim Anwendungsstart initialisieren
client, err := mdflag.New(".mdflag")
if err != nil {
    log.Fatal(err)
}

// Flag im Request-Handler auswerten
func handleCheckout(w http.ResponseWriter, r *http.Request) {
    userID := getUserID(r)

    if client.Enabled("new-checkout-flow", userID) {
        // Neuer Pfad
        renderOnePageCheckout(w, r)
    } else {
        // Bisheriger stabiler Pfad
        renderThreeStepCheckout(w, r)
    }
}
```

### 3. Schrittweise Einführung

```bash
# Für 5% der Benutzer aktivieren
mdflag rollout --name new-checkout-flow --percentage 5

# Metriken überwachen, dann auf 25% erhöhen
mdflag rollout --name new-checkout-flow --percentage 25

# Wenn alles stabil ist, für alle aktivieren (100%)
mdflag rollout --name new-checkout-flow --percentage 100

# Bei Fehlern sofort deaktivieren (0%)
mdflag rollout --name new-checkout-flow --percentage 0
```

### 4. Integrität prüfen

```bash
mdflag verify
```

---

## Integration mit KI-Agenten

MDFLAG verbindet sich mit KI-Agenten über MCP (Model Context Protocol).

### Für Cursor

In `.cursor/mcp.json` im Projektverzeichnis hinzufügen:

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

### Für Claude Code

In `.claude/mcp.json` hinzufügen:

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

## Verfügbare MCP-Tools

| Tool | Agentenzugriff | Beschreibung |
|------|----------------|--------------|
| `mdflag_create` | Ja | Neues Feature-Flag erstellen |
| `mdflag_list` | Ja | Alle Flags und deren Status auflisten |
| `mdflag_verify` | Ja | Integrität aller Flags prüfen |
| `mdflag rollout` | Nein | Rollout-Prozentsatz ändern (Nur CLI) |

---

## Lizenz

MIT-Lizenz — siehe [LICENSE](../LICENSE) für Details.

---

<div align="center">

<img src="../.github/assets/seal.svg" alt="MDFLAG Seal Mascot" width="600"/>

</div>
