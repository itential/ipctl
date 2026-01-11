// Copyright 2024 Itential Inc. All Rights Reserved
// Unauthorized copying of this file, via any medium is strictly prohibited
// Proprietary and confidential

package app

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing/object"
)

func TestGetInfo(t *testing.T) {
	tests := []struct {
		name          string
		setupVersion  string
		setupBuild    string
		expectedName  string
		wantVersion   string
		wantBuild     string
		checkRelease  bool
		expectedShort string
		expectedFull  string
	}{
		{
			name:          "release build with version and build",
			setupVersion:  "v1.2.3",
			setupBuild:    "abc123",
			expectedName:  "ipctl",
			wantVersion:   "v1.2.3",
			wantBuild:     "abc123",
			checkRelease:  true,
			expectedShort: "v1.2.3",
			expectedFull:  "ipctl v1.2.3 (abc123)",
		},
		{
			name:          "development build with no version",
			setupVersion:  "",
			setupBuild:    "",
			expectedName:  "ipctl",
			wantVersion:   "",
			wantBuild:     "",
			checkRelease:  false,
			expectedShort: "development",
			expectedFull:  "ipctl development",
		},
		{
			name:          "partial build info - version only",
			setupVersion:  "v2.0.0",
			setupBuild:    "",
			expectedName:  "ipctl",
			wantVersion:   "v2.0.0",
			wantBuild:     "",
			checkRelease:  false,
			expectedShort: "v2.0.0",
			expectedFull:  "ipctl v2.0.0",
		},
		{
			name:          "partial build info - build only",
			setupVersion:  "",
			setupBuild:    "def456",
			expectedName:  "ipctl",
			wantVersion:   "",
			wantBuild:     "def456",
			checkRelease:  false,
			expectedShort: "development",
			expectedFull:  "ipctl development",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Save original values
			origVersion := Version
			origBuild := Build
			defer func() {
				Version = origVersion
				Build = origBuild
			}()

			// Set test values
			Version = tt.setupVersion
			Build = tt.setupBuild

			// Get info
			info := GetInfo()

			// Verify Info fields
			if info.Name != tt.expectedName {
				t.Errorf("GetInfo().Name = %q, want %q", info.Name, tt.expectedName)
			}
			if info.Version != tt.wantVersion {
				t.Errorf("GetInfo().Version = %q, want %q", info.Version, tt.wantVersion)
			}
			if info.Build != tt.wantBuild {
				t.Errorf("GetInfo().Build = %q, want %q", info.Build, tt.wantBuild)
			}

			// Verify IsRelease
			if got := info.IsRelease(); got != tt.checkRelease {
				t.Errorf("Info.IsRelease() = %v, want %v", got, tt.checkRelease)
			}

			// Verify ShortVersion
			if got := info.ShortVersion(); got != tt.expectedShort {
				t.Errorf("Info.ShortVersion() = %q, want %q", got, tt.expectedShort)
			}

			// Verify FullVersion
			if got := info.FullVersion(); got != tt.expectedFull {
				t.Errorf("Info.FullVersion() = %q, want %q", got, tt.expectedFull)
			}

			// Verify String() returns same as FullVersion()
			if got := info.String(); got != info.FullVersion() {
				t.Errorf("Info.String() = %q, want %q", got, info.FullVersion())
			}
		})
	}
}

