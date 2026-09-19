# Спецификация MCP-инструментов mdflag

## Доступные инструменты для AI-агентов

### `mdflag_create`

Создаёт новый фича-флаг в папке `.mdflag/`.

#### Параметры:
- `name` (string, обязательное): Уникальный ID флага.
- `hypothesis` (string): Проверяемая гипотеза.
- `description` (string): Человеко-читаемое описание флага и затронутых файлов.
- `author` (string): Имя или ID агента.
- `metrics` (array of strings): Отслеживаемые метрики.

#### Поведение:
Инструмент автоматически ставит `percentage: 0`, `status: draft` и вычисляет базовые хеши `agent_section_hash` и `human_section_hash`.
