// Copyright 2024 Itential Inc. All Rights Reserved
// Unauthorized copying of this file, via any medium is strictly prohibited
// Proprietary and confidential

package flags

import (
	"github.com/spf13/cobra"
)

type AutomationCreateOptions struct {
	Description string
	Replace     bool
}

func (o *AutomationCreateOptions) Flags(cmd *cobra.Command) {
	cmd.Flags().StringVar(&o.Description, "description", o.Description, "Description of the automation")
	cmd.Flags().BoolVar(&o.Replace, "replace", o.Replace, "Replace the exist automation if it exists")
}
