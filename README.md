<div align="center">

<img src=".github/assets/banner.svg" alt="Agent-First Feature Flags" width="100%"/>

<br/>

[![Status](https://img.shields.io/badge/🚀%20Status-v1.0.0%20Ready-10b981?style=flat-square&labelColor=1c2128)](#)
[![License](https://img.shields.io/badge/License-MIT-79c0ff?style=flat-square&labelColor=1c2128)](LICENSE)
[![Go](https://img.shields.io/badge/Go-1.27+-00ADD8?style=flat-square&logo=go&labelColor=1c2128)](#)

<br/>

[English](README.md) | [Русский](docs/README.ru.md) | [Español](docs/README.es.md) | [中文](docs/README.zh.md) | [日本語](docs/README.ja.md) | [Français](docs/README.fr.md) | [Deutsch](docs/README.de.md)

</div>

---

# MDFLAG

**Feature flags for AI agents. No SaaS. No API keys. Just Markdown files in your Git repo.**

AI coding agents (Cursor, Claude Code, Copilot, Windsurf) write code and insert it directly into your project. If the code is bad, everything breaks. Rolling back is expensive and slow.

Existing feature flag tools (LaunchDarkly, Statsig, Flipt) are inaccessible to agents — they require website registration, API keys, SDK configuration. Agents don't have browsers, email addresses, or credit cards.

MDFLAG solves this by making feature flags plain Markdown files in your project directory. AI agents already know how to create files — so they can use flags.

---

## Key Idea

**AI can create files in Git, but cannot register for SaaS.**

Therefore:
- Traditional flags (LaunchDarkly) = inaccessible for agents
- Git-based flags (MDFLAG) = accessible for agents

MDFLAG is the first solution that works with what AI already does, rather than trying to teach it something new.

---

## How It Works

**Without a flag:**
AI builds a new room by demolishing the old one. If the new room is crooked, you have no room at all.

**With a flag:**
AI builds a new room next to the old one. You decide which one to live in. If the new room is crooked, just close the door and live in the old one.

---

## Features

- **Agent-native**: Works through MCP (Model Context Protocol) — no registration, no API keys
- **Git-based**: Flags are Markdown files in your repository, versioned with your code
- **Integrity protection**: Hash-based tamper detection prevents unauthorized modifications
- **Gradual rollout**: Enable flags for 5% → 25% → 100% of users
- **Instant rollback**: Change one number in a file to disable the feature
- **Human control**: Agents create flags, humans manage rollout

---

## Installation

### From GitHub Releases

Download the latest binary for your platform from [Releases](https://github.com/legends-deadlines/Agent-First-Markdown-Feature-Flag/releases):

```bash
# Linux/macOS
curl -L https://github.com/legends-deadlines/Agent-First-Markdown-Feature-Flag/releases/latest/download/mdflag-linux-amd64 -o mdflag
chmod +x mdflag
sudo mv mdflag /usr/local/bin/

# Verify installation
mdflag --help
```

### From Source

```bash
git clone https://github.com/legends-deadlines/Agent-First-Markdown-Feature-Flag.git
cd Agent-First-Markdown-Feature-Flag
go build -o mdflag ./cmd/mdflag/
```

---

## Quick Start

### 1. Create a flag

```bash
mdflag create \
  --name new-checkout-flow \
  --percentage 0 \
  --hypothesis "Increase checkout conversion by 5%" \
  --metrics "conversion_rate,checkout_time_seconds" \
  --author "claude-code-agent" \
  --description "New one-page checkout flow replacing the 3-step process"
```

This creates `.mdflag/new-checkout-flow.md`:

```markdown
---
name: new-checkout-flow
percentage: 0
targeting: user_id
status: active
created: 2026-09-17T10:00:00Z
author: claude-code-agent
hypothesis: Increase checkout conversion by 5%
metrics:
    - conversion_rate
    - checkout_time_seconds
agent_section_hash: "sha256:a1b2c3..."
human_section_hash: "sha256:d4e5f6..."
---

## Description
New one-page checkout flow replacing the 3-step process
```

### 2. Use the flag in your code

```go
import "github.com/legends-deadlines/Agent-First-Markdown-Feature-Flag/pkg/mdflag"

// Initialize client once at application startup
client, err := mdflag.New(".mdflag")
if err != nil {
    log.Fatal(err)
}

// Check the flag in your request handler
func handleCheckout(w http.ResponseWriter, r *http.Request) {
    userID := getUserID(r)

    if client.Enabled("new-checkout-flow", userID) {
        // New code path
        renderOnePageCheckout(w, r)
    } else {
        // Old code path
        renderThreeStepCheckout(w, r)
    }
}
```

### 3. Gradually roll out

```bash
# Enable for 5% of users
mdflag rollout --name new-checkout-flow --percentage 5

# Monitor metrics, then increase to 25%
mdflag rollout --name new-checkout-flow --percentage 25

# If everything is good, enable for everyone
mdflag rollout --name new-checkout-flow --percentage 100

# If something is wrong, disable instantly
mdflag rollout --name new-checkout-flow --percentage 0
```

### 4. Verify integrity

```bash
mdflag verify
```

Expected output:
```
[OK]   .mdflag/new-checkout-flow.md

Results: 1 valid, 0 invalid, 1 total
```

---

## AI Agent Integration

MDFLAG integrates with AI agents through MCP (Model Context Protocol).

### For Cursor

Add to `.cursor/mcp.json` in your project root:

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

### For Claude Code

Add to `.claude/mcp.json`:

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

### For Windsurf

Add to `.windsurf/mcp.json`:

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

### Agent Instructions

Add to your agent's system prompt or `CLAUDE.md` / `.cursorrules`:

```
When making experimental code changes:
1. Create a feature flag using mdflag_create BEFORE writing the code
2. Wrap new code with: if client.Enabled("flag-name", entityID) { ... }
3. Set initial percentage to 0
4. Do NOT modify flag files directly — use mdflag tools only
5. Commit both the flag file and your code changes together
```

---

## Available MCP Tools

| Tool | Agent Access | Description |
|------|--------------|-------------|
| `mdflag_create` | Yes | Create a new feature flag |
| `mdflag_list` | Yes | List all flags with their status |
| `mdflag_verify` | Yes | Check integrity of all flags |
| `mdflag rollout` | No | Change flag percentage (CLI only) |

This access model ensures agents can create flags but cannot enable them or modify existing ones.

---

## Security Model

MDFLAG uses a three-layer protection system:

### Layer 1: Hash-based integrity

Each flag file contains SHA-256 hashes of its content sections:
- `agent_section_hash`: covers fields the agent sets at creation
- `human_section_hash`: covers fields that control behavior (percentage, status)

If anyone modifies the file directly without using MDFLAG tools, the hashes won't match and `mdflag verify` will detect it.

### Layer 2: Tool separation

Agents have access to `mdflag_create`, `mdflag_list`, `mdflag_verify`.
Agents do NOT have access to `mdflag_rollout` or delete operations.
Only humans can change flag percentages through the CLI.

### Layer 3: Git hooks (optional)

Pre-commit hooks can block commits with integrity violations. See [docs/git-hooks.md](docs/git-hooks.md).

---

## CLI Reference

### `mdflag create`

Create a new feature flag.

```bash
mdflag create --name <name> --hypothesis <text> [options]

Options:
  --name          Flag name (required, kebab-case)
  --hypothesis    Experiment hypothesis (required)
  --percentage    Enable percentage 0-100 (default: 0)
  --targeting     Targeting: user_id, session_id, random (default: user_id)
  --author        Author name
  --metrics       Comma-separated metric names
  --expires       Expiration date (RFC3339)
  --description   Human-readable description
  --dir           Flags directory (default: .mdflag)
```

### `mdflag rollout`

Change the enable percentage of a flag.

```bash
mdflag rollout --name <name> --percentage <0-100>

Options:
  --name          Flag name (required)
  --percentage    New percentage 0-100 (required)
  --dir           Flags directory (default: .mdflag)
```

### `mdflag list`

List all flags in the directory.

```bash
mdflag list [dir]
```

### `mdflag verify`

Check integrity of all flags.

```bash
mdflag verify [dir]
```

### `mdflag serve`

Start MCP server for AI agents.

```bash
mdflag serve [--dir <path>]
```

---

## Runtime Library

The Go runtime library provides flag evaluation for your application.

```go
package main

import (
    "log"
    "github.com/legends-deadlines/Agent-First-Markdown-Feature-Flag/pkg/mdflag"
)

func main() {
    // Initialize once at startup
    client, err := mdflag.New(".mdflag")
    if err != nil {
        log.Fatal(err)
    }

    // Optional: enable hot-reload
    stop, err := client.StartWatcher()
    if err == nil {
        defer stop()
    }

    // Use in request handlers
    userID := "user_12345"
    if client.Enabled("new-checkout-flow", userID) {
        // New code path
    } else {
        // Old code path
    }
}
```

### Deterministic Targeting

The same entity ID always gets the same result for a given flag. This ensures:
- Users don't see features flicker between requests
- A/B tests have consistent groups
- Debugging is reproducible

### Percentage Evaluation

For flags with percentage between 1 and 99, MDFLAG uses FNV-1a hashing:

```
bucket = FNV-1a(flagName + ":" + entityID) % 100
enabled = bucket < percentage
```

This provides uniform distribution across entities.

---

## Flag Lifecycle

1. **Created** (percentage: 0) — Agent creates the flag and wraps code
2. **Review** — Human reviews the PR including the flag file
3. **Gradual rollout** — Human increases percentage: 5% → 25% → 100%
4. **Monitoring** — Track metrics defined in the flag
5. **Cleanup** — Either:
   - Remove the flag and keep new code (if successful)
   - Set percentage to 0 and revert code (if failed)

---

## Comparison with Alternatives

| Tool | How It Works | Problem for AI Agents |
|------|--------------|----------------------|
| LaunchDarkly / Statsig | SaaS service, requires API keys | Agent cannot register for accounts |
| Flipt v2 | GitOps, but requires local server | Agent cannot start a server |
| dif.sh | Markdown flags | No protection against agent modification |
| **MDFLAG** | Markdown flags + integrity protection | **Works autonomously** |

Key difference: MDFLAG not only allows agents to create flags, but also **protects them from modification** — agents cannot enable flags themselves, cannot remove markers, cannot change rules.

---

## Philosophy

**Don't try to make AI smarter.** AI will still write bad code sometimes.

**Instead:** create a system where bad AI code doesn't break the project. Isolate AI code, enable gradually, rollback instantly.

This is not about "replacing humans", but about "giving humans control over AI speed".

---

## Market Validation

- **dif.sh** reached #1 on Product Hunt (201 votes) — shows demand for the idea
- **GitClear research**: code churn in AI projects increased to 28-40%
- **46.4%** of agent-generated PRs are rejected in review

The problem is real, the solution is needed.

---

## Development

### Build

```bash
go build ./...
```

### Test

```bash
go test ./... -v
```

### Lint

```bash
go vet ./...
gofmt -l .
```

---

## License

MIT License - see [LICENSE](LICENSE) for details.

---

## Roadmap

- [x] Core: file format, parser, hasher, validator
- [x] CLI commands: create, rollout, verify, list
- [x] Runtime library with deterministic targeting
- [x] MCP server for AI agents
- [ ] Git hooks for commit-time protection
- [ ] Multi-language runtime support (TypeScript, Python)
- [ ] Flag analytics and metrics integration
- [ ] Team features: CODEOWNERS integration, CI checks
- [ ] Web dashboard for flag management

---

<div align="center">

<img src=".github/assets/seal.svg" alt="MDFLAG Seal Mascot" width="600"/>

</div>


