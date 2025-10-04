package config

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/spf13/viper"
)

// Config represents the application configuration
type Config struct {
	Theme              string   `yaml:"theme"`
	Platforms          []string `yaml:"platforms"`
	ConfirmDestructive bool     `yaml:"confirm_destructive"`
	Clipboard          bool     `yaml:"clipboard"`
	Pager              string   `yaml:"pager"`
	Keymap             Keymap   `yaml:"keymap"`
	CacheTTLHours      int      `yaml:"cache_ttl_hours"`
	CacheDir           string   `yaml:"cache_dir"`
	DevMode            bool     `yaml:"dev_mode"`
}

// Keymap represents keyboard shortcuts configuration
type Keymap struct {
	Run   string `yaml:"run"`
	Copy  string `yaml:"copy"`
	Paste string `yaml:"paste"`
}

// DefaultConfig returns the default configuration
func DefaultConfig() *Config {
	return &Config{
		Theme:              "dark",
		Platforms:          []string{"common", "linux"},
		ConfirmDestructive: true,
		Clipboard:          true,
		Pager:              "less -R",
		Keymap: Keymap{
			Run:   "ctrl+enter",
			Copy:  "y",
			Paste: "p",
		},
		CacheTTLHours: 72,
		CacheDir:      getDefaultCacheDir(),
		DevMode:       false,
	}
}

// Load loads the configuration from file or returns default
func Load() (*Config, error) {
	configDir := getConfigDir()
	configFile := filepath.Join(configDir, "config.yml")

	// Set up viper
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath(configDir)

	// Set defaults
	cfg := DefaultConfig()
	viper.SetDefault("theme", cfg.Theme)
	viper.SetDefault("platforms", cfg.Platforms)
	viper.SetDefault("confirm_destructive", cfg.ConfirmDestructive)
	viper.SetDefault("clipboard", cfg.Clipboard)
	viper.SetDefault("pager", cfg.Pager)
	viper.SetDefault("keymap.run", cfg.Keymap.Run)
	viper.SetDefault("keymap.copy", cfg.Keymap.Copy)
	viper.SetDefault("keymap.paste", cfg.Keymap.Paste)
	viper.SetDefault("cache_ttl_hours", cfg.CacheTTLHours)
	viper.SetDefault("cache_dir", cfg.CacheDir)

	// Try to read config file
	if err := viper.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			// Config file not found, create default
			if err := createDefaultConfig(configFile); err != nil {
				return cfg, fmt.Errorf("failed to create default config: %w", err)
			}
		} else {
			return cfg, fmt.Errorf("failed to read config: %w", err)
		}
	}

	// Unmarshal into struct
	if err := viper.Unmarshal(cfg); err != nil {
		return cfg, fmt.Errorf("failed to unmarshal config: %w", err)
	}

	// Ensure cache directory exists
	if err := os.MkdirAll(cfg.CacheDir, 0755); err != nil {
		return cfg, fmt.Errorf("failed to create cache directory: %w", err)
	}

	return cfg, nil
}

// Save saves the configuration to file
func (c *Config) Save() error {
	configDir := getConfigDir()
	configFile := filepath.Join(configDir, "config.yml")

	// Ensure config directory exists
	if err := os.MkdirAll(configDir, 0755); err != nil {
		return fmt.Errorf("failed to create config directory: %w", err)
	}

	// Set viper values
	viper.Set("theme", c.Theme)
	viper.Set("platforms", c.Platforms)
	viper.Set("confirm_destructive", c.ConfirmDestructive)
	viper.Set("clipboard", c.Clipboard)
	viper.Set("pager", c.Pager)
	viper.Set("keymap.run", c.Keymap.Run)
	viper.Set("keymap.copy", c.Keymap.Copy)
	viper.Set("keymap.paste", c.Keymap.Paste)
	viper.Set("cache_ttl_hours", c.CacheTTLHours)
	viper.Set("cache_dir", c.CacheDir)

	return viper.WriteConfigAs(configFile)
}

// getConfigDir returns the configuration directory
func getConfigDir() string {
	if homeDir, err := os.UserHomeDir(); err == nil {
		return filepath.Join(homeDir, ".config", "tldrpp")
	}
	return filepath.Join(".", ".config", "tldrpp")
}

// getDefaultCacheDir returns the default cache directory
func getDefaultCacheDir() string {
	if homeDir, err := os.UserHomeDir(); err == nil {
		return filepath.Join(homeDir, ".cache", "tldrpp", "pages")
	}
	return filepath.Join(".", ".cache", "tldrpp", "pages")
}

// createDefaultConfig creates a default configuration file
func createDefaultConfig(configFile string) error {
	configDir := filepath.Dir(configFile)
	if err := os.MkdirAll(configDir, 0755); err != nil {
		return err
	}

	cfg := DefaultConfig()
	return cfg.Save()
}