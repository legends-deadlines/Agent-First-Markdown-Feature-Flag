package flag

import "time"

// Flag описывает полный файл флага: мета-данные + Markdown-тело
type Flag struct {
	// Meta содержит машинно-читаемые поля из YAML фронтматтера
	Meta FlagMeta

	// Body содержит человекочитаемый Markdown-текст
	Body string
}

// FlagMeta содержит все поля из YAML фронтматтера
type FlagMeta struct {
	// Имя флага (обязательное, уникальное)
	Name string `yaml:"name"`

	// Процент включения: 0 = выключен, 100 = включён для всех
	Percentage int `yaml:"percentage"`

	// Способ таргетинга: "user_id", "session_id", "random"
	Targeting string `yaml:"targeting"`

	// Статус флага: "active", "paused", "archived"
	Status string `yaml:"status"`

	// Дата и время создания флага
	Created time.Time `yaml:"created"`

	// Автор флага: имя агента или человека
	Author string `yaml:"author"`

	// Срок действия: после этой даты флаг считается выключенным
	Expires time.Time `yaml:"expires"`

	// Гипотеза: что должен показать эксперимент
	Hypothesis string `yaml:"hypothesis"`

	// Метрики для отслеживания успеха
	Metrics []string `yaml:"metrics"`

	// Хеш агентской секции для проверки целостности
	AgentSectionHash string `yaml:"agent_section_hash"`

	// Хеш секции управления для проверки целостности
	HumanSectionHash string `yaml:"human_section_hash"`
}

// Константы для валидных значений статуса
const (
	StatusActive   = "active"
	StatusPaused   = "paused"
	StatusArchived = "archived"
)

// IsValidStatus проверяет, является ли статус допустимым
func IsValidStatus(s string) bool {
	switch s {
	case StatusActive, StatusPaused, StatusArchived:
		return true
	default:
		return false
	}
}

// IsValidTargeting проверяет допустимость способа таргетинга
func IsValidTargeting(t string) bool {
	switch t {
	case "user_id", "session_id", "random":
		return true
	default:
		return false
	}
}
