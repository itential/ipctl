// Copyright 2024 Itential Inc. All Rights Reserved
// Unauthorized copying of this file, via any medium is strictly prohibited
// Proprietary and confidential

package terminal

import "os"

// Config holds the configuration options for the terminal package.
// It defines terminal output formatting and behavior settings.
type Config struct {
	// NoColor disables colored output in the terminal
	NoColor bool
	// DefaultOutput specifies the default output format (human, json, yaml)
	DefaultOutput string
	// Pager enables paging for long output
	Pager bool
}

// DefaultConfig returns a Config with sensible defaults.
// Default output is human-readable with colors enabled and paging disabled.
func DefaultConfig() Config {
	return Config{
		NoColor:       false,
		DefaultOutput: "human",
		Pager:         false,
	}
}

// LoadFromEnv creates a Config by loading values from environment variables.
// Supported environment variables:
//   - IPCTL_TERMINAL_NO_COLOR: Disable colored output (true/false)
//   - IPCTL_TERMINAL_DEFAULT_OUTPUT: Default output format (human, json, yaml)
//   - IPCTL_TERMINAL_PAGER: Enable pager for long output (true/false)
//
// Returns a Config with defaults for any unset environment variables.
func LoadFromEnv() Config {
	cfg := DefaultConfig()

	if noColor := os.Getenv("IPCTL_TERMINAL_NO_COLOR"); noColor == "true" {
		cfg.NoColor = true
	}

	if defaultOutput := os.Getenv("IPCTL_TERMINAL_DEFAULT_OUTPUT"); defaultOutput != "" {
		cfg.DefaultOutput = defaultOutput
	}

	if pager := os.Getenv("IPCTL_TERMINAL_PAGER"); pager == "true" {
		cfg.Pager = true
	}

	return cfg
}