func TestInfo_IsRelease(t *testing.T) {
	tests := []struct {
		name    string
		version string
		build   string
		want    bool
	}{
		{
			name:    "both version and build set",
			version: "v1.0.0",
			build:   "abc123",
			want:    true,
		},
		{
			name:    "only version set",
			version: "v1.0.0",
			build:   "",
			want:    false,
		},
		{
			name:    "only build set",
			version: "",
			build:   "abc123",
			want:    false,
		},
		{
			name:    "neither set",
			version: "",
			build:   "",
			want:    false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			info := Info{
				Name:    "ipctl",
				Version: tt.version,
				Build:   tt.build,
			}
			if got := info.IsRelease(); got != tt.want {
				t.Errorf("Info.IsRelease() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestInfo_ShortVersion(t *testing.T) {
	tests := []struct {
		name    string
		version string
		want    string
	}{
		{
			name:    "version set",
			version: "v1.2.3",
			want:    "v1.2.3",
		},
		{
			name:    "version empty",
			version: "",
			want:    "development",
		},
		{
			name:    "semantic version",
			version: "v2.0.0-beta.1",
			want:    "v2.0.0-beta.1",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			info := Info{
				Name:    "ipctl",
				Version: tt.version,
				Build:   "abc123",
			}
			if got := info.ShortVersion(); got != tt.want {
				t.Errorf("Info.ShortVersion() = %q, want %q", got, tt.want)
			}
		})
	}
}

func TestInfo_FullVersion(t *testing.T) {
	tests := []struct {
		name string
		info Info
		want string
	}{
		{
			name: "release build",
			info: Info{
				Name:    "ipctl",
				Version: "v1.2.3",
				Build:   "abc123",
			},
			want: "ipctl v1.2.3 (abc123)",
		},
		{
			name: "development build",
			info: Info{
				Name:    "ipctl",
				Version: "",
				Build:   "",
			},
			want: "ipctl development",
		},
		{
			name: "version only",
			info: Info{
				Name:    "ipctl",
				Version: "v2.0.0",
				Build:   "",
			},
			want: "ipctl v2.0.0",
		},
		{
			name: "build only",
			info: Info{
				Name:    "ipctl",
				Version: "",
				Build:   "def456",
			},
			want: "ipctl development",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.info.FullVersion(); got != tt.want {
				t.Errorf("Info.FullVersion() = %q, want %q", got, tt.want)
			}
		})
	}
}

func TestInfo_String(t *testing.T) {
	tests := []struct {
		name string
		info Info
		want string
	}{
		{
			name: "release build",
			info: Info{
				Name:    "ipctl",
				Version: "v1.0.0",
				Build:   "abc123",
			},
			want: "ipctl v1.0.0 (abc123)",
		},
		{
			name: "development build",
			info: Info{
				Name:    "ipctl",
				Version: "",
				Build:   "",
			},
			want: "ipctl development",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.info.String(); got != tt.want {
				t.Errorf("Info.String() = %q, want %q", got, tt.want)
			}
		})
	}
}

func TestGetCurrentSha(t *testing.T) {
	tests := []struct {
		name        string
		setup       func(t *testing.T) string // Returns temp dir path
		wantErr     bool
		errContains string
		validateSHA func(t *testing.T, sha string)
	}{
		{
			name: "valid git repository",
			setup: func(t *testing.T) string {
				// Create a temporary directory with a git repository
				tempDir := t.TempDir()
				repo, err := git.PlainInit(tempDir, false)
				if err != nil {
					t.Fatalf("Failed to init git repo: %v", err)
				}

				// Create a test file
				testFile := filepath.Join(tempDir, "test.txt")
				if err := os.WriteFile(testFile, []byte("test"), 0644); err != nil {
					t.Fatalf("Failed to create test file: %v", err)
				}

				// Add and commit the file
				worktree, err := repo.Worktree()
				if err != nil {
					t.Fatalf("Failed to get worktree: %v", err)
				}

				if _, err := worktree.Add("test.txt"); err != nil {
					t.Fatalf("Failed to add file: %v", err)
				}

				_, err = worktree.Commit("Initial commit", &git.CommitOptions{
					Author: &object.Signature{
						Name:  "Test User",
						Email: "test@example.com",
					},
				})
				if err != nil {
					t.Fatalf("Failed to commit: %v", err)
				}

				return tempDir
			},
			wantErr: false,
			validateSHA: func(t *testing.T, sha string) {
				// SHA should be 40 hex characters
				if len(sha) != 40 {
					t.Errorf("Expected SHA length of 40, got %d", len(sha))
				}
				// Verify it's all hex characters
				for _, c := range sha {
					if !strings.ContainsRune("0123456789abcdef", c) {
						t.Errorf("SHA contains non-hex character: %c", c)
						break
					}
				}
			},
		},
		{
			name: "not a git repository",
			setup: func(t *testing.T) string {
				// Create a temporary directory without git
				return t.TempDir()
			},
			wantErr:     true,
			errContains: "failed to open git repository",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Setup test environment
			testDir := tt.setup(t)

			// Save current directory
			origDir, err := os.Getwd()
			if err != nil {
				t.Fatalf("Failed to get current directory: %v", err)
			}
			defer os.Chdir(origDir)

			// Change to test directory
			if err := os.Chdir(testDir); err != nil {
				t.Fatalf("Failed to change to test directory: %v", err)
			}

			// Execute test
			sha, err := GetCurrentSha()

			// Validate error expectations
			if tt.wantErr {
				if err == nil {
					t.Errorf("GetCurrentSha() expected error, got nil")
				} else if tt.errContains != "" && !strings.Contains(err.Error(), tt.errContains) {
					t.Errorf("GetCurrentSha() error = %q, want error containing %q", err.Error(), tt.errContains)
				}
				return
			}

			// Validate success expectations
			if err != nil {
				t.Errorf("GetCurrentSha() unexpected error: %v", err)
				return
			}

			if tt.validateSHA != nil {
				tt.validateSHA(t, sha)
			}
		})
	}
}

