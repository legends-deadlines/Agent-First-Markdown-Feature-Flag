package mdflag

import (
	"fmt"
	"testing"
	"time"

	"github.com/legends-deadlines/mdflag/internal/flag"
)

// setupTestFlag создаёт тестовый флаг в указанной директории
func setupTestFlag(t *testing.T, dir, name string, percentage int, status string) {
	t.Helper()

	meta := flag.FlagMeta{
		Name:       name,
		Percentage: percentage,
		Targeting:  "user_id",
		Status:     status,
		Created:    time.Now(),
		Author:     "test",
		Hypothesis: "Test hypothesis",
		Metrics:    []string{"test_metric"},
	}

	store, err := flag.NewStore(dir)
	if err != nil {
		t.Fatalf("NewStore failed: %v", err)
	}

	_, err = store.Create(meta, "## Description\nTest")
	if err != nil {
		t.Fatalf("Create flag failed: %v", err)
	}
}

func TestNew_LoadsFlags(t *testing.T) {
	dir := t.TempDir()

	setupTestFlag(t, dir, "test-flag", 50, flag.StatusActive)

	client, err := New(dir)
	if err != nil {
		t.Fatalf("New failed: %v", err)
	}

	names := client.List()
	if len(names) != 1 {
		t.Errorf("expected 1 flag, got %d", len(names))
	}
	if names[0] != "test-flag" {
		t.Errorf("expected flag name 'test-flag', got %q", names[0])
	}
}

func TestEnabled_FlagNotFound(t *testing.T) {
	dir := t.TempDir()

	client, err := New(dir)
	if err != nil {
		t.Fatalf("New failed: %v", err)
	}

	// Несуществующий флаг должен возвращать false
	if client.Enabled("nonexistent", "user123") {
		t.Error("nonexistent flag should return false")
	}
}

func TestEnabled_ZeroPercent(t *testing.T) {
	dir := t.TempDir()

	setupTestFlag(t, dir, "zero-flag", 0, flag.StatusActive)

	client, err := New(dir)
	if err != nil {
		t.Fatalf("New failed: %v", err)
	}

	// 0% = выключен для всех
	for _, uid := range []string{"user1", "user2", "user3"} {
		if client.Enabled("zero-flag", uid) {
			t.Errorf("zero percent flag should be disabled for %s", uid)
		}
	}
}

func TestEnabled_HundredPercent(t *testing.T) {
	dir := t.TempDir()

	setupTestFlag(t, dir, "full-flag", 100, flag.StatusActive)

	client, err := New(dir)
	if err != nil {
		t.Fatalf("New failed: %v", err)
	}

	// 100% = включён для всех
	for _, uid := range []string{"user1", "user2", "user3"} {
		if !client.Enabled("full-flag", uid) {
			t.Errorf("100 percent flag should be enabled for %s", uid)
		}
	}
}

func TestEnabled_PausedStatus(t *testing.T) {
	dir := t.TempDir()

	setupTestFlag(t, dir, "paused-flag", 100, flag.StatusPaused)

	client, err := New(dir)
	if err != nil {
		t.Fatalf("New failed: %v", err)
	}

	// Пауза = выключен, даже если процент 100
	if client.Enabled("paused-flag", "user123") {
		t.Error("paused flag should be disabled")
	}
}

func TestEnabled_ExpiredFlag(t *testing.T) {
	dir := t.TempDir()

	// Создаём флаг с истёкшим сроком действия
	meta := flag.FlagMeta{
		Name:       "expired-flag",
		Percentage: 100,
		Targeting:  "user_id",
		Status:     flag.StatusActive,
		Created:    time.Now().Add(-24 * time.Hour),
		Author:     "test",
		Hypothesis: "Test",
		Expires:    time.Now().Add(-1 * time.Hour), // Истёк час назад
	}

	store, err := flag.NewStore(dir)
	if err != nil {
		t.Fatalf("NewStore failed: %v", err)
	}

	_, err = store.Create(meta, "Test")
	if err != nil {
		t.Fatalf("Create flag failed: %v", err)
	}

	client, err := New(dir)
	if err != nil {
		t.Fatalf("New failed: %v", err)
	}

	// Истёкший флаг должен быть выключен
	if client.Enabled("expired-flag", "user123") {
		t.Error("expired flag should be disabled")
	}
}

