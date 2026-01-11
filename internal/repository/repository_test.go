// Copyright 2024 Itential Inc. All Rights Reserved
// Unauthorized copying of this file, via any medium is strictly prohibited
// Proprietary and confidential

package repository

import (
	"os"
	"path/filepath"
	"testing"
)

// TestDefault verifies that Default returns an empty repository.
func TestDefault(t *testing.T) {
	r := Default()

	if r.Url != "" {
		t.Errorf("expected Url to be empty, got %q", r.Url)
	}
	if r.PrivateKey != "" {
		t.Errorf("expected PrivateKey to be empty, got %q", r.PrivateKey)
	}
	if r.PrivateKeyFile != "" {
		t.Errorf("expected PrivateKeyFile to be empty, got %q", r.PrivateKeyFile)
	}
	if r.Reference != "" {
		t.Errorf("expected Reference to be empty, got %q", r.Reference)
	}
}

// TestNewManager verifies that NewManager creates a manager with expected initial state.
func TestNewManager(t *testing.T) {
	m := NewManager()

	if m.repositories == nil {
		t.Fatal("expected repositories map to be initialized")
	}
	if len(m.repositories) != 0 {
		t.Errorf("expected repositories map to be empty, got %d entries", len(m.repositories))
	}
}

// TestManagerAdd verifies that Add correctly adds repositories to the manager.
func TestManagerAdd(t *testing.T) {
	m := NewManager()

	// Add repositories
	r1 := &Repository{Url: "https://github.com/example/repo1.git"}
	r2 := &Repository{Url: "https://github.com/example/repo2.git"}

	m.Add("repo1", r1)
	m.Add("repo2", r2)

	if m.repositories["repo1"] != r1 {
		t.Error("expected repo1 to be added to repositories map")
	}
	if m.repositories["repo2"] != r2 {
		t.Error("expected repo2 to be added to repositories map")
	}
}

// TestManagerGet verifies that Get retrieves repositories correctly.
func TestManagerGet(t *testing.T) {
	m := NewManager()

	// Add test repository
	r1 := &Repository{Url: "https://github.com/example/repo.git"}
	m.Add("test", r1)

	// Test getting existing repository
	got, err := m.Get("test")
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}
	if got != r1 {
		t.Error("expected to get the test repository")
	}

	// Test getting non-existent repository
	_, err = m.Get("nonexistent")
	if err == nil {
		t.Error("expected error for non-existent repository")
	}
}

// TestManagerList verifies that List returns all repository names.
func TestManagerList(t *testing.T) {
	m := NewManager()

	// Test empty manager
	names := m.List()
	if len(names) != 0 {
		t.Errorf("expected empty list, got %d names", len(names))
	}

	// Add repositories
	m.Add("repo1", &Repository{})
	m.Add("repo2", &Repository{})
	m.Add("repo3", &Repository{})

	names = m.List()
	if len(names) != 3 {
		t.Errorf("expected 3 names, got %d", len(names))
	}

	// Verify all names are present
	nameMap := make(map[string]bool)
	for _, name := range names {
		nameMap[name] = true
	}

	if !nameMap["repo1"] || !nameMap["repo2"] || !nameMap["repo3"] {
		t.Errorf("expected all repository names to be in list, got %v", names)
	}
}

// TestNewLoader verifies that NewLoader creates a loader with correct initialization.
func TestNewLoader(t *testing.T) {
	values := map[string]interface{}{"url": "https://github.com/example/repo.git"}
	overrides := map[string]interface{}{"reference": "main"}

	l := NewLoader(values, overrides)

	if l.values == nil || l.overrides == nil {
		t.Error("expected all maps to be initialized")
	}
}

