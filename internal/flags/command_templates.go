// Copyright 2024 Itential Inc. All Rights Reserved
// Unauthorized copying of this file, via any medium is strictly prohibited
// Proprietary and confidential

package flags

import (
	"github.com/spf13/cobra"
)

type CommandTemplateCreateOptions struct {
	Replace bool
}

func (o *CommandTemplateCreateOptions) Flags(cmd *cobra.Command) {
	cmd.Flags().BoolVar(&o.Replace, "replace", o.Replace, "Replace the exist command template if it exists")
}

type CommandTemplateGetOptions struct {
	All bool
}

func (o *CommandTemplateGetOptions) Flags(cmd *cobra.Command) {
	cmd.Flags().BoolVar(&o.All, "all", o.All, "Include all command templates")
}
