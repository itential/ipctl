// Copyright 2024 Itential Inc. All Rights Reserved
// Unauthorized copying of this file, via any medium is strictly prohibited
// Proprietary and confidential

package profile

import (
	"testing"
)

// TestDefault verifies that Default returns a profile with expected default values.
func TestDefault(t *testing.T) {
	p := Default()

	if p.Host != "localhost" {
		t.Errorf("expected Host to be 'localhost', got %q", p.Host)
	}
	if p.Port != 0 {
		t.Errorf("expected Port to be 0, got %d", p.Port)
	}
	if !p.UseTLS {
		t.Error("expected UseTLS to be true")
	}
	if !p.Verify {
		t.Error("expected Verify to be true")
	}
	if p.Username != "" {
		t.Errorf("expected Username to be empty, got %q", p.Username)
	}
	if p.Password != "" {
		t.Errorf("expected Password to be empty, got %q", p.Password)
	}
	if p.ClientID != "" {
		t.Errorf("expected ClientID to be empty, got %q", p.ClientID)
	}
	if p.ClientSecret != "" {
		t.Errorf("expected ClientSecret to be empty, got %q", p.ClientSecret)
	}
	if p.MongoUrl != "" {
		t.Errorf("expected MongoUrl to be empty, got %q", p.MongoUrl)
	}
	if p.Timeout != 0 {
		t.Errorf("expected Timeout to be 0, got %d", p.Timeout)
	}
}

// TestNewManager verifies that NewManager creates a manager with expected initial state.
func TestNewManager(t *testing.T) {
	m := NewManager()

	if m.profiles == nil {
		t.Fatal("expected profiles map to be initialized")
	}
	if m.defaultProfile == nil {
		t.Fatal("expected defaultProfile to be initialized")
	}
	if m.activeProfile != "default" {
		t.Errorf("expected activeProfile to be 'default', got %q", m.activeProfile)
	}
}

// TestManagerAdd verifies that Add correctly adds profiles to the manager.
func TestManagerAdd(t *testing.T) {
	m := NewManager()

	// Add a named profile
	p1 := &Profile{Host: "example.com", Port: 443}
	m.Add("test", p1)

	if m.profiles["test"] != p1 {
		t.Error("expected profile to be added to profiles map")
	}

	// Add a default profile
	p2 := &Profile{Host: "default.com", Port: 8080}
	m.Add("default", p2)

	if m.defaultProfile != p2 {
		t.Error("expected defaultProfile to be updated")
	}
	if m.profiles["default"] != p2 {
		t.Error("expected default profile to also be in profiles map")
	}
}

// TestManagerGet verifies that Get retrieves profiles correctly.
func TestManagerGet(t *testing.T) {
	m := NewManager()

	// Add test profiles
	p1 := &Profile{Host: "example.com", Port: 443}
	m.Add("test", p1)

	// Test getting named profile
	got, err := m.Get("test")
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}
	if got != p1 {
		t.Error("expected to get the test profile")
	}

	// Test getting default profile with empty name
	got, err = m.Get("")
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}
	if got == nil {
		t.Error("expected to get default profile")
	}

	// Test getting default profile explicitly
	got, err = m.Get("default")
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}
	if got == nil {
		t.Error("expected to get default profile")
	}

	// Test getting non-existent profile
	got, err = m.Get("nonexistent")
	if err == nil {
		t.Error("expected error for non-existent profile")
	}
	if got == nil {
		t.Error("expected to get default profile as fallback")
	}
}

// TestManagerSetActiveAndActive verifies SetActive and Active methods.
func TestManagerSetActiveAndActive(t *testing.T) {
	m := NewManager()

	// Add test profiles
	p1 := &Profile{Host: "test1.com"}
	p2 := &Profile{Host: "test2.com"}
	m.Add("test1", p1)
	m.Add("test2", p2)

	// Set active profile
	m.SetActive("test1")
	if m.activeProfile != "test1" {
		t.Errorf("expected activeProfile to be 'test1', got %q", m.activeProfile)
	}

	// Get active profile
	got, err := m.Active()
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}
	if got != p1 {
		t.Error("expected to get test1 profile")
	}

	// Set different active profile
	m.SetActive("test2")
	got, err = m.Active()
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}
	if got != p2 {
		t.Error("expected to get test2 profile")
	}

	// Set active to non-existent profile
	m.SetActive("nonexistent")
	_, err = m.Active()
	if err == nil {
		t.Error("expected error for non-existent active profile")
	}
}

