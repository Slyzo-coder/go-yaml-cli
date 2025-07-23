package main

import (
	"strings"
	"testing"
)

// TestJSONMode tests the JSON formatting mode (compact by default)
func TestJSONMode(t *testing.T) {
	tests := []struct {
		name        string
		input       string
		expected    []string
		expectError bool
	}{
		{
			"simple key value",
			"key: value",
			[]string{`"key":"value"`, "{"},
			false,
		},
		{
			"sequence",
			"items:\n  - one\n  - two",
			[]string{`"items":[`, `"one"`, `"two"`, "]"},
			false,
		},
		{
			"mapping",
			"person:\n  name: John\n  age: 30",
			[]string{`"person":{`, `"name":"John"`, `"age":30`, "}"},
			false,
		},
		{
			"mixed types",
			"data:\n  string: hello\n  number: 42\n  boolean: true\n  null: null",
			[]string{},
			true, // This should error because map[interface{}]interface{} can't be JSON encoded
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			stdout, stderr, err := runCommand(tt.input, "-j")

			if tt.expectError {
				if err == nil {
					t.Errorf("Expected error, got none")
				}
				if !strings.Contains(stderr, "failed to encode JSON") {
					t.Errorf("Expected JSON encoding error, got %q", stderr)
				}
			} else {
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
			}
		})
	}
}

// TestComplexYAMLJSON tests complex YAML structures with JSON mode
func TestComplexYAMLJSON(t *testing.T) {
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
		`"person":{`, `"name":"John Doe"`, `"age":30`, `"hobbies":[`,
		`"reading"`, `"hiking"`, `"running"`,
	}

	stdout, stderr, err := runCommand(complexYAML, "-j")
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

// TestPrettyJSONMode tests the pretty JSON formatting mode
func TestPrettyJSONMode(t *testing.T) {
	tests := []struct {
		name        string
		input       string
		expected    []string
		expectError bool
	}{
		{
			"simple key value",
			"key: value",
			[]string{`"key": "value"`, "{"},
			false,
		},
		{
			"sequence",
			"items:\n  - one\n  - two",
			[]string{`"items": [`, `"one"`, `"two"`, "]"},
			false,
		},
		{
			"mapping",
			"person:\n  name: John\n  age: 30",
			[]string{`"person": {`, `"name": "John"`, `"age": 30`, "}"},
			false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			stdout, stderr, err := runCommand(tt.input, "-J")

			if tt.expectError {
				if err == nil {
					t.Errorf("Expected error, got none")
				}
				if !strings.Contains(stderr, "failed to encode JSON") {
					t.Errorf("Expected JSON encoding error, got %q", stderr)
				}
			} else {
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
			}
		})
	}
}

// TestCapitalJFlag tests the -J flag (pretty JSON)
func TestCapitalJFlag(t *testing.T) {
	tests := []struct {
		name        string
		input       string
		expected    []string
		expectError bool
	}{
		{
			"simple key value",
			"key: value",
			[]string{`"key": "value"`, "{"},
			false,
		},
		{
			"sequence",
			"items:\n  - one\n  - two",
			[]string{`"items": [`, `"one"`, `"two"`, "]"},
			false,
		},
		{
			"mapping",
			"person:\n  name: John\n  age: 30",
			[]string{`"person": {`, `"name": "John"`, `"age": 30`, "}"},
			false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			stdout, stderr, err := runCommand(tt.input, "-J")

			if tt.expectError {
				if err == nil {
					t.Errorf("Expected error, got none")
				}
				if !strings.Contains(stderr, "failed to encode JSON") {
					t.Errorf("Expected JSON encoding error, got %q", stderr)
				}
			} else {
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
			}
		})
	}
}
