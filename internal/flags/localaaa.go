// Copyright 2024 Itential Inc. All Rights Reserved
// Unauthorized copying of this file, via any medium is strictly prohibited
// Proprietary and confidential

package flags

import (
	"github.com/spf13/cobra"
)

type LocalAAAOptions struct {
	Groups []string
}

func (o *LocalAAAOptions) Flags(cmd *cobra.Command) {
	cmd.Flags().StringArrayVar(&o.Groups, "group", o.Groups, "Group to insert user into")
}
