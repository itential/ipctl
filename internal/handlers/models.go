// Copyright 2024 Itential Inc. All Rights Reserved
// Unauthorized copying of this file, via any medium is strictly prohibited
// Proprietary and confidential

package handlers

import (
	"github.com/itential/ipctl/internal/flags"
	"github.com/itential/ipctl/internal/runners"
)

func NewModelHandler(rt *Runtime, desc Descriptors) AssetHandler {
	return NewAssetHandler(
		runners.NewModelRunner(rt.GetClient(), rt.GetConfig()),
		desc[modelsDescriptor],
		&AssetHandlerFlags{
			Create: &flags.ModelCreateOptions{},
			Delete: &flags.ModelDeleteOptions{},
			Import: &flags.ModelImportOptions{},
			Export: &flags.ModelExportOptions{},
		},
	)
}
