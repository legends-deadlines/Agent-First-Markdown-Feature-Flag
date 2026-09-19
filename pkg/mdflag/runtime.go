// Пакет mdflag предоставляет публичную библиотеку для проверки
// флагов в рантайме приложения.
//
// Базовое использование:
//
//	client, err := mdflag.New(".mdflag")
//	if err != nil {
//	    log.Fatal(err)
//	}
//
//	if client.Enabled("new-checkout", userID) {
//	    // Новая ветка кода
//	} else {
//	    // Старая ветка кода
//	}
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

// Client предоставляет API для проверки флагов в приложении.
// Потокобезопасен: можно использовать из нескольких горутин.
type Client struct {
	mu       sync.RWMutex
	flags    map[string]*flag.Flag
	flagsDir string

	// nowFunc позволяет переопределить время в тестах
	nowFunc func() time.Time
}

// Option задаёт опциональную настройку клиента
type Option func(*Client)

// WithNowFunc переопределяет функцию времени (для тестов)
func WithNowFunc(f func() time.Time) Option {
	return func(c *Client) {
		c.nowFunc = f
	}
}

// New создаёт новый клиент, читая все флаги из указанной директории.
// Возвращает ошибку, если директория не существует или не содержит валидных флагов.
func New(flagsDir string, opts ...Option) (*Client, error) {
	c := &Client{
		flags:    make(map[string]*flag.Flag),
		flagsDir: flagsDir,
		nowFunc:  time.Now,
	}

	// Применяем опции
	for _, opt := range opts {
		opt(c)
	}

	// Загружаем все флаги
	if err := c.loadAll(); err != nil {
		return nil, fmt.Errorf("load flags from %s: %w", flagsDir, err)
	}

	return c, nil
}

// Enabled проверяет, включён ли флаг для данной сущности.
//
// Параметры:
//   - flagName: имя флага
//   - entityID: идентификатор сущности для таргетинга (например, user_id или session_id)
//
// Возвращает true, если флаг включён для данной сущности, иначе false.
// Если флаг не найден, возвращает false (безопасное поведение по умолчанию).
func (c *Client) Enabled(flagName, entityID string) bool {
	c.mu.RLock()
	f, ok := c.flags[flagName]
	c.mu.RUnlock()

	// Флаг не найден = выключен (безопасное поведение)
	if !ok {
		return false
	}

	// Проверяем статус
	if f.Meta.Status != flag.StatusActive {
		return false
	}

	// Проверяем срок действия
	if !f.Meta.Expires.IsZero() && f.Meta.Expires.Before(c.nowFunc()) {
		return false
	}

	// 0% = выключен для всех
	if f.Meta.Percentage == 0 {
		return false
	}

	// 100% = включён для всех
	if f.Meta.Percentage == 100 {
		return true
	}

	// Детерминированный хеш для процентного таргетинга
	return c.inPercentage(flagName, entityID, f.Meta.Percentage)
}

// Get возвращает мета-данные флага по имени.
// Возвращает ошибку, если флаг не найден.
func (c *Client) Get(flagName string) (*flag.Flag, error) {
	c.mu.RLock()
	defer c.mu.RUnlock()

	f, ok := c.flags[flagName]
	if !ok {
		return nil, fmt.Errorf("flag %q not found", flagName)
	}

	return f, nil
}

// List возвращает имена всех загруженных флагов
func (c *Client) List() []string {
	c.mu.RLock()
	defer c.mu.RUnlock()

	names := make([]string, 0, len(c.flags))
	for name := range c.flags {
		names = append(names, name)
	}
	return names
}

// Reload перечитывает все флаги из директории.
// Вызывается при изменении файлов (если включён file watcher)
// или вручную.
func (c *Client) Reload() error {
	return c.loadAll()
}

// inPercentage проверяет, попадает ли сущность в процент включения.
// Использует FNV-1a хеш для детерминизма: одна и та же пара
// (flagName, entityID) всегда даёт один и тот же результат.
func (c *Client) inPercentage(flagName, entityID string, percentage int) bool {
	// FNV-1a: быстрый, детерминированный, равномерный
	h := fnv.New32a()

	// Записываем имя флага
	_, _ = h.Write([]byte(flagName))
	// Разделитель для предотвращения коллизий
	_, _ = h.Write([]byte(":"))
	// Записываем идентификатор сущности
	_, _ = h.Write([]byte(entityID))

	// Хеш -> число 0-99 (бакет)
	bucket := int(h.Sum32() % 100)

	// Если бакет меньше процента, сущность включена
	return bucket < percentage
}

// loadAll читает все флаги из директории и обновляет кэш.
// Флаги с ошибками парсинга или нарушенной целостностью пропускаются.
func (c *Client) loadAll() error {
	pattern := filepath.Join(c.flagsDir, "*.md")
	matches, err := filepath.Glob(pattern)
	if err != nil {
		return fmt.Errorf("glob flags dir: %w", err)
	}

	// Pre-allocate map нужного размера для избежания реаллокаций
	newFlags := make(map[string]*flag.Flag, len(matches))

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

		newFlags[f.Meta.Name] = f
	}

	// Атомарно заменяем кэш
	c.mu.Lock()
	c.flags = newFlags
	c.mu.Unlock()

	return nil
}
