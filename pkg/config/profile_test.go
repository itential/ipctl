// Copyright 2024 Itential Inc. All Rights Reserved
// Unauthorized copying of this file, via any medium is strictly prohibited
// Proprietary and confidential

package config

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

var profileFields = []string{
	"host",
	"port",
	"use_tls",
	"verify",
	"username",
	"password",
	"client_id",
	"client_secret",
	"timeout",
	"mongo_url",
}

func TestGetProfileFields(t *testing.T) {
	for _, ele := range getProfileFields() {
		assert.Contains(t, profileFields, ele, "unknown field detected in profile")

	}
}

func TestLoadProfileWithValues(t *testing.T) {
	values := map[string]interface{}{
		"host":          "test",
		"port":          1000,
		"use_tls":       false,
		"verify":        true,
		"username":      "test",
		"password":      "test",
		"client_id":     "test",
		"client_secret": "test",
		"timeout":       1000,
		"mongo_url":     "test",
	}

	p := loadProfile(values, map[string]interface{}{}, map[string]interface{}{})

	for key, value := range values {
		switch key {
		case "host":
			assert.Equal(t, p.Host, value)
		case "port":
			assert.Equal(t, p.Port, value)
		case "use_tls":
			assert.Equal(t, p.UseTLS, value)
		case "verify":
			assert.Equal(t, p.Verify, value)
		case "username":
			assert.Equal(t, p.Username, value)
		case "password":
			assert.Equal(t, p.Password, value)
		case "client_id":
			assert.Equal(t, p.ClientID, value)
		case "client_secret":
			assert.Equal(t, p.ClientSecret, value)
		case "timeout":
			assert.Equal(t, p.Timeout, value)
		case "mongo_url":
			assert.Equal(t, p.MongoUrl, value)
		}
	}
}

func TestLoadProfileWithDefaults(t *testing.T) {
	// Test loading profile with only defaults (empty values)
	defaults := map[string]interface{}{
		"host":          "default.example.com",
		"port":          8443,
		"use_tls":       true,
		"verify":        false,
		"username":      "defaultuser",
		"password":      "defaultpass",
		"client_id":     "defaultclient",
		"client_secret": "defaultsecret",
		"timeout":       30000,
		"mongo_url":     "mongodb://default:27017",
	}

	p := loadProfile(map[string]interface{}{}, defaults, map[string]interface{}{})

	assert.Equal(t, "default.example.com", p.Host)
	assert.Equal(t, 8443, p.Port)
	assert.True(t, p.UseTLS)
	assert.False(t, p.Verify)
	assert.Equal(t, "defaultuser", p.Username)
	assert.Equal(t, "defaultpass", p.Password)
	assert.Equal(t, "defaultclient", p.ClientID)
	assert.Equal(t, "defaultsecret", p.ClientSecret)
	assert.Equal(t, 30000, p.Timeout)
	assert.Equal(t, "mongodb://default:27017", p.MongoUrl)
}

func TestLoadProfileWithOverrides(t *testing.T) {
	values := map[string]interface{}{
		"host":     "values.example.com",
		"port":     9000,
		"username": "valuesuser",
	}

	defaults := map[string]interface{}{
		"host":     "default.example.com",
		"port":     8443,
		"username": "defaultuser",
		"password": "defaultpass",
	}

	overrides := map[string]interface{}{
		"host":     "override.example.com",
		"password": "overridepass",
	}

	p := loadProfile(values, defaults, overrides)

	// Overrides should take precedence
	assert.Equal(t, "override.example.com", p.Host)
	assert.Equal(t, "overridepass", p.Password)

	// Values should be used when no override exists
	assert.Equal(t, 9000, p.Port)
	assert.Equal(t, "valuesuser", p.Username)
}

func TestLoadProfileWithStringConversions(t *testing.T) {
	testCases := []struct {
		name     string
		values   map[string]interface{}
		expected func(*Profile) bool
	}{
		{
			name: "String values for boolean fields",
			values: map[string]interface{}{
				"use_tls": "true",
				"verify":  "false",
			},
			expected: func(p *Profile) bool {
				return p.UseTLS == true && p.Verify == false
			},
		},
		{
			name: "String values for integer fields",
			values: map[string]interface{}{
				"port":    "8080",
				"timeout": "5000",
			},
			expected: func(p *Profile) bool {
				return p.Port == 8080 && p.Timeout == 5000
			},
		},
		{
			name: "Boolean values for boolean fields",
			values: map[string]interface{}{
				"use_tls": true,
				"verify":  false,
			},
			expected: func(p *Profile) bool {
				return p.UseTLS == true && p.Verify == false
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			p := loadProfile(tc.values, map[string]interface{}{}, map[string]interface{}{})
			assert.True(t, tc.expected(p))
		})
	}
}

func TestLoadProfileWithEmptyValues(t *testing.T) {
	// Test that empty or nil values fall back to defaults correctly
	values := map[string]interface{}{
		"host":     "",
		"username": "",
		"password": "",
		"port":     nil,
		"use_tls":  nil,
	}

	p := loadProfile(values, map[string]interface{}{}, map[string]interface{}{})

	// Should use built-in defaults
	assert.Equal(t, defaultHost, p.Host)
	assert.Equal(t, defaultUsername, p.Username)
	assert.Equal(t, defaultPassword, p.Password)
	assert.Equal(t, defaultUseTLS, p.UseTLS)
}

func TestGetProfileFromFlag(t *testing.T) {
	testCases := []struct {
		name            string
		args            []string
		expectedProfile string
	}{
		{
			name:            "No profile flag",
			args:            []string{"ipctl", "version"},
			expectedProfile: "",
		},
		{
			name:            "Profile flag with value",
			args:            []string{"ipctl", "--profile", "test", "version"},
			expectedProfile: "test",
		},
		{
			name:            "Profile flag with different command",
			args:            []string{"ipctl", "--profile", "production", "workflows", "list"},
			expectedProfile: "production",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Save original args
			originalArgs := os.Args
			defer func() { os.Args = originalArgs }()

			// Set test args
			os.Args = tc.args

			profile := getProfileFromFlag()
			assert.Equal(t, tc.expectedProfile, profile)
		})
	}
}

// NOTE: Commented out edge case tests that would call os.Exit
// These would need to be tested in integration tests or with process isolation
//
// func TestLoadProfileEdgeCases(t *testing.T) {
// 	// Tests for invalid port/timeout strings that call handleError -> os.Exit
// 	// Cannot be easily unit tested without process isolation
// }

func TestProfileFieldsConsistency(t *testing.T) {
	// Verify that getProfileFields returns all expected fields
	fields := getProfileFields()
	expectedFields := []string{
		"host", "port", "use_tls", "verify",
		"username", "password", "client_id", "client_secret",
		"mongo_url", "timeout",
	}

	assert.Equal(t, len(expectedFields), len(fields))

	for _, expected := range expectedFields {
		assert.Contains(t, fields, expected, "Expected field %s not found in profile fields", expected)
	}
}
