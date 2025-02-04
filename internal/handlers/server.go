// Copyright 2024 Itential Inc. All Rights Reserved
// Unauthorized copying of this file, via any medium is strictly prohibited
// Proprietary and confidential

package handlers

import (
	"github.com/itential/ipctl/internal/runners"
	"github.com/itential/ipctl/pkg/client"
	"github.com/itential/ipctl/pkg/config"
)

func NewServerHandler(c client.Client, cfg *config.Config, desc Descriptors) AssetHandler {
	return NewAssetHandler(
		runners.NewServerRunner(c, cfg),
		desc[serverDescriptor],
		nil,
	)
}
