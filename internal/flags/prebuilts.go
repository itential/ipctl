// Copyright 2024 Itential Inc. All Rights Reserved
// Unauthorized copying of this file, via any medium is strictly prohibited
// Proprietary and confidential

package flags

import (
	"github.com/spf13/cobra"
)

type PrebuiltDeleteOptions struct {
	All bool
}

func (o *PrebuiltDeleteOptions) Flags(cmd *cobra.Command) {
	cmd.Flags().BoolVar(&o.All, "all", o.All, "Delete all prebuilt assets")
}

type PrebuiltImportOptions struct {
	Path string
}

func (o *PrebuiltImportOptions) Flags(cmd *cobra.Command) {
	cmd.Flags().StringVar(&o.Path, "path", o.Path, "Path to Prebuilt")
}

type PrebuiltExportOptions struct {
	Expand bool
}

func (o *PrebuiltExportOptions) Flags(cmd *cobra.Command) {
	cmd.Flags().BoolVar(&o.Expand, "expand", o.Expand, "Expand the project assets")
}
