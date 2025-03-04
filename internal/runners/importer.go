// Copyright 2024 Itential Inc. All Rights Reserved
// Unauthorized copying of this file, via any medium is strictly prohibited
// Proprietary and confidential

package runners

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"

	"github.com/itential/ipctl/internal/flags"
	"github.com/itential/ipctl/internal/utils"
	"github.com/itential/ipctl/pkg/logger"
	"github.com/mitchellh/go-homedir"
)

// This function accepts the request object as the first argument and
// will extract the path, repository information (if provided) and load the
// data from disk into the ptr argument
func importUnmarshalFromRequest(in Request, ptr any) error {
	logger.Trace()

	path, err := importGetPathFromRequest(in)
	if err != nil {
		return err
	}
	defer os.RemoveAll(path)

	return importLoadFromDisk(path, ptr)
}

// importLoadFromDisk will take a path string and pointer to a struct and load
// the data form disk.  This function will read the data from the file provided
// by path and unmarshal the data into the struct pointer.
func importLoadFromDisk(path string, ptr any) error {
	logger.Trace()

	if !utils.PathExists(path) {
		return errors.New(fmt.Sprintf("import path `%s` does not exist", path))
	}

	b, err := os.ReadFile(path)
	if err != nil {
		return err
	}

	utils.UnmarshalData(b, ptr)

	return nil
}

// This function accepts a single required argument which is the incoming
// request object.  It will extrace the value for path and will also check to
// see if the import should come from a Git repository.  If the `--repository`
// argument is specified, this function will clone the Git repository and
// return the full path.
func importGetPathFromRequest(in Request) (string, error) {
	logger.Trace()

	common := in.Common.(*flags.AssetImportCommon)

	path := in.Args[0]

	if common.Repository != "" {
		r := &Repository{
			Url:            common.Repository,
			PrivateKeyFile: common.PrivateKeyFile,
			Reference:      common.Reference,
		}

		p, err := CloneRepository(r)
		if err != nil {
			return "", err
		}

		path = filepath.Join(p, path)
	}

	return homedir.Expand(path)
}
