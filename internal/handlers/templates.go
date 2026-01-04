// Copyright 2024 Itential Inc. All Rights Reserved
// Unauthorized copying of this file, via any medium is strictly prohibited
// Proprietary and confidential

package handlers

import (
	"github.com/itential/ipctl/internal/flags"
	"github.com/itential/ipctl/internal/runners"
)

func NewTemplateHandler(rt *Runtime, desc Descriptors) AssetHandler {
	return NewAssetHandler(
		runners.NewTemplateRunner(rt.GetClient(), rt.GetConfig()),
		desc[templatesDescriptor],
		&AssetHandlerFlags{
			Create: &flags.TemplateCreateOptions{},
			Get:    &flags.TemplateGetOptions{},
			Load:   &flags.TemplateLoadOptions{},
		},
	)
}
