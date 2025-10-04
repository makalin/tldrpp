package plugin

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/makalin/tldrpp/internal/types"
)

// Plugin represents a tldr++ plugin
type Plugin interface {
	Name() string
	Description() string
	Execute(args []string) error
}

// SubmitPlugin handles submission to tldr-pages
type SubmitPlugin struct {
	page    *types.Page
	example *types.Example
}

// NewSubmitPlugin creates a new submit plugin
func NewSubmitPlugin(page *types.Page, example *types.Example) *SubmitPlugin {
	return &SubmitPlugin{
		page:    page,
		example: example,
	}
}

// Name returns the plugin name
func (p *SubmitPlugin) Name() string {
	return "submit"
}

// Description returns the plugin description
func (p *SubmitPlugin) Description() string {
	return "Submit example to tldr-pages repository"
}

// Execute executes the submit plugin
func (p *SubmitPlugin) Execute(args []string) error {
	if len(args) == 0 {
		return fmt.Errorf("no command specified")
	}

	switch args[0] {
	case "init":
		return p.initSubmission()
	case "validate":
		return p.validateExample()
	case "create-pr":
		return p.createPullRequest()
	default:
		return fmt.Errorf("unknown command: %s", args[0])
	}
}

// initSubmission initializes a new submission
func (p *SubmitPlugin) initSubmission() error {
	fmt.Println("Initializing tldr-pages submission...")
	fmt.Printf("Page: %s (%s)\n", p.page.Name, p.page.Platform)
	fmt.Printf("Example: %s\n", p.example.Description)
	fmt.Printf("Command: %s\n", p.example.Command)
	fmt.Println()

	// Check if git is available
	if !p.isGitAvailable() {
		return fmt.Errorf("git is not available. Please install git to submit to tldr-pages")
	}

	// Check if gh CLI is available
	if !p.isGitHubCLIAvailable() {
		fmt.Println("Warning: GitHub CLI (gh) is not available.")
		fmt.Println("You'll need to manually create a pull request.")
	}

	// Create submission directory
	submissionDir := filepath.Join(os.TempDir(), "tldrpp-submission")
	if err := os.MkdirAll(submissionDir, 0755); err != nil {
		return fmt.Errorf("failed to create submission directory: %w", err)
	}

	// Generate markdown content
	content := p.generateMarkdown()
	contentFile := filepath.Join(submissionDir, fmt.Sprintf("%s.md", p.page.Name))
	if err := os.WriteFile(contentFile, []byte(content), 0644); err != nil {
		return fmt.Errorf("failed to write content file: %w", err)
	}

	fmt.Printf("Submission files created in: %s\n", submissionDir)
	fmt.Println("Next steps:")
	fmt.Println("1. Review the generated markdown file")
	fmt.Println("2. Run 'tldrpp plugin submit validate' to check for issues")
	fmt.Println("3. Run 'tldrpp plugin submit create-pr' to create a pull request")

	return nil
}

// validateExample validates the example against tldr-pages standards
func (p *SubmitPlugin) validateExample() error {
	fmt.Println("Validating example against tldr-pages standards...")

	var issues []string

	// Check description length
	if len(p.example.Description) > 80 {
		issues = append(issues, "Description is too long (>80 characters)")
	}

	// Check command length
	if len(p.example.Command) > 100 {
		issues = append(issues, "Command is too long (>100 characters)")
	}

	// Check for common issues
	if strings.Contains(p.example.Command, "sudo") {
		issues = append(issues, "Avoid using 'sudo' in examples")
	}

	if strings.Contains(p.example.Command, "&&") {
		issues = append(issues, "Avoid chaining commands with '&&'")
	}

	// Check placeholder usage
	for _, placeholder := range p.example.Placeholders {
		if placeholder.Name == "" {
			issues = append(issues, "Empty placeholder name found")
		}
		if len(placeholder.Name) > 20 {
			issues = append(issues, fmt.Sprintf("Placeholder name '%s' is too long", placeholder.Name))
		}
	}

	if len(issues) == 0 {
		fmt.Println("✓ Example validation passed!")
		return nil
	}

	fmt.Println("✗ Validation issues found:")
	for _, issue := range issues {
		fmt.Printf("  - %s\n", issue)
	}

	return fmt.Errorf("validation failed with %d issues", len(issues))
}

