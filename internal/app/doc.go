// Copyright 2024 Itential Inc. All Rights Reserved
// Unauthorized copying of this file, via any medium is strictly prohibited
// Proprietary and confidential

// Package app provides application-level metadata and information for ipctl.
//
// This package consolidates all application-level configuration and metadata,
// including version information, build details, and Git repository information.
// It serves as a single source of truth for application identity and metadata.
//
// # Version Information
//
// The package exposes version and build information that is set at compile time
// through linker flags:
//
//	go build -ldflags="-X 'github.com/itential/ipctl/internal/app.Version=v1.0.0' \
//	                    -X 'github.com/itential/ipctl/internal/app.Build=abc123'"
//
// Access version information using the Info type:
//
//	info := app.GetInfo()
//	fmt.Printf("Version: %s\n", info.Version)
//	fmt.Printf("Build: %s\n", info.Build)
//
// # Git Repository Information
//
// When running from source (not a released binary), the package can detect
// the current Git commit SHA:
//
//	sha, err := app.GetCurrentSha()
//	if err != nil {
//	    log.Printf("Not running from a git repository: %v", err)
//	} else {
//	    fmt.Printf("Running from commit: %s\n", sha)
//	}
//
// # Application Information
//
// The Info type provides a comprehensive view of application metadata:
//
//	type Info struct {
//	    Name    string  // Application name
//	    Version string  // Semantic version (e.g., "v1.2.3")
//	    Build   string  // Git commit SHA (short form)
//	}
//
// # Display Helpers
//
// The package provides helper methods for common display operations:
//
//	info := app.GetInfo()
//
//	// Display full version information
//	fmt.Println(info.FullVersion())  // "ipctl v1.2.3 (abc123)"
//
//	// Display short version
//	fmt.Println(info.ShortVersion()) // "v1.2.3"
//
//	// Check if version is set
//	if info.IsRelease() {
//	    fmt.Println("Running a release build")
//	} else {
//	    fmt.Println("Running a development build")
//	}
//
// # Integration with Build Systems
//
// Makefile integration:
//
//	GIT_COMMIT  := $(shell git rev-parse --short HEAD)
//	GIT_VERSION := $(shell git tag --sort=-v:refname | head -n 1)
//
//	build:
//	    go build -ldflags="-X 'github.com/itential/ipctl/internal/app.Build=$(GIT_COMMIT)' \
//	                       -X 'github.com/itential/ipctl/internal/app.Version=$(GIT_VERSION)'"
//
// Goreleaser integration:
//
//	builds:
//	  - ldflags:
//	      - -X github.com/itential/ipctl/internal/app.Build={{ .Env.BUILD }}
//	      - -X github.com/itential/ipctl/internal/app.Version={{ .Tag }}
//
// # Thread Safety
//
// All functions and methods in this package are safe for concurrent access.
// Version and build information is read-only after initialization and can be
// accessed from multiple goroutines without synchronization.
//
// # Error Handling
//
// Functions that interact with the filesystem (like GetCurrentSha) return
// errors that should be handled appropriately. The error messages provide
// context about what operation failed:
//
//	sha, err := app.GetCurrentSha()
//	if err != nil {
//	    if errors.Is(err, git.ErrRepositoryNotExists) {
//	        log.Println("Not in a git repository")
//	    } else {
//	        log.Printf("Failed to get git SHA: %v", err)
//	    }
//	}
//
// # Example Usage
//
// Complete example showing typical usage in a CLI application:
//
//	package main
//
//	import (
//	    "fmt"
//	    "github.com/itential/ipctl/internal/app"
//	)
//
//	func main() {
//	    info := app.GetInfo()
//
//	    if info.IsRelease() {
//	        // Running a released version
//	        fmt.Printf("%s %s (%s)\n", info.Name, info.Version, info.Build)
//	    } else {
//	        // Running from source
//	        sha, err := app.GetCurrentSha()
//	        if err == nil {
//	            fmt.Printf("%s running from commit %s\n", info.Name, sha)
//	        } else {
//	            fmt.Printf("%s (development build)\n", info.Name)
//	        }
//	    }
//	}
package app
