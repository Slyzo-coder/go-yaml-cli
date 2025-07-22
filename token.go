// Package main provides YAML token formatting utilities for the go-yaml tool.
package main

import (
	"fmt"
	"io"
	"os"
	"strings"

	"go.yaml.in/yaml/v3"
)

// Token represents a YAML token with comment information
type Token struct {
	Type        string
	Value       string
	Style       string
	StartLine   int
	StartColumn int
	EndLine     int
	EndColumn   int
	HeadComment string
	LineComment string
	FootComment string
}

// ProcessTokens reads YAML from stdin and outputs token information using the internal scanner
func ProcessTokens(profuse, compact bool) error {
	decoder := yaml.NewDecoder(os.Stdin)

	if compact {
		first := true
		for {
			var node yaml.Node
			err := decoder.Decode(&node)
			if err != nil {
				if err == io.EOF {
					break
				}
				return fmt.Errorf("failed to decode YAML: %v", err)
			}

			tokens := processNodeToTokens(&node, profuse)
			for _, token := range tokens {
				if !first {
					fmt.Println()
				}
				first = false

				fmt.Print("- ")
				printTokenCompact(token, profuse)
			}
		}
		// Add final newline for compact mode
		if !first {
			fmt.Println()
		}
	} else {
		for {
			var node yaml.Node
			err := decoder.Decode(&node)
			if err != nil {
				if err == io.EOF {
					break
				}
				return fmt.Errorf("failed to decode YAML: %v", err)
			}

			tokens := processNodeToTokens(&node, profuse)
			for _, token := range tokens {
				printToken(token, profuse)
			}
		}
	}

	return nil
}

// printTokenCompact prints a token in compact flow style format
func printTokenCompact(token *Token, profuse bool) {
	fmt.Print("{Type: ", token.Type)
	if token.Value != "" {
		fmt.Printf(", Value: %q", token.Value)
	}
	if token.Style != "" && token.Style != "Plain" {
		fmt.Printf(", Style: %s", token.Style)
	}
	if token.HeadComment != "" {
		fmt.Printf(", Head: %q", token.HeadComment)
	}
	if token.LineComment != "" {
		fmt.Printf(", Line: %q", token.LineComment)
	}
	if token.FootComment != "" {
		fmt.Printf(", Foot: %q", token.FootComment)
	}
	if profuse {
		if token.StartLine == token.EndLine && token.StartColumn == token.EndColumn {
			fmt.Printf(", Pos: {%d: %d}", token.StartLine, token.StartColumn)
		} else {
			fmt.Printf(", Pos: {%d: %d, %d: %d}", token.StartLine, token.StartColumn, token.EndLine, token.EndColumn)
		}
	}
	fmt.Print("}")
}

// printToken prints a token in the expected format
func printToken(token *Token, profuse bool) {
	fmt.Printf("- Token: %v\n", token.Type)

	if token.Value != "" {
		fmt.Printf("  Value: %q\n", token.Value)
	}

	if token.Style != "" && token.Style != "Plain" {
		fmt.Printf("  Style: %s\n", token.Style)
	}

	if token.HeadComment != "" {
		fmt.Printf("  Head: %q\n", token.HeadComment)
	}
	if token.LineComment != "" {
		fmt.Printf("  Line: %q\n", token.LineComment)
	}
	if token.FootComment != "" {
		fmt.Printf("  Foot: %q\n", token.FootComment)
	}

	if profuse {
		if token.StartLine == token.EndLine && token.StartColumn == token.EndColumn {
			fmt.Printf("  Pos: {%d: %d}\n", token.StartLine, token.StartColumn)
		} else {
			fmt.Printf("  Pos: {%d: %d, %d: %d}\n", token.StartLine, token.StartColumn, token.EndLine, token.EndColumn)
		}
	}
	fmt.Println()
}

