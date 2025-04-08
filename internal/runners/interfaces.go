// Copyright 2024 Itential Inc. All Rights Reserved
// Unauthorized copying of this file, via any medium is strictly prohibited
// Proprietary and confidential

package runners

import (
	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing"
)

type Runner interface {
}

type RunnerFunc func(Request) (*Response, error)

type Copier interface {
	Copy(Request) (*Response, error)
	CopyFrom(string, string) (any, error)
	CopyTo(string, any, bool) (any, error)
}

type Reader interface {
	Get(Request) (*Response, error)
	Describe(Request) (*Response, error)
}

type Writer interface {
	Create(Request) (*Response, error)
	Delete(Request) (*Response, error)
	Clear(Request) (*Response, error)
}

type Editor interface {
	Edit(Request) (*Response, error)
}

type Importer interface {
	Import(Request) (*Response, error)
}

type Exporter interface {
	Export(Request) (*Response, error)
}

type Controller interface {
	Start(Request) (*Response, error)
	Stop(Request) (*Response, error)
	Restart(Request) (*Response, error)
}

type Inspector interface {
	Inspect(Request) (*Response, error)
}

type Dumper interface {
	Dump(Request) (*Response, error)
}

type Loader interface {
	Load(Request) (*Response, error)
}

type FileReader interface {
	Read(path string) ([]byte, error)
}

type Cloner interface {
	Clone(r RepositoryPayload) (string, error)
}

type GitRepository interface {
	Worktree() (GitWorktree, error)
	Push(options *PushOptions) error
}

type GitWorktree interface {
	AddGlob(pattern string) error
	Status() (GitStatus, error)
	Commit(msg string, opts *CommitOptions) (Hash, error)
}

type GitStatus interface {
	IsClean() bool
}

type GitProvider interface {
	Open(path string) (GitRepository, error)
}

type PushOptions = git.PushOptions
type CommitOptions = git.CommitOptions
type Hash = plumbing.Hash
