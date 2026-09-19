package mdflag

import (
	"fmt"
	"sync"
	"testing"
	"time"

	"github.com/legends-deadlines/mdflag/internal/flag"
)

// Тест: параллельные вызовы Enabled
func TestClient_ConcurrentEnabled(t *testing.T) {
	dir := t.TempDir()

	// Создаём флаг
	store, err := flag.NewStore(dir)
	if err != nil {
		t.Fatalf("NewStore failed: %v", err)
	}

	meta := flag.FlagMeta{
		Name:       "concurrent-flag",
		Percentage: 50,
		Targeting:  "user_id",
		Status:     flag.StatusActive,
		Hypothesis: "Test",
	}

	if _, err := store.Create(meta, ""); err != nil {
		t.Fatalf("Create failed: %v", err)
	}

	// Создаём клиент
	client, err := New(dir)
	if err != nil {
		t.Fatalf("New failed: %v", err)
	}

	// Параллельные вызовы
	var wg sync.WaitGroup
	errors := make(chan error, 100)

	for i := 0; i < 100; i++ {
		wg.Add(1)
		go func(id int) {
			defer wg.Done()
			entityID := fmt.Sprintf("user_%d", id)
			_ = client.Enabled("concurrent-flag", entityID)
		}(i)
	}

	wg.Wait()
	close(errors)

	for err := range errors {
		t.Errorf("concurrent error: %v", err)
	}
}

// Тест: горячая перезагрузка
func TestClient_HotReload(t *testing.T) {
	dir := t.TempDir()

	// Создаём флаг с 0%
	store, err := flag.NewStore(dir)
	if err != nil {
		t.Fatalf("NewStore failed: %v", err)
	}

	meta := flag.FlagMeta{
		Name:       "hot-reload-flag",
		Percentage: 0,
		Targeting:  "user_id",
		Status:     flag.StatusActive,
		Hypothesis: "Test",
	}

	if _, err := store.Create(meta, ""); err != nil {
		t.Fatalf("Create failed: %v", err)
	}

	// Создаём клиент
	client, err := New(dir)
	if err != nil {
		t.Fatalf("New failed: %v", err)
	}

	// Изначально выключен
	if client.Enabled("hot-reload-flag", "user1") {
		t.Error("flag should be disabled at 0%")
	}

	// Меняем процент
	if err := store.UpdatePercentage("hot-reload-flag", 100); err != nil {
		t.Fatalf("UpdatePercentage failed: %v", err)
	}

	// Без перезагрузки всё ещё выключен (кэш)
	if client.Enabled("hot-reload-flag", "user1") {
		t.Error("flag should still be disabled before reload")
	}

	// Перезагружаем
	if err := client.Reload(); err != nil {
		t.Fatalf("Reload failed: %v", err)
	}

	// Теперь включён
	if !client.Enabled("hot-reload-flag", "user1") {
		t.Error("flag should be enabled after reload")
	}
}

// Тест: распределение процентов на большой выборке
func TestClient_PercentageDistribution(t *testing.T) {
	dir := t.TempDir()

	store, err := flag.NewStore(dir)
	if err != nil {
		t.Fatalf("NewStore failed: %v", err)
	}

	testCases := []struct {
		percentage int
		tolerance  float64
	}{
		{10, 0.05},
		{25, 0.05},
		{50, 0.05},
		{75, 0.05},
		{90, 0.05},
	}

	for _, tc := range testCases {
		name := fmt.Sprintf("dist-%d", tc.percentage)
		meta := flag.FlagMeta{
			Name:       name,
			Percentage: tc.percentage,
			Targeting:  "user_id",
			Status:     flag.StatusActive,
			Hypothesis: "Distribution test",
		}

		if _, err := store.Create(meta, ""); err != nil {
			t.Fatalf("Create failed for %d%%: %v", tc.percentage, err)
		}
	}

	client, err := New(dir)
	if err != nil {
		t.Fatalf("New failed: %v", err)
	}

	totalEntities := 10000

	for _, tc := range testCases {
		name := fmt.Sprintf("dist-%d", tc.percentage)
		enabledCount := 0

		for i := 0; i < totalEntities; i++ {
			entityID := fmt.Sprintf("user_%d", i)
			if client.Enabled(name, entityID) {
				enabledCount++
			}
		}

		actualRatio := float64(enabledCount) / float64(totalEntities)
		expectedRatio := float64(tc.percentage) / 100.0

		diff := actualRatio - expectedRatio
		if diff < 0 {
			diff = -diff
		}

		if diff > tc.tolerance {
			t.Errorf("percentage %d%%: expected ~%.2f, got %.2f (diff %.2f > tolerance %.2f)",
				tc.percentage, expectedRatio, actualRatio, diff, tc.tolerance)
		}
	}
}

// Тест: срок действия флага
func TestClient_Expiration(t *testing.T) {
	dir := t.TempDir()

	store, err := flag.NewStore(dir)
	if err != nil {
		t.Fatalf("NewStore failed: %v", err)
	}

	// Флаг, который истёк час назад
	pastTime := time.Now().Add(-1 * time.Hour)
	meta := flag.FlagMeta{
		Name:       "expired-flag",
		Percentage: 100,
		Targeting:  "user_id",
		Status:     flag.StatusActive,
		Hypothesis: "Test",
		Expires:    pastTime,
	}

	if _, err := store.Create(meta, ""); err != nil {
		t.Fatalf("Create failed: %v", err)
	}

	client, err := New(dir)
	if err != nil {
		t.Fatalf("New failed: %v", err)
	}

	// Истёкший флаг должен быть выключен, даже если процент 100
	if client.Enabled("expired-flag", "user1") {
		t.Error("expired flag should be disabled")
	}
}