// TestManagerList verifies that List returns all profile names.
func TestManagerList(t *testing.T) {
	m := NewManager()

	// Test empty manager
	names := m.List()
	if len(names) != 0 {
		t.Errorf("expected empty list, got %d names", len(names))
	}

	// Add profiles
	m.Add("test1", &Profile{})
	m.Add("test2", &Profile{})
	m.Add("test3", &Profile{})

	names = m.List()
	if len(names) != 3 {
		t.Errorf("expected 3 names, got %d", len(names))
	}

	// Verify all names are present
	nameMap := make(map[string]bool)
	for _, name := range names {
		nameMap[name] = true
	}

	if !nameMap["test1"] || !nameMap["test2"] || !nameMap["test3"] {
		t.Errorf("expected all profile names to be in list, got %v", names)
	}
}

// TestNewLoader verifies that NewLoader creates a loader with correct initialization.
func TestNewLoader(t *testing.T) {
	values := map[string]interface{}{"host": "example.com"}
	defaults := map[string]interface{}{"port": 443}
	overrides := map[string]interface{}{"username": "admin"}

	l := NewLoader(values, defaults, overrides)

	if l.values == nil || l.defaults == nil || l.overrides == nil {
		t.Error("expected all maps to be initialized")
	}
}

// TestLoaderLoad verifies that Load correctly merges values with proper precedence.
func TestLoaderLoad(t *testing.T) {
	tests := []struct {
		name      string
		values    map[string]interface{}
		defaults  map[string]interface{}
		overrides map[string]interface{}
		expected  *Profile
	}{
		{
			name:      "empty loader uses built-in defaults",
			values:    map[string]interface{}{},
			defaults:  map[string]interface{}{},
			overrides: map[string]interface{}{},
			expected: &Profile{
				Host:   "localhost",
				Port:   0,
				UseTLS: true,
				Verify: true,
			},
		},
		{
			name: "values override built-in defaults",
			values: map[string]interface{}{
				"host":     "example.com",
				"port":     8443,
				"username": "user",
			},
			defaults:  map[string]interface{}{},
			overrides: map[string]interface{}{},
			expected: &Profile{
				Host:     "example.com",
				Port:     8443,
				Username: "user",
				UseTLS:   true,
				Verify:   true,
			},
		},
		{
			name:   "defaults provide base values",
			values: map[string]interface{}{},
			defaults: map[string]interface{}{
				"host":     "default.com",
				"port":     443,
				"username": "defaultuser",
			},
			overrides: map[string]interface{}{},
			expected: &Profile{
				Host:     "default.com",
				Port:     443,
				Username: "defaultuser",
				UseTLS:   true,
				Verify:   true,
			},
		},
		{
			name: "values override defaults",
			values: map[string]interface{}{
				"host": "value.com",
				"port": 9000,
			},
			defaults: map[string]interface{}{
				"host":     "default.com",
				"port":     443,
				"username": "defaultuser",
			},
			overrides: map[string]interface{}{},
			expected: &Profile{
				Host:     "value.com",
				Port:     9000,
				Username: "defaultuser",
				UseTLS:   true,
				Verify:   true,
			},
		},
		{
			name: "overrides take highest precedence",
			values: map[string]interface{}{
				"host":     "value.com",
				"port":     9000,
				"username": "valueuser",
			},
			defaults: map[string]interface{}{
				"host":     "default.com",
				"port":     443,
				"username": "defaultuser",
			},
			overrides: map[string]interface{}{
				"host": "override.com",
				"port": 10000,
			},
			expected: &Profile{
				Host:     "override.com",
				Port:     10000,
				Username: "valueuser",
				UseTLS:   true,
				Verify:   true,
			},
		},
		{
			name: "boolean fields handled correctly",
			values: map[string]interface{}{
				"use_tls": false,
				"verify":  false,
			},
			defaults:  map[string]interface{}{},
			overrides: map[string]interface{}{},
			expected: &Profile{
				Host:   "localhost",
				UseTLS: false,
				Verify: false,
			},
		},
		{
			name: "all profile fields",
			values: map[string]interface{}{
				"host":          "api.example.com",
				"port":          8443,
				"use_tls":       true,
				"verify":        false,
				"username":      "admin",
				"password":      "secret",
				"client_id":     "client123",
				"client_secret": "clientsecret",
				"mongo_url":     "mongodb://localhost:27017",
				"timeout":       30,
			},
			defaults:  map[string]interface{}{},
			overrides: map[string]interface{}{},
			expected: &Profile{
				Host:         "api.example.com",
				Port:         8443,
				UseTLS:       true,
				Verify:       false,
				Username:     "admin",
				Password:     "secret",
				ClientID:     "client123",
				ClientSecret: "clientsecret",
				MongoUrl:     "mongodb://localhost:27017",
				Timeout:      30,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			l := NewLoader(tt.values, tt.defaults, tt.overrides)
			got := l.Load()

			if got.Host != tt.expected.Host {
				t.Errorf("Host: expected %q, got %q", tt.expected.Host, got.Host)
			}
			if got.Port != tt.expected.Port {
				t.Errorf("Port: expected %d, got %d", tt.expected.Port, got.Port)
			}
			if got.UseTLS != tt.expected.UseTLS {
				t.Errorf("UseTLS: expected %v, got %v", tt.expected.UseTLS, got.UseTLS)
			}
			if got.Verify != tt.expected.Verify {
				t.Errorf("Verify: expected %v, got %v", tt.expected.Verify, got.Verify)
			}
			if got.Username != tt.expected.Username {
				t.Errorf("Username: expected %q, got %q", tt.expected.Username, got.Username)
			}
			if got.Password != tt.expected.Password {
				t.Errorf("Password: expected %q, got %q", tt.expected.Password, got.Password)
			}
			if got.ClientID != tt.expected.ClientID {
				t.Errorf("ClientID: expected %q, got %q", tt.expected.ClientID, got.ClientID)
			}
			if got.ClientSecret != tt.expected.ClientSecret {
				t.Errorf("ClientSecret: expected %q, got %q", tt.expected.ClientSecret, got.ClientSecret)
			}
			if got.MongoUrl != tt.expected.MongoUrl {
				t.Errorf("MongoUrl: expected %q, got %q", tt.expected.MongoUrl, got.MongoUrl)
			}
			if got.Timeout != tt.expected.Timeout {
				t.Errorf("Timeout: expected %d, got %d", tt.expected.Timeout, got.Timeout)
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

// TestGetIntTypeConversions verifies getInt handles different types correctly.
func TestGetIntTypeConversions(t *testing.T) {
	tests := []struct {
		name     string
		value    interface{}
		fallback int
		expected int
	}{
		{"nil returns default", nil, 42, 42},
		{"int returns value", 123, 42, 123},
		{"int64 returns value", int64(456), 42, 456},
		{"float64 returns truncated value", 789.5, 42, 789},
		{"string number returns parsed value", "999", 42, 999},
		{"invalid string returns default", "invalid", 42, 42},
		{"bool returns default", true, 42, 42},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := getInt(tt.value, tt.fallback)
			if got != tt.expected {
				t.Errorf("expected %d, got %d", tt.expected, got)
			}
		})
	}
}

// TestGetBoolTypeConversions verifies getBool handles different types correctly.
func TestGetBoolTypeConversions(t *testing.T) {
	tests := []struct {
		name     string
		value    interface{}
		fallback bool
		expected bool
	}{
		{"nil returns default", nil, true, true},
		{"bool true returns true", true, false, true},
		{"bool false returns false", false, true, false},
		{"string 'true' returns true", "true", false, true},
		{"string 'false' returns false", "false", true, false},
		{"string 'yes' returns false", "yes", true, false},
		{"int returns default", 1, false, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := getBool(tt.value, tt.fallback)
			if got != tt.expected {
				t.Errorf("expected %v, got %v", tt.expected, got)
			}
		})
	}
}
