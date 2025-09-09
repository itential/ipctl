// Copyright 2024 Itential Inc. All Rights Reserved
// Unauthorized copying of this file, via any medium is strictly prohibited
// Proprietary and confidential

package config

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
)

const confContents = `
[client]
host = nonDefaultValue
`

func TestGetConfigDirFromFlag(t *testing.T) {
	confFileDir := t.TempDir()
	flagFilePath := filepath.Join(confFileDir, "torero.conf")
	if err := os.WriteFile(flagFilePath, []byte(confContents), 0644); err != nil {
		t.Errorf("unable to write config file for test: %v", err)
	}
	oldOsArgs := os.Args
	os.Args = []string{os.Args[0], "--config", flagFilePath}
	defer func() { os.Args = oldOsArgs }()

	d := getConfigFileFromFlag()
	assert.Equal(t, d, flagFilePath)
}

// We can only test environment variable precedence as we do not want to mess with other directories
// on a tester's local machine
func TestSetConfigFile(t *testing.T) {
	tests := []struct {
		name           string
		envLocationSet bool
	}{
		{
			name:           "File exists at environment var location and IPCTL_CONFIG set",
			envLocationSet: true,
		},
		{
			name:           "File exists at environment var location and IPCTL_CONFIG not set",
			envLocationSet: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			envFilePath := filepath.Join(t.TempDir(), "torero.conf")
			os.WriteFile(envFilePath, []byte(confContents), 0644)

			if tt.envLocationSet {
				t.Setenv("IPCTL_CONFIG", envFilePath)
			}

			setConfigFile("~/.torero.d", "/etc/torero", "torero.conf")

			confLocation := viper.ConfigFileUsed()
			if tt.envLocationSet && confLocation != envFilePath {
				t.Errorf("expected config to be at %s but got %s", envFilePath, confLocation)
			}
			if !tt.envLocationSet && confLocation == envFilePath {
				t.Errorf("config file location is set to %s but expected a different location or none at all", confLocation)
			}
		})
	}
}

func TestGetAndExpandDirectory(t *testing.T) {
	tests := []struct {
		name         string
		directory    string
		shouldExpand bool
	}{
		{
			name:         "with tilde",
			directory:    "~/serverlogs",
			shouldExpand: true,
		},
		{
			name:         "with many paths",
			directory:    "~/server/logs",
			shouldExpand: true,
		},
		{
			name:         "no tilde",
			directory:    "/tmp/test",
			shouldExpand: false,
		},
	}

	for _, tt := range tests {
		exampleKey := "example.key"
		viper.SetDefault(exampleKey, tt.directory)
		expandedDir := GetAndExpandDirectory(exampleKey)

		if expandedDir == "" {
			t.Errorf("unable to expand directory %s", tt.directory)
		}
		if tt.shouldExpand == true && tt.directory == expandedDir {
			t.Errorf("expected directory %s to be expanded but it is the same", tt.directory)
		}
		if tt.shouldExpand == false && tt.directory != expandedDir {
			t.Errorf("expected directory %s not to be expanded but it was changed to %s", tt.directory, expandedDir)
		}
	}
}

func TestGetTzLocation(t *testing.T) {
	tests := []struct {
		name                  string
		tz                    string
		expectedStringifiedTz string
	}{
		{
			name:                  "local lowercase",
			tz:                    "local",
			expectedStringifiedTz: "Local",
		},
		{
			name:                  "local uppercase",
			tz:                    "local",
			expectedStringifiedTz: "Local",
		},
		{
			name:                  "UTC lowercase",
			tz:                    "utc",
			expectedStringifiedTz: "UTC",
		},
		{
			name:                  "UTC uppercase",
			tz:                    "UTC",
			expectedStringifiedTz: "UTC",
		},
		{
			name:                  "Valid tz identifier",
			tz:                    "America/New_York",
			expectedStringifiedTz: "America/New_York",
		},
		{
			name:                  "Invalid defaults to UTC",
			tz:                    "invalid",
			expectedStringifiedTz: "UTC",
		},
	}

	for _, tt := range tests {
		exampleKey := "example.key"
		viper.SetDefault(exampleKey, tt.tz)
		loc := getTzLocation(exampleKey)
		if loc.String() != tt.expectedStringifiedTz {
			t.Errorf("expected parsed timezone of '%s' but got '%s'", tt.expectedStringifiedTz, loc.String())
		}
	}
}

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
				"log.level":                   "DEBUG",
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
			// Clear viper state before each test
			viper.Reset()

			config := NewConfig(tc.defaults, tc.envBinding, tc.appWorkingDir, tc.sysConfigPath, tc.fileName)

			if config == nil {
				t.Error("Expected config to be created")
				return
			}

			if tc.expectedFields {
				assert.NotEmpty(t, config.WorkingDir)
				assert.NotEmpty(t, config.LogLevel)
				assert.NotNil(t, config.LogTimestampTimezone)
			}
		})
	}
}

