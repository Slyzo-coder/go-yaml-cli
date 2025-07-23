package main

import (
	"strings"
	"testing"
)

// TestYAMLMode tests the YAML formatting mode (without preserve)
func TestYAMLMode(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []string
	}{
		{
			"simple key value",
			"key: value",
			[]string{"key: value"},
		},
		{
			"sequence",
			"items:\n  - one\n  - two",
			[]string{"items:", "  - one", "  - two"},
		},
		{
			"mapping",
			"person:\n  name: John\n  age: 30",
			[]string{"person:", "  age: 30", "  name: John"},
		},
		{
			"complex nested structure",
			"data:\n  user:\n    name: Alice\n    roles:\n      - admin\n      - user",
			[]string{"data:", "  user:", "    name: Alice", "    roles:", "      - admin", "      - user"},
		},
		{
			"with comments (should be stripped)",
			"# comment\nkey: value  # inline",
			[]string{"key: value"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			stdout, stderr, err := runCommand(tt.input, "-y")
			if err != nil {
				t.Errorf("Expected no error, got %v", err)
			}
			if stderr != "" {
				t.Errorf("Expected no stderr, got %q", stderr)
			}
			for _, expected := range tt.expected {
				if !strings.Contains(stdout, expected) {
					t.Errorf("Expected output to contain %q, got %q", expected, stdout)
				}
			}
		})
	}
}

// TestYAMLLongFlag tests the long flag version
func TestYAMLLongFlag(t *testing.T) {
	stdout, stderr, err := runCommand("key: value", "--yaml")
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}
	if stderr != "" {
		t.Errorf("Expected no stderr, got %q", stderr)
	}
	if !strings.Contains(stdout, "key: value") {
		t.Errorf("Expected output to contain 'key: value', got %q", stdout)
	}
}

// TestYAMLCapitalYFlag tests the -Y flag (YAML with preserve)
func TestYAMLCapitalYFlag(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []string
	}{
		{
			"with comments (should be preserved)",
			"# comment\nkey: value  # inline",
			[]string{"# comment", "key: value", "# inline"},
		},
		{
			"with styles (should be preserved)",
			"quoted: 'single'\ndouble: \"double\"",
			[]string{"quoted: 'single'", "double: \"double\""},
		},
		{
			"with complex structure",
			"# header\nperson:\n  name: \"John\"  # inline\n  age: 30",
			[]string{"# header", "person:", "name: \"John\"", "# inline", "age: 30"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			stdout, stderr, err := runCommand(tt.input, "-Y")
			if err != nil {
				t.Errorf("Expected no error, got %v", err)
			}
			if stderr != "" {
				t.Errorf("Expected no stderr, got %q", stderr)
			}
			for _, expected := range tt.expected {
				if !strings.Contains(stdout, expected) {
					t.Errorf("Expected output to contain %q, got %q", expected, stdout)
				}
			}
		})
	}
}

// TestYAMLPreserveMode tests the YAML formatting mode with preserve flag
func TestYAMLPreserveMode(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []string
	}{
		{
			"with comments (should be preserved)",
			"# comment\nkey: value  # inline",
			[]string{"# comment", "key: value", "# inline"},
		},
		{
			"with styles (should be preserved)",
			"quoted: 'single'\ndouble: \"double\"",
			[]string{"quoted: 'single'", "double: \"double\""},
		},
		{
			"with complex structure",
			"# header\nperson:\n  name: \"John\"  # inline\n  age: 30",
			[]string{"# header", "person:", "name: \"John\"", "# inline", "age: 30"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			stdout, stderr, err := runCommand(tt.input, "-Y")
			if err != nil {
				t.Errorf("Expected no error, got %v", err)
			}
			if stderr != "" {
				t.Errorf("Expected no stderr, got %q", stderr)
			}
			for _, expected := range tt.expected {
				if !strings.Contains(stdout, expected) {
					t.Errorf("Expected output to contain %q, got %q", expected, stdout)
				}
			}
		})
	}
}

// TestYAMLFormatting tests that YAML is properly formatted
func TestYAMLFormatting(t *testing.T) {
	// Test that the output is properly indented with valid YAML
	input := `person:
  name: John
  age: 30
  hobbies:
    - reading
    - hiking`

	stdout, stderr, err := runCommand(input, "-y")
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}
	if stderr != "" {
		t.Errorf("Expected no stderr, got %q", stderr)
	}

	// Check that the output contains the expected content (order may vary)
	expectedContent := []string{
		"person:", "name: John", "age: 30", "hobbies:", "reading", "hiking",
	}

	for _, expected := range expectedContent {
		if !strings.Contains(stdout, expected) {
			t.Errorf("Expected output to contain %q, got %q", expected, stdout)
		}
	}

	// Check that the output is properly formatted with consistent indentation
	if !strings.Contains(stdout, "  ") {
		t.Errorf("Expected output to have proper indentation, got: %q", stdout)
	}
}

// TestComplexYAMLYAML tests complex YAML structures with YAML mode
func TestComplexYAMLYAML(t *testing.T) {
	complexYAML := `
# Complex YAML test
---
person:
  name: &name John Doe
  age: 30
  hobbies:
    - reading
    - hiking
    - &sport running
  address:
    street: 123 Main St
    city: Anytown
    zip: "12345"
  aliases:
    - *name
    - *sport
---
# Second document
settings:
  debug: true
  log_level: INFO
  features: [feature1, feature2, feature3]
`

	// Test that both documents are processed
	expected := []string{
		"person:", "name: John Doe", "age: 30", "hobbies:", "reading", "hiking", "running",
		"address:", "street: 123 Main St", "city: Anytown", "zip: \"12345\"",
		"aliases:", "settings:", "debug: true", "log_level: INFO", "features:",
	}

	stdout, stderr, err := runCommand(complexYAML, "-y")
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}
	if stderr != "" {
		t.Errorf("Expected no stderr, got %q", stderr)
	}

	// Check that all expected content is present (order may vary)
	for _, expected := range expected {
		if !strings.Contains(stdout, expected) {
			t.Errorf("Expected output to contain %q, got %q", expected, stdout)
		}
	}
}
