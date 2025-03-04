// Copyright 2024 Itential Inc. All Rights Reserved
// Unauthorized copying of this file, via any medium is strictly prohibited
// Proprietary and confidential

package flags

import (
	"github.com/spf13/cobra"
)

type CommandTemplateGetOptions struct {
	All bool
}

func (o *CommandTemplateGetOptions) Flags(cmd *cobra.Command) {
	cmd.Flags().BoolVar(&o.All, "all", o.All, "Include all command templates")
}
