// Copyright 2024 Itential Inc. All Rights Reserved
// Unauthorized copying of this file, via any medium is strictly prohibited
// Proprietary and confidential

package app

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewSettings(t *testing.T) {
	settings := NewSettings()

	// Verify all default values are set correctly
	assert.Equal(t, defaultWorkingDir, settings.WorkingDir, "WorkingDir should match default")
	assert.Equal(t, defaultDefaultProfile, settings.DefaultProfile, "DefaultProfile should match default")
	assert.Equal(t, defaultDefaultRepository, settings.DefaultRepository, "DefaultRepository should match default")
	assert.Equal(t, defaultDatasetsEnabled, settings.Features.DatasetsEnabled, "DatasetsEnabled should match default")
	assert.Equal(t, defaultGitName, settings.Git.Name, "Git.Name should match default")
	assert.Equal(t, defaultGitEmail, settings.Git.Email, "Git.Email should match default")
	assert.Equal(t, defaultGitUser, settings.Git.User, "Git.User should match default")
}

func TestSettings_GetWorkingDir(t *testing.T) {
	testCases := []struct {
		name       string
		workingDir string
	}{
		{
			name:       "Default working directory",
			workingDir: defaultWorkingDir,
		},
		{
			name:       "Custom working directory",
			workingDir: "/custom/path",
		},
		{
			name:       "Empty working directory",
			workingDir: "",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			settings := &Settings{WorkingDir: tc.workingDir}
			assert.Equal(t, tc.workingDir, settings.GetWorkingDir())
		})
	}
}

func TestSettings_GetDefaultProfile(t *testing.T) {
	testCases := []struct {
		name           string
		defaultProfile string
	}{
		{
			name:           "Default profile (empty)",
			defaultProfile: "",
		},
		{
			name:           "Custom default profile",
			defaultProfile: "production",
		},
		{
			name:           "Development profile",
			defaultProfile: "dev",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			settings := &Settings{DefaultProfile: tc.defaultProfile}
			assert.Equal(t, tc.defaultProfile, settings.GetDefaultProfile())
		})
	}
}

func TestSettings_GetDefaultRepository(t *testing.T) {
	testCases := []struct {
		name              string
		defaultRepository string
	}{
		{
			name:              "Default repository (empty)",
			defaultRepository: "",
		},
		{
			name:              "Custom default repository",
			defaultRepository: "main-repo",
		},
		{
			name:              "Named repository",
			defaultRepository: "my-artifacts",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			settings := &Settings{DefaultRepository: tc.defaultRepository}
			assert.Equal(t, tc.defaultRepository, settings.GetDefaultRepository())
		})
	}
}

func TestSettings_IsDatasetsEnabled(t *testing.T) {
	testCases := []struct {
		name            string
		datasetsEnabled bool
	}{
		{
			name:            "Datasets disabled (default)",
			datasetsEnabled: false,
		},
		{
			name:            "Datasets enabled",
			datasetsEnabled: true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			settings := &Settings{
				Features: Features{
					DatasetsEnabled: tc.datasetsEnabled,
				},
			}
			assert.Equal(t, tc.datasetsEnabled, settings.IsDatasetsEnabled())
		})
	}
}

func TestSettings_GetGitName(t *testing.T) {
	testCases := []struct {
		name    string
		gitName string
	}{
		{
			name:    "Default git name (empty)",
			gitName: "",
		},
		{
			name:    "Custom git name",
			gitName: "John Doe",
		},
		{
			name:    "Git name with special characters",
			gitName: "José García-López",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			settings := &Settings{
				Git: GitConfig{Name: tc.gitName},
			}
			assert.Equal(t, tc.gitName, settings.GetGitName())
		})
	}
}

func TestSettings_GetGitEmail(t *testing.T) {
	testCases := []struct {
		name     string
		gitEmail string
	}{
		{
			name:     "Default git email (empty)",
			gitEmail: "",
		},
		{
			name:     "Custom git email",
			gitEmail: "john.doe@example.com",
		},
		{
			name:     "Git email with plus addressing",
			gitEmail: "user+tag@example.com",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			settings := &Settings{
				Git: GitConfig{Email: tc.gitEmail},
			}
			assert.Equal(t, tc.gitEmail, settings.GetGitEmail())
		})
	}
}

