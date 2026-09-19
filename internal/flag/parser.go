package flag

import (
	"bytes"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"gopkg.in/yaml.v3"
)

// ParseFile читает файл флага с диска и парсит его
func ParseFile(path string) (*Flag, error) {
	// Защита от path traversal: проверяем, что путь не содержит ..
	cleanPath := filepath.Clean(path)
	if strings.Contains(cleanPath, "..") {
		return nil, fmt.Errorf("path traversal detected: %s", path)
	}

	data, err := os.ReadFile(cleanPath)
	if err != nil {
		return nil, fmt.Errorf("read flag file: %w", err)
	}

	return Parse(data)
}

// Parse разбирает байтовое содержимое файла флага
func Parse(data []byte) (*Flag, error) {
	// Извлекаем фронтматтер и тело
	frontmatter, body, err := splitFrontmatter(data)
	if err != nil {
		return nil, fmt.Errorf("split frontmatter: %w", err)
	}

	// Парсим YAML фронтматтер
	var meta FlagMeta
	if err := yaml.Unmarshal(frontmatter, &meta); err != nil {
		return nil, fmt.Errorf("parse frontmatter yaml: %w", err)
	}

	// Валидация обязательных полей
	if meta.Name == "" {
		return nil, fmt.Errorf("flag name is required")
	}
	if meta.Percentage < 0 || meta.Percentage > 100 {
		return nil, fmt.Errorf("percentage must be 0-100, got %d", meta.Percentage)
	}
	if !IsValidStatus(meta.Status) {
		return nil, fmt.Errorf("invalid status: %s", meta.Status)
	}

	return &Flag{
		Meta: meta,
		Body: string(body),
	}, nil
}

// splitFrontmatter извлекает YAML-блок и Markdown-тело из файла.
// Формат: файл начинается с "---", затем YAML, затем "---", затем тело.
func splitFrontmatter(data []byte) ([]byte, []byte, error) {
	// Normalization for line endings
	data = bytes.ReplaceAll(data, []byte("\r\n"), []byte("\n"))

	// Проверяем, что файл начинается с "---"
	if !bytes.HasPrefix(data, []byte("---")) {
		return nil, nil, fmt.Errorf("flag file must start with ---")
	}

	// Пропускаем начальный "---" и ищем закрывающий разделитель
	rest := data[3:]

	// Ищем "\n---" как разделитель
	sep := bytes.Index(rest, []byte("\n---"))
	if sep < 0 {
		return nil, nil, fmt.Errorf("missing closing --- for frontmatter")
	}

	frontmatter := rest[:sep]
	// Пропускаем "\n---" и берём остальное как тело
	body := rest[sep+4:]
	if len(body) > 0 && body[0] == '\n' {
		body = body[1:]
	}

	return frontmatter, body, nil
}

// WriteFile записывает флаг обратно в файл с пересчётом хешей.
// Использует атомарную запись: сначала во временный файл, потом rename.
func WriteFile(path string, f *Flag) error {
	// Защита от path traversal
	cleanPath := filepath.Clean(path)
	if strings.Contains(cleanPath, "..") {
		return fmt.Errorf("path traversal detected: %s", path)
	}

	// Пересчитываем хеши перед записью
	f.Meta.AgentSectionHash = ComputeAgentHash(f.Meta, f.Body)
	f.Meta.HumanSectionHash = ComputeHumanHash(f.Meta)

	// Сериализуем мета-данные в YAML
	metaBytes, err := yaml.Marshal(&f.Meta)
	if err != nil {
		return fmt.Errorf("marshal meta: %w", err)
	}

	// Собираем содержимое файла
	var buf bytes.Buffer
	buf.WriteString("---\n")
	buf.Write(metaBytes)
	buf.WriteString("---\n")
	buf.WriteString(f.Body)

	// Атомарная запись: создаём временный файл, потом переименовываем
	tmpPath := cleanPath + ".tmp"
	if err := os.WriteFile(tmpPath, buf.Bytes(), 0644); err != nil {
		return fmt.Errorf("write temp file: %w", err)
	}
	if err := os.Rename(tmpPath, cleanPath); err != nil {
		// Пытаемся удалить временный файл при ошибке rename
		_ = os.Remove(tmpPath)
		return fmt.Errorf("rename temp file: %w", err)
	}

	return nil
}