func TestEnabled_Deterministic(t *testing.T) {
	dir := t.TempDir()

	setupTestFlag(t, dir, "det-flag", 50, flag.StatusActive)

	client, err := New(dir)
	if err != nil {
		t.Fatalf("New failed: %v", err)
	}

	// Одна и та же сущность должна всегда получать один результат
	for i := 0; i < 100; i++ {
		result1 := client.Enabled("det-flag", "user123")
		result2 := client.Enabled("det-flag", "user123")
		if result1 != result2 {
			t.Error("Enabled should be deterministic for same entity")
		}
	}
}

func TestEnabled_Distribution(t *testing.T) {
	dir := t.TempDir()

	setupTestFlag(t, dir, "dist-flag", 50, flag.StatusActive)

	client, err := New(dir)
	if err != nil {
		t.Fatalf("New failed: %v", err)
	}

	// Проверяем распределение на большом количестве сущностей
	enabledCount := 0
	totalEntities := 10000

	for i := 0; i < totalEntities; i++ {
		entityID := fmt.Sprintf("user%d", i)
		if client.Enabled("dist-flag", entityID) {
			enabledCount++
		}
	}

	// Ожидаем примерно 50% (допустимое отклонение 10%)
	ratio := float64(enabledCount) / float64(totalEntities)
	if ratio < 0.4 || ratio > 0.6 {
		t.Errorf("expected ~50%% enabled, got %.2f%%", ratio*100)
	}
}

func TestReload_PicksUpChanges(t *testing.T) {
	dir := t.TempDir()

	setupTestFlag(t, dir, "reload-flag", 0, flag.StatusActive)

	client, err := New(dir)
	if err != nil {
		t.Fatalf("New failed: %v", err)
	}

	// Изначально выключен (0%)
	if client.Enabled("reload-flag", "user1") {
		t.Error("flag should be disabled at 0%")
	}

	// Меняем процент через хранилище
	store, err := flag.NewStore(dir)
	if err != nil {
		t.Fatalf("NewStore failed: %v", err)
	}
	if err := store.UpdatePercentage("reload-flag", 100); err != nil {
		t.Fatalf("UpdatePercentage failed: %v", err)
	}

	// Перезагружаем
	if err := client.Reload(); err != nil {
		t.Fatalf("Reload failed: %v", err)
	}

	// Теперь должен быть включён (100%)
	if !client.Enabled("reload-flag", "user1") {
		t.Error("flag should be enabled at 100% after reload")
	}
}

func TestGet_ReturnsFlag(t *testing.T) {
	dir := t.TempDir()

	setupTestFlag(t, dir, "get-flag", 50, flag.StatusActive)

	client, err := New(dir)
	if err != nil {
		t.Fatalf("New failed: %v", err)
	}

	f, err := client.Get("get-flag")
	if err != nil {
		t.Fatalf("Get failed: %v", err)
	}

	if f.Meta.Name != "get-flag" {
		t.Errorf("expected name 'get-flag', got %q", f.Meta.Name)
	}
	if f.Meta.Percentage != 50 {
		t.Errorf("expected percentage 50, got %d", f.Meta.Percentage)
	}
}

func TestGet_NotFound(t *testing.T) {
	dir := t.TempDir()

	client, err := New(dir)
	if err != nil {
		t.Fatalf("New failed: %v", err)
	}

	_, err = client.Get("nonexistent")
	if err == nil {
		t.Error("expected error for nonexistent flag")
	}
}
