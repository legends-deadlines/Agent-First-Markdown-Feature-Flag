package flag

import (
	"fmt"
	"path/filepath"
)

// ValidateResult содержит результат проверки целостности файла флага
type ValidateResult struct {
	// AgentHashValid true, если агентская секция не была изменена
	AgentHashValid bool

	// HumanHashValid true, если секция управления не была изменена
	HumanHashValid bool

	// Errors содержит список обнаруженных проблем
	Errors []string
}

// IsValid возвращает true, если все проверки прошли успешно
func (r ValidateResult) IsValid() bool {
	return r.AgentHashValid && r.HumanHashValid && len(r.Errors) == 0
}

// Validate проверяет целостность файла флага, пересчитывая хеши
// и сравнивая их с записанными в файле.
func Validate(f *Flag) ValidateResult {
	result := ValidateResult{
		AgentHashValid: true,
		HumanHashValid: true,
	}

	// Проверяем обязательные поля
	if f.Meta.Name == "" {
		result.Errors = append(result.Errors, "flag name is required")
	}
	if f.Meta.Percentage < 0 || f.Meta.Percentage > 100 {
		result.Errors = append(result.Errors, fmt.Sprintf("percentage out of range: %d", f.Meta.Percentage))
	}
	if !IsValidStatus(f.Meta.Status) {
		result.Errors = append(result.Errors, fmt.Sprintf("invalid status: %s", f.Meta.Status))
	}

	// Проверяем хеш агентской секции
	computedAgent := ComputeAgentHash(f.Meta, f.Body)
	if computedAgent != f.Meta.AgentSectionHash {
		result.AgentHashValid = false
		result.Errors = append(result.Errors, fmt.Sprintf(
			"agent section modified: expected %s, computed %s",
			f.Meta.AgentSectionHash, computedAgent,
		))
	}

	// Проверяем хеш секции управления
	computedHuman := ComputeHumanHash(f.Meta)
	if computedHuman != f.Meta.HumanSectionHash {
		result.HumanHashValid = false
		result.Errors = append(result.Errors, fmt.Sprintf(
			"human section modified: expected %s, computed %s",
			f.Meta.HumanSectionHash, computedHuman,
		))
	}

	return result
}

// ValidateDir проверяет все файлы флагов в директории
func ValidateDir(dir string) (map[string]ValidateResult, error) {
	pattern := filepath.Join(dir, "*.md")
	matches, err := filepath.Glob(pattern)
	if err != nil {
		return nil, fmt.Errorf("glob flags dir: %w", err)
	}

	results := make(map[string]ValidateResult, len(matches))

	for _, path := range matches {
		f, err := ParseFile(path)
		if err != nil {
			results[path] = ValidateResult{
				Errors: []string{fmt.Sprintf("parse error: %v", err)},
			}
			continue
		}

		results[path] = Validate(f)
	}

	return results, nil
}
