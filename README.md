<div align="center">

<img src=".github/assets/banner.svg" alt="Agent-First Feature Flags" width="100%"/>

<br/>

[![Status](https://img.shields.io/badge/⚙%20Status-In%20Development-f78166?style=flat-square&labelColor=1c2128)](#)
[![License](https://img.shields.io/badge/License-MIT-79c0ff?style=flat-square&labelColor=1c2128)](LICENSE)

**🇷🇺 [Читать на русском](docs/README.ru.md)**

</div>

---

> [!WARNING]
> **🚧 This project is under active development.** The concept is locked in, tech stack is being decided.

## The Problem

AI coding agents (**Cursor, Claude Code, Copilot**) generate code with zero awareness of feature flags. They don't know which features are rolled out to 10% of users, which flags are mutually exclusive, or what guardrails exist.

The result: **28–40% excess code churn**, **46.4% AI-generated PRs rejected**, **1.7× more defects**.

Existing feature flag tools (LaunchDarkly, Flipt, dif.sh) were built for humans clicking dashboards — not for autonomous agents reading repos.

## The Solution

**Feature flags as Markdown files in your Git repo.** A format that both humans and AI agents can read natively.

```
flags/*.md ──► CLI Compiler ──► .cursorrules, CLAUDE.md, context.json
                   ▲
            codebase AST scan
```

The CLI scans your flag definitions, analyzes where flags are used in code, detects conflicts and dependencies, then generates rules that AI agents follow automatically.

## How a Flag Looks

```markdown
---
id: new_payment_flow
status: active
rollout: 25
exclusion_group: checkout_v2
dependencies:
  - user_service_v2
guardrails:
  - "Do NOT modify legacy_payment.go inside this feature."
---

# New Payment Flow

Progressive rollout of the updated payment pipeline.
```

After `dif compile`, your AI agent sees:

```
ACTIVE: new_payment_flow (25%), user_service_v2 (100%)
INACTIVE: old_checkout, legacy_payment
RULE: new_payment_flow ⊕ old_checkout (mutually exclusive)
RULE: Do NOT modify legacy_payment.go inside new_payment_flow
```

## Why Not Just Use X?

|  | **This project** | dif.sh | Flipt | LaunchDarkly |
|:---|:---:|:---:|:---:|:---:|
| Flags live in Git | ✅ | ✅ | ❌ | ❌ |
| AI agents read flags autonomously | ✅ | ❌ | ❌ | ❌ |
| Auto-generated agent guardrails | ✅ | ❌ | ❌ | ❌ |
| Multi-language codebase analysis | ✅ | ❌ | ❌ | ❌ |
| Offline deterministic targeting | ✅ | ❌ | ❌ | ❌ |

## Roadmap

- [x] Concept & architecture
- [ ] Flag parser & schema validation
- [ ] Codebase AST scanner
- [ ] Dependency graph & conflict detection
- [ ] Agent rules generator (`.cursorrules`, `CLAUDE.md`)
- [ ] Runtime targeting engine
- [ ] CLI & CI/CD integration

## License

MIT

---

<div align="center">

<img src=".github/assets/seal.svg" alt="Seal mascot" width="600"/>

<sub>by <a href="https://github.com/legends-deadlines">legends-deadlines</a></sub>

</div>
