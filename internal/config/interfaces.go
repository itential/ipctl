// Copyright 2024 Itential Inc. All Rights Reserved
// Unauthorized copying of this file, via any medium is strictly prohibited
// Proprietary and confidential

package config

import (
	"github.com/itential/ipctl/internal/app"
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

// Provider is a combined interface providing all configuration capabilities.
// It combines profile/repository providers from this package with application-level
// providers from the app package.
//
// The Config type implements this interface, allowing it to be passed
// wherever any of the specific provider interfaces are needed.
//
// Example:
//
//	func Initialize(cfg Provider) error {
//	    // Can use cfg as ProfileProvider, app.ApplicationProvider, etc.
//	    profile, err := cfg.ActiveProfile()
//	    workDir := cfg.GetWorkingDir()
//	    // ...
//	}
//
// For better encapsulation and testability, prefer using specific interfaces
// (ProfileProvider, app.ApplicationProvider, etc.) instead of Provider when possible.
type Provider interface {
	ProfileProvider
	RepositoryProvider
	app.ApplicationProvider
	app.FeaturesProvider
	app.GitProvider
}
