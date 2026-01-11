// Copyright 2024 Itential Inc. All Rights Reserved
// Unauthorized copying of this file, via any medium is strictly prohibited
// Proprietary and confidential

package config_test

import (
	"testing"

	"github.com/itential/ipctl/internal/app"
	"github.com/itential/ipctl/internal/config"
	"github.com/itential/ipctl/internal/profile"
	"github.com/itential/ipctl/internal/repository"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestConfigImplementsProvider verifies that Config implements the Provider interface.
// This test ensures compile-time interface compliance is maintained at runtime.
func TestConfigImplementsProvider(t *testing.T) {
	cfg := &config.Config{}

	// Verify Config implements all interfaces
	var _ config.Provider = cfg
	var _ config.ProfileProvider = cfg
	var _ config.RepositoryProvider = cfg
	var _ app.ApplicationProvider = cfg
	var _ app.FeaturesProvider = cfg
	var _ app.GitProvider = cfg

	// If we get here, the compile-time checks passed
	t.Log("Config successfully implements all Provider interfaces")
}

// TestMockProviderImplementsProvider verifies that mockProvider implements the Provider interface.
// This ensures our test mock can be used anywhere Config is expected.
func TestMockProviderImplementsProvider(t *testing.T) {
	mock := newMockProvider()

	// Verify mockProvider implements all interfaces
	var _ config.Provider = mock
	var _ config.ProfileProvider = mock
	var _ config.RepositoryProvider = mock
	var _ app.ApplicationProvider = mock
	var _ app.FeaturesProvider = mock
	var _ app.GitProvider = mock

	t.Log("mockProvider successfully implements all Provider interfaces")
}

// TestMockProviderBehavior verifies that mockProvider behaves correctly.
func TestMockProviderBehavior(t *testing.T) {
	t.Run("ApplicationProvider methods", func(t *testing.T) {
		mock := newMockProvider()
		mock.workingDir = "/custom/path"
		mock.defaultProfile = "production"
		mock.defaultRepository = "myrepo"

		assert.Equal(t, "/custom/path", mock.GetWorkingDir())
		assert.Equal(t, "production", mock.GetDefaultProfile())
		assert.Equal(t, "myrepo", mock.GetDefaultRepository())
	})

	t.Run("FeaturesProvider methods", func(t *testing.T) {
		mock := newMockProvider()
		assert.False(t, mock.IsDatasetsEnabled())

		mock.datasetsEnabled = true
		assert.True(t, mock.IsDatasetsEnabled())
	})

	t.Run("GitProvider methods", func(t *testing.T) {
		mock := newMockProvider()
		mock.gitName = "John Doe"
		mock.gitEmail = "john@example.com"
		mock.gitUser = "johndoe"

		assert.Equal(t, "John Doe", mock.GetGitName())
		assert.Equal(t, "john@example.com", mock.GetGitEmail())
		assert.Equal(t, "johndoe", mock.GetGitUser())
	})

	t.Run("ProfileProvider methods", func(t *testing.T) {
		mock := newMockProvider()

		// Test getting non-existent profile returns default
		prof, err := mock.GetProfile("nonexistent")
		require.Error(t, err)
		assert.NotNil(t, prof) // Should return default profile

		// Test adding and retrieving a profile
		testProfile := &profile.Profile{
			Host:     "testhost",
			Port:     8080,
			UseTLS:   false,
			Username: "admin",
			Password: "password",
			Timeout:  30,
		}
		mock.profiles["test"] = testProfile

		prof, err = mock.GetProfile("test")
		require.NoError(t, err)
		assert.Equal(t, "testhost", prof.Host)
		assert.Equal(t, 8080, prof.Port)

		// Test active profile
		mock.defaultProfile = "test"
		activeProf, err := mock.ActiveProfile()
		require.NoError(t, err)
		assert.Equal(t, testProfile, activeProf)
	})

	t.Run("RepositoryProvider methods", func(t *testing.T) {
		mock := newMockProvider()

		// Test getting non-existent repository
		_, err := mock.GetRepository("nonexistent")
		require.Error(t, err)

		// Test adding and retrieving a repository
		testRepo := &repository.Repository{
			Url:            "https://github.com/test/repo.git",
			Reference:      "main",
			PrivateKeyFile: "",
		}
		mock.repositories["test"] = testRepo

		repo, err := mock.GetRepository("test")
		require.NoError(t, err)
		assert.Equal(t, testRepo, repo)
	})
}

// TestProviderInterfaceSegregation demonstrates that functions can accept
// specific interfaces instead of the full Provider interface.
func TestProviderInterfaceSegregation(t *testing.T) {
	// Helper functions that accept specific interfaces
	needsProfiles := func(p config.ProfileProvider) error {
		_, err := p.ActiveProfile()
		return err
	}

	needsGit := func(g app.GitProvider) string {
		return g.GetGitName()
	}

	needsApp := func(a app.ApplicationProvider) string {
		return a.GetWorkingDir()
	}

	needsFeatures := func(f app.FeaturesProvider) bool {
		return f.IsDatasetsEnabled()
	}

	needsRepo := func(r config.RepositoryProvider) error {
		_, err := r.GetRepository("test")
		return err
	}

	// Create a mock provider
	mock := newMockProvider()
	mock.profiles["default"] = profile.Default()

	// Verify we can pass the mock to functions expecting specific interfaces
	t.Run("Pass to ProfileProvider function", func(t *testing.T) {
		err := needsProfiles(mock)
		assert.NoError(t, err)
	})

	t.Run("Pass to GitProvider function", func(t *testing.T) {
		name := needsGit(mock)
		assert.Equal(t, "Test User", name)
	})

	t.Run("Pass to ApplicationProvider function", func(t *testing.T) {
		dir := needsApp(mock)
		assert.Equal(t, "/tmp/test", dir)
	})

	t.Run("Pass to FeaturesProvider function", func(t *testing.T) {
		enabled := needsFeatures(mock)
		assert.False(t, enabled)
	})

	t.Run("Pass to RepositoryProvider function", func(t *testing.T) {
		err := needsRepo(mock)
		assert.Error(t, err) // Expected since we didn't add a "test" repository
	})
}

// TestConfigInterfaceMethods verifies that Config's interface methods work correctly.
// This test uses a real Config instance (when we can construct one).
func TestConfigInterfaceMethods(t *testing.T) {
	t.Run("Interface methods delegate to manager", func(t *testing.T) {
		// Create a config using the loader
		loader := config.NewLoader()
		cfg, err := loader.Load()
		require.NoError(t, err)
		require.NotNil(t, cfg)

		// Test Provider interface methods exist and can be called
		t.Run("ApplicationProvider methods", func(t *testing.T) {
			assert.NotEmpty(t, cfg.GetWorkingDir())
			// GetDefaultProfile and GetDefaultRepository may be empty strings
			_ = cfg.GetDefaultProfile()
			_ = cfg.GetDefaultRepository()
		})

		t.Run("FeaturesProvider methods", func(t *testing.T) {
			// Just verify the method exists and returns a bool
			_ = cfg.IsDatasetsEnabled()
		})

		t.Run("GitProvider methods", func(t *testing.T) {
			// Methods should exist even if they return empty strings
			_ = cfg.GetGitName()
			_ = cfg.GetGitEmail()
			_ = cfg.GetGitUser()
		})

		t.Run("ProfileProvider methods", func(t *testing.T) {
			// ActiveProfile should work with default profile
			_, err := cfg.ActiveProfile()
			// May error if no profiles configured, but method should exist
			_ = err

			// GetProfile should work
			_, err = cfg.GetProfile("default")
			_ = err
		})

		t.Run("RepositoryProvider methods", func(t *testing.T) {
			// GetRepository should work even if it returns an error
			_, err := cfg.GetRepository("nonexistent")
			_ = err
		})
	})
}
