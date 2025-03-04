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

func NewProjectHandler(iapClient client.Client, cfg *config.Config, desc Descriptors) AssetHandler {
	return NewAssetHandler(
		runners.NewProjectRunner(iapClient, cfg),
		desc[projectsDescriptor],
		&AssetHandlerFlags{
			Import: &flags.ProjectImportOptions{},
			Export: &flags.ProjectExportOptions{},
			Copy:   &flags.ProjectCopyOptions{},
		},
	)
}
