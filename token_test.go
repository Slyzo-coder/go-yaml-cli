package main

import (
	"strings"
	"testing"
)

// TestTokenMode tests the token formatting mode
func TestTokenMode(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []string
	}{
		{
			"simple key value",
			"key: value",
			[]string{"Token: STREAM-START", "Token: BLOCK-MAPPING-START", "Token: KEY", "Token: SCALAR", "Value: \"key\"", "Token: VALUE", "Value: \"value\"", "Token: BLOCK-END", "Token: STREAM-END"},
		},
		{
			"sequence",
			"items:\n  - one\n  - two",
			[]string{"Token: STREAM-START", "Token: BLOCK-MAPPING-START", "Token: BLOCK-SEQUENCE-START", "Token: BLOCK-ENTRY", "Value: \"one\"", "Value: \"two\"", "Token: BLOCK-END"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			stdout, stderr, err := runCommand(tt.input, "-t")
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

// TestComplexYAMLToken tests complex YAML structures with token mode
func TestComplexYAMLToken(t *testing.T) {
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
		"Token: STREAM-START", "Token: DOCUMENT-START", "Token: BLOCK-MAPPING-START",
		"Token: KEY", "Token: SCALAR", "Token: VALUE", "Token: BLOCK-END", "Token: STREAM-END",
	}

	stdout, stderr, err := runCommand(complexYAML, "-t")
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
