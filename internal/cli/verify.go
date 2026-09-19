package cli

import (
	"flag"
	"fmt"

	mdflag "github.com/legends-deadlines/mdflag/internal/flag"
)

// Verify обрабатывает команду `mdflag verify [dir]`
func Verify(args []string) error {
	fs := flag.NewFlagSet("verify", flag.ContinueOnError)
	dir := fs.String("dir", ".mdflag", "Директория хранения флагов")

	if err := fs.Parse(args); err != nil {
		return err
	}

	// Если передан неименованный аргумент пути
	if fs.NArg() > 0 {
		*dir = fs.Arg(0)
	}

	results, err := mdflag.ValidateDir(*dir)
	if err != nil {
		return fmt.Errorf("ошибка верификации директории: %w", err)
	}

	if len(results) == 0 {
		fmt.Printf("Файлы флагов в директории '%s' не найдены.\n", *dir)
		return nil
	}

	hasErrors := false
	for path, res := range results {
		if !res.IsValid() {
			hasErrors = true
			fmt.Printf("[FAIL] %s:\n", path)
			for _, errStr := range res.Errors {
				fmt.Printf("  - %s\n", errStr)
			}
		} else {
			fmt.Printf("[OK] %s — хеши валидны\n", path)
		}
	}

	if hasErrors {
		return fmt.Errorf("проверка целостности завершилась с ошибками")
	}

	fmt.Println("Все файлы флагов успешно прошли проверку целостности.")
	return nil
}
