// Copyright 2024 Itential Inc. All Rights Reserved
// Unauthorized copying of this file, via any medium is strictly prohibited
// Proprietary and confidential

package output

import (
	"encoding/json"
	"strings"
	"testing"

	"gopkg.in/yaml.v2"
)

// TestNewJSONFormatter tests the creation of a JSON formatter.
func TestNewJSONFormatter(t *testing.T) {
	formatter := NewJSONFormatter()
	if formatter == nil {
		t.Error("NewJSONFormatter should return a non-nil formatter")
	}
}

// TestJSONFormatter_Format tests JSON formatting with various inputs.
func TestJSONFormatter_Format(t *testing.T) {
	tests := []struct {
		name      string
		input     any
		wantError bool
		validate  func(t *testing.T, output string)
	}{
		{
			name: "format simple map",
			input: map[string]string{
				"name":   "test",
				"status": "active",
			},
			wantError: false,
			validate: func(t *testing.T, output string) {
				var result map[string]string
				if err := json.Unmarshal([]byte(output), &result); err != nil {
					t.Errorf("output is not valid JSON: %v", err)
				}
				if result["name"] != "test" {
					t.Errorf("expected name=test, got %s", result["name"])
				}
			},
		},
		{
			name: "format struct",
			input: struct {
				ID   int    `json:"id"`
				Name string `json:"name"`
			}{
				ID:   42,
				Name: "example",
			},
			wantError: false,
			validate: func(t *testing.T, output string) {
				if !strings.Contains(output, `"id"`) {
					t.Error("output should contain id field")
				}
				if !strings.Contains(output, `"name"`) {
					t.Error("output should contain name field")
				}
			},
		},
		{
			name:      "format nil data",
			input:     nil,
			wantError: true,
		},
		{
			name: "format slice",
			input: []string{
				"item1",
				"item2",
			},
			wantError: false,
			validate: func(t *testing.T, output string) {
				var result []string
				if err := json.Unmarshal([]byte(output), &result); err != nil {
					t.Errorf("output is not valid JSON: %v", err)
				}
				if len(result) != 2 {
					t.Errorf("expected 2 items, got %d", len(result))
				}
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			formatter := NewJSONFormatter()
			output, err := formatter.Format(tt.input)

			if tt.wantError {
				if err == nil {
					t.Error("expected error but got none")
				}
				return
			}

			if err != nil {
				t.Errorf("unexpected error: %v", err)
				return
			}

			if tt.validate != nil {
				tt.validate(t, output)
			}
		})
	}
}

// TestNewYAMLFormatter tests the creation of a YAML formatter.
func TestNewYAMLFormatter(t *testing.T) {
	formatter := NewYAMLFormatter()
	if formatter == nil {
		t.Error("NewYAMLFormatter should return a non-nil formatter")
	}
}

// TestYAMLFormatter_Format tests YAML formatting with various inputs.
func TestYAMLFormatter_Format(t *testing.T) {
	tests := []struct {
		name      string
		input     any
		wantError bool
		validate  func(t *testing.T, output string)
	}{
		{
			name: "format simple map",
			input: map[string]string{
				"name":   "test",
				"status": "active",
			},
			wantError: false,
			validate: func(t *testing.T, output string) {
				var result map[string]string
				if err := yaml.Unmarshal([]byte(output), &result); err != nil {
					t.Errorf("output is not valid YAML: %v", err)
				}
				if result["name"] != "test" {
					t.Errorf("expected name=test, got %s", result["name"])
				}
			},
		},
		{
			name: "format struct",
			input: struct {
				ID   int    `yaml:"id"`
				Name string `yaml:"name"`
			}{
				ID:   42,
				Name: "example",
			},
			wantError: false,
			validate: func(t *testing.T, output string) {
				if !strings.Contains(output, "id:") {
					t.Error("output should contain id field")
				}
				if !strings.Contains(output, "name:") {
					t.Error("output should contain name field")
				}
			},
		},
		{
			name:      "format nil data",
			input:     nil,
			wantError: true,
		},
		{
			name: "format slice",
			input: []map[string]string{
				{"name": "item1"},
				{"name": "item2"},
			},
			wantError: false,
			validate: func(t *testing.T, output string) {
				var result []map[string]string
				if err := yaml.Unmarshal([]byte(output), &result); err != nil {
					t.Errorf("output is not valid YAML: %v", err)
				}
				if len(result) != 2 {
					t.Errorf("expected 2 items, got %d", len(result))
				}
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			formatter := NewYAMLFormatter()
			output, err := formatter.Format(tt.input)

			if tt.wantError {
				if err == nil {
					t.Error("expected error but got none")
				}
				return
			}

			if err != nil {
				t.Errorf("unexpected error: %v", err)
				return
			}

			if tt.validate != nil {
				tt.validate(t, output)
			}
		})
	}
}

