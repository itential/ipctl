// Copyright 2024 Itential Inc. All Rights Reserved
// Unauthorized copying of this file, via any medium is strictly prohibited
// Proprietary and confidential

package handlers

import (
	"github.com/itential/ipctl/internal/flags"
	"github.com/itential/ipctl/internal/runners"
	"github.com/itential/ipctl/pkg/client"
	"github.com/itential/ipctl/pkg/config"
)

func NewRoleHandler(c client.Client, cfg *config.Config, desc Descriptors) AssetHandler {
	return NewAssetHandler(
		runners.NewRoleRunner(c, cfg),
		desc[rolesDescriptor],
		&AssetHandlerFlags{
			Get:      &flags.RoleGetOptions{},
			Describe: &flags.RoleDescribeOptions{},
			Create:   &flags.RoleCreateOptions{},
		},
	)
}
