// Copyright 2024 Itential Inc. All Rights Reserved
// Unauthorized copying of this file, via any medium is strictly prohibited
// Proprietary and confidential

package runners

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/itential/ipctl/internal/flags"
	"github.com/itential/ipctl/internal/logging"
	"github.com/itential/ipctl/internal/utils"
	giturls "github.com/whilp/git-urls"
)

// validateRepositoryURL validates that a repository URL is well-formed and uses a supported scheme.
func validateRepositoryURL(url string) error {
	if url == "" {
		return fmt.Errorf("repository URL cannot be empty")
	}

	validSchemes := []string{"file", "git", "https", "ssh", "git+ssh"}

	u, err := giturls.Parse(url)
	if err != nil {
		return fmt.Errorf("invalid repository URL %q: %w", url, err)
	}

	schemeValid := false
	for _, scheme := range validSchemes {
		if u.Scheme == scheme || strings.HasPrefix(u.Scheme, scheme) {
			schemeValid = true
			break
		}
	}

	if !schemeValid {
		return fmt.Errorf("unsupported URL scheme %q (supported: %v)", u.Scheme, validSchemes)
	}

	return nil
}

// exportNewRepository will create a new Repository object from an the incoming
// Request object.
func exportNewRepositoryFromRequest(in Request) (*Repository, error) {
	logging.Trace()

	common := in.Common.(flags.Gitter)

	url := common.GetRepository()
	privateKeyFile := common.GetPrivateKeyFile()
	reference := common.GetReference()

	if err := validateRepositoryURL(url); err != nil {
		return nil, err
	}

	u, err := giturls.Parse(url)
	if err != nil {
		return nil, fmt.Errorf("failed to parse repository URL %q: %w", url, err)
	}

	if u.Scheme == "file" && strings.HasPrefix(u.Path, "@") {
		r, err := in.Config.GetRepository(u.Path[1:])
		if err != nil {
			return nil, err
		}

		url = r.Url

		if privateKeyFile == "" {
			privateKeyFile = r.PrivateKeyFile
		}

		if reference == "" {
			reference = r.Reference
		}
	}

	return NewRepository(
		url,
		WithReference(reference),
		WithPrivateKeyFile(privateKeyFile),
		WithName(in.Config.GitName),
		WithEmail(in.Config.GitEmail),
	), nil
}

// exportAssetFromRequest will take a request object and instance of an asset
// and write it to disk.  If the Git command line options where invoked, it
// will write the asset to the repository and commit it.  If not, this function
// will simply write the asset to the local disk.
func exportAssetFromRequest(in Request, o any, fn string) error {
	return exportAssets(in, map[string]interface{}{fn: o})
}

// exportAssets accepts the Request object and a map of the assets and will
// write them to disk.  If the request object includes repository settings,
// this function will push the assets into the repository.  The assets argument
// must be a map where the key is the filename and the value is the asset to
// write to disk.
func exportAssets(in Request, assets map[string]interface{}) error {
	logging.Trace()

	path := in.Common.(flags.Committer).GetPath()

	var repo *Repository
	var repoPath string

	if in.Common.(flags.Gitter).GetRepository() != "" {
		var e error

		repo, e = exportNewRepositoryFromRequest(in)
		if e != nil {
			return e
		}

		repoPath, e = repo.Clone(
			&FileReaderImpl{},
			&ClonerImpl{},
		)
		if e != nil {
			return e
		}
		defer os.RemoveAll(repoPath)

		path = filepath.Join(repoPath, in.Common.(flags.Committer).GetPath())
	}

	for key, value := range assets {
		if err := utils.WriteJsonToDisk(value, key, path); err != nil {
			return err
		}
	}

	if repo != nil {
		msg := in.Common.(flags.Committer).GetMessage()
		if err := repo.CommitAndPush(repoPath, msg); err != nil {
			return err
		}
	}

	return nil
}
