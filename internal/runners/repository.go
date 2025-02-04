// Copyright 2024 Itential Inc. All Rights Reserved
// Unauthorized copying of this file, via any medium is strictly prohibited
// Proprietary and confidential

package runners

import (
	"time"

	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing/object"
	"github.com/itential/ipctl/internal/utils"
	"github.com/itential/ipctl/pkg/config"
	"github.com/itential/ipctl/pkg/logger"
	"github.com/itential/ipctl/pkg/repositories"
)

type Repository struct {
	Url            string
	PrivateKeyFile string
	Reference      string
}

func GetRepository(name string, config *config.Config) (*Repository, error) {
	repo, err := config.GetRepository(name)
	if err != nil {
		return nil, err
	}

	return &Repository{
		Url:            repo.Url,
		PrivateKeyFile: repo.PrivateKeyFile,
		Reference:      repo.Reference,
	}, nil
}

func CloneRepository(in *Repository) (string, error) {
	logger.Trace()

	r := repositories.Repository{
		Url: in.Url,
	}
	if in.PrivateKeyFile != "" {
		pk, err := utils.ReadStringFromFile(in.PrivateKeyFile)
		if err != nil {
			return "", err
		}
		r.PrivateKey = pk
	}
	if in.Reference != "" {
		r.Reference = in.Reference
	}
	p, err := r.Clone()
	if err != nil {
		return "", err
	}

	return p, nil
}

func CommitAndPushRepo(in *Repository, path, msg string) error {
	r, err := git.PlainOpen(path)
	if err != nil {
		return err
	}

	w, err := r.Worktree()
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
				Name:  "ipctl",
				Email: "ipctl@localhost",
				When:  time.Now(),
			},
		})
		if err != nil {
			return err
		}

		logger.Info("%v", commit)

		err = r.Push(&git.PushOptions{})
		if err != nil {
			return err
		}
	}

	return nil
}
