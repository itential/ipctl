// Copyright 2024 Itential Inc. All Rights Reserved
// Unauthorized copying of this file, via any medium is strictly prohibited
// Proprietary and confidential

package flags

import (
	"github.com/spf13/cobra"
)

type ProjectImportOptions struct {
	Members []string
}

func (o *ProjectImportOptions) Flags(cmd *cobra.Command) {
	cmd.Flags().StringArrayVar(&o.Members, "member", o.Members, "Configure one or more project members")
}

type ProjectExportOptions struct {
	Expand bool
}

func (o *ProjectExportOptions) Flags(cmd *cobra.Command) {
	cmd.Flags().BoolVar(&o.Expand, "expand", o.Expand, "Expand the project assets")
}
