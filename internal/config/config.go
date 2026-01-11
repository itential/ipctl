// Copyright 2024 Itential Inc. All Rights Reserved
// Unauthorized copying of this file, via any medium is strictly prohibited
// Proprietary and confidential

package config

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/itential/ipctl/internal/profile"
	"github.com/itential/ipctl/internal/repository"
)

const (
	defaultFileName = "config"

	defaultAppWorkingDir        = "~/.platform.d"
	defaultAppDefaultProfile    = ""
	defaultAppDefaultRepository = ""

	defaultFeaturesDatasetsEnabled = false

	defaultGitName  = ""
	defaultGitEmail = ""
	defaultGitUser  = "git"
)

// Features holds feature flag configuration.
type Features struct {
	DatasetsEnabled bool `json:"datasets_enabled"`
}

// GitConfig holds Git-related configuration.
type GitConfig struct {
	Name  string `json:"name"`
	Email string `json:"email"`
	User  string `json:"user"`
}

// Config holds application-level configuration.
// It focuses on core application settings, feature flags, and profile/repository management.
// Domain-specific configuration (logging, terminal) should be managed by their respective packages.
type Config struct {
	// Application settings
	WorkingDir        string `json:"working_dir"`
	DefaultProfile    string `json:"default_profile"`
	DefaultRepository string `json:"default_repository"`

	// Feature flags
	Features Features `json:"features"`

	// Git settings
	Git GitConfig `json:"git"`

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
// Implements ApplicationProvider interface.
func (c *Config) GetWorkingDir() string {
	return c.WorkingDir
}

// GetDefaultProfile returns the name of the default profile.
// Implements ApplicationProvider interface.
func (c *Config) GetDefaultProfile() string {
	return c.DefaultProfile
}

// GetDefaultRepository returns the name of the default repository.
// Implements ApplicationProvider interface.
func (c *Config) GetDefaultRepository() string {
	return c.DefaultRepository
}

// IsDatasetsEnabled returns whether the datasets feature is enabled.
// Implements FeaturesProvider interface.
func (c *Config) IsDatasetsEnabled() bool {
	return c.Features.DatasetsEnabled
}

// GetGitName returns the git user name for commits.
// Implements GitProvider interface.
func (c *Config) GetGitName() string {
	return c.Git.Name
}

// GetGitEmail returns the git user email for commits.
// Implements GitProvider interface.
func (c *Config) GetGitEmail() string {
	return c.Git.Email
}

// GetGitUser returns the git username for authentication.
// Implements GitProvider interface.
func (c *Config) GetGitUser() string {
	return c.Git.User
}

// DumpConfig returns a JSON representation of the configuration.
// Useful for debugging and inspection.
func (ac *Config) DumpConfig() string {
	bs, _ := json.Marshal(ac)
	return string(bs)
}

var defaultValues = map[string]interface{}{
	"application.working_dir":        defaultAppWorkingDir,
	"application.default_profile":    defaultAppDefaultProfile,
	"application.default_repository": defaultAppDefaultRepository,

	"features.datasets_enabled": defaultFeaturesDatasetsEnabled,

	"git.name":  defaultGitName,
	"git.email": defaultGitEmail,
	"git.user":  defaultGitUser,
}

var defaultEnvVarBindings = map[string]string{
	"application.working_dir":        "IPCTL_APPLICATION_WORKING_DIR",
	"application.default_profile":    "IPCTL_APPLICATION_DEFAULT_PROFILE",
	"application.default_repository": "IPCTL_APPLICATION_DEFAULT_REPOSITORY",

	"features.datasets_enabled": "IPCTL_FEATURES_DATASETS_ENABLED",

	"git.name":  "IPCTL_GIT_NAME",
	"git.email": "IPCTL_GIT_EMAIL",
	"git.user":  "IPCTL_GIT_USER",
}
