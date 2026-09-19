package cli

import (
	"flag"
	"fmt"

	mdflag "github.com/legends-deadlines/mdflag/internal/flag"
)

// List обрабатывает команду `mdflag list [dir]`
func List(args []string) error {
	fs := flag.NewFlagSet("list", flag.ContinueOnError)
	dir := fs.String("dir", ".mdflag", "Директория хранения флагов")

	if err := fs.Parse(args); err != nil {
		return err
	}

	if fs.NArg() > 0 {
		*dir = fs.Arg(0)
	}

	st, err := mdflag.NewStore(*dir)
	if err != nil {
		return fmt.Errorf("ошибка инициализации хранилища: %w", err)
	}

	flagsList, err := st.List()
	if err != nil {
		return fmt.Errorf("ошибка чтения списка флагов: %w", err)
	}

	if len(flagsList) == 0 {
		fmt.Printf("Флаги в директории '%s' отсутствуют.\n", *dir)
		return nil
	}

	fmt.Printf("%-25s %-10s %-10s %-15s %-20s\n", "NAME", "STATUS", "ROLLOUT", "TARGETING", "AUTHOR")
	fmt.Println("--------------------------------------------------------------------------------")
	for _, flg := range flagsList {
		fmt.Printf("%-25s %-10s %d%%       %-15s %-20s\n",
			flg.Meta.Name,
			flg.Meta.Status,
			flg.Meta.Percentage,
			flg.Meta.Targeting,
			flg.Meta.Author,
		)
	}
	return nil
}
