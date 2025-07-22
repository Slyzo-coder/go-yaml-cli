// Package main provides YAML event formatting utilities for the go-yaml tool.
package main

import (
	"fmt"
	"io"
	"os"
	"strings"

	"go.yaml.in/yaml/v3"
)

// EventType represents the type of a YAML event
type EventType string

const (
	EventDocumentStart EventType = "DOCUMENT-START"
	EventDocumentEnd   EventType = "DOCUMENT-END"
	EventScalar        EventType = "SCALAR"
	EventSequenceStart EventType = "SEQUENCE-START"
	EventSequenceEnd   EventType = "SEQUENCE-END"
	EventMappingStart  EventType = "MAPPING-START"
	EventMappingEnd    EventType = "MAPPING-END"
)

// Event represents a YAML event
type Event struct {
	Type        EventType
	Value       string
	Anchor      string
	Tag         string
	Style       string
	Implicit    bool
	StartLine   int
	StartColumn int
	EndLine     int
	EndColumn   int
	HeadComment string
	LineComment string
	FootComment string
}

// ProcessEvents reads YAML from stdin and outputs event information
func ProcessEvents(profuse, compact bool) error {
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

			events := processNodeToEvents(&node, profuse)
			for _, event := range events {
				if !first {
					fmt.Println()
				}
				first = false

				fmt.Print("- ")
				printEventCompact(event, profuse)
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

			events := processNodeToEvents(&node, profuse)
			for _, event := range events {
				printEvent(event, profuse)
			}
		}
	}

	return nil
}

// processNode recursively processes a YAML node and generates events
func processNode(node *yaml.Node, depth int, profuse bool) error {
	// Generate start event for this node
	event := createEvent(node, depth)
	if event != nil {
		printEvent(event, profuse)
	}

	// Process children if this node has content
	if node.Content != nil {
		for _, child := range node.Content {
			if err := processNode(child, depth+1, profuse); err != nil {
				return err
			}
		}
	}

	// Generate end event for sequence and mapping nodes
	if node.Kind == yaml.SequenceNode || node.Kind == yaml.MappingNode {
		endEvent := &Event{
			Type: getEndEventType(node.Kind),
		}
		printEvent(endEvent, profuse)
	}

	return nil
}

// createEvent creates an event from a YAML node
func createEvent(node *yaml.Node, depth int) *Event {
	event := &Event{
		StartLine:   node.Line,
		StartColumn: node.Column,
		EndLine:     node.Line,
		EndColumn:   node.Column,
		HeadComment: node.HeadComment,
		LineComment: node.LineComment,
		FootComment: node.FootComment,
	}

	switch node.Kind {
	case yaml.DocumentNode:
		event.Type = EventDocumentStart
		event.Implicit = true
	case yaml.ScalarNode:
		event.Type = EventScalar
		event.Value = node.Value
		event.Anchor = node.Anchor
		event.Tag = formatTag(node.Tag)
		event.Style = formatStyle(node.Style)
		event.Implicit = true
	case yaml.SequenceNode:
		event.Type = EventSequenceStart
		event.Anchor = node.Anchor
		event.Tag = formatTag(node.Tag)
		event.Style = formatStyle(node.Style)
		event.Implicit = true
	case yaml.MappingNode:
		event.Type = EventMappingStart
		event.Anchor = node.Anchor
		event.Tag = formatTag(node.Tag)
		event.Style = formatStyle(node.Style)
		event.Implicit = true
	case yaml.AliasNode:
		event.Type = EventScalar
		event.Anchor = node.Anchor
		event.Value = "*" + node.Anchor
	default:
		return nil
	}

	return event
}

// getEndEventType returns the end event type for a node kind
func getEndEventType(kind yaml.Kind) EventType {
	switch kind {
	case yaml.SequenceNode:
		return EventSequenceEnd
	case yaml.MappingNode:
		return EventMappingEnd
	default:
		return EventDocumentEnd
	}
}

// printEventCompact prints an event in compact flow style format
func printEventCompact(event *Event, profuse bool) {
	fmt.Print("{Event: ", event.Type)
	if event.Value != "" {
		fmt.Printf(", Value: %q", event.Value)
	}
	if event.Style != "" {
		fmt.Printf(", Style: %s", event.Style)
	}
	if event.HeadComment != "" {
		fmt.Printf(", Head: %q", event.HeadComment)
	}
	if event.LineComment != "" {
		fmt.Printf(", Line: %q", event.LineComment)
	}
	if event.FootComment != "" {
		fmt.Printf(", Foot: %q", event.FootComment)
	}
	if profuse {
		if event.StartLine == event.EndLine && event.StartColumn == event.EndColumn {
			fmt.Printf(", Pos: {%d: %d}", event.StartLine, event.StartColumn)
		} else {
			fmt.Printf(", Pos: {%d: %d, %d: %d}", event.StartLine, event.StartColumn, event.EndLine, event.EndColumn)
		}
	}
	fmt.Print("}")
}