// TestLoaderLoad verifies that Load correctly merges values with proper precedence.
func TestLoaderLoad(t *testing.T) {
	tests := []struct {
		name      string
		values    map[string]interface{}
		overrides map[string]interface{}
		expected  *Repository
	}{
		{
			name:      "empty loader returns default",
			values:    map[string]interface{}{},
			overrides: map[string]interface{}{},
			expected: &Repository{
				Url:            "",
				PrivateKey:     "",
				PrivateKeyFile: "",
				Reference:      "",
			},
		},
		{
			name: "values are loaded",
			values: map[string]interface{}{
				"url":              "https://github.com/example/repo.git",
				"private_key":      "key123",
				"private_key_file": "/path/to/key",
				"reference":        "main",
			},
			overrides: map[string]interface{}{},
			expected: &Repository{
				Url:            "https://github.com/example/repo.git",
				PrivateKey:     "key123",
				PrivateKeyFile: "/path/to/key",
				Reference:      "main",
			},
		},
		{
			name: "overrides take precedence",
			values: map[string]interface{}{
				"url":       "https://github.com/example/repo.git",
				"reference": "main",
			},
			overrides: map[string]interface{}{
				"url":       "https://github.com/example/override.git",
				"reference": "develop",
			},
			expected: &Repository{
				Url:       "https://github.com/example/override.git",
				Reference: "develop",
			},
		},
		{
			name: "overrides partial fields",
			values: map[string]interface{}{
				"url":         "https://github.com/example/repo.git",
				"private_key": "originalkey",
				"reference":   "main",
			},
			overrides: map[string]interface{}{
				"reference": "feature-branch",
			},
			expected: &Repository{
				Url:        "https://github.com/example/repo.git",
				PrivateKey: "originalkey",
				Reference:  "feature-branch",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			l := NewLoader(tt.values, tt.overrides)
			got := l.Load()

			if got.Url != tt.expected.Url {
				t.Errorf("Url: expected %q, got %q", tt.expected.Url, got.Url)
			}
			if got.PrivateKey != tt.expected.PrivateKey {
				t.Errorf("PrivateKey: expected %q, got %q", tt.expected.PrivateKey, got.PrivateKey)
			}
			if got.PrivateKeyFile != tt.expected.PrivateKeyFile {
				t.Errorf("PrivateKeyFile: expected %q, got %q", tt.expected.PrivateKeyFile, got.PrivateKeyFile)
			}
			if got.Reference != tt.expected.Reference {
				t.Errorf("Reference: expected %q, got %q", tt.expected.Reference, got.Reference)
			}
		})
	}
}

// TestLoaderLoadHomeDirectoryExpansion verifies that home directory paths are expanded.
func TestLoaderLoadHomeDirectoryExpansion(t *testing.T) {
	home, err := os.UserHomeDir()
	if err != nil {
		t.Skip("could not determine home directory")
	}

	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "tilde is expanded",
			input:    "~/.ssh/id_rsa",
			expected: filepath.Join(home, ".ssh", "id_rsa"),
		},
		{
			name:     "absolute path is not changed",
			input:    "/absolute/path/to/key",
			expected: "/absolute/path/to/key",
		},
		{
			name:     "relative path is not changed",
			input:    "relative/path/to/key",
			expected: "relative/path/to/key",
		},
		{
			name:     "empty path is not changed",
			input:    "",
			expected: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			values := map[string]interface{}{
				"private_key_file": tt.input,
			}
			l := NewLoader(values, map[string]interface{}{})
			got := l.Load()

			if got.PrivateKeyFile != tt.expected {
				t.Errorf("expected %q, got %q", tt.expected, got.PrivateKeyFile)
			}
		})
	}
}

// TestGetStringTypeConversions verifies getString handles different types correctly.
func TestGetStringTypeConversions(t *testing.T) {
	tests := []struct {
		name     string
		value    interface{}
		fallback string
		expected string
	}{
		{"nil returns default", nil, "default", "default"},
		{"string returns value", "test", "default", "test"},
		{"int returns default", 123, "default", "default"},
		{"bool returns default", true, "default", "default"},
		{"empty string returns empty", "", "default", ""},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := getString(tt.value, tt.fallback)
			if got != tt.expected {
				t.Errorf("expected %q, got %q", tt.expected, got)
			}
		})
	}
}
