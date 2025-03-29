// Copyright 2024 Itential Inc. All Rights Reserved
// Unauthorized copying of this file, via any medium is strictly prohibited
// Proprietary and confidential

package flags

import (
	"github.com/spf13/cobra"
)

type ModelCreateOptions struct {
	Description string
	Schema      string
	Replace     bool
}

func (o *ModelCreateOptions) Flags(cmd *cobra.Command) {
	cmd.Flags().StringVar(&o.Description, "description", o.Description, "Short description of model")
	cmd.Flags().StringVar(&o.Schema, "schema", o.Schema, "JSON Schema of the resource model")
	cmd.Flags().BoolVar(&o.Replace, "replace", o.Replace, "Overwrite existing model if it exists")
}

type ModelDeleteOptions struct {
	All             bool
	DeleteInstances bool
}

func (o *ModelDeleteOptions) Flags(cmd *cobra.Command) {
	cmd.Flags().BoolVar(&o.All, "all", o.All, "Include all action assets")
	cmd.Flags().BoolVar(&o.DeleteInstances, "delete-instances", o.DeleteInstances, "Delete all instances for this model")
}

type ModelExportOptions struct {
	Expand bool
}

func (o *ModelExportOptions) Flags(cmd *cobra.Command) {
	cmd.Flags().BoolVar(&o.Expand, "expand", o.Expand, "Expand export to include all assets")
}

type ModelImportOptions struct {
	All        bool
	SkipChecks bool
}

func (o *ModelImportOptions) Flags(cmd *cobra.Command) {
	cmd.Flags().BoolVar(&o.All, "all", o.All, "Include all action assets")
	cmd.Flags().BoolVar(&o.SkipChecks, "skip-checks", o.SkipChecks, "Skip checking for other assets")
}
