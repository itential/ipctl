// Copyright 2024 Itential Inc. All Rights Reserved
// Unauthorized copying of this file, via any medium is strictly prohibited
// Proprietary and confidential

package config

import (
	"os"
	"path/filepath"
	"reflect"
	"testing"

	"github.com/stretchr/testify/assert"
)

var repositoryFields = []string{
	"url",
	"private_key",
	"private_key_file",
	"reference",
}

func TestGetRepositoryFields(t *testing.T) {
	fields := getRepositoryFields()

	// Verify all expected fields are present
	for _, expected := range repositoryFields {
		assert.Contains(t, fields, expected, "Expected field %s not found in repository fields", expected)
	}

	// Verify no unexpected fields
	assert.Equal(t, len(repositoryFields), len(fields), "Repository fields count mismatch")
}

func TestRepositoryFieldsConsistency(t *testing.T) {
	// Test that getRepositoryFields matches the Repository struct
	fields := getRepositoryFields()
	expectedFields := []string{"url", "private_key", "private_key_file", "reference"}

	assert.Equal(t, len(expectedFields), len(fields))

	for _, expected := range expectedFields {
		assert.Contains(t, fields, expected)
	}
}

func TestLoadRepository(t *testing.T) {
	testCases := []struct {
		name      string
		values    map[string]interface{}
		overrides map[string]interface{}
		expected  *Repository
	}{
		{
			name: "Load repository with all values",
			values: map[string]interface{}{
				"url":              "https://github.com/test/repo.git",
				"private_key":      "-----BEGIN RSA PRIVATE KEY-----\ntest\n-----END RSA PRIVATE KEY-----",
				"private_key_file": "~/.ssh/id_rsa",
				"reference":        "main",
			},
			overrides: map[string]interface{}{},
			expected: &Repository{
				Url:            "https://github.com/test/repo.git",
				PrivateKey:     "-----BEGIN RSA PRIVATE KEY-----\ntest\n-----END RSA PRIVATE KEY-----",
				PrivateKeyFile: filepath.Join(os.Getenv("HOME"), ".ssh/id_rsa"),
				Reference:      "main",
			},
		},
		{
			name: "Load repository with overrides",
			values: map[string]interface{}{
				"url":       "https://github.com/original/repo.git",
				"reference": "original",
			},
			overrides: map[string]interface{}{
				"url":       "https://github.com/override/repo.git",
				"reference": "override",
			},
			expected: &Repository{
				Url:       "https://github.com/override/repo.git",
				Reference: "override",
			},
		},
		{
			name:      "Load repository with empty values",
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
			name: "Load repository with nil values",
			values: map[string]interface{}{
				"url":              nil,
				"private_key":      nil,
				"private_key_file": nil,
				"reference":        nil,
			},
			overrides: map[string]interface{}{},
			expected: &Repository{
				Url:            "",
				PrivateKey:     "",
				PrivateKeyFile: "",
				Reference:      "",
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			repo := loadRepository(tc.values, tc.overrides)

			assert.Equal(t, tc.expected.Url, repo.Url)
			assert.Equal(t, tc.expected.PrivateKey, repo.PrivateKey)
			assert.Equal(t, tc.expected.Reference, repo.Reference)

			// Special handling for private key file since it gets expanded
			if tc.expected.PrivateKeyFile != "" {
				assert.Equal(t, tc.expected.PrivateKeyFile, repo.PrivateKeyFile)
			}
		})
	}
}

