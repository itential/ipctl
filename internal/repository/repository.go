// Copyright 2024 Itential Inc. All Rights Reserved
// Unauthorized copying of this file, via any medium is strictly prohibited
// Proprietary and confidential

package repository

import (
	"fmt"

	"github.com/mitchellh/go-homedir"
)

// Repository represents configuration for a Git repository containing Itential assets.
// It includes the repository URL, authentication credentials, and the reference (branch/tag) to use.
type Repository struct {
	Url            string
	PrivateKey     string
	PrivateKeyFile string
	Reference      string
}

// Default returns a Repository with default (empty) values.
func Default() *Repository {
	return &Repository{}
}

// Manager manages a collection of named repositories.
type Manager struct {
	repositories map[string]*Repository
}

// NewManager creates a new repository manager with an empty collection.
func NewManager() *Manager {
	return &Manager{
		repositories: make(map[string]*Repository),
	}
}

// Add adds a repository to the manager with the given name.
func (m *Manager) Add(name string, repo *Repository) {
	m.repositories[name] = repo
}

// Get retrieves a repository by name.
// Returns an error if the named repository doesn't exist.
func (m *Manager) Get(name string) (*Repository, error) {
	if r, exists := m.repositories[name]; exists {
		return r, nil
	}
	return nil, fmt.Errorf("repository %q does not exist", name)
}

// List returns the names of all registered repositories.
func (m *Manager) List() []string {
	names := make([]string, 0, len(m.repositories))
	for name := range m.repositories {
		names = append(names, name)
	}
	return names
}

// Loader loads repository configuration from multiple sources with precedence.
// Values are loaded with the following precedence (highest to lowest):
// 1. Overrides (typically from environment variables or CLI flags)
// 2. Values (from configuration file)
type Loader struct {
	values    map[string]interface{}
	overrides map[string]interface{}
}

// NewLoader creates a new repository loader with the specified value sources.
func NewLoader(values, overrides map[string]interface{}) *Loader {
	return &Loader{
		values:    values,
		overrides: overrides,
	}
}

// Load constructs a Repository by merging values from values and overrides.
// Overrides take precedence over values.
// The PrivateKeyFile field supports home directory expansion (~/file becomes /home/user/file).
func (l *Loader) Load() *Repository {
	r := &Repository{}

	// Helper to get value with precedence: overrides -> values
	getValue := func(key string) interface{} {
		if val, ok := l.overrides[key]; ok {
			return val
		}
		return l.values[key]
	}

	r.Url = getString(getValue("url"), "")
	r.PrivateKey = getString(getValue("private_key"), "")

	// Expand home directory for private key file
	keyFile := getString(getValue("private_key_file"), "")
	if keyFile != "" {
		if expanded, err := homedir.Expand(keyFile); err == nil {
			keyFile = expanded
		}
	}
	r.PrivateKeyFile = keyFile

	r.Reference = getString(getValue("reference"), "")

	return r
}

// getString converts an interface value to a string with a fallback default.
func getString(val interface{}, defaultVal string) string {
	if val == nil {
		return defaultVal
	}
	if s, ok := val.(string); ok {
		return s
	}
	return defaultVal
}
