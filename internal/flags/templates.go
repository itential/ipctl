// Copyright 2024 Itential Inc. All Rights Reserved
// Unauthorized copying of this file, via any medium is strictly prohibited
// Proprietary and confidential

package flags

import (
	"github.com/spf13/cobra"
)

type TemplateCreateOptions struct {
	Description string
	Group       string
	Type        string
	Replace     bool
}

func (o *TemplateCreateOptions) Flags(cmd *cobra.Command) {
	cmd.Flags().StringVar(&o.Description, "description", o.Description, "Description of JSON form")
	cmd.Flags().BoolVar(&o.Replace, "replace", o.Replace, "Replace the exist form if it exists")
	cmd.Flags().StringVar(&o.Group, "group", o.Group, "Group name (REQUIRED)")
	cmd.MarkFlagRequired("group")
	cmd.Flags().StringVar(&o.Type, "type", o.Type, "Type of template to create (REQUIRED)")
	cmd.MarkFlagRequired("type")
}

type TemplateGetOptions struct {
	All bool
}

func (o *TemplateGetOptions) Flags(cmd *cobra.Command) {
	cmd.Flags().BoolVar(&o.All, "all", o.All, "Include all workflows")
}
