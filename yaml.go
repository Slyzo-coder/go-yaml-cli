// Package main provides YAML formatting utilities for the go-yaml tool.
package main

import (
	"fmt"
	"io"
	"os"

	"go.yaml.in/yaml/v3"
)

// ProcessYAML reads YAML from stdin and outputs formatted YAML
func ProcessYAML(preserve bool) error {
	if preserve {
		// Preserve comments and styles by using yaml.Node
		decoder := yaml.NewDecoder(os.Stdin)
		firstDoc := true

		for {
			var node yaml.Node
			err := decoder.Decode(&node)
			if err != nil {
				if err == io.EOF || err.Error() == "EOF" {
					break
				}
				return fmt.Errorf("failed to decode YAML: %v", err)
			}

			// Add document separator for all documents except the first
			if !firstDoc {
				fmt.Println("---")
			}
			firstDoc = false

			// If the node is not a DocumentNode, wrap it in one
			var outNode *yaml.Node
			if node.Kind == yaml.DocumentNode {
				outNode = &node
			} else {
				outNode = &yaml.Node{
					Kind:    yaml.DocumentNode,
					Content: []*yaml.Node{&node},
				}
			}

			encoder := yaml.NewEncoder(os.Stdout)
			encoder.SetIndent(2)
			if err := encoder.Encode(outNode); err != nil {
				encoder.Close()
				return fmt.Errorf("failed to encode YAML: %v", err)
			}
			encoder.Close()
		}
	} else {
		// Don't preserve comments and styles - use interface{} for clean output
		decoder := yaml.NewDecoder(os.Stdin)
		firstDoc := true

		for {
			var data interface{}
			err := decoder.Decode(&data)
			if err != nil {
				if err == io.EOF || err.Error() == "EOF" {
					break
				}
				return fmt.Errorf("failed to decode YAML: %v", err)
			}

			// Add document separator for all documents except the first
			if !firstDoc {
				fmt.Println("---")
			}
			firstDoc = false

			encoder := yaml.NewEncoder(os.Stdout)
			encoder.SetIndent(2)
			if err := encoder.Encode(data); err != nil {
				encoder.Close()
				return fmt.Errorf("failed to encode YAML: %v", err)
			}
			encoder.Close()
		}
	}

	return nil
}
