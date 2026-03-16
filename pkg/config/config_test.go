package config

import (
	"os"
	"path/filepath"
	"testing"
)

func TestConfigManager(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "orez-test")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(tmpDir)

	cfgPath := filepath.Join(tmpDir, "config.json")
	m, err := NewConfigManager(cfgPath)
	if err != nil {
		t.Fatal(err)
	}

	// Test Load (should create defaults)
	if err := m.Load(); err != nil {
		t.Fatalf("Failed to load config: %v", err)
	}

	cfg := m.Get()
	if cfg.Provider.Model != "llama3" {
		t.Errorf("Expected default model llama3, got %s", cfg.Provider.Model)
	}

	// Test Save
	cfg.Provider.Model = "gpt-4"
	// t.Logf("Saving cfg with model: %s", cfg.Provider.Model)
	if err := m.Save(cfg); err != nil {
		t.Fatalf("Failed to save config: %v", err)
	}

	// Reload to verify
	m2, _ := NewConfigManager(cfgPath)
	if err := m2.Load(); err != nil {
		data, _ := os.ReadFile(cfgPath)
		t.Fatalf("Failed to load m2: %v. Content: %s", err, string(data))
	}
	// t.Logf("Reloaded model: %s", m2.Get().Provider.Model)
	if m2.Get().Provider.Model != "gpt-4" {
		t.Errorf("Expected saved model gpt-4, got %s", m2.Get().Provider.Model)
	}
}

func TestStripComments(t *testing.T) {
	input := []byte(`{
		// line comment
		"key": "value" /* block
		comment */
	}`)
	
	output := stripComments(input)
	// We don't need exact match of whitespace, but it should be valid JSON
	if len(output) > len(input) {
		t.Error("stripComments increased size")
	}
}
