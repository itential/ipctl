// Copyright 2024 Itential Inc. All Rights Reserved
// Unauthorized copying of this file, via any medium is strictly prohibited
// Proprietary and confidential

package terminal

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestDefaultConfig(t *testing.T) {
	cfg := DefaultConfig()

	assert.False(t, cfg.NoColor)
	assert.Equal(t, "human", cfg.DefaultOutput)
	assert.True(t, cfg.Pager)
}

func TestLoadFromEnv(t *testing.T) {
	tests := []struct {
		name     string
		env      map[string]string
		expected Config
	}{
		{
			name: "all defaults",
			env:  map[string]string{},
			expected: Config{
				NoColor:       false,
				DefaultOutput: "human",
				Pager:         true,
			},
		},
		{
			name: "disable color",
			env: map[string]string{
				"IPCTL_TERMINAL_NO_COLOR": "true",
			},
			expected: Config{
				NoColor:       true,
				DefaultOutput: "human",
				Pager:         true,
			},
		},
		{
			name: "set default output to json",
			env: map[string]string{
				"IPCTL_TERMINAL_DEFAULT_OUTPUT": "json",
			},
			expected: Config{
				NoColor:       false,
				DefaultOutput: "json",
				Pager:         true,
			},
		},
		{
			name: "disable pager",
			env: map[string]string{
				"IPCTL_TERMINAL_PAGER": "false",
			},
			expected: Config{
				NoColor:       false,
				DefaultOutput: "human",
				Pager:         false,
			},
		},
		{
			name: "all settings custom",
			env: map[string]string{
				"IPCTL_TERMINAL_NO_COLOR":       "true",
				"IPCTL_TERMINAL_DEFAULT_OUTPUT": "yaml",
				"IPCTL_TERMINAL_PAGER":          "false",
			},
			expected: Config{
				NoColor:       true,
				DefaultOutput: "yaml",
				Pager:         false,
			},
		},
		{
			name: "no color set to false keeps default",
			env: map[string]string{
				"IPCTL_TERMINAL_NO_COLOR": "false",
			},
			expected: Config{
				NoColor:       false,
				DefaultOutput: "human",
				Pager:         true,
			},
		},
		{
			name: "pager set to true keeps enabled",
			env: map[string]string{
				"IPCTL_TERMINAL_PAGER": "true",
			},
			expected: Config{
				NoColor:       false,
				DefaultOutput: "human",
				Pager:         true,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Set environment variables
			for k, v := range tt.env {
				t.Setenv(k, v)
			}

			cfg := LoadFromEnv()
			assert.Equal(t, tt.expected, cfg)
		})
	}
}
