// Copyright 2024 Itential Inc. All Rights Reserved
// Unauthorized copying of this file, via any medium is strictly prohibited
// Proprietary and confidential

package handlers

import (
	"github.com/itential/ipctl/internal/flags"
	"github.com/itential/ipctl/internal/runners"
)

func NewProfileHandler(rt *Runtime, desc Descriptors) AssetHandler {
	return NewAssetHandler(
		runners.NewProfileRunner(rt.GetClient(), rt.GetConfig()),
		desc[profilesDescriptor],
		&AssetHandlerFlags{
			Create: &flags.ProfileCreateOptions{},
		},
	)
}