// processNodeToTokens converts a node to a slice of tokens
func processNodeToTokens(node *yaml.Node, profuse bool) []*Token {
	var tokens []*Token

	// Add stream start token
	tokens = append(tokens, &Token{
		Type: "STREAM-START",
	})

	// Add document start token
	tokens = append(tokens, &Token{
		Type: "DOCUMENT-START",
	})

	// Process the node content
	tokens = append(tokens, processNodeToTokensRecursive(node, profuse)...)

	// Add document end token
	tokens = append(tokens, &Token{
		Type: "DOCUMENT-END",
	})

	// Add stream end token
	tokens = append(tokens, &Token{
		Type: "STREAM-END",
	})

	return tokens
}

// processNodeToTokensRecursive recursively converts a node to tokens
func processNodeToTokensRecursive(node *yaml.Node, profuse bool) []*Token {
	var tokens []*Token

	switch node.Kind {
	case yaml.DocumentNode:
		for _, child := range node.Content {
			tokens = append(tokens, processNodeToTokensRecursive(child, profuse)...)
		}
	case yaml.MappingNode:
		tokens = append(tokens, &Token{
			Type:        "BLOCK-MAPPING-START",
			StartLine:   node.Line,
			StartColumn: node.Column,
			EndLine:     node.Line,
			EndColumn:   node.Column,
			HeadComment: node.HeadComment,
			LineComment: node.LineComment,
			FootComment: node.FootComment,
		})
		for i := 0; i < len(node.Content); i += 2 {
			if i+1 < len(node.Content) {
				// Key
				tokens = append(tokens, &Token{
					Type:        "KEY",
					StartLine:   node.Content[i].Line,
					StartColumn: node.Content[i].Column,
					EndLine:     node.Content[i].Line,
					EndColumn:   node.Content[i].Column,
				})
				keyTokens := processNodeToTokensRecursive(node.Content[i], profuse)
				tokens = append(tokens, keyTokens...)
				// Value
				tokens = append(tokens, &Token{
					Type:        "VALUE",
					StartLine:   node.Content[i+1].Line,
					StartColumn: node.Content[i+1].Column,
					EndLine:     node.Content[i+1].Line,
					EndColumn:   node.Content[i+1].Column,
				})
				valueTokens := processNodeToTokensRecursive(node.Content[i+1], profuse)
				tokens = append(tokens, valueTokens...)
			}
		}
		tokens = append(tokens, &Token{
			Type:        "BLOCK-END",
			StartLine:   node.Line,
			StartColumn: node.Column,
			EndLine:     node.Line,
			EndColumn:   node.Column,
		})
	case yaml.SequenceNode:
		tokens = append(tokens, &Token{
			Type:        "BLOCK-SEQUENCE-START",
			StartLine:   node.Line,
			StartColumn: node.Column,
			EndLine:     node.Line,
			EndColumn:   node.Column,
			HeadComment: node.HeadComment,
			LineComment: node.LineComment,
			FootComment: node.FootComment,
		})
		for _, child := range node.Content {
			tokens = append(tokens, &Token{
				Type:        "BLOCK-ENTRY",
				StartLine:   child.Line,
				StartColumn: child.Column,
				EndLine:     child.Line,
				EndColumn:   child.Column,
			})
			childTokens := processNodeToTokensRecursive(child, profuse)
			tokens = append(tokens, childTokens...)
		}
		tokens = append(tokens, &Token{
			Type:        "BLOCK-END",
			StartLine:   node.Line,
			StartColumn: node.Column,
			EndLine:     node.Line,
			EndColumn:   node.Column,
		})
	case yaml.ScalarNode:
		// Calculate end position for scalars based on value length
		endLine := node.Line
		endColumn := node.Column
		if node.Value != "" {
			// For single-line values, add the length to the column
			if !strings.Contains(node.Value, "\n") {
				endColumn += len(node.Value)
			} else {
				// For multi-line values, we'd need more complex logic
				// For now, just use the start position
				endColumn = node.Column
			}
		}
		tokens = append(tokens, &Token{
			Type:        "SCALAR",
			Value:       node.Value,
			StartLine:   node.Line,
			StartColumn: node.Column,
			EndLine:     endLine,
			EndColumn:   endColumn,
			Style:       formatStyle(node.Style),
			HeadComment: node.HeadComment,
			LineComment: node.LineComment,
			FootComment: node.FootComment,
		})
	}

	return tokens
}
