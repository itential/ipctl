// Copyright 2024 Itential Inc. All Rights Reserved
// Unauthorized copying of this file, via any medium is strictly prohibited
// Proprietary and confidential

package config

import (
	"github.com/itential/ipctl/internal/profile"
	"github.com/itential/ipctl/internal/repository"
)

// ProfileProvider provides access to connection profiles.
// Use this interface when you need to:
// - Get a specific profile by name
// - Access the currently active profile
//
// Example:
//
//	func ConnectToServer(profiles ProfileProvider) error {
//	    profile, err := profiles.ActiveProfile()
//	    if err != nil {
//	        return err
//	    }
//	    // ... use profile ...
//	}
type ProfileProvider interface {
	// GetProfile returns a profile by name.
	// Returns an error if the profile doesn't exist.
	GetProfile(name string) (*profile.Profile, error)

	// ActiveProfile returns the currently active profile.
	// Returns an error if the active profile doesn't exist.
	ActiveProfile() (*profile.Profile, error)
}

// RepositoryProvider provides access to git repositories.
// Use this interface when you need to:
// - Get a specific repository by name
//
// Example:
//
//	func CloneRepository(repos RepositoryProvider) error {
//	    repo, err := repos.GetRepository("default")
//	    if err != nil {
//	        return err
//	    }
//	    // ... use repository ...
//	}
type RepositoryProvider interface {
	// GetRepository returns a repository by name.
	// Returns an error if the repository doesn't exist.
	GetRepository(name string) (*repository.Repository, error)
}

// ApplicationProvider provides application-level settings.
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

// FeaturesProvider provides feature flag access.
// Use this interface when you need to:
// - Check if specific features are enabled
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

// GitProvider provides git configuration.
// Use this interface when you need to:
// - Configure git commits with user information
// - Access git username for operations
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

// Provider is a combined interface providing all configuration capabilities.
// Use this interface when a component needs access to multiple configuration aspects.
//
// The Config type implements this interface, allowing it to be passed
// wherever any of the specific provider interfaces are needed.
//
// Example:
//
//	func Initialize(cfg Provider) error {
//	    // Can use cfg as ProfileProvider, ApplicationProvider, etc.
//	    profile, err := cfg.ActiveProfile()
//	    workDir := cfg.GetWorkingDir()
//	    // ...
//	}
//
// For better encapsulation and testability, prefer using specific interfaces
// (ProfileProvider, ApplicationProvider, etc.) instead of Provider when possible.
type Provider interface {
	ProfileProvider
	RepositoryProvider
	ApplicationProvider
	FeaturesProvider
	GitProvider
}