// TestNewHumanFormatter tests the creation of a human formatter.
func TestNewHumanFormatter(t *testing.T) {
	tests := []struct {
		name     string
		usePager bool
	}{
		{
			name:     "with pager enabled",
			usePager: true,
		},
		{
			name:     "with pager disabled",
			usePager: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			formatter := NewHumanFormatter(tt.usePager)
			if formatter == nil {
				t.Error("NewHumanFormatter should return a non-nil formatter")
			}
			if formatter.UsePager() != tt.usePager {
				t.Errorf("expected UsePager()=%v, got %v", tt.usePager, formatter.UsePager())
			}
		})
	}
}

// TestHumanFormatter_Format tests human formatting with various inputs.
func TestHumanFormatter_Format(t *testing.T) {
	tests := []struct {
		name      string
		input     any
		wantError bool
		validate  func(t *testing.T, output string)
	}{
		{
			name:      "format string",
			input:     "Hello, World!",
			wantError: false,
			validate: func(t *testing.T, output string) {
				if output != "Hello, World!" {
					t.Errorf("expected 'Hello, World!', got '%s'", output)
				}
			},
		},
		{
			name:      "format empty string",
			input:     "",
			wantError: true,
		},
		{
			name:      "format nil data",
			input:     nil,
			wantError: true,
		},
		{
			name:      "format non-string data",
			input:     map[string]string{"key": "value"},
			wantError: true,
		},
		{
			name: "format multiline string",
			input: `NAME	STATUS
item1	active
item2	inactive`,
			wantError: false,
			validate: func(t *testing.T, output string) {
				if !strings.Contains(output, "NAME") {
					t.Error("output should contain NAME header")
				}
				if !strings.Contains(output, "item1") {
					t.Error("output should contain item1")
				}
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			formatter := NewHumanFormatter(false)
			output, err := formatter.Format(tt.input)

			if tt.wantError {
				if err == nil {
					t.Error("expected error but got none")
				}
				return
			}

			if err != nil {
				t.Errorf("unexpected error: %v", err)
				return
			}

			if tt.validate != nil {
				tt.validate(t, output)
			}
		})
	}
}

// TestNewFormatter tests the formatter factory function.
func TestNewFormatter(t *testing.T) {
	tests := []struct {
		name         string
		format       string
		usePager     bool
		wantType     string
		wantError    bool
		validateType func(t *testing.T, f Formatter)
	}{
		{
			name:      "create json formatter",
			format:    "json",
			usePager:  false,
			wantType:  "*output.JSONFormatter",
			wantError: false,
			validateType: func(t *testing.T, f Formatter) {
				if _, ok := f.(*JSONFormatter); !ok {
					t.Errorf("expected JSONFormatter, got %T", f)
				}
			},
		},
		{
			name:      "create yaml formatter",
			format:    "yaml",
			usePager:  false,
			wantType:  "*output.YAMLFormatter",
			wantError: false,
			validateType: func(t *testing.T, f Formatter) {
				if _, ok := f.(*YAMLFormatter); !ok {
					t.Errorf("expected YAMLFormatter, got %T", f)
				}
			},
		},
		{
			name:      "create human formatter without pager",
			format:    "human",
			usePager:  false,
			wantType:  "*output.HumanFormatter",
			wantError: false,
			validateType: func(t *testing.T, f Formatter) {
				hf, ok := f.(*HumanFormatter)
				if !ok {
					t.Errorf("expected HumanFormatter, got %T", f)
					return
				}
				if hf.UsePager() {
					t.Error("expected pager to be disabled")
				}
			},
		},
		{
			name:      "create human formatter with pager",
			format:    "human",
			usePager:  true,
			wantType:  "*output.HumanFormatter",
			wantError: false,
			validateType: func(t *testing.T, f Formatter) {
				hf, ok := f.(*HumanFormatter)
				if !ok {
					t.Errorf("expected HumanFormatter, got %T", f)
					return
				}
				if !hf.UsePager() {
					t.Error("expected pager to be enabled")
				}
			},
		},
		{
			name:      "uppercase format name",
			format:    "JSON",
			usePager:  false,
			wantError: false,
			validateType: func(t *testing.T, f Formatter) {
				if _, ok := f.(*JSONFormatter); !ok {
					t.Errorf("expected JSONFormatter, got %T", f)
				}
			},
		},
		{
			name:      "mixed case format name",
			format:    "YaMl",
			usePager:  false,
			wantError: false,
			validateType: func(t *testing.T, f Formatter) {
				if _, ok := f.(*YAMLFormatter); !ok {
					t.Errorf("expected YAMLFormatter, got %T", f)
				}
			},
		},
		{
			name:      "unsupported format",
			format:    "xml",
			usePager:  false,
			wantError: true,
		},
		{
			name:      "empty format",
			format:    "",
			usePager:  false,
			wantError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			formatter, err := NewFormatter(tt.format, tt.usePager)

			if tt.wantError {
				if err == nil {
					t.Error("expected error but got none")
				}
				return
			}

			if err != nil {
				t.Errorf("unexpected error: %v", err)
				return
			}

			if formatter == nil {
				t.Error("expected non-nil formatter")
				return
			}

			if tt.validateType != nil {
				tt.validateType(t, formatter)
			}
		})
	}
}
