# Спецификация формата файла флага mdflag

## Формат файла

Файл содержит две части: YAML frontmatter (заголовок между `---`) и Markdown-содержимое.

```markdown
---
name: new-checkout-flow
percentage: 5
targeting: user_id
status: active
created: 2026-09-17T10:00:00Z
author: claude-code-agent
expires: 2026-12-17T00:00:00Z
hypothesis: "Конверсия в покупку вырастет на 5%"
metrics:
  - conversion_rate
  - checkout_time_seconds
agent_section_hash: "sha256:e3b0c44298fc1c149afbf4c8996fb92427ae41e4"
human_section_hash: "sha256:a7ffc6f8bf1ed76651c14756a061d662f580ff4d"
---

# Заголовок флага

## Description
Описание функционала.
```

## Правила расчета хешей

- `agent_section_hash`: Вычисляется SHA-256 от строк полей `name`, `hypothesis`, `metrics` и полного текста Markdown body.
- `human_section_hash`: Вычисляется SHA-256 от полей `percentage`, `status`, `targeting`, `expires`.
- При вызове `mdflag rollout` измеряется только `human_section_hash`.
- При `mdflag verify` пересчитываются оба хеша для проверки целостности.
