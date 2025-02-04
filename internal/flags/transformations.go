// Copyright 2024 Itential Inc. All Rights Reserved
// Unauthorized copying of this file, via any medium is strictly prohibited
// Proprietary and confidential

package flags

import (
	"github.com/spf13/cobra"
)

type TransformationCreateOptions struct {
	Description string
	Replace     bool
}

func (o *TransformationCreateOptions) Flags(cmd *cobra.Command) {
	cmd.Flags().StringVar(&o.Description, "description", o.Description, "Description of the transformation")
	cmd.Flags().BoolVar(&o.Replace, "replace", o.Replace, "Replace the exist transformation if it exists")
}

type TransformationGetOptions struct {
	All bool
}

func (o *TransformationGetOptions) Flags(cmd *cobra.Command) {
	cmd.Flags().BoolVar(&o.All, "all", o.All, "Include all transformations")
}
