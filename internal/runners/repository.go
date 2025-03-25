// Copyright 2024 Itential Inc. All Rights Reserved
// Unauthorized copying of this file, via any medium is strictly prohibited
// Proprietary and confidential

package runners

import (
	"fmt"
	"os/user"
	"time"

	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing/object"
	"github.com/itential/ipctl/internal/utils"
	"github.com/itential/ipctl/pkg/logger"
	"github.com/itential/ipctl/pkg/repositories"
)

type RepositoryOption func(r *Repository)

type Repository struct {
	Url            string
	PrivateKeyFile string
	Reference      string
	Name           string
	Email          string
}

func NewRepository(url string, opts ...RepositoryOption) *Repository {
	logger.Trace()

	currentUser, err := user.Current()
	if err != nil {
		logger.Fatal(err, "failed get current user")
	}

	r := &Repository{
		Url:   url,
		Name:  currentUser.Username,
		Email: fmt.Sprintf("%s@users.ipctl", currentUser.Username),
	}

	for _, opt := range opts {
		opt(r)
	}

	return r
}

func WithReference(v string) RepositoryOption {
	return func(r *Repository) {
		r.Reference = v
	}
}

func WithPrivateKeyFile(v string) RepositoryOption {
	return func(r *Repository) {
		r.PrivateKeyFile = v
	}
}

func WithName(v string) RepositoryOption {
	return func(r *Repository) {
		if v != "" {
			r.Name = v
		}
	}
}

func WithEmail(v string) RepositoryOption {
	return func(r *Repository) {
		if v != "" {
			r.Name = v
		}
	}
}

func (r Repository) Clone() (string, error) {
	logger.Trace()

	repo := repositories.Repository{
		Url:  r.Url,
		User: "git",
	}
	if r.PrivateKeyFile != "" {
		pk, err := utils.ReadFromFile(r.PrivateKeyFile)
		if err != nil {
			return "", err
		}
		repo.PrivateKey = pk
	}
	if r.Reference != "" {
		repo.Reference = r.Reference
	}
	p, err := repo.Clone()
	if err != nil {
		return "", err
	}

	return p, nil
}

func (r Repository) CommitAndPush(path, msg string) error {
	logger.Trace()

	repo, err := git.PlainOpen(path)
	if err != nil {
		return err
	}

	w, err := repo.Worktree()
	if err != nil {
		return err
	}

	if err = w.AddGlob("*"); err != nil {
		return err
	}

	status, err := w.Status()
	if err != nil {
		return err
	}

	if !status.IsClean() {
		commit, err := w.Commit(msg, &git.CommitOptions{
			Author: &object.Signature{
				Name:  r.Name,
				Email: r.Email,
				When:  time.Now(),
			},
		})
		if err != nil {
			return err
		}

		logger.Info("%v", commit)

		if err := repo.Push(&git.PushOptions{}); err != nil {
			return err
		}
	}

	return nil
}

func CloneRepository(in *Repository) (string, error) {
	return in.Clone()
}

func CommitAndPushRepo(in *Repository, path, msg string) error {
	return in.CommitAndPush(path, msg)
}
