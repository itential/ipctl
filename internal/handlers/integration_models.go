// Copyright 2024 Itential Inc. All Rights Reserved
// Unauthorized copying of this file, via any medium is strictly prohibited
// Proprietary and confidential

package handlers

import (
	"github.com/itential/ipctl/internal/runners"
)

func NewIntegrationModelHandler(r Runtime, desc Descriptors) AssetHandler {
	return NewAssetHandler(
		runners.NewIntegrationModelRunner(r.Client, r.Config),
		desc[integrationModels],
		nil,
	)
}