func TestSettings_GetGitUser(t *testing.T) {
	testCases := []struct {
		name    string
		gitUser string
	}{
		{
			name:    "Default git user",
			gitUser: defaultGitUser,
		},
		{
			name:    "Custom git user",
			gitUser: "johndoe",
		},
		{
			name:    "Empty git user",
			gitUser: "",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			settings := &Settings{
				Git: GitConfig{User: tc.gitUser},
			}
			assert.Equal(t, tc.gitUser, settings.GetGitUser())
		})
	}
}

func TestDefaultValues(t *testing.T) {
	defaults := DefaultValues()

	// Verify all expected keys are present
	expectedKeys := []string{
		"application.working_dir",
		"application.default_profile",
		"application.default_repository",
		"features.datasets_enabled",
		"git.name",
		"git.email",
		"git.user",
	}

	for _, key := range expectedKeys {
		_, exists := defaults[key]
		assert.True(t, exists, "Key %s should exist in defaults", key)
	}

	// Verify specific values
	assert.Equal(t, defaultWorkingDir, defaults["application.working_dir"])
	assert.Equal(t, defaultDefaultProfile, defaults["application.default_profile"])
	assert.Equal(t, defaultDefaultRepository, defaults["application.default_repository"])
	assert.Equal(t, defaultDatasetsEnabled, defaults["features.datasets_enabled"])
	assert.Equal(t, defaultGitName, defaults["git.name"])
	assert.Equal(t, defaultGitEmail, defaults["git.email"])
	assert.Equal(t, defaultGitUser, defaults["git.user"])
}

func TestDefaultEnvBindings(t *testing.T) {
	bindings := DefaultEnvBindings()

	// Verify all expected keys are present
	expectedKeys := []string{
		"application.working_dir",
		"application.default_profile",
		"application.default_repository",
		"features.datasets_enabled",
		"git.name",
		"git.email",
		"git.user",
	}

	for _, key := range expectedKeys {
		_, exists := bindings[key]
		assert.True(t, exists, "Key %s should exist in bindings", key)
	}

	// Verify environment variable names follow the IPCTL_* convention
	testCases := []struct {
		key    string
		envVar string
	}{
		{"application.working_dir", "IPCTL_APPLICATION_WORKING_DIR"},
		{"application.default_profile", "IPCTL_APPLICATION_DEFAULT_PROFILE"},
		{"application.default_repository", "IPCTL_APPLICATION_DEFAULT_REPOSITORY"},
		{"features.datasets_enabled", "IPCTL_FEATURES_DATASETS_ENABLED"},
		{"git.name", "IPCTL_GIT_NAME"},
		{"git.email", "IPCTL_GIT_EMAIL"},
		{"git.user", "IPCTL_GIT_USER"},
	}

	for _, tc := range testCases {
		t.Run(tc.key, func(t *testing.T) {
			assert.Equal(t, tc.envVar, bindings[tc.key], "Environment variable for %s should be %s", tc.key, tc.envVar)
		})
	}
}

func TestSettings_ImplementsInterfaces(t *testing.T) {
	settings := NewSettings()

	// Verify Settings implements all provider interfaces
	var _ ApplicationProvider = settings
	var _ FeaturesProvider = settings
	var _ GitProvider = settings
	var _ SettingsProvider = settings

	// This test will fail to compile if Settings doesn't implement the interfaces
	t.Log("Settings implements all required interfaces")
}

func TestSettings_CompleteConfiguration(t *testing.T) {
	// Test a fully configured Settings instance
	settings := &Settings{
		WorkingDir:        "/home/user/.platform.d",
		DefaultProfile:    "production",
		DefaultRepository: "main-repo",
		Features: Features{
			DatasetsEnabled: true,
		},
		Git: GitConfig{
			Name:  "John Doe",
			Email: "john.doe@example.com",
			User:  "jdoe",
		},
	}

	// Verify all getters return the correct values
	assert.Equal(t, "/home/user/.platform.d", settings.GetWorkingDir())
	assert.Equal(t, "production", settings.GetDefaultProfile())
	assert.Equal(t, "main-repo", settings.GetDefaultRepository())
	assert.True(t, settings.IsDatasetsEnabled())
	assert.Equal(t, "John Doe", settings.GetGitName())
	assert.Equal(t, "john.doe@example.com", settings.GetGitEmail())
	assert.Equal(t, "jdoe", settings.GetGitUser())
}
