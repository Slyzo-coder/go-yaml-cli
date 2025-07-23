package main

import (
	"strings"
	"testing"
)

// TestAllModesWithInlineYAML tests all modes with inline YAML data
func TestAllModesWithInlineYAML(t *testing.T) {
	testYAML := "key: value\nlist: [1, 2, 3]\nquoted: \"text\""
	modes := []struct {
		name  string
		flags []string
	}{
		{"node", []string{"-n"}},
		{"event", []string{"-e"}},
		{"token", []string{"-t"}},
		{"json", []string{"-j"}},
		{"json_pretty", []string{"-J"}},
		{"json_capital_j", []string{"-J"}},
		{"yaml", []string{"-y"}},
		{"yaml_preserve", []string{"-Y"}},
		{"yaml_capital_y", []string{"-Y"}},
	}

	for _, mode := range modes {
		t.Run(mode.name, func(t *testing.T) {
			stdout, stderr, err := runCommand(testYAML, mode.flags...)
			if err != nil {
				t.Errorf("Expected no error, got %v", err)
			}
			if stderr != "" {
				t.Errorf("Expected no stderr, got %q", stderr)
			}
			if stdout == "" {
				t.Errorf("Expected non-empty output")
			}
		})
	}
}

// TestComplexYAMLAllModes tests complex YAML structures with all modes
func TestComplexYAMLAllModes(t *testing.T) {
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

	tests := []struct {
		name     string
		flags    []string
		expected []string
	}{
		{
			"node mode",
			[]string{"-n"},
			[]string{"kind: Document", "kind: Mapping", "text: person", "text: name", "text: age", "text: hobbies", "text: address", "text: aliases"},
		},
		{
			"event mode",
			[]string{"-e"},
			[]string{"Event: DOCUMENT-START", "Event: MAPPING-START", "Event: SCALAR", "Event: SEQUENCE-START", "Event: SEQUENCE-END", "Event: MAPPING-END"},
		},
		{
			"token mode",
			[]string{"-t"},
			[]string{"Token: STREAM-START", "Token: DOCUMENT-START", "Token: BLOCK-MAPPING-START", "Token: KEY", "Token: SCALAR", "Token: VALUE", "Token: BLOCK-END", "Token: STREAM-END"},
		},
		{
			"json mode",
			[]string{"-j"},
			[]string{`"person":{`, `"name":"John Doe"`, `"age":30`, `"hobbies":[`, `"reading"`, `"hiking"`, `"running"`},
		},
		{
			"json mode with pretty",
			[]string{"-J"},
			[]string{`"person": {`, `"name": "John Doe"`, `"age": 30`, `"hobbies": [`, `"reading"`, `"hiking"`, `"running"`},
		},
		{
			"json mode with capital J",
			[]string{"-J"},
			[]string{`"person": {`, `"name": "John Doe"`, `"age": 30`, `"hobbies": [`, `"reading"`, `"hiking"`, `"running"`},
		},
		{
			"yaml mode",
			[]string{"-y"},
			[]string{"person:", "name: John Doe", "age: 30", "hobbies:", "reading", "hiking", "running"},
		},
		{
			"yaml mode with preserve",
			[]string{"-Y"},
			[]string{"person:", "name: &name John Doe", "age: 30", "hobbies:", "reading", "hiking", "running"},
		},
		{
			"yaml mode with capital Y",
			[]string{"-Y"},
			[]string{"person:", "name: &name John Doe", "age: 30", "hobbies:", "reading", "hiking", "running"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			stdout, stderr, err := runCommand(complexYAML, tt.flags...)
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
