package flag

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"sort"
	"strings"
	"time"
)

// ComputeAgentHash вычисляет хеш секции, которую задаёт агент при создании флага.
// Включает: имя, гипотезу, метрики и Markdown-тело.
// Метрики сортируются для канонической формы.
func ComputeAgentHash(meta FlagMeta, body string) string {
	// Копируем и сортируем метрики для стабильности хеша
	metrics := make([]string, len(meta.Metrics))
	copy(metrics, meta.Metrics)
	sort.Strings(metrics)

	// Собираем каноническую строку в фиксированном порядке полей
	canonical := fmt.Sprintf(
		"name=%s|hypothesis=%s|metrics=%s|body=%s",
		meta.Name,
		meta.Hypothesis,
		strings.Join(metrics, ","),
		body,
	)

	return hashString(canonical)
}

// ComputeHumanHash вычисляет хеш секции управления, которую может менять только человек.
// Включает: процент, статус, таргетинг и срок действия.
func ComputeHumanHash(meta FlagMeta) string {
	// Форматируем время в UTC для канонической формы
	var expiresStr string
	if meta.Expires.IsZero() {
		expiresStr = ""
	} else {
		expiresStr = meta.Expires.UTC().Format(time.RFC3339)
	}

	canonical := fmt.Sprintf(
		"percentage=%d|status=%s|targeting=%s|expires=%s",
		meta.Percentage,
		meta.Status,
		meta.Targeting,
		expiresStr,
	)

	return hashString(canonical)
}

// hashString возвращает строку формата "sha256:<hex>"
func hashString(input string) string {
	h := sha256.Sum256([]byte(input))
	return "sha256:" + hex.EncodeToString(h[:])
}
