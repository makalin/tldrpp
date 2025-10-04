package types

import (
	"testing"
)

func TestParsePage(t *testing.T) {
	content := `# tar

> Archive utility.

- Extract archive:
  \`tar -xf {{file}}\`

- List contents:
  \`tar -tf {{file}}\`
`

	entry := IndexEntry{
		Name:        "tar",
		Description: "Archive utility",
		Platform:    "linux",
	}

	page, err := ParsePage(content, entry)
	if err != nil {
		t.Fatalf("ParsePage failed: %v", err)
	}

	if page.Name != "tar" {
		t.Errorf("Expected name 'tar', got '%s'", page.Name)
	}

	if page.Description != "Archive utility" {
		t.Errorf("Expected description 'Archive utility', got '%s'", page.Description)
	}

	if page.Platform != "linux" {
		t.Errorf("Expected platform 'linux', got '%s'", page.Platform)
	}

	if len(page.Examples) != 2 {
		t.Errorf("Expected 2 examples, got %d", len(page.Examples))
	}

	if page.Examples[0].Description != "Extract archive" {
		t.Errorf("Expected first example description 'Extract archive', got '%s'", page.Examples[0].Description)
	}

	if page.Examples[0].Command != "tar -xf {{file}}" {
		t.Errorf("Expected first example command 'tar -xf {{file}}', got '%s'", page.Examples[0].Command)
	}
}

func TestFindBestExample(t *testing.T) {
	examples := []Example{
		{Description: "Extract archive", Command: "tar -xf {{file}}"},
		{Description: "List contents", Command: "tar -tf {{file}}"},
	}

	page := &Page{
		Name:        "tar",
		Description: "Archive utility",
		Platform:    "linux",
		Examples:    examples,
	}

	// Test with generic query
	best := page.FindBestExample("tar")
	if best == nil {
		t.Fatal("Expected to find an example")
	}

	if best.Description != "Extract archive" {
		t.Errorf("Expected 'Extract archive', got '%s'", best.Description)
	}

	// Test with empty examples
	emptyPage := &Page{
		Name:        "empty",
		Description: "Empty page",
		Platform:    "linux",
		Examples:    []Example{},
	}

	best = emptyPage.FindBestExample("query")
	if best != nil {
		t.Error("Expected nil for empty examples")
	}
}

func TestExampleRender(t *testing.T) {
	example := Example{
		Description: "Extract archive",
		Command:     "tar -xf {{file}}",
		Placeholders: []Placeholder{
			{Name: "file", Type: "file", Default: "archive.tar.gz"},
		},
	}

	vars := map[string]string{
		"file": "test.tar.gz",
	}

	result := example.Render(vars)
	expected := "tar -xf test.tar.gz"

	if result != expected {
		t.Errorf("Expected '%s', got '%s'", expected, result)
	}
}

func TestExampleRenderWithDefault(t *testing.T) {
	example := Example{
		Description: "Extract archive",
		Command:     "tar -xf {{file}}",
		Placeholders: []Placeholder{
			{Name: "file", Type: "file", Default: "archive.tar.gz"},
		},
	}

	result := example.Render(map[string]string{})
	expected := "tar -xf archive.tar.gz"

	if result != expected {
		t.Errorf("Expected '%s', got '%s'", expected, result)
	}
}

func TestExtractPlaceholders(t *testing.T) {
	tests := []struct {
		command     string
		expected    []Placeholder
		description string
	}{
		{
			command:     "tar -xf {{file}}",
			expected:    []Placeholder{{Name: "file", Type: "file"}},
			description: "single placeholder",
		},
		{
			command:     "cp {{src}} {{dest}}",
			expected:    []Placeholder{{Name: "src", Type: "text"}, {Name: "dest", Type: "text"}},
			description: "multiple placeholders",
		},
		{
			command:     "ls -la",
			expected:    []Placeholder{},
			description: "no placeholders",
		},
		{
			command:     "tar -xf {{file}} --output={{dest}}",
			expected:    []Placeholder{{Name: "file", Type: "file"}, {Name: "dest", Type: "text"}},
			description: "placeholders with other text",
		},
	}

	for _, test := range tests {
		t.Run(test.description, func(t *testing.T) {
			placeholders := extractPlaceholders(test.command)

			if len(placeholders) != len(test.expected) {
				t.Errorf("Expected %d placeholders, got %d", len(test.expected), len(placeholders))
				return
			}

			for i, expected := range test.expected {
				if placeholders[i].Name != expected.Name {
					t.Errorf("Expected placeholder name '%s', got '%s'", expected.Name, placeholders[i].Name)
				}
				if placeholders[i].Type != expected.Type {
					t.Errorf("Expected placeholder type '%s', got '%s'", expected.Type, placeholders[i].Type)
				}
			}
		})
	}
}

func TestInferPlaceholderType(t *testing.T) {
	tests := []struct {
		name     string
		expected string
	}{
		{"file", "file"},
		{"input_file", "file"},
		{"file_path", "file"},
		{"directory", "directory"},
		{"dir", "directory"},
		{"port", "port"},
		{"port_number", "port"},
		{"number", "number"},
		{"count", "number"},
		{"num", "number"},
		{"url", "url"},
		{"link", "url"},
		{"ip", "ip"},
		{"address", "ip"},
		{"username", "username"},
		{"user", "username"},
		{"password", "password"},
		{"pass", "password"},
		{"email", "email"},
		{"unknown", "text"},
		{"random", "text"},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			result := inferPlaceholderType(test.name)
			if result != test.expected {
				t.Errorf("Expected type '%s' for '%s', got '%s'", test.expected, test.name, result)
			}
		})
	}
}