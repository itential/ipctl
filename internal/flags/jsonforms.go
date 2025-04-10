// Copyright 2024 Itential Inc. All Rights Reserved
// Unauthorized copying of this file, via any medium is strictly prohibited
// Proprietary and confidential

package flags

import (
	"github.com/spf13/cobra"
)

type JsonFormCreateOptions struct {
	Description string
	Replace     bool
}

func (o *JsonFormCreateOptions) Flags(cmd *cobra.Command) {
	cmd.Flags().StringVar(&o.Description, "description", o.Description, "Description of JSON form")
	cmd.Flags().BoolVar(&o.Replace, "replace", o.Replace, "Replace the existing form if it exists")
}

type JsonFormGetOptions struct {
	All bool
}

func (o *JsonFormGetOptions) Flags(cmd *cobra.Command) {
	cmd.Flags().BoolVar(&o.All, "all", o.All, "Include all JSON Forms")
}
