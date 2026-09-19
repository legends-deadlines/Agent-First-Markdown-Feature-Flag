# MDFLAG File Format Specification

Version: 1.0
Last updated: 2026-09-17

## Overview

MDFLAG files use Markdown with YAML frontmatter. The frontmatter contains machine-readable configuration, while the Markdown body contains human-readable documentation.

## File Structure

```
---
<YAML frontmatter>
---

<Markdown body>
```

## Frontmatter Fields

### Required Fields

| Field | Type | Description |
|-------|------|-------------|
| `name` | string | Unique flag identifier. Must be kebab-case. |
| `percentage` | integer | Enable percentage (0-100). |
| `status` | string | Flag status: `active`, `paused`, or `archived`. |

### Optional Fields

| Field | Type | Default | Description |
|-------|------|---------|-------------|
| `targeting` | string | `user_id` | Targeting method: `user_id`, `session_id`, `random` |
| `created` | timestamp | current time | Flag creation time (RFC3339) |
| `author` | string | empty | Who created the flag |
| `expires` | timestamp | never | Expiration time (RFC3339) |
| `hypothesis` | string | empty | Experiment hypothesis |
| `metrics` | array | empty | List of metrics to track |

### Integrity Fields

| Field | Type | Description |
|-------|------|-------------|
| `agent_section_hash` | string | SHA-256 hash of agent-controlled content |
| `human_section_hash` | string | SHA-256 hash of human-controlled content |

## Hash Computation

### Agent Section Hash

Covers fields that the agent sets at creation time:

```
canonical = "name={name}|hypothesis={hypothesis}|metrics={sorted_metrics}|body={body}"
agent_section_hash = "sha256:" + SHA256(canonical)
```

Where:
- `sorted_metrics` is the comma-joined list of metrics, sorted alphabetically
- `body` is the entire Markdown body content

### Human Section Hash

Covers fields that control flag behavior:

```
canonical = "percentage={percentage}|status={status}|targeting={targeting}|expires={expires}"
human_section_hash = "sha256:" + SHA256(canonical)
```

Where:
- `expires` is formatted as RFC3339 in UTC, or empty string if not set

## Example

```markdown
---
name: new-checkout-flow
percentage: 5
targeting: user_id
status: active
created: 2026-09-17T10:00:00Z
author: claude-code-agent
expires: 2026-12-17T00:00:00Z
hypothesis: Increase checkout conversion by 5%
metrics:
    - conversion_rate
    - checkout_time_seconds
agent_section_hash: "sha256:a1b2c3d4..."
human_section_hash: "sha256:e5f6a7b8..."
---

## Description
New one-page checkout flow replacing the 3-step process.

## Implementation Notes
Changes affected:
- frontend/src/checkout/OnePageCheckout.tsx
- backend/api/checkout.go

## Rollout Plan
1. 5% — internal team (2 days)
2. 25% — all new users (5 days)
3. 100% — all users (after metrics confirmation)
```

## File Naming

Flag files must be named `{flag-name}.md` and stored in the flags directory (default: `.mdflag/`).

The flag name in the frontmatter must match the filename without extension.

## Validation Rules

1. `name` must be non-empty and contain only lowercase letters, numbers, and hyphens
2. `percentage` must be between 0 and 100 inclusive
3. `status` must be one of: `active`, `paused`, `archived`
4. `targeting` must be one of: `user_id`, `session_id`, `random`
5. Hashes must match computed values for the current content
6. File must not be a symbolic link
7. File must not contain path traversal characters

## Modification Rules

### Agent Modifications
- Can create new flags through `mdflag_create`
- Cannot modify existing flags
- Cannot delete flags

### Human Modifications
- Can create flags through CLI or MCP
- Can change `percentage` through `mdflag rollout`
- Can change `status` through CLI
- Can delete flags through file system

### Automatic Recalculation
When a flag is modified through MDFLAG tools, hashes are automatically recalculated and written back to the file.

Direct file modifications will cause hash mismatches, detected by `mdflag verify`.
