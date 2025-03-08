// Copyright 2024 Itential Inc. All Rights Reserved
// Unauthorized copying of this file, via any medium is strictly prohibited
// Proprietary and confidential

package runners

import (
	"os"
	"path/filepath"

	"github.com/itential/ipctl/internal/flags"
	"github.com/itential/ipctl/internal/utils"
	"github.com/itential/ipctl/pkg/logger"
)

// exportNewRepository will create a new Repository object from an export
// request.
func exportNewRepositoryFromRequest(in Request) *Repository {
	common := in.Common.(*flags.AssetExportCommon)
	return NewRepository(
		common.Repository,
		WithReference(common.Reference),
		WithPrivateKeyFile(common.PrivateKeyFile),
		WithName(in.Config.GitName),
		WithEmail(in.Config.GitEmail),
	)
}

// exportAssetFromRequest will take a request object and instance of an asset
// and write it to disk.  If the Git command line options where invoked, it
// will write the asset to the repository and commit it.  If not, this function
// will simply write the asset to the local disk.
func exportAssetFromRequest(in Request, o any, fn string) error {
	logger.Trace()

	path := in.Common.(flags.Committer).GetPath()

	var repo *Repository
	var repoPath string

	if in.Common.(flags.Gitter).GetRepository() != "" {
		repo = NewRepository(
			in.Common.(flags.Gitter).GetRepository(),
			WithReference(in.Common.(flags.Gitter).GetReference()),
			WithPrivateKeyFile(in.Common.(flags.Gitter).GetPrivateKeyFile()),
			WithName(in.Config.GitName),
			WithEmail(in.Config.GitEmail),
		)

		var e error

		repoPath, e = repo.Clone()
		if e != nil {
			return e
		}
		defer os.RemoveAll(repoPath)

		path = filepath.Join(repoPath, in.Common.(flags.Committer).GetPath())
	}

	if err := utils.WriteJsonToDisk(o, fn, path); err != nil {
		return err
	}

	if in.Common.(flags.Gitter).GetRepository() != "" {
		msg := in.Common.(flags.Committer).GetMessage()
		if err := repo.CommitAndPush(repoPath, msg); err != nil {
			return err
		}
	}

	return nil
}
