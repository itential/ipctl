// Copyright 2024 Itential Inc. All Rights Reserved
// Unauthorized copying of this file, via any medium is strictly prohibited
// Proprietary and confidential

package flags

import (
	"github.com/spf13/cobra"
)

type ModelCreateOptions struct {
	Description string
	Schema      string
	Replace     bool
}

func (o *ModelCreateOptions) Flags(cmd *cobra.Command) {
	cmd.Flags().StringVar(&o.Description, "description", o.Description, "Short description of model")
	cmd.Flags().StringVar(&o.Schema, "schema", o.Schema, "JSON Schema of the resource model")
	cmd.Flags().BoolVar(&o.Replace, "replace", o.Replace, "Overwrite existing model if it exists")
}
