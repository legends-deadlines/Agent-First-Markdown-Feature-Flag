# MDFLAG MCP Tools

MDFLAG provides MCP (Model Context Protocol) tools for AI agent integration. This document describes each tool's parameters, behavior, and usage.

## Server Configuration

The MCP server runs over stdio transport. Configure your agent to start it:

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

## Tools

### mdflag_create

Create a new feature flag.

**Access**: Available to agents

**Parameters**:

| Parameter | Type | Required | Default | Description |
|-----------|------|----------|---------|-------------|
| `name` | string | Yes | — | Unique flag name (kebab-case) |
| `hypothesis` | string | Yes | — | Expected outcome of the experiment |
| `percentage` | number | No | 0 | Initial enable percentage (0-100) |
| `targeting` | string | No | `user_id` | Targeting method |
| `author` | string | No | empty | Creator identifier |
| `description` | string | No | empty | Human-readable description |
| `metrics` | string | No | empty | Comma-separated metric names |
| `expires` | string | No | never | Expiration date (RFC3339) |

**Example Call**:

```json
{
  "name": "mdflag_create",
  "arguments": {
    "name": "new-checkout-flow",
    "hypothesis": "Increase checkout conversion by 5%",
    "percentage": 0,
    "author": "claude-code-agent",
    "description": "New one-page checkout flow",
    "metrics": "conversion_rate,checkout_time_seconds",
    "expires": "2026-12-17T00:00:00Z"
  }
}
```

**Example Response**:

```
Flag 'new-checkout-flow' created successfully.

File: .mdflag/new-checkout-flow.md
Percentage: 0%
Targeting: user_id
Hypothesis: Increase checkout conversion by 5%

NEXT STEPS:
1. Wrap your new code with a runtime check:
   if client.Enabled("new-checkout-flow", entityID) { /* new code */ }
2. Commit both the flag file and your code changes.
3. A human reviewer will gradually enable the flag using 'mdflag rollout'.

IMPORTANT: Do NOT modify the flag file directly. Integrity hashes protect it.
```

**Errors**:
- `parameter 'name' is required` — name not provided
- `parameter 'hypothesis' is required` — hypothesis not provided
- `flag "X" already exists` — duplicate flag name
- `percentage must be 0-100` — invalid percentage
- `invalid targeting "X"` — unknown targeting method

---

### mdflag_list

List all feature flags in the project.

**Access**: Available to agents

**Parameters**: None

**Example Call**:

```json
{
  "name": "mdflag_list",
  "arguments": {}
}
```

**Example Response**:

```
Found 2 flag(s) in .mdflag:

### new-checkout-flow
- Percentage: 5%
- Status: active
- Targeting: user_id
- Hypothesis: Increase checkout conversion by 5%
- Author: claude-code-agent
- Metrics: conversion_rate, checkout_time_seconds
- Expires: 2026-12-17

### dark-mode
- Percentage: 100%
- Status: active
- Targeting: user_id
- Hypothesis: Improve user satisfaction
- Author: cursor-agent
```

---

### mdflag_verify

Verify integrity of all feature flags.

**Access**: Available to agents

**Parameters**: None

**Example Call**:

```json
{
  "name": "mdflag_verify",
  "arguments": {}
}
```

**Example Response (success)**:

```
[OK]   .mdflag/new-checkout-flow.md
[OK]   .mdflag/dark-mode.md

Results: 2 valid, 0 invalid, 2 total

All flags are intact. No integrity violations detected.
```

**Example Response (violation)**:

```
[OK]   .mdflag/new-checkout-flow.md
[FAIL] .mdflag/dark-mode.md
       - human section modified: expected sha256:abc..., computed sha256:def...

Results: 1 valid, 1 invalid, 2 total

WARNING: Integrity violations detected. Flag files may have been modified outside of MDFLAG tools.
Do not trust the affected flags until the issue is resolved.
```

---

## Tools NOT Available to Agents

The following operations are intentionally excluded from MCP tools and are only available through CLI:

| Operation | CLI Command | Reason |
|-----------|-------------|--------|
| Change percentage | `mdflag rollout --name X --percentage N` | Only humans should enable features |
| Delete flag | Manual file deletion | Prevents accidental data loss |
| Modify status | Manual file edit + hash recalculation | Prevents unauthorized state changes |

This access model ensures:
- Agents can propose experiments (create flags)
- Humans control rollout (change percentages)
- Integrity is verifiable at any time (verify)

---

## Agent Integration Best Practices

### For Agent Developers

Add these instructions to your agent's system prompt:

```
When making experimental code changes:

1. BEFORE writing code, create a feature flag:
   - Use mdflag_create with a descriptive name and hypothesis
   - Set percentage to 0 (disabled by default)
   - Include relevant metrics to track

2. Write your code with flag checks:
   - Import the MDFLAG runtime library
   - Wrap new code in: if client.Enabled("flag-name", entityID) { ... }
   - Keep old code in the else branch

3. Commit everything together:
   - Flag file (.mdflag/your-flag.md)
   - Code changes
   - Both in the same commit

4. NEVER:
   - Modify existing flag files directly
   - Change flag percentages
   - Delete flags
   - Bypass integrity checks
```

### For Human Reviewers

When reviewing PRs that include MDFLAG files:

1. Verify the flag file exists and is valid: `mdflag verify`
2. Check the hypothesis is clear and testable
3. Review the metrics are appropriate
4. Ensure percentage starts at 0
5. After merging, gradually roll out:
   ```bash
   mdflag rollout --name <flag> --percentage 5
   # Monitor for 1-2 days
   mdflag rollout --name <flag> --percentage 25
   # Monitor for 3-5 days
   mdflag rollout --name <flag> --percentage 100
   ```
