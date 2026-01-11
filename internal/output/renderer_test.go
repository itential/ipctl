// Copyright 2024 Itential Inc. All Rights Reserved
// Unauthorized copying of this file, via any medium is strictly prohibited
// Proprietary and confidential

package output

import (
	"testing"

	"github.com/itential/ipctl/internal/runners"
)

// TestNewRenderer tests the creation of renderers with various configurations.
func TestNewRenderer(t *testing.T) {
	tests := []struct {
		name      string
		format    string
		usePager  bool
		wantError bool
		validate  func(t *testing.T, r *Renderer)
	}{
		{
			name:      "create json renderer",
			format:    "json",
			usePager:  false,
			wantError: false,
			validate: func(t *testing.T, r *Renderer) {
				if _, ok := r.formatter.(*JSONFormatter); !ok {
					t.Errorf("expected JSONFormatter, got %T", r.formatter)
				}
			},
		},
		{
			name:      "create yaml renderer",
			format:    "yaml",
			usePager:  false,
			wantError: false,
			validate: func(t *testing.T, r *Renderer) {
				if _, ok := r.formatter.(*YAMLFormatter); !ok {
					t.Errorf("expected YAMLFormatter, got %T", r.formatter)
				}
			},
		},
		{
			name:      "create human renderer with pager",
			format:    "human",
			usePager:  true,
			wantError: false,
			validate: func(t *testing.T, r *Renderer) {
				hf, ok := r.formatter.(*HumanFormatter)
				if !ok {
					t.Errorf("expected HumanFormatter, got %T", r.formatter)
					return
				}
				if !hf.UsePager() {
					t.Error("expected pager to be enabled")
				}
			},
		},
		{
			name:      "create human renderer without pager",
			format:    "human",
			usePager:  false,
			wantError: false,
			validate: func(t *testing.T, r *Renderer) {
				hf, ok := r.formatter.(*HumanFormatter)
				if !ok {
					t.Errorf("expected HumanFormatter, got %T", r.formatter)
					return
				}
				if hf.UsePager() {
					t.Error("expected pager to be disabled")
				}
			},
		},
		{
			name:      "unsupported format",
			format:    "xml",
			usePager:  false,
			wantError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			renderer, err := NewRenderer(tt.format, tt.usePager)

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

			if renderer == nil {
				t.Error("expected non-nil renderer")
				return
			}

			if tt.validate != nil {
				tt.validate(t, renderer)
			}
		})
	}
}

// TestRenderer_Render_JSON tests rendering responses in JSON format.
func TestRenderer_Render_JSON(t *testing.T) {
	tests := []struct {
		name      string
		response  *runners.Response
		wantError bool
	}{
		{
			name: "render simple object",
			response: &runners.Response{
				Object: map[string]string{
					"name":   "test",
					"status": "active",
				},
			},
			wantError: false,
		},
		{
			name: "render struct",
			response: &runners.Response{
				Object: struct {
					ID   int    `json:"id"`
					Name string `json:"name"`
				}{
					ID:   42,
					Name: "example",
				},
			},
			wantError: false,
		},
		{
			name: "render nil object",
			response: &runners.Response{
				Object: nil,
			},
			wantError: true,
		},
		{
			name:      "render nil response",
			response:  nil,
			wantError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			renderer, err := NewRenderer("json", false)
			if err != nil {
				t.Fatalf("failed to create renderer: %v", err)
			}

			err = renderer.Render(tt.response)

			if tt.wantError {
				if err == nil {
					t.Error("expected error but got none")
				}
				return
			}

			if err != nil {
				t.Errorf("unexpected error: %v", err)
			}
		})
	}
}

// TestRenderer_Render_YAML tests rendering responses in YAML format.
func TestRenderer_Render_YAML(t *testing.T) {
	tests := []struct {
		name      string
		response  *runners.Response
		wantError bool
	}{
		{
			name: "render simple object",
			response: &runners.Response{
				Object: map[string]string{
					"name":   "test",
					"status": "active",
				},
			},
			wantError: false,
		},
		{
			name: "render struct",
			response: &runners.Response{
				Object: struct {
					ID   int    `yaml:"id"`
					Name string `yaml:"name"`
				}{
					ID:   42,
					Name: "example",
				},
			},
			wantError: false,
		},
		{
			name: "render nil object",
			response: &runners.Response{
				Object: nil,
			},
			wantError: true,
		},
		{
			name:      "render nil response",
			response:  nil,
			wantError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			renderer, err := NewRenderer("yaml", false)
			if err != nil {
				t.Fatalf("failed to create renderer: %v", err)
			}

			err = renderer.Render(tt.response)

			if tt.wantError {
				if err == nil {
					t.Error("expected error but got none")
				}
				return
			}

			if err != nil {
				t.Errorf("unexpected error: %v", err)
			}
		})
	}
}

// TestRenderer_Render_Human tests rendering responses in human format.
func TestRenderer_Render_Human(t *testing.T) {
	tests := []struct {
		name      string
		response  *runners.Response
		wantError bool
	}{
		{
			name: "render text",
			response: &runners.Response{
				Text: "Hello, World!",
			},
			wantError: false,
		},
		{
			name: "render tabular data",
			response: &runners.Response{
				Object: []map[string]any{
					{"name": "item1", "status": "active"},
					{"name": "item2", "status": "inactive"},
				},
				Keys: []string{"name", "status"},
			},
			wantError: false,
		},
		{
			name: "render with template",
			response: &runners.Response{
				Object: map[string]string{
					"name": "test",
				},
				Template: "Name: {{.name}}",
			},
			wantError: false,
		},
		{
			name: "render empty response",
			response: &runners.Response{
				Text:   "",
				Object: nil,
			},
			wantError: true,
		},
		{
			name:      "render nil response",
			response:  nil,
			wantError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			renderer, err := NewRenderer("human", false)
			if err != nil {
				t.Fatalf("failed to create renderer: %v", err)
			}

			err = renderer.Render(tt.response)

			if tt.wantError {
				if err == nil {
					t.Error("expected error but got none")
				}
				return
			}

			if err != nil {
				t.Errorf("unexpected error: %v", err)
			}
		})
	}
}

// TestRenderer_Render_HumanWithPager tests human rendering with pager enabled.
func TestRenderer_Render_HumanWithPager(t *testing.T) {
	tests := []struct {
		name      string
		response  *runners.Response
		wantError bool
	}{
		{
			name: "render tabular data with pager",
			response: &runners.Response{
				Object: []map[string]any{
					{"name": "item1", "status": "active"},
					{"name": "item2", "status": "inactive"},
				},
				Keys: []string{"name", "status"},
			},
			wantError: false,
		},
		{
			name: "render text with pager",
			response: &runners.Response{
				Text: "Long output that might need paging...",
			},
			wantError: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			renderer, err := NewRenderer("human", true)
			if err != nil {
				t.Fatalf("failed to create renderer: %v", err)
			}

			err = renderer.Render(tt.response)

			if tt.wantError {
				if err == nil {
					t.Error("expected error but got none")
				}
				return
			}

			if err != nil {
				t.Errorf("unexpected error: %v", err)
			}
		})
	}
}
