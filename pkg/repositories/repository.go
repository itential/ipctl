// Copyright 2024 Itential Inc. All Rights Reserved
// Unauthorized copying of this file, via any medium is strictly prohibited
// Proprietary and confidential

package repositories

import (
	"fmt"
	"net/url"
	"os"

	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing"
	"github.com/itential/ipctl/pkg/logger"
	giturls "github.com/whilp/git-urls"
)

type Repository struct {
	Url        string
	Reference  string
	PrivateKey string
}

func (r Repository) Clone() (string, error) {
	logger.Trace()

	target, err := os.MkdirTemp("", "tmp")
	if err != nil {
		logger.Fatal(err, "failed to create temp dir")
	}
	logger.Info("temporary folder is %s", target)
	defer os.Remove(target)

	logger.Debug("source repository url is %s", r.Url)
	logger.Debug("source reference is %s", r.Reference)

	cloneOptions := &git.CloneOptions{
		URL: r.Url,
		//Auth: r.PrivateKey,
	}

	// Git will default to main/master if no repo is specified
	if r.Reference != "" {
		cloneOptions.ReferenceName = plumbing.NewBranchReferenceName(r.Reference)
	}

	u, err := parseRepositoryUrl(r.Url)
	if err != nil {
		return "", err
	}
	if u == nil {
		return "", fmt.Errorf("invalid repository url: %s", r.Url)
	}
	logger.Debug("uri schema is %s", u.Scheme)

	res, err := git.PlainClone(target, false, cloneOptions)
	if err != nil {
		return "", fmt.Errorf("failed to clone the repository: %s", err)
	}

	ref, err := res.Head()
	if err != nil {
		return "", err
	}

	_, err = res.CommitObject(ref.Hash())
	if err != nil {
		return "", err
	}

	logger.Debug("clone repository completed successfully to %v", target)
	logger.Debug("target is %s", target)

	return target, nil
}

func parseRepositoryUrl(uri string) (*url.URL, error) {
	u, err := giturls.Parse(uri)
	if err != nil {
		return nil, fmt.Errorf("failed to parse git url: %w", err)
	}
	return u, nil
}
