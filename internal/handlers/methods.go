// Copyright 2024 Itential Inc. All Rights Reserved
// Unauthorized copying of this file, via any medium is strictly prohibited
// Proprietary and confidential

package handlers

import (
	"github.com/itential/ipctl/internal/runners"
)

func NewMethodHandler(r Runtime, desc Descriptors) AssetHandler {
	return NewAssetHandler(
		runners.NewMethodRunner(r.Client, r.Config),
		desc[methodsDescriptor],
		nil,
	)
}
