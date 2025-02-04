// Copyright 2024 Itential Inc. All Rights Reserved
// Unauthorized copying of this file, via any medium is strictly prohibited
// Proprietary and confidential

package runners

import (
	"github.com/itential/ipctl/internal/utils"
	"github.com/itential/ipctl/pkg/logger"
)

func ReadImportFromFile(in Request, ptr any) error {
	logger.Trace()

	path, err := NormalizePath(in)
	if err != nil {
		return err
	}

	return utils.ReadObjectFromDisk(path, ptr)
}

type Asset struct {
	Request Request
}
