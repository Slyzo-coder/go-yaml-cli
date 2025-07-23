// Package main provides a YAML node inspection tool that reads YAML from stdin
// and outputs a detailed analysis of its node structure, including comments
// and content organization.
package main

import (
	"bytes"
	"errors"
	"flag"
	"fmt"
	"go.yaml.in/yaml/v3"
	"io"
	"log"
	"os"
)

const version = "3.0.3.1"

// main reads YAML from stdin, parses it, and outputs the node structure
func main() {
	// Parse command line flags
	showHelp := flag.Bool("h", false, "Show this help information")
	showVersion := flag.Bool("version", false, "Show version information")

	// YAML modes
	yamlMode := flag.Bool("y", false, "YAML encoding output")
	yamlPreserveMode := flag.Bool("Y", false, "YAML style and comments preserved")

	// JSON modes
	jsonMode := flag.Bool("j", false, "JSON compact output")
	jsonPrettyMode := flag.Bool("J", false, "JSON pretty output")

	// Token modes
	tokenMode := flag.Bool("t", false, "Token output")
	tokenProfuseMode := flag.Bool("T", false, "Token with line info")

	// Event modes
	eventMode := flag.Bool("e", false, "Event output")
	eventProfuseMode := flag.Bool("E", false, "Event with line info")

	// Node mode
	nodeMode := flag.Bool("n", false, "Node representation output")

	// Shared flags
	longMode := flag.Bool("l", false, "Long (block) formatted output")

	// Long flag aliases
	flag.BoolVar(showHelp, "help", false, "Show this help information")
	flag.BoolVar(yamlMode, "yaml", false, "YAML encoding output")
	flag.BoolVar(yamlPreserveMode, "YAML", false, "YAML style and comments preserved")
	flag.BoolVar(jsonMode, "json", false, "JSON compact output")
	flag.BoolVar(jsonPrettyMode, "JSON", false, "JSON pretty output")
	flag.BoolVar(tokenMode, "token", false, "Token output")
	flag.BoolVar(tokenProfuseMode, "TOKEN", false, "Token with line info")
	flag.BoolVar(eventMode, "event", false, "Event output")
	flag.BoolVar(eventProfuseMode, "EVENT", false, "Event with line info")
	flag.BoolVar(nodeMode, "node", false, "Node representation output")
	flag.BoolVar(longMode, "long", false, "Long (block) formatted output")

	flag.Parse()

	// Show version and exit
	if *showVersion {
		fmt.Printf("go-yaml version %s\n", version)
		return
	}

	// Show help and exit
	if *showHelp {
		printHelp()
		return
	}

	// Check if stdin has data
	stat, err := os.Stdin.Stat()
	if err != nil {
		log.Fatal("Failed to stat stdin:", err)
	}

	// If no stdin and no flags, show help
	if (stat.Mode()&os.ModeCharDevice) != 0 && !*nodeMode && !*eventMode && !*eventProfuseMode && !*tokenMode && !*tokenProfuseMode && !*jsonMode && !*jsonPrettyMode && !*yamlMode && !*yamlPreserveMode && !*longMode {
		printHelp()
		return
	}

	// Error if stdin has data but no mode flags are provided
	if (stat.Mode()&os.ModeCharDevice) == 0 && !*nodeMode && !*eventMode && !*eventProfuseMode && !*tokenMode && !*tokenProfuseMode && !*jsonMode && !*jsonPrettyMode && !*yamlMode && !*yamlPreserveMode && !*longMode {
		fmt.Fprintf(os.Stderr, "Error: stdin has data but no mode specified. Use -n/--node, -e/--event, -E/--EVENT, -t/--token, -T/--TOKEN, -j/--json, -J/--JSON, -y/--yaml, -Y/--YAML flag.\n")
		os.Exit(1)
	}

	// Process YAML input
	if *eventMode {
		// Use event formatting mode (compact by default)
		compact := !*longMode // compact is default, long mode negates it
		if err := ProcessEvents(false, compact); err != nil {
			log.Fatal("Failed to process events:", err)
		}
	} else if *eventProfuseMode {
		// Use event formatting mode with profuse output
		compact := !*longMode // compact is default, long mode negates it
		if err := ProcessEvents(true, compact); err != nil {
			log.Fatal("Failed to process events:", err)
		}
	} else if *tokenMode {
		// Use token formatting mode (compact by default)
		compact := !*longMode // compact is default, long mode negates it
		if err := ProcessTokens(false, compact); err != nil {
			log.Fatal("Failed to process tokens:", err)
		}
	} else if *tokenProfuseMode {
		// Use token formatting mode with profuse output
		compact := !*longMode // compact is default, long mode negates it
		if err := ProcessTokens(true, compact); err != nil {
			log.Fatal("Failed to process tokens:", err)
		}
	} else if *jsonMode {
		// Use JSON formatting mode (compact by default)
		if err := ProcessJSON(false); err != nil {
			log.Fatal("Failed to process JSON:", err)
		}
	} else if *jsonPrettyMode {
		// Use pretty JSON formatting mode
		if err := ProcessJSON(true); err != nil {
			log.Fatal("Failed to process JSON:", err)
		}
	} else if *yamlMode {
		// Use YAML formatting mode (clean by default)
		if err := ProcessYAML(false); err != nil {
			log.Fatal("Failed to process YAML:", err)
		}
	} else if *yamlPreserveMode {
		// Use YAML formatting mode with preserve
		if err := ProcessYAML(true); err != nil {
			log.Fatal("Failed to process YAML:", err)
		}
	} else {
		// Use node formatting mode (default)
		reader := io.Reader(os.Stdin)
		dec := yaml.NewDecoder(reader)
		firstDoc := true

		for {
			var node yaml.Node
			err := dec.Decode(&node)
			if errors.Is(err, io.EOF) {
				break
			}
			if err != nil {
				log.Fatal("Failed to load YAML node:", err)
			}

			// Add document separator for all documents except the first
			if !firstDoc {
				fmt.Println("---")
			}
			firstDoc = false

			info := FormatNode(node)

			// Use encoder with 2-space indentation
			var buf bytes.Buffer
			enc := yaml.NewEncoder(&buf)
			enc.SetIndent(2)
			if err := enc.Encode(info); err != nil {
				log.Fatal("Failed to marshal node info:", err)
			}
			enc.Close()
			fmt.Print(buf.String())
		}
	}
}

// printHelp displays the help information for the program
func printHelp() {
	fmt.Printf(`go-yaml version %s

A tool to show how go.yaml.in/yaml/v3 handles YAML both internally and
externally.

Usage:
  go-yaml [options] < input.yaml

Options:
  -y, --yaml       YAML encoding output
  -Y, --YAML       YAML style and comments preserved

  -j, --json       JSON compact output
  -J, --JSON       JSON pretty output

  -t, --token      Token output
  -T, --TOKEN      Token with line info

  -e, --event      Event output
  -E, --EVENT      Event with line info

  -n, --node       Node representation output

  -l, --long       Long (block) formatted output

  -h, --help       Show this help information
  --version        Show version information

`, version)
}
