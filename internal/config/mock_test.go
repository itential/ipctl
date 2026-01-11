// Copyright 2024 Itential Inc. All Rights Reserved
// Unauthorized copying of this file, via any medium is strictly prohibited
// Proprietary and confidential

package config_test

import (
	"fmt"

	"github.com/itential/ipctl/internal/config"
	"github.com/itential/ipctl/internal/profile"
	"github.com/itential/ipctl/internal/repository"
)

// mockProvider is a test implementation of config.Provider interface.
// It provides a simple way to mock configuration for testing without
// requiring file I/O or complex setup.
//
// Example usage:
//
//	mock := newMockProvider()
//	mock.profiles["test"] = profile.New("localhost", 8080, false, "admin", "password", 30)
//	mock.defaultProfile = "test"
//
//	profile, err := mock.ActiveProfile()
//	// use profile in tests...
type mockProvider struct {
	workingDir        string
	defaultProfile    string
	defaultRepository string
	datasetsEnabled   bool
	gitName           string
	gitEmail          string
	gitUser           string
	profiles          map[string]*profile.Profile
	repositories      map[string]*repository.Repository
}

// newMockProvider creates a new mockProvider with sensible defaults.
func newMockProvider() *mockProvider {
	return &mockProvider{
		workingDir:        "/tmp/test",
		defaultProfile:    "default",
		defaultRepository: "default",
		datasetsEnabled:   false,
		gitName:           "Test User",
		gitEmail:          "test@example.com",
		gitUser:           "git",
		profiles:          make(map[string]*profile.Profile),
		repositories:      make(map[string]*repository.Repository),
	}
}

// GetWorkingDir implements ApplicationProvider interface.
func (m *mockProvider) GetWorkingDir() string {
	return m.workingDir
}

// GetDefaultProfile implements ApplicationProvider interface.
func (m *mockProvider) GetDefaultProfile() string {
	return m.defaultProfile
}

// GetDefaultRepository implements ApplicationProvider interface.
func (m *mockProvider) GetDefaultRepository() string {
	return m.defaultRepository
}

// IsDatasetsEnabled implements FeaturesProvider interface.
func (m *mockProvider) IsDatasetsEnabled() bool {
	return m.datasetsEnabled
}

// GetGitName implements GitProvider interface.
func (m *mockProvider) GetGitName() string {
	return m.gitName
}

// GetGitEmail implements GitProvider interface.
func (m *mockProvider) GetGitEmail() string {
	return m.gitEmail
}

// GetGitUser implements GitProvider interface.
func (m *mockProvider) GetGitUser() string {
	return m.gitUser
}

// GetProfile implements ProfileProvider interface.
func (m *mockProvider) GetProfile(name string) (*profile.Profile, error) {
	if p, ok := m.profiles[name]; ok {
		return p, nil
	}
	// Return default profile similar to real implementation
	return profile.Default(), fmt.Errorf("profile %q not found", name)
}

// ActiveProfile implements ProfileProvider interface.
func (m *mockProvider) ActiveProfile() (*profile.Profile, error) {
	return m.GetProfile(m.defaultProfile)
}

// GetRepository implements RepositoryProvider interface.
func (m *mockProvider) GetRepository(name string) (*repository.Repository, error) {
	if r, ok := m.repositories[name]; ok {
		return r, nil
	}
	return nil, fmt.Errorf("repository %q not found", name)
}

// Ensure mockProvider implements Provider interface at compile time.
var _ config.Provider = (*mockProvider)(nil)
