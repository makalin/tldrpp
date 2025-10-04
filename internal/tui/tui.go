package tui

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/makalin/tldrpp/internal/cache"
	"github.com/makalin/tldrpp/internal/config"
	"github.com/makalin/tldrpp/internal/types"
)

// App represents the main TUI application
type App struct {
	config      *config.Config
	cache       *cache.Manager
	state       AppState
	searchQuery string
	pages       []*types.Page
	selectedIdx int
	platforms   []string
	theme        Theme
}

// AppState represents the current state of the application
type AppState int

const (
	StateSearch AppState = iota
	StatePages
	StateExamples
	StateEdit
	StateHelp
)

// Theme represents the UI theme
type Theme struct {
	Background   lipgloss.Color
	Foreground   lipgloss.Color
	Accent       lipgloss.Color
	Success      lipgloss.Color
	Warning      lipgloss.Color
	Error        lipgloss.Color
	Border       lipgloss.Color
	Highlight    lipgloss.Color
}

// New creates a new TUI application
func New(cfg *config.Config, cacheManager *cache.Manager) *App {
	app := &App{
		config:    cfg,
		cache:     cacheManager,
		state:     StateSearch,
		platforms: cfg.Platforms,
		theme:     getTheme(cfg.Theme),
	}
	
	return app
}

// Run starts the TUI application
func (a *App) Run(searchQuery string) error {
	a.searchQuery = searchQuery
	
	// Load initial pages
	if err := a.loadPages(); err != nil {
		return fmt.Errorf("failed to load pages: %w", err)
	}

	// Create and run the bubbletea program
	p := bubbletea.NewProgram(a, bubbletea.WithAltScreen())
	_, err := p.Run()
	return err
}

// Init initializes the bubbletea model
func (a *App) Init() bubbletea.Cmd {
	return nil
}

// Update handles bubbletea updates
func (a *App) Update(msg bubbletea.Msg) (bubbletea.Model, bubbletea.Cmd) {
	switch msg := msg.(type) {
	case bubbletea.KeyMsg:
		return a.handleKeyPress(msg)
	case bubbletea.WindowSizeMsg:
		return a.handleResize(msg)
	}
	return a, nil
}

// View renders the TUI
func (a *App) View() string {
	switch a.state {
	case StateSearch:
		return a.renderSearch()
	case StatePages:
		return a.renderPages()
	case StateExamples:
		return a.renderExamples()
	case StateEdit:
		return a.renderEdit()
	case StateHelp:
		return a.renderHelp()
	default:
		return a.renderSearch()
	}
}

// handleKeyPress handles keyboard input
func (a *App) handleKeyPress(msg bubbletea.KeyMsg) (bubbletea.Model, bubbletea.Cmd) {
	switch msg.String() {
	case "ctrl+c", "q":
		return a, bubbletea.Quit
	case "?":
		if a.state == StateHelp {
			a.state = StateSearch
		} else {
			a.state = StateHelp
		}
	case "enter":
		if a.state == StateSearch {
			a.state = StatePages
		} else if a.state == StatePages {
			a.state = StateExamples
		}
	case "esc":
		switch a.state {
		case StatePages:
			a.state = StateSearch
		case StateExamples:
			a.state = StatePages
		case StateEdit:
			a.state = StateExamples
		case StateHelp:
			a.state = StateSearch
		}
	case "tab":
		if a.state == StateExamples {
			a.state = StateEdit
		}
	case "ctrl+enter":
		if a.state == StateExamples || a.state == StateEdit {
			return a.executeCommand()
		}
	case "y":
		if a.state == StateExamples || a.state == StateEdit {
			return a.copyCommand()
		}
	case "p":
		if a.state == StateExamples || a.state == StateEdit {
			return a.pasteCommand()
		}
	case "r":
		if a.state == StateSearch {
			return a.refreshCache()
		}
	case "o":
		if a.state == StateExamples {
			return a.openInPager()
		}
	case "a":
		if a.state == StatePages {
			a.toggleAllPlatforms()
		}
	case "1", "2", "3", "4", "5", "6":
		if a.state == StatePages {
			a.togglePlatform(msg.String())
		}
	case "up", "k":
		if a.selectedIdx > 0 {
			a.selectedIdx--
		}
	case "down", "j":
		if a.selectedIdx < len(a.pages)-1 {
			a.selectedIdx++
		}
	}

	return a, nil
}

// handleResize handles window resize events
func (a *App) handleResize(msg bubbletea.WindowSizeMsg) (bubbletea.Model, bubbletea.Cmd) {
	return a, nil
}

// loadPages loads pages based on current search query and platforms
func (a *App) loadPages() error {
	pages, err := a.cache.SearchPages(a.searchQuery, a.platforms)
	if err != nil {
		return err
	}
	a.pages = pages
	a.selectedIdx = 0
	return nil
}

