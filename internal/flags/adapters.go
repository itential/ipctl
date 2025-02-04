// Copyright 2024 Itential Inc. All Rights Reserved
// Unauthorized copying of this file, via any medium is strictly prohibited
// Proprietary and confidential

package flags

import "github.com/spf13/cobra"

type AdapterCreateOptions struct {
	Template  string
	Variables []string
}

func (o *AdapterCreateOptions) Flags(cmd *cobra.Command) {
	cmd.Flags().StringVar(&o.Template, "template", o.Template, "Adapter template configuration")
	cmd.Flags().StringArrayVar(&o.Variables, "set", o.Variables, "One or more template values")
}
