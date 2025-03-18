// Copyright 2024 Itential Inc. All Rights Reserved
// Unauthorized copying of this file, via any medium is strictly prohibited
// Proprietary and confidential

package metadata

import (
	"fmt"
	"os"

	"github.com/go-git/go-git/v5"
)

// Represents the Git SHA (short) the build was compiled against
var Build string

// Represent the version of the build in the binary
var Version string

func GetCurrentSha() (string, error) {
	cwd, err := os.Getwd()
	if err != nil {
		return "", err
	}

	// Open the Git repository
	repo, err := git.PlainOpen(cwd)
	if err != nil {
		return "", fmt.Errorf("failed to open repository: %w", err)
	}

	// Get the HEAD reference
	head, err := repo.Head()
	if err != nil {
		return "", fmt.Errorf("failed to get HEAD: %w", err)
	}

	// Resolve the commit
	commit, err := repo.CommitObject(head.Hash())
	if err != nil {
		return "", fmt.Errorf("failed to get commit object: %w", err)
	}

	// Extract the SHA
	sha := commit.Hash.String()
	return sha, nil
}