// renderSearch renders the search interface
func (a *App) renderSearch() string {
	var content strings.Builder
	
	// Title
	title := lipgloss.NewStyle().
		Foreground(a.theme.Accent).
		Bold(true).
		Render("tldr++ - Interactive Cheat-Sheets")
	
	content.WriteString(title + "\n\n")
	
	// Search box
	searchBox := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(a.theme.Border).
		Padding(1, 2).
		Render(fmt.Sprintf("Search: %s", a.searchQuery))
	
	content.WriteString(searchBox + "\n\n")
	
	// Instructions
	instructions := lipgloss.NewStyle().
		Foreground(a.theme.Foreground).
		Render("Press Enter to search, ? for help, q to quit")
	
	content.WriteString(instructions)
	
	return content.String()
}

// renderPages renders the pages list
func (a *App) renderPages() string {
	var content strings.Builder
	
	// Header
	header := lipgloss.NewStyle().
		Foreground(a.theme.Accent).
		Bold(true).
		Render(fmt.Sprintf("Pages (%d found)", len(a.pages)))
	
	content.WriteString(header + "\n\n")
	
	// Platform filters
	platforms := lipgloss.NewStyle().
		Foreground(a.theme.Foreground).
		Render(fmt.Sprintf("Platforms: %s", strings.Join(a.platforms, ", ")))
	
	content.WriteString(platforms + "\n\n")
	
	// Pages list
	for i, page := range a.pages {
		style := lipgloss.NewStyle().Foreground(a.theme.Foreground)
		if i == a.selectedIdx {
			style = style.Background(a.theme.Highlight).Foreground(a.theme.Background)
		}
		
		pageText := fmt.Sprintf("%s - %s (%s)", page.Name, page.Description, page.Platform)
		content.WriteString(style.Render(pageText) + "\n")
	}
	
	// Footer
	footer := lipgloss.NewStyle().
		Foreground(a.theme.Foreground).
		Render("↑↓ Navigate, Enter Select, Esc Back, ? Help")
	
	content.WriteString("\n" + footer)
	
	return content.String()
}

// renderExamples renders the examples for the selected page
func (a *App) renderExamples() string {
	if len(a.pages) == 0 || a.selectedIdx >= len(a.pages) {
		return "No pages available"
	}
	
	page := a.pages[a.selectedIdx]
	var content strings.Builder
	
	// Header
	header := lipgloss.NewStyle().
		Foreground(a.theme.Accent).
		Bold(true).
		Render(fmt.Sprintf("%s - %s", page.Name, page.Description))
	
	content.WriteString(header + "\n\n")
	
	// Examples
	for i, example := range page.Examples {
		style := lipgloss.NewStyle().Foreground(a.theme.Foreground)
		if i == 0 { // Highlight first example
			style = style.Background(a.theme.Highlight).Foreground(a.theme.Background)
		}
		
		exampleText := fmt.Sprintf("%s\n  %s", example.Description, example.Command)
		content.WriteString(style.Render(exampleText) + "\n\n")
	}
	
	// Footer
	footer := lipgloss.NewStyle().
		Foreground(a.theme.Foreground).
		Render("Tab Edit, Ctrl+Enter Run, y Copy, p Paste, Esc Back")
	
	content.WriteString(footer)
	
	return content.String()
}

// renderEdit renders the placeholder editing interface
func (a *App) renderEdit() string {
	if len(a.pages) == 0 || a.selectedIdx >= len(a.pages) {
		return "No pages available"
	}
	
	page := a.pages[a.selectedIdx]
	if len(page.Examples) == 0 {
		return "No examples available"
	}
	
	example := page.Examples[0] // Use first example for now
	var content strings.Builder
	
	// Header
	header := lipgloss.NewStyle().
		Foreground(a.theme.Accent).
		Bold(true).
		Render(fmt.Sprintf("Edit: %s", example.Description))
	
	content.WriteString(header + "\n\n")
	
	// Command with placeholders
	command := example.Command
	for _, placeholder := range example.Placeholders {
		placeholderText := fmt.Sprintf("{{%s}}", placeholder.Name)
		highlighted := lipgloss.NewStyle().
			Background(a.theme.Warning).
			Foreground(a.theme.Background).
			Render(placeholderText)
		command = strings.Replace(command, placeholderText, highlighted, 1)
	}
	
	commandBox := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(a.theme.Border).
		Padding(1, 2).
		Render(command)
	
	content.WriteString(commandBox + "\n\n")
	
	// Placeholders
	if len(example.Placeholders) > 0 {
		placeholders := lipgloss.NewStyle().
			Foreground(a.theme.Foreground).
			Render("Placeholders:")
		content.WriteString(placeholders + "\n")
		
		for _, placeholder := range example.Placeholders {
			placeholderText := fmt.Sprintf("  %s (%s): %s", 
				placeholder.Name, placeholder.Type, placeholder.Default)
			content.WriteString(placeholderText + "\n")
		}
	}
	
	// Footer
	footer := lipgloss.NewStyle().
		Foreground(a.theme.Foreground).
		Render("Ctrl+Enter Run, y Copy, p Paste, Esc Back")
	
	content.WriteString("\n" + footer)
	
	return content.String()
}

