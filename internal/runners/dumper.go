// Copyright 2024 Itential Inc. All Rights Reserved
// Unauthorized copying of this file, via any medium is strictly prohibited
// Proprietary and confidential

package runners

import (
	"github.com/itential/ipctl/internal/logging"
)

// dumpAssets takes the Request object and a map of objects and dumps all
// assetst to disk.  The Request object is used to dump either to local disk or
// to a Git repository.  The objects arguments provides a map of assets to dump
// to disk.  The key of the map must be the filename for the asset and the
// value must be the object instance.
func dumpAssets(in Request, objects map[string]interface{}) error {
	logging.Trace()
	return exportAssets(in, objects)
}