// createPullRequest creates a pull request to tldr-pages
func (p *SubmitPlugin) createPullRequest() error {
	fmt.Println("Creating pull request to tldr-pages...")

	// Check if gh CLI is available
	if !p.isGitHubCLIAvailable() {
		return fmt.Errorf("GitHub CLI (gh) is not available. Please install it or create a PR manually")
	}

	// Generate branch name
	branchName := fmt.Sprintf("tldrpp-%s-%s", p.page.Name, p.page.Platform)

	// Create markdown content
	content := p.generateMarkdown()

	// Create a temporary file for the content
	tempFile := filepath.Join(os.TempDir(), fmt.Sprintf("%s.md", p.page.Name))
	if err := os.WriteFile(tempFile, []byte(content), 0644); err != nil {
		return fmt.Errorf("failed to create temporary file: %w", err)
	}
	defer os.Remove(tempFile)

	// Create PR using gh CLI
	title := fmt.Sprintf("Add example for %s (%s)", p.page.Name, p.page.Platform)
	body := fmt.Sprintf("This PR adds a new example for the `%s` command on the `%s` platform.\n\nExample: %s\n\nCommand: `%s`", 
		p.page.Name, p.page.Platform, p.example.Description, p.example.Command)

	cmd := exec.Command("gh", "pr", "create", 
		"--repo", "tldr-pages/tldr",
		"--title", title,
		"--body", body,
		"--file", tempFile)

	if err := cmd.Run(); err != nil {
		return fmt.Errorf("failed to create pull request: %w", err)
	}

	fmt.Println("✓ Pull request created successfully!")
	return nil
}

// generateMarkdown generates markdown content for the submission
func (p *SubmitPlugin) generateMarkdown() string {
	var content strings.Builder

	// Title
	content.WriteString(fmt.Sprintf("# %s\n\n", p.page.Name))

	// Description
	content.WriteString(fmt.Sprintf("> %s.\n\n", p.page.Description))

	// Example
	content.WriteString(fmt.Sprintf("- %s:\n", p.example.Description))
	content.WriteString(fmt.Sprintf("  `%s`\n", p.example.Command))

	return content.String()
}

// isGitAvailable checks if git is available
func (p *SubmitPlugin) isGitAvailable() bool {
	_, err := exec.LookPath("git")
	return err == nil
}

// isGitHubCLIAvailable checks if GitHub CLI is available
func (p *SubmitPlugin) isGitHubCLIAvailable() bool {
	_, err := exec.LookPath("gh")
	return err == nil
}

// PluginManager manages plugins
type PluginManager struct {
	plugins map[string]Plugin
}

// NewPluginManager creates a new plugin manager
func NewPluginManager() *PluginManager {
	return &PluginManager{
		plugins: make(map[string]Plugin),
	}
}

// RegisterPlugin registers a plugin
func (pm *PluginManager) RegisterPlugin(plugin Plugin) {
	pm.plugins[plugin.Name()] = plugin
}

// ExecutePlugin executes a plugin
func (pm *PluginManager) ExecutePlugin(name string, args []string) error {
	plugin, exists := pm.plugins[name]
	if !exists {
		return fmt.Errorf("plugin '%s' not found", name)
	}

	return plugin.Execute(args)
}

// ListPlugins lists all registered plugins
func (pm *PluginManager) ListPlugins() []Plugin {
	var plugins []Plugin
	for _, plugin := range pm.plugins {
		plugins = append(plugins, plugin)
	}
	return plugins
}

// InteractiveMode runs the plugin in interactive mode
func (pm *PluginManager) InteractiveMode() error {
	scanner := bufio.NewScanner(os.Stdin)

	for {
		fmt.Print("tldrpp plugin> ")
		if !scanner.Scan() {
			break
		}

		line := strings.TrimSpace(scanner.Text())
		if line == "" {
			continue
		}

		if line == "exit" || line == "quit" {
			break
		}

		if line == "help" {
			pm.showHelp()
			continue
		}

		if line == "list" {
			pm.listPlugins()
			continue
		}

		// Parse command
		parts := strings.Fields(line)
		if len(parts) < 2 {
			fmt.Println("Usage: <plugin> <command> [args...]")
			continue
		}

		pluginName := parts[0]
		args := parts[1:]

		if err := pm.ExecutePlugin(pluginName, args); err != nil {
			fmt.Printf("Error: %v\n", err)
		}
	}

	return nil
}

// showHelp shows help information
func (pm *PluginManager) showHelp() {
	fmt.Println("tldr++ Plugin System")
	fmt.Println()
	fmt.Println("Available commands:")
	fmt.Println("  help                    Show this help")
	fmt.Println("  list                    List available plugins")
	fmt.Println("  <plugin> <command>     Execute plugin command")
	fmt.Println("  exit/quit              Exit plugin mode")
	fmt.Println()
	fmt.Println("Available plugins:")
	for _, plugin := range pm.plugins {
		fmt.Printf("  %-10s %s\n", plugin.Name(), plugin.Description())
	}
}

// listPlugins lists all plugins
func (pm *PluginManager) listPlugins() {
	fmt.Println("Available plugins:")
	for _, plugin := range pm.plugins {
		fmt.Printf("  %-10s %s\n", plugin.Name(), plugin.Description())
	}
}