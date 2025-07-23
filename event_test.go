package main

import (
	"strings"
	"testing"
)

// TestEventMode tests the event formatting mode
func TestEventMode(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []string
	}{
		{
			"simple key value",
			"key: value",
			[]string{"Event: DOCUMENT-START", "Event: MAPPING-START", "Event: SCALAR", "Value: key", "Value: value", "Event: MAPPING-END"},
		},
		{
			"sequence",
			"items:\n  - one\n  - two",
			[]string{"Event: DOCUMENT-START", "Event: MAPPING-START", "Event: SEQUENCE-START", "Event: SCALAR", "Value: one", "Value: two", "Event: SEQUENCE-END"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			stdout, stderr, err := runCommand(tt.input, "-e")
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

// TestComplexYAMLEvent tests complex YAML structures with event mode
func TestComplexYAMLEvent(t *testing.T) {
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
		"Event: DOCUMENT-START", "Event: MAPPING-START", "Event: SCALAR",
		"Event: SEQUENCE-START", "Event: SEQUENCE-END", "Event: MAPPING-END",
	}

	stdout, stderr, err := runCommand(complexYAML, "-e")
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
