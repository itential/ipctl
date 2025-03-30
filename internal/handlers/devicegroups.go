// Copyright 2024 Itential Inc. All Rights Reserved
// Unauthorized copying of this file, via any medium is strictly prohibited
// Proprietary and confidential

package handlers

import (
	"github.com/itential/ipctl/internal/flags"
	"github.com/itential/ipctl/internal/runners"
)

func NewDeviceGroupHandler(r Runtime, desc Descriptors) AssetHandler {
	return NewAssetHandler(
		runners.NewDeviceGroupRunner(r.Client, r.Config),
		desc[deviceGroupsDescriptor],
		&AssetHandlerFlags{
			Create: &flags.DeviceGroupCreateOptions{},
		},
	)
}
