// Copyright 2024 Itential Inc. All Rights Reserved
// Unauthorized copying of this file, via any medium is strictly prohibited
// Proprietary and confidential

package handlers

import (
	"github.com/itential/ipctl/internal/flags"
	"github.com/itential/ipctl/internal/runners"
)

func NewProjectHandler(r Runtime, desc Descriptors) AssetHandler {
	return NewAssetHandler(
		runners.NewProjectRunner(r.Client, r.Config),
		desc[projectsDescriptor],
		&AssetHandlerFlags{
			Import: &flags.ProjectImportOptions{},
			Export: &flags.ProjectExportOptions{},
			Copy:   &flags.ProjectCopyOptions{},
		},
	)
}
