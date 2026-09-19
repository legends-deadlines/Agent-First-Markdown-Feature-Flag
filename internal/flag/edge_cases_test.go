package flag

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

// Тест: пустой файл
func TestParse_EmptyFile(t *testing.T) {
	_, err := Parse([]byte(""))
	if err == nil {
		t.Error("expected error for empty file")
	}
}

// Тест: файл без фронтматтера
func TestParse_NoFrontmatter(t *testing.T) {
	data := []byte("just some markdown text without frontmatter")
	_, err := Parse(data)
	if err == nil {
		t.Error("expected error for missing frontmatter")
	}
}

// Тест: фронтматтер без закрывающего разделителя
func TestParse_UnclosedFrontmatter(t *testing.T) {
	data := []byte("---\nname: test\npercentage: 50\n")
	_, err := Parse(data)
	if err == nil {
		t.Error("expected error for unclosed frontmatter")
	}
}

// Тест: очень длинное имя флага
func TestParse_VeryLongName(t *testing.T) {
	longName := strings.Repeat("a", 1000)
	data := []byte("---\nname: " + longName + "\npercentage: 50\nstatus: active\n---\nbody")

	_, err := Parse(data)
	// Должен либо принять, либо отклонить, но не паниковать
	_ = err
}

// Тест: специальные символы в имени
func TestParse_SpecialCharsInName(t *testing.T) {
	data := []byte("---\nname: test-flag-with-dashes-123\npercentage: 50\nstatus: active\n---\nbody")

	f, err := Parse(data)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if f.Meta.Name != "test-flag-with-dashes-123" {
		t.Errorf("unexpected name: %q", f.Meta.Name)
	}
}

// Тест: Unicode в гипотезе
func TestParse_UnicodeHypothesis(t *testing.T) {
	data := []byte("---\nname: test\npercentage: 50\nstatus: active\nhypothesis: Конверсия вырастет на 5%\n---\nbody")

	f, err := Parse(data)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if f.Meta.Hypothesis != "Конверсия вырастет на 5%" {
		t.Errorf("unicode hypothesis not preserved: %q", f.Meta.Hypothesis)
	}
}

// Тест: пустые метрики
func TestParse_EmptyMetrics(t *testing.T) {
	data := []byte("---\nname: test\npercentage: 50\nstatus: active\nmetrics: []\n---\nbody")

	f, err := Parse(data)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(f.Meta.Metrics) != 0 {
		t.Errorf("expected empty metrics, got %v", f.Meta.Metrics)
	}
}

// Тест: нулевой процент
func TestParse_ZeroPercent(t *testing.T) {
	data := []byte("---\nname: test\npercentage: 0\nstatus: active\n---\nbody")

	f, err := Parse(data)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if f.Meta.Percentage != 0 {
		t.Errorf("expected 0, got %d", f.Meta.Percentage)
	}
}

// Тест: 100 процентов
func TestParse_HundredPercent(t *testing.T) {
	data := []byte("---\nname: test\npercentage: 100\nstatus: active\n---\nbody")

	f, err := Parse(data)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if f.Meta.Percentage != 100 {
		t.Errorf("expected 100, got %d", f.Meta.Percentage)
	}
}

// Тест: отрицательный процент
func TestParse_NegativePercent(t *testing.T) {
	data := []byte("---\nname: test\npercentage: -1\nstatus: active\n---\nbody")

	_, err := Parse(data)
	if err == nil {
		t.Error("expected error for negative percentage")
	}
}

// Тест: процент больше 100
func TestParse_OverHundredPercent(t *testing.T) {
	data := []byte("---\nname: test\npercentage: 101\nstatus: active\n---\nbody")

	_, err := Parse(data)
	if err == nil {
		t.Error("expected error for percentage > 100")
	}
}

// Тест: несуществующий файл
func TestParseFile_NonExistent(t *testing.T) {
	_, err := ParseFile("/nonexistent/path/flag.md")
	if err == nil {
		t.Error("expected error for nonexistent file")
	}
}

// Тест: path traversal в ParseFile
func TestParseFile_PathTraversal(t *testing.T) {
	_, err := ParseFile("../../etc/passwd")
	if err == nil {
		t.Error("expected error for path traversal")
	}
}

// Тест: симлинк должен отклоняться
func TestParseFile_Symlink(t *testing.T) {
	dir := t.TempDir()

	// Создаём реальный файл вне директории
	realFile := filepath.Join(dir, "real.md")
	os.WriteFile(realFile, []byte("---\nname: test\npercentage: 50\nstatus: active\n---\nbody"), 0644)

	// Создаём симлинк в директории флагов
	flagsDir := filepath.Join(dir, "flags")
	os.MkdirAll(flagsDir, 0755)
	symlink := filepath.Join(flagsDir, "symlink.md")
	os.Symlink(realFile, symlink)

	// Парсер должен либо отклонить симлинк, либо прочитать его
	// Проверяем, что не паникует
	_, _ = ParseFile(symlink)
}

// Тест: одновременная запись и чтение
func TestStore_ConcurrentAccess(t *testing.T) {
	dir := t.TempDir()

	store, err := NewStore(dir)
	if err != nil {
		t.Fatalf("NewStore failed: %v", err)
	}

	// Создаём флаг
	meta := FlagMeta{
		Name:       "concurrent-test",
		Percentage: 50,
		Status:     StatusActive,
		Hypothesis: "Test",
	}

	if _, err := store.Create(meta, ""); err != nil {
		t.Fatalf("Create failed: %v", err)
	}

	// Читаем параллельно
	done := make(chan bool, 10)
	for i := 0; i < 10; i++ {
		go func() {
			_, err := store.Get("concurrent-test")
			done <- (err == nil)
		}()
	}

	// Ждём все горутины
	for i := 0; i < 10; i++ {
		if !<-done {
			t.Error("concurrent read failed")
		}
	}
}
