// Copyright 2024 Itential Inc. All Rights Reserved
// Unauthorized copying of this file, via any medium is strictly prohibited
// Proprietary and confidential

package runners

import (
	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing"
)

// GitProviderImpl is the real implementation that uses go-git
type GitProviderImpl struct{}

func (p *GitProviderImpl) Open(path string) (GitRepository, error) {
	repo, err := git.PlainOpen(path)
	if err != nil {
		return nil, err
	}
	return &GitRepositoryImpl{repo}, nil
}

// GitRepositoryImpl wraps go-git's Repository
type GitRepositoryImpl struct {
	repo *git.Repository
}

func (r *GitRepositoryImpl) Worktree() (GitWorktree, error) {
	w, err := r.repo.Worktree()
	if err != nil {
		return nil, err
	}
	return &GitWorktreeImpl{w}, nil
}

func (r *GitRepositoryImpl) Push(opts *PushOptions) error {
	return r.repo.Push(opts)
}

// GitWorktreeImpl wraps go-git's Worktree
type GitWorktreeImpl struct {
	w *git.Worktree
}

func (w *GitWorktreeImpl) AddGlob(pattern string) error {
	return w.w.AddGlob(pattern)
}

func (w *GitWorktreeImpl) Status() (GitStatus, error) {
	status, err := w.w.Status()
	if err != nil {
		return nil, err
	}
	return &GitStatusImpl{status}, nil
}

func (w *GitWorktreeImpl) Commit(msg string, opts *CommitOptions) (Hash, error) {
	commitHash, err := w.w.Commit(msg, opts)
	return plumbing.Hash(commitHash), err
}

// GitStatusImpl wraps go-git's Status object
type GitStatusImpl struct {
	status git.Status
}

func (s *GitStatusImpl) IsClean() bool {
	return s.status.IsClean()
}
