package mdflag

import (
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/fsnotify/fsnotify"
)

// StartWatcher запускает наблюдение за директорией флагов.
// При изменении файла флага автоматически вызывает Reload.
//
// Использование:
//
//	client, _ := mdflag.New(".mdflag")
//	stop := client.StartWatcher()
//	defer stop()
//
// Возвращает функцию для остановки наблюдателя.
func (c *Client) StartWatcher() (func(), error) {
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		return nil, fmt.Errorf("create fsnotify watcher: %w", err)
	}

	// Наблюдаем за директорией флагов
	if err := watcher.Add(c.flagsDir); err != nil {
		_ = watcher.Close()
		return nil, fmt.Errorf("watch directory %s: %w", c.flagsDir, err)
	}

	// Канал для остановки
	stopCh := make(chan struct{})

	// Горутина для обработки событий
	go func() {
		// Debounce: если файлы меняются часто, перезагружаем один раз в 500мс
		var debounceTimer *time.Timer
		debounceInterval := 500 * time.Millisecond

		for {
			select {
			case event, ok := <-watcher.Events:
				if !ok {
					return
				}

				// Нас интересуют только файлы .md
				if filepath.Ext(event.Name) != ".md" {
					continue
				}

				// Интересуют события записи, создания, удаления
				if event.Op&(fsnotify.Write|fsnotify.Create|fsnotify.Remove|fsnotify.Rename) != 0 {
					// Debounce: откладываем перезагрузку
					if debounceTimer != nil {
						debounceTimer.Stop()
					}
					debounceTimer = time.AfterFunc(debounceInterval, func() {
						if err := c.Reload(); err != nil {
							fmt.Fprintf(os.Stderr, "mdflag: reload error: %v\n", err)
						}
					})
				}

			case err, ok := <-watcher.Errors:
				if !ok {
					return
				}
				fmt.Fprintf(os.Stderr, "mdflag: watcher error: %v\n", err)

			case <-stopCh:
				if debounceTimer != nil {
					debounceTimer.Stop()
				}
				_ = watcher.Close()
				return
			}
		}
	}()

	// Возвращаем функцию остановки
	stop := func() {
		close(stopCh)
	}

	return stop, nil
}