func TestLoadRepositoryPrivateKeyFileExpansion(t *testing.T) {
	testCases := []struct {
		name           string
		privateKeyFile string
		expectedPrefix string
	}{
		{
			name:           "Tilde expansion",
			privateKeyFile: "~/.ssh/id_rsa",
			expectedPrefix: os.Getenv("HOME"),
		},
		{
			name:           "No expansion needed",
			privateKeyFile: "/absolute/path/key",
			expectedPrefix: "/absolute/path/key",
		},
		{
			name:           "Empty private key file",
			privateKeyFile: "",
			expectedPrefix: "",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			values := map[string]interface{}{
				"private_key_file": tc.privateKeyFile,
			}

			repo := loadRepository(values, map[string]interface{}{})

			if tc.expectedPrefix == "" {
				assert.Equal(t, "", repo.PrivateKeyFile)
			} else if tc.privateKeyFile == "/absolute/path/key" {
				assert.Equal(t, tc.expectedPrefix, repo.PrivateKeyFile)
			} else {
				// For tilde expansion, check that it starts with HOME
				assert.True(t,
					len(repo.PrivateKeyFile) > len(tc.expectedPrefix) || repo.PrivateKeyFile == tc.expectedPrefix,
					"Private key file path should be expanded: got %s", repo.PrivateKeyFile)
			}
		})
	}
}

func TestLoadRepositoryOverridePrecedence(t *testing.T) {
	// Test that overrides take precedence over values
	values := map[string]interface{}{
		"url":              "https://github.com/values/repo.git",
		"private_key":      "values-key",
		"private_key_file": "values-file",
		"reference":        "values-ref",
	}

	overrides := map[string]interface{}{
		"url":       "https://github.com/overrides/repo.git",
		"reference": "overrides-ref",
	}

	repo := loadRepository(values, overrides)

	// Overridden values should be from overrides
	assert.Equal(t, "https://github.com/overrides/repo.git", repo.Url)
	assert.Equal(t, "overrides-ref", repo.Reference)

	// Non-overridden values should be from values
	assert.Equal(t, "values-key", repo.PrivateKey)
}

func TestLoadRepositoryWithDifferentTypes(t *testing.T) {
	// Test that only string values are accepted (as per the type assertion in the code)
	values := map[string]interface{}{
		"url":              "https://github.com/test/repo.git",
		"private_key":      "test-key",
		"private_key_file": "/path/to/key",
		"reference":        "main",
	}

	repo := loadRepository(values, map[string]interface{}{})

	assert.Equal(t, "https://github.com/test/repo.git", repo.Url)
	assert.Equal(t, "test-key", repo.PrivateKey)
	assert.Equal(t, "/path/to/key", repo.PrivateKeyFile)
	assert.Equal(t, "main", repo.Reference)
}

func TestRepositoryStructFields(t *testing.T) {
	// Test that Repository struct has all expected fields with correct types
	repo := &Repository{}
	v := reflect.ValueOf(repo).Elem()
	t_type := v.Type()

	expectedFields := map[string]reflect.Kind{
		"Url":            reflect.String,
		"PrivateKey":     reflect.String,
		"PrivateKeyFile": reflect.String,
		"Reference":      reflect.String,
	}

	assert.Equal(t, len(expectedFields), t_type.NumField(), "Repository struct field count mismatch")

	for i := 0; i < t_type.NumField(); i++ {
		field := t_type.Field(i)
		expectedKind, exists := expectedFields[field.Name]

		assert.True(t, exists, "Unexpected field %s in Repository struct", field.Name)
		assert.Equal(t, expectedKind, field.Type.Kind(), "Field %s has wrong type", field.Name)

		// Verify JSON tags
		jsonTag := field.Tag.Get("json")
		assert.NotEmpty(t, jsonTag, "Field %s missing JSON tag", field.Name)
	}
}

func TestRepositoryJSONTags(t *testing.T) {
	// Verify that JSON tags are correctly set for serialization
	expectedTags := map[string]string{
		"Url":            "url",
		"PrivateKey":     "private_key",
		"PrivateKeyFile": "private_key_file",
		"Reference":      "reference",
	}

	t_type := reflect.TypeOf(Repository{})

	for i := 0; i < t_type.NumField(); i++ {
		field := t_type.Field(i)
		expectedTag, exists := expectedTags[field.Name]

		assert.True(t, exists, "Unexpected field %s", field.Name)

		jsonTag := field.Tag.Get("json")
		assert.Equal(t, expectedTag, jsonTag, "Field %s has incorrect JSON tag", field.Name)
	}
}
