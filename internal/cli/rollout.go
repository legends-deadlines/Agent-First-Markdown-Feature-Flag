package cli

import (
	"flag"
	"fmt"

	mdflag "github.com/legends-deadlines/mdflag/internal/flag"
)

// Rollout обрабатывает команду `mdflag rollout <args>`
func Rollout(args []string) error {
	fs := flag.NewFlagSet("rollout", flag.ContinueOnError)
	name := fs.String("name", "", "Имя флага (обязательное)")
	percentage := fs.Int("percentage", -1, "Новый процент включения (0-100)")
	dir := fs.String("dir", ".mdflag", "Директория хранения флагов")

	if err := fs.Parse(args); err != nil {
		return err
	}

	if *name == "" {
		return fmt.Errorf("параметр --name обязателен")
	}
	if *percentage < 0 || *percentage > 100 {
		return fmt.Errorf("параметр --percentage должен быть от 0 до 100")
	}

	st, err := mdflag.NewStore(*dir)
	if err != nil {
		return fmt.Errorf("ошибка инициализации хранилища: %w", err)
	}

	if err := st.UpdatePercentage(*name, *percentage); err != nil {
		return fmt.Errorf("не удалось обновить раскат: %w", err)
	}

	fmt.Printf("Процент включения для '%s' обновлен до %d%%\n", *name, *percentage)
	return nil
}
