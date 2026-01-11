// Copyright 2024 Itential Inc. All Rights Reserved
// Unauthorized copying of this file, via any medium is strictly prohibited
// Proprietary and confidential

package runners

import (
	"fmt"
	"os/user"
	"time"

	"github.com/go-git/go-git/v5/plumbing/object"
	"github.com/itential/ipctl/internal/logging"
	"github.com/itential/ipctl/internal/utils"
	"github.com/itential/ipctl/pkg/repositories"
)

// FileReaderImpl is a concrete file reader using utils.ReadFromFile
type FileReaderImpl struct{}

func (f *FileReaderImpl) Read(path string) ([]byte, error) {
	return utils.ReadFromFile(path)
}

// ClonerImpl is a real cloner using repositories.Repository
type ClonerImpl struct{}

func (c *ClonerImpl) Clone(p RepositoryPayload) (string, error) {
	repo := repositories.Repository{
		Url:        p.Url,
		User:       p.User,
		Reference:  p.Reference,
		PrivateKey: p.PrivateKey,
	}
	return repo.Clone()
}

// RepositoryOption provides options for configuring an instance of Repository
type RepositoryOption func(r *Repository)

type Repository struct {
	Url            string
	PrivateKeyFile string
	Reference      string
	Name           string
	Email          string
}

// RepositoryPayload is a simplified struct used to pass data into the cloner
// interface
type RepositoryPayload struct {
	Url        string
	User       string
	Reference  string
	PrivateKey []byte
}

// userProvider is a function that returns the current OS user
type userProvider func() (*user.User, error)

// newRepository is the internal constructor that accepts a user provider (used for testing)
func newRepository(url string, getUser userProvider, opts ...RepositoryOption) *Repository {
	logging.Trace()

	currentUser, err := getUser()
	if err != nil {
		logging.Fatal(err, "failed to get current user")
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

// NewRepository will create a new instance of a Repository argument.  The
// required argument `url` defines the URL for the reposiotry.  This function
// will accept any valid URL format.  This function will also accept one or
// more options for configuring the repository.
//
// The following options are supported:
//   - WithReference
//   - WithPrivateKeyFile
//   - WithName
//   - WithEmail
//
// See the optional function for details about each implemenation.  If an
// option is not passed, a default value is set for repository object.
func NewRepository(url string, opts ...RepositoryOption) *Repository {
	logging.Trace()
	return newRepository(url, user.Current, opts...)
}

// WithReference sets the Git reference (branch, tag, etc.)
func WithReference(v string) RepositoryOption {
	return func(r *Repository) {
		r.Reference = v
	}
}

// WithPrivateKeyFile sets the path to the SSH private key file
func WithPrivateKeyFile(v string) RepositoryOption {
	return func(r *Repository) {
		r.PrivateKeyFile = v
	}
}

// WithName overrides the Git username if provided
func WithName(v string) RepositoryOption {
	return func(r *Repository) {
		if v != "" {
			r.Name = v
		}
	}
}

// WithEmail overrides the Git email if provided
func WithEmail(v string) RepositoryOption {
	return func(r *Repository) {
		if v != "" {
			r.Email = v
		}
	}
}

func (r *Repository) Clone(reader FileReader, cloner Cloner) (string, error) {
	logging.Trace()

	payload := RepositoryPayload{
		Url:  r.Url,
		User: "git",
	}

	if r.PrivateKeyFile != "" {
		key, err := reader.Read(r.PrivateKeyFile)
		if err != nil {
			return "", err
		}
		payload.PrivateKey = key
	}

	if r.Reference != "" {
		payload.Reference = r.Reference
	}

	return cloner.Clone(payload)
}

func (r *Repository) CommitAndPush(path, msg string) error {
	return r.commitAndPush(path, msg, &GitProviderImpl{})
}

/*
*******************************************************************************
Private functions
*******************************************************************************
*/

// commitAndPush handles adding the working tree to the repository, commiting
// the working tree and pushing it to the remote repository.  The `path`
// argument is the path to the repo that contains the clond.  The `msg`
// argument is the commit message.  The `provider` argument specifies the
// provider to use.
func (r *Repository) commitAndPush(path, msg string, provider GitProvider) error {
	logging.Trace()

	repo, err := provider.Open(path)
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
		commit, err := w.Commit(msg, &CommitOptions{
			Author: &object.Signature{
				Name:  r.Name,
				Email: r.Email,
				When:  time.Now(),
			},
		})
		if err != nil {
			return err
		}

		logging.Info("%v", commit)

		if err := repo.Push(&PushOptions{}); err != nil {
			return err
		}
	}

	return nil
}
