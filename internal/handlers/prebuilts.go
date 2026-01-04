// Copyright 2024 Itential Inc. All Rights Reserved
// Unauthorized copying of this file, via any medium is strictly prohibited
// Proprietary and confidential

package handlers

import (
	"github.com/itential/ipctl/internal/flags"
	"github.com/itential/ipctl/internal/runners"
)

func NewPrebuiltHandler(rt *Runtime, desc Descriptors) AssetHandler {
	return NewAssetHandler(
		runners.NewPrebuiltRunner(rt.GetClient(), rt.GetConfig()),
		desc[prebuiltsDescriptor],
		&AssetHandlerFlags{
			Delete: &flags.PrebuiltDeleteOptions{},
			Import: &flags.PrebuiltImportOptions{},
			Export: &flags.PrebuiltExportOptions{},
		},
	)
}
