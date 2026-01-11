// Copyright 2024 Itential Inc. All Rights Reserved
// Unauthorized copying of this file, via any medium is strictly prohibited
// Proprietary and confidential

package profile

import (
	"fmt"
	"strconv"
)

// Profile represents the configuration for connecting to an Itential Platform instance.
// It includes connection details, authentication credentials, and operational settings.
type Profile struct {
	Host         string
	Port         int
	UseTLS       bool
	Verify       bool
	Username     string
	Password     string
	ClientID     string
	ClientSecret string
	MongoUrl     string
	Timeout      int
}

// Default returns a Profile with default values.
// The default profile uses localhost with TLS enabled and verification enabled.
func Default() *Profile {
	return &Profile{
		Host:    "localhost",
		Port:    0,
		UseTLS:  true,
		Verify:  true,
		Timeout: 0,
	}
}

// Manager manages a collection of named profiles.
// It maintains a default profile and tracks which profile is currently active.
type Manager struct {
	profiles       map[string]*Profile
	defaultProfile *Profile
	activeProfile  string
}

// NewManager creates a new profile manager with an empty collection.
// The active profile is initialized to "default".
func NewManager() *Manager {
	return &Manager{
		profiles:       make(map[string]*Profile),
		defaultProfile: Default(),
		activeProfile:  "default",
	}
}

// Add adds a profile to the manager with the given name.
// If the name is "default", it replaces the default profile.
func (m *Manager) Add(name string, profile *Profile) {
	if name == "default" {
		m.defaultProfile = profile
	}
	m.profiles[name] = profile
}

// Get retrieves a profile by name.
// If name is empty or "default", returns the default profile.
// Returns an error if the named profile doesn't exist.
func (m *Manager) Get(name string) (*Profile, error) {
	if name == "" {
		name = "default"
	}

	if p, exists := m.profiles[name]; exists {
		return p, nil
	}

	if name == "default" {
		return m.defaultProfile, nil
	}

	return Default(), fmt.Errorf("profile %q not found, using defaults", name)
}

// SetActive sets the name of the currently active profile.
func (m *Manager) SetActive(name string) {
	m.activeProfile = name
}

// Active returns the currently active profile.
// Returns an error if the active profile doesn't exist.
func (m *Manager) Active() (*Profile, error) {
	return m.Get(m.activeProfile)
}

// List returns the names of all registered profiles.
func (m *Manager) List() []string {
	names := make([]string, 0, len(m.profiles))
	for name := range m.profiles {
		names = append(names, name)
	}
	return names
}

// Loader loads profile configuration from multiple sources with precedence.
// Values are loaded with the following precedence (highest to lowest):
// 1. Overrides (typically from environment variables or CLI flags)
// 2. Values (from configuration file)
// 3. Defaults (from default profile or built-in defaults)
type Loader struct {
	defaults  map[string]interface{}
	values    map[string]interface{}
	overrides map[string]interface{}
}

// NewLoader creates a new profile loader with the specified value sources.
func NewLoader(values, defaults, overrides map[string]interface{}) *Loader {
	return &Loader{
		defaults:  defaults,
		values:    values,
		overrides: overrides,
	}
}

// Load constructs a Profile by merging values from defaults, values, and overrides.
// Overrides take precedence over values, which take precedence over defaults.
func (l *Loader) Load() *Profile {
	p := &Profile{}

	// Helper to get value with precedence: overrides -> values -> defaults
	getValue := func(key string) interface{} {
		if val, ok := l.overrides[key]; ok {
			return val
		}
		if val, ok := l.values[key]; ok {
			return val
		}
		return l.defaults[key]
	}

	// Load each field with appropriate type conversion
	p.Host = getString(getValue("host"), "localhost")
	p.Port = getInt(getValue("port"), 0)
	p.UseTLS = getBool(getValue("use_tls"), true)
	p.Verify = getBool(getValue("verify"), true)
	p.Username = getString(getValue("username"), "")
	p.Password = getString(getValue("password"), "")
	p.ClientID = getString(getValue("client_id"), "")
	p.ClientSecret = getString(getValue("client_secret"), "")
	p.MongoUrl = getString(getValue("mongo_url"), "")
	p.Timeout = getInt(getValue("timeout"), 0)

	return p
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

// getInt converts an interface value to an int with a fallback default.
// Supports both int types and string representations of integers.
func getInt(val interface{}, defaultVal int) int {
	if val == nil {
		return defaultVal
	}
	switch v := val.(type) {
	case int:
		return v
	case int64:
		return int(v)
	case float64:
		return int(v)
	case string:
		if i, err := strconv.Atoi(v); err == nil {
			return i
		}
	}
	return defaultVal
}

// getBool converts an interface value to a bool with a fallback default.
// Supports both bool types and string representations ("true", "false").
func getBool(val interface{}, defaultVal bool) bool {
	if val == nil {
		return defaultVal
	}
	switch v := val.(type) {
	case bool:
		return v
	case string:
		return v == "true"
	}
	return defaultVal
}
