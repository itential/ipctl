// Copyright 2024 Itential Inc. All Rights Reserved
// Unauthorized copying of this file, via any medium is strictly prohibited
// Proprietary and confidential

package flags

import "github.com/spf13/cobra"

type GroupCreateOptions struct {
	Description string
}

func (o *GroupCreateOptions) Flags(cmd *cobra.Command) {
	cmd.Flags().StringVar(&o.Description, "description", o.Description, "Short description of the group")
}
