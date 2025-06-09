// Package main provides YAML to JSON conversion utilities for the go-yaml tool.
package main

import (
	"encoding/json"
	"fmt"
	"os"

	"go.yaml.in/yaml/v3"
)

// ProcessJSON reads YAML from stdin and outputs JSON encoding
func ProcessJSON(pretty bool) error {
	decoder := yaml.NewDecoder(os.Stdin)

	for {
		// Read each document
		var data interface{}
		err := decoder.Decode(&data)
		if err != nil {
			if err.Error() == "EOF" {
				break
			}
			return fmt.Errorf("failed to decode YAML: %v", err)
		}

		// Encode as JSON
		encoder := json.NewEncoder(os.Stdout)
		if pretty {
			encoder.SetIndent("", "  ")
		}
		if err := encoder.Encode(data); err != nil {
			return fmt.Errorf("failed to encode JSON: %v", err)
		}
	}

	return nil
}