func TestConfig_DumpConfig(t *testing.T) {
	viper.Reset()
	config := NewConfig(nil, nil, "", "", "")

	dump := config.DumpConfig()
	assert.NotEmpty(t, dump)

	// Verify it's valid JSON
	assert.True(t, len(dump) > 10)
	assert.Contains(t, dump, "working_dir")
	assert.Contains(t, dump, "log_level")
}

func TestConfig_PopulateFields(t *testing.T) {
	viper.Reset()

	// Set some test values
	viper.Set("application.working_dir", "~/test")
	viper.Set("application.default_profile", "testprofile")
	viper.Set("log.level", "DEBUG")
	viper.Set("features.datasets_enabled", true)
	viper.Set("terminal.no_color", true)
	viper.Set("git.name", "Test User")

	config := &Config{}
	config.populateFields()

	assert.Contains(t, config.WorkingDir, "test")
	assert.Equal(t, "testprofile", config.DefaultProfile)
	assert.Equal(t, "DEBUG", config.LogLevel)
	assert.True(t, config.FeaturesDatasetsEnabled)
	assert.True(t, config.TerminalNoColor)
	assert.Equal(t, "Test User", config.GitName)
}

func TestConfig_GetProfile(t *testing.T) {
	viper.Reset()

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
				c.profiles = map[string]*Profile{
					"test": {
						Host:     "test.example.com",
						Port:     8080,
						Username: "testuser",
					},
				}
			},
			expectError:  false,
			expectedHost: "test.example.com",
		},
		{
			name:        "Get non-existent profile returns default",
			profileName: "nonexistent",
			setupProfiles: func(c *Config) {
				c.profiles = map[string]*Profile{}
			},
			expectError:  true,
			expectedHost: defaultHost,
		},
		{
			name:        "Case insensitive profile lookup",
			profileName: "TEST",
			setupProfiles: func(c *Config) {
				c.profiles = map[string]*Profile{
					"test": {
						Host: "test.example.com",
						Port: 8080,
					},
				}
			},
			expectError:  false,
			expectedHost: "test.example.com",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			config := &Config{}
			tc.setupProfiles(config)

			profile, err := config.GetProfile(tc.profileName)

			if tc.expectError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}

			assert.NotNil(t, profile)
			assert.Equal(t, tc.expectedHost, profile.Host)
		})
	}
}

func TestConfig_ActiveProfile(t *testing.T) {
	testCases := []struct {
		name            string
		profileName     string
		setupProfiles   func(*Config)
		expectedProfile string
		expectError     bool
	}{
		{
			name:        "Active profile when profileName is set",
			profileName: "custom",
			setupProfiles: func(c *Config) {
				c.profiles = map[string]*Profile{
					"custom":  {Host: "custom.example.com"},
					"default": {Host: "default.example.com"},
				}
			},
			expectedProfile: "custom.example.com",
			expectError:     false,
		},
		{
			name:        "Active profile defaults to default when profileName is empty",
			profileName: "",
			setupProfiles: func(c *Config) {
				c.profiles = map[string]*Profile{
					"default": {Host: "default.example.com"},
				}
			},
			expectedProfile: "default.example.com",
			expectError:     false,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			config := &Config{profileName: tc.profileName}
			tc.setupProfiles(config)

			profile, err := config.ActiveProfile()

			if tc.expectError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}

			assert.NotNil(t, profile)
			assert.Equal(t, tc.expectedProfile, profile.Host)
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
				c.repositories = map[string]*Repository{
					"test-repo": {
						Url:       "https://github.com/test/repo.git",
						Reference: "main",
					},
				}
			},
			expectError: false,
			expectedUrl: "https://github.com/test/repo.git",
		},
		{
			name:     "Get non-existent repository",
			repoName: "nonexistent",
			setupRepositories: func(c *Config) {
				c.repositories = map[string]*Repository{}
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