// renderHelp renders the help screen
func (a *App) renderHelp() string {
	var content strings.Builder
	
	// Title
	title := lipgloss.NewStyle().
		Foreground(a.theme.Accent).
		Bold(true).
		Render("tldr++ Help")
	
	content.WriteString(title + "\n\n")
	
	// Keybindings
	keybindings := []struct {
		key, description string
	}{
		{"Enter", "Accept example / Select page"},
		{"Tab", "Edit placeholders"},
		{"Ctrl+Enter", "Run command (safe)"},
		{"y", "Copy to clipboard"},
		{"p", "Paste to terminal"},
		{"1-6", "Toggle platform filters"},
		{"a", "Toggle all platforms"},
		{"r", "Refresh cache"},
		{"o", "Open in pager"},
		{"?", "Show/hide help"},
		{"Esc", "Go back"},
		{"q", "Quit"},
	}
	
	for _, kb := range keybindings {
		key := lipgloss.NewStyle().
			Foreground(a.theme.Accent).
			Bold(true).
			Render(kb.key)
		desc := lipgloss.NewStyle().
			Foreground(a.theme.Foreground).
			Render(kb.description)
		content.WriteString(fmt.Sprintf("%-15s %s\n", key, desc))
	}
	
	// Footer
	footer := lipgloss.NewStyle().
		Foreground(a.theme.Foreground).
		Render("Press ? to close help")
	
	content.WriteString("\n" + footer)
	
	return content.String()
}

// executeCommand executes the current command
func (a *App) executeCommand() (bubbletea.Model, bubbletea.Cmd) {
	// This would execute the command
	// For now, just show a message
	return a, bubbletea.Quit
}

// copyCommand copies the current command to clipboard
func (a *App) copyCommand() (bubbletea.Model, bubbletea.Cmd) {
	// This would copy to clipboard
	// For now, just show a message
	return a, bubbletea.Quit
}

// pasteCommand pastes the current command to terminal
func (a *App) pasteCommand() (bubbletea.Model, bubbletea.Cmd) {
	// This would paste to terminal
	// For now, just show a message
	return a, bubbletea.Quit
}

// refreshCache refreshes the pages cache
func (a *App) refreshCache() (bubbletea.Model, bubbletea.Cmd) {
	// This would refresh the cache
	// For now, just reload pages
	a.loadPages()
	return a, nil
}

// openInPager opens the current page in a pager
func (a *App) openInPager() (bubbletea.Model, bubbletea.Cmd) {
	// This would open in pager
	// For now, just show a message
	return a, bubbletea.Quit
}

// toggleAllPlatforms toggles all platform filters
func (a *App) toggleAllPlatforms() {
	allPlatforms := []string{"common", "linux", "osx", "sunos", "windows", "android"}
	if len(a.platforms) == len(allPlatforms) {
		a.platforms = []string{"common"}
	} else {
		a.platforms = allPlatforms
	}
	a.loadPages()
}

// togglePlatform toggles a specific platform filter
func (a *App) togglePlatform(platformNum string) {
	platformMap := map[string]string{
		"1": "common",
		"2": "linux", 
		"3": "osx",
		"4": "sunos",
		"5": "windows",
		"6": "android",
	}
	
	platform := platformMap[platformNum]
	if platform == "" {
		return
	}
	
	// Toggle platform
	var newPlatforms []string
	found := false
	for _, p := range a.platforms {
		if p == platform {
			found = true
		} else {
			newPlatforms = append(newPlatforms, p)
		}
	}
	
	if !found {
		newPlatforms = append(newPlatforms, platform)
	}
	
	a.platforms = newPlatforms
	a.loadPages()
}

// getTheme returns the theme configuration
func getTheme(themeName string) Theme {
	switch themeName {
	case "light":
		return Theme{
			Background: lipgloss.Color("#ffffff"),
			Foreground: lipgloss.Color("#000000"),
			Accent:     lipgloss.Color("#0066cc"),
			Success:    lipgloss.Color("#00aa00"),
			Warning:    lipgloss.Color("#ffaa00"),
			Error:      lipgloss.Color("#cc0000"),
			Border:     lipgloss.Color("#cccccc"),
			Highlight:  lipgloss.Color("#e6f3ff"),
		}
	case "solarized":
		return Theme{
			Background: lipgloss.Color("#002b36"),
			Foreground: lipgloss.Color("#839496"),
			Accent:     lipgloss.Color("#268bd2"),
			Success:    lipgloss.Color("#859900"),
			Warning:    lipgloss.Color("#b58900"),
			Error:      lipgloss.Color("#dc322f"),
			Border:     lipgloss.Color("#586e75"),
			Highlight:  lipgloss.Color("#073642"),
		}
	default: // dark
		return Theme{
			Background: lipgloss.Color("#1e1e1e"),
			Foreground: lipgloss.Color("#ffffff"),
			Accent:     lipgloss.Color("#007acc"),
			Success:    lipgloss.Color("#00aa00"),
			Warning:    lipgloss.Color("#ffaa00"),
			Error:      lipgloss.Color("#cc0000"),
			Border:     lipgloss.Color("#333333"),
			Highlight:  lipgloss.Color("#2d2d30"),
		}
	}
}