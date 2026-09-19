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
