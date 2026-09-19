package flag

import (
	"fmt"
	"os"
	"path/filepath"
)

type Store struct {
	dir string
}

func NewStore(dir string) (*Store, error) {
	// Создаём директорию, если её нет
	if err := os.MkdirAll(dir, 0755); err != nil {
		return nil, fmt.Errorf("create flags dir: %w", err)
	}

	return &Store{dir: dir}, nil
}

func (s *Store) Create(meta FlagMeta, body string) (*Flag, error) {
	f := &Flag{
		Meta: meta,
		Body: body,
	}
	if err := WriteFile(s.flagPath(meta.Name), f); err != nil {
		return nil, err
	}
	return f, nil
}

func (s *Store) Get(name string) (*Flag, error) {
	return ParseFile(s.flagPath(name))
}

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
			continue
		}
		flags = append(flags, f)
	}

	return flags, nil
}

func (s *Store) flagPath(name string) string {
	return filepath.Join(s.dir, name+".md")
}
