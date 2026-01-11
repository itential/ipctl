// Copyright 2024 Itential Inc. All Rights Reserved
// Unauthorized copying of this file, via any medium is strictly prohibited
// Proprietary and confidential

package app

import (
	"fmt"
	"os"

	"github.com/go-git/go-git/v5"
)

const (
	// Name is the application name.
	Name = "ipctl"
)

// Build represents the Git SHA (short form) that the binary was compiled from.
// This value is set at build time using linker flags:
//
//	-ldflags "-X github.com/itential/ipctl/internal/app.Build=<commit-sha>"
var Build string

// Version represents the semantic version of the binary.
// This value is set at build time using linker flags:
//
//	-ldflags "-X github.com/itential/ipctl/internal/app.Version=<version>"
var Version string

// Info contains application metadata including name, version, and build information.
// Use GetInfo() to obtain an Info instance populated with current build metadata.
type Info struct {
	// Name is the application name (always "ipctl").
	Name string

	// Version is the semantic version string (e.g., "v1.2.3").
	// Empty string indicates a development build.
	Version string

	// Build is the Git commit SHA (short form) the binary was built from.
	// Empty string indicates version information is not available.
	Build string
}

// GetInfo returns application information including name, version, and build details.
// The returned Info is populated with values from the Build and Version package variables,
// which are typically set at compile time via linker flags.
func GetInfo() Info {
	return Info{
		Name:    Name,
		Version: Version,
		Build:   Build,
	}
}

// IsRelease returns true if the application is running a release build.
// A release build is one where both Version and Build are set (non-empty).
// Returns false for development builds where version information is not available.
func (i Info) IsRelease() bool {
	return i.Version != "" && i.Build != ""
}

// ShortVersion returns just the version string without build information.
// Returns "development" if version is not set.
func (i Info) ShortVersion() string {
	if i.Version == "" {
		return "development"
	}
	return i.Version
}

// FullVersion returns a complete version string including build information.
// Format: "ipctl v1.2.3 (abc123)" for release builds.
// Format: "ipctl development" for development builds.
func (i Info) FullVersion() string {
	if i.IsRelease() {
		return fmt.Sprintf("%s %s (%s)", i.Name, i.Version, i.Build)
	}
	return fmt.Sprintf("%s %s", i.Name, i.ShortVersion())
}

// String returns a human-readable representation of the application info.
// This is equivalent to calling FullVersion().
func (i Info) String() string {
	return i.FullVersion()
}

// GetCurrentSha retrieves the current Git commit SHA from the repository
// in the current working directory. This is useful when running from source
// to display the exact commit being executed.
//
// Returns the full SHA hash as a string, or an error if:
//   - The current directory is not part of a Git repository
//   - Unable to read the current working directory
//   - Unable to access HEAD reference
//   - Unable to read the commit object
//
// Example:
//
//	sha, err := app.GetCurrentSha()
//	if err != nil {
//	    log.Printf("Not in a git repository: %v", err)
//	    return
//	}
//	fmt.Printf("Running from commit: %s\n", sha)
func GetCurrentSha() (string, error) {
	cwd, err := os.Getwd()
	if err != nil {
		return "", fmt.Errorf("failed to get current working directory: %w", err)
	}

	// Open the Git repository in the current directory
	repo, err := git.PlainOpen(cwd)
	if err != nil {
		return "", fmt.Errorf("failed to open git repository: %w", err)
	}

	// Get the HEAD reference
	head, err := repo.Head()
	if err != nil {
		return "", fmt.Errorf("failed to get HEAD reference: %w", err)
	}

	// Resolve the commit object
	commit, err := repo.CommitObject(head.Hash())
	if err != nil {
		return "", fmt.Errorf("failed to get commit object: %w", err)
	}

	// Return the full SHA hash
	return commit.Hash.String(), nil
}
