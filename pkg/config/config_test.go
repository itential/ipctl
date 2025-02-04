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
