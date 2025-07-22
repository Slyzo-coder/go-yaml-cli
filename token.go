// Package main provides YAML token formatting utilities for the go-yaml tool.
package main

import (
	"fmt"
	"os"

	"go.yaml.in/yaml/v3"
)

// ProcessTokens reads YAML from stdin and outputs token information using the internal scanner
func ProcessTokens(profuse, compact bool) error {
	parser, err := yaml.NewParser(os.Stdin)
	if err != nil {
		return fmt.Errorf("error creating parser: %v", err)
	}
	defer parser.Close()

	if compact {
		first := true
		for {
			token, err := parser.Next()
			if err != nil {
				return fmt.Errorf("parser error: %v", err)
			}
			if token == nil {
				break
			}

			if !first {
				fmt.Println()
			}
			first = false

			fmt.Print("- ")
			printTokenCompact(token, profuse)
		}
		// Add final newline for compact mode
		if !first {
			fmt.Println()
		}
	} else {
		for {
			token, err := parser.Next()
			if err != nil {
				return fmt.Errorf("parser error: %v", err)
			}
			if token == nil {
				break
			}

			printToken(token, profuse)
		}
	}

	return nil
}

// printTokenCompact prints a token in compact flow style format
func printTokenCompact(token *yaml.Token, profuse bool) {
	fmt.Print("{Type: ", token.Type)
	if token.Value != "" {
		fmt.Printf(", Value: %q", token.Value)
	}
	if token.Style != "" && token.Style != "Plain" {
		fmt.Printf(", Style: %s", token.Style)
	}
	if profuse {
		if token.StartLine == token.EndLine && token.StartCol == token.EndCol {
			fmt.Printf(", Pos: {%d: %d}", token.StartLine, token.StartCol)
		} else {
			fmt.Printf(", Pos: {%d: %d, %d: %d}", token.StartLine, token.StartCol, token.EndLine, token.EndCol)
		}
	}
	fmt.Print("}")
}

// printToken prints a token in the expected format
func printToken(token *yaml.Token, profuse bool) {
	fmt.Printf("- Token: %v\n", token.Type)

	if token.Value != "" {
		fmt.Printf("  Value: %q\n", token.Value)
	}

	if token.Style != "" && token.Style != "Plain" {
		fmt.Printf("  Style: %s\n", token.Style)
	}

	if profuse {
		if token.StartLine == token.EndLine && token.StartCol == token.EndCol {
			fmt.Printf("  Pos: {%d: %d}\n", token.StartLine, token.StartCol)
		} else {
			fmt.Printf("  Pos: {%d: %d, %d: %d}\n", token.StartLine, token.StartCol, token.EndLine, token.EndCol)
		}
	}
	fmt.Println()
}
