package main

import (
	"strings"
	"testing"
)

// TestNodeMode tests the node formatting mode
func TestNodeMode(t *testing.T) {
	stdout, stderr, err := runCommand("key: value\nlist: [1, 2, 3]", "-n")
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}
	if stderr != "" {
		t.Errorf("Expected no stderr, got %q", stderr)
	}

	// Check for expected content
	expected := []string{
		"kind: Document",
		"kind: Mapping",
		"text: key",
		"text: value",
		"kind: Sequence",
		"text: \"1\"",
		"text: \"2\"",
		"text: \"3\"",
	}

	for _, exp := range expected {
		if !strings.Contains(stdout, exp) {
			t.Errorf("Expected output to contain %q", exp)
		}
	}
}

// TestComplexYAMLNode tests complex YAML with node mode
func TestComplexYAMLNode(t *testing.T) {
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

	expected := []string{
		"kind: Document", "kind: Mapping", "text: person", "text: name",
		"text: age", "text: hobbies", "text: address", "text: aliases",
	}

	stdout, stderr, err := runCommand(complexYAML, "-n")
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}
	if stderr != "" {
		t.Errorf("Expected no stderr, got %q", stderr)
	}
	for _, expected := range expected {
		if !strings.Contains(stdout, expected) {
			t.Errorf("Expected output to contain %q, got %q", expected, stdout)
		}
	}
}
