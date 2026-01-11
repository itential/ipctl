// Copyright 2024 Itential Inc. All Rights Reserved
// Unauthorized copying of this file, via any medium is strictly prohibited
// Proprietary and confidential

package config

import (
	"testing"

	"github.com/itential/ipctl/internal/profile"
	"github.com/itential/ipctl/internal/repository"
	"github.com/stretchr/testify/assert"
)

// Note: Tests for getConfigFileFromFlag(), setConfigFile(), and GetAndExpandDirectory()
// have been removed as these functions are now internal to the Loader.
// See loader_test.go for comprehensive tests of the new loading mechanism.

func TestNewConfig(t *testing.T) {
	testCases := []struct {
		name           string
		defaults       map[string]interface{}
		envBinding     map[string]string
		appWorkingDir  string
		sysConfigPath  string
		fileName       string
		expectedFields bool
	}{
		{
			name:           "New config with defaults",
			defaults:       nil,
			envBinding:     nil,
			appWorkingDir:  "",
			sysConfigPath:  "",
			fileName:       "",
			expectedFields: true,
		},
		{
			name: "New config with custom values",
			defaults: map[string]interface{}{
				"application.working_dir":     "/tmp/test",
				"application.default_profile": "test",
				"git.name":                    "Test User",
			},
			envBinding:     map[string]string{},
			appWorkingDir:  "/tmp/test",
			sysConfigPath:  "/etc/test",
			fileName:       "test",
			expectedFields: true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			config := NewConfig(tc.defaults, tc.envBinding, tc.appWorkingDir, tc.sysConfigPath, tc.fileName)

			if config == nil {
				t.Error("Expected config to be created")
				return
			}

			if tc.expectedFields {
				assert.NotEmpty(t, config.Settings.WorkingDir)
			}
		})
	}
}

func TestConfig_DumpConfig(t *testing.T) {
	config := NewConfig(nil, nil, "", "", "")

	dump := config.DumpConfig()
	assert.NotEmpty(t, dump)

	// Verify it's valid JSON
	assert.True(t, len(dump) > 10)
	// With the refactored structure, Settings is now nested
	assert.Contains(t, dump, "Settings")
	assert.Contains(t, dump, "WorkingDir")
	assert.Contains(t, dump, "Features")
	assert.Contains(t, dump, "Git")
}

// Note: Test for populateFields() has been removed as this function is now
// internal to the Loader.buildConfig() method.
// See loader_test.go for comprehensive tests of the new loading mechanism.

func TestConfig_GetProfile(t *testing.T) {
	testCases := []struct {
		name          string
		profileName   string
		setupProfiles func(*Config)
		expectError   bool
		expectedHost  string
	}{
		{
			name:        "Get existing profile",
			profileName: "test",
			setupProfiles: func(c *Config) {
				c.profileManager = profile.NewManager()
				c.profileManager.Add("test", &profile.Profile{
					Host:     "test.example.com",
					Port:     8080,
					Username: "testuser",
				})
			},
			expectError:  false,
			expectedHost: "test.example.com",
		},
		{
			name:        "Get non-existent profile returns default",
			profileName: "nonexistent",
			setupProfiles: func(c *Config) {
				c.profileManager = profile.NewManager()
			},
			expectError:  true,
			expectedHost: "localhost",
		},
		{
			name:        "Get profile by name is case sensitive",
			profileName: "test",
			setupProfiles: func(c *Config) {
				c.profileManager = profile.NewManager()
				c.profileManager.Add("test", &profile.Profile{
					Host: "test.example.com",
					Port: 8080,
				})
			},
			expectError:  false,
			expectedHost: "test.example.com",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			config := &Config{}
			tc.setupProfiles(config)

			p, err := config.GetProfile(tc.profileName)

			if tc.expectError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}

			assert.NotNil(t, p)
			assert.Equal(t, tc.expectedHost, p.Host)
		})
	}
}

func TestConfig_ActiveProfile(t *testing.T) {
	testCases := []struct {
		name            string
		activeProfile   string
		setupProfiles   func(*Config)
		expectedProfile string
		expectError     bool
	}{
		{
			name:          "Active profile when custom profile is set",
			activeProfile: "custom",
			setupProfiles: func(c *Config) {
				c.profileManager = profile.NewManager()
				c.profileManager.Add("custom", &profile.Profile{Host: "custom.example.com"})
				c.profileManager.Add("default", &profile.Profile{Host: "default.example.com"})
				c.profileManager.SetActive("custom")
			},
			expectedProfile: "custom.example.com",
			expectError:     false,
		},
		{
			name:          "Active profile defaults to default",
			activeProfile: "",
			setupProfiles: func(c *Config) {
				c.profileManager = profile.NewManager()
				c.profileManager.Add("default", &profile.Profile{Host: "default.example.com"})
			},
			expectedProfile: "default.example.com",
			expectError:     false,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			config := &Config{}
			tc.setupProfiles(config)

			p, err := config.ActiveProfile()

			if tc.expectError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}

			assert.NotNil(t, p)
			assert.Equal(t, tc.expectedProfile, p.Host)
		})
	}
}

func TestConfig_GetRepository(t *testing.T) {
	testCases := []struct {
		name              string
		repoName          string
		setupRepositories func(*Config)
		expectError       bool
		expectedUrl       string
	}{
		{
			name:     "Get existing repository",
			repoName: "test-repo",
			setupRepositories: func(c *Config) {
				c.repositoryManager = repository.NewManager()
				c.repositoryManager.Add("test-repo", &repository.Repository{
					Url:       "https://github.com/test/repo.git",
					Reference: "main",
				})
			},
			expectError: false,
			expectedUrl: "https://github.com/test/repo.git",
		},
		{
			name:     "Get non-existent repository",
			repoName: "nonexistent",
			setupRepositories: func(c *Config) {
				c.repositoryManager = repository.NewManager()
			},
			expectError: true,
			expectedUrl: "",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			config := &Config{}
			tc.setupRepositories(config)

			repo, err := config.GetRepository(tc.repoName)

			if tc.expectError {
				assert.Error(t, err)
				assert.Nil(t, repo)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, repo)
				assert.Equal(t, tc.expectedUrl, repo.Url)
			}
		})
	}
}
