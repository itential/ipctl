// Copyright 2024 Itential Inc. All Rights Reserved
// Unauthorized copying of this file, via any medium is strictly prohibited
// Proprietary and confidential

package app

const (
	// Default application settings
	defaultWorkingDir        = "~/.platform.d"
	defaultDefaultProfile    = ""
	defaultDefaultRepository = ""

	// Default feature flags
	defaultDatasetsEnabled = false

	// Default git settings
	defaultGitName  = ""
	defaultGitEmail = ""
	defaultGitUser  = "git"
)

// Settings holds application-level configuration.
// This includes core application settings, feature flags, and git configuration.
// Settings is focused on application behavior and does not manage connection
// profiles or repositories (those are handled by internal/config).
//
// Settings values are typically loaded from configuration files, environment
// variables, or command-line flags. The precedence order is:
//  1. Command-line flags (highest)
//  2. Environment variables
//  3. Configuration file
//  4. Default values (lowest)
type Settings struct {
	// WorkingDir is the application's working directory.
	// This is where profiles, repositories, and other application data are stored.
	// Defaults to ~/.platform.d
	WorkingDir string

	// DefaultProfile is the name of the profile to use when none is specified.
	// Empty string means no default profile is configured.
	DefaultProfile string

	// DefaultRepository is the name of the repository to use when none is specified.
	// Empty string means no default repository is configured.
	DefaultRepository string

	// Features contains feature flag settings.
	Features Features

	// Git contains git-related configuration.
	Git GitConfig
}

// Features holds feature flag configuration.
// Feature flags allow enabling/disabling functionality at runtime without
// code changes. They're useful for gradual rollouts, A/B testing, and
// managing experimental features.
type Features struct {
	// DatasetsEnabled controls whether dataset management features are available.
	// When false, dataset commands are hidden from the CLI.
	DatasetsEnabled bool
}

// GitConfig holds Git-related configuration.
// These settings are used when the application interacts with git repositories
// for import/export operations or when committing changes.
type GitConfig struct {
	// Name is the git user name used for commits.
	// This corresponds to git config user.name
	Name string

	// Email is the git user email used for commits.
	// This corresponds to git config user.email
	Email string

	// User is the git username for authentication.
	// This is used for git operations that require authentication.
	User string
}

// NewSettings creates a new Settings instance with default values.
// Use SettingsBuilder for more control over initialization.
func NewSettings() *Settings {
	return &Settings{
		WorkingDir:        defaultWorkingDir,
		DefaultProfile:    defaultDefaultProfile,
		DefaultRepository: defaultDefaultRepository,
		Features: Features{
			DatasetsEnabled: defaultDatasetsEnabled,
		},
		Git: GitConfig{
			Name:  defaultGitName,
			Email: defaultGitEmail,
			User:  defaultGitUser,
		},
	}
}

// GetWorkingDir returns the application working directory.
// Implements ApplicationProvider interface.
func (s *Settings) GetWorkingDir() string {
	return s.WorkingDir
}

// GetDefaultProfile returns the name of the default profile.
// Implements ApplicationProvider interface.
func (s *Settings) GetDefaultProfile() string {
	return s.DefaultProfile
}

// GetDefaultRepository returns the name of the default repository.
// Implements ApplicationProvider interface.
func (s *Settings) GetDefaultRepository() string {
	return s.DefaultRepository
}

// IsDatasetsEnabled returns whether the datasets feature is enabled.
// Implements FeaturesProvider interface.
func (s *Settings) IsDatasetsEnabled() bool {
	return s.Features.DatasetsEnabled
}

// GetGitName returns the git user name for commits.
// Implements GitProvider interface.
func (s *Settings) GetGitName() string {
	return s.Git.Name
}

// GetGitEmail returns the git user email for commits.
// Implements GitProvider interface.
func (s *Settings) GetGitEmail() string {
	return s.Git.Email
}

// GetGitUser returns the git username for authentication.
// Implements GitProvider interface.
func (s *Settings) GetGitUser() string {
	return s.Git.User
}

// DefaultValues returns a map of default values for all settings.
// This map uses dotted key notation (e.g., "application.working_dir")
// and is suitable for use with configuration loaders like viper.
func DefaultValues() map[string]interface{} {
	return map[string]interface{}{
		"application.working_dir":        defaultWorkingDir,
		"application.default_profile":    defaultDefaultProfile,
		"application.default_repository": defaultDefaultRepository,

		"features.datasets_enabled": defaultDatasetsEnabled,

		"git.name":  defaultGitName,
		"git.email": defaultGitEmail,
		"git.user":  defaultGitUser,
	}
}

// DefaultEnvBindings returns a map of environment variable bindings for all settings.
// This map uses dotted key notation for keys and IPCTL_* prefixed environment
// variable names for values.
func DefaultEnvBindings() map[string]string {
	return map[string]string{
		"application.working_dir":        "IPCTL_APPLICATION_WORKING_DIR",
		"application.default_profile":    "IPCTL_APPLICATION_DEFAULT_PROFILE",
		"application.default_repository": "IPCTL_APPLICATION_DEFAULT_REPOSITORY",

		"features.datasets_enabled": "IPCTL_FEATURES_DATASETS_ENABLED",

		"git.name":  "IPCTL_GIT_NAME",
		"git.email": "IPCTL_GIT_EMAIL",
		"git.user":  "IPCTL_GIT_USER",
	}
}
