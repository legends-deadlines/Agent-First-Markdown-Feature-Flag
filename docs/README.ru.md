<div align="center">

<img src="../.github/assets/banner.svg" alt="Agent-First Feature Flags" width="100%"/>

<br/>

[![Статус](https://img.shields.io/badge/⚙%20Статус-В%20разработке-f78166?style=flat-square&labelColor=1c2128)](#)
[![Лицензия](https://img.shields.io/badge/Лицензия-MIT-79c0ff?style=flat-square&labelColor=1c2128)](../LICENSE)

**🇬🇧 [Read in English](../README.md)**

</div>

---

> [!WARNING]
> **🚧 Проект в активной разработке.** Концепция зафиксирована, технический стек определяется.

## Проблема

AI-агенты (**Cursor, Claude Code, Copilot**) генерируют код, ничего не зная о фича-флагах. Они не видят, какие фичи раскатаны на 10%, какие флаги взаимоисключающие, какие guardrails существуют.

Результат: **28–40% лишнего code churn**, **46.4% AI-PR отклоняются**, **в 1.7× больше дефектов**.

Существующие системы флагов (LaunchDarkly, Flipt, dif.sh) сделаны для людей, которые кликают по дашбордам — не для автономных агентов, читающих репозиторий.

## Решение

**Фича-флаги как Markdown-файлы прямо в Git-репозитории.** Формат, одинаково понятный человеку и AI-агенту.

```
flags/*.md ──► CLI-компилятор ──► .cursorrules, CLAUDE.md, context.json
                    ▲
             AST-анализ кодовой базы
```

CLI сканирует определения флагов, анализирует где в коде они используются, находит конфликты и зависимости, затем генерирует правила, которым AI-агенты следуют автоматически.

## Как выглядит флаг

```markdown
---
id: new_payment_flow
status: active
rollout: 25
exclusion_group: checkout_v2
dependencies:
  - user_service_v2
guardrails:
  - "НЕ трогать legacy_payment.go в рамках этой фичи."
---

# New Payment Flow

Постепенная раскатка обновлённого платёжного пайплайна.
```

После `dif compile` AI-агент видит:

```
ACTIVE: new_payment_flow (25%), user_service_v2 (100%)
INACTIVE: old_checkout, legacy_payment
RULE: new_payment_flow ⊕ old_checkout (взаимоисключающие)
RULE: НЕ трогать legacy_payment.go в рамках new_payment_flow
```

## Почему не просто X?

|  | **Этот проект** | dif.sh | Flipt | LaunchDarkly |
|:---|:---:|:---:|:---:|:---:|
| Флаги живут в Git | ✅ | ✅ | ❌ | ❌ |
| AI-агент читает флаги автономно | ✅ | ❌ | ❌ | ❌ |
| Автогенерация guardrails для агентов | ✅ | ❌ | ❌ | ❌ |
| Мультиязычный анализ кодовой базы | ✅ | ❌ | ❌ | ❌ |
| Оффлайн-таргетинг без SDK | ✅ | ❌ | ❌ | ❌ |

## Дорожная карта

- [x] Концепция и архитектура
- [ ] Парсер флагов и валидация схемы
- [ ] AST-сканер кодовой базы
- [ ] Граф зависимостей и детекция конфликтов
- [ ] Генератор правил для агентов (`.cursorrules`, `CLAUDE.md`)
- [ ] Runtime-движок таргетинга
- [ ] CLI и CI/CD интеграция

## Лицензия

MIT

---

<div align="center">

<img src="../.github/assets/seal.svg" alt="Тюлень-маскот" width="600"/>

<sub>by <a href="https://github.com/legends-deadlines">legends-deadlines</a></sub>

</div>
