// Copyright 2024 Itential Inc. All Rights Reserved
// Unauthorized copying of this file, via any medium is strictly prohibited
// Proprietary and confidential

package output

import (
	"encoding/json"
	"fmt"
	"strings"

	"gopkg.in/yaml.v2"
)

// Formatter defines the interface for formatting output data.
//
// Implementations of this interface handle converting data into specific
// output formats like JSON, YAML, or human-readable text.
type Formatter interface {
	// Format converts the provided data into a formatted string.
	// Returns an error if the data cannot be formatted.
	Format(data any) (string, error)
}

// JSONFormatter formats data as indented JSON.
type JSONFormatter struct{}

// NewJSONFormatter creates a new JSON formatter.
func NewJSONFormatter() *JSONFormatter {
	return &JSONFormatter{}
}

// Format implements the Formatter interface for JSON output.
//
// The JSON is formatted with 4-space indentation for readability.
func (f *JSONFormatter) Format(data any) (string, error) {
	if data == nil {
		return "", fmt.Errorf("cannot format nil data as JSON")
	}

	b, err := json.MarshalIndent(data, "", "    ")
	if err != nil {
		return "", fmt.Errorf("failed to marshal data to JSON: %w", err)
	}

	return string(b), nil
}

// YAMLFormatter formats data as YAML.
type YAMLFormatter struct{}

// NewYAMLFormatter creates a new YAML formatter.
func NewYAMLFormatter() *YAMLFormatter {
	return &YAMLFormatter{}
}

// Format implements the Formatter interface for YAML output.
func (f *YAMLFormatter) Format(data any) (string, error) {
	if data == nil {
		return "", fmt.Errorf("cannot format nil data as YAML")
	}

	b, err := yaml.Marshal(data)
	if err != nil {
		return "", fmt.Errorf("failed to marshal data to YAML: %w", err)
	}

	return string(b), nil
}

// HumanFormatter formats data in a human-readable format.
//
// This formatter can handle both tabular data (when keys are provided)
// and text-based output (when using templates or plain text).
type HumanFormatter struct {
	// usePager indicates whether output should be displayed through a pager
	usePager bool
}

// NewHumanFormatter creates a new human-readable formatter.
//
// The usePager parameter controls whether the output will be piped
// through a pager for easier viewing of long output.
func NewHumanFormatter(usePager bool) *HumanFormatter {
	return &HumanFormatter{
		usePager: usePager,
	}
}

// Format implements the Formatter interface for human-readable output.
//
// For HumanFormatter, the data parameter should be a pre-formatted string
// (from Response.String()). This formatter primarily passes through the
// string and indicates whether paging should be used.
func (f *HumanFormatter) Format(data any) (string, error) {
	if data == nil {
		return "", fmt.Errorf("cannot format nil data for human output")
	}

	// For human output, we expect the data to already be formatted as a string
	if str, ok := data.(string); ok {
		if str == "" {
			return "", fmt.Errorf("cannot format empty string for human output")
		}
		return str, nil
	}

	return "", fmt.Errorf("human formatter expects string data, got %T", data)
}

// UsePager returns whether this formatter should use a pager for output.
func (f *HumanFormatter) UsePager() bool {
	return f.usePager
}

// NewFormatter creates a formatter based on the specified format string.
//
// Supported formats:
//   - "json": Returns a JSONFormatter
//   - "yaml": Returns a YAMLFormatter
//   - "human": Returns a HumanFormatter
//
// Returns an error if the format is not recognized.
func NewFormatter(format string, usePager bool) (Formatter, error) {
	switch strings.ToLower(format) {
	case "json":
		return NewJSONFormatter(), nil
	case "yaml":
		return NewYAMLFormatter(), nil
	case "human":
		return NewHumanFormatter(usePager), nil
	default:
		return nil, fmt.Errorf("unsupported output format: %s", format)
	}
}
