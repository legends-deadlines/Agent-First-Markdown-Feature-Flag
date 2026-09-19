---
name: new-checkout-flow
percentage: 5
targeting: user_id
status: active
created: 2026-09-17T10:00:00Z
author: claude-code-agent
expires: 2026-12-17T00:00:00Z
hypothesis: Конверсия в покупку вырастет на 5%
metrics:
    - conversion_rate
    - checkout_time_seconds
agent_section_hash: "sha256:placeholder"
human_section_hash: "sha256:placeholder"
---

## Description
Новый поток оформления заказа с одним экраном вместо трёх.

## Implementation Notes
Изменения затронули:
- `frontend/src/checkout/OnePageCheckout.tsx`
- `backend/api/checkout.go`
