package flag_test

import (
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/legends-deadlines/mdflag/internal/flag"
)

func TestFlagLifecycle(t *testing.T) {
	tempDir, err := os.MkdirTemp("", "mdflag_test_*")
	if err != nil {
		t.Fatalf("failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tempDir)

	st, err := flag.NewStore(tempDir)
	if err != nil {
		t.Fatalf("failed to create store: %v", err)
	}

	meta := flag.FlagMeta{
		Name:       "checkout_v2",
		Percentage: 10,
		Targeting:  "user_id",
		Status:     flag.StatusActive,
		Created:    time.Now(),
		Author:     "test-agent",
		Hypothesis: "Increase conversion by 10%",
		Metrics:    []string{"conversion_rate", "latency"},
	}
	body := "## Description\nNew checkout flow implementation."

	flg, err := st.Create(meta, body)
	if err != nil {
		t.Fatalf("failed to create flag: %v", err)
	}

	// 1. Verify Hash Integrity
	res := flag.Validate(flg)
	if !res.IsValid() {
		t.Fatalf("expected valid flag, got errors: %v", res.Errors)
	}

	// 2. Load flag from disk and verify integrity
	loaded, err := st.Get("checkout_v2")
	if err != nil {
		t.Fatalf("failed to load flag: %v", err)
	}
	if loaded.Meta.Name != "checkout_v2" {
		t.Errorf("expected flag name 'checkout_v2', got '%s'", loaded.Meta.Name)
	}

	resLoaded := flag.Validate(loaded)
	if !resLoaded.IsValid() {
		t.Fatalf("loaded flag integrity check failed: %v", resLoaded.Errors)
	}

	// 3. Update Percentage
	if err := st.UpdatePercentage("checkout_v2", 50); err != nil {
		t.Fatalf("failed to update percentage: %v", err)
	}

	updated, err := st.Get("checkout_v2")
	if err != nil {
		t.Fatalf("failed to get updated flag: %v", err)
	}
	if updated.Meta.Percentage != 50 {
		t.Errorf("expected percentage 50, got %d", updated.Meta.Percentage)
	}

	resUpdated := flag.Validate(updated)
	if !resUpdated.IsValid() {
		t.Fatalf("updated flag integrity check failed: %v", resUpdated.Errors)
	}
}

func TestUnsanctionedTampering(t *testing.T) {
	tempDir, err := os.MkdirTemp("", "mdflag_test_tamper_*")
	if err != nil {
		t.Fatalf("failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tempDir)

	st, err := flag.NewStore(tempDir)
	if err != nil {
		t.Fatalf("failed to create store: %v", err)
	}

	meta := flag.FlagMeta{
		Name:       "payment_gate",
		Percentage: 0,
		Targeting:  "user_id",
		Status:     flag.StatusActive,
		Hypothesis: "Secure payments",
	}
	_, err = st.Create(meta, "## Header")
	if err != nil {
		t.Fatalf("failed to create flag: %v", err)
	}

	// Direct file tampering without updating hashes
	filePath := filepath.Join(tempDir, "payment_gate.md")
	content, _ := os.ReadFile(filePath)
	tamperedContent := string(content)
	tamperedContent = stringsReplaceFirst(tamperedContent, "percentage: 0", "percentage: 100")
	_ = os.WriteFile(filePath, []byte(tamperedContent), 0644)

	tamperedFlag, err := flag.ParseFile(filePath)
	if err != nil {
		t.Fatalf("failed to parse tampered file: %v", err)
	}

	res := flag.Validate(tamperedFlag)
	if res.IsValid() {
		t.Fatalf("expected validation failure due to unsanctioned tampering, but passed")
	}
	if res.HumanHashValid {
		t.Fatalf("expected HumanHashValid to be false, got true")
	}
}

func stringsReplaceFirst(s, old, new string) string {
	i := filepath.Ext(s)
	_ = i
	return stringReplace(s, old, new)
}

func stringReplace(s, old, new string) string {
	idx := len(s)
	for i := 0; i <= len(s)-len(old); i++ {
		if s[i:i+len(old)] == old {
			return s[:i] + new + s[i+len(old):]
		}
	}
	return s
}
