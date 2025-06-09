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
	showHelp := flag.Bool("h", false, "Show help information")
	showVersion := flag.Bool("version", false, "Show version information")

	// YAML modes
	yamlMode := flag.Bool("y", false, "YAML formatting mode")
	yamlPreserveMode := flag.Bool("Y", false, "YAML Preserve; short for -y -p")

	// JSON modes
	jsonMode := flag.Bool("j", false, "JSON formatting mode (compact)")
	jsonPrettyMode := flag.Bool("J", false, "JSON Pretty; short for -j -p")

	// Token modes
	tokenMode := flag.Bool("t", false, "Token formatting mode")
	tokenProfuseMode := flag.Bool("T", false, "Token Profuse; short for -t -p")

	// Event modes
	eventMode := flag.Bool("e", false, "Event formatting mode")
	eventProfuseMode := flag.Bool("E", false, "Event Profuse; short for -e -p")

	// Node mode
	nodeMode := flag.Bool("n", false, "Node formatting mode")

	// Shared flags
	preserveMode := flag.Bool("p", false, "Preserve comments and styles (with -y)")
	prettyMode := flag.Bool("pretty", false, "Pretty JSON output (with -j)")
	profuseMode := flag.Bool("profuse", false, "Show line info for --token and --event")
	compactMode := flag.Bool("c", false, "Compact output (flow style, no blank lines)")

	// Long flag aliases
	flag.BoolVar(showHelp, "help", false, "Show help information")
	flag.BoolVar(yamlMode, "yaml", false, "YAML formatting mode")
	flag.BoolVar(yamlPreserveMode, "YAML", false, "YAML Preserve; short for -y -p")
	flag.BoolVar(jsonMode, "json", false, "JSON formatting mode (compact)")
	flag.BoolVar(jsonPrettyMode, "JSON", false, "JSON Pretty; short for -j -p")
	flag.BoolVar(tokenMode, "token", false, "Token formatting mode")
	flag.BoolVar(tokenProfuseMode, "TOKEN", false, "Token Profuse; short for -t -p")
	flag.BoolVar(eventMode, "event", false, "Event formatting mode")
	flag.BoolVar(eventProfuseMode, "EVENT", false, "Event Profuse; short for -e -p")
	flag.BoolVar(nodeMode, "node", false, "Node formatting mode")
	flag.BoolVar(preserveMode, "preserve", false, "Preserve comments and styles (with -y)")
	flag.BoolVar(compactMode, "compact", false, "Compact output (flow style, no blank lines)")

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
	if (stat.Mode()&os.ModeCharDevice) != 0 && !*nodeMode && !*eventMode && !*eventProfuseMode && !*tokenMode && !*tokenProfuseMode && !*jsonMode && !*jsonPrettyMode && !*yamlMode && !*yamlPreserveMode && !*preserveMode && !*compactMode {
		printHelp()
		return
	}

	// Error if stdin has data but no mode flags are provided
	if (stat.Mode()&os.ModeCharDevice) == 0 && !*nodeMode && !*eventMode && !*eventProfuseMode && !*tokenMode && !*tokenProfuseMode && !*jsonMode && !*jsonPrettyMode && !*yamlMode && !*yamlPreserveMode && !*preserveMode && !*compactMode {
		fmt.Fprintf(os.Stderr, "Error: stdin has data but no mode specified. Use -n/--node, -e/--event, -E/--EVENT, -t/--token, -T/--TOKEN, -j/--json, -J/--JSON, -y/--yaml, -Y/--YAML, or --preserve flag.\n")
		os.Exit(1)
	}

	// Process YAML input
	if *eventMode {
		// Use event formatting mode
		profuse := *profuseMode || *preserveMode // -p means profuse for event mode
		if err := ProcessEvents(profuse, *compactMode); err != nil {
			log.Fatal("Failed to process events:", err)
		}
	} else if *eventProfuseMode {
		// Use event formatting mode with profuse output
		if err := ProcessEvents(true, *compactMode); err != nil {
			log.Fatal("Failed to process events:", err)
		}
	} else if *tokenMode {
		// Use token formatting mode
		profuse := *profuseMode || *preserveMode // -p means profuse for token mode
		if err := ProcessTokens(profuse, *compactMode); err != nil {
			log.Fatal("Failed to process tokens:", err)
		}
	} else if *tokenProfuseMode {
		// Use token formatting mode with profuse output
		if err := ProcessTokens(true, *compactMode); err != nil {
			log.Fatal("Failed to process tokens:", err)
		}
	} else if *jsonMode {
		// Use JSON formatting mode (compact by default)
		pretty := *prettyMode || *preserveMode // -p means pretty for JSON mode
		if err := ProcessJSON(pretty); err != nil {
			log.Fatal("Failed to process JSON:", err)
		}
	} else if *jsonPrettyMode {
		// Use pretty JSON formatting mode
		if err := ProcessJSON(true); err != nil {
			log.Fatal("Failed to process JSON:", err)
		}
	} else if *yamlMode {
		// Use YAML formatting mode (clean by default)
		preserve := *preserveMode // -p means preserve for YAML mode
		if err := ProcessYAML(preserve); err != nil {
			log.Fatal("Failed to process YAML:", err)
		}
	} else if *yamlPreserveMode {
		// Use YAML formatting mode with preserve
		if err := ProcessYAML(true); err != nil {
			log.Fatal("Failed to process YAML:", err)
		}
	} else if *preserveMode {
		// Use YAML formatting mode with preserve (--preserve flag)
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
  -y, --yaml       YAML formatting mode
  -Y, --YAML       YAML Preserve; short for -y -p

  -j, --json       JSON formatting mode (compact)
  -J, --JSON       JSON Pretty; short for -j -p

  -t, --token      Token formatting mode
  -T, --TOKEN      Token Profuse; short for -t -p

  -e, --event      Event formatting mode
  -E, --EVENT      Event Profuse; short for -e -p

  -n, --node       Node formatting mode

  -p, --preserve   Preserve comments and styles (with -y)
  -p, --pretty     Pretty JSON output (with -j)
  -p, --profuse    Show line info for --token and --event
  -c, --compact    Remove blanks lines in formatting

  -h, --help       Show this help information
  --version        Show version information

`, version)
}