func TestGetCurrentSha_NotInGitRepo(t *testing.T) {
	// Create a temp directory that's not a git repository
	tempDir := t.TempDir()

	// Save current directory
	origDir, err := os.Getwd()
	if err != nil {
		t.Fatalf("Failed to get current directory: %v", err)
	}
	defer os.Chdir(origDir)

	// Change to temp directory
	if err := os.Chdir(tempDir); err != nil {
		t.Fatalf("Failed to change directory: %v", err)
	}

	// Attempt to get SHA
	sha, err := GetCurrentSha()

	// Should return an error
	if err == nil {
		t.Errorf("GetCurrentSha() in non-git directory expected error, got nil")
	}

	// SHA should be empty
	if sha != "" {
		t.Errorf("GetCurrentSha() returned SHA %q, want empty string on error", sha)
	}

	// Error should mention repository
	if !strings.Contains(err.Error(), "repository") {
		t.Errorf("GetCurrentSha() error = %q, want error mentioning 'repository'", err.Error())
	}
}

// TestConstants verifies that package constants have expected values.
func TestConstants(t *testing.T) {
	tests := []struct {
		name string
		got  string
		want string
	}{
		{
			name: "application name",
			got:  Name,
			want: "ipctl",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.got != tt.want {
				t.Errorf("constant value = %q, want %q", tt.got, tt.want)
			}
		})
	}
}

// TestPackageVariables verifies that package variables can be set correctly.
// This simulates what happens during build with linker flags.
func TestPackageVariables(t *testing.T) {
	// Save original values
	origVersion := Version
	origBuild := Build
	defer func() {
		Version = origVersion
		Build = origBuild
	}()

	tests := []struct {
		name        string
		setVersion  string
		setBuild    string
		wantVersion string
		wantBuild   string
	}{
		{
			name:        "set version and build",
			setVersion:  "v1.2.3",
			setBuild:    "abc123",
			wantVersion: "v1.2.3",
			wantBuild:   "abc123",
		},
		{
			name:        "set to empty strings",
			setVersion:  "",
			setBuild:    "",
			wantVersion: "",
			wantBuild:   "",
		},
		{
			name:        "set long version string",
			setVersion:  "v2.0.0-beta.1+20230101",
			setBuild:    "1234567890abcdef",
			wantVersion: "v2.0.0-beta.1+20230101",
			wantBuild:   "1234567890abcdef",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Set package variables
			Version = tt.setVersion
			Build = tt.setBuild

			// Verify they were set correctly
			if Version != tt.wantVersion {
				t.Errorf("Version = %q, want %q", Version, tt.wantVersion)
			}
			if Build != tt.wantBuild {
				t.Errorf("Build = %q, want %q", Build, tt.wantBuild)
			}

			// Verify GetInfo reflects the changes
			info := GetInfo()
			if info.Version != tt.wantVersion {
				t.Errorf("GetInfo().Version = %q, want %q", info.Version, tt.wantVersion)
			}
			if info.Build != tt.wantBuild {
				t.Errorf("GetInfo().Build = %q, want %q", info.Build, tt.wantBuild)
			}
		})
	}
}
