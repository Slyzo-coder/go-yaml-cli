package main

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"strings"
	"testing"
)

// TestMain sets up the test environment
func TestMain(m *testing.M) {
	// Build the program before running tests
	cmd := exec.Command("go", "build")
	if err := cmd.Run(); err != nil {
		fmt.Printf("Failed to build program: %v\n", err)
		os.Exit(1)
	}

	os.Exit(m.Run())
}

// runCommand runs the go-yaml program with given input and flags
func runCommand(input string, flags ...string) (string, string, error) {
	cmd := exec.Command("./go-yaml", flags...)

	var stdin bytes.Buffer
	var stdout bytes.Buffer
	var stderr bytes.Buffer

	if input != "" {
		stdin.WriteString(input)
		cmd.Stdin = &stdin
	}

	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	err := cmd.Run()
	return stdout.String(), stderr.String(), err
}

// TestHelpFlag tests the help flag functionality
func TestHelpFlag(t *testing.T) {
	tests := []struct {
		name     string
		flags    []string
		expected string
	}{
		{"-h flag", []string{"-h"}, "go-yaml version 3.0.3.1"},
		{"--help flag", []string{"--help"}, "go-yaml version 3.0.3.1"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			stdout, stderr, err := runCommand("", tt.flags...)
			if err != nil {
				t.Errorf("Expected no error, got %v", err)
			}
			if !strings.Contains(stdout, tt.expected) {
				t.Errorf("Expected output to contain %q, got %q", tt.expected, stdout)
			}
			if stderr != "" {
				t.Errorf("Expected no stderr, got %q", stderr)
			}
		})
	}
}

// TestVersionFlag tests the version flag functionality
func TestVersionFlag(t *testing.T) {
	stdout, stderr, err := runCommand("", "--version")
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}
	if !strings.Contains(stdout, "3.0.3.1") {
		t.Errorf("Expected output to contain version, got %q", stdout)
	}
	if stderr != "" {
		t.Errorf("Expected no stderr, got %q", stderr)
	}
}

// TestNoInputShowsHelp tests that running without input shows help
func TestNoInputShowsHelp(t *testing.T) {
	stdout, stderr, err := runCommand("")
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}
	if !strings.Contains(stdout, "go-yaml version 3.0.3.1") {
		t.Errorf("Expected help output, got %q", stdout)
	}
	if stderr != "" {
		t.Errorf("Expected no stderr, got %q", stderr)
	}
}

// TestErrorWhenInputButNoMode tests error when input is provided but no mode flag
func TestErrorWhenInputButNoMode(t *testing.T) {
	stdout, stderr, err := runCommand("key: value")
	if err == nil {
		t.Errorf("Expected error, got none")
	}
	if !strings.Contains(stderr, "Error: stdin has data but no mode specified") {
		t.Errorf("Expected error message, got %q", stderr)
	}
	if stdout != "" {
		t.Errorf("Expected no stdout, got %q", stdout)
	}
}

// TestLongFlags tests the long flag versions
func TestLongFlags(t *testing.T) {
	tests := []struct {
		name     string
		flags    []string
		expected string
	}{
		{"--node flag", []string{"--node"}, "kind: Document"},
		{"--event flag", []string{"--event"}, "Event: DOCUMENT-START"},
		{"--token flag", []string{"--token"}, "Token: STREAM-START"},
		{"--json flag", []string{"--json"}, `"key":"value"`},
		{"--yaml flag", []string{"--yaml"}, "key: value"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			stdout, stderr, err := runCommand("key: value", tt.flags...)
			if err != nil {
				t.Errorf("Expected no error, got %v", err)
			}
			if stderr != "" {
				t.Errorf("Expected no stderr, got %q", stderr)
			}
			if !strings.Contains(stdout, tt.expected) {
				t.Errorf("Expected output to contain %q, got %q", tt.expected, stdout)
			}
		})
	}
}

// TestInvalidYAML tests error handling with invalid YAML
func TestInvalidYAML(t *testing.T) {
	invalidYAML := "key: value\n  invalid: indentation"

	_, _, err := runCommand(invalidYAML, "-n")
	// We expect this to fail due to invalid YAML
	if err == nil {
		t.Errorf("Expected error for invalid YAML, got none")
	}
}
