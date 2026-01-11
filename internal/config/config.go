// Copyright 2024 Itential Inc. All Rights Reserved
// Unauthorized copying of this file, via any medium is strictly prohibited
// Proprietary and confidential

package config

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/itential/ipctl/internal/app"
	"github.com/itential/ipctl/internal/profile"
	"github.com/itential/ipctl/internal/repository"
)

const (
	defaultFileName = "config"
)

// Config holds configuration for the ipctl CLI application.
// It combines application-level settings (from app.Settings) with
// profile and repository management capabilities.
//
// Config serves as the central configuration provider for the application,
// delegating application-level concerns to app.Settings while managing
// connection profiles and repositories locally.
type Config struct {
	// Settings holds application-level configuration.
	// This includes working directory, default profile/repository, feature flags,
	// and git configuration.
	Settings *app.Settings

	// Managers for profiles and repositories
	profileManager    *profile.Manager
	repositoryManager *repository.Manager
}

// NewConfig creates a new Config instance using the provided settings.
// This function maintains backward compatibility with the legacy API.
// For new code, consider using NewLoader() for better testability and error handling.
//
// Parameters:
//   - defaults: Default values for configuration keys (uses defaultValues if nil)
//   - envBinding: Environment variable bindings (uses defaultEnvVarBindings if nil)
//   - appWorkingDir: User configuration directory (uses ~/.platform.d if empty)
//   - sysConfigPath: System configuration directory (uses /etc/ipctl if empty)
//   - fileName: Configuration file base name (uses "config" if empty)
//
// The function calls os.Exit(1) on error for backward compatibility.
// Use NewLoader() if you need proper error handling.
func NewConfig(defaults map[string]interface{}, envBinding map[string]string, appWorkingDir, sysConfigPath, fileName string) *Config {
	loader := NewLoader()

	if defaults != nil {
		loader.defaults = defaults
	}

	if envBinding != nil {
		loader.envBindings = envBinding
	}

	if appWorkingDir != "" {
		loader.workingDir = appWorkingDir
	}

	if sysConfigPath != "" {
		loader.sysConfigPath = sysConfigPath
	}

	if fileName != "" {
		loader.fileName = fileName
	}

	cfg, err := loader.Load()
	if err != nil {
		// For backward compatibility, exit on error
		fmt.Fprintf(os.Stderr, "Error: Failed to load configuration: %s\n", err)
		os.Exit(1)
	}

	return cfg
}

// Ensure Config implements Provider interface at compile time.
var _ Provider = (*Config)(nil)

// GetProfile returns a profile by name.
// If the profile doesn't exist, returns an error and a default profile.
// Implements ProfileProvider interface.
func (c *Config) GetProfile(name string) (*profile.Profile, error) {
	return c.profileManager.Get(name)
}

// ActiveProfile returns the currently active profile.
// Returns an error if the active profile doesn't exist.
// Implements ProfileProvider interface.
func (c *Config) ActiveProfile() (*profile.Profile, error) {
	return c.profileManager.Active()
}

// GetRepository returns a repository by name.
// Returns an error if the repository doesn't exist.
// Implements RepositoryProvider interface.
func (c *Config) GetRepository(name string) (*repository.Repository, error) {
	return c.repositoryManager.Get(name)
}

// GetWorkingDir returns the application working directory.
// Implements app.ApplicationProvider interface.
// Delegates to the embedded Settings.
func (c *Config) GetWorkingDir() string {
	return c.Settings.GetWorkingDir()
}

// GetDefaultProfile returns the name of the default profile.
// Implements app.ApplicationProvider interface.
// Delegates to the embedded Settings.
func (c *Config) GetDefaultProfile() string {
	return c.Settings.GetDefaultProfile()
}

// GetDefaultRepository returns the name of the default repository.
// Implements app.ApplicationProvider interface.
// Delegates to the embedded Settings.
func (c *Config) GetDefaultRepository() string {
	return c.Settings.GetDefaultRepository()
}

// IsDatasetsEnabled returns whether the datasets feature is enabled.
// Implements app.FeaturesProvider interface.
// Delegates to the embedded Settings.
func (c *Config) IsDatasetsEnabled() bool {
	return c.Settings.IsDatasetsEnabled()
}

// GetGitName returns the git user name for commits.
// Implements app.GitProvider interface.
// Delegates to the embedded Settings.
func (c *Config) GetGitName() string {
	return c.Settings.GetGitName()
}

// GetGitEmail returns the git user email for commits.
// Implements app.GitProvider interface.
// Delegates to the embedded Settings.
func (c *Config) GetGitEmail() string {
	return c.Settings.GetGitEmail()
}

// GetGitUser returns the git username for authentication.
// Implements app.GitProvider interface.
// Delegates to the embedded Settings.
func (c *Config) GetGitUser() string {
	return c.Settings.GetGitUser()
}

// DumpConfig returns a JSON representation of the configuration.
// Useful for debugging and inspection.
func (ac *Config) DumpConfig() string {
	bs, _ := json.Marshal(ac)
	return string(bs)
}

// defaultValues returns the default configuration values.
// This delegates to app.DefaultValues() for application-level settings.
var defaultValues = app.DefaultValues()

// defaultEnvVarBindings returns the environment variable bindings.
// This delegates to app.DefaultEnvBindings() for application-level settings.
var defaultEnvVarBindings = app.DefaultEnvBindings()