// processNodeToEvents converts a node to a slice of events for compact output
func processNodeToEvents(node *yaml.Node, profuse bool) []*Event {
	var events []*Event

	// Add document start event
	events = append(events, &Event{
		Type:        "DOCUMENT-START",
		StartLine:   node.Line,
		StartColumn: node.Column,
		EndLine:     node.Line,
		EndColumn:   node.Column,
	})

	// Process the node content
	events = append(events, processNodeToEventsRecursive(node, profuse)...)

	// Add document end event
	events = append(events, &Event{
		Type:        "DOCUMENT-END",
		StartLine:   node.Line,
		StartColumn: node.Column,
		EndLine:     node.Line,
		EndColumn:   node.Column,
	})

	return events
}

// processNodeToEventsRecursive recursively converts a node to events
func processNodeToEventsRecursive(node *yaml.Node, profuse bool) []*Event {
	var events []*Event

	switch node.Kind {
	case yaml.DocumentNode:
		for _, child := range node.Content {
			events = append(events, processNodeToEventsRecursive(child, profuse)...)
		}
	case yaml.MappingNode:
		events = append(events, &Event{
			Type:        "MAPPING-START",
			StartLine:   node.Line,
			StartColumn: node.Column,
			EndLine:     node.Line,
			EndColumn:   node.Column,
			Style:       formatStyle(node.Style),
			HeadComment: node.HeadComment,
			LineComment: node.LineComment,
			FootComment: node.FootComment,
		})
		for i := 0; i < len(node.Content); i += 2 {
			if i+1 < len(node.Content) {
				// Key
				keyEvents := processNodeToEventsRecursive(node.Content[i], profuse)
				events = append(events, keyEvents...)
				// Value
				valueEvents := processNodeToEventsRecursive(node.Content[i+1], profuse)
				events = append(events, valueEvents...)
			}
		}
		events = append(events, &Event{
			Type:        "MAPPING-END",
			StartLine:   node.Line,
			StartColumn: node.Column,
			EndLine:     node.Line,
			EndColumn:   node.Column,
		})
	case yaml.SequenceNode:
		events = append(events, &Event{
			Type:        "SEQUENCE-START",
			StartLine:   node.Line,
			StartColumn: node.Column,
			EndLine:     node.Line,
			EndColumn:   node.Column,
			Style:       formatStyle(node.Style),
			HeadComment: node.HeadComment,
			LineComment: node.LineComment,
			FootComment: node.FootComment,
		})
		for _, child := range node.Content {
			childEvents := processNodeToEventsRecursive(child, profuse)
			events = append(events, childEvents...)
		}
		events = append(events, &Event{
			Type:        "SEQUENCE-END",
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
		events = append(events, &Event{
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

	return events
}

// printEvent prints an event in the expected format
func printEvent(event *Event, profuse bool) {
	fmt.Printf("- Event: %v\n", event.Type)

	switch event.Type {
	case EventScalar:
		if event.Value != "" {
			fmt.Printf("  Value: %q\n", event.Value)
		}
		if event.Style != "" {
			fmt.Printf("  Style: %s\n", event.Style)
		}
		if event.Tag != "" {
			fmt.Printf("  Tag: %s\n", event.Tag)
		}
		if event.Anchor != "" {
			fmt.Printf("  Anchor: %s\n", event.Anchor)
		}
	case EventSequenceStart, EventMappingStart:
		if event.Style != "" {
			fmt.Printf("  Style: %s\n", event.Style)
		}
		if event.Tag != "" {
			fmt.Printf("  Tag: %s\n", event.Tag)
		}
		if event.Anchor != "" {
			fmt.Printf("  Anchor: %s\n", event.Anchor)
		}
	}

	if event.HeadComment != "" {
		fmt.Printf("  Head: %q\n", event.HeadComment)
	}
	if event.LineComment != "" {
		fmt.Printf("  Line: %q\n", event.LineComment)
	}
	if event.FootComment != "" {
		fmt.Printf("  Foot: %q\n", event.FootComment)
	}

	if profuse {
		if event.StartLine == event.EndLine && event.StartColumn == event.EndColumn {
			fmt.Printf("  Pos: {%d: %d}\n", event.StartLine, event.StartColumn)
		} else {
			fmt.Printf("  Pos: {%d: %d, %d: %d}\n", event.StartLine, event.StartColumn, event.EndLine, event.EndColumn)
		}
	}
	fmt.Println()
}
