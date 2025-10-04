package app

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/makalin/tldrpp/internal/cache"
	"github.com/makalin/tldrpp/internal/config"
	"github.com/makalin/tldrpp/internal/tui"
	"github.com/spf13/viper"
)

// Initialize downloads the tldr pages index and sets up the cache
func Initialize() error {
	cfg, err := config.Load()
	if err != nil {
		return fmt.Errorf("failed to load config: %w", err)
	}

	cacheManager := cache.New(cfg.CacheDir)
	return cacheManager.Initialize()
}

// UpdateCache refreshes the tldr pages cache
func UpdateCache() error {
	cfg, err := config.Load()
	if err != nil {
		return fmt.Errorf("failed to load config: %w", err)
	}

	cacheManager := cache.New(cfg.CacheDir)
	return cacheManager.Update()
}

// RunTUI starts the terminal user interface
func RunTUI(searchQuery, platform, theme string, dev bool) error {
	cfg, err := config.Load()
	if err != nil {
		return fmt.Errorf("failed to load config: %w", err)
	}

	// Override config with command line flags
	if platform != "" {
		cfg.Platforms = []string{platform}
	}
	if theme != "" {
		cfg.Theme = theme
	}
	if dev {
		cfg.DevMode = true
	}

	cacheManager := cache.New(cfg.CacheDir)
	
	// Ensure cache is initialized
	if !cacheManager.IsInitialized() {
		if err := cacheManager.Initialize(); err != nil {
			return fmt.Errorf("failed to initialize cache: %w", err)
		}
	}

	app := tui.New(cfg, cacheManager)
	return app.Run(searchQuery)
}

// RenderCommand renders a command with placeholders filled
func RenderCommand(command string, vars map[string]string) error {
	cfg, err := config.Load()
	if err != nil {
		return fmt.Errorf("failed to load config: %w", err)
	}

	cacheManager := cache.New(cfg.CacheDir)
	if !cacheManager.IsInitialized() {
		if err := cacheManager.Initialize(); err != nil {
			return fmt.Errorf("failed to initialize cache: %w", err)
		}
	}

	page, err := cacheManager.FindPage(command)
	if err != nil {
		return fmt.Errorf("command not found: %w", err)
	}

	// Find the best matching example
	example := page.FindBestExample(command)
	if example == nil {
		return fmt.Errorf("no suitable example found for command: %s", command)
	}

	// Render the command with variables
	rendered := example.Render(vars)
	fmt.Println(rendered)
	return nil
}

// ExecuteCommand executes a command with placeholders filled
func ExecuteCommand(command string, vars map[string]string) error {
	cfg, err := config.Load()
	if err != nil {
		return fmt.Errorf("failed to load config: %w", err)
	}

	cacheManager := cache.New(cfg.CacheDir)
	if !cacheManager.IsInitialized() {
		if err := cacheManager.Initialize(); err != nil {
			return fmt.Errorf("failed to initialize cache: %w", err)
		}
	}

	page, err := cacheManager.FindPage(command)
	if err != nil {
		return fmt.Errorf("command not found: %w", err)
	}

	// Find the best matching example
	example := page.FindBestExample(command)
	if example == nil {
		return fmt.Errorf("no suitable example found for command: %s", command)
	}

	// Render the command with variables
	rendered := example.Render(vars)
	
	// Check if command is destructive
	if isDestructiveCommand(rendered) && cfg.ConfirmDestructive {
		fmt.Printf("This command appears destructive: %s\n", rendered)
		fmt.Print("Are you sure you want to execute it? (y/N): ")
		var response string
		fmt.Scanln(&response)
		if strings.ToLower(response) != "y" && strings.ToLower(response) != "yes" {
			fmt.Println("Command cancelled.")
			return nil
		}
	}

	// Execute the command
	cmd := exec.Command("sh", "-c", rendered)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Stdin = os.Stdin

	// Log the execution
	if err := logExecution(rendered); err != nil {
		fmt.Fprintf(os.Stderr, "Warning: failed to log execution: %v\n", err)
	}

	return cmd.Run()
}

// SubmitToTldr opens the plugin for submitting examples to tldr-pages
func SubmitToTldr() error {
	cfg, err := config.Load()
	if err != nil {
		return fmt.Errorf("failed to load config: %w", err)
	}

	cacheManager := cache.New(cfg.CacheDir)
	if !cacheManager.IsInitialized() {
		if err := cacheManager.Initialize(); err != nil {
			return fmt.Errorf("failed to initialize cache: %w", err)
		}
	}

	// For now, we'll just print a message
	// In a real implementation, this would open the plugin interface
	fmt.Println("Plugin system initialized. Use 'tldrpp plugin submit init' to start a submission.")
	return nil
}

// isDestructiveCommand checks if a command is potentially destructive
func isDestructiveCommand(command string) bool {
	destructiveVerbs := []string{
		"rm", "rmdir", "del", "erase",
		"dd", "mkfs", "fdisk", "parted",
		"iptables", "ufw", "firewall-cmd",
		"chmod", "chown", "chattr",
		"kill", "killall", "pkill",
		"shutdown", "reboot", "halt",
		"mv", "move", "rename",
		"cp", "copy", "xcopy",
		"tar", "zip", "unzip",
		"git", "svn", "hg",
	}

	command = strings.ToLower(command)
	for _, verb := range destructiveVerbs {
		if strings.HasPrefix(command, verb+" ") || command == verb {
			return true
		}
	}
	return false
}

// logExecution logs command execution to audit log
func logExecution(command string) error {
	cfg, err := config.Load()
	if err != nil {
		return err
	}

	logDir := filepath.Join(cfg.CacheDir, "..")
	if err := os.MkdirAll(logDir, 0755); err != nil {
		return err
	}

	logFile := filepath.Join(logDir, "exec.log")
	f, err := os.OpenFile(logFile, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return err
	}
	defer f.Close()

	_, err = fmt.Fprintf(f, "%s: %s\n", viper.GetString("timestamp"), command)
	return err
}