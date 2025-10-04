package config

import (
	"os"
	"path/filepath"
	"testing"
)

func TestDefaultConfig(t *testing.T) {
	cfg := DefaultConfig()

	if cfg.Theme != "dark" {
		t.Errorf("Expected theme 'dark', got '%s'", cfg.Theme)
	}

	if len(cfg.Platforms) != 2 || cfg.Platforms[0] != "common" || cfg.Platforms[1] != "linux" {
		t.Errorf("Expected platforms ['common', 'linux'], got %v", cfg.Platforms)
	}

	if !cfg.ConfirmDestructive {
		t.Error("Expected ConfirmDestructive to be true")
	}

	if !cfg.Clipboard {
		t.Error("Expected Clipboard to be true")
	}

	if cfg.Pager != "less -R" {
		t.Errorf("Expected pager 'less -R', got '%s'", cfg.Pager)
	}

	if cfg.Keymap.Run != "ctrl+enter" {
		t.Errorf("Expected keymap run 'ctrl+enter', got '%s'", cfg.Keymap.Run)
	}

	if cfg.Keymap.Copy != "y" {
		t.Errorf("Expected keymap copy 'y', got '%s'", cfg.Keymap.Copy)
	}

	if cfg.Keymap.Paste != "p" {
		t.Errorf("Expected keymap paste 'p', got '%s'", cfg.Keymap.Paste)
	}

	if cfg.CacheTTLHours != 72 {
		t.Errorf("Expected cache TTL 72 hours, got %d", cfg.CacheTTLHours)
	}

	if cfg.DevMode {
		t.Error("Expected DevMode to be false")
	}
}

func TestLoadConfig(t *testing.T) {
	// Test loading config when file doesn't exist
	cfg, err := Load()
	if err != nil {
		t.Fatalf("Load failed: %v", err)
	}

	// Should return default config
	if cfg.Theme != "dark" {
		t.Errorf("Expected theme 'dark', got '%s'", cfg.Theme)
	}

	if len(cfg.Platforms) != 2 {
		t.Errorf("Expected 2 platforms, got %d", len(cfg.Platforms))
	}
}

func TestSaveAndLoadConfig(t *testing.T) {
	// Create temporary directory
	tempDir := t.TempDir()
	configDir := filepath.Join(tempDir, ".config", "tldrpp")
	configFile := filepath.Join(configDir, "config.yml")

	// Override config directory for testing
	originalGetConfigDir := getConfigDir
	getConfigDir = func() string {
		return configDir
	}
	defer func() {
		getConfigDir = originalGetConfigDir
	}()

	// Create and save config
	cfg := DefaultConfig()
	cfg.Theme = "light"
	cfg.Platforms = []string{"linux", "osx"}

	err := cfg.Save()
	if err != nil {
		t.Fatalf("Save failed: %v", err)
	}

	// Check if file was created
	if _, err := os.Stat(configFile); os.IsNotExist(err) {
		t.Fatal("Config file was not created")
	}

	// Load config
	loadedCfg, err := Load()
	if err != nil {
		t.Fatalf("Load failed: %v", err)
	}

	// Check loaded values
	if loadedCfg.Theme != "light" {
		t.Errorf("Expected theme 'light', got '%s'", loadedCfg.Theme)
	}

	if len(loadedCfg.Platforms) != 2 || loadedCfg.Platforms[0] != "linux" || loadedCfg.Platforms[1] != "osx" {
		t.Errorf("Expected platforms ['linux', 'osx'], got %v", loadedCfg.Platforms)
	}
}

func TestGetConfigDir(t *testing.T) {
	dir := getConfigDir()
	if dir == "" {
		t.Error("Expected non-empty config directory")
	}

	// Should contain .config/tldrpp
	if !filepath.IsAbs(dir) {
		t.Error("Expected absolute path for config directory")
	}
}

func TestGetDefaultCacheDir(t *testing.T) {
	dir := getDefaultCacheDir()
	if dir == "" {
		t.Error("Expected non-empty cache directory")
	}

	// Should contain .cache/tldrpp/pages
	if !filepath.IsAbs(dir) {
		t.Error("Expected absolute path for cache directory")
	}
}

func TestCreateDefaultConfig(t *testing.T) {
	// Create temporary directory
	tempDir := t.TempDir()
	configFile := filepath.Join(tempDir, "config.yml")

	err := createDefaultConfig(configFile)
	if err != nil {
		t.Fatalf("createDefaultConfig failed: %v", err)
	}

	// Check if file was created
	if _, err := os.Stat(configFile); os.IsNotExist(err) {
		t.Fatal("Default config file was not created")
	}
}