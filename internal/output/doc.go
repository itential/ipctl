// Copyright 2024 Itential Inc. All Rights Reserved
// Unauthorized copying of this file, via any medium is strictly prohibited
// Proprietary and confidential

// Package output provides output formatting and rendering capabilities for CLI commands.
//
// This package separates the concerns of output formatting from command execution,
// allowing handlers to focus on business logic while formatters handle presentation.
//
// The package defines a Formatter interface that can be implemented for different
// output formats (JSON, YAML, human-readable). A Renderer coordinates between
// formatters and the terminal display layer based on configuration.
//
// Example usage:
//
//	renderer, err := output.NewRenderer(config.TerminalDefaultOutput, config)
//	if err != nil {
//	    return err
//	}
//
//	resp := &runners.Response{
//	    Object: data,
//	    Keys:   []string{"name", "status"},
//	}
//
//	if err := renderer.Render(resp); err != nil {
//	    return err
//	}
//
// The package supports three output formats:
//
//  1. JSON: Structured JSON output with indentation
//  2. YAML: YAML-formatted output
//  3. Human: Human-readable tabular or templated output with optional pager
//
// Formatters are easily testable in isolation and can be extended to support
// additional output formats without modifying command handlers.
package output
