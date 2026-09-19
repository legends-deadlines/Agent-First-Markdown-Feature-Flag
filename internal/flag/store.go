package flag

import (
	"fmt"
	"os"
	"path/filepath"
	"time"
)

// Store управляет файлами флагов в директории проекта
type Store struct {
	dir string
}

// NewStore создаёт новое хранилище флагов
func NewStore(dir string) (*Store, error) {
	// Создаём директорию, если её нет
	if err := os.MkdirAll(dir, 0755); err != nil {
		return nil, fmt.Errorf("create flags dir: %w", err)
	}

	return &Store{dir: dir}, nil
}

// Create создаёт новый файл флага.
// Возвращает ошибку, если флаг с таким именем уже существует.
func (s *Store) Create(meta FlagMeta, body string) (*Flag, error) {
	// Проверяем, что флаг не существует
	path := s.flagPath(meta.Name)
	if _, err := os.Stat(path); err == nil {
		return nil, fmt.Errorf("flag %q already exists", meta.Name)
	}

	// Устанавливаем значения по умолчанию
	if meta.Status == "" {
		meta.Status = StatusActive
	}
	if meta.Targeting == "" {
		meta.Targeting = "user_id"
	}
	if meta.Created.IsZero() {
		meta.Created = time.Now()
	}

	f := &Flag{
		Meta: meta,
		Body: body,
	}

	// Записываем файл (хеши пересчитаются автоматически)
	if err := WriteFile(path, f); err != nil {
		return nil, fmt.Errorf("write flag: %w", err)
	}

	return f, nil
}

// Get читает флаг по имени
func (s *Store) Get(name string) (*Flag, error) {
	path := s.flagPath(name)
	return ParseFile(path)
}

// UpdatePercentage меняет процент включения флага.
// Пересчитывает только хеш секции управления.
func (s *Store) UpdatePercentage(name string, percentage int) error {
	if percentage < 0 || percentage > 100 {
		return fmt.Errorf("percentage must be 0-100, got %d", percentage)
	}

	f, err := s.Get(name)
	if err != nil {
		return fmt.Errorf("get flag: %w", err)
	}

	f.Meta.Percentage = percentage

	return WriteFile(s.flagPath(name), f)
}

// List возвращает все флаги в директории
func (s *Store) List() ([]*Flag, error) {
	pattern := filepath.Join(s.dir, "*.md")
	matches, err := filepath.Glob(pattern)
	if err != nil {
		return nil, fmt.Errorf("glob flags dir: %w", err)
	}

	// Pre-allocate слайс нужного размера
	flags := make([]*Flag, 0, len(matches))

	for _, path := range matches {
		f, err := ParseFile(path)
		if err != nil {
			// Логируем ошибку, но пропускаем файл
			fmt.Fprintf(os.Stderr, "skip %s: %v\n", path, err)
			continue
		}
		flags = append(flags, f)
	}

	return flags, nil
}

// flagPath возвращает путь к файлу флага по его имени
func (s *Store) flagPath(name string) string {
    return filepath.Join(s.dir, name+".md")
}
