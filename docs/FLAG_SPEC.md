# Спецификация формата файла фича-флага (.mdflag/*.md)

## Обзор

Фича-флаги хранятся в централизованной директории `.mdflag/` в корне репозитория. Файл флага разделен на четко разграниченные секции для безопасной совместной работы человека и AI-агента.

## Структура файла флага

Файл содержит три основные секции:

1. **`AGENT_WRITABLE`**: Секция, доступная для заполнения AI-агентом при создании флага (имя, описание, гипотеза, предлагаемые ограничения).
2. **`HUMAN_CONTROLLED`**: Секция под управлением человека (статус, процент раскатки, финальные guardrails).
3. **`INTEGRITY`**: Защищенная секция с SHA-256 хешами верхних секций, контролируемая бинарником `mdflag`.

## Пример формата файла

```markdown
## AGENT_WRITABLE
id: new_payment_flow
description: Новый поток оформления заказа
hypothesis: Конверсия вырастет на 5%

## HUMAN_CONTROLLED
status: active
rollout: 25
exclusion_group: checkout_v2
dependencies:
  - user_service_v2
guardrails:
  - "Do NOT modify legacy_payment.go inside this feature."

## INTEGRITY
agent_section_hash: sha256:a1b2c3d4e5f67890123456789abcdef0123456789abcdef0123456789abcdef0
human_section_hash: sha256:e5f6a7b8c9d0123456789abcdef0123456789abcdef0123456789abcdef01234
```

## Таблица полей

| Поле | Секция | Обязательное | Разрешено агенту | Описание |
|------|--------|--------------|------------------|----------|
| `id` | AGENT_WRITABLE | Да | Да (при `create`) | Уникальный ID флага (snake_case) |
| `description` | AGENT_WRITABLE | Да | Да | Описание функционала |
| `hypothesis` | AGENT_WRITABLE | Нет | Да | Проверяемая гипотеза |
| `status` | HUMAN_CONTROLLED | Да | Нет | Состояние: `draft`, `active`, `deprecated` |
| `rollout` | HUMAN_CONTROLLED | Да | Нет | Процент раскатки (0 до 100) |
| `exclusion_group` | HUMAN_CONTROLLED | Нет | Нет | Группа взаимоисключающих флагов |
| `dependencies` | HUMAN_CONTROLLED | Нет | Нет | Список ID зависимых флагов |
| `guardrails` | HUMAN_CONTROLLED | Нет | Нет | Защитные правила для агента |
| `agent_section_hash` | INTEGRITY | Да | Нет | SHA-256 от секции AGENT_WRITABLE |
| `human_section_hash` | INTEGRITY | Да | Нет | SHA-256 от секции HUMAN_CONTROLLED |
