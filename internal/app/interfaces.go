// Copyright 2024 Itential Inc. All Rights Reserved
// Unauthorized copying of this file, via any medium is strictly prohibited
// Proprietary and confidential

package app

// ApplicationProvider provides access to application-level settings.
// Use this interface when you need to:
// - Access the working directory
// - Get default profile or repository names
//
// Example:
//
//	func InitializeApp(app ApplicationProvider) error {
//	    workDir := app.GetWorkingDir()
//	    // ... initialize application ...
//	}
type ApplicationProvider interface {
	// GetWorkingDir returns the application working directory.
	GetWorkingDir() string

	// GetDefaultProfile returns the name of the default profile.
	GetDefaultProfile() string

	// GetDefaultRepository returns the name of the default repository.
	GetDefaultRepository() string
}

// FeaturesProvider provides access to feature flags.
// Use this interface when you need to:
// - Check if specific features are enabled
//
// Feature flags allow enabling/disabling functionality at runtime,
// useful for gradual rollouts and managing experimental features.
//
// Example:
//
//	func ShowDatasets(features FeaturesProvider) bool {
//	    return features.IsDatasetsEnabled()
//	}
type FeaturesProvider interface {
	// IsDatasetsEnabled returns whether the datasets feature is enabled.
	IsDatasetsEnabled() bool
}

// GitProvider provides access to git configuration.
// Use this interface when you need to:
// - Configure git commits with user information
// - Access git username for authentication operations
//
// Example:
//
//	func ConfigureGit(git GitProvider) error {
//	    name := git.GetGitName()
//	    email := git.GetGitEmail()
//	    // ... configure git ...
//	}
type GitProvider interface {
	// GetGitName returns the git user name for commits.
	GetGitName() string

	// GetGitEmail returns the git user email for commits.
	GetGitEmail() string

	// GetGitUser returns the git username for authentication.
	GetGitUser() string
}

// SettingsProvider is a combined interface providing all application-level
// configuration capabilities. It groups ApplicationProvider, FeaturesProvider,
// and GitProvider into a single interface.
//
// The Settings type implements this interface, allowing it to be passed
// wherever any of the specific provider interfaces are needed.
//
// Example:
//
//	func Initialize(settings SettingsProvider) error {
//	    // Can use settings as ApplicationProvider, FeaturesProvider, or GitProvider
//	    workDir := settings.GetWorkingDir()
//	    datasetsEnabled := settings.IsDatasetsEnabled()
//	    gitName := settings.GetGitName()
//	    // ...
//	}
//
// For better encapsulation and testability, prefer using specific interfaces
// (ApplicationProvider, FeaturesProvider, GitProvider) instead of SettingsProvider
// when possible. This makes dependencies explicit and components easier to test.
type SettingsProvider interface {
	ApplicationProvider
	FeaturesProvider
	GitProvider
}

// Ensure Settings implements SettingsProvider interface at compile time.
var _ SettingsProvider = (*Settings)(nil)
