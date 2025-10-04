package types

import (
	"regexp"
	"strings"
)

// IndexEntry represents an entry in the tldr pages index
type IndexEntry struct {
	Name        string `json:"name"`
	Description string `json:"description"`
	Platform    string `json:"platform"`
}

// Page represents a tldr page
type Page struct {
	Name        string    `json:"name"`
	Description string    `json:"description"`
	Platform    string    `json:"platform"`
	Examples    []Example `json:"examples"`
	RawContent  string    `json:"raw_content"`
}

// Example represents a command example
type Example struct {
	Description string `json:"description"`
	Command     string `json:"command"`
	Placeholders []Placeholder `json:"placeholders"`
}

// Placeholder represents a placeholder in a command
type Placeholder struct {
	Name        string `json:"name"`
	Type        string `json:"type"`
	Description string `json:"description"`
	Default     string `json:"default"`
}

// ParsePage parses a tldr page from markdown content
func ParsePage(content string, entry IndexEntry) (*Page, error) {
	page := &Page{
		Name:        entry.Name,
		Description: entry.Description,
		Platform:    entry.Platform,
		RawContent:  content,
	}

	lines := strings.Split(content, "\n")
	var currentExample *Example
	var inExample bool

	for _, line := range lines {
		line = strings.TrimSpace(line)
		
		if strings.HasPrefix(line, "# ") {
			// Skip title
			continue
		} else if strings.HasPrefix(line, "> ") {
			// Description
			page.Description = strings.TrimPrefix(line, "> ")
		} else if strings.HasPrefix(line, "- ") {
			// Start new example
			if currentExample != nil {
				page.Examples = append(page.Examples, *currentExample)
			}
			currentExample = &Example{
				Description: strings.TrimPrefix(line, "- "),
			}
			inExample = true
		} else if strings.HasPrefix(line, "`") && strings.HasSuffix(line, "`") && inExample {
			// Command
			command := strings.Trim(line, "`")
			currentExample.Command = command
			currentExample.Placeholders = extractPlaceholders(command)
		} else if line == "" {
			// Empty line ends example
			inExample = false
		}
	}

	// Add last example
	if currentExample != nil {
		page.Examples = append(page.Examples, *currentExample)
	}

	return page, nil
}

// FindBestExample finds the best matching example for a command
func (p *Page) FindBestExample(query string) *Example {
	if len(p.Examples) == 0 {
		return nil
	}

	query = strings.ToLower(query)
	
	// Look for exact match in description
	for _, example := range p.Examples {
		if strings.Contains(strings.ToLower(example.Description), query) {
			return &example
		}
	}

	// Look for partial match
	for _, example := range p.Examples {
		if strings.Contains(strings.ToLower(example.Description), query) {
			return &example
		}
	}

	// Return first example as fallback
	return &p.Examples[0]
}

// Render renders a command with placeholders filled
func (e *Example) Render(vars map[string]string) string {
	command := e.Command
	
	// Replace placeholders with variables
	for _, placeholder := range e.Placeholders {
		value := vars[placeholder.Name]
		if value == "" {
			value = placeholder.Default
		}
		if value == "" {
			value = placeholder.Name // Use placeholder name as fallback
		}
		
		placeholderPattern := regexp.MustCompile(`\{\{` + regexp.QuoteMeta(placeholder.Name) + `\}\}`)
		command = placeholderPattern.ReplaceAllString(command, value)
	}
	
	return command
}

// extractPlaceholders extracts placeholders from a command string
func extractPlaceholders(command string) []Placeholder {
	var placeholders []Placeholder
	
	// Regex to find {{placeholder}} patterns
	re := regexp.MustCompile(`\{\{([^}]+)\}\}`)
	matches := re.FindAllStringSubmatch(command, -1)
	
	seen := make(map[string]bool)
	for _, match := range matches {
		if len(match) > 1 {
			name := match[1]
			if !seen[name] {
				seen[name] = true
				placeholder := Placeholder{
					Name: name,
					Type: inferPlaceholderType(name),
				}
				placeholders = append(placeholders, placeholder)
			}
		}
	}
	
	return placeholders
}

// inferPlaceholderType infers the type of a placeholder based on its name
func inferPlaceholderType(name string) string {
	name = strings.ToLower(name)
	
	switch {
	case strings.Contains(name, "file") || strings.Contains(name, "path"):
		return "file"
	case strings.Contains(name, "dir") || strings.Contains(name, "directory"):
		return "directory"
	case strings.Contains(name, "port"):
		return "port"
	case strings.Contains(name, "num") || strings.Contains(name, "number") || strings.Contains(name, "count"):
		return "number"
	case strings.Contains(name, "url") || strings.Contains(name, "link"):
		return "url"
	case strings.Contains(name, "ip") || strings.Contains(name, "address"):
		return "ip"
	case strings.Contains(name, "user") || strings.Contains(name, "username"):
		return "username"
	case strings.Contains(name, "pass") || strings.Contains(name, "password"):
		return "password"
	case strings.Contains(name, "email"):
		return "email"
	default:
		return "text"
	}
}