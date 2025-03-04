// Copyright 2024 Itential Inc. All Rights Reserved
// Unauthorized copying of this file, via any medium is strictly prohibited
// Proprietary and confidential

package flags

import (
	"github.com/spf13/cobra"
)

// Command line options for `import project ...`
type ProjectImportOptions struct {
	Members []string
}

func (o *ProjectImportOptions) Flags(cmd *cobra.Command) {
	cmd.Flags().StringArrayVar(&o.Members, "member", o.Members, "Configure one or more project members")
}

// Command line options for `export project ...`
type ProjectExportOptions struct {
	Expand bool
}

func (o *ProjectExportOptions) Flags(cmd *cobra.Command) {
	cmd.Flags().BoolVar(&o.Expand, "expand", o.Expand, "Expand the project assets")
}

// Command line options for `copy project ...`
type ProjectCopyOptions struct {
	Members []string
}

func (o *ProjectCopyOptions) Flags(cmd *cobra.Command) {
	cmd.Flags().StringArrayVar(&o.Members, "member", o.Members, "Configure one or more project members")
}
