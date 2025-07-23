package main

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

// TestFileBasedTests runs tests using the test directory structure
func TestFileBasedTests(t *testing.T) {
	testDirs, err := filepath.Glob("test/*")
	if err != nil {
		t.Fatalf("Failed to find test directories: %v", err)
	}

	for _, testDir := range testDirs {
		// Skip if not a directory
		if info, err := os.Stat(testDir); err != nil || !info.IsDir() {
			continue
		}

		testName := filepath.Base(testDir)
		t.Run(testName, func(t *testing.T) {
			runFileBasedTest(t, testDir)
		})
	}
}

// runFileBasedTest runs tests for a specific test directory
func runFileBasedTest(t *testing.T, testDir string) {
	// Read input file
	inputFile := filepath.Join(testDir, "in.yaml")
	inputData, err := os.ReadFile(inputFile)
	if err != nil {
		t.Fatalf("Failed to read input file %s: %v", inputFile, err)
	}

	// Find all output files
	outputFiles, err := filepath.Glob(filepath.Join(testDir, "out-*.yaml"))
	if err != nil {
		t.Fatalf("Failed to find output files in %s: %v", testDir, err)
	}

	for _, outputFile := range outputFiles {
		// Extract flags from filename (e.g., "out-e-p.yaml" -> ["-e", "-p"])
		baseName := filepath.Base(outputFile)
		flagPart := strings.TrimPrefix(baseName, "out-")
		flagPart = strings.TrimSuffix(flagPart, ".yaml")

		flags := parseFlagsFromFilename(flagPart)

		t.Run(flagPart, func(t *testing.T) {
			// Read expected output
			expectedData, err := os.ReadFile(outputFile)
			if err != nil {
				t.Fatalf("Failed to read expected output file %s: %v", outputFile, err)
			}
			expected := strings.TrimSpace(string(expectedData))

			// Run the command
			stdout, stderr, err := runCommand(string(inputData), flags...)
			if err != nil {
				t.Errorf("Command failed with flags %v: %v", flags, err)
				return
			}
			if stderr != "" {
				t.Errorf("Command produced stderr with flags %v: %q", flags, stderr)
				return
			}

			actual := strings.TrimSpace(stdout)
			if actual != expected {
				t.Errorf("Output mismatch for flags %v\nExpected:\n%s\n\nGot:\n%s", flags, expected, actual)
			}
		})
	}
}

// parseFlagsFromFilename converts a filename part like "e-p-c" to flags like ["-e", "-p", "-c"]
func parseFlagsFromFilename(flagPart string) []string {
	if flagPart == "" {
		return []string{}
	}

	parts := strings.Split(flagPart, "-")
	var flags []string

	for _, part := range parts {
		switch part {
		case "n":
			flags = append(flags, "-n")
		case "e":
			flags = append(flags, "-e")
		case "t":
			flags = append(flags, "-t")
		case "j":
			flags = append(flags, "-j")
		case "y":
			flags = append(flags, "-y")
		case "p":
			flags = append(flags, "-p")
		case "l":
			flags = append(flags, "-l")
		case "E":
			flags = append(flags, "-E")
		case "T":
			flags = append(flags, "-T")
		case "J":
			flags = append(flags, "-J")
		case "Y":
			flags = append(flags, "-Y")
		}
	}

	return flags
}
