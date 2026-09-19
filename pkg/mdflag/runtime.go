package mdflag

import (
	"fmt"
	"hash/fnv"
	"os"
	"path/filepath"
	"sync"
	"time"

	"github.com/legends-deadlines/mdflag/internal/flag"
)

// Client предоставляет API для проверки флагов в приложении
type Client struct {
	mu       sync.RWMutex
	flags    map[string]*flag.Flag
	flagsDir string
}

// New создаёт новый клиент, читая все флаги из директории
func New(flagsDir string) (*Client, error) {
	c := &Client{
		flags:    make(map[string]*flag.Flag),
		flagsDir: flagsDir,
	}

	if err := c.loadAll(); err != nil {
		return nil, fmt.Errorf("load flags: %w", err)
	}

	return c, nil
}

// Enabled проверяет, включён ли флаг для данной сущности
// entityID — идентификатор для детерминированного таргетинга (например, user_id)
func (c *Client) Enabled(flagName, entityID string) bool {
	c.mu.RLock()
	f, ok := c.flags[flagName]
	c.mu.RUnlock()

	if !ok {
		return false // флаг не найден = выключен
	}

	// Проверяем статус
	if f.Meta.Status != flag.StatusActive {
		return false
	}

	// Проверяем срок действия
	if !f.Meta.Expires.IsZero() && f.Meta.Expires.Before(now()) {
		return false
	}

	// 0% = выключен для всех
	if f.Meta.Percentage == 0 {
		return false
	}

	// 100% = включен для всех
	if f.Meta.Percentage == 100 {
		return true
	}

	// Детерминированный хеш для процента
	return c.inPercentage(flagName, entityID, f.Meta.Percentage)
}

// inPercentage проверяет, попадает ли сущность в процент включения
func (c *Client) inPercentage(flagName, entityID string, percentage int) bool {
	// Используем FNV-1a для скорости и детерминизма
	h := fnv.New32a()
	_, _ = h.Write([]byte(flagName))
	_, _ = h.Write([]byte(":"))
	_, _ = h.Write([]byte(entityID))

	// Хеш -> число 0-99
	bucket := int(h.Sum32() % 100)
	return bucket < percentage
}

// loadAll читает все флаги из директории
func (c *Client) loadAll() error {
	pattern := filepath.Join(c.flagsDir, "*.md")
	matches, err := filepath.Glob(pattern)
	if err != nil {
		return fmt.Errorf("glob flags dir: %w", err)
	}

	// Pre-allocate для производительности
	flagsMap := make(map[string]*flag.Flag, len(matches))

	for _, path := range matches {
		f, err := flag.ParseFile(path)
		if err != nil {
			// Логируем ошибку, но не прерываем загрузку других флагов
			fmt.Fprintf(os.Stderr, "mdflag: skip %s: %v\n", path, err)
			continue
		}

		// Проверяем целостность
		result := flag.Validate(f)
		if !result.AgentHashValid || !result.HumanHashValid {
			fmt.Fprintf(os.Stderr, "mdflag: integrity violation in %s: %v\n", path, result.Errors)
			continue
		}

		flagsMap[f.Meta.Name] = f
	}

	c.mu.Lock()
	c.flags = flagsMap
	c.mu.Unlock()
	return nil
}

// now — функция времени, которую можно переопределить в тестах
var now = func() time.Time { return time.Now() }
