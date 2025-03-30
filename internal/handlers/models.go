// Copyright 2024 Itential Inc. All Rights Reserved
// Unauthorized copying of this file, via any medium is strictly prohibited
// Proprietary and confidential

package handlers

import (
	"github.com/itential/ipctl/internal/flags"
	"github.com/itential/ipctl/internal/runners"
)

func NewModelHandler(r Runtime, desc Descriptors) AssetHandler {
	return NewAssetHandler(
		runners.NewModelRunner(r.Client, r.Config),
		desc[modelsDescriptor],
		&AssetHandlerFlags{
			Create: &flags.ModelCreateOptions{},
			Delete: &flags.ModelDeleteOptions{},
			Import: &flags.ModelImportOptions{},
			Export: &flags.ModelExportOptions{},
		},
	)
}
