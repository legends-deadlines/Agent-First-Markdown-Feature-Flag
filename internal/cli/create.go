package cli

import (
	"flag"
	"fmt"
	"strings"
	"time"

	mdflag "github.com/legends-deadlines/mdflag/internal/flag"
)

// Create обрабатывает команду `mdflag create <args>`
func Create(args []string) error {
	fs := flag.NewFlagSet("create", flag.ContinueOnError)
	name := fs.String("name", "", "Имя флага (обязательное)")
	percentage := fs.Int("percentage", 0, "Начальный процент включения (0-100)")
	hypothesis := fs.String("hypothesis", "", "Проверяемая гипотеза")
	author := fs.String("author", "agent", "Автор флага")
	targeting := fs.String("targeting", "user_id", "Способ таргетинга (user_id, session_id, random)")
	metricsRaw := fs.String("metrics", "", "Метрики через запятую")
	dir := fs.String("dir", ".mdflag", "Директория хранения флагов")

	if err := fs.Parse(args); err != nil {
		return err
	}

	if *name == "" {
		return fmt.Errorf("параметр --name обязателен")
	}

	if *percentage < 0 || *percentage > 100 {
		return fmt.Errorf("процент включения должен быть от 0 до 100, получено: %d", *percentage)
	}

	st, err := mdflag.NewStore(*dir)
	if err != nil {
		return fmt.Errorf("ошибка инициализации хранилища: %w", err)
	}

	var metrics []string
	if *metricsRaw != "" {
		for _, m := range strings.Split(*metricsRaw, ",") {
			trimmed := strings.TrimSpace(m)
			if trimmed != "" {
				metrics = append(metrics, trimmed)
			}
		}
	}

	meta := mdflag.FlagMeta{
		Name:       *name,
		Percentage: *percentage,
		Targeting:  *targeting,
		Status:     mdflag.StatusActive,
		Created:    time.Now(),
		Author:     *author,
		Expires:    time.Now().AddDate(0, 3, 0),
		Hypothesis: *hypothesis,
		Metrics:    metrics,
	}

	body := fmt.Sprintf("## Description\nФлаг %s создан через CLI.\n", *name)

	flg, err := st.Create(meta, body)
	if err != nil {
		return fmt.Errorf("не удалось создать флаг: %w", err)
	}

	fmt.Printf("Флаг '%s' успешно создан (%s)\n", flg.Meta.Name, *dir)
	return nil
}
