// Copyright 2024 Itential Inc. All Rights Reserved
// Unauthorized copying of this file, via any medium is strictly prohibited
// Proprietary and confidential

package flags

import "github.com/spf13/cobra"

type AdapterCreateOptions struct {
	Model      string
	Template   string
	Variables  []string
	Path       string
	Reference  string
	Repository string
}

func (o *AdapterCreateOptions) Flags(cmd *cobra.Command) {
	cmd.Flags().StringVar(&o.Model, "model", o.Model, "Adapter model")
	cmd.Flags().StringVar(&o.Template, "template", o.Template, "Adapter template configuration")
	cmd.Flags().StringArrayVar(&o.Variables, "set", o.Variables, "One or more template values")
	cmd.Flags().StringVar(&o.Path, "path", o.Path, "Path to the file in the repository")
	cmd.Flags().StringVar(&o.Reference, "reference", o.Reference, "Git reference")
	cmd.Flags().StringVar(&o.Repository, "repository", o.Repository, "Git repository URL")
}
